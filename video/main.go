package video

import (
	"os/exec"
	"runtime"
)

type Player interface {
	Play(video_url string, title string, params ...string)
	isInstalled() bool
}

type VLC struct{}

func (v VLC) isInstalled() bool {
	if runtime.GOOS == "windows" {
		err := exec.Command("vlc.exe", "--version").Run()
		return err == nil
	} else {
		err := exec.Command("vlc", "--version").Run()
		return err == nil
	}
}
