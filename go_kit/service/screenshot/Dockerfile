FROM golang

RUN export PATH="$PATH:$GOPATH/bin"

RUN mkdir --parent $GOPATH/src/github.com/kmanuel/simple_microservices/go_kit/service/screenshot
WORKDIR $GOPATH/src/github.com/kmanuel/simple_microservices/go_kit/service/screenshot

COPY . .

RUN go get ./
RUN go build

ENV FAKTORY_URL tcp://faktory:7419

EXPOSE 8080

ENTRYPOINT ["go", "run", "screenshot.go"]
