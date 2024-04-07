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
  go build -v -ldflags="-w -s" -o /tmp/rpiwatchdog && \
  chown 1000 /tmp/rpiwatchdog

FROM scratch AS final
WORKDIR /app
USER 1000

COPY --link --from=build /tmp/rpiwatchdog /app/rpiwatchdog

# Port 1111 is only used if RPIW_SERVEHEALTHSOURCE is set to true
EXPOSE 1111


ENV RPIW_DEVICEPATH=/dev/watchdog \
  RPIW_SERVEHEALTHSOURCE=false \
  RPIW_USEHEALTHSOURCE= \
  RPIW_VERBOSELOGGING=false \
  RPIW_HEALTHCHECKINTERVAL=30 \
  RPIW_HEALTHCHECKTIMEOUT=3 

ENTRYPOINT [ "/app/rpiwatchdog" ]
CMD [ "watch" ]

HEALTHCHECK --interval=10s --timeout=2s --start-period=5s --retries=3 \
  CMD [ "/app/rpiwatchdog", "healthcheck" ]