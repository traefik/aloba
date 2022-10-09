# syntax=docker/dockerfile:1.2
FROM alpine

RUN apk --no-cache --no-progress add ca-certificates \
    && rm -rf /var/cache/apk/*

COPY aloba /

ENTRYPOINT ["/aloba"]
