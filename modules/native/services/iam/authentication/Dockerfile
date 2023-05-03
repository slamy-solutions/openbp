FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.18-alpine as build
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
COPY modules/native/services/iam/authentication/go.mod /src/modules/native/services/iam/authentication/go.mod
COPY modules/native/services/iam/authentication/go.sum /src/modules/native/services/iam/authentication/go.sum
RUN cd /src && go work use ./modules/native/services/iam/authentication
RUN cd /src/modules/native/services/iam/authentication && go mod download

# Code
WORKDIR /src/modules/native/services/iam/authentication/src
COPY modules/native/services/iam/authentication/src/main.go ./main.go
COPY modules/native/services/iam/authentication/src/services ./services/
COPY modules/native/services/iam/authentication/src/eventHandling ./eventHandling/

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -ldflags="-w -s" -a -o app ./main.go
RUN chmod +x app

FROM scratch
COPY --from=build /src/modules/native/services/iam/authentication/src/app /app
CMD ["/app"]