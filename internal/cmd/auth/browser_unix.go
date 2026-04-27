//go:build !windows

package auth

import (
	"os/exec"
	"runtime"
)

func browserCmd(url string) *exec.Cmd {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url)
	default:
		return exec.Command("xdg-open", url)
	}
}
