package valj

import (
	"github.com/VincentBrodin/valj/stage"
)

type List struct {
	Stage *stage.Stage
	Items []string
	Size  int

	start   int
	end     int
	window  int
	index   int
	size    int
}

func (l *List) Prompt(prompt string) (int, error) {
	l.Stage.SetStyle(stage.HideCode)
	l.size = min(l.Size, len(l.Items))
	l.index = 0
	l.window = 0

	l.Stage.Write(prompt + "\n")
	l.Stage.ResetStyle()
	l.Stage.SetStyle(stage.MutedStyle)
	l.Stage.Writef("%d/%d use ↑↓ or jk to move, ↵ to select\n", l.index+1, len(l.Items))

	l.draw()

	finalRow, _, err := l.Stage.GetPos()
	if err != nil {
		return -1, err
	}


	l.start = finalRow - l.size
	l.end = finalRow
	l.Stage.MoveTo(l.start, 1)

	// l.end = l.Stage.Row
	// l.start = l.end - len(l.Items)

	// Reset to the top
	l.Stage.MoveTo(l.start, 1)
	// Run the selection code on the first
	l.move(0)

	for {
		l.Stage.ResetStyle()
		key, err := l.Stage.ReadKey()
		if err != nil {
			continue
		}
		if key == "j" || key == "downarrow" {
			l.move(1)
		} else if key == "k" || key == "uparrow" {
			l.move(-1)
		} else if key == "q" {
			break
		} else if key == "enter" {
			l.Stage.ResetStyle()
			l.Stage.SetStyle(stage.ShowCode)
			l.Stage.MoveTo(l.end, 1)
			return l.index, nil
		}
	}
	l.Stage.ResetStyle()
	l.Stage.SetStyle(stage.ShowCode)
	l.Stage.MoveTo(l.end, 1)
	return -1, nil
}

func (l *List) draw() {
	l.Stage.ResetStyle()
	l.Stage.SetStyle(stage.MutedStyle)
	for i, item := range l.Items[l.window : l.size+l.window] {
		l.Stage.SetStyle(stage.ClearLineCode)
		l.Stage.Writef("%d - %s\n", i+1+l.window, item)
	}
}

func (l *List) update() {
	l.Stage.MoveTo(l.start-1, 1)
	l.Stage.ResetStyle()
	l.Stage.SetStyle(stage.MutedStyle)
	l.Stage.Writef("%d/%d\n", l.index+1, len(l.Items))
}

func (l *List) move(dir int) {
	row := l.Stage.Row
	l.Stage.ResetStyle()
	l.Stage.SetStyle(stage.MutedStyle)
	l.Stage.Writef("%d - %s", l.index+1, l.Items[l.index])

	if l.index+dir < 0 {
		l.Stage.MoveTo(l.end-1, 1)
		l.index = len(l.Items) - 1
		l.window = len(l.Items) - l.size - 1
		row = l.Stage.Row
		l.Stage.MoveTo(l.start, 1)
		l.draw()
		l.Stage.MoveTo(row+1, 1)
	} else if l.index+dir >= len(l.Items) {
		l.Stage.MoveTo(l.start, 1)
		l.index = 0
		l.window = 0
		row = l.Stage.Row
		l.Stage.MoveTo(l.start, 1)
		l.draw()
		l.Stage.MoveTo(row, 1)

	} else {
		l.Stage.MoveTo(row+dir, 1)
		l.index += dir
	}

	// Overflow
	d := l.index - l.window
	if d >= l.size {
		l.window++
		row = l.Stage.Row
		l.Stage.MoveTo(l.start, 1)
		l.draw()
		l.Stage.MoveTo(row-1, 1)
	} else if d < 0 {
		l.window--
		row = l.Stage.Row
		l.Stage.MoveTo(l.start, 1)
		l.draw()
		l.Stage.MoveTo(row+1, 1)
	}

	row = l.Stage.Row
	l.Stage.ResetStyle()
	l.Stage.SetStyle(stage.BoldStyle)
	l.Stage.Writef("%d - %s", l.index+1, l.Items[l.index])
	l.update()
	l.Stage.MoveTo(row, 1)
}
