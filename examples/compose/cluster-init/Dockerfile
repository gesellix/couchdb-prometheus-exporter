FROM alpine:edge

RUN apk add --no-cache curl

COPY wait-for /wait-for
COPY await-nodes.sh /await-nodes.sh
COPY init-cluster.sh /init-cluster.sh
COPY docker-entrypoint.sh /docker-entrypoint.sh

CMD [ "/docker-entrypoint.sh" ]
