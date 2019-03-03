FROM golang

RUN export PATH="$PATH:$GOPATH/bin"

RUN mkdir --parent $GOPATH/src/github.com/kmanuel/simple_microservices/self_implemented/service/portrait
WORKDIR $GOPATH/src/github.com/kmanuel/simple_microservices/self_implemented/service/portrait

COPY . .

RUN go get ./
RUN go build

ENV FAKTORY_URL tcp://faktory:7419

EXPOSE 8080

ENTRYPOINT ["go", "run", "main.go"]
