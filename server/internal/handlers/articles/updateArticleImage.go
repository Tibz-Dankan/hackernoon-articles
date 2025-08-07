package articles

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

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

	filename, err := filepath.Abs("./20250802-213734-hn-bitcoin-articles.json")
	// filename, err := filepath.Abs("./20250807-125125-hackernoon-bitcoin-articles.json")
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

// func init() {
// 	log.Println("App initialized. Scheduling UpdateArticleImage() to run in 15 seconds...")

// 	go func() {
// 		time.Sleep(15 * time.Second)
// 		start := time.Now()
// 		UpdateArticleImage()
// 		fmt.Printf("Total Update Duration : %s\n", time.Since(start))
// 	}()
// }
