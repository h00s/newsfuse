FROM golang:alpine AS backend
WORKDIR /app
ENV CGO_ENABLED=0 
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend ./
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/newsfuse

FROM oven/bun:latest AS frontend
WORKDIR /app
COPY frontend/package.json frontend/bun.lockb ./
RUN bun install --frozen-lockfile
COPY frontend ./
RUN bun run build

FROM gcr.io/distroless/static-debian13:latest
WORKDIR /app
COPY --from=backend /out/newsfuse ./
COPY --from=frontend /app/build ./public

EXPOSE 3000

ENTRYPOINT ["./newsfuse"]