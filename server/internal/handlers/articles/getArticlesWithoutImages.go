package articles

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func GetArticlesWithoutImages() {
	filename, err := filepath.Abs("./20250803-004514-hn-bitcoin-articles.json")
	if err != nil {
		log.Println("Error finding absolute path:", err)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var scrapedData ScrapedData
	var scrapedArticles []ScrapedArticle
	err = json.Unmarshal(data, &scrapedData)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	log.Printf("Successfully loaded %d articles from %s\n\n", len(scrapedData.Articles), filename)

	for _, currArticle := range scrapedData.Articles {
		if currArticle.ImageUrl == "" {
			log.Printf("Article : %s has no ImageURL", currArticle.Title)
			scrapedArticles = append(scrapedArticles, currArticle)
		}
	}
	log.Printf("Total Articles Without Images: %v", len(scrapedArticles))
	log.Println("Gotten all articles without images successfully")
}

// func init() {
// 	log.Println("App initialized. Scheduling GetArticlesWithoutImages() to run in 15 seconds...")

// 	go func() {
// 		time.Sleep(15 * time.Second)
// 		start := time.Now()
// 		GetArticlesWithoutImages()
// 		fmt.Printf("Total Get Duration : %s\n", time.Since(start))
// 	}()
// }
