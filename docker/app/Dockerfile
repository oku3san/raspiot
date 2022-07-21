##
## Build
##
FROM golang:latest as build

WORKDIR /go/src
COPY ./src /go/src

RUN env GOOS=linux GOARCH=arm GOARM=5 go build -o /go/src/raspiot /go/src/main.go

##
## Deploy
##
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /go/src/raspiot /raspiot

ENTRYPOINT ["/raspiot"]
