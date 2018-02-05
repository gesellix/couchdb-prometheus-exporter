FROM alpine:edge AS builder
LABEL builder=true

ENV CGO_ENABLED=0
ENV GOPATH /go
ENV APPPATH $GOPATH/src/github.com/gesellix/couchdb-prometheus-exporter

RUN adduser -DH user
RUN apk add --update -t build-deps go git mercurial libc-dev gcc libgcc
COPY . $APPPATH
RUN cd $APPPATH && go get -d \
 && go build \
    -a \
    -ldflags '-extldflags "-static"' \
    -o /bin/couchdb-prometheus-exporter

FROM scratch
LABEL maintainer="Tobias Gesellchen <tobias@gesellix.de> (@gesellix)"

EXPOSE 9984
ENTRYPOINT [ "/couchdb-prometheus-exporter", "-telemetry.address=0.0.0.0:9984" ]
CMD [ "-logtostderr" ]

COPY --from=builder /etc/passwd /etc/passwd
USER user

COPY --from=builder /bin/couchdb-prometheus-exporter /couchdb-prometheus-exporter
