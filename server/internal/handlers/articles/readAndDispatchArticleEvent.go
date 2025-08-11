package articles

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"time"

	"github.com/Tibz-Dankan/hackernoon-articles/internal/constants"
	"github.com/Tibz-Dankan/hackernoon-articles/internal/events"
	"github.com/Tibz-Dankan/hackernoon-articles/internal/models"
)

type ScrapedData struct {
	ScrapedAt     string           `json:"scraped_at"`
	TotalArticles int              `json:"total_articles"`
	Source        string           `json:"source"`
	Category      string           `json:"category"`
	Articles      []ScrapedArticle `json:"articles"`
}

func ProcessArticles() error {
	filename, err := filepath.Abs("./20250810-114840-hackernoon-bitcoin-articles.json")
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

func ProcessArticlesWithoutImages() error {
	article := models.Article{}
	filename, err := filepath.Abs("./20250803-004514-hn-bitcoin-articles.json")
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

	for _, currArticle := range scrapedData.Articles {
		if currArticle.ImageUrl == "" && currArticle.AuthorName != "" && currArticle.Title != "" {

			var isBlackListed bool
			for _, blArticle := range blackListArticles {
				if blArticle == currArticle.Title {
					log.Printf("BlackListed: %s", currArticle.Title)
					isBlackListed = true
				}
			}
			if isBlackListed {
				continue
			}

			savedArticle, err := article.FindByTitle(currArticle.Title)
			if err != nil && err.Error() != constants.RECORD_NOT_FOUND_ERROR {
				log.Printf("Error finding the saved article: %v", err)
				continue
			}
			if savedArticle.ID != "" {
				log.Printf("Article is already saved: %s ", savedArticle.Title)
				continue
			}

			log.Printf("Publishing article: %s", currArticle.Title)
			events.EB.Publish("SCRAPE_SINGLE_ARTICLE", currArticle)
		}
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

func init() {
	log.Println("App initialized. Scheduling ProcessArticlesWithoutImages() to run in 15 seconds...")

	go func() {
		time.Sleep(15 * time.Second)
		ProcessArticlesWithoutImages()
	}()
}

var blackListArticles = []string{
	"Decentralized Applications Will Take Cryptocurrency to the Mainstream",
	"What is The Bitcoin Halving and What Impact Will It Have on the Crypto Market?",
}
