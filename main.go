package main

import (
	"SOMAS2023/internal/server"
	"fmt"
)

func main() {
	fmt.Println("Hello Agents")
<<<<<<< HEAD
	server.Initialize(500).Start()
=======
	s := server.Initialize(100)
	s.FoundingInstitutions()
	s.UpdateGameStates()
	s.Start()
>>>>>>> dee9ec203d4c60a8651299808ce6a6d04c17b859
}
