FROM golang

RUN mkdir --parent $GOPATH/src/github.com/kmanuel/simple_microservices/self_implemented/gateway
WORKDIR $GOPATH/src/github.com/kmanuel/simple_microservices/self_implemented/gateway

COPY . .

RUN go get ./
RUN go build

ENV FAKTORY_URL tcp://faktory:7419

ENTRYPOINT ["go", "run", "main.go"]