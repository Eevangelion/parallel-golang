package webscanner

import (
	"fmt"
	"sort"
	"testing"
)

type scanTest struct {
	urls     []string
	expected URLInfos
}

var errorScanTests = []scanTest{
	{[]string{"https://sci-hub.ru/", "/foo/bar/", "http//google.com", ""}, URLInfos{}},
	{[]string{"dsadsad/ffdf", "http://google*com/", "https://google.com/", "google.com/"}, URLInfos{}},
}

var scanTests = []scanTest{
	{[]string{"https://codeforces.com/problemset/problem/1847/B", "https://www.youtube.com/"}, URLInfos{
		{"https://codeforces.com/problemset/problem/1847/B", 32},
		{"https://www.youtube.com/", 14},
	}},
	{[]string{"https://google.com/", "https://www.litres.ru/"}, URLInfos{
		{"https://google.com/", 18},
		{"https://www.litres.ru/", 736},
	}},
}

func TestScanURLS(t *testing.T) {
	for num, test := range errorScanTests {
		t.Run(fmt.Sprintf("Test %d", num+1),
			func(t *testing.T) {
				_, err := ScanURLS(test.urls)
				if err == nil {
					t.Errorf("Should have panicked because of an incorrect URL!")
				}
			})
	}
	for num, test := range scanTests {
		t.Run(fmt.Sprintf("Test %d", num+1),
			func(t *testing.T) {
				result, err := ScanURLS(test.urls)
				if err != nil {
					t.Errorf("%s", err)
				}
				if len(test.expected) != len(result) {
					t.Errorf("Length of expected: %d is not equal to length of given: %d", len(test.expected), len(result))
					return
				}
				sort.Sort(result)

				for numUrl, uInfo := range test.expected {
					if uInfo.NumOfLinks != result[numUrl].NumOfLinks ||
						uInfo.Url != test.expected[numUrl].Url {
						t.Errorf("Expected %d, actual %d", uInfo.NumOfLinks, result[numUrl].NumOfLinks)
						break
					}
				}
			})
	}
}
