package main 

import (
	"fmt"

	"github.com/VincentBrodin/whale/confirm"
)

func exampleConfirm() {
	c := confirm.New(confirm.DefualtConfig())
	c.Config.Lable = "Do you like cats?"
	res, err := c.Prompt()
	if err != nil {
		panic(err)
	}

	if res {
		fmt.Println("User likes cats")
	} else {
		fmt.Println("User is wrong")
	}
}
