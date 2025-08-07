package articles

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Tibz-Dankan/hackernoon-articles/internal/events"
	"github.com/chromedp/chromedp"
)

type ScrapedArticle struct {
	Title           string
	URL             string
	ImageUrl        string
	PostedAt        time.Time
	AuthorName      string
	AuthorPageURL   string
	AuthorAvatarUrl string
	Summary         string
	Tag             string   // Single tag from the tag div
	Tags            []string // Keep for backward compatibility
	ReadDuration    string   // Read duration like "4m", "2h", etc.
}

type HackerNoonScraper struct {
	ctx context.Context
}

func NewHackerNoonScraper() *HackerNoonScraper {
	// Create chrome context with options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, _ := chromedp.NewContext(allocCtx)

	return &HackerNoonScraper{ctx: ctx}
}

func (h *HackerNoonScraper) ScrapeBitcoinArticles(maxArticles int, scrolls int) ([]ScrapedArticle, error) {
	var htmlContent string
	var articles []ScrapedArticle

	log.Println("Navigating to Hacker Noon Bitcoin articles...")

	err := chromedp.Run(h.ctx,
		// Navigate to Bitcoin tagged articles
		chromedp.Navigate("https://hackernoon.com/tagged/bitcoin"),

		// Wait for the page to load
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Sleep(3*time.Second),

		// Perform infinite scrolling to load more articles
		// h.performInfiniteScroll(scrolls, maxArticles),
		h.performInfiniteScrollOptimized(scrolls, maxArticles),
		// Get the final HTML content
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scrape Hacker Noon: %v", err)
	}

	// Parse HTML with goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	fmt.Println("Extracting article information...")

	// Find the infinite scroll container and extract articles
	doc.Find(".infinite-scroll-component article").Each(func(i int, s *goquery.Selection) {
		if len(articles) >= maxArticles {
			return
		}

		article := h.extractArticleData(s)
		if article.Title != "" {
			articles = append(articles, article)
			log.Printf("Found article %d: %s\n", len(articles), article.Title)
		}
	})

	return articles, nil
}

func (h *HackerNoonScraper) extractArticleData(s *goquery.Selection) ScrapedArticle {
	article := ScrapedArticle{}

	// Extract title from title-wrapper h2 a
	titleLink := s.Find(".title-wrapper h2 a").First()
	article.Title = strings.TrimSpace(titleLink.Text())

	// Extract article URL from title link
	if href, exists := titleLink.Attr("href"); exists {
		if strings.HasPrefix(href, "/") {
			article.URL = "https://hackernoon.com" + href
		} else if strings.HasPrefix(href, "https://hackernoon.com") {
			article.URL = href
		} else if strings.HasPrefix(href, "http") {
			article.URL = href
		}
	}

	// Extract image URL from image-wrapper a span img src
	imageLink := s.Find(".image-wrapper a span img").First()
	if src, exists := imageLink.Attr("src"); exists && src != "" {
		if strings.Contains(src, "http") {
			article.ImageUrl = src
		} else if strings.HasPrefix(src, "/") {
			article.ImageUrl = "https://hackernoon.com" + src
		}
	}

	// Extract author information from card-info .author .author-info
	authorInfo := s.Find(".card-info .author .author-info").First()

	// Extract author name and page URL from author-link
	authorLink := authorInfo.Find("a.author-link").First()
	article.AuthorName = strings.TrimSpace(authorLink.Text())

	if href, exists := authorLink.Attr("href"); exists {
		if strings.HasPrefix(href, "/") {
			article.AuthorPageURL = "https://hackernoon.com" + href
		} else if strings.HasPrefix(href, "https://hackernoon.com") {
			article.AuthorPageURL = href
		} else if strings.HasPrefix(href, "http") {
			article.AuthorPageURL = href
		}
	}

	// Extract author avatar from author section span img
	authorAvatar := s.Find(".card-info .author span img").First()
	if src, exists := authorAvatar.Attr("src"); exists && src != "" {
		if strings.Contains(src, "http") {
			article.AuthorAvatarUrl = src
		} else if strings.HasPrefix(src, "/") {
			article.AuthorAvatarUrl = "https://hackernoon.com" + src
		}
	}

	// Set default avatar if none found
	if article.AuthorAvatarUrl == "" {
		article.AuthorAvatarUrl = "https://hackernoon.com/default-avatar.png"
	}

	// Extract publish date and read duration from .author-info .date
	dateDiv := authorInfo.Find(".date").First()

	// First extract the read duration from the inner div
	readDurationDiv := dateDiv.Find("div").First()
	article.ReadDuration = strings.TrimSpace(readDurationDiv.Text())

	// Get only the direct text content of dateDiv (excluding inner divs)
	var dateText string
	dateDiv.Contents().Each(func(i int, s *goquery.Selection) {
		// Only get text nodes (not element nodes)
		if goquery.NodeName(s) == "#text" {
			dateText += s.Text()
		}
	})
	dateText = strings.TrimSpace(dateText)

	log.Println("readDuration:", article.ReadDuration)
	log.Println("dateText:", dateText)

	if dateText != "" {
		if parsedTime, err := h.parseDateTime(dateText); err == nil {
			article.PostedAt = parsedTime
			log.Println("parsed date:", parsedTime)
		} else {
			log.Println("date parsing error:", err)
		}
	}

	// If no date found, use current time as fallback
	if article.PostedAt.IsZero() {
		article.PostedAt = time.Now()
		log.Println("using fallback date:", article.PostedAt)
	}

	// Extract tag from image-wrapper .tag a
	tagLink := s.Find(".image-wrapper .tag a").First()
	article.Tag = strings.TrimSpace(tagLink.Text())

	// Also add to Tags array for backward compatibility
	if article.Tag != "" {
		article.Tags = append(article.Tags, article.Tag)
	}

	// Extract summary - this might need adjustment based on actual structure
	// Since you didn't mention summary location, keeping flexible approach
	summarySelectors := []string{".summary", ".description", ".excerpt", ".snippet"}
	for _, selector := range summarySelectors {
		if summary := s.Find(selector).First().Text(); summary != "" && len(summary) > 50 {
			article.Summary = strings.TrimSpace(summary)
			if len(article.Summary) > 300 {
				article.Summary = article.Summary[:300] + "..."
			}
			break
		}
	}

	return article
}

