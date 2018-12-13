FROM golang

RUN mkdir --parent /go/src/github.com/kmanuel/dockertest

WORKDIR /go/src/github.com/kmanuel/dockertest

COPY . .

RUN go get ./

CMD ["go", "run", "main.go"]