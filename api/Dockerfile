FROM golang:1.24-alpine AS builder

ENV GOTOOLCHAIN=auto
WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o acc-laptime-tracker-api ./api

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/acc-laptime-tracker-api .

EXPOSE 80

ENTRYPOINT ["./acc-laptime-tracker-api"]

