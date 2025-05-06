package main

import (
	"fmt"
	"sergeyanosov/go_todo_project/pkg/server"
)

func main() {
	err := server.Run()
	if err != nil {
		fmt.Println("Завершаем работу")
		fmt.Println(err)
		panic(err)
	}
}
