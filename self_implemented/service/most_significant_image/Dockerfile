FROM golang

RUN mkdir --parent $GOPATH/src/github.com/kmanuel/simple_microservices/self_implemented/src/service/most_significant_image
WORKDIR $GOPATH/src/github.com/kmanuel/simple_microservices/self_implemented/src/service/most_significant_image

COPY . .

RUN go get ./
RUN go build

ENV FAKTORY_URL tcp://faktory:7419

ENTRYPOINT ["go", "run", "main.go"]