FROM golang:1.25.5-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/cmd/url-shortener-app/exe /app/cmd/url-shortener-app