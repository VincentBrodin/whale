## Whale 🐋

**Whale** is a lightweight, pluggable Terminal UI for Go.
Named after the Swedish word **“val”** (*choice* / *whale*), it helps you build interactive prompts with lists, confirmations, fuzzy search, and more—effortlessly.

* [x] List
* [x] List (Searchable)
* [x] Confirmation
* [ ] Text
* [ ] Multi-choice list

---

## Usage

### List

The `List` component allows users to scroll through a list of choices and select one. It also supports optional fuzzy searching and full keybinding configuration.

#### Basic Usage

```go
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
```

#### Configuration

You can customize how the list behaves and appears using `list.Config`:

```go
type Config struct {
	Lable string // Text displayed at the top

	AllowSearch bool // Enables search mode
	ViewSize    int  // Max number of items to display at once

	UpKeys     []string // Keys to move up
	DownKeys   []string // Keys to move down
	SelectKeys []string // Keys to confirm a choice
	SearchKeys []string // Keys to enter search mode
	ExitKeys   []string // Keys to exit search mode
	AbortKeys  []string // Keys to cancel/abort the prompt

	// Custom render logic
	RenderItem   func(item string, selected bool, config Config) string
	RenderInfo   func(index, size int, config Config) string
	RenderSearch func(start, end string, config Config) string
}
```

Customize these fields to match your application’s needs—for example, adjusting the `Prompt` text, setting `ViewSize` to limit the number of displayed items,
or defining your own keyboard controls.

### Confirm

The `Confirm` component prompts users to answer between 2 choices [y/n].

#### Basic Usage
```go
package main

import (
	"fmt"
	"github.com/VincentBrodin/whale/confirm"
)

func main() {
	c := confirm.New(confirm.DefualtConfig())
	res, err := c.Prompt()
	if err != nil {
		panic(err)
	}
	if res {
		fmt.Println("User said yes")
	} else {
		fmt.Println("User said no")
	}
}
```

#### Configuration

You can customize how the confirm behaves and appears using `confirm.Config`:

```go
type Config struct {
	Lable string // The text at the top

	TrueOption    string // Most commanly y
	FalseOption   string // Most commanly n
	CaseSensative bool

	AllowDefuatValue bool // Allows the user to not enter anything and the defulat value will be used
	DefualtValue     bool // What the defualt value will be, true or false

	SelectKeys []string // Keys to confirm a choice
	AbortKeys  []string // Keys to confirm a choice

    // Custom render logic
	RenderLable func(config Config) string
}
```
<!-- <p align="center"> -->
<!--   <img src="https://github.com/VincentBrodin/whale/assets/demo/list.gif" width="600" alt="Whale List Demo"> -->
<!-- </p> -->
