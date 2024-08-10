package network

import (
	"fmt"
	"io"
	"net/http"
	URL "net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/briandowns/spinner"
)

func GetAnimeList() ([][]string, string, string, error) {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	defer s.Stop()
	response, err := http.Get("https://mavanimes.cc/tous-les-animes-en-vostfr")
	if err != nil {
		return nil, "", "", err
	}
	defer response.Body.Close()
	body, error := io.ReadAll(response.Body)
	if error != nil {
		return nil, "", "", error
	}
	html := string(body)

	regex := regexp.MustCompile(`<a href="anime/(.*?)">(.*?)</a>`)
	matches := regex.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		match[1] = "https://mavanimes.cc/anime/" + match[1]
	}
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
		return nil, "", "", fmt.Errorf("couldn't find the cookies")
	}
	decoded_xsrf, err := URL.QueryUnescape(xsrf)
	if err != nil {
		return nil, "", "", err
	}
	return matches, decoded_xsrf, mavanimes_session, nil
}

func GetRecentEpisode() ([][]string, string, string, error) {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	defer s.Stop()
	res, err := http.Get("https://www.mavanimes.cc/")

	if err != nil {
		return nil, "", "", err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, "", "", err
	}
	episode_list := [][]string{}
	doc.Find("div.animes-grid div.item a").Each(func(i int, s *goquery.Selection) {
		ep_url, _ := s.Attr("href")
		ep_title := s.Text()
		ep_title = regexp.MustCompile(`\s+`).ReplaceAllString(ep_title, " ")
		if !strings.Contains(ep_url, "anime") {
			episode_list = append(episode_list, []string{"https://mavanimes.cc/" + ep_url, ep_title})
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
		return nil, "", "", fmt.Errorf("couldn't get the cookies")
	}
	decoded_xsrf, err := URL.QueryUnescape(xsrf)
	if err != nil {
		return nil, "", "", err
	}
	return episode_list, decoded_xsrf, mavanimes_session, nil
}

func GetEpisodeList(url string) ([][]string, error) {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()

	response, err := http.Get(url)
	s.Stop()
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
	for _, match := range matches {
		match[1] = "https://mavanimes.cc" + match[1]
		match[1] = strings.TrimSpace(match[1])
		match[2] = strings.TrimSpace(match[2])

	}
	return matches, nil
}

func GetVideoLink(url string, mavToken string, xcrf string) ([]string, error) {
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
	req.Header.Set("X-XSRF-TOKEN", xcrf)
	req.Header.Set("Cookie", fmt.Sprintf("XSRF-TOKEN=%s; mavanimes_session=%s", xcrf, mavToken))
	req.Header.Set("Referer", url)
	req.Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	match := regexp.MustCompile(`src=\\"(.*?)\\"`).FindAllString(string(body), -1)
	player_url := []string{}
	for _, url := range match {
		// Only support d0000d.com and streamwsh.click
		if regexp.MustCompile(`d0000d.com|streamwsh.click`).MatchString(url) {
			player_url = append(player_url, cleanUrl(url))
		}
	}
	return player_url, err
}

func cleanUrl(url string) string {
	url = url[6 : len(url)-2]
	url = regexp.MustCompile(`\\\/`).ReplaceAllString(url, "/")
	url = regexp.MustCompile(`;.*`).ReplaceAllString(url, "")
	return url
}
