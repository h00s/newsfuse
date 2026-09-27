package services

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/go-raptor/raptor/v4"
	"github.com/h00s/newsfuse/app/models"
	"github.com/uptrace/bun"
)

const (
	headlinesPageSize = 30
	// The first page of each topic is cached and invalidated on ingest; the TTL only bounds a
	// race between a read that loaded before an ingest and stored after it.
	headlinesTTL = 5 * time.Minute
)

// inTopic scopes headlines to a topic through their source.
const inTopic = `headlines.source_id IN (SELECT id FROM sources WHERE topic_id = ?)`

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

type HeadlinesService struct {
	raptor.Service

	DB      *DatabaseService
	Cache   *CacheService
	Topics  *TopicsService
	Sources *SourcesService
}

func topicHeadlinesKey(topicID int64) string {
	return fmt.Sprintf("headlines:topic:%d", topicID)
}

// List pages through a topic's headlines, newest first. sourceID narrows the list to one of the
// topic's sources; beforeID continues after the last headline of the previous page (0 for the
// first page).
func (s *HeadlinesService) List(topicID int64, sourceID *int64, beforeID int64) (models.Headlines, error) {
	if err := s.Topics.Verify(topicID); err != nil {
		return nil, err
	}
	if sourceID != nil {
		if err := s.Sources.VerifyInTopic(*sourceID, topicID); err != nil {
			return nil, err
		}
	}

	load := func() (models.Headlines, error) {
		return s.page(beforeID, func(q *bun.SelectQuery) *bun.SelectQuery {
			q = q.Where(inTopic, topicID)
			if sourceID != nil {
				q = q.Where("headlines.source_id = ?", *sourceID)
			}
			return q
		})
	}
	if sourceID == nil && beforeID == 0 {
		return cached(s.Cache, topicHeadlinesKey(topicID), headlinesTTL, load)
	}
	return load()
}

// Search pages through the headlines whose title contains query, newest first. LIKE wildcards
// in the query match themselves.
func (s *HeadlinesService) Search(query string, beforeID int64) (models.Headlines, error) {
	pattern := "%" + likeEscaper.Replace(query) + "%"
	return s.page(beforeID, func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where(`headlines.title ILIKE ? ESCAPE '\'`, pattern)
	})
}

func (s *HeadlinesService) page(beforeID int64, filter func(*bun.SelectQuery) *bun.SelectQuery) (models.Headlines, error) {
	var headlines models.Headlines
	q := s.DB.Conn().NewSelect().
		Model(&headlines).
		Relation("Source")
	q = filter(q)
	if beforeID > 0 {
		q = q.Where("headlines.id < ?", beforeID)
	}
	err := q.Order("headlines.id DESC").
		Limit(headlinesPageSize).
		Scan(s.DB.Ctx)
	return headlines, s.DB.HandleError(err)
}

// CountSince counts a topic's headlines published after since.
func (s *HeadlinesService) CountSince(topicID int64, since time.Time) (int, error) {
	if err := s.Topics.Verify(topicID); err != nil {
		return 0, err
	}
	count, err := s.DB.Conn().NewSelect().
		Model((*models.Headline)(nil)).
		Where(inTopic, topicID).
		Where("headlines.published_at > ?", since).
		Count(s.DB.Ctx)
	return count, s.DB.HandleError(err)
}

func (s *HeadlinesService) Get(id int64) (*models.Headline, error) {
	headline := new(models.Headline)
	err := s.DB.Conn().NewSelect().
		Model(headline).
		Relation("Source").
		Where("headlines.id = ?", id).
		Scan(s.DB.Ctx)
	return headline, s.DB.HandleErrorNotFound(err, "Headline not found")
}

// Ingest stores a scraped batch, skipping headlines whose url is already known, and returns how
// many were new. Scrapers list the newest first; storing the batch reversed gives newer
// headlines higher ids, which is the order every list uses.
func (s *HeadlinesService) Ingest(batch models.Headlines) (int, error) {
	if len(batch) == 0 {
		return 0, nil
	}
	rows := slices.Clone(batch)
	slices.Reverse(rows)

	res, err := s.DB.Conn().NewInsert().
		Model(&rows).
		On("CONFLICT (url) DO NOTHING").
		Exec(s.DB.Ctx)
	if err != nil {
		return 0, s.DB.HandleError(err)
	}
	inserted, err := res.RowsAffected()
	if err != nil {
		return 0, s.DB.HandleError(err)
	}
	if inserted > 0 {
		s.invalidateTopicsOf(rows)
	}
	return int(inserted), nil
}

func (s *HeadlinesService) invalidateTopicsOf(headlines models.Headlines) {
	seen := map[int64]bool{}
	for _, h := range headlines {
		if seen[h.SourceID] {
			continue
		}
		seen[h.SourceID] = true
		source, err := s.Sources.Get(h.SourceID)
		if err != nil {
			s.Log.Warn("Ingested headlines of an unknown source", "source", h.SourceID)
			continue
		}
		s.Cache.Invalidate(topicHeadlinesKey(source.TopicID))
	}
}
