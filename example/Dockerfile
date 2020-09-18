FROM golang:1.15.2

RUN go get github.com/fogfish/sideshell

ENTRYPOINT go run $GOPATH/src/github.com/fogfish/sideshell/example/main.go
