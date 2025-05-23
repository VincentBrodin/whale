package list

import (
	"fmt"
	"strings"

	"github.com/VincentBrodin/whale/codes"
)

type Config struct {
	Prompt string

	AllowSearch bool
	ViewSize    int // The max amount of items to be shown at one time

	UpKeys     []string
	DownKeys   []string
	SelectKeys []string
	SearchKeys []string
	ExitKeys   []string
	AbortKeys  []string

	RenderItem func(item string, selected bool, config Config) string
	RenderInfo func(index, size int, config Config) string
	// Start is the string before the cursor and end is after the cursor
	RenderSearch func(start, end string, config Config) string
}

func DefualtConfig() Config {
	return Config{
		Prompt:      "Select option",
		AllowSearch: true,
		ViewSize:    4,

		UpKeys:     []string{"arrowup", "k"},
		DownKeys:   []string{"arrowdown", "j"},
		SelectKeys: []string{"enter"},
		SearchKeys: []string{"/"},
		ExitKeys:   []string{"esc"},
		AbortKeys:  []string{"ctrl+c"},

		RenderItem: func(item string, selected bool, config Config) string {
			if selected {
				return fmt.Sprintf("  > %s", item)
			}
			return fmt.Sprintf("%s    %s", codes.Muted, item)
		},
		RenderInfo: func(index, size int, config Config) string {
			keys := strings.Builder{}
			if _, err := keys.WriteString("up:"); err != nil {
				return "error"
			}
			for _, key := range config.UpKeys {
				if _, err := keys.WriteString(keyToSymbol(key)); err != nil {
					return "error"
				}
			}
			if _, err := keys.WriteString(" | down:"); err != nil {
				return "error"
			}
			for _, key := range config.UpKeys {
				if _, err := keys.WriteString(keyToSymbol(key)); err != nil {
					return "error"
				}
			}

			if _, err := keys.WriteString(" | search:"); err != nil {
				return "error"
			}
			for _, key := range config.SearchKeys {
				if _, err := keys.WriteString(keyToSymbol(key)); err != nil {
					return "error"
				}
			}

			if _, err := keys.WriteString(" | select:"); err != nil {
				return "error"
			}
			for _, key := range config.SelectKeys {
				if _, err := keys.WriteString(keyToSymbol(key)); err != nil {
					return "error"
				}
			}

			return fmt.Sprintf("%s%d/%d | %s |", codes.Muted, index, size, keys.String())
		},
		RenderSearch: func(start, end string, config Config) string {
			keys := strings.Builder{}
			if _, err := keys.WriteString("exit: "); err != nil {
				return "error"
			}
			for _, key := range config.ExitKeys {
				if _, err := keys.WriteString(keyToSymbol(key)); err != nil {
					return "error"
				}
			}

			return fmt.Sprintf("%sSearch:%s%s█%s%s | %s |", codes.Muted, codes.Reset, start, end, codes.Muted, keys.String())
		},
	}
}

func keyToSymbol(key string) string {
	switch key {
	case "arrowup":
		return "↑"
	case "arrowdown":
		return "↓"
	case "arrowleft":
		return "←"
	case "arrowright":
		return "→"
	case "enter":
		return "↵"
	default:
		return key
	}
}
