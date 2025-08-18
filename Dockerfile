FROM alpine:latest

RUN apk --no-cache add ca-certificates aom-dev libwebp-dev

COPY go-media-resizer /usr/local/bin/

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/go-media-resizer", "start"]
