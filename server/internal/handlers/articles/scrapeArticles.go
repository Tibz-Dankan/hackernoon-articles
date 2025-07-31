package articles

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Tibz-Dankan/hackernoon-articles/internal/events"
	"github.com/chromedp/chromedp"
)

type ScrapedArticleData struct {
	Title           string
	URL             string
	ImageUrl        string
	PostedAt        time.Time
	AuthorName      string
	AuthorPageURL   string // Added author's page URL
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

func (h *HackerNoonScraper) ScrapeBitcoinArticles(maxArticles int, scrolls int) ([]ScrapedArticleData, error) {
	var htmlContent string
	var articles []ScrapedArticleData

	fmt.Println("Navigating to Hacker Noon Bitcoin articles...")

	err := chromedp.Run(h.ctx,
		// Navigate to Bitcoin tagged articles
		chromedp.Navigate("https://hackernoon.com/tagged/bitcoin"),

		// Wait for the page to load
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Sleep(3*time.Second),

		// Perform infinite scrolling to load more articles
		h.performInfiniteScroll(scrolls, maxArticles),

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
			fmt.Printf("Found article %d: %s\n", len(articles), article.Title)
		}
	})

	return articles, nil
}

func (h *HackerNoonScraper) extractArticleData(s *goquery.Selection) ScrapedArticleData {
	article := ScrapedArticleData{}

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

	// Extract image URL from image-wrapper a href
	// imageLink := s.Find(".image-wrapper a").First()
	imageLink := s.Find(".image-wrapper  a span img").First()
	if src, exists := imageLink.Attr("src"); exists && src != "" {
		if strings.Contains(src, "http") {
			article.ImageUrl = src
		} else if strings.HasPrefix(src, "/") {
			article.ImageUrl = "https://hackernoon.com" + src
		}
	}

	// If no image found, set a default placeholder
	if article.ImageUrl == "" {
		article.ImageUrl = "https://hackernoon.com/hn-logo.png"
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

	// Extract publish date from .author-info .date
	dateDiv := authorInfo.Find(".date").First()
	dateText := strings.TrimSpace(dateDiv.Text())

	if dateText != "" {
		if parsedTime, err := h.parseDateTime(dateText); err == nil {
			article.PostedAt = parsedTime
		}
	}

	// If no date found, use current time as fallback
	if article.PostedAt.IsZero() {
		article.PostedAt = time.Now()
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

func (h *HackerNoonScraper) performInfiniteScroll(scrolls int, maxArticles int) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		log.Printf("Starting infinite scroll to load more articles...\n")

		for i := 0; i < scrolls; i++ {
			// Check current number of articles in the infinite-scroll-component
			var articleCount int
			err := chromedp.Evaluate(`
				document.querySelectorAll('.infinite-scroll-component article').length
			`, &articleCount).Do(ctx)
			if err == nil && articleCount >= maxArticles {
				log.Printf("Found enough articles (%d), stopping scroll\n", articleCount)
				break
			}

			// Scroll to bottom to trigger auto-loading
			err = chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight)`, nil).Do(ctx)
			if err != nil {
				return fmt.Errorf("scroll failed: %v", err)
			}

			// Wait for new content to auto-load
			time.Sleep(3 * time.Second)

			log.Printf("Completed scroll %d/%d - Current articles: %d\n", i+1, scrolls, articleCount)
		}

		return nil
	})
}

func (h *HackerNoonScraper) Close() {
	chromedp.Cancel(h.ctx)
}

// Save scraped articles to JSON file
func saveToJSON(articles []ScrapedArticleData) error {
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
func ScrapeHackerNoonBitcoinArticlesOnly(maxArticles, scrolls int) ([]ScrapedArticleData, error) {
	scraper := NewHackerNoonScraper()
	defer scraper.Close()

	log.Println("=== Hacker Noon Bitcoin Articles Scraper (Data Only) ===")
	log.Println("Starting JavaScript-aware scraping of Hacker Noon...")

	return scraper.ScrapeBitcoinArticles(maxArticles, scrolls)
}

func init() {
	log.Println("App initialized. Scheduling ScrapeHackerNoonBitcoinArticles() to run in 15 seconds...")

	go func() {
		time.Sleep(15 * time.Second)
		ScrapeHackerNoonBitcoinArticles(200, 24)
	}()
}

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
