FROM golang

RUN mkdir --parent $GOPATH/src/github.com/kmanuel/simple_microservices/go_kit/gateway
WORKDIR $GOPATH/src/github.com/kmanuel/simple_microservices/go_kit/gateway

COPY . .

RUN go get ./
RUN go build

ENV FAKTORY_URL tcp://faktory:7419

ENTRYPOINT ["go", "run", "gateway.go"]
