//go:build !windows

package auth

import (
	"os/exec"
	"runtime"
)

func browserCmd(targetURL string) *exec.Cmd {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", targetURL)
	default:
		return exec.Command("xdg-open", targetURL)
	}
}
