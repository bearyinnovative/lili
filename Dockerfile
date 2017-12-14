FROM golang:1.9.2

MAINTAINER RoCry <rocry@bearyinnovative.com>

COPY . /go/src/github.com/bearyinnovative/lili
WORKDIR /go/src/github.com/bearyinnovative/lili

RUN go get -v github.com/kardianos/govendor && govendor sync

RUN go build -o main
CMD ["./main"]