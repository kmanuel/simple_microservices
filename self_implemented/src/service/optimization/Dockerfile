FROM golang

RUN apt-get update && \
    apt-get install curl -y

RUN apt-get install ruby -y && \
    gem install image_optim && \
    gem install image_optim_pack && \
    apt-get install -y libjpeg-turbo-progs advancecomp gifsicle jhead jpegoptim optipng pngcrush pngquant && \
    apt-get install build-essential -y && \
    apt-get install libjpeg-dev -y

WORKDIR /tmp
COPY lib .
RUN cd ./jpegoptim-1.4.6 && ./configure && make && make install




RUN mkdir $GOPATH/src/app
WORKDIR $GOPATH/src/app

COPY . .

RUN go get ./
RUN go build

ENV FAKTORY_URL tcp://faktory:7419

ENTRYPOINT ["go", "run", "main.go"]]