func (h *HackerNoonScraper) parseDateTime(dateStr string) (time.Time, error) {
	formats := []string{
		"Jan 2, 2006",     // Primary format based on your example
		"January 2, 2006", // Full month name variant
		time.RFC3339,
		time.RFC822,
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"02/01/2006",
		"01/02/2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

// func (h *HackerNoonScraper) performInfiniteScroll(scrolls int, maxArticles int) chromedp.Action {
// 	return chromedp.ActionFunc(func(ctx context.Context) error {
// 		log.Printf("Starting infinite scroll with manual image loading...\n")

// 		for i := 0; i < scrolls; i++ {
// 			// Check current number of articles
// 			var articleCount int
// 			err := chromedp.Evaluate(`
// 				document.querySelectorAll('.infinite-scroll-component article').length
// 			`, &articleCount).Do(ctx)
// 			if err == nil && articleCount >= maxArticles {
// 				log.Printf("Found enough articles (%d), stopping scroll\n", articleCount)
// 				break
// 			}

// 			log.Printf("Before scroll %d: %d articles found\n", i+1, articleCount)

// 			// Load images for current articles
// 			err = h.loadAllCurrentImages().Do(ctx)
// 			if err != nil {
// 				log.Printf("Error loading images: %v\n", err)
// 			}

// 			// Scroll to load more articles
// 			err = chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight)`, nil).Do(ctx)
// 			if err != nil {
// 				return fmt.Errorf("scroll failed: %v", err)
// 			}

// 			// Wait for new articles to load
// 			time.Sleep(3 * time.Second)

// 			var newArticleCount int
// 			chromedp.Evaluate(`
// 				document.querySelectorAll('.infinite-scroll-component article').length
// 			`, &newArticleCount).Do(ctx)

// 			log.Printf("After scroll %d: %d articles found (added %d new)\n", i+1, newArticleCount, newArticleCount-articleCount)
// 		}

// 		// Final image loading pass
// 		log.Println("Final image loading pass...")
// 		err := h.loadAllCurrentImages().Do(ctx)
// 		if err != nil {
// 			log.Printf("Error in final image loading: %v\n", err)
// 		}

// 		return nil
// 	})
// }

// // Load all images currently on the page with 100% reliability
// func (h *HackerNoonScraper) loadAllCurrentImages() chromedp.Action {
// 	return chromedp.ActionFunc(func(ctx context.Context) error {
// 		log.Println("Loading all current images with 100% reliability...")

// 		// First, get total image count
// 		var totalImages int
// 		err := chromedp.Evaluate(`
// 			(() => {
// 				return document.querySelectorAll('.infinite-scroll-component img[data-nimg]').length;
// 			})();
// 		`, &totalImages).Do(ctx)
// 		if err != nil {
// 			return fmt.Errorf("error counting images: %v", err)
// 		}

// 		log.Printf("Found %d total images to load\n", totalImages)

// 		// Load images in batches by article
// 		var articleCount int
// 		err = chromedp.Evaluate(`
// 			(() => {
// 				return document.querySelectorAll('.infinite-scroll-component article').length;
// 			})();
// 		`, &articleCount).Do(ctx)
// 		if err != nil {
// 			return fmt.Errorf("error counting articles: %v", err)
// 		}

// 		log.Printf("Processing %d articles for image loading\n", articleCount)

// 		// Process each article individually for maximum reliability
// 		for articleIndex := 0; articleIndex < articleCount; articleIndex++ {
// 			log.Printf("Loading images for article %d/%d\n", articleIndex+1, articleCount)

// 			// Load images for this specific article with persistence
// 			err = h.loadImagesForSingleArticle(articleIndex).Do(ctx)
// 			if err != nil {
// 				log.Printf("Error loading images for article %d: %v\n", articleIndex+1, err)
// 			}

// 			// Small delay between articles to prevent overwhelming
// 			time.Sleep(100 * time.Millisecond)
// 		}

// 		// Final verification - keep trying until ALL images are loaded
// 		log.Println("Final verification: ensuring 100% image loading...")
// 		maxGlobalAttempts := 30 // Up to 15 seconds of additional attempts

// 		for attempt := 0; attempt < maxGlobalAttempts; attempt++ {
// 			var stats map[string]interface{}
// 			err = chromedp.Evaluate(`
// 				(() => {
// 					const allImages = document.querySelectorAll('.infinite-scroll-component img[data-nimg]');
// 					const stats = {
// 						total: allImages.length,
// 						loaded: 0,
// 						placeholder: 0
// 					};

// 					allImages.forEach(img => {
// 						if (img.src.startsWith('data:image/gif')) {
// 							stats.placeholder++;
// 						} else if (img.src.startsWith('http')) {
// 							stats.loaded++;
// 						}
// 					});

// 					return stats;
// 				})();
// 			`, &stats).Do(ctx)

// 			if err == nil {
// 				loaded := int(stats["loaded"].(float64))
// 				placeholder := int(stats["placeholder"].(float64))
// 				total := int(stats["total"].(float64))

// 				log.Printf("Global verification attempt %d: %d/%d loaded, %d still placeholder\n", attempt+1, loaded, total, placeholder)

// 				// Check if ALL images are loaded (100%)
// 				if placeholder == 0 {
// 					log.Printf("SUCCESS: 100%% image loading achieved! (%d/%d)\n", loaded, total)
// 					break
// 				}
// 			}

// 			// Aggressively try to load remaining placeholder images
// 			chromedp.Evaluate(`
// 				(() => {
// 					const remainingPlaceholders = document.querySelectorAll('img[data-nimg][src^="data:image/gif"]');
// 					console.log('Forcing load for', remainingPlaceholders.length, 'remaining images');

// 					remainingPlaceholders.forEach((img, index) => {
// 						// Multiple aggressive strategies
// 						img.scrollIntoView({behavior: 'auto', block: 'center'});
// 						img.loading = 'eager';

// 						// Try to trigger intersection observer
// 						const rect = img.getBoundingClientRect();
// 						window.scrollTo(0, window.pageYOffset + rect.top - window.innerHeight/2);

// 						// Dispatch events
// 						img.dispatchEvent(new Event('load'));
// 						img.dispatchEvent(new Event('scroll'));
// 					});

// 					// Global triggers
// 					window.dispatchEvent(new Event('scroll'));
// 					window.dispatchEvent(new Event('resize'));
// 					window.dispatchEvent(new Event('load'));
// 				})();
// 			`, nil).Do(ctx)

// 			time.Sleep(500 * time.Millisecond)
// 		}

// 		// Final report
// 		var finalStats map[string]interface{}
// 		chromedp.Evaluate(`
// 			(() => {
// 				const allImages = document.querySelectorAll('.infinite-scroll-component img[data-nimg]');
// 				const stats = {
// 					total: allImages.length,
// 					loaded: 0,
// 					placeholder: 0
// 				};

// 				allImages.forEach(img => {
// 					if (img.src.startsWith('data:image/gif')) {
// 						stats.placeholder++;
// 					} else if (img.src.startsWith('http')) {
// 						stats.loaded++;
// 					}
// 				});

// 				return stats;
// 			})();
// 		`, &finalStats).Do(ctx)

// 		if finalStats != nil {
// 			loaded := int(finalStats["loaded"].(float64))
// 			total := int(finalStats["total"].(float64))
// 			placeholder := int(finalStats["placeholder"].(float64))

// 			successRate := float64(loaded) / float64(total) * 100
// 			log.Printf("FINAL RESULT: %.1f%% success rate (%d/%d loaded, %d failed)\n", successRate, loaded, total, placeholder)

// 			if placeholder > 0 {
// 				log.Printf("WARNING: %d images still have placeholder sources\n", placeholder)
// 			}
// 		}

// 		return nil
// 	})
// }

// // Load images for a single article with maximum persistence
// func (h *HackerNoonScraper) loadImagesForSingleArticle(articleIndex int) chromedp.Action {
// 	return chromedp.ActionFunc(func(ctx context.Context) error {
// 		// Get image count for this article
// 		var imageCount int
// 		err := chromedp.Evaluate(fmt.Sprintf(`
// 			(() => {
// 				const article = document.querySelectorAll('.infinite-scroll-component article')[%d];
// 				if (!article) return 0;
// 				return article.querySelectorAll('img[data-nimg]').length;
// 			})();
// 		`, articleIndex), &imageCount).Do(ctx)

// 		if err != nil || imageCount == 0 {
// 			return nil // Skip if no images or error
// 		}

// 		log.Printf("Article %d has %d images to load\n", articleIndex+1, imageCount)

// 		// Scroll article into view first
// 		err = chromedp.Evaluate(fmt.Sprintf(`
// 			(() => {
// 				const article = document.querySelectorAll('.infinite-scroll-component article')[%d];
// 				if (article) {
// 					article.scrollIntoView({behavior: 'auto', block: 'center'});
// 				}
// 			})();
// 		`, articleIndex), nil).Do(ctx)
// 		if err != nil {
// 			return err
// 		}

// 		time.Sleep(500 * time.Millisecond)

// 		// Now wait for each image in this article to load
// 		maxAttempts := 30 // 15 seconds max per article
// 		for attempt := 0; attempt < maxAttempts; attempt++ {
// 			var articleStats map[string]interface{}
// 			err = chromedp.Evaluate(fmt.Sprintf(`
// 				(() => {
// 					const article = document.querySelectorAll('.infinite-scroll-component article')[%d];
// 					if (!article) return {total: 0, loaded: 0, placeholder: 0};

// 					const articleImages = article.querySelectorAll('img[data-nimg]');
// 					const stats = {
// 						total: articleImages.length,
// 						loaded: 0,
// 						placeholder: 0
// 					};

// 					articleImages.forEach(img => {
// 						if (img.src.startsWith('data:image/gif')) {
// 							stats.placeholder++;
// 						} else if (img.src.startsWith('http')) {
// 							stats.loaded++;
// 						}
// 					});

// 					return stats;
// 				})();
// 			`, articleIndex), &articleStats).Do(ctx)

// 			if err == nil {
// 				loaded := int(articleStats["loaded"].(float64))
// 				placeholder := int(articleStats["placeholder"].(float64))
// 				total := int(articleStats["total"].(float64))

// 				if total > 0 {
// 					if placeholder == 0 {
// 						log.Printf("Article %d: ALL images loaded (%d/%d)\n", articleIndex+1, loaded, total)
// 						break
// 					} else {
// 						log.Printf("Article %d progress: %d/%d loaded, %d remaining\n", articleIndex+1, loaded, total, placeholder)
// 					}
// 				}
// 			}

// 			// Aggressively load remaining images in this article
// 			chromedp.Evaluate(fmt.Sprintf(`
// 				(() => {
// 					const article = document.querySelectorAll('.infinite-scroll-component article')[%d];
// 					if (!article) return;

// 					const placeholderImgs = article.querySelectorAll('img[data-nimg][src^="data:image/gif"]');
// 					placeholderImgs.forEach((img, imgIndex) => {
// 						// Multiple loading strategies
// 						img.scrollIntoView({behavior: 'auto', block: 'center'});
// 						img.loading = 'eager';

// 						// Force intersection
// 						const rect = img.getBoundingClientRect();
// 						if (rect.top < window.innerHeight && rect.bottom > 0) {
// 							img.dispatchEvent(new Event('load'));
// 						}
// 					});

// 					// Trigger global events
// 					window.dispatchEvent(new Event('scroll'));
// 					window.dispatchEvent(new Event('resize'));
// 				})();
// 			`, articleIndex), nil).Do(ctx)

// 			time.Sleep(500 * time.Millisecond)
// 		}

// 		return nil
// 	})
// }

// BatchImageResult tracks the success of image loading for a batch
type BatchImageResult struct {
	BatchIndex      int
	ArticlesInBatch []int
	LoadedImages    int
	FailedImages    int
	Success         bool
}

func (h *HackerNoonScraper) performInfiniteScrollOptimized(scrolls int, maxArticles int) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		log.Printf("Starting optimized infinite scroll for up to %d articles...\n", maxArticles)

		var failedBatches []BatchImageResult
		var mu sync.Mutex // Protect failedBatches slice

		batchIndex := 0

		for i := 0; i < scrolls; i++ {
			// Check current number of articles
			var currentArticleCount int
			err := chromedp.Evaluate(`
				document.querySelectorAll('.infinite-scroll-component article').length
			`, &currentArticleCount).Do(ctx)
			if err == nil && currentArticleCount >= maxArticles {
				log.Printf("Reached target articles (%d), stopping scroll\n", currentArticleCount)
				break
			}

			previousCount := currentArticleCount
			log.Printf("Scroll %d: Starting with %d articles\n", i+1, previousCount)

			// Perform scroll to load more articles
			err = chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight)`, nil).Do(ctx)
			if err != nil {
				return fmt.Errorf("scroll failed: %v", err)
			}

			// Wait for new articles to load
			time.Sleep(3 * time.Second)

			// Get new article count
			var newArticleCount int
			chromedp.Evaluate(`
				document.querySelectorAll('.infinite-scroll-component article').length
			`, &newArticleCount).Do(ctx)

			newArticlesAdded := newArticleCount - previousCount
			log.Printf("Scroll %d: Loaded %d new articles (total: %d)\n", i+1, newArticlesAdded, newArticleCount)

			// If new articles were added, load their images immediately
			if newArticlesAdded > 0 {
				batchIndex++

				// Create batch info for the newly loaded articles
				var newArticleIndices []int
				for idx := previousCount; idx < newArticleCount; idx++ {
					newArticleIndices = append(newArticleIndices, idx)
				}

				log.Printf("Batch %d: Loading images for articles %d-%d (%d articles)\n",
					batchIndex, previousCount, newArticleCount-1, newArticlesAdded)

				// Load images for this batch
				batchResult := h.loadImagesBatch(ctx, newArticleIndices, batchIndex)

				// If batch failed, store it for retry later
				if !batchResult.Success {
					mu.Lock()
					failedBatches = append(failedBatches, batchResult)
					mu.Unlock()
					log.Printf("Batch %d: Image loading incomplete, will retry after all articles loaded\n", batchIndex)
				} else {
					log.Printf("Batch %d: All images loaded successfully!\n", batchIndex)
				}
			}

			// Small delay before next scroll
			time.Sleep(1 * time.Second)
		}

		// After all articles are loaded, retry failed batches
		if len(failedBatches) > 0 {
			log.Printf("\n=== RETRY PHASE ===\n")
			log.Printf("Retrying image loading for %d failed batches...\n", len(failedBatches))

			err := h.retryFailedBatches(ctx, failedBatches)
			if err != nil {
				log.Printf("Error during retry phase: %v\n", err)
			}
		}

		// Final comprehensive check
		log.Printf("\n=== FINAL VERIFICATION ===\n")
		err := h.performFinalImageVerification(ctx)
		if err != nil {
			log.Printf("Error during final verification: %v\n", err)
		}

		return nil
	})
}

// Load images for a specific batch of articles
func (h *HackerNoonScraper) loadImagesBatch(ctx context.Context, articleIndices []int, batchIndex int) BatchImageResult {
	result := BatchImageResult{
		BatchIndex:      batchIndex,
		ArticlesInBatch: articleIndices,
	}

	if len(articleIndices) == 0 {
		result.Success = true
		return result
	}

	log.Printf("Processing batch %d with %d articles...\n", batchIndex, len(articleIndices))

	// First, scroll the first article of this batch into view
	err := chromedp.Evaluate(fmt.Sprintf(`
		(() => {
			const article = document.querySelectorAll('.infinite-scroll-component article')[%d];
			if (article) {
				article.scrollIntoView({behavior: 'auto', block: 'start'});
			}
		})();
	`, articleIndices[0]), nil).Do(ctx)
	if err != nil {
		log.Printf("Error scrolling batch into view: %v\n", err)
	}

	time.Sleep(500 * time.Millisecond)

	// Use goroutines to process articles in parallel (with controlled concurrency)
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 3) // Limit to 3 concurrent articles
	var totalLoaded, totalFailed int
	var mu sync.Mutex

	for _, articleIndex := range articleIndices {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			loaded, failed := h.loadSingleArticleImagesOptimized(ctx, idx)

			mu.Lock()
			totalLoaded += loaded
			totalFailed += failed
			mu.Unlock()
		}(articleIndex)
	}

	wg.Wait()

	result.LoadedImages = totalLoaded
	result.FailedImages = totalFailed
	result.Success = totalFailed == 0

	log.Printf("Batch %d complete: %d images loaded, %d failed\n", batchIndex, totalLoaded, totalFailed)
	return result
}

// Load images for a single article with optimized approach
func (h *HackerNoonScraper) loadSingleArticleImagesOptimized(ctx context.Context, articleIndex int) (loaded, failed int) {
	// Get image count for this article
	var imageCount int
	err := chromedp.Evaluate(fmt.Sprintf(`
		(() => {
			const article = document.querySelectorAll('.infinite-scroll-component article')[%d];
			if (!article) return 0;
			return article.querySelectorAll('img[data-nimg]').length;
		})();
	`, articleIndex), &imageCount).Do(ctx)

	if err != nil || imageCount == 0 {
		return 0, 0
	}

	// Scroll this article into view
	chromedp.Evaluate(fmt.Sprintf(`
		(() => {
			const article = document.querySelectorAll('.infinite-scroll-component article')[%d];
			if (article) {
				article.scrollIntoView({behavior: 'auto', block: 'center'});
			}
		})();
	`, articleIndex), nil).Do(ctx)

	time.Sleep(200 * time.Millisecond)

	// Trigger image loading for this article
	chromedp.Evaluate(fmt.Sprintf(`
		(() => {
			const article = document.querySelectorAll('.infinite-scroll-component article')[%d];
			if (!article) return;
			
			const images = article.querySelectorAll('img[data-nimg]');
			images.forEach(img => {
				img.scrollIntoView({behavior: 'auto', block: 'nearest'});
				img.loading = 'eager';
				img.dispatchEvent(new Event('load'));
			});
		})();
	`, articleIndex), nil).Do(ctx)

	// Wait a bit for images to start loading
	time.Sleep(1 * time.Second)

	// Check results - single attempt for batch processing
	var stats map[string]interface{}
	err = chromedp.Evaluate(fmt.Sprintf(`
		(() => {
			const article = document.querySelectorAll('.infinite-scroll-component article')[%d];
			if (!article) return {total: 0, loaded: 0, placeholder: 0};
			
			const articleImages = article.querySelectorAll('img[data-nimg]');
			const stats = {
				total: articleImages.length,
				loaded: 0,
				placeholder: 0
			};
			
			articleImages.forEach(img => {
				if (img.src.startsWith('data:image/gif')) {
					stats.placeholder++;
				} else if (img.src.startsWith('http')) {
					stats.loaded++;
				}
			});
			
			return stats;
		})();
	`, articleIndex), &stats).Do(ctx)

	if err == nil && stats != nil {
		loadedCount := int(stats["loaded"].(float64))
		placeholderCount := int(stats["placeholder"].(float64))
		return loadedCount, placeholderCount
	}

	return 0, imageCount
}

// Retry failed batches with more aggressive approach
func (h *HackerNoonScraper) retryFailedBatches(ctx context.Context, failedBatches []BatchImageResult) error {
	log.Printf("Starting retry phase for %d failed batches...\n", len(failedBatches))

	for _, batch := range failedBatches {
		log.Printf("Retrying batch %d (%d articles)...\n", batch.BatchIndex, len(batch.ArticlesInBatch))

		// Use more aggressive loading for failed batches
		for _, articleIndex := range batch.ArticlesInBatch {
			err := h.aggressivelyLoadArticleImages(ctx, articleIndex)
			if err != nil {
				log.Printf("Error in aggressive loading for article %d: %v\n", articleIndex, err)
			}
		}
	}

	return nil
}

// Aggressively load images for articles that failed in batch processing
func (h *HackerNoonScraper) aggressivelyLoadArticleImages(ctx context.Context, articleIndex int) error {
	maxAttempts := 10 // More focused attempts

	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Scroll article into view
		chromedp.Evaluate(fmt.Sprintf(`
			(() => {
				const article = document.querySelectorAll('.infinite-scroll-component article')[%d];
				if (article) {
					article.scrollIntoView({behavior: 'auto', block: 'center'});
				}
			})();
		`, articleIndex), nil).Do(ctx)

		time.Sleep(300 * time.Millisecond)

		// Check current status
		var stats map[string]interface{}
		err := chromedp.Evaluate(fmt.Sprintf(`
			(() => {
				const article = document.querySelectorAll('.infinite-scroll-component article')[%d];
				if (!article) return {total: 0, loaded: 0, placeholder: 0};
				
				const articleImages = article.querySelectorAll('img[data-nimg]');
				const stats = {
					total: articleImages.length,
					loaded: 0,
					placeholder: 0
				};
				
				articleImages.forEach(img => {
					if (img.src.startsWith('data:image/gif')) {
						stats.placeholder++;
					} else if (img.src.startsWith('http')) {
						stats.loaded++;
					}
				});
				
				return stats;
			})();
		`, articleIndex), &stats).Do(ctx)

		if err == nil && stats != nil {
			placeholder := int(stats["placeholder"].(float64))
			if placeholder == 0 {
				// All images loaded for this article
				break
			}
		}

		// Aggressively trigger loading
		chromedp.Evaluate(fmt.Sprintf(`
			(() => {
				const article = document.querySelectorAll('.infinite-scroll-component article')[%d];
				if (!article) return;
				
				const placeholderImgs = article.querySelectorAll('img[data-nimg][src^="data:image/gif"]');
				placeholderImgs.forEach(img => {
					img.scrollIntoView({behavior: 'auto', block: 'center'});
					img.loading = 'eager';
					img.dispatchEvent(new Event('load'));
					img.dispatchEvent(new Event('scroll'));
				});
				
				window.dispatchEvent(new Event('scroll'));
				window.dispatchEvent(new Event('resize'));
			})();
		`, articleIndex), nil).Do(ctx)

		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

// Perform final verification of all images
func (h *HackerNoonScraper) performFinalImageVerification(ctx context.Context) error {
	log.Println("Performing final comprehensive image verification...")

	var finalStats map[string]interface{}
	err := chromedp.Evaluate(`
		(() => {
			const allImages = document.querySelectorAll('.infinite-scroll-component img[data-nimg]');
			const stats = {
				total: allImages.length,
				loaded: 0,
				placeholder: 0
			};
			
			allImages.forEach(img => {
				if (img.src.startsWith('data:image/gif')) {
					stats.placeholder++;
				} else if (img.src.startsWith('http')) {
					stats.loaded++;
				}
			});
			
			return stats;
		})();
	`, &finalStats).Do(ctx)

	if err == nil && finalStats != nil {
		loaded := int(finalStats["loaded"].(float64))
		total := int(finalStats["total"].(float64))
		placeholder := int(finalStats["placeholder"].(float64))

		successRate := float64(loaded) / float64(total) * 100
		log.Printf("FINAL RESULTS:\n")
		log.Printf("  Total images: %d\n", total)
		log.Printf("  Successfully loaded: %d\n", loaded)
		log.Printf("  Still placeholder: %d\n", placeholder)
		log.Printf("  Success rate: %.1f%%\n", successRate)

		if placeholder > 0 {
			log.Printf("WARNING: %d images still have placeholder sources\n", placeholder)
		} else {
			log.Printf("SUCCESS: All images loaded successfully!\n")
		}
	}

	return nil
}

func (h *HackerNoonScraper) Close() {
	chromedp.Cancel(h.ctx)
}

// Save scraped articles to JSON file
func saveToJSON(articles []ScrapedArticle) error {
	// Create filename with current date
	now := time.Now()
	filename := fmt.Sprintf("%s-hackernoon-bitcoin-articles.json", now.Format("20060102-150405"))

	// Create JSON data
	data := map[string]interface{}{
		"scraped_at":     now.Format("2006-01-02T15:04:05Z"),
		"total_articles": len(articles),
		"source":         "hackernoon.com",
		"category":       "bitcoin",
		"articles":       articles,
	}

	// Marshal to JSON with indentation
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Write to file
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	for _, article := range articles {
		events.EB.Publish("SAVE_SCRAPED_ARTICLES", article)
	}

	log.Printf("✅ Saved %d articles to %s", len(articles), filename)
	return nil
}

// Main scraping function that handles both JSON export and event publishing
func (h *HackerNoonScraper) ScrapeAndSave(maxArticles int, scrolls int) error {
	// Scrape articles - removed Bitcoin filtering since we're already on Bitcoin page
	articles, err := h.ScrapeBitcoinArticles(maxArticles, scrolls)
	if err != nil {
		return fmt.Errorf("scraping failed: %v", err)
	}

	if len(articles) == 0 {
		log.Println("⚠️  No articles found")
		return nil
	}

	// Save to JSON file and publish to event bus
	return saveToJSON(articles)
}

// Main scraping function that can be called from your application
func ScrapeHackerNoonBitcoinArticles(maxArticles, scrolls int) error {
	scraper := NewHackerNoonScraper()
	defer scraper.Close()

	log.Println("=== Hacker Noon Bitcoin Articles Scraper ===")
	log.Println("Starting JavaScript-aware scraping of Hacker Noon...")

	return scraper.ScrapeAndSave(maxArticles, scrolls)
}

// Alternative function that returns articles without saving (for testing)
func ScrapeHackerNoonBitcoinArticlesOnly(maxArticles, scrolls int) ([]ScrapedArticle, error) {
	scraper := NewHackerNoonScraper()
	defer scraper.Close()

	log.Println("=== Hacker Noon Bitcoin Articles Scraper (Data Only) ===")
	log.Println("Starting JavaScript-aware scraping of Hacker Noon...")

	return scraper.ScrapeBitcoinArticles(maxArticles, scrolls)
}

// func init() {
// 	log.Println("App initialized. Scheduling ScrapeHackerNoonBitcoinArticles() to run in 15 seconds...")

// 	go func() {
// 		time.Sleep(15 * time.Second)
// 		start := time.Now()
// 		ScrapeHackerNoonBitcoinArticles(20, 5)
// 		// ScrapeHackerNoonBitcoinArticles(200, 15)
// 		// ScrapeHackerNoonBitcoinArticles(2000, 120)
// 		// ScrapeHackerNoonBitcoinArticles(6500, 250)
// 		fmt.Printf(
// 			"Total Scraping Duration : %s\n",
// 			time.Since(start),
// 		)
// 	}()

// 	// log.Println("App initialized. Scheduling ScrapeHackerNoonBitcoinArticles() to run in  2 minutes...")

// 	// go func() {
// 	// 	time.Sleep(2 * time.Minute)
// 	// 	ScrapeHackerNoonBitcoinArticles(200, 24)
// 	// }()
// }

// func init() {
// 	fmt.Println("App initialized. Scheduling ScrapeHackerNoonBitcoinArticles() to run daily at 6pm EAT...")
// 	go func() {
// 		// Set timezone to East Africa Time (EAT - UTC+3)
// 		eat, err := time.LoadLocation("Africa/Nairobi")
// 		if err != nil {
// 			log.Printf("Error loading EAT timezone: %v\n", err)
// 			return
// 		}

// 		// Calculate initial delay to next 6pm EAT
// 		now := time.Now().In(eat)
// 		next6pm := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, eat)

// 		// If it's already past 6pm today, schedule for 6pm tomorrow
// 		if now.After(next6pm) || now.Equal(next6pm) {
// 			next6pm = next6pm.Add(24 * time.Hour)
// 		}

// 		initialDelay := next6pm.Sub(now)
// 		log.Printf("Initial execution scheduled for: %v (in %v)\n", next6pm.Format("2006-01-02 15:04:05 MST"), initialDelay)

// 		// Wait for initial 6pm
// 		initialTimer := time.NewTimer(initialDelay)
// 		<-initialTimer.C

// 		// Execute first time
// 		log.Printf("Executing ScrapeHackerNoonBitcoinArticles() at %v\n", time.Now().In(eat).Format("2006-01-02 15:04:05 MST"))
// 		// ScrapeHackerNoonBitcoinArticles(30, 8)
// 		ScrapeHackerNoonBitcoinArticles(200, 24)

// 		ticker := time.NewTicker(24 * time.Hour)
// 		defer ticker.Stop()

// 		log.Println("Daily ticker started. Next executions every 24 hours at 6pm EAT")

// 		// Execute daily at 6pm
// 		for range ticker.C {
// 			currentTime := time.Now().In(eat)
// 			log.Printf("Executing ScrapeHackerNoonBitcoinArticles() at %v\n", currentTime.Format("2006-01-02 15:04:05 MST"))
// 			ScrapeHackerNoonBitcoinArticles(200, 24)
// 		}
// 	}()
// }
