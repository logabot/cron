# Cron: simple cron binary for container

[![CI State](https://github.com/go-co-op/gocron/actions/workflows/go_test.yml/badge.svg?branch=v2&event=push)](https://github.com/go-co-op/gocron/actions)

## Motivation

> [!WARNING]
> Sometimes tou need to run schedule tasks inside big image, or project with long startup time and you can't use kubernetes cronJob
> This is dirty hack, but sometimes it's needed

## Usage

Copy binary to your image and use

```docker
FROM ghcr.io/logabot/cron:v0.1.0 as CRON

....

COPY --from=CRON /app/cron ./

CMD["/app/cron"]
```

## Parameters

cli
```txt
Usage of cron:
  -config string
        crontab file. env: CONFIG (default "config")
  -entrypoint string
        shell script filename that runs before cron start. env: ENTRYPOINT
  -shell string
        define shell for run cron command. env: SHELL (default "/bin/bash")
```
env

## Config

Config is crontab, but with name

```crontab
2 9 * * * MyJob go build -o ./out && date
````

> [!WARNING]
> You can define many crontabs via many lines, but it not tested


## Entrypoint

Use `-entrypoint` for correspond shell script that run before start `cron`. e.g `npm install`
