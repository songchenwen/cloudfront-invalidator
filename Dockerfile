FROM golang:alpine as builder

WORKDIR /src
COPY . /src
RUN go mod download && \
    go build -o bin/cfi && \
    mv bin/cfi /cfi

FROM alpine:latest
LABEL org.opencontainers.image.source="https://github.com/songchenwen/cloudfront-invalidator"

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /cfi /cfi
ENTRYPOINT ["/cfi"]
