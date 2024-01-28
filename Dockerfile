FROM golang:alpine AS backend

WORKDIR /app

COPY backend ./

RUN go mod download && \
    go build -o /out/newsfuse

FROM oven/bun:latest AS frontend

WORKDIR /app

COPY frontend ./

RUN bun install --frozen-lockfile && \
    bun run build

FROM alpine:latest

WORKDIR /app

COPY --from=backend /out/newsfuse ./
COPY --from=frontend /app/build ./public

EXPOSE 3000

CMD ["./newsfuse"]