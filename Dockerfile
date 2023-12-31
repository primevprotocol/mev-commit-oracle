FROM golang:1.21.1 AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o mev-commit-oracle ./cmd/main.go

FROM alpine:latest

COPY --from=builder /app/mev-commit-oracle /usr/local/bin/mev-commit-oracle

EXPOSE 8080

ENTRYPOINT ["mev-commit-oracle"]
