package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Version() string {
	return "1.0.0"
}

type LatestVersion struct {
	Name       string `json:"name"`
	ZipballURL string `json:"zipball_url"`
	TarballURL string `json:"tarball_url"`
	Commit     struct {
		SHA string `json:"sha"`
		URL string `json:"url"`
	} `json:"commit"`
	NodeID string `json:"node_id"`
}

func getLatestVersion() string {
	res, err := http.Get("https://api.github.com/repos/StevenSermeus/anime-tui/tags")
	if err != nil {
		return ""
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ""
	}
	json_res := []LatestVersion{}
	err = json.Unmarshal(body, &json_res)
	if err != nil {
		return ""
	}
	return json_res[0].Name
}

func CheckForUpdate() {
	latest := getLatestVersion()
	if latest == "" {
		return
	}
	if strings.Contains("dev", Version()) {
		return
	}
	if latest == Version() {
		return
	}
	fmt.Printf("A new version of anime-tui is available: %s\n", latest)
}
