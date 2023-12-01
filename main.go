package main

import (
	"SOMAS2023/internal/server"
	"fmt"
)

func main() {
	fmt.Println("Hello Agents")
	server.Initialize(100).Start()
}
