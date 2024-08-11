# Anime tui

It's a simple terminal application that allows you to browse and watch anime in VOSTFR from your terminal.

![demo](.github/assets/demo.gif)

## Requirements

You only need to have `vlc` installed on your system.

## Installation

Adapt the following commands to your system.

### Mac arm

```bash
wget https://github.com/StevenSermeus/anime-tui/releases/download/{version}/anime-tui_Darwin_arm64.tar.gz
tar -xvf anime-tui_Darwin_arm64.tar
mv anime-tui /usr/local/bin
```

### Mac amd64

```bash
wget  https://github.com/StevenSermeus/anime-tui/releases/download/{version}/anime-tui_Darwin_x86_64.tar.gz
tar -xvf anime-tui_Darwin_x86_64.tar
mv anime-tui /usr/local/bin
```

### Linux arm

```bash
wget https://github.com/StevenSermeus/anime-tui/releases/download/{version}/anime-tui_Linux_arm64.tar.gz
tar -xvf anime-tui_Linux_arm64.tar
mv anime-tui /usr/local/bin
```

### Linux amd64

```bash
wget https://github.com/StevenSermeus/anime-tui/releases/download/{version}/anime-tui_Linux_x86_64.tar.gz
tar -xvf anime-tui_Linux_x86_64.tar
mv anime-tui /usr/local/bin
```

### From source

Requires go 1.22.5 installed on your system.

```bash
make install
```

## Usage

```bash
anime-tui
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## Start and build

Have go 1.22.5 installed on your system and make if you want to have an easy way to build the project.

```bash
make
```

To build for all platforms you need to have goreleaser installed on your system.

```bash
make build-all
```

## Hoverview of the project

```
custom_error -> CustomError is the package to handle custom errors.
├── st_platform -> StPlatform is the package to interact with the streaming platform.
├── tui -> TUI is the package to interact with the terminal.
├── video_player -> VideoPlayer is the package to interact with the video player (VLC, ...).
│   └── vlc
└── video_provider -> VideoProvider is the package to interact with the video provider present on streaming platforms.
```

## Roadmap

- [ ] Add autoupdate
- [ ] Add download feature

## License

By contributing, you agree that your contributions will be licensed under its GNU AFFERO GENERAL PUBLIC LICENSE Version 3.
