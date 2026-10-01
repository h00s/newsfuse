# NewsFuse

The News Aggregator Application is a web-based platform that aggregates news articles from various internet news sites and displays them in a compact format. This application aims to simplify the process of keeping up with the latest news by collecting all the relevant articles in one place.

## Features

- ✅ Scrapes news articles from multiple sources
- ✅ Summarizes articles with AI for quick and easy reading
- ✅ Stores articles in a database for easy retrieval and display
- ✅ Displays articles in a compact format for easy reading
- ✅ Allows users to filter articles by source
- ✅ Provides a search function for finding specific articles
- ✅ Automatically updates with the latest news articles on a regular basis
- ✅ Keyboard shortcuts (press `?` in the app)

## Technologies Used

- Go with [Raptor](https://github.com/go-raptor/raptor), Bun and Goose on PostgreSQL
- SvelteKit with Svelte 5 runes, as a static single-page app
- Tailwind CSS 4 with daisyUI 5
- OpenAI for summaries
- Docker

## Development

The Go server serves the API under `/api/v1`. In development, Vite serves the frontend and proxies `/api` to the Go server, so the browser talks to one origin.

```bash
# terminal 1 — backend on :3000 (needs backend/.raptor.dev.yaml, see below)
cd backend && raptor db migrate up && raptor dev

# terminal 2 — frontend on :5173
cd frontend && bun install && bun run dev
```

`backend/.raptor.dev.yaml` is not tracked, because it holds credentials:

```yaml
database:
  host: localhost
  port: 5432
  username: dev
  password: …
  name: newsfuse_development

app:
  spa_optional: "true" # boot without a frontend build
  openai_key: …
  # scrapers_enabled: "false" # a boolean: don't scrape the news sites
  # llm_provider: "stub"      # summarize without calling OpenAI
```

## Testing

```bash
# backend: integration tests against a separate database (scrapers off, stub LLM)
cd backend
DATABASE_NAME=newsfuse_test raptor db migrate up
DATABASE_PASSWORD=… go test ./...

# frontend: type check and unit tests
cd frontend && bun run check && bun run test
```

## Deployment

The Docker image builds both halves; the Go binary serves the frontend build from `public/`. Configure it with environment variables:

| Variable | |
|---|---|
| `DATABASE_HOST`, `DATABASE_PORT`, `DATABASE_USERNAME`, `DATABASE_PASSWORD`, `DATABASE_NAME` | PostgreSQL |
| `DATABASE_AUTO_MIGRATE=true` | apply pending migrations at boot |
| `APP_OPENAI_KEY` | required: the server won't start without it |
| `SERVER_IP_EXTRACTOR=x-forwarded-for` | behind a reverse proxy, so rate limits see client addresses |
| `SERVER_SHUTDOWN_DELAY` | seconds to keep serving, with `/readyz` failing, before draining on stop; the image sets 2 |
| `APP_SECURE_HSTS_MAX_AGE` | HSTS max-age in seconds; the image sets 300, raise it to 31536000 once HTTPS is settled |

Serve it over HTTPS, or keep the `Host` header at the proxy: the server rejects cross-origin writes by comparing `Origin` with `Host`.

`GET /healthz` answers 200 while the process runs. `GET /readyz` answers 503 once shutdown begins, or when the database doesn't answer within 2 seconds; give a readiness probe a timeout of at least 3 seconds. Every response carries an `X-Request-Id`, which the log lines repeat as `request_id`.

## Contributing

Contributions to the News Aggregator Application are always welcome! If you would like to contribute, please fork this repository and submit a pull request.

## License

This project is licensed under the [MIT License](https://opensource.org/licenses/MIT).
