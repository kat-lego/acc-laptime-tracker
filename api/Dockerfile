# Use Alpine-based Go image for a smaller build environment
FROM golang:1.24-alpine AS builder

# Enable module-aware mode and Go toolchain
ENV GOTOOLCHAIN=auto
WORKDIR /app

# Install git (required by some Go modules), ca-certificates for HTTPS
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o acc-laptime-tracker-api ./api

# Use Alpine as the final image
FROM alpine:latest

# Install certificates if your app makes HTTPS requests
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/acc-laptime-tracker-api .

EXPOSE 80

ENTRYPOINT ["./acc-laptime-tracker-api"]

