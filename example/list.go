package main

import (
	"fmt"

	"github.com/VincentBrodin/whale/list"
)

func exampleList() {
	foods := []string{
		"Pizza",
		"Sushi",
		"Tacos",
		"Burger",
		"Pasta",
		"Salad",
		"Ramen",
		"Steak",
		"Ice Cream",
		"Falafel",
	}

	l := list.New(list.DefualtConfig())
	l.Config.Lable = "Select your favorit food"
	l.Config.AllowSearch = false
	res, err := l.Prompt(foods)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Your favorit food is: %s\n", foods[res])
}
