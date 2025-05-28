package confirm

type Config struct {
	Lable string

	TrueOption    string
	FalseOption   string
	CaseSensative bool

	AllowDefuatValue bool // Allows the user to not enter anything and the defulat value will be used
	DefualtValue     bool // What the defualt value will be, true or false

	SelectKeys []string // Keys to confirm a choice
	AbortKeys  []string // Keys to confirm a choice
}

func DefualtConfig() Config {
	return Config{
		Lable:            "Select option",
		TrueOption:       "y",
		FalseOption:      "n",
		CaseSensative:    false,
		AllowDefuatValue: true,
		DefualtValue:     true,
		SelectKeys:       []string{"enter"},
		AbortKeys:        []string{"ctrl+c"},
	}
}
