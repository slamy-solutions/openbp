###### Build softhsm for simulating HSM
FROM --platform=${BUILDPLATFORM:-linux/amd64} alpine: as hsm-build

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
        sqlite-dev

RUN git clone https://github.com/opendnssec/SoftHSMv2.git ${SOFTHSM2_SOURCES}
WORKDIR ${SOFTHSM2_SOURCES}

RUN git checkout ${SOFTHSM2_VERSION} -b ${SOFTHSM2_VERSION} \
    && sh autogen.sh \
    && ./configure --prefix=/usr/local --with-objectstore-backend-db --with-crypto-backend=openssl \
    && make \
    && make install

RUN softhsm2-util --init-token --free --label "openbp" --pin "12345678" --so-pin "12345678"

###### Build vault code
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.18-alpine as build
RUN apk add --no-cache git

ENV GOPATH=/src
RUN mkdir /src && cd /src && go work init

# Copy and install libs
WORKDIR /src/modules/system/services/vault
COPY modules/system/services/vault/go.mod ./go.mod
COPY modules/system/services/vault/go.sum ./go.sum
RUN go mod download

# Copy code
COPY modules/system/services/vault/src/main.go ./main.go
COPY modules/system/services/vault/src/services ./services/

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -ldflags="-w -s" -a -o app ./main.go
RUN chmod +x app

###### Building final image
FROM --platform=${BUILDPLATFORM:-linux/amd64} alpine:3.17.3

RUN apk add --no-cache sqlite openssl

ENV HSM_PROVIDER=softhsm2
COPY modules/system/services/vault/softhsm2.conf /softhsm2.conf
ENV SOFTHSM2_CONF=/softhsm2.conf
ENV SOFTHSM2_PKCS11_LIBRARY_PATH=/libsofthsm2.so
ENV SOFTHSM2_PKCS11_TOKEN_LABEL=openbp

ENV DYNAMIC_PKCS11_LIBRARY_PATH=/pkcs/pkcs.so
ENV DYNAMIC_PKCS11_SLOT=0

COPY --from=hsm-build /usr/lib/softhsm/libsofthsm2.so /libsofthsm2.so

COPY --from=build /src/modules/native/services/namespace/src/app /app
CMD ["/app"]