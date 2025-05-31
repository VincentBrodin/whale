package confirm

import (
	"fmt"
	"slices"
	"strings"

	"github.com/VincentBrodin/whale/codes"
	"github.com/VincentBrodin/whale/screen"
	"github.com/VincentBrodin/whale/text"
)

type Confirm struct {
	Config   Config
	text     *text.Text
	screen   *screen.Screen // The screen
	startPos int            // The screen position of the first element
}

func New(config Config) *Confirm {
	return &Confirm{
		Config: config,
		text:   &text.Text{},
		screen: screen.New(),
	}
}

func (c *Confirm) Prompt() (bool, error) {
	defer func() {
		// Moves the users cursor to the end of the screen
		_ = c.screen.Print("\n")
	}()
	return c.prompt()
}

func (c *Confirm) prompt() (bool, error) {
	if err := c.screen.Printf("%s%s%s%s", codes.Reset, c.Config.RenderLable(c.Config), codes.Reset, c.text.Start()); err != nil {
		return false, err
	}

	row, col, err := c.screen.GetPos()
	if err != nil {
		return false, err
	}
	c.startPos = row

	if err := c.screen.Printf("%s%s", codes.Reset, c.text.End()); err != nil {
		return false, err
	}

	if err := c.screen.SetPos(row, col); err != nil {
		return false, err
	}

	for {
		key, err := c.screen.ReadKey()
		if err != nil {
			return false, err
		}
		if slices.Contains(c.Config.AbortKeys, key) {
			return false, fmt.Errorf("User aborted")
		}
		if slices.Contains(c.Config.SelectKeys, key) {
			return c.confirm()
		} else {
			c.text.Update(key)
			c.screen.SetPos(c.startPos, 1)
			if err := c.screen.Printf("%s%s%s%s%s", codes.ClearLine, codes.Reset, c.Config.RenderLable(c.Config), codes.Reset, c.text.Start()); err != nil {
				return false, err
			}

			row, col, err := c.screen.GetPos()
			if err != nil {
				return false, err
			}

			if err := c.screen.Printf("%s%s", codes.Reset, c.text.End()); err != nil {
				return false, err
			}

			if err := c.screen.SetPos(row, col); err != nil {
				return false, err
			}
		}
	}
}

func (c *Confirm) confirm() (bool, error) {
	value := c.text.Value
	if !c.Config.CaseSensative {
		value = strings.ToLower(value)
	}
	if value == c.Config.TrueOption {
		return true, nil
	} else if value == c.Config.FalseOption {
		return false, nil
	} else if c.Config.AllowDefuatValue {
		return c.Config.DefualtValue, nil
	} else {
		if err := c.screen.Printf("%s%s\nInvalid input\n", codes.Reset, codes.Error); err != nil {
			return false, err
		}
		c.text.Reset()
		return c.prompt()
	}
}

func (c *Confirm) Reset() {
	c.text.Reset()
}
