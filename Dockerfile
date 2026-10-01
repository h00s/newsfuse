FROM golang:alpine AS backend
WORKDIR /app
ENV CGO_ENABLED=0
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend ./
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/newsfuse

FROM oven/bun:latest AS frontend
WORKDIR /app
COPY frontend/package.json frontend/bun.lock ./
RUN bun install --frozen-lockfile
COPY frontend ./
RUN bun run build

FROM gcr.io/distroless/static-debian13:latest
WORKDIR /app
COPY --from=backend /out/newsfuse ./
COPY --from=frontend /app/build ./public
# Raptor binds 127.0.0.1 by default, which is unreachable from outside the container.
ENV SERVER_ADDRESS=0.0.0.0
# On stop, keep serving with /readyz failing for 2 seconds before draining. Kept small because
# Docker kills a container 10 seconds after asking it to stop.
ENV SERVER_SHUTDOWN_DELAY=2
# The site is served over HTTPS only. Start HSTS small; raise it to 31536000 once that's settled.
ENV APP_SECURE_HSTS_MAX_AGE=300

EXPOSE 3000

ENTRYPOINT ["./newsfuse"]
