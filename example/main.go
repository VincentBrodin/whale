package main

import (
	"fmt"

	"github.com/VincentBrodin/whale/list"
)

func main() {
	examples := []string{
		"List",
		"Searchable List",
		"Confirm",
	}

	l := list.New(list.DefualtConfig())
	l.Config.Lable = "Select example"
	res, err := l.Prompt(examples)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Starting %s\n", examples[res])

	switch res {
	case 0:
		exampleList()
		break
	case 1:
		exampleListSearch()
		break
	case 2:
		exampleConfirm()
		break
	default:
		break
	}
}
