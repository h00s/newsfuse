package collector_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/h00s/newsfuse/internal/collector"
	"golang.org/x/text/encoding/charmap"
)

func serve(t *testing.T, handler http.HandlerFunc) string {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return server.URL
}

func serveHTML(t *testing.T, body string) string {
	return serve(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(body))
	})
}

func TestElementHelpers(t *testing.T) {
	url := serveHTML(t, `<html><body>
		<div class="card" id="c1" data-x=" raw ">
			<h2> Naslov <span>prvi</span> </h2>
			<a href=" /a ">A</a><a href="/b">B</a>
			<p>jedan</p><p>dva</p>
		</div>
	</body></html>`)

	var got *collector.HTMLElement
	var paragraphs []string
	var indexes []int
	c := collector.New()
	c.OnHTML("div.card", func(e *collector.HTMLElement) {
		got = e
		e.ForEach("p", func(i int, el *collector.HTMLElement) {
			indexes = append(indexes, i)
			paragraphs = append(paragraphs, el.Text)
		})
	})
	if err := c.Visit(t.Context(), url); err != nil {
		t.Fatal(err)
	}

	if got == nil {
		t.Fatal("callback did not run")
	}
	if got.Name != "div" {
		t.Errorf("Name = %q, want div", got.Name)
	}
	if got.Attr("id") != "c1" || got.Attr("data-x") != " raw " || got.Attr("missing") != "" {
		t.Errorf("Attr = %q, %q, %q; want c1, untrimmed \" raw \", empty", got.Attr("id"), got.Attr("data-x"), got.Attr("missing"))
	}
	if !strings.Contains(got.Text, "Naslov") || !strings.HasPrefix(got.Text, "\n") {
		t.Errorf("Text = %q, want the untrimmed text content", got.Text)
	}
	if v := got.ChildText("h2"); v != "Naslov prvi" {
		t.Errorf("ChildText(h2) = %q, want trimmed %q", v, "Naslov prvi")
	}
	if v := got.ChildText("a"); v != "AB" {
		t.Errorf("ChildText(a) = %q, want every match joined, %q", v, "AB")
	}
	if v := got.ChildAttr("a", "href"); v != "/a" {
		t.Errorf("ChildAttr(a, href) = %q, want the first match trimmed, %q", v, "/a")
	}
	if v := got.ChildAttr("nav", "href"); v != "" {
		t.Errorf("ChildAttr on no match = %q, want empty", v)
	}
	if len(paragraphs) != 2 || paragraphs[0] != "jedan" || paragraphs[1] != "dva" || indexes[1] != 1 {
		t.Errorf("ForEach yielded %q at %v, want [jedan dva] at [0 1]", paragraphs, indexes)
	}
}

func TestCallbacksRunInRegistrationOrder(t *testing.T) {
	url := serveHTML(t, `<html><body><h2>a</h2><h3>b</h3><h2>c</h2><h3>d</h3></body></html>`)

	var order []string
	c := collector.New()
	c.OnHTML("h3", func(e *collector.HTMLElement) { order = append(order, e.Text) })
	c.OnHTML("h2", func(e *collector.HTMLElement) { order = append(order, e.Text) })
	if err := c.Visit(t.Context(), url); err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(order, ""); got != "bdac" {
		t.Errorf("callback order = %q, want %q", got, "bdac")
	}
}

func TestVisitSendsHeaders(t *testing.T) {
	var ua, accept string
	url := serve(t, func(w http.ResponseWriter, r *http.Request) {
		ua, accept = r.UserAgent(), r.Header.Get("Accept")
	})

	if err := collector.New(collector.UserAgent("TestBot/1.0")).Visit(t.Context(), url); err != nil {
		t.Fatal(err)
	}
	if ua != "TestBot/1.0" || accept != "*/*" {
		t.Errorf("User-Agent = %q, Accept = %q; want TestBot/1.0, */*", ua, accept)
	}
}

func TestVisitFailsOnErrorStatus(t *testing.T) {
	url := serve(t, func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "<html><body><h2>x</h2></body></html>", http.StatusServiceUnavailable)
	})

	ran := false
	c := collector.New()
	c.OnHTML("h2", func(*collector.HTMLElement) { ran = true })
	err := c.Visit(t.Context(), url)

	statusErr, ok := errors.AsType[*collector.StatusError](err)
	if !ok || statusErr.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("Visit error = %v, want a StatusError with 503", err)
	}
	if ran {
		t.Error("callback ran on an error response")
	}
}

func TestVisitSkipsNonHTML(t *testing.T) {
	url := serve(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"h2": "<h2>x</h2>"}`))
	})

	ran := false
	c := collector.New()
	c.OnHTML("h2", func(*collector.HTMLElement) { ran = true })
	if err := c.Visit(t.Context(), url); err != nil {
		t.Fatal(err)
	}
	if ran {
		t.Error("callback ran on a JSON response")
	}
}

func TestVisitDecodesDeclaredCharset(t *testing.T) {
	page, err := charmap.Windows1250.NewEncoder().String(`<html><body><h1>Čćšžđ</h1></body></html>`)
	if err != nil {
		t.Fatal(err)
	}
	url := serve(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=windows-1250")
		_, _ = w.Write([]byte(page))
	})

	var title string
	c := collector.New()
	c.OnHTML("h1", func(e *collector.HTMLElement) { title = e.Text })
	if err := c.Visit(t.Context(), url); err != nil {
		t.Fatal(err)
	}
	if title != "Čćšžđ" {
		t.Errorf("title = %q, want %q", title, "Čćšžđ")
	}
}

func TestVisitTruncatesAtMaxBodySize(t *testing.T) {
	url := serveHTML(t, `<html><body><h2>prvi</h2>`+strings.Repeat(" ", 1000)+`<h2>drugi</h2></body></html>`)

	var titles []string
	c := collector.New(collector.MaxBodySize(500))
	c.OnHTML("h2", func(e *collector.HTMLElement) { titles = append(titles, e.Text) })
	if err := c.Visit(t.Context(), url); err != nil {
		t.Fatal(err)
	}
	if len(titles) != 1 || titles[0] != "prvi" {
		t.Errorf("titles = %q, want only the one inside the cap", titles)
	}
}

func TestVisitKeepsCookiesAcrossRedirects(t *testing.T) {
	url := serve(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.SetCookie(w, &http.Cookie{Name: "consent", Value: "yes", Path: "/"})
			http.Redirect(w, r, "/page", http.StatusFound)
			return
		}
		if c, err := r.Cookie("consent"); err != nil || c.Value != "yes" {
			http.Error(w, "no cookie", http.StatusForbidden)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`<html><body><h2>ok</h2></body></html>`))
	})

	var title string
	c := collector.New()
	c.OnHTML("h2", func(e *collector.HTMLElement) { title = e.Text })
	if err := c.Visit(t.Context(), url); err != nil {
		t.Fatal(err)
	}
	if title != "ok" {
		t.Errorf("title = %q, want ok", title)
	}
}

func TestVisitHonorsContextAndTimeout(t *testing.T) {
	url := serve(t, func(_ http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-time.After(5 * time.Second):
		}
	})

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if err := collector.New().Visit(ctx, url); !errors.Is(err, context.Canceled) {
		t.Errorf("Visit with a canceled context = %v, want context.Canceled", err)
	}

	start := time.Now()
	if err := collector.New(collector.Timeout(100*time.Millisecond)).Visit(t.Context(), url); err == nil {
		t.Error("Visit past the timeout succeeded, want an error")
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("Visit took %v, want it cut off near the 100ms timeout", elapsed)
	}
}
