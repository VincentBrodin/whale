package main

import (
	"fmt"

	"github.com/VincentBrodin/whale/confirm"
	// "github.com/VincentBrodin/whale/list"
)

func main() {
	// items := []string{
	// 	"Apple",
	// 	"Banana",
	// 	"Cherry",
	// 	"Durian",
	// 	"Elderberry",
	// 	"Fig",
	// 	"Grape",
	// 	"Honeydew",
	// 	"Kiwi",
	// 	"Lemon",
	// }

	// l := list.New(list.DefualtConfig())

	c := confirm.New(confirm.DefualtConfig())
	// c.Config.AllowDefuatValue = false
	// c.Config.DefualtValue = false

	for {
		// l.Reset()
		c.Reset()

		// i, err := l.Prompt(items)
		// if err != nil {
		// 	panic(err)
		// }

		c.Config.Lable = fmt.Sprintf("Are you sure?")
		res, err := c.Prompt()
		if err != nil {
			panic(err)
		}

		if res {
			break
		}
	}
}
