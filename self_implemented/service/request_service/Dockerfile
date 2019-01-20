FROM golang

RUN mkdir --parent $GOPATH/src/github.com/kmanuel/simple_microservices/self_implemented/src/request_service
WORKDIR $GOPATH/src/github.com/kmanuel/simple_microservices/self_implemented/src/request_service

COPY . .

RUN go get ./
RUN go build

ENV FAKTORY_URL tcp://faktory:7419

ENTRYPOINT ["go", "run", "main.go"]