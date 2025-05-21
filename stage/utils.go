package stage

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

const (
	ResetCode      = "\033[0m"
	BoldStyle      = "\033[1m"
	UnderlineStyle = "\033[4m"

	ErrorStyle   = "\033[31m"
	WarnStyle    = "\033[33m"
	SuccessStyle = "\033[32m"
	InfoStyle    = "\033[36m"
	HintStyle    = "\033[34m"
	MutedStyle   = "\033[90m"

	LinkStyle = UnderlineStyle + InfoStyle

	ClearCode     = "\033[2J\033[H"
	ClearLineCode = "\033[2K"

	UpCode       = "\033[%dA"
	DownCode     = "\033[%dB"
	RightCode    = "\033[%dC"
	LeftCode     = "\033[%dD"
	MoveToCode   = "\033[%d;%dH"
	SavePosCode  = "\033[s"
	ResetPosCode = "\033[u"

	HideCode = "\033[?25l"
	ShowCode = "\033[?25h"
)

func (s *Stage) updatePos(str string) {
	for i := 0; i < len(str); i++ {
		switch str[i] {
		case '\n':
			s.Row++
			s.Col = 1
		case '\r':
			s.Col = 1
			if i+1 < len(str) && str[i+1] == '\n' {
				i++
			}
		case '\t':
			const tabWidth = 8
			s.Col = ((s.Col-1)/tabWidth+1)*tabWidth + 1
		case '\b':
			if s.Col > 1 {
				s.Col--
			}
		default:
			s.Col++
		}
	}
}

func (s *Stage) write(str string) error {
	_, err := s.out.Write([]byte(str))
	return err
}

// Ask for the cursor position
func (s *Stage) GetPos() (row, col int, err error) {
	restore, err := enableRawMode()
	if err != nil {
		return 0, 0, err
	}
	defer restore()

	if err := s.write("\033[6n"); err != nil {
		return 0, 0, err
	}

	res, err := s.reader.ReadString('R')
	if err != nil {
		return 0, 0, err
	}

	res = strings.TrimPrefix(res, "\033[")
	res = strings.TrimSuffix(res, "R")

	parts := strings.Split(res, ";")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("unexpected cursor position response: %q", res)
	}

	row, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	col, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return row, col, nil
}

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
