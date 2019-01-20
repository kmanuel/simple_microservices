FROM golang

RUN mkdir --parent $GOPATH/src/github.com/kmanuel/simple_microservices/self_implemented/src/service/crop
WORKDIR $GOPATH/src/github.com/kmanuel/simple_microservices/self_implemented/src/service/crop

COPY . .

RUN go get ./

ENV FAKTORY_URL tcp://faktory:7419

EXPOSE 8080

ENTRYPOINT ["go", "run", "main.go"]
