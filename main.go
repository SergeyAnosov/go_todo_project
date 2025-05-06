package main

import (
	"fmt"

	"go1f/pkg/server"
)

func main() {
	err := server.Run()
	if err != nil {
		fmt.Println("Завершаем работу")
		fmt.Println(err)
		panic(err)
	}
}
