package webscanner

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func getNumberOfCharacters(url string, ch chan URLInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	if IsURL(url) == false {
		return
	}

	r, err := http.Get(url)
	if err != nil {
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	s := string(body)

	rdr := strings.NewReader(s)
	z := html.NewTokenizer(rdr)
	prevStartToken := z.Token()

	res := URLInfo{url, 0}

loopCountChars:
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			break loopCountChars
		case tt == html.StartTagToken:
			prevStartToken = z.Token()
		case tt == html.TextToken:
			if prevStartToken.Data == "script" || prevStartToken.Data == "style" {
				continue
			}
			txt := strings.TrimSpace(html.UnescapeString(string(z.Text())))
			res.NumOfCharacters += len(txt)
		}
	}
	ch <- res
	return
}

func ScanURL(link string) ([]URLInfo, error) {
	if IsURL(link) == false {
		return []URLInfo{}, errors.New("Incorrect URL!")
	}
	r, err := http.Get(link)
	if err != nil {
		return []URLInfo{}, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []URLInfo{}, err
	}
	s := string(body)

	rdr := strings.NewReader(s)
	z := html.NewTokenizer(rdr)

	result := []URLInfo{}

	var wg sync.WaitGroup

	ch := make(chan URLInfo)

loopSearchForAnchors:
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			break loopSearchForAnchors
		case tt == html.StartTagToken:
			token := z.Token()

			if token.Data == "a" {
				for _, a := range token.Attr {
					if a.Key == "href" {
						wg.Add(1)
						go getNumberOfCharacters(a.Val, ch, &wg)
					}
				}
			}
		}
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for data := range ch {
		if data.NumOfCharacters > 0 {
			result = append(result, data)
		}
	}

	return result, nil
}
