FROM golang:1.21.1 AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o mev-oracle ./main.go

FROM alpine:latest

RUN apk --no-cache add curl
RUN apk add --no-cache jq

COPY --from=builder /app/mev-oracle /mev-oracle
COPY --from=builder /app/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 13522 13523 13524 8080

ENTRYPOINT ["/entrypoint.sh"]