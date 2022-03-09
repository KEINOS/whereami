# =============================================================================
#  Docker Container for Local Build
# =============================================================================
#  This Dockerfile will create a light weight container image that only contains
#  the `whereami` binary.
#  Use this container if you don't have Go installed or if you don't want to
#  install the command locally.

# -----------------------------------------------------------------------------
#  Build Stage
# -----------------------------------------------------------------------------
FROM golang:alpine AS build

RUN apk add --no-cache \
#    alpine-sdk \
#    build-base \
    ca-certificates

COPY . /workspace

WORKDIR /workspace

ENV CGO_ENABLED 0

RUN ls -lah

RUN \
    go build \
        # Static linking and shrink size
        -ldflags="-s -w -extldflags \"-static\"" \
        # Outpath
        -o /go/bin/whereami \
        # Path to main
        ./cmd/whereami/main.go \
    # Smoke test
    && /go/bin/whereami -h

# -----------------------------------------------------------------------------
#  Main Stage
# -----------------------------------------------------------------------------
FROM scratch

COPY --from=build /go/bin/whereami /usr/bin/whereami
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/usr/bin/whereami"]
