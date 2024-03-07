FROM golang:1.22-alpine as builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
RUN echo "I am running on ${BUILDPLATFORM}, building for ${TARGETPLATFORM}"
WORKDIR /src
RUN export GOOS=$(echo $OS_ARCH | cut -d'/' -f1)
RUN export GOARCH=$(echo $OS_ARCH | cut -d'/' -f2)
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
