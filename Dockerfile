FROM golang:1.7

RUN go get github.com/vtrifonov/http-api-mock

RUN mkdir /config
VOLUME /config

EXPOSE 8082 8083

ENTRYPOINT ["/go/bin/http-api-mock","-config-path","/config"]