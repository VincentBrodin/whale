package list

import (
	"fmt"
	"sort"
	"unicode/utf8"

	"github.com/VincentBrodin/suddig/configs"
	"github.com/VincentBrodin/suddig/matcher"
	"github.com/VincentBrodin/whale/codes"
	"github.com/VincentBrodin/whale/screen"
)

type List struct {
	Items  []string // The items
	Config Config

	index    int
	winPos   int // The start of the window
	startPos int // The screen position of the first element
	endPos   int // The screen position of the end of the list

	searching bool   // True if we are in search mode
	search    string // The search string
	insertPos int    // The users cursor position in the string
	results   []int  // Maps the Items positions to there screen position

	screen  *screen.Screen   // The screen
	matcher *matcher.Matcher // The fuzzy matcher, look at github.com/VincentBrodin/suddig for examples of how to customize it
}

// Creates a new list with defualt configuration
func New(config Config) *List {
	return &List{
		Config:  config,
		matcher: matcher.New(configs.DamerauLevenshtein()),
		screen:  screen.New(),
	}
}

// Shows the list and listens for the users to give an answer, returns the index of the answer.
// if an error occurs the answer will be -1
func (l *List) Prompt(items []string) (int, error) {
	if len(items) == 0 {
		return -1, fmt.Errorf("No items to select\n")
	}

	defer func() {
		_ = l.screen.Print(codes.ShowCursor) // We just show the cursor, we dont care if it fails
	}()

	if err := l.screen.Printf("%s%s\n", codes.HideCursor, l.Config.Prompt); err != nil {
		return -1, err
	}

	// Set the initial results state
	l.Items = items
	l.results = make([]int, len(items))
	for i := range items {
		l.results[i] = i
	}

	if err := l.render(true); err != nil {
		return -1, err
	}

	row, _, err := l.screen.GetPos()
	if err != nil {
		return -1, err
	}

	l.endPos = row - 1
	l.startPos = row - min(l.Config.ViewSize, len(l.Items))
	if err := l.listen(); err != nil {
		return -1, err
	}
	return l.results[l.index], nil
}

// Handels rendering and re-rendering the list
func (l *List) render(init bool) error {
	if !init {
		l.screen.SetPos(l.startPos-1, 1)
	}

	if l.searching {
		start := string([]rune(l.search)[:l.insertPos])
		end := string([]rune(l.search)[l.insertPos:])
		search := l.Config.RenderSearch(start, end, l.Config)
		if err := l.screen.Printf("%s%s%s\n", codes.Reset, codes.ClearLine, search); err != nil {
			return err
		}
	} else {
		helper := l.Config.RenderInfo(l.index+1, len(l.Items), l.Config)
		if err := l.screen.Printf("%s%s%s\n", codes.Reset, codes.ClearLine, helper); err != nil {
			return err
		}
	}

	l.adjustWindow()

	// Loop through the window
	size := min(l.Config.ViewSize, len(l.Items))
	for _i := range l.Items[l.winPos : l.winPos+size] {
		i := _i + l.winPos // Get the real index
		line := l.Config.RenderItem(l.Items[l.results[i]], l.index == i, l.Config)
		if err := l.screen.Printf("%s%s%s\n", codes.Reset, codes.ClearLine, line); err != nil {
			return err
		}
	}

	return l.screen.Print(codes.Reset)
}

// Listens for user input and response accordingly
func (l *List) listen() error {
	for {
		key, err := l.screen.ReadKey()
		if err != nil {
			return err
		}
		if key == "ctrl+c" { // Abort
			return fmt.Errorf("User exited program")
		} else if key == "enter" { // Confirm
			return nil
		} else if l.searching { // Search
			l.updateSearch(key)
		} else { // Normal mode
			l.updateMove(key)
		}
		if err := l.render(false); err != nil {
			return err
		}
	}
}

// Updates the view window if needed
func (l *List) adjustWindow() {
	size := min(l.Config.ViewSize, len(l.Items))

	d := l.index - l.winPos
	if d >= size {
		l.winPos++
	} else if d < 0 {
		l.winPos--
	}

	l.winPos = max(l.winPos, 0)
	l.winPos = min(l.winPos, len(l.Items)-1)
}

func (l *List) updateMove(key string) {
	switch key {
	// Scroll down
	case "j", "arrowdown":
		l.index++
		if l.index >= len(l.Items) {
			l.index = 0
			l.winPos = 0
		}
		break
	// Scroll up
	case "k", "arrowup":
		l.index--
		if l.index < 0 {
			l.index = len(l.Items) - 1
			l.winPos = l.index - min(l.Config.ViewSize, len(l.Items))
		}
		break
	// Search
	case "/":
		if !l.Config.AllowSearch {
			break
		}
		l.searching = true
		l.insertPos = 0
		l.search = ""
		l.index = 0
		l.winPos = 0
		for i := range l.Items {
			l.results[i] = i
		}
		l.updateSearchResults()
		break
	case "esc":
		l.searching = false
		break
	}
}

func (l *List) updateSearch(key string) {
	size := utf8.RuneCountInString(l.search)
	if key == "esc" {
		l.searching = false
	} else if utf8.RuneCountInString(key) == 1 {
		r := []rune(l.search)[:l.insertPos]
		r = append(r, []rune(key)...)
		r = append(r, []rune(l.search)[l.insertPos:]...)
		l.search = string(r)
		l.insertPos++
		l.updateSearchResults()
	} else {
		switch key {
		case "backspace":
			if l.insertPos >= 1 {
				r := []rune(l.search)[:l.insertPos-1]
				r = append(r, []rune(l.search)[l.insertPos:]...)
				l.search = string(r) // Remove the last rune
				l.insertPos--
				l.updateSearchResults()
			}
			break
		case "arrowleft":
			l.insertPos--
			break
		case "arrowright":
			l.insertPos++
			break
		}
		l.insertPos = min(l.insertPos, size)
		l.insertPos = max(l.insertPos, 0)
	}
}

func (l *List) updateSearchResults() {
	result := l.matcher.ParallelRank(l.search, l.Items)
	sort.Slice(l.results, func(i, j int) bool {
		return result[l.results[i]] > result[l.results[j]]
	})
}
