package confirm

import (
	"fmt"

	"github.com/VincentBrodin/whale/codes"
)

type Config struct {
	Lable string // The text at the top

	TrueOption    string // Most commanly y
	FalseOption   string // Most commanly n
	CaseSensative bool

	AllowDefuatValue bool // Allows the user to not enter anything and the defulat value will be used
	DefualtValue     bool // What the defualt value will be, true or false

	SelectKeys []string // Keys to confirm a choice
	AbortKeys  []string // Keys to confirm a choice

	RenderLable func(config Config) string
}

func DefualtConfig() Config {
	return Config{
		Lable: "Select option",

		TrueOption:    "y",
		FalseOption:   "n",
		CaseSensative: false,

		AllowDefuatValue: true,
		DefualtValue:     true,

		SelectKeys: []string{"enter"},
		AbortKeys:  []string{"ctrl+c"},

		RenderLable: func(config Config) string {
			options := ""
			if config.AllowDefuatValue {
				if config.DefualtValue {
					options = fmt.Sprintf("%s[%s%s%s/%s]", codes.Reset, codes.Success, config.TrueOption, codes.Reset, config.FalseOption)
				} else {
					options = fmt.Sprintf("%s[%s/%s%s%s]", codes.Reset, config.TrueOption, codes.Success, config.FalseOption, codes.Reset)
				}
			} else {
				options = fmt.Sprintf("%s[%s/%s]", codes.Reset, config.TrueOption, config.FalseOption)
			}
			return fmt.Sprintf("%s %s: ", config.Lable, options)

		},
	}
}
