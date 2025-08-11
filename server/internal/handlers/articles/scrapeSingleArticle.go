package articles

import (
	"context"
	"fmt"
	"log"
	"strings"
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
	var finalURL string
	var statusCode int64
	var pageTitle string

	log.Printf("Scraping image URL from article: %s", articleURL)

	err := chromedp.Run(ctx,
		// Navigate to the article
		chromedp.Navigate(articleURL),

		// Wait for page to load
		chromedp.WaitVisible("body", chromedp.ByQuery),

		// Check the final URL after any redirects and get status code
		chromedp.Evaluate(`window.location.href`, &finalURL),
		chromedp.Evaluate(`
			(() => {
				// Try to get status from performance API
				const entries = performance.getEntriesByType('navigation');
				if (entries.length > 0) {
					return entries[0].responseStatus || 200;
				}
				return 200; // Default to 200 if we can't determine
			})();
		`, &statusCode),

		// Get page title to help detect 404 pages
		chromedp.Title(&pageTitle),
	)

	if err != nil {
		return "", fmt.Errorf("failed to navigate to article: %v", err)
	}

	// Check if URL redirected to a 404 page
	if strings.Contains(strings.ToLower(finalURL), "/404") {
		log.Printf("❌ Article URL redirected to 404 page: %s -> %s", articleURL, finalURL)
		return "", fmt.Errorf("article not found - redirected to 404 page: %s", finalURL)
	}

	// Check if page title indicates a 404
	lowerTitle := strings.ToLower(pageTitle)
	if strings.Contains(lowerTitle, "404") ||
		strings.Contains(lowerTitle, "not found") ||
		strings.Contains(lowerTitle, "page not found") {
		log.Printf("❌ Article appears to be 404 based on title: %s", pageTitle)
		return "", fmt.Errorf("article not found - page title indicates 404: %s", pageTitle)
	}

	// Check status code (though this might not always be reliable in Chrome)
	if statusCode >= 400 {
		log.Printf("❌ Article returned error status code: %d", statusCode)
		return "", fmt.Errorf("article returned error status code: %d", statusCode)
	}

	// Additional check: look for common 404 page indicators in the DOM
	var is404Page bool
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`
			(() => {
				const bodyText = document.body.innerText.toLowerCase();
				const headingText = document.querySelector('h1, h2, h3') ? 
					document.querySelector('h1, h2, h3').innerText.toLowerCase() : '';
				
				// Check for common 404 indicators
				const indicators = ['404', 'not found', 'page not found', 'page does not exist', 'oops'];
				
				// Check if any indicator is prominently displayed
				for (let indicator of indicators) {
					if (headingText.includes(indicator) || 
						(bodyText.includes(indicator) && bodyText.split(' ').length < 100)) {
						return true;
					}
				}
				
				// Check for HackerNoon specific 404 patterns
				if (document.querySelector('.error-page') || 
					document.querySelector('[class*="404"]') ||
					document.querySelector('[id*="404"]')) {
					return true;
				}
				
				return false;
			})();
		`, &is404Page),
	)

	if err == nil && is404Page {
		log.Printf("❌ Article appears to be a 404 page based on content analysis")
		return "", fmt.Errorf("article not found - page content indicates 404")
	}

	// Now proceed with the original image scraping logic
	err = chromedp.Run(ctx,
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
		// If we can't find the download button, it might be because this is a 404 page
		// that loaded but doesn't have the expected content structure
		log.Printf("❌ Failed to find download button - this might be a 404 page or invalid article: %v", err)
		return "", fmt.Errorf("failed to scrape article image (possibly 404 or invalid article): %v", err)
	}

	if imageURL == "" {
		log.Printf("⚠️  No image URL found in download button for article: %s", articleURL)
		return "", fmt.Errorf("no image URL found in download button")
	}

	// Final URL comparison log
	if finalURL != articleURL {
		log.Printf("ℹ️  URL redirected: %s -> %s", articleURL, finalURL)
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
			// log.Printf("Article to be scraped %+v:", scrapedArticle)
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

// func init() {
// 	log.Println("App initialized. initialized ScrapeSingleArticle()")
// 	go func() {
// 		ScrapeSingleArticle()
// 	}()
// }
