package main

import (
	"fmt"
	"github.com/google/uuid"
)

func main() {
	println("started")
	uuidStr := uuid.New().String()
	fmt.Println(uuidStr)
}