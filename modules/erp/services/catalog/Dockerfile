FROM golang:1.20-alpine as build
RUN apk add --no-cache git

ENV GOPATH=/src

COPY go.work /src/go.work
COPY go.work.sum /src/go.work.sum

COPY lib /src/lib
COPY catalog/go.sum /src/catalog/go.sum
COPY catalog/go.mod /src/catalog/go.mod

RUN cd /src/lib && go mod download
RUN cd /src/catalog && go mod download

WORKDIR /src/catalog

COPY catalog/main.go .
COPY catalog/services ./services/
COPY catalog/grpc ./grpc/

RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app

FROM scratch
COPY --from=build /src/catalog/app /app
CMD ["./app"]