//go:build windows

package videoplayer

func (v VLC) PlayAnime() error {
	panic("vlc is not supported on windows yet")
}

func (v VLC) IsInstalled() bool {
	panic("vlc is not supported on windows yet")
}
