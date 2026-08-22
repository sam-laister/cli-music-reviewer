package services

import (
	"fmt"
	"os/exec"
	"runtime"
)

type BrowserServiceImpl struct {
}

func (b BrowserServiceImpl) Open(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

func NewBrowserService() *BrowserServiceImpl {
	return &BrowserServiceImpl{}
}
