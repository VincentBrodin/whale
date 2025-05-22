package main

import (
	"fmt"
	"github.com/VincentBrodin/whale"
)

func main() {
	items := []string {"Test", "Hello", "Cool","World","Golang", "House","Boat","Other stuff"}

	list := whale.NewList(items)
	i, err := list.Prompt("Select something")
	if err != nil {
		panic(err)
	}
	fmt.Println(items[i])
}

