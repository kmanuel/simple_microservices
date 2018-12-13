FROM golang

RUN apt-get update
RUN apt-get install libfontconfig -y

RUN mkdir $GOPATH/src/app
WORKDIR $GOPATH/src/app

COPY . .

RUN go get ./
RUN go build

EXPOSE 8080

ENTRYPOINT ["go", "run", "main.go"]
