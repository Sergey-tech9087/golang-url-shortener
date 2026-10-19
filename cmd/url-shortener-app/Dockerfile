FROM golang:1.26.0-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/cmd/url-shortener-app/exe /app/cmd/url-shortener-app

FROM alpine:3.23.4
WORKDIR /app
COPY --from=builder /app/cmd/url-shortener-app/exe /app
CMD ["/app/exe"]