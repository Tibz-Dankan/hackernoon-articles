package articles

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"time"

	"github.com/Tibz-Dankan/hackernoon-articles/internal/constants"
	"github.com/Tibz-Dankan/hackernoon-articles/internal/events"
	"github.com/Tibz-Dankan/hackernoon-articles/internal/models"
)

type ScrapedData struct {
	ScrapedAt     string           `json:"scraped_at"`
	TotalArticles int              `json:"total_articles"`
	Source        string           `json:"source"`
	Category      string           `json:"category"`
	Articles      []ScrapedArticle `json:"articles"`
}

func ProcessArticles() error {
	filename, err := filepath.Abs("./20250810-114840-hackernoon-bitcoin-articles.json")
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
	fmt.Printf("Successfully loaded %d articles from %s\n\n", len(scrapedData.Articles), filename)

	for _, article := range scrapedData.Articles {
		events.EB.Publish("SAVE_SCRAPED_ARTICLES", article)
	}

	return nil
}

func ProcessArticlesWithoutImages() error {
	article := models.Article{}
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
	fmt.Printf("Successfully loaded %d articles from %s\n\n", len(scrapedData.Articles), filename)

	for _, currArticle := range scrapedData.Articles {
		if currArticle.ImageUrl == "" && currArticle.AuthorName != "" && currArticle.Title != "" {

			var isBlackListed bool
			for _, blArticle := range blackListArticles {
				if blArticle == currArticle.Title {
					log.Printf("BlackListed: %s", currArticle.Title)
					isBlackListed = true
				}
			}
			if isBlackListed {
				continue
			}

			savedArticle, err := article.FindByTitle(currArticle.Title)
			if err != nil && err.Error() != constants.RECORD_NOT_FOUND_ERROR {
				log.Printf("Error finding the saved article: %v", err)
				continue
			}
			if savedArticle.ID != "" {
				log.Printf("Article is already saved: %s ", savedArticle.Title)
				continue
			}

			log.Printf("Publishing article: %s", currArticle.Title)
			events.EB.Publish("SCRAPE_SINGLE_ARTICLE", currArticle)
		}
	}
	return nil
}

// func init() {
// 	log.Println("App initialized. Scheduling ProcessArticles() to run in 15 seconds...")

// 	go func() {
// 		time.Sleep(15 * time.Second)
// 		ProcessArticles()
// 	}()
// }

// To more into single article image selection logic
func init() {
	log.Println("App initialized. Scheduling ProcessArticlesWithoutImages() to run in 15 seconds...")

	go func() {
		time.Sleep(15 * time.Second)
		ProcessArticlesWithoutImages()
	}()
}

var blackListArticles = []string{
	"Decentralized Applications Will Take Cryptocurrency to the Mainstream",
	"What is The Bitcoin Halving and What Impact Will It Have on the Crypto Market?",
	"A (Very) Basic Intro To Elliptic Curve Cryptography",
	"How to generate a Bitcoin address — Technical address generation explanation",
	"Top Crypto Exchange and Blockchain Companies to Watch for in Canada: 2020 Edition",
	"Will Blockchain Produce a New Generation of Retail Algorithmic Traders?",
	"Rising From the Ashes — A Tale of Bitcoin Crashes",
	"Trace Mayer on Why You Must Own Your Bitcoin Private Keys",
	"The Weaknesses of Blockchain and Decentralization",
	"Is the Next Generation of Blockchain Technologies Already Upon Us?",
	"An Overview of MakerDAO",
	"The Universal Crypto Exchange APIs",
	"Vijay Boyapati’s Bullish Case for Bitcoin",
	"Coin Center’s Peter Van Valkenburg on Preserving the Freedom to Innovate with Public Blockchains",
	"Jesse Powell is Building a Culture of Crypto Values at Kraken",
	"Adam Back on a Decade of Bitcoin",
	" Will Bitcoin enchain the world",
	"Why are CBD & Kratom Vendors are Switching to Cryptocurrency?",
	"Constructing Cryptocurrency Indices — Performance & Methodology",
	"How to make money on arbitrage with cryptocurrencies",
	"Questioning the Obsession with Blockchains and On-Chain Governance with Nic Carter",
	"Francis Pouliot on the Network Effect of Money and Why Tokens Are Scams",
	"Brave’s Brendan Eich on Fixing Online Advertising",
	"Will Bitcoin enchain the world?",
	"OTC crypto deals, part 2: Minimize your risks",
	"What Bitcoin, Ethereum and other digital assets will become",
	"Will We Ever Run Out of Bitcoin Wallets?",
	"The Weaknesses of Blockchain and Decentralization.",
	"A crypto-trader’s diary — week 13; TRON",
}
