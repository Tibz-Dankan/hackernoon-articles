package articles

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	// "time"

	"github.com/Tibz-Dankan/hackernoon-articles/internal/events"
)

type ScrapedData struct {
	ScrapedAt     string           `json:"scraped_at"`
	TotalArticles int              `json:"total_articles"`
	Source        string           `json:"source"`
	Category      string           `json:"category"`
	Articles      []ScrapedArticle `json:"articles"`
}

func ProcessArticles() error {
	filename, err := filepath.Abs("./20250802-213734-hn-bitcoin-articles.json")
	if err != nil {
		log.Println("Error finding absolute path:", err)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var scrapedData ScrapedData
	err = json.Unmarshal(data, &scrapedData)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	fmt.Printf("Successfully loaded %d articles from %s\n\n", len(scrapedData.Articles), filename)

	for _, article := range scrapedData.Articles {
		events.EB.Publish("SAVE_SCRAPED_ARTICLES", article)
	}

	return nil
}

// func init() {
// 	log.Println("App initialized. Scheduling ProcessArticles() to run in 15 seconds...")

// 	go func() {
// 		time.Sleep(15 * time.Second)
// 		ProcessArticles()
// 	}()
// }
