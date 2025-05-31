package list

import (
	"fmt"
	"strings"

	"github.com/VincentBrodin/whale/codes"
)

type Config struct {
	Lable string // Text displayed at the top

	AllowSearch bool // Enables search mode
	ViewSize    int  // Max number of items to display at once

	UpKeys         []string // Keys to move up
	DownKeys       []string // Keys to move down
	SelectKeys     []string // Keys to confirm a choice
	SearchKeys     []string // Keys to enter search mode
	ExitSearchKeys []string // Keys to exit search mode
	AbortKeys      []string // Keys to cancel/abort the prompt

	// Custom render logic
	RenderItem         func(item string, selected bool, config Config) string
	RenderInfo         func(index, size int, config Config) string
	RenderSearchPrefix func(config Config) string // This is the logic for the text that goes before the search input
	RenderSearchSuffix func(config Config) string // And this is the text after
}

func DefualtConfig() Config {
	return Config{
		Lable:       "Select option",
		AllowSearch: true,
		ViewSize:    4,

		UpKeys:         []string{"arrowup", "k"},
		DownKeys:       []string{"arrowdown", "j"},
		SelectKeys:     []string{"enter"},
		SearchKeys:     []string{"/"},
		ExitSearchKeys: []string{"esc"},
		AbortKeys:      []string{"ctrl+c"},

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

			if config.AllowSearch {
				if _, err := keys.WriteString(" | search:"); err != nil {
					return "error"
				}
				for _, key := range config.SearchKeys {
					if _, err := keys.WriteString(keyToSymbol(key)); err != nil {
						return "error"
					}
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
		RenderSearchPrefix: func(config Config) string {
			return fmt.Sprintf("%s%sSearch:", codes.Reset, codes.Muted)
		},
		RenderSearchSuffix: func(config Config) string {
			keys := strings.Builder{}
			if _, err := keys.WriteString("exit: "); err != nil {
				return "error"
			}
			for _, key := range config.ExitSearchKeys {
				if _, err := keys.WriteString(keyToSymbol(key)); err != nil {
					return "error"
				}
			}

			return fmt.Sprintf("%s%s | %s |", codes.Reset, codes.Muted, keys.String())
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
