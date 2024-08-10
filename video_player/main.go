package videoplayer

type VideoPlayer interface {
	IsInstalled() bool
	PlayAnime() error
}
