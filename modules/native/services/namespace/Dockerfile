FROM golang:1.18-alpine as build
RUN apk add --no-cache git

ENV GOPATH=/src
COPY go.work /src/go.work
COPY go.work.sum /src/go.work.sum

COPY modules/system/libs/go /src/modules/system/libs/go
COPY modules/native/services/namespace/go.mod /src/modules/native/services/namespace/go.mod
COPY modules/native/services/namespace/go.sum /src/modules/native/services/namespace/go.sum


RUN cd /src/modules/system/libs/go/cache && go mod download
RUN cd /src/modules/system/libs/go/mongodb && go mod download
RUN cd /src/modules/system/libs/go/telemetry && go mod download
RUN cd /src/modules/native/services/namespace && go mod download

# Code
WORKDIR /src/modules/native/services/namespace/src
COPY modules/native/services/namespace/src/main.go ./main.go
COPY modules/native/services/namespace/src/services ./services/
COPY modules/native/services/namespace/src/grpc ./grpc/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -o app ./main.go
RUN chmod +x app

FROM scratch
COPY --from=build /src/modules/native/services/namespace/src/app /app
CMD ["/app"]