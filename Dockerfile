FROM golang:1.23.2 as BUILD

WORKDIR /app

COPY . .

RUN make build

FROM scratch

WORKDIR /app

COPY --from=BUILD /app/out/cron /app
