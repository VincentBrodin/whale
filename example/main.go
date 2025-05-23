package main

import (
	"fmt"
	"github.com/VincentBrodin/whale/list"
)

func main() {
	items := []string{"Test","tuff", "Hello", "Cool", "World", "Golang", "House", "Boat", "Other stuff"}

	list := list.New(list.DefualtConfig())
	i, err := list.Prompt(items)
	if err != nil {
		panic(err)
	}
	fmt.Println(items)
	fmt.Println(items[i])
}
