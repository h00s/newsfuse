package scrapers

import "github.com/h00s/newsfuse/api/mmc/models"

type RadioDaruvar struct {
	ID      int    `json:"id"`
	BaseURL string `json:"base_url"`
}

func NewRadioDaruvar() *RadioDaruvar {
	return &RadioDaruvar{
		ID:      1,
		BaseURL: "http://www.radio-daruvar.hr/",
	}
}

func (r *RadioDaruvar) GetHeadlines() ([]models.Headline, error) {
	return nil, nil
}
