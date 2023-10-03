FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.21-alpine as build
RUN apk add --no-cache git

ENV GOPATH=/src
RUN mkdir /src && cd /src && go work init

# Setup libraries
COPY modules/system/libs/golang /src/modules/system/libs/golang
RUN cd /src && go work use ./modules/system/libs/golang
RUN cd /src/modules/system/libs/golang && go mod download

COPY modules/runtime/libs/golang /src/modules/runtime/libs/golang
RUN cd /src && go work use ./modules/runtime/libs/golang
RUN cd /src/modules/runtime/libs/golang && go mod download

# Setup workspace
COPY modules/runtime/services/manager/go.mod /src/modules/runtime/services/manager/go.mod
COPY modules/runtime/services/manager/go.sum /src/modules/runtime/services/manager/go.sum
RUN cd /src && go work use ./modules/runtime/services/manager
RUN cd /src/modules/runtime/services/manager && go mod download

# Copy code
WORKDIR /src/modules/runtime/services/manager/src
COPY modules/runtime/services/manager/src /src/modules/runtime/services/manager/src

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -ldflags="-w -s" -a -o app ./main.go
RUN chmod +x app

FROM scratch
ENV PORT=80
EXPOSE 80
COPY --from=build /src/modules/runtime/services/manager/src/app /app
CMD ["/app"]