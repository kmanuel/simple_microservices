#!/bin/sh
go build -ldflags "-linkmode external -extldflags -static" -o service/crop/app service/crop/main.go
go build -ldflags "-linkmode external -extldflags -static" -o service/most_significant_image/app ./service/most_significant_image/main.go
go build -ldflags "-linkmode external -extldflags -static" -o service/optimization/app ./service/optimization/main.go
go build -ldflags "-linkmode external -extldflags -static" -o service/portrait/app ./service/portrait/main.go
go build -ldflags "-linkmode external -extldflags -static" -o service/screenshot/app ./service/screenshot/main.go
go build -ldflags "-linkmode external -extldflags -static" -o gateway/app ./gateway/main.go
go build -ldflags "-linkmode external -extldflags -static" -o service/request_service/app ./service/request_service/main.go
