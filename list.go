package whale

import (
	"fmt"

	"github.com/VincentBrodin/whale/codes"
	"github.com/VincentBrodin/whale/screen"
)

type List struct {
	Screen *screen.Screen // The screen
	Items  []string       // The items
	View   int            // The max amount of items to be shown at one time

	RenderLine   func(index int, item string, selected bool) string // How the output line should look
	RenderHelper func(index, size int) string                       // How the helper line should look

	index    int
	winPos   int
	startPos int
	endPos   int
}

// Creates a new list with defualt configuration
func NewList(items []string) *List {
	return &List{
		Screen: screen.New(),
		Items:  items,
		View:   4,

		RenderLine: func(index int, item string, selected bool) string {
			if selected {
				return fmt.Sprintf("  → %s", item)
			}
			return fmt.Sprintf("%s    %s", codes.Muted, item)
		},

		RenderHelper: func(index, size int) string {
			return fmt.Sprintf("%d/%d | up:↑k down:↓j  select:↵", index, size)
		},
	}
}

// Shows the list and listens for the users to give an answer, returns the index of the answer.
// if an error occurs the answer will be -1
func (l *List) Prompt(prompt string) (int, error) {
	if err := l.Screen.Printf("%s%s\n", codes.HideCursor, prompt); err != nil {
		return -1, err
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
	return l.index, nil
}

// Handels rendering and re-rendering the list
func (l *List) render(init bool) error {
	if !init {
		l.Screen.SetPos(l.startPos-1, 1)
	}

	helper := l.RenderHelper(l.index+1, len(l.Items))
	if err := l.Screen.Printf("%s%s%s%s\n",codes.Reset, codes.ClearLine, codes.Muted, helper ); err != nil {
		return err
	}

	l.adjustWindow()

	// Loop through the window
	size := min(l.View, len(l.Items))
	for _i, item := range l.Items[l.winPos : l.winPos+size] {
		i := _i + l.winPos // Get the real index
		line := l.RenderLine(l.index, item, l.index == i)
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

		switch key {
		// Exit
		case "ctrl+c":
			return fmt.Errorf("User exited program")
		// Scroll down
		case "j", "arrowdown":
			l.index++
			if l.index >= len(l.Items) {
				l.index = 0
				l.winPos = 0
			}
			if err := l.render(false); err != nil {
				return err
			}
			break
		// Scroll up
		case "k", "arrowup":
			l.index--
			if l.index < 0 {
				l.index = len(l.Items) - 1
				l.winPos = l.index - min(l.View, len(l.Items))
			}
			if err := l.render(false); err != nil {
				return err
			}
			break
		case "enter":
			return nil
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
