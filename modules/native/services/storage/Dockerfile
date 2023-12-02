FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.21.4-alpine as build
RUN apk add --no-cache git

ENV GOPATH=/src
RUN mkdir /src && cd /src && go work init

# Setup libraries

COPY modules/system/libs/golang /src/modules/system/libs/golang
RUN cd /src && go work use ./modules/system/libs/golang
RUN cd /src/modules/system/libs/golang && go mod download

COPY modules/native/libs/golang /src/modules/native/libs/golang
RUN cd /src && go work use ./modules/native/libs/golang
RUN cd /src/modules/native/libs/golang && go mod download

# Setup workspace
COPY modules/native/services/storage/go.mod /src/modules/native/services/storage/go.mod
COPY modules/native/services/storage/go.sum /src/modules/native/services/storage/go.sum
RUN cd /src && go work use ./modules/native/services/storage
RUN cd /src/modules/native/services/storage && go mod download

# Copy code
WORKDIR /src/modules/native/services/storage/src
COPY modules/native/services/storage/src/main.go ./main.go
COPY modules/native/services/storage/src/eventHandler.go ./eventHandler.go
COPY modules/native/services/storage/src/services ./services/

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -ldflags="-w -s" -a -o app ./main.go ./eventHandler.go
RUN chmod +x app

FROM scratch
COPY --from=build /src/modules/native/services/storage/src/app /app
CMD ["/app"]