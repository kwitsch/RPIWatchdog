# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM ghcr.io/kwitsch/ziggoimg:dev AS build

RUN --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=bind,source=go.mod,target=go.mod \
  --mount=type=cache,target=/root/.cache/go-build \ 
  --mount=type=cache,target=/go/pkg \
  go mod download

RUN --mount=type=bind,target=. \
  --mount=type=cache,target=/root/.cache/go-build \ 
  --mount=type=cache,target=/go/pkg \
  mkdir -p /app && \
  go build -v -ldflags="-w -s" -o /app/rpiwatchdog

FROM scratch AS final

COPY --link --from=build --chmod=755  /app /app

# Port 1111 is only used if RPIW_SERVEHEALTHSOURCE is set to true
EXPOSE 1111


ENV RPIW_SERVEHEALTHSOURCE=false \
  RPIW_USEHEALTHSOURCE= \
  RPIW_VERBOSELOGGING=false \
  RPIW_HEALTHCHECKTIMEOUT=3 

ENTRYPOINT [ "/app/rpiwatchdog" ]
CMD [ "watch" ]

HEALTHCHECK --interval=10s --timeout=2s --start-period=5s --retries=3 \
  CMD [ "/app/rpiwatchdog", "healthcheck" ]