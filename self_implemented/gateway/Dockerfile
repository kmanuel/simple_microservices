FROM golang

RUN mkdir $GOPATH/src/app
WORKDIR $GOPATH/src/app

COPY . .

RUN go get ./
RUN go build

ENTRYPOINT ["go", "run", "main.go"]