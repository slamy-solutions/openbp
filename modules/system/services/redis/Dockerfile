FROM redis:7.0-alpine

LABEL version="1.0"
LABEL description="Multipurpose redis instance for OpenBP"

ENV MAXMEMORY 512mb
COPY entrypoint.sh /entrypoint.sh
COPY redis.conf /usr/local/etc/redis/redis.conf
ENTRYPOINT [ "sh", "/entrypoint.sh" ]