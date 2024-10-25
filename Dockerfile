FROM golang:1.23.2 as BUILD

WORKDIR /app

COPY . .

RUN make build

FROM scratch

LABEL org.opencontainers.image.source=https://github.com/logabot/cron
LABEL org.opencontainers.image.description="cron docker image"
LABEL org.opencontainers.image.licenses=APACHE

WORKDIR /app

COPY --from=BUILD /app/out/cron /app
