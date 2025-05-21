package main

import (
	"fmt"

	"github.com/VincentBrodin/valj"
)

func main() {
	items := []string {"Test", "Hello", "Cool","World","Golang", "House","Boat","Other stuff"}

	v, err := valj.New()
	if err != nil {
		return
	}

	list := v.NewList(items)
	list.Size = 4

	i, _ := list.Prompt("Select something")
	fmt.Println(items[i])
}

