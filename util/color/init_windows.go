package color

import (
	"golang.org/x/sys/windows"
)

func init() {
	var (
		dwMode  uint32 = 0
		hOut, _        = windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
	)
	if windows.GetConsoleMode(hOut, &dwMode) != nil {
		return
	}
	dwMode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING

	// Here is the only place an error would be unexpected.
	if err := windows.SetConsoleMode(hOut, dwMode); err != nil {
		panic(err)
	}
}
