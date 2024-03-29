FROM golang:1.17-alpine3.13 AS build

WORKDIR /app

ADD . ./
RUN apk add --no-cache git make build-base
RUN go mod download
RUN go build -o /mast-server ./cmd/server

FROM alpine:3.13
WORKDIR /
COPY --from=build /mast-server /mast-server
ENTRYPOINT ["/mast-server"]
