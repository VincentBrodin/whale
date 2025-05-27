package main

import (
	"fmt"
	"github.com/VincentBrodin/whale/list"
)

func main() {
	items := []string{
		"Apple",
		"Banana",
		"Cherry",
		"Durian",
		"Elderberry",
		"Fig",
		"Grape",
		"Honeydew",
		"Kiwi",
		"Lemon",
	}

	list := list.New(list.DefualtConfig())
	i, err := list.Prompt(items)
	if err != nil {
		panic(err)
	}
	fmt.Printf("User selected %s\n", items[i])
}
