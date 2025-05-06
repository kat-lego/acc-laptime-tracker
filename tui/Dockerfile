FROM golang:1.24 AS builder

ENV GOTOOLCHAIN=auto

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o acc-laptime-tracker ./tui

FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    ca-certificates \
    openssh-client \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/acc-laptime-tracker .

EXPOSE 8080

ENTRYPOINT ["./acc-laptime-tracker"]
