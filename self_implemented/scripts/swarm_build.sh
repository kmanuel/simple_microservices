#!/bin/sh

go build -ldflags "-linkmode external -extldflags -static" -o service/crop/app service/crop/main.go
cd service/crop
docker build -t swarm_crop .
cd ../..


go build -ldflags "-linkmode external -extldflags -static" -o service/most_significant_image/app service/most_significant_image/main.go
cd service/most_significant_image
docker build -t swarm_most_significant_image .
cd ../..


go build -ldflags "-linkmode external -extldflags -static" -o service/optimization/app service/optimization/main.go
cd service/optimization
docker build -t swarm_optimization .
cd ../..


go build -ldflags "-linkmode external -extldflags -static" -o service/portrait/app service/portrait/main.go
cd service/portrait
docker build -t swarm_portrait .
cd ../..


go build -ldflags "-linkmode external -extldflags -static" -o service/screenshot/app service/screenshot/main.go
cd service/screenshot
docker build -t swarm_screenshot .
cd ../..

go build -ldflags "-linkmode external -extldflags -static" -o service/screenshot/app service/screenshot/main.go
cd service/scale_service
docker build -t swarm_scale .
cd ../..



go build -ldflags "-linkmode external -extldflags -static" -o service/gateway/app service/gateway/main.go
cd service/gateway
docker build -t swarm_gateway .
cd ../..
