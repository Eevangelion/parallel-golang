package webscanner

import (
	"bytes"
	"errors"
	"io"
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

func ScanURL(link string) (result URLInfo, err error) {
	if IsURL(link) == false {
		result = URLInfo{}
		err = errors.New("Incorrect URL!")
		return
	}
	r, err := http.Get(link)
	if err != nil {
		result = URLInfo{}
		return
	}
	defer r.Body.Close()
	buffer := &bytes.Buffer{}
	_, err = io.Copy(buffer, r.Body)
	if err != nil {
		result = URLInfo{}
		return
	}
	s := buffer.String()

	rdr := strings.NewReader(s)
	z := html.NewTokenizer(rdr)

	count := 0

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
						count += 1
					}
				}
			}
		}
	}
	result = URLInfo{link, count}
	return
}

func ScanURLS(urls []string) (r URLInfos, e error) {
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			uInfo, err := ScanURL(url)
			if err != nil {
				e = err
			}
			r = append(r, uInfo)
		}(url)
	}
	wg.Wait()
	return r, e
}
