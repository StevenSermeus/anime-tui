package stplatform

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	customerror "github.com/StevenSermeus/anime-tui/custom_error"
	videoprovider "github.com/StevenSermeus/anime-tui/video_provider"
	"github.com/briandowns/spinner"
)

type Mavanimes struct {
	xsrf    string
	session string
}

func (m *Mavanimes) GetAnimeList() ([]Anime, error) {
	spinner := spinner.New(spinner.CharSets[30], 100*time.Millisecond)
	spinner.Start()
	defer spinner.Stop()
	response, err := http.Get("https://mavanimes.cc/tous-les-animes-en-vostfr")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, error := io.ReadAll(response.Body)
	if error != nil {
		return nil, error
	}
	html := string(body)
	regex := regexp.MustCompile(`<a href="anime/(.*?)">(.*?)</a>`)
	matches := regex.FindAllStringSubmatch(html, -1)
	cookies := response.Cookies()
	xsrf := ""
	mavanimes_session := ""
	for _, cookie := range cookies {
		if cookie.Name == "XSRF-TOKEN" {
			xsrf = cookie.Value
		}
		if cookie.Name == "mavanimes_session" {
			mavanimes_session = cookie.Value
		}
	}
	if xsrf == "" || mavanimes_session == "" {
		return nil, customerror.XSRFTokenNotFound{Err: "couldn't find the cookies"}
	}
	decoded_xsrf, err := url.QueryUnescape(xsrf)
	if err != nil {
		return nil, err
	}
	m.session = mavanimes_session
	m.xsrf = decoded_xsrf
	animes := make([]Anime, len(matches))
	for i, match := range matches {
		animes[i] = Anime{
			Title: match[2],
			URL:   "https://mavanimes.cc/anime/" + match[1],
		}
	}
	return animes, nil
}

func (m *Mavanimes) GetRecentEpisode() ([]Episode, error) {
	spinner := spinner.New(spinner.CharSets[30], 100*time.Millisecond)
	spinner.Start()
	defer spinner.Stop()
	res, err := http.Get("https://www.mavanimes.cc/")

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	episode_list := []Episode{}
	doc.Find("div.animes-grid div.item a").Each(func(i int, s *goquery.Selection) {
		ep_url, _ := s.Attr("href")
		ep_title := s.Text()
		ep_title = regexp.MustCompile(`\s+`).ReplaceAllString(ep_title, " ")
		if !strings.Contains(ep_url, "anime") {
			episode_list = append(episode_list, Episode{URL: "https://mavanimes.cc/" + ep_url, Title: ep_title})
		}
	})
	cookies := res.Cookies()
	xsrf := ""
	mavanimes_session := ""
	for _, cookie := range cookies {
		if cookie.Name == "XSRF-TOKEN" {
			xsrf = cookie.Value
		}
		if cookie.Name == "mavanimes_session" {
			mavanimes_session = cookie.Value
		}
	}
	if xsrf == "" || mavanimes_session == "" {
		return nil, customerror.XSRFTokenNotFound{Err: "couldn't find the cookies"}
	}
	decoded_xsrf, err := url.QueryUnescape(xsrf)
	if err != nil {
		return nil, err
	}
	m.session = mavanimes_session
	m.xsrf = decoded_xsrf
	return episode_list, nil
}

func (m *Mavanimes) GetEpisodeList(url string) ([]Episode, error) {
	spinner := spinner.New(spinner.CharSets[30], 100*time.Millisecond)
	spinner.Start()
	defer spinner.Stop()
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	html := string(body)
	re := regexp.MustCompile(`<a href="([^"]+)">\s*:•s*([^<]+)`)
	matches := re.FindAllStringSubmatch(html, -1)
	sort.Slice(matches, func(i, j int) bool {
		return matches[i][2] < matches[j][2]
	})
	episode_list := make([]Episode, len(matches))
	for i, match := range matches {
		match[1] = "https://mavanimes.cc" + match[1]
		match[1] = strings.TrimSpace(match[1])
		match[2] = strings.TrimSpace(match[2])
		episode_list[i] = Episode{URL: match[1], Title: match[2]}
	}
	return episode_list, nil
}

func (m *Mavanimes) GetVideoURL(url string) ([]videoprovider.VideoProvider, error) {
	spinner := spinner.New(spinner.CharSets[30], 100*time.Millisecond)
	spinner.Start()
	defer spinner.Stop()
	res, err := m.httpPostWithCookies(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	match := regexp.MustCompile(`src=\\"(.*?)\\"`).FindAllString(string(body), -1)
	player_url := []videoprovider.VideoProvider{}
	for _, ma := range match {
		link := m.cleanUrl(ma)
		if regexp.MustCompile(`d0000d.com`).MatchString(link) {
			player_url = append(player_url, videoprovider.D000D{BaseUrl: link})
		}
	}
	if len(player_url) == 0 {
		return nil, customerror.VideoURLNotFound{Err: "couldn't find the video url"}
	}
	return player_url, nil
}

func (m *Mavanimes) httpPostWithCookies(url string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-CH-UA", `"Chromium";v="127", "Not)A;Brand";v="99"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("X-XSRF-TOKEN", m.xsrf)
	req.Header.Set("Cookie", fmt.Sprintf("XSRF-TOKEN=%s; mavanimes_session=%s", m.session, m.session))
	req.Header.Set("Referer", url)
	req.Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	client := &http.Client{}
	return client.Do(req)
}

func (m *Mavanimes) cleanUrl(url string) string {
	url = url[6 : len(url)-2]
	url = regexp.MustCompile(`\\\/`).ReplaceAllString(url, "/")
	url = regexp.MustCompile(`;.*`).ReplaceAllString(url, "")
	return url
}
