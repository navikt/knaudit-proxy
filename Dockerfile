FROM golang:1.22-alpine as builder

WORKDIR /src
COPY go.sum go.sum
COPY go.mod go.mod
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o knaudit-proxy .

RUN echo "Built binary architecture: $(go env GOARCH)"

FROM alpine:3

WORKDIR /app
COPY --from=builder /src/knaudit-proxy /app/knaudit-proxy
RUN uname -a
RUN chmod +x /app/knaudit-proxy
CMD ["/app/knaudit-proxy"]
