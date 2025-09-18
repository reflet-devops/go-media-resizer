FROM golang:1.24-bookworm AS builder
ARG VERSION="main"
ARG COMMIT_SHORT="snapshot"

RUN apt-get update \
  && apt-get install --force-yes -y \
  libwebp-dev \
  pkg-config

RUN mkdir /app
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
ENV LDFLAGS=" -s -w -X github.com/reflet-devops/go-media-resizer/version.Version=${VERSION} -X github.com/reflet-devops/go-media-resizer/version.Commit=${COMMIT_SHORT}"
COPY . .
RUN CGO_ENABLED=1 go build -tags=nodynamic -ldflags "${LDFLAGS}" -o go-media-resizer

FROM debian:12

RUN apt-get update \
  && apt-get install --force-yes -y \
  libaom-dev \
  libavif-dev \
  libwebp-dev \
  && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /etc/go-media-resizer /var/log/go-media-resizer /var/run/go-media-resizer && \
    chmod 774 /etc/go-media-resizer /var/log/go-media-resizer /var/run/go-media-resizer

COPY docker-entrypoint /usr/local/bin/
COPY --from=builder /app/go-media-resizer /usr/local/bin/

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/docker-entrypoint"]
