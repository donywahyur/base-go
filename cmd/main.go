package main

import (
	"fmt"
	"go_base/internal/config"
)

func main() {
	application := config.InitializeApp()

	config.LoadRoute(application)

	err := application.FiberApp.Listen(":8080")

	if err != nil {
		panic(err)
	}

	fmt.Println("Server is running on port 8080")
}
