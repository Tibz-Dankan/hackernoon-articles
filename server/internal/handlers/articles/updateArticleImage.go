package articles

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Tibz-Dankan/hackernoon-articles/internal/events"
	"github.com/Tibz-Dankan/hackernoon-articles/internal/models"
	"github.com/Tibz-Dankan/hackernoon-articles/internal/pkg"
)

func UpdateArticleImage() {
	article := models.Article{}
	imageProcessor := pkg.ImageProcessor{}

	s3Client := pkg.S3Client{}
	ctx := context.Background()

	newS3Client, err := s3Client.NewS3Client(ctx)
	if err != nil {
		log.Printf("Error creating newS3Client: %v", err)
	}

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

	// articles, _, err := article.FindAll(6000, "")
	articles, count, err := article.FindAllWithWrongImage(6000, "")
	if err != nil {
		log.Printf("Error finding articles: %v", err)
	}
	log.Printf("Article With Wrong Images: %v", count)

	for _, currArticle := range articles {
		if !strings.Contains(currArticle.ImageUrl, "?") {
			log.Printf("Article : %s has correct ImageURL", currArticle.Title)
			continue
		}
		var scrapedArticleImageURL string

		for _, scrapedArticle := range scrapedData.Articles {
			if currArticle.Title == scrapedArticle.Title {
				scrapedArticleImageURL = scrapedArticle.ImageUrl
				break
			}
		}

		if scrapedArticleImageURL == "" {
			log.Printf("Article : %s has no ImageURL", currArticle.Title)
			log.Printf("Article Image URL : %s has no ImageURL", currArticle.ImageUrl)
			continue
		}

		articleImgBuf, getImgErr := imageProcessor.GetImageFromURL(scrapedArticleImageURL)
		if getImgErr != nil {
			log.Println("Error getting author's avatar image from url : ", err)
			continue
		}

		contentType, err := imageProcessor.GetContentTypeFromBinary(articleImgBuf)
		if err != nil {
			log.Println("Error getting author avatar image content type : ", err)
			continue
		}
		log.Println("Content type:", contentType)

		imgFile := imageProcessor.BinaryToReader(articleImgBuf)

		uploadImageResp, err := newS3Client.UploadFile(
			ctx,
			imgFile,
			scrapedArticleImageURL,
			contentType,
			0,
		)
		if err != nil {
			log.Println("Error uploading file to s3 : ", err)
			continue
		}

		err = newS3Client.DeleteFile(ctx, currArticle.ImageFilename)
		if err != nil {
			log.Println("Error deleting old file to s3 : ", err)
			continue
		}

		currArticle.ImageUrl = uploadImageResp.URL
		currArticle.ImageFilename = uploadImageResp.Filename

		updatedArticle, err := currArticle.Update()
		if err != nil {
			log.Println("Error creating article : ", err)
			continue
		}
		log.Println("Updated Article successfully: ", updatedArticle.Title)
	}
}

func UpdateArticleImageV2() {
	// // Update many Articles
	// article := models.Article{}
	// articles, count, err := article.FindAllWithWrongImage(6000, "")
	// if err != nil {
	// 	log.Printf("Error finding articles: %v", err)
	// }
	// log.Printf("Article With Wrong Images: %v", count)

	// for _, currArticle := range articles {
	// 	if !strings.Contains(currArticle.ImageUrl, "?") {
	// 		log.Printf("Article : %s has correct ImageURL", currArticle.Title)
	// 		continue
	// 	}
	// 	events.EB.Publish("SCRAPE_SINGLE_ARTICLE_v2", currArticle)
	// 	log.Println("Updated Article Image Initiated: ", currArticle.Title)
	// }

	// Update one Article
	article := models.Article{}
	savedArticle, err := article.FindByTitle("Bitcoin Mining Could Make Our Electricity Grids Smarter")
	if err != nil {
		log.Printf("Error finding article: %v", err)
	}
	if savedArticle.ID != "" {
		events.EB.Publish("SCRAPE_SINGLE_ARTICLE_v2", savedArticle)
		log.Println("Updated Article Image Initiated: ", savedArticle.Title)
	}
}

// func init() {
// 	log.Println("App initialized. Scheduling UpdateArticleImageV2() to run in 15 seconds...")

// 	go func() {
// 		time.Sleep(15 * time.Second)
// 		start := time.Now()
// 		// UpdateArticleImage()
// 		UpdateArticleImageV2()
// 		fmt.Printf("Total Update Duration : %s\n", time.Since(start))
// 	}()
// }
