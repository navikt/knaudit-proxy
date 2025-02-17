FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.24-alpine as builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app/
ADD . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -ldflags="-w -s" -o knaudit-proxy cmd/knaudit-proxy/main.go
RUN chmod +x knaudit-proxy

FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:3
WORKDIR /app
COPY --from=builder /app/knaudit-proxy /app/knaudit-proxy
CMD ["/app/knaudit-proxy", "-backend-type", "oracle"]
