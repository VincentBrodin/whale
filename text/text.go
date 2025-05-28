package text

import "unicode/utf8"

// Simple text  that
type Text struct {
	Value     string
	insertPos int
}

func (t *Text) Update(key string) {
	size := utf8.RuneCountInString(t.Value)
	if utf8.RuneCountInString(key) == 1 {
		r := []rune(t.Value)[:t.insertPos]
		r = append(r, []rune(key)...)
		r = append(r, []rune(t.Value)[t.insertPos:]...)
		t.Value = string(r)
		t.insertPos++
	} else {
		switch key {
		case "backspace":
			if t.insertPos >= 1 {
				r := []rune(t.Value)[:t.insertPos-1]
				r = append(r, []rune(t.Value)[t.insertPos:]...)
				t.Value = string(r) // Remove the last rune
				t.insertPos--
			}
			break
		case "arrowleft":
			t.insertPos--
			break
		case "arrowright":
			t.insertPos++
			break
		}
		t.insertPos = min(t.insertPos, size)
		t.insertPos = max(t.insertPos, 0)
	}
}

// The string infront of the cursor
func(t *Text) Start() string {
	return string([]rune(t.Value)[:t.insertPos])
}

// The string behind the cursor
func(t *Text) End() string {
	return string([]rune(t.Value)[t.insertPos:])
}

func (t *Text) Reset() {
	t.Value = ""
	t.insertPos = 0
}
