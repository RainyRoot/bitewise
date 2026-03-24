FROM golang:1.26-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /build

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /server ./cmd/server/

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -h /app appuser

WORKDIR /app

COPY --from=builder /server .
COPY --from=builder /build/migrations/ ./migrations/

RUN mkdir -p /data && chown appuser:appuser /data

USER appuser

ENV PORT=8080
ENV DATABASE_PATH=/data/bitewise.db

EXPOSE 8080

VOLUME ["/data"]

ENTRYPOINT ["./server"]
