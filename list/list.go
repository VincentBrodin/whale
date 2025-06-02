package list

import (
	"fmt"
	"slices"
	"sort"

	"github.com/VincentBrodin/suddig/configs"
	"github.com/VincentBrodin/suddig/matcher"
	"github.com/VincentBrodin/whale/codes"
	"github.com/VincentBrodin/whale/screen"
	"github.com/VincentBrodin/whale/text"
)

type List struct {
	Items  []string // The items
	Config Config

	index    int
	winPos   int // The start of the window
	startPos int // The screen position of the first element
	endPos   int // The screen position of the end of the list

	searching bool  // True if we are in search mode
	results   []int // Maps the Items positions to there screen position

	screen  *screen.Screen   // The screen
	matcher *matcher.Matcher // The fuzzy matcher, look at github.com/VincentBrodin/suddig for examples of how to customize it
	text    *text.Text       // The search text input
}

// Creates a new list with defualt configuration
func New(config Config) *List {
	return &List{
		Config:  config,
		matcher: matcher.New(configs.DamerauLevenshtein()),
		screen:  screen.New(),
		text:    &text.Text{},
	}
}

// Shows the list and listens for the users to give an answer, returns the index of the answer.
// if an error occurs the answer will be -1
func (l *List) Prompt(items []string) (int, error) {
	if len(items) == 0 {
		return -1, fmt.Errorf("No items to select\n")
	}

	defer func() {
		// We just show the cursor, we dont care if it fails
		if err := l.screen.Print(codes.ShowCursor); err != nil {
			return
		}
		// This makes sure that the users cursor will always be at the end of the list when we exit
		if err := l.screen.SetPos(l.endPos, 1); err != nil {
			return
		}
	}()

	if err := l.screen.Printf("%s%s\n", codes.HideCursor, l.Config.Lable); err != nil {
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

	cRow, cCol := 0, 0

	if l.searching {
		if err := l.screen.Printf("%s", codes.ShowCursor); err != nil {
			return err
		}
		prefix := l.Config.RenderSearchPrefix(l.Config)
		if err := l.screen.Printf("%s%s%s%s%s", codes.Reset, codes.ClearLine, prefix, codes.Reset, l.text.Start()); err != nil {
			return err
		}

		row, col, err := l.screen.GetPos()
		if err != nil {
			cRow = l.endPos 
			cCol = 1
		} else {
			cRow = row
			cCol = col
		}

		suffix := l.Config.RenderSearchSuffix(l.Config)
		if err := l.screen.Printf("%s%s%s\n", codes.Reset, l.text.End(), suffix); err != nil {
			return err
		}

	} else {
		if err := l.screen.Printf("%s", codes.HideCursor); err != nil {
			return err
		}
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

	if l.searching {
		if err := l.screen.SetPos(cRow, cCol); err != nil {
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
		if slices.Contains(l.Config.AbortKeys, key) { // Abort
			return fmt.Errorf("User aborted")
		} else if slices.Contains(l.Config.SelectKeys, key) { // Select
			if l.searching {
				l.searching = false
			} else {
				return nil
			}
		} else if slices.Contains(l.Config.SearchKeys, key) && l.Config.AllowSearch { // Search
			l.searching = true
			l.text.Reset()
			l.index = 0
			l.winPos = 0
			for i := range l.Items {
				l.results[i] = i
			}
			l.search()
		} else if l.searching { // Searching
			if slices.Contains(l.Config.ExitSearchKeys, key) {
				l.searching = false
			} else {
				l.text.Update(key)
				l.search()
			}
		} else { // Normal mode
			l.move(key)
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

func (l *List) move(key string) {
	if slices.Contains(l.Config.DownKeys, key) {
		l.index++
		if l.index >= len(l.Items) {
			l.index = 0
			l.winPos = 0
		}
	} else if slices.Contains(l.Config.UpKeys, key) {
		l.index--
		if l.index < 0 {
			l.index = len(l.Items) - 1
			l.winPos = l.index - min(l.Config.ViewSize, len(l.Items))
		}

	}
}

func (l *List) search() {
	result := l.matcher.ParallelRank(l.text.Value, l.Items)
	sort.Slice(l.results, func(i, j int) bool {
		return result[l.results[i]] > result[l.results[j]]
	})
}

func (l *List) Reset() {
	l.text.Reset()
	l.index = 0
	l.winPos = 0
}
