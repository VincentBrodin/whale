package codes

const (
	ResetCode      = "\033[0m"
	BoldStyle      = "\033[1m"
	UnderlineStyle = "\033[4m"

	ErrorCode   = "\033[31m"
	WarnCode    = "\033[33m"
	SuccessCode = "\033[32m"
	InfoCode    = "\033[36m"
	HintCode    = "\033[34m"
	MutedCode   = "\033[90m"

	LinkCode = UnderlineStyle + InfoCode

	ClearCode     = "\033[2J\033[H"
	ClearLineCode = "\033[2K"

	UpCode       = "\033[%dA"
	DownCode     = "\033[%dB"
	RightCode    = "\033[%dC"
	LeftCode     = "\033[%dD"
	MoveToCode   = "\033[%d;%dH"
	SavePosCode  = "\033[s"
	ResetPosCode = "\033[u"

	HideCursorCode = "\033[?25l"
	ShowCursorCode = "\033[?25h"
)
