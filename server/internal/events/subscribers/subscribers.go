package subscribers

import (
	"log"

	"github.com/Tibz-Dankan/hackernoon-articles/internal/handlers/articles"
)

func InitEventSubscribers() {
	log.Println("Initiating global event subscribers...")

	go articles.SaveScrapedArticles()
	go articles.ScrapeSingleArticle()
}
