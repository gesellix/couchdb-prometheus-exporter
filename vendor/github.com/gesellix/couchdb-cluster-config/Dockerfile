FROM alpine:3.9 AS builder
LABEL builder=true

ENV CGO_ENABLED=0
ENV GOPATH /go
ENV APPPATH $GOPATH/src/github.com/gesellix/couchdb-cluster-config

RUN adduser -DH user
RUN apk add --update -t build-deps go git mercurial libc-dev gcc libgcc
COPY . $APPPATH
RUN cd $APPPATH && go get -d \
 && go build \
    -a \
    -ldflags '-extldflags "-static"' \
    -o /bin/couchdb-cluster-config

FROM scratch
LABEL maintainer="Tobias Gesellchen <tobias@gesellix.de> (@gesellix)"

ENTRYPOINT [ "/couchdb-cluster-config" ]
CMD [ "--help" ]

COPY --from=builder /etc/passwd /etc/passwd
USER user

COPY --from=builder /bin/couchdb-cluster-config /couchdb-cluster-config
