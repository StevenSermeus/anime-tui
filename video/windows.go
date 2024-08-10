//go:build windows

package video

import "os/exec"

func (v VLC) Play(video_url string, title string, params ...string) {
	if !v.isInstalled() {
		panic("VLC is not installed")
	}
	cmd := exec.Command("vlc.exe", "--http-referrer", "https://d0000d.com/", "--fullscreen", "--play-and-exit", "--meta-title="+title, video_url)
	cmd.Stdout = nil
	cmd.Stderr = nil
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
}
