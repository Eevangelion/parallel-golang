package webscanner

import (
	"fmt"
	"sort"
	"testing"
)

type scanTest struct {
	url        string
	expected   []URLInfo
	errMessage string
}

var errorScanTests = []scanTest{
	{"/foo/bar/", []URLInfo{}, "Should have panicked because of an incorrect URL!"},
	{"http//google.com", []URLInfo{}, "Should have panicked because of an incorrect URL!"},
	{"dsadsad/ffdf", []URLInfo{}, "Should have panicked because of an incorrect URL!"},
	{"", []URLInfo{}, "Should have panicked because of an incorrect URL!"},
	{"https:://google.com/", []URLInfo{}, "Should have panicked because of an incorrect URL!"},
	{"http://google*com/", []URLInfo{}, "Should have panicked because of an incorrect URL!"},
	{"google.com/", []URLInfo{}, "Should have panicked because of an incorrect URL!"},
}

var scanTests = []scanTest{
	{"https://sci-hub.ru/", []URLInfo{}, "Could not handle correct URL!"},
	{"https://google.com", []URLInfo{
		{"https://mail.google.com/mail/&amp;ogbl", 3029},
		{"https://www.google.com/imghp?hl=ru&amp;ogbl", 118},
		{"https://www.google.ru/intl/ru/about/products", 4508},
		{"https://accounts.google.com/ServiceLogin?hl=ru&amp;passive=true&amp;continue=https://www.google.com/&amp;ec=GAZAmgQ", 162},
		{"https://www.google.com/setprefs?sig=0_a8BK0oZ7Y94hsrDNGzQ1rLVDQ1Q%3D&amp;hl=en&amp;source=homepage&amp;sa=X&amp;ved=0ahUKEwjWuruin_b_AhUPyYsKHSuWDcYQ2ZgBCBA", 153},
		{"https://about.google/?utm_source=google-RU&amp;utm_medium=referral&amp;utm_campaign=hp-footer&amp;fg=1", 2506},
		{"https://www.google.com/intl/ru_ru/ads/?subid=ww-ww-et-g-awa-a-g_hpafoot1_1!o2&amp;utm_source=google.com&amp;utm_medium=referral&amp;utm_campaign=google_hpafooter&amp;fg=1", 2902},
		{"https://www.google.com/services/?subid=ww-ww-et-g-awa-a-g_hpbfoot1_1!o2&amp;utm_source=google.com&amp;utm_medium=referral&amp;utm_campaign=google_hpbfooter&amp;fg=1", 4200},
		{"https://google.com/search/howsearchworks/?fg=1", 1921},
		{"https://sustainability.google/intl/ru/carbon-free/?utm_source=googlehpfooter&amp;utm_medium=housepromos&amp;utm_campaign=bottom-footer&amp;utm_content=", 8096},
		{"https://policies.google.com/privacy?hl=ru&amp;fg=1", 34636},
		{"https://policies.google.com/terms?hl=ru&amp;fg=1", 24532},
		{"https://www.google.com/preferences?hl=ru&amp;fg=1", 930},
		{"https://support.google.com/websearch/?p=ws_results_help&amp;hl=ru&amp;fg=1", 886},
	}, "Could not handle correct URL!"},
	// {"https://github.com/", []URLInfo{}, "Could not handle correct URL!"},
	// {"https://shikimori.me/", []URLInfo{}, "Could not handle correct URL!"},
	// {"https://tex.stackexchange.com/questions", []URLInfo{}, "Could not handle correct URL!"},
	// {"https://codeforces.com/gym/104460/standings", []URLInfo{}, "Could not handle correct URL!"},
}

func TestScanURLS(t *testing.T) {
	for num, test := range errorScanTests {
		t.Run(fmt.Sprintf("Test %d", num+1),
			func(t *testing.T) {
				_, err := ScanURL(test.url)
				if err == nil {
					t.Errorf(fmt.Sprintf(test.errMessage))
				}
			})
	}
	for num, test := range scanTests {
		t.Run(fmt.Sprintf("Test %d", num+1),
			func(t *testing.T) {
				result, err := ScanURL(test.url)
				if err != nil {
					t.Errorf(fmt.Sprintf("%s", err))
				}
				sort.Sort(URLInfos(result))
				sort.Sort(URLInfos(test.expected))
				for numUrl, uInfo := range result {
					if uInfo.NumOfCharacters != test.expected[numUrl].NumOfCharacters ||
						uInfo.Url != test.expected[numUrl].Url {
						t.Errorf(fmt.Sprintf("Expected %d, actual %d", test.expected[numUrl].NumOfCharacters, uInfo.NumOfCharacters))
						break
					}
				}
			})
	}
}
