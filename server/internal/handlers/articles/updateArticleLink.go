package articles

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/Tibz-Dankan/hackernoon-articles/internal/models"
)

func UpdateArticleLink() {
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
	log.Printf("Successfully loaded %d articles from %s\n\n", len(scrapedData.Articles), filename)

	articles, _, err := article.FindAll(6000, "")
	if err != nil {
		log.Printf("Error finding articles: %v", err)
	}

	for _, currArticle := range articles {
		if currArticle.Href != "" {
			log.Printf("Article : %s has correct link", currArticle.Title)
			continue
		}
		var scrapedArticleURL string

		for _, scrapedArticle := range scrapedData.Articles {
			if currArticle.Title == scrapedArticle.Title {
				scrapedArticleURL = scrapedArticle.URL
				break
			}
		}

		if scrapedArticleURL == "" {
			log.Printf("Article : %s has no article Link", currArticle.Title)
			continue
		}

		currArticle.Href = scrapedArticleURL

		updatedArticle, err := currArticle.Update()
		if err != nil {
			log.Println("Error updating article link: ", err)
			continue
		}
		log.Println("Updated Article successfully: ", updatedArticle.Title)
	}
}

// func init() {
// 	log.Println("App initialized. Scheduling UpdateArticleLink() to run in 15 seconds...")

// 	go func() {
// 		time.Sleep(15 * time.Second)
// 		start := time.Now()
// 		UpdateArticleLink()
// 		fmt.Printf("Total Article Link Update Duration : %s\n", time.Since(start))
// 	}()
// }
