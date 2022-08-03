FROM golang:1.18-alpine as build
RUN apk add --no-cache git

COPY modules/system/libs/go /src/modules/system/libs/go
COPY modules/native/services/iam/authentication/password/go.mod /src/modules/native/services/iam/authentication/password/go.mod
COPY modules/native/services/iam/authentication/password/go.sum /src/modules/native/services/iam/authentication/password/go.sum

ENV GOPATH=/src
RUN cd /src && \
    go work init && \
	go work use ./modules/native/services/iam/authentication/password && \
	go work use ./modules/system/libs/go/mongodb && \
	go work use ./modules/system/libs/go/telemetry

RUN cd /src/modules/system/libs/go/mongodb && go mod download
RUN cd /src/modules/system/libs/go/telemetry && go mod download
RUN cd /src/modules/native/services/iam/authentication/password && go mod download

# Code
WORKDIR /src/modules/native/services/iam/authentication/password/src
COPY modules/native/services/iam/authentication/password/src/main.go ./main.go
COPY modules/native/services/iam/authentication/password/src/services ./services/
COPY modules/native/services/iam/authentication/password/src/grpc ./grpc/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -o app ./main.go
RUN chmod +x app

FROM scratch
COPY --from=build /src/modules/native/services/iam/authentication/password/src/app /app
CMD ["/app"]