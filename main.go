package main

import (
	"SOMAS2023/internal/server"
	"fmt"
)

func main() {
	fmt.Println("Hello Agents")
	s := server.Initialize(10)
	s.UpdateGameStates()
	s.Start()
}
