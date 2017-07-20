FROM alpine:edge
MAINTAINER Tobias Gesellchen <tobias@gesellix.de> (@gesellix)
EXPOSE 9984

ENV GOPATH /go
ENV APPPATH $GOPATH/src/github.com/gesellix/couchdb-prometheus-exporter
COPY . $APPPATH

RUN apk add --update -t build-deps go git mercurial libc-dev gcc libgcc \
    && cd $APPPATH && go get -d && go build -o /bin/couchdb-prometheus-exporter \
    && apk del --purge build-deps && rm -rf $GOPATH

ENTRYPOINT [ "/bin/couchdb-prometheus-exporter", "-telemetry.address=0.0.0.0:9984" ]
CMD [ "-logtostderr" ]
