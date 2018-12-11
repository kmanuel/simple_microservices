FROM golang

RUN export PATH="$PATH:$GOPATH/bin"

RUN go get -u github.com/esimov/caire/cmd/caire

RUN mkdir $GOPATH/src/app
WORKDIR $GOPATH/src/app

COPY . .

RUN go get ./
RUN go build

ENTRYPOINT ["go", "run", "main.go"]