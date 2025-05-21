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

	index  int
	winPos int

	startPos int
	endPos   int

	row int
}

func NewList(items []string) *List {
	return &List{
		Screen:   screen.New(),
		Items:    items,
		View:     4,
		index:    0,
		winPos:   0,
		startPos: 0,
		endPos:   0,
		row:      0,
	}
}

func (l *List) Prompt(prompt string) (int, error) {
	if err := l.Screen.Printf("%s\n", prompt); err != nil {
		return -1, err
	}

	if err := l.render(true); err != nil {
		return -1, err
	}

	row, _, err := l.Screen.GetPos()
	if err != nil {
		return -1, err
	}

	l.endPos = row - 1
	l.startPos = row - min(l.View, len(l.Items))
	l.row = l.startPos
	if err := l.listen(); err != nil {
		return -1, err
	}
	return l.index, nil
}

func (l *List) render(init bool) error {
	if !init {
		l.Screen.SetPos(l.startPos-1, 1)
	}

	// Clear the line
	if err := l.Screen.SetCode(codes.ClearLineCode); err != nil {
		return err
	}
	if err := l.Screen.SetCode(codes.MutedCode); err != nil {
		return err
	}
	if err := l.Screen.Printf("%d/%d | up:↑k down:↓j  select:↵\n", l.index+1, len(l.Items)); err != nil {
		return err
	}

	size := min(l.View, len(l.Items))

	// This code updates the window to follow the cursor
	d := l.index - l.winPos
	if d >= size {
		l.winPos++
	} else if d < 0 {
		l.winPos--
	}
	l.winPos = max(l.winPos, 0)
	l.winPos = min(l.winPos, len(l.Items)-1)

	// Loop through the window
	for _i, item := range l.Items[l.winPos : l.winPos+size] {
		i := _i + l.winPos // Get the real index

		// Clear the line
		if err := l.Screen.SetCode(codes.ClearLineCode); err != nil {
			return err
		}

		// If we are on the index we want to write with full color
		if i == l.index {
			if err := l.Screen.ResetCode(); err != nil {
				return err
			}
			if err := l.Screen.Printf("  → %s\n", item); err != nil {
				return err
			}
		} else { // Else we draw it muted
			if err := l.Screen.SetCode(codes.MutedCode); err != nil {
				return err
			}
			if err := l.Screen.Printf("    %s\n", item); err != nil {
				return err
			}
		}
	}
	if err := l.Screen.ResetCode(); err != nil {
		return err
	}
	return nil
}

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
