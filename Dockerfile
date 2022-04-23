ARG GOLANG_VERSION=1.17
FROM golang:${GOLANG_VERSION}-buster as builder
ARG GOPROXY=https://goproxy.cn
WORKDIR ${GOPATH}/src/github.com/projectxpolaris/youauth

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ${GOPATH}/bin/youauth ./main.go

FROM debian:buster-slim

COPY --from=builder /usr/local/lib /usr/local/lib
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY ./static /app/static
COPY ./templates /app/templates
COPY --from=builder /go/bin/youauth /usr/local/bin/youauth
WORKDIR /app
ENTRYPOINT ["/usr/local/bin/youauth","run"]