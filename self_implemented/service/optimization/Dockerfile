FROM golang:1.12.0 as builder

RUN apt-get update && apt-get install -y unzip --no-install-recommends && \
    apt-get autoremove -y && apt-get clean -y && \
    wget -O dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && \
    mv dep /usr/bin && \
    chmod a+x /usr/bin/dep

RUN mkdir -p /go/src/github.com
WORKDIR /go/src/github.com

COPY Gopkg.toml Gopkg.lock ./

RUN dep ensure -vendor-only

RUN mv vendor/github.com/* .

RUN mkdir -p /go/src/github.com/kmanuel/simple_microservices/self_implemented/service/optimization/src
WORKDIR /go/src/github.com/kmanuel/simple_microservices/self_implemented/service/optimization/src
COPY src/ .
RUN go build -ldflags "-linkmode external -extldflags -static" -o app main.go && \
    chmod +x app



FROM debian

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

WORKDIR /app

COPY --from=builder /go/src/github.com/kmanuel/simple_microservices/self_implemented/service/optimization/src/app .

EXPOSE 8080

ENTRYPOINT ["/app/app"]
