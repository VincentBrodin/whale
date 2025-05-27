## Whale 🐋

**Whale** is a lightweight, pluggable Terminal UI for Go.
Named after the Swedish word **“val”** (*choice* / *whale*), it helps you build interactive prompts with lists, confirmations, fuzzy search, and more—effortlessly.

* [x] List
* [x] List (Searchable)
* [ ] Confirmation
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
	Prompt string // Text displayed at the top

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

<!-- <p align="center"> -->
<!--   <img src="https://github.com/VincentBrodin/whale/assets/demo/list.gif" width="600" alt="Whale List Demo"> -->
<!-- </p> -->
