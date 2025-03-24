package utils

import (
	"fmt"
)

// ANSI color escape codes
const (
	// Colors
	Reset   = "\033[0m"
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	// Bright colors
	BrightBlack   = "\033[90m"
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"

	// Background colors
	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"

	// Text formatting
	Bold      = "\033[1m"
	Underline = "\033[4m"
	Blink     = "\033[5m"
	Reverse   = "\033[7m"
)

// PrintCyan prints text in cyan color and adds a newline
func PrintCyan(format string, a ...interface{}) {
	fmt.Printf(Cyan+format+Reset+"\n", a...)
}

// PrintCyanNoNewline prints text in cyan color without adding a newline
func PrintCyanNoNewline(format string, a ...interface{}) {
	fmt.Printf(Cyan+format+Reset, a...)
}

// PrintYellow prints text in yellow color and adds a newline
func PrintYellow(format string, a ...interface{}) {
	fmt.Printf(Yellow+format+Reset+"\n", a...)
}

// PrintYellowNoNewline prints text in yellow color without adding a newline
func PrintYellowNoNewline(format string, a ...interface{}) {
	fmt.Printf(Yellow+format+Reset, a...)
}

// PrintGreen prints text in green color and adds a newline
func PrintGreen(format string, a ...interface{}) {
	fmt.Printf(Green+format+Reset+"\n", a...)
}

// PrintRed prints text in red color and adds a newline
func PrintRed(format string, a ...interface{}) {
	fmt.Printf(Red+format+Reset+"\n", a...)
}

// ColorizeText wraps text in the specified color and returns the colorized string
func ColorizeText(color string, text string) string {
	return color + text + Reset
}
