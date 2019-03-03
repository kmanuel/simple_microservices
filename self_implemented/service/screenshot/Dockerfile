FROM golang

RUN apt-get update
RUN apt-get install libfontconfig -y

RUN mkdir --parent $GOPATH/src/github.com/kmanuel/simple_microservices/self_implemented/service/screenshot
WORKDIR $GOPATH/src/github.com/kmanuel/simple_microservices/self_implemented/service/screenshot

COPY . .

RUN go get ./
RUN go build

ENV FAKTORY_URL tcp://faktory:7419

EXPOSE 8080

ENTRYPOINT ["go", "run", "main.go"]
