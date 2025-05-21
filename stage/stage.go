package stage

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Stage struct {
	Row    int
	Col    int
	in     io.Reader
	out    io.Writer
	reader *bufio.Reader
	buffer []byte
}

func New() (*Stage, error) {
	in := os.Stdin
	s := &Stage{
		in:     in,
		out:    os.Stdout,
		reader: bufio.NewReader(in),
	}

	var err error
	s.Row, s.Col, err = s.GetPos()

	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Stage) Write(str string) error {
	s.updatePos(str)
	_, err := s.out.Write([]byte(str))
	return err
}

func (s *Stage) Writef(format string, a ...any) error {
	str := fmt.Sprintf(format, a...)
	s.updatePos(str)
	_, err := s.out.Write([]byte(str))
	return err
}

func (s *Stage) Read() (string, error) {
	return s.reader.ReadString('\n')
}

func (s *Stage) ReadKey() (string, error) {
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

func (s *Stage) Clear() error {
	s.buffer = s.buffer[:0]
	return s.write(ClearCode)
}

func (s *Stage) MoveTo(row, col int) error {
	s.Row = row
	s.Col = col
	return s.write(fmt.Sprintf(MoveToCode, row, col))
}
