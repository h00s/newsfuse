package scrapers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/h00s/newsfuse/app/utils/scrapers"
)

func servePage(t *testing.T, body string) string {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	return server.URL
}

// Markup from the November 2026 redesign of mojportal.hr, trimmed.
func TestMojportalhrHeadlines(t *testing.T) {
	url := servePage(t, `<html><body>
		<article class="se-card se-card--news se-card--story se-card--web-obicni se-card--kolumne">
			<a class="se-card__wrap" href="/kolumne/toni-i-1-b/">
				<div class="se-card__content">
					<p class="se-card__meta"><span class="se-card__tag-section">Kolumne</span></p>
					<h3 class="se-card__head">Toni i 1. B pokazali roditeljima što znači sjećanje na Vukovar</h3>
					<p class="se-card__lead">Toni, koji je u dva mjeseca...</p>
				</div>
			</a>
		</article>
		<article class="se-card se-card--news se-card--event">
			<a class="se-card__wrap" href="/kalendar-dogadanja/pakrac/igrani-film-260-dana/">
				<h3 class="se-card__head">Igrani film 260 dana</h3>
			</a>
		</article>
		<article class="se-card se-card--classified se-card--oglasi">
			<a class="se-card__wrap" href="/oglasi/1/"><h3 class="se-card__head">Prodajem stan</h3></a>
		</article>
		<article class="se-card se-card--news se-card--column">
			<a class="se-card__wrap" href="/crna-kronika/zena-iz-bjelovara/">
				<h2 class="se-card__head"> Žena iz Bjelovara ostala bez 61.000 eura </h2>
			</a>
		</article>
	</body></html>`)
	s := scrapers.NewMojportalhr(7)
	s.URL = url

	got, err := s.Scrape(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	want := []struct{ title, url string }{
		{"Toni i 1. B pokazali roditeljima što znači sjećanje na Vukovar", "https://mojportal.hr/kolumne/toni-i-1-b/"},
		{"Žena iz Bjelovara ostala bez 61.000 eura", "https://mojportal.hr/crna-kronika/zena-iz-bjelovara/"},
	}
	if len(got) != len(want) {
		t.Fatalf("Scrape = %+v, want the %d news cards only", got, len(want))
	}
	for i, w := range want {
		if got[i].Title != w.title || got[i].URL != w.url || got[i].SourceID != 7 {
			t.Errorf("headline %d = %q %q (source %d), want %q %q (source 7)", i, got[i].Title, got[i].URL, got[i].SourceID, w.title, w.url)
		}
	}
}

func TestMojportalhrStory(t *testing.T) {
	url := servePage(t, `<html><body><div class="se-article__text">
		<p>60-godišnjakinja s bjelovarskog područja prijavila je policiji prijevaru.</p>
		<h2>Kako funkcioniraju investicijske prijevare?</h2>
		<p>Investicijske prijevare najčešće započinju oglasima.</p>
		<p class="se-caption">Foto: MojPortal.hr</p>
		<p><em><strong>Dobijte informaciju odmah, zapratite nas na <a href="https://www.facebook.com/MojPortal.hr">Facebooku </a>i <a href="https://www.tiktok.com/@mojportal.hr">TikToku</a>!</strong></em></p>
	</div></body></html>`)

	got, err := scrapers.NewMojportalhr(7).ScrapeStory(t.Context(), url)
	if err != nil {
		t.Fatal(err)
	}
	want := "<p>60-godišnjakinja s bjelovarskog područja prijavila je policiji prijevaru.</p><p>Investicijske prijevare najčešće započinju oglasima.</p>"
	if got != want {
		t.Errorf("ScrapeStory = %q, want %q", got, want)
	}
}
