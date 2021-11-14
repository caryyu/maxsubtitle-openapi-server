FROM golang:1.17-alpine3.13 AS build

WORKDIR /app

ADD . ./
RUn go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
RUN go build -o /mast-server ./cmd/server

FROM alpine:3.13
WORKDIR /
COPY --from=build /mast-server /mast-server
ENTRYPOINT ["/mast-server"]
