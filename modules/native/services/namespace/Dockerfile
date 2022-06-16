FROM golang:1.18-alpine as build
RUN apk add --no-cache git

ENV GOPATH=/src

COPY go.sum /src/go.sum
COPY go.mod /src/go.mod
RUN cd /src/&& go mod download

WORKDIR /src

COPY main.go .
COPY services ./services/
COPY grpc ./grpc/

RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app

FROM scratch
COPY --from=build /src/app /app
CMD ["./app"]