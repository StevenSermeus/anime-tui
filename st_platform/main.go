package stplatform

import (
	videoprovider "github.com/StevenSermeus/anime-tui/video_provider"
)

type Anime struct {
	Title string
	URL   string
}

type Episode struct {
	Title string
	URL   string
}

type StreamingPlatform interface {
	GetAnimeList() ([]Anime, error)
	GetRecentEpisode() ([]Episode, error)
	GetEpisodeList(string) ([]Episode, error)
	GetVideoURL(string) ([]videoprovider.VideoProvider, error)
}
