FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.18-alpine as build
RUN apk add --no-cache git

# Copy libraries
COPY modules/system/libs/golang /src/modules/system/libs/golang
COPY modules/native/libs/golang /src/modules/native/libs/golang

COPY modules/native/services/iam/oauth/go.mod /src/modules/native/services/iam/oauth/go.mod
COPY modules/native/services/iam/oauth/go.sum /src/modules/native/services/iam/oauth/go.sum

ENV GOPATH=/src
RUN cd /src && \
    go work init && \
	go work use ./modules/native/services/iam/oauth && \
	go work use ./modules/system/libs/golang && \
	go work use ./modules/native/libs/golang

RUN cd /src/modules/system/libs/golang && go mod download
RUN cd /src/modules/native/libs/golang && go mod download
RUN cd /src/modules/native/services/iam/oauth && go mod download

# Code
WORKDIR /src/modules/native/services/iam/oauth/src
COPY modules/native/services/iam/oauth/src/main.go ./main.go
COPY modules/native/services/iam/oauth/src/services ./services/

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -ldflags="-w -s" -a -o app ./main.go
RUN chmod +x app

FROM scratch
COPY --from=build /src/modules/native/services/iam/oauth/src/app /app
CMD ["/app"]