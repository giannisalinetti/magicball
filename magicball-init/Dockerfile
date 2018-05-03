FROM golang:1.8

WORKDIR /go/src/magicball-init

COPY main.go .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["magicball-init"]
