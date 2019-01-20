FROM golang

RUN export PATH="$PATH:$GOPATH/bin"

RUN go get -u github.com/esimov/caire/cmd/caire

RUN mkdir --parent $GOPATH/src/github.com/kmanuel/simple_microservices/self_implemented/src/service/portrait
WORKDIR $GOPATH/src/github.com/kmanuel/simple_microservices/self_implemented/src/service/portrait

COPY . .

RUN go get ./
RUN go build

ENV FAKTORY_URL tcp://faktory:7419

EXPOSE 8080

ENTRYPOINT ["go", "run", "main.go"]