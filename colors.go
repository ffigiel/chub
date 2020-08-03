package main

const (
	ansiReset   = "\u001b[0m"
	ansiRed     = "\u001b[31m"
	ansiGreen   = "\u001b[32m"
	ansiYellow  = "\u001b[33m"
	ansiBlue    = "\u001b[34m"
	ansiMagenta = "\u001b[35m"
	ansiCyan    = "\u001b[36m"
)

var colors = []string{
	ansiRed,
	ansiGreen,
	ansiYellow,
	ansiBlue,
	ansiMagenta,
	ansiCyan,
}

func getColor(n int) string {
	return colors[n%6]
}

func withColor(s, color string) string {
	return color + s + ansiReset
}
