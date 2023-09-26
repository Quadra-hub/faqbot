package main

import (
	"fmt"

	"github.com/Quadra-hub/go-chatgpt/database"
	"github.com/Quadra-hub/go-chatgpt/router"
)

func main() {
	fmt.Println("Hello World")
	database.Start()
	router.Run()
}
