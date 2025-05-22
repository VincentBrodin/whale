package screen

import (
	"os"

	"golang.org/x/term"
)

func enableRawMode() (func() error, error) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return nil, err
	}
	return func() error {
		return term.Restore(fd, oldState)
	}, nil
}
