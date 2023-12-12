package main

import (
	"SOMAS2023/internal/common/utils"
	"SOMAS2023/internal/server"
	"fmt"
)

func main() {
	fmt.Print(int(utils.Leadership))
	fmt.Println("Hello Agents")
	s := server.Initialize(3)
	s.UpdateGameStates()
	s.Start()

}
