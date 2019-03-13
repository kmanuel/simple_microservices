#!/bin/sh
go build -ldflags "-linkmode external -extldflags -static" -o service/crop/app service/crop/src/main.go
go build -ldflags "-linkmode external -extldflags -static" -o service/most_significant_image/app service/most_significant_image/src/main.go
go build -ldflags "-linkmode external -extldflags -static" -o service/optimization/app service/optimization/src/main.go
go build -ldflags "-linkmode external -extldflags -static" -o service/portrait/app service/portrait/src/main.go
go build -ldflags "-linkmode external -extldflags -static" -o service/screenshot/app service/screenshot/src/main.go
go build -ldflags "-linkmode external -extldflags -static" -o service/gateway/app service/gateway/src/main.go
