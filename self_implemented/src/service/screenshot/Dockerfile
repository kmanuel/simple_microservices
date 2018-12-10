FROM golang
RUN apt-get update
RUN go get github.com/gorilla/mux
RUN apt-get install libfontconfig -y
RUN mkdir /app
RUN mkdir /app/output
ADD . /app
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]
