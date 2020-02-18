FROM golang:alpine

RUN apk update --no-cache && \
    apk upgrade --no-cache && \
    apk add curl docker git su-exec --no-cache && \
    go get -u golang.org/x/lint/golint && \
    chmod u+s /sbin/su-exec

ENV CGO_ENABLED 0

WORKDIR /go-mongodb

CMD /bin/sh -c '\
    su-exec root addgroup -g ${DEV_GID} dev || true && \
    su-exec root adduser -u ${DEV_UID} -D -G dev dev || true && \
    tail -f /dev/null'
