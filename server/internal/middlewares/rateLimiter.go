package middlewares

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RateLimiter struct {
	requests     map[string]int
	blockedUntil map[string]time.Time
	resetAfter   map[string]time.Time
	limit        int
	window       time.Duration
	mutex        sync.Mutex
}

var rateLimiter = &RateLimiter{
	requests:     make(map[string]int),
	blockedUntil: make(map[string]time.Time),
	resetAfter:   make(map[string]time.Time),
	window:       60 * time.Second, // 1 min window
	limit:        20,
}

// AllowRequest blocks requests until the window expires
func (rl *RateLimiter) AllowRequest(clientIp string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	if unblockTime, blocked := rl.blockedUntil[clientIp]; blocked {
		if time.Now().Before(unblockTime) {
			return false
		}

		delete(rl.blockedUntil, clientIp)
	}

	if resetTime, exists := rl.resetAfter[clientIp]; exists && time.Now().After(resetTime) {
		delete(rl.requests, clientIp)
		delete(rl.resetAfter, clientIp)
	}

	count := rl.requests[clientIp]
	if count <= rl.limit {
		if count == 0 {
			rl.resetAfter[clientIp] = time.Now().Add(rl.window)
		}
		rl.requests[clientIp]++
		return true
	}

	// Block further requests
	rl.blockedUntil[clientIp] = time.Now().Add(rl.window)
	return false
}

// CleanupExpiredEntries removes expired IPs periodically
func (rl *RateLimiter) CleanupExpiredEntries() {
	log.Println("IP cleanup initiated...")
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()

		for ip, resetTime := range rl.resetAfter {
			if now.After(resetTime) {
				delete(rl.requests, ip)
				delete(rl.resetAfter, ip)
			}
		}

		for ip, unblockTime := range rl.blockedUntil {
			if now.After(unblockTime) {
				delete(rl.blockedUntil, ip)
			}
		}

		rl.mutex.Unlock()
	}
}

func RateLimit(c *fiber.Ctx) error {
	var clientIP string = ""

	if os.Getenv("GO_ENV") == "production" {
		clientIP = c.Get("CF-Connecting-IP") // Cloudflare's real IP header
		log.Println("IP Cloudflare header: ", clientIP)
	}

	if clientIP == "" {
		forwardedFor := c.Get("X-Forwarded-For")
		if forwardedFor != "" {
			clientIP = strings.Split(forwardedFor, ",")[0]
			log.Println("IP X-Forwarded-For: ", clientIP)
		}
	}

	if clientIP == "" {
		clientIP = c.Get("X-Real-IP")
		log.Println("IP X-Real-IP: ", clientIP)
	}

	if clientIP == "" {
		clientIP = c.IP()
		log.Println("IP default: ", clientIP)
	}

	log.Println("clientIP address:", clientIP)

	if !rateLimiter.AllowRequest(clientIP) {
		log.Println("Request denied for:", clientIP)
		return fiber.NewError(fiber.StatusTooManyRequests,
			"You have made too many requests! Please try again later.")
	}

	c.Locals("clientIP", clientIP)

	return c.Next()
}

func init() {
	go rateLimiter.CleanupExpiredEntries()
}
