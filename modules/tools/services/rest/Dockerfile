FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.18-alpine as build
RUN apk add --no-cache git

# Copy libraries
COPY modules/system/libs/golang /src/modules/system/libs/golang
COPY modules/native/libs/golang /src/modules/native/libs/golang

# Setup workspace
COPY modules/tools/services/rest/go.mod /src/modules/tools/services/rest/go.mod
COPY modules/tools/services/rest/go.sum /src/modules/tools/services/rest/go.sum

ENV GOPATH=/src
RUN cd /src && \
    go work init && \
	go work use ./modules/tools/services/rest && \
	go work use ./modules/system/libs/golang && \
	go work use ./modules/native/libs/golang

# Download dependencies for libraries and code
RUN cd /src/modules/system/libs/golang && go mod download
RUN cd /src/modules/native/libs/golang && go mod download
RUN cd /src/modules/tools/services/rest && go mod download

# Copy code
WORKDIR /src/modules/tools/services/rest/src
COPY modules/tools/services/rest/src /src/modules/tools/services/rest/src

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -ldflags="-w -s" -a -o app ./main.go
RUN chmod +x app

FROM scratch
ENV GIN_MODE=release
ENV PORT=80
EXPOSE 80
COPY --from=build /src/modules/tools/services/rest/src/app /app
CMD ["/app"]