package main

import (
	"fmt"

	"github.com/VincentBrodin/whale/confirm"
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

	l := list.New(list.DefualtConfig())
	i, err := l.Prompt(items)
	if err != nil {
		panic(err)
	}
	c := confirm.New(confirm.DefualtConfig())
	c.Config.Lable = fmt.Sprintf("Do you want to select %s", items[i])
	c.Config.AllowDefuatValue = false
	res,err := c.Prompt()
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
