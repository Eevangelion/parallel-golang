package extract_data

import (
	"log"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/Sakagam1/parallel-golang/extract_data/models"
	"github.com/joho/godotenv"
)

var logger *log.Logger

func TestGetUrlByDocsDataType(t *testing.T) {
	logger.Println("Start testing GetUrlByDocsDataType function")
	if name := GetUrlByDocsDataType[models.Character](); name != "character" {
		t.Errorf("Wrong URL detected: expected character, got %s", name)
	}
}

type fetchDataTest[DDType DocsDataType] struct {
	page     uint16
	limit    uint16
	expected []DDType
}

var fetchCharactersDataTests = []fetchDataTest[models.Character]{
	{
		page:  1,
		limit: 1,
		expected: []models.Character{
			{
				Height:  "",
				Race:    "Human",
				Gender:  "Female",
				Birth:   "",
				Spouse:  "Belemir",
				Death:   "",
				Realm:   "",
				Hair:    "",
				Name:    "Adanel",
				WikiUrl: "http://lotr.wikia.com//wiki/Adanel",
			},
		},
	},
	{
		page:  8,
		limit: 3,
		expected: []models.Character{
			{
				Height:  "",
				Race:    "Human",
				Gender:  "Female",
				Birth:   "FA 361",
				Spouse:  "Loved ,Aegnor, but they never married",
				Death:   "FA 455",
				Realm:   "",
				Hair:    "Dark brown",
				Name:    "Andreth",
				WikiUrl: "http://lotr.wikia.com//wiki/Andreth",
			},
			{
				Height:  "",
				Race:    "Human",
				Gender:  "Male",
				Birth:   "FA 440",
				Spouse:  "Unnamed wife",
				Death:   "FA 489",
				Realm:   "",
				Hair:    "",
				Name:    "Andróg",
				WikiUrl: "http://lotr.wikia.com//wiki/Andr%C3%B3g",
			},
			{
				Height:  "",
				Race:    "Elf",
				Gender:  "Male",
				Birth:   "Years of the Trees",
				Spouse:  "None",
				Death:   "FA 538",
				Realm:   "Estolad",
				Hair:    "Dark red",
				Name:    "Amras",
				WikiUrl: "http://lotr.wikia.com//wiki/Amras",
			},
		},
	},
}

func TestFetchData(t *testing.T) {
	logger.Println("Start testing FetchData function")
	for testNum, test := range fetchCharactersDataTests {
		logger.Printf("Test %d started\n", testNum+1)
		var ch = make(chan []models.Character)
		var wg sync.WaitGroup
		wg.Add(1)
		client := &http.Client{}
		go FetchData[models.Character](client, &wg, ch, test.page, test.limit)
		go func() {
			wg.Wait()
			close(ch)
			logger.Printf("Test %d finished\n", testNum+1)
		}()
		result := <-ch
		if len(test.expected) != len(result) {
			t.Errorf("Wrong length of result: expected %d, got %d", len(test.expected), len(result))
		}
		for num, character := range test.expected {
			if character != result[num] {
				t.Error("Wrong character info")
				break
			}
		}
	}
}

func TestMain(m *testing.M) {
	file, err := os.OpenFile("extract_data.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	API_TOKEN = os.Getenv("API_TOKEN")

	logger = log.New(file, "[TEST] ", log.LstdFlags)

	log.SetOutput(file)
	log.SetPrefix("[Extract Data] ")

	logger.Println("Testing started")

	exitCode := m.Run()

	logger.Printf("Testing finished with code %d\n", exitCode)

	os.Exit(exitCode)
}
