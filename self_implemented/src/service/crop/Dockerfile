FROM golang

RUN mkdir $GOPATH/src/app
WORKDIR $GOPATH/src/app

COPY . .

RUN go get ./
RUN go build
