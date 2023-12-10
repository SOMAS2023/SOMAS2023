package main

import (
	"SOMAS2023/internal/server"
)

func main() {
	//** fmt.Println("Hello Agents")
	s := server.Initialize(10)
	s.FoundingInstitutions()
	s.UpdateGameStates()
	s.Start()
}
