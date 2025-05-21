package screen

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/VincentBrodin/whale/codes"
	"golang.org/x/term"
)

type Screen struct {
	in     io.Reader
	out    io.Writer
	reader *bufio.Reader
}

func New() *Screen {
	in := os.Stdin
	s := &Screen{
		in:     in,
		out:    os.Stdout,
		reader: bufio.NewReader(in),
	}
	return s
}

// Prints output to the screen
func (s *Screen) Printf(format string, a ...any) error {
	str := fmt.Sprintf(format, a...)
	_, err := s.out.Write([]byte(str))
	return err
}

// Writes output to the screen without updating cursor position, great for styling
func (s *Screen) Print(str string) error {
	_, err := s.out.Write([]byte(str))
	return err
}

func (s *Screen) Read() (string, error) {
	return s.reader.ReadString('\n')
}

func (s *Screen) ReadKey() (string, error) {
	restore, err := enableRawMode()
	if err != nil {
		return "", err
	}
	defer restore()

	buf := make([]byte, 3)
	n, err := s.in.Read(buf)
	if err != nil {
		return "", err
	}

	if n == 1 {
		if n == 1 {
			switch buf[0] {
			case 13:
				return "enter", nil
			case 3:
				return "ctrl+c", nil
			case 27:
				return "esc", nil
			default:
				return string(buf[0]), nil
			}
		}
	} else if n == 3 && buf[0] == 27 && buf[1] == 91 {
		switch buf[2] {
		case 65:
			return "uparrow", nil
		case 66:
			return "downarrow", nil
		case 67:
			return "rightarrow", nil
		case 68:
			return "leftarrow", nil
		}
	}
	return "", fmt.Errorf("Unknown sequence: %v\n", buf[:n])
}

func (s *Screen) Clear() error {
	return s.Print(codes.ClearCode)
}

func (s *Screen) SetPos(row, col int) error {
	return s.Printf(codes.MoveToCode, row, col)
}

// Ask for the cursor position
func (s *Screen) GetPos() (row, col int, err error) {
	restore, err := enableRawMode()
	if err != nil {
		return 0, 0, err
	}
	defer restore()

	if err := s.Print("\033[6n"); err != nil {
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

func (s *Screen) SetCode(style string) error {
	if err := s.ResetCode(); err != nil {
		return err
	}
	return s.Print(style)
}

func (s *Screen) ResetCode() error {
	return s.Print(codes.ResetCode)
}
