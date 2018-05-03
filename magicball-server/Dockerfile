FROM golang:1.8

WORKDIR /go/src/magicball-server

COPY main.go .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8080

CMD ["magicball-server"]
