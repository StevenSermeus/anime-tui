package videoprovider

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type VideoUrl struct {
	Url    string
	Reffer string
}

type VideoProvider interface {
	GetVideoUrl() (VideoUrl, error)
}

type D000D struct {
	BaseUrl string
}

func (d D000D) GetVideoUrl() (VideoUrl, error) {
	link, err := d.step1()
	if err != nil {
		return VideoUrl{}, err
	}
	video_link, err := d.step2(link)
	if err != nil {
		return VideoUrl{}, err
	}
	fmt.Println(video_link)
	return VideoUrl{Url: video_link, Reffer: strings.Split(d.BaseUrl, "/")[2]}, nil
}

func (d D000D) step1() (string, error) {
	url := regexp.MustCompile(`/e/`).ReplaceAllString(d.BaseUrl, "/d/")
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	additional_link := ""
	doc.Find("div.download-content a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		additional_link = link
	})
	return "https://www.d0000d.com" + additional_link, nil
}

func (d D000D) step2(link string) (string, error) {
	//Sleep 1 second to avoid getting blocked
	time.Sleep(1 * time.Second)
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Cookie", "lang=1;")
	req.Header.Set("DNT", "1")
	req.Header.Set("Priority", "u=0, i")
	req.Header.Set("Sec-CH-UA", `"Chromium";v="127", "Not)A;Brand";v="99"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	match := regexp.MustCompile(`href="(.*?)"`).FindAllString(string(body), -1)
	video_link := ""
	for _, url := range match {
		if regexp.MustCompile(`mp4`).MatchString(url) {
			video_link = url
		}
	}
	if video_link == "" {
		panic("Couldn't find the video link")
	}
	video_link = regexp.MustCompile(`href="`).ReplaceAllString(video_link, "")
	video_link = video_link[:len(video_link)-1]
	return video_link, nil
}
