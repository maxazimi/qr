//go:build !android && !ios

package platform

import (
	"fmt"
	"gioui.org/io/event"
	"os/exec"
	"runtime"
)

func HandleEvent(e event.Event) {}

func RequestCameraPermission() chan error {
	errChan := make(chan error, 1)
	defer close(errChan)
	errChan <- nil
	return errChan
}

func OpenURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	default:
		return fmt.Errorf("unsupported platform")
	}
	return exec.Command(cmd, args...).Start()
}
