FROM golang

RUN mkdir $GOPATH/src/app
WORKDIR $GOPATH/src/app

COPY . .

RUN go get ./
RUN go build

ENV FAKTORY_URL tcp://faktory:7419

ENTRYPOINT ["go", "run", "main.go"]