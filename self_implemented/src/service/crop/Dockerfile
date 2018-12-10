FROM golang

RUN mkdir /app
WORKDIR /app

COPY . .

RUN go get ./
RUN go build
