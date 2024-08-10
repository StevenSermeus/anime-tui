//go:build !windows

package video

import (
	"os/exec"
	"syscall"
)

func (v VLC) Play(video_url string, title string, params ...string) {
	if !v.isInstalled() {
		panic("VLC is not installed")
	}
	cmd := exec.Command("vlc", "--http-referrer", "https://d0000d.com/", "--fullscreen", "--play-and-exit", "--meta-title="+title, video_url)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}
	cmd.Stdout = nil
	cmd.Stderr = nil
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
}
