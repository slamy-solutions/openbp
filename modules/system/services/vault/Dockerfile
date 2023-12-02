###### Build softhsm for simulating HSM
FROM --platform=${BUILDPLATFORM:-linux/amd64} alpine:3.17.3 as hsm-build

ARG SOFTHSM2_VERSION=2.6.1

ENV SOFTHSM2_VERSION=${SOFTHSM2_VERSION} \
    SOFTHSM2_SOURCES=/tmp/softhsm2

RUN apk --no-cache --update add \
        alpine-sdk \
        autoconf \
        automake \
        git \
        libtool \
        openssl-dev \
        sqlite-dev \
        sqlite

RUN git clone https://github.com/opendnssec/SoftHSMv2.git ${SOFTHSM2_SOURCES}
WORKDIR ${SOFTHSM2_SOURCES}

RUN git checkout ${SOFTHSM2_VERSION} -b ${SOFTHSM2_VERSION} \
    && sh autogen.sh \
    && ./configure --prefix=/usr/local --with-objectstore-backend-db --with-crypto-backend=openssl \
    && make -j 16 \
    && make install

###### Build vault code
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.21.4-alpine as build
# gcc, g++, make, build-base are needed for CGO
RUN apk add --no-cache git gcc g++ make build-base

ENV GOPATH=/src
RUN mkdir /src && cd /src && go work init

# Copy system libraries
COPY modules/system/libs/golang /src/modules/system/libs/golang
RUN cd /src && go work use ./modules/system/libs/golang
RUN cd /src/modules/system/libs/golang && go mod download

# Copy workspace
COPY modules/system/services/vault/go.mod /src/modules/system/services/vault/go.mod
COPY modules/system/services/vault/go.sum /src/modules/system/services/vault/go.sum
RUN cd /src && go work use ./modules/system/services/vault
RUN cd /src/modules/system/services/vault && go mod download

# Copy code
WORKDIR /src/modules/system/services/vault/src
COPY modules/system/services/vault/src/main.go ./main.go
COPY modules/system/services/vault/src/service ./service/
COPY modules/system/services/vault/src/pkcs ./pkcs/

# You have to enable CGO, because PKCS11 will not compile without it.
# CGO is used to directly communicate with PKCS11 provider C library.
RUN CGO_ENABLED=1 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -ldflags="-w -s" -a -o app ./main.go
RUN chmod +x app

###### Building final image
FROM --platform=${BUILDPLATFORM:-linux/amd64} alpine:3.17.3

RUN apk add --no-cache sqlite-libs openssl gcc

ENV HSM_PROVIDER=softhsm2
COPY modules/system/services/vault/softhsm2.conf /softhsm2.conf
ENV SOFTHSM2_CONF=/softhsm2.conf
ENV SOFTHSM2_PKCS11_LIBRARY_PATH=/libsofthsm2.so
ENV SOFTHSM2_PKCS11_TOKEN_LABEL=openbp

COPY --from=hsm-build /tmp/softhsm2/src/lib/.libs/libsofthsm2.so /libsofthsm2.so
COPY --from=hsm-build /tmp/softhsm2/src/bin/util/softhsm2-util /softhsm2-util
RUN mkdir /data && /softhsm2-util --module /libsofthsm2.so --init-token --free --label "openbp" --pin "12345678" --so-pin "12345678"

ENV DYNAMIC_PKCS11_LIBRARY_PATH=/pkcs/pkcs.so
ENV DYNAMIC_PKCS11_SLOT=0

COPY --from=build /src/modules/system/services/vault/src/app /app
CMD ["/app"]