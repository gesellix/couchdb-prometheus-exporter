FROM alpine:3.11 AS builder
LABEL builder=true

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV APPPATH /app

RUN adduser -DH user
RUN apk add --update -t build-deps go git mercurial libc-dev gcc libgcc
COPY . $APPPATH
RUN cd $APPPATH && go get -d \
 && go test -short ./... \
 && go build \
    -a \
    -ldflags '-s -w -extldflags "-static"' \
    -o /bin/couchdb-prometheus-exporter

FROM scratch
LABEL maintainer="Tobias Gesellchen <tobias@gesellix.de> (@gesellix)"

ENV TELEMETRY.ADDRESS="0.0.0.0:9984"
ENV LOGTOSTDERR="true"

EXPOSE 9984
ENTRYPOINT [ "/couchdb-prometheus-exporter" ]
CMD [ ]

COPY --from=builder /etc/passwd /etc/passwd
USER user

COPY --from=builder /bin/couchdb-prometheus-exporter /couchdb-prometheus-exporter
