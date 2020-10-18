package colour

var (
	black   = "\033[0;30m"
	red     = "\033[0;31m"
	green   = "\033[0;32m"
	yellow  = "\033[0;33m"
	blue    = "\033[0;34m"
	magenta = "\033[0;35m"
	cyan    = "\033[0;36m"

	b_black   = "\033[1;30m"
	b_red     = "\033[1;31m"
	b_green   = "\033[1;32m"
	b_yellow  = "\033[1;33m"
	b_blue    = "\033[1;34m"
	b_magenta = "\033[1;35m"
	b_cyan    = "\033[1;36m"

	none = "\033[0m"
	bold = "\033[1m"
)

func Bold(s string) string {
	return bold + s + none
}

func Black(s string) string {
	return black + s + none
}

func Red(s string) string {
	return red + s + none
}

func Green(s string) string {
	return green + s + none
}

func Yellow(s string) string {
	return yellow + s + none
}

func Blue(s string) string {
	return blue + s + none
}

func Magenta(s string) string {
	return magenta + s + none
}

func Cyan(s string) string {
	return cyan + s + none
}

func BBlack(s string) string {
	return b_black + s + none
}

func BRed(s string) string {
	return b_red + s + none
}

func BGreen(s string) string {
	return b_green + s + none
}

func BYellow(s string) string {
	return b_yellow + s + none
}

func BBlue(s string) string {
	return b_blue + s + none
}

func BMagenta(s string) string {
	return b_magenta + s + none
}

func BCyan(s string) string {
	return b_cyan + s + none
}
