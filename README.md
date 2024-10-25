# Cron

## Motivation

[!WARNING]
> Sometimes tou need to run schedule tasks inside big image, or project with long startup time and you can't use kubernetes cronJob
> This is dirty hack, but sometimes it's needed

## Usage

Copy binary to image 

```
FROM logabot/cron:latest as CRON

....

COPY --from=CRON /app/cron ./
```
