package articles

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Tibz-Dankan/hackernoon-articles/internal/events"
	"github.com/chromedp/chromedp"
)

func ScrapeSingleArticleImage(articleURL string) (string, error) {
	// Create optimized context for single article scraping
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-plugins", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var imageURL string

	log.Printf("Scraping image URL from article: %s", articleURL)

	err := chromedp.Run(ctx,
		// Navigate to the article
		chromedp.Navigate(articleURL),

		// Wait for page to load
		chromedp.WaitVisible("body", chromedp.ByQuery),

		// Wait for download button to be present (with timeout)
		chromedp.WaitVisible(".download-button", chromedp.ByQuery),

		// Extract the image URL from download button
		chromedp.Evaluate(`
			(() => {
				const downloadButton = document.querySelector('.download-button a');
				if (downloadButton && downloadButton.href) {
					console.log('Found download button with URL:', downloadButton.href);
					return downloadButton.href;
				}
				
				// Alternative selector in case structure is different
				const altButton = document.querySelector('button.download-button a');
				if (altButton && altButton.href) {
					console.log('Found alternative download button with URL:', altButton.href);
					return altButton.href;
				}
				
				console.log('No download button found');
				return '';
			})();
		`, &imageURL),
	)

	if err != nil {
		return "", fmt.Errorf("failed to scrape article image: %v", err)
	}

	if imageURL == "" {
		log.Printf("⚠️  No image URL found in download button for article: %s", articleURL)
		return "", fmt.Errorf("no image URL found in download button")
	}

	// Log the scraped image URL
	log.Printf("✅ Scraped image URL: %s", imageURL)

	return imageURL, nil
}

func ScrapeSingleArticle() {
	go func() {
		scrapeSingleArticleChan := make(chan events.DataEvent)
		events.EB.Subscribe("SCRAPE_SINGLE_ARTICLE", scrapeSingleArticleChan)

		for {
			start := time.Now()

			scrapedArticleEvent := <-scrapeSingleArticleChan
			scrapedArticle, ok := scrapedArticleEvent.Data.(ScrapedArticle)
			if !ok {
				log.Printf("Invalid articleData type received: %T", scrapedArticle)
				continue
			}
			if scrapedArticle.AuthorName == "" || scrapedArticle.Title == "" {
				log.Printf("Article is has no title or author's name ")
				continue
			}
			log.Printf("Scraping single article in progress %s:", scrapedArticle.Title)

			imageURL, err := ScrapeSingleArticleImage(scrapedArticle.URL)
			if err != nil {
				log.Printf("Error scraping article: %v", err)
				continue
			}

			scrapedArticle.ImageUrl = imageURL

			log.Printf("scrapedArticle.ImageUrl: %s", scrapedArticle.ImageUrl)

			events.EB.Publish("SAVE_SCRAPED_ARTICLES", scrapedArticle)

			log.Println("Successfully createdArticle: ", scrapedArticle.Title)

			fmt.Printf(
				"Total Scraping Duration : %s\n",
				time.Since(start),
			)
		}
	}()
}

func init() {
	log.Println("App initialized. initialized ScrapeSingleArticle()")
	go func() {
		ScrapeSingleArticle()
	}()
}
