# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM ghcr.io/kwitsch/ziggoimg AS build

RUN --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=bind,source=go.mod,target=go.mod \
  --mount=type=cache,target=/root/.cache/go-build \ 
  --mount=type=cache,target=/go/pkg \
  go mod download

RUN --mount=type=bind,target=. \
  --mount=type=cache,target=/root/.cache/go-build \ 
  --mount=type=cache,target=/go/pkg \
  go build -v -ldflags="-w -s" -o /tmp/watchdog && \
  chown 1000 /tmp/watchdog

FROM scratch
WORKDIR /var/run
WORKDIR /app
USER 1000

COPY --link --from=build /tmp/watchdog /app/watchdog

ENTRYPOINT [ "/app/watchdog" ]