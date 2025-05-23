package whale

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
	Screen      *screen.Screen   // The screen
	Items       []string         // The items
	View        int              // The max amount of items to be shown at one time
	AllowSearch bool             // True if list can be searched
	Matcher     *matcher.Matcher // The fuzzy matcher, look at github.com/VincentBrodin/suddig for examples of how to customize it

	RenderLine   func(item string, selected bool) string // How the output line should look
	RenderHelper func(index, size int) string            // How the helper line should look

	index    int // The current index of items the user is on
	winPos   int // The start of the window
	startPos int // The screen position of the first element
	endPos   int // The screen position of the end of the list

	searching bool   // True if we are in search mode
	search    string // The search string
	insertPos int    // The users cursor position in the string
	results   []int  // Maps the Items positions to there screen position
}

// Creates a new list with defualt configuration
func NewList(items []string) *List {
	return &List{
		Screen:      screen.New(),
		Items:       items,
		View:        4,
		AllowSearch: true,
		Matcher:     matcher.New(configs.DamerauLevenshtein()),

		RenderLine: func(item string, selected bool) string {
			if selected {
				return fmt.Sprintf("  → %s", item)
			}
			return fmt.Sprintf("%s    %s", codes.Muted, item)
		},

		RenderHelper: func(index, size int) string {
			return fmt.Sprintf("%d/%d | up:↑k down:↓j select:↵ search:/", index, size)
		},
	}
}

// Shows the list and listens for the users to give an answer, returns the index of the answer.
// if an error occurs the answer will be -1
func (l *List) Prompt(prompt string) (int, error) {
	if err := l.Screen.Printf("%s%s\n", codes.HideCursor, prompt); err != nil {
		return -1, err
	}

	// Set the initial results state
	l.results = make([]int, len(l.Items))
	for i := range l.Items {
		l.results[i] = i
	}

	defer func() {
		_ = l.Screen.Print(codes.ShowCursor) // We just show the cursor, we dont care if it fails
	}()

	if err := l.render(true); err != nil {
		return -1, err
	}

	row, _, err := l.Screen.GetPos()
	if err != nil {
		return -1, err
	}

	l.endPos = row - 1
	l.startPos = row - min(l.View, len(l.Items))
	if err := l.listen(); err != nil {
		return -1, err
	}
	return l.results[l.index], nil
}

// Handels rendering and re-rendering the list
func (l *List) render(init bool) error {
	if !init {
		l.Screen.SetPos(l.startPos-1, 1)
	}

	if l.searching {
		start := string([]rune(l.search)[:l.insertPos])
		end := string([]rune(l.search)[l.insertPos:])
		if err := l.Screen.Printf("%s%sSearch: %s█%s\n", codes.Reset, codes.ClearLine, start, end); err != nil {
			return err
		}
	} else {
		helper := l.RenderHelper(l.index+1, len(l.Items))
		if err := l.Screen.Printf("%s%s%s%s\n", codes.Reset, codes.ClearLine, codes.Muted, helper); err != nil {
			return err
		}
	}

	l.adjustWindow()

	// Loop through the window
	size := min(l.View, len(l.Items))
	for _i := range l.Items[l.winPos : l.winPos+size] {
		i := _i + l.winPos // Get the real index
		line := l.RenderLine(l.Items[l.results[i]], l.index == i)
		if err := l.Screen.Printf("%s%s%s\n", codes.Reset, codes.ClearLine, line); err != nil {
			return err
		}
	}

	return l.Screen.Print(codes.Reset)
}

// Listens for user input and response accordingly
func (l *List) listen() error {
	for {
		key, err := l.Screen.ReadKey()
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
	size := min(l.View, len(l.Items))

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
			l.winPos = l.index - min(l.View, len(l.Items))
		}
		break
	// Search
	case "/":
		if !l.AllowSearch {
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
	result := l.Matcher.ParallelRank(l.search, l.Items)
	sort.Slice(l.results, func(i, j int) bool {
		return result[l.results[i]] > result[l.results[j]]
	})
}
