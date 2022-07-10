##
## Build
##
FROM golang:latest as build

WORKDIR /app
COPY . /app

RUN env GOOS=linux GOARCH=arm GOARM=5 go build -o /app/raspiot /app/cmd/main.go

##
## Deploy
##
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /app/raspiot /raspiot

ENTRYPOINT ["/raspiot"]
