package screen

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/VincentBrodin/whale/codes"
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

	// Read up to 8 bytes for potential escape sequences
	buf := make([]byte, 8)
	n, err := s.in.Read(buf)
	if err != nil {
		return "", err
	}
	buf = buf[:n]

	// Single-byte control and printable characters
	if n == 1 {
		b := buf[0]
		switch {
		case b == 13:
			return "enter", nil
		case b == 27:
			return "esc", nil
		case b == 3:
			return "ctrl+c", nil
		case b == 8 || b == 127:
			return "backspace", nil
		case b >= 1 && b <= 26:
			return fmt.Sprintf("ctrl+%c", 'a'+b-1), nil
		}
	} else if n == 3 { // Arrows
		a := buf[0]
		b := buf[1]
		c := buf[2]
		if a == 27 && b == 91 {
			switch c {
			case 65:
				return "arrowup", nil
			case 66:
				return "arrowdown", nil

			case 67:
				return "arrowright", nil

			case 68:
				return "arrowleft", nil
			}
		}
	}

	// UTF-8 multibyte rune
	if utf8.FullRune(buf) {
		r, _ := utf8.DecodeRune(buf)
		return string(r), nil
	}

	return "", fmt.Errorf("unknown key sequence: %v", buf)
}

func (s *Screen) Clear() error {
	return s.Print(codes.Clear)
}

func (s *Screen) SetPos(row, col int) error {
	return s.Printf(codes.MoveTo, row, col)
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
