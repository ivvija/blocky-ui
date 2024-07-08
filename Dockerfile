
# Most of the logic is taken from
# https://github.com/0xERR0R/blocky/blob/main/Dockerfile

FROM --platform=$BUILDPLATFORM ghcr.io/kwitsch/ziggoimg AS build

ARG VERSION

# download packages
# bind mount go.mod and go.sum
# use cache for go packages
RUN --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    go mod download

# build binary
# bind mount source code
# use cache for go packages
RUN --mount=type=bind,target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    go build \
      -tags static  \
      -v  \
      -o /bin/blocky-ui


FROM scratch

ARG VERSION

LABEL org.opencontainers.image.title="blocky-ui" \
      org.opencontainers.image.vendor="ivvija" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.description="Simple web interface for blocky" \
      org.opencontainers.image.url="https://github.com/ivvija/blocky-ui#readme" \
      org.opencontainers.image.source="https://github.com/ivvija/blocky"

USER 100
WORKDIR /app

COPY --link --from=build /bin/blocky-ui /app/blocky-ui

ENV HOST=0.0.0.0
ENV PORT=3000
ENV API_BASE_URL=http://blocky:4000/api
ENV PAUSE_DURATION=5m

ENTRYPOINT ["/app/blocky-ui"]
