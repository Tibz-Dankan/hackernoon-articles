package articles

import (
	"log"

	"github.com/Tibz-Dankan/hackernoon-articles/internal/models"
	"github.com/Tibz-Dankan/hackernoon-articles/internal/pkg"
)

func UpdateArticleTagIndex() {
	article := models.Article{}

	articles, _, err := article.FindAllByPostedAtInAsc(8000)
	if err != nil {
		log.Printf("Error finding articles: %v", err)
	}
	log.Printf("Total Articles to be index : %d", len(articles))

	for index, currArticle := range articles {
		if currArticle.TagIndex != "" {
			log.Printf("Article : %s already has tag index", currArticle.TagIndex)
			continue
		}
		tagIndex := pkg.BuildTag(index + 1)
		currArticle.TagIndex = tagIndex
		// currArticle.TagIndex = pkg.BuildTag(index + 1)
		log.Println("tagIndex:", tagIndex)

		updatedArticle, err := currArticle.Update()
		if err != nil {
			log.Println("Error updating article Tag Index: ", err)
			continue
		}
		log.Println("Updated Article Tag Index successfully: ", updatedArticle.Title)
	}
}

// func init() {
// 	log.Println("App initialized. Scheduling UpdateArticleTagIndex() to run in 15 seconds...")

// 	go func() {
// 		time.Sleep(15 * time.Second)
// 		start := time.Now()
// 		UpdateArticleTagIndex()
// 		fmt.Printf("Total Article Tag Index Update Duration : %s\n", time.Since(start))
// 	}()
// }
