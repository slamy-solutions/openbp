FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.21.4-alpine as build
RUN apk add --no-cache git

ENV GOPATH=/src
RUN mkdir /src && cd /src && go work init

# Setup libraries
COPY modules/system/libs/golang /src/modules/system/libs/golang
RUN cd /src && go work use ./modules/system/libs/golang
RUN cd /src/modules/system/libs/golang && go mod download

#COPY modules/native/libs/golang /src/modules/native/libs/golang
#RUN cd /src && go work use ./modules/native/libs/golang
#RUN cd /src/modules/native/libs/golang && go mod download

COPY modules/erp/libs/golang /src/modules/erp/libs/golang
RUN cd /src && go work use ./modules/erp/libs/golang
RUN cd /src/modules/erp/libs/golang && go mod download

# Setup workspace
COPY modules/erp/services/core/go.mod /src/modules/erp/services/core/go.mod
COPY modules/erp/services/core/go.sum /src/modules/erp/services/core/go.sum
RUN cd /src && go work use ./modules/erp/services/core
RUN cd /src/modules/erp/services/core && go mod download

# Copy code
WORKDIR /src/modules/erp/services/core/src
COPY modules/erp/services/core/src /src/modules/erp/services/core/src

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -ldflags="-w -s" -a -o app ./main.go
RUN chmod +x app

FROM scratch
ENV GIN_MODE=release
ENV PORT=80
EXPOSE 80
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /src/modules/erp/services/core/src/app /app
CMD ["/app"]