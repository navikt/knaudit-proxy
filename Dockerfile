FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.22-alpine as builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /src
COPY go.sum go.sum
COPY go.mod go.mod
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o knaudit-proxy .

FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:3

WORKDIR /app
COPY --from=builder /src/knaudit-proxy /app/knaudit-proxy
RUN uname -a
RUN chmod +x /app/knaudit-proxy
CMD ["/app/knaudit-proxy"]
