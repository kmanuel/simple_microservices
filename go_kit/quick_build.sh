#!/bin/sh
go build -ldflags "-linkmode external -extldflags -static" -o service/crop/app service/crop/crop.go
go build -ldflags "-linkmode external -extldflags -static" -o service/most_significant_image/app ./service/most_significant_image/most_significant_image.go
go build -ldflags "-linkmode external -extldflags -static" -o service/optimization/app ./service/optimization/optimization.go
go build -ldflags "-linkmode external -extldflags -static" -o service/portrait/app ./service/portrait/portrait.go
go build -ldflags "-linkmode external -extldflags -static" -o service/screenshot/app ./service/screenshot/screenshot.go
go build -ldflags "-linkmode external -extldflags -static" -o gateway/app ./gateway/gateway.go
go build -ldflags "-linkmode external -extldflags -static" -o service/request_service/app ./service/request_service/request_service.go
