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

RUN mkdir -p /go/src/github.com/kmanuel/simple_microservices/self_implemented/service/gateway/src
WORKDIR /go/src/github.com/kmanuel/simple_microservices/self_implemented/service/gateway/src
COPY src/ .
RUN go build -ldflags "-linkmode external -extldflags -static" -o app main.go && \
    chmod +x app



FROM alpine:3.8
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /go/src/github.com/kmanuel/simple_microservices/self_implemented/service/gateway/src/app .

ENV FAKTORY_URL tcp://faktory:7419
EXPOSE 8080

ENTRYPOINT ["/app/app"]
