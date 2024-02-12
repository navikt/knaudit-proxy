FROM golang:1.22-alpine as builder

WORKDIR /src
COPY go.sum go.sum
COPY go.mod go.mod
RUN go mod download
COPY . .
RUN go build -o knaudit-proxy .

FROM alpine:3

WORKDIR /app
COPY --from=builder /src/knaudit-proxy /app/knaudit-proxy
CMD ["/app/knaudit-proxy"]