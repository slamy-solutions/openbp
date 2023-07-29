FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.20-alpine as build
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

COPY modules/iot/libs/golang /src/modules/iot/libs/golang
RUN cd /src && go work use ./modules/iot/libs/golang
RUN cd /src/modules/iot/libs/golang && go mod download

# Setup workspace
COPY modules/iot/services/core/go.mod /src/modules/iot/services/core/go.mod
COPY modules/iot/services/core/go.sum /src/modules/iot/services/core/go.sum
RUN cd /src && go work use ./modules/iot/services/core
RUN cd /src/modules/iot/services/core && go mod download

# Copy code
WORKDIR /src/modules/iot/services/core/src
COPY modules/iot/services/core/src /src/modules/iot/services/core/src

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -ldflags="-w -s" -a -o app ./main.go ./events.go
RUN chmod +x app

FROM scratch
ENV GIN_MODE=release
ENV PORT=80
EXPOSE 80
COPY --from=build /src/modules/iot/services/core/src/app /app
CMD ["/app"]