FROM mongo:5.0.5-focal AS rs-mongo

# Make MongoDB a replica set to support transactions. Based on https://stackoverflow.com/a/68621185/1952977
RUN apt-get update && apt-get install patch

# How to create scripts/docker-entrypoint.sh.patch
# 1. Download the original file:
#    wget https://github.com/docker-library/mongo/raw/master/5.0/docker-entrypoint.sh
#    ( wget https://github.com/docker-library/mongo/raw/b5c0cd58cb5626fee4d963ce05ba4d9026deb265/5.0/docker-entrypoint.sh )
# 2. Make a copy of it:
#    cp docker-entrypoint.sh docker-entrypoint-patched.sh
# 3. Add required modifications to docker-entrypoint-patched.sh
# 4. Create patch:
#    diff -u docker-entrypoint.sh docker-entrypoint-patched.sh > scripts/docker-entrypoint.sh.patch
# 5. Clean up:
#    rm docker-entrypoint.sh docker-entrypoint-patched.sh
COPY docker-entrypoint.sh.patch .
RUN patch /usr/local/bin/docker-entrypoint.sh docker-entrypoint.sh.patch
RUN mkdir -p /etc/mongo-key && chown mongodb:mongodb /etc/mongo-key

CMD ["--replSet", "rs", "--keyFile", "/etc/mongo-key/mongodb.key"]