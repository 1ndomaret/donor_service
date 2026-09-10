FROM golang:1.26.4-alpine AS builder

WORKDIR /app

COPY . .

ENV GOFLAGS=-mod=mod
ENV GOMODCACHE=/tmp/gomodcache
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN GODEBUG=gomodindex=0 go build -v -x -o /app/donor-service ./app/cmd

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/donor-service /app/donor-service
COPY secret.env ./secret.env

EXPOSE 1324

CMD ["/app/donor-service"]