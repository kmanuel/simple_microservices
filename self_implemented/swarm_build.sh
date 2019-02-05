#!/bin/sh

go build -ldflags "-linkmode external -extldflags -static" -o service/crop/app service/crop/main.go
cd service/crop
docker build -t swarm_crop -f PrecompDockerfile .
cd ../..


go build -ldflags "-linkmode external -extldflags -static" -o service/most_significant_image/app service/most_significant_image/main.go
cd service/most_significant_image
docker build -t swarm_most_significant_image -f PrecompDockerfile .
cd ../..


go build -ldflags "-linkmode external -extldflags -static" -o service/optimization/app service/optimization/main.go
cd service/optimization
docker build -t swarm_optimization -f PrecompDockerfile .
cd ../..


go build -ldflags "-linkmode external -extldflags -static" -o service/portrait/app service/portrait/main.go
cd service/portrait
docker build -t swarm_portrait -f PrecompDockerfile .
cd ../..


go build -ldflags "-linkmode external -extldflags -static" -o service/screenshot/app service/screenshot/main.go
cd service/screenshot
docker build -t swarm_screenshot -f PrecompDockerfile .
cd ../..


go build -ldflags "-linkmode external -extldflags -static" -o gateway/app gateway/main.go
cd gateway
docker build -t swarm_gateway -f PrecompDockerfile .
cd ..


go build -ldflags "-linkmode external -extldflags -static" -o service/request_service/app service/request_service/main.go
cd service/request_service
docker build -t swarm_request_service -f PrecompDockerfile .
cd ../..
