FROM golang:1.25.5-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/shucompress ./cmd

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    ghostscript \
    libreoffice \
    fonts-dejavu \
  && rm -rf /var/lib/apt/lists/*

RUN useradd -r -u 10001 -g nogroup appuser

COPY --from=builder /out/shucompress /app/shucompress

RUN mkdir -p /app/tmp && chown -R appuser:nogroup /app

USER appuser

EXPOSE 8085

CMD ["/app/shucompress"]
