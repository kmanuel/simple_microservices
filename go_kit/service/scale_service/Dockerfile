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

RUN mkdir -p /go/src/github.com/kmanuel/simple_microservices/self_implemented/service/crop/src/
WORKDIR /go/src/github.com/kmanuel/simple_microservices/self_implemented/service/crop/src/
COPY src/ .
RUN go build -ldflags "-linkmode external -extldflags -static" -o app main.go && \
    chmod +x app


FROM golang

RUN apt-get update

RUN apt-get install \
        apt-transport-https \
        ca-certificates \
        curl \
        gnupg2 \
        software-properties-common -y

RUN curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add -

RUN add-apt-repository \
       "deb [arch=amd64] https://download.docker.com/linux/debian \
       $(lsb_release -cs) \
       stable"

RUN apt-get update

RUN apt-get install docker-ce docker-ce-cli containerd.io -y

WORKDIR /app

COPY --from=builder /go/src/github.com/kmanuel/simple_microservices/self_implemented/service/crop/src/app .

EXPOSE 8085

ENTRYPOINT ["/app/app"]
