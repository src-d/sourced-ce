package format

import (
	"fmt"
	"runtime"
)

// Color represents a color code
type Color string

const (
	// Red for errors
	Red Color = "31"
	// Yellow for warnings
	Yellow Color = "33"
)

// Colorize returns the passed string with the passed color
func Colorize(color Color, s string) string {
	if runtime.GOOS == "windows" {
		return s
	}

	return fmt.Sprintf("\x1b[%sm%s\x1b[0m", color, s)
}
