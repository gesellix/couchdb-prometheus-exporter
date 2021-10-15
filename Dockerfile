FROM golang:1.17-alpine AS builder
LABEL builder=true

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV APPPATH /app

#RUN apk add --update -t build-deps go git mercurial libc-dev gcc libgcc
COPY . $APPPATH
RUN cd $APPPATH && go get -d \
 && go test -short ./... \
 && go build \
    -a \
    -ldflags '-s -w -extldflags "-static"' \
    -o /bin/main

FROM alpine:3.14.2
LABEL maintainer="Tobias Gesellchen <tobias@gesellix.de> (@gesellix)"

ENV TELEMETRY_ADDRESS="0.0.0.0:9984"
ENV LOGTOSTDERR="true"

EXPOSE 9984
ENTRYPOINT [ "/couchdb-prometheus-exporter" ]
CMD [ ]

RUN apk --no-cache add ca-certificates \
 && adduser -DH user
USER user

COPY --from=builder /bin/main /couchdb-prometheus-exporter
