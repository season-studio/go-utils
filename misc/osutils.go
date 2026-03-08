package misc

import (
	"os"
	"runtime"
)

func GetShell() (string, string) {
	if runtime.GOOS == "windows" {
		shell := os.Getenv("ComSpec")
		if shell == "" {
			shell = "cmd.exe"
		}
		return shell, "/c"
	}
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "sh"
	}
	return shell, "-c"
}
