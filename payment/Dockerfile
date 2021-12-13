# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17-buster AS build

WORKDIR /app

COPY go.mod ./
#COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /payment

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /payment /payment

EXPOSE 8082

USER nonroot:nonroot

ENTRYPOINT ["/payment"]