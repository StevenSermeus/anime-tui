//go:build !windows

package videoplayer

import (
	"os/exec"
	"syscall"

	customerror "github.com/StevenSermeus/anime-tui/custom_error"
)

func (v VLC) PlayAnime() error {
	if v.VideoUrl == "" {
		return customerror.MissingVideoUrl{Err: "video URL is missing"}
	}
	if v.Title == "" {
		v.Title = "Provider by anime-tui"
	}
	if v.Refferer == "" {
		v.Refferer = "https://www.google.com"
	}
	cmd := exec.Command("vlc", "--http-referrer", v.Refferer, "--fullscreen", "--play-and-exit", "--meta-title="+v.Title, v.VideoUrl)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}
	cmd.Stdout = nil
	cmd.Stderr = nil
	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

func (v VLC) IsInstalled() bool {
	return exec.Command("vlc", "--version").Run() == nil
}
