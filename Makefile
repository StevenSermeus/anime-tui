build:
	CGO_ENABLED=0 go build -o bin/anime-tui 

install:
	make build
	cp bin/anime-tui ~/.local/bin/anime-tui

dev:
	go mod tidy
	go run anime-tui.go

run:
	go run anime-tui.go

build-all:
	go mod tidy
	goreleaser build --snapshot --clean