package codes

const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Underline = "\033[4m"

	Error   = "\033[31m"
	Warn    = "\033[33m"
	Success = "\033[32m"
	Info    = "\033[36m"
	Hint    = "\033[34m"
	Muted   = "\033[90m"

	Link = Underline + Info

	Clear     = "\033[2J\033[H"
	ClearLine = "\033[2K"

	Up       = "\033[%dA"
	Down     = "\033[%dB"
	Right    = "\033[%dC"
	Left     = "\033[%dD"
	MoveTo   = "\033[%d;%dH"
	SavePos  = "\033[s"
	ResetPos = "\033[u"

	HideCursor = "\033[?25l"
	ShowCursor = "\033[?25h"
)
