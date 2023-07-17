package extract_data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
	"sync"
)

func GetUrlByDocsDataType[DDType DocsDataType]() string {
	dataType := reflect.TypeOf((*DDType)(nil)).Elem()
	return strings.ToLower(dataType.Name())
}

const urlPrefix string = "https://the-one-api.dev/v2"

var API_TOKEN string

func FetchData[DDType DocsDataType](client *http.Client, wg *sync.WaitGroup, ch chan<- []DDType, page uint16, limit uint16) {
	defer wg.Done()
	url := fmt.Sprintf("%s/%s?page=%d&limit=%d", urlPrefix, GetUrlByDocsDataType[DDType](), page, limit)
	log.Printf("URL created: %s\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Request creation has failed: %s\n", err)
		return
	}
	log.Println("Request created")

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+API_TOKEN)

	res, err := client.Do(req)

	if err != nil {
		log.Printf("Request execution has failed: %s\n", err)
		return
	}
	defer res.Body.Close()
	log.Println("Request completed")

	data := Response[DDType]{}

	buffer := &bytes.Buffer{}
	_, err = io.Copy(buffer, res.Body)

	jsonData := buffer.Bytes()

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Printf("JSON parsing has failed: %s\n", err)
		return
	}
	log.Println("JSON parsed")

	ch <- data.Docs

	log.Println("Data has written to channel")
	return
}
