package main

import (
	"SOMAS2023/internal/server"
	"fmt"
)

func main() {
	fmt.Println("Hello Agents")
	s := server.Initialize(100)
	s.FoundingInstitutions()
	s.UpdateGameStates()
	s.Start()
}
