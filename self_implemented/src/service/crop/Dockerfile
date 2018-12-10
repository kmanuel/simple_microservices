FROM golang

RUN mkdir --parent /go/src/github.com/kmanuel/simple_microservices/self_implemented/src/minioconnector
COPY minioconnector /go/src/github.com/kmanuel/simple_microservices/self_implemented/src/minioconnector

RUN mkdir --parent /go/src/github.com/kmanuel/crop

WORKDIR /go/src/github.com/kmanuel/crop

COPY . .

RUN go get ./

EXPOSE 8082

CMD ["go", "run", "main.go"]