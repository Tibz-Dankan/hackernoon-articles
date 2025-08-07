package articles

import (
	"context"
	"log"

	"github.com/Tibz-Dankan/hackernoon-articles/internal/constants"
	"github.com/Tibz-Dankan/hackernoon-articles/internal/events"
	"github.com/Tibz-Dankan/hackernoon-articles/internal/models"
	"github.com/Tibz-Dankan/hackernoon-articles/internal/pkg"
)

func SaveScrapedArticles() {
	go func() {
		scrapedArticleChan := make(chan events.DataEvent)
		events.EB.Subscribe("SAVE_SCRAPED_ARTICLES", scrapedArticleChan)
		article := models.Article{}
		author := models.Author{}

		imageProcessor := pkg.ImageProcessor{}

		s3Client := pkg.S3Client{}

		ctx := context.Background()

		newS3Client, err := s3Client.NewS3Client(ctx)
		if err != nil {
			log.Printf("Error creating newS3Client: %v", err)
		}

		for {
			scrapedArticleEvent := <-scrapedArticleChan
			scrapedArticle, ok := scrapedArticleEvent.Data.(ScrapedArticle)
			if !ok {
				log.Printf("Invalid articleData type received: %T", scrapedArticle)
				continue
			}
			if scrapedArticle.AuthorName == "" || scrapedArticle.Title == "" {
				log.Printf("Article is has no title or author's name ")
				continue
			}
			log.Printf("Saving article in progress %s:", scrapedArticle.Title)

			savedArticle, err := article.FindByTitle(scrapedArticle.Title)
			if err != nil && err.Error() != constants.RECORD_NOT_FOUND_ERROR {
				log.Printf("Error finding the saved article: %v", err)
				continue
			}
			if savedArticle.ID != "" {
				log.Printf("Article is ready saved: %s ", scrapedArticle.Title)
				continue
			}

			articleAuthor, err := author.FindByName(scrapedArticle.AuthorName)
			if err != nil && err.Error() != constants.RECORD_NOT_FOUND_ERROR {
				log.Printf("Error finding article's author: %v", err)
				continue
			}

			// Create author if doesn't exist
			if articleAuthor.ID == "" {
				avatarImgBuf, getImgErr := imageProcessor.GetImageFromURL(scrapedArticle.AuthorAvatarUrl)
				if getImgErr != nil {
					log.Println("Error getting author's avatar image from url : ", err)
					continue
				}

				// Upload author's avatar Image
				// if len(avatarImgBuf) > 0 && getImgErr == nil {
				if len(avatarImgBuf) > 0 {
					contentType, err := imageProcessor.GetContentTypeFromBinary(avatarImgBuf)
					if err != nil {
						log.Println("Error getting author avatar image content type : ", err)
						continue
					}
					log.Println("Content type:", contentType)

					imgFile := imageProcessor.BinaryToReader(avatarImgBuf)

					uploadAvatarResp, err := newS3Client.UploadFile(
						ctx,
						imgFile,
						scrapedArticle.AuthorAvatarUrl,
						contentType,
						0,
					)
					if err != nil {
						log.Println("Error uploading file to s : ", err)
						continue
					}
					articleAuthor, err = author.Create(models.Author{
						Name:           scrapedArticle.AuthorName,
						PageUrl:        scrapedArticle.AuthorPageURL,
						AvatarUrl:      uploadAvatarResp.URL,
						AvatarFilename: uploadAvatarResp.Filename,
					})
					if err != nil {
						log.Println("Error creating author : ", err)
						continue
					}
				}

			}

			article.AuthorID = articleAuthor.ID
			article.Tag = scrapedArticle.Tag
			article.Title = scrapedArticle.Title
			article.Href = scrapedArticle.URL
			article.PostedAt = scrapedArticle.PostedAt
			article.ReadDuration = scrapedArticle.ReadDuration
			article.ImageFilename = "ImageFilename.jpeg"
			article.ImageUrl = scrapedArticle.ImageUrl

			articleImgBuf, getImgErr := imageProcessor.GetImageFromURL(scrapedArticle.ImageUrl)
			if getImgErr != nil {
				log.Println("Error getting author's avatar image from url : ", err)
				continue
			}

			// Upload article Image
			// if len(articleImgBuf) > 0 && getImgErr == nil {
			if len(articleImgBuf) > 0 {
				contentType, err := imageProcessor.GetContentTypeFromBinary(articleImgBuf)
				if err != nil {
					log.Println("Error getting article image content type : ", err)
					continue
				}
				log.Println("Content type:", contentType)

				imgFile := imageProcessor.BinaryToReader(articleImgBuf)

				uploadImageResp, err := newS3Client.UploadFile(
					ctx,
					imgFile,
					scrapedArticle.ImageUrl,
					contentType,
					0,
				)

				if err != nil {
					log.Println("Error uploading article image to s3", err)
				}

				article.ImageUrl = uploadImageResp.URL
				article.ImageFilename = uploadImageResp.Filename
			}

			createdArticle, err := article.Create(article)
			if err != nil {
				log.Println("Error creating article : ", err)
				continue
			}
			log.Println("Successfully createdArticle: ", createdArticle.Title)
		}
	}()
}
