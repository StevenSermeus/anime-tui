//go:build windows

package videoplayer

import (
	"os/exec"
)

func (v VLC) PlayAnime() error {
	panic("vlc is not supported on windows yet")
}

func (v VLC) IsInstalled() bool {
	panic("vlc is not supported on windows yet")
	return exec.Command("vlc.exe", "--version").Run() == nil
}
