// package main

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/Tibz-Dankan/hackernoon-articles/internal/events/subscribers"
// 	"github.com/Tibz-Dankan/hackernoon-articles/internal/handlers/articles"
// 	"github.com/Tibz-Dankan/hackernoon-articles/internal/handlers/status"
// 	"github.com/Tibz-Dankan/hackernoon-articles/internal/handlers/uploads"
// 	"github.com/Tibz-Dankan/hackernoon-articles/internal/middlewares"
// 	"github.com/Tibz-Dankan/hackernoon-articles/internal/pkg"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/fiber/v2/middleware/cors"
// 	"github.com/gofiber/fiber/v2/middleware/logger"
// 	"github.com/joho/godotenv"
// )

// func main() {
// 	app := fiber.New(fiber.Config{
// 		ErrorHandler: pkg.DefaultErrorHandler,
// 	})

// 	app.Use(cors.New(cors.Config{
// 		AllowOrigins:  "*",
// 		AllowHeaders:  "Origin, Content-Type, Accept, Authorization",
// 		ExposeHeaders: "Content-Length",
// 	}))

// 	app.Use(logger.New())

// 	app.Use(middlewares.RateLimit)

// 	// Load dev .env file
// 	env := os.Getenv("GO_ENV")
// 	if env == "development" {
// 		err := godotenv.Load()
// 		if err != nil {
// 			log.Fatalf("Error loading .env file")
// 		}
// 		log.Println("Loaded .env var file")
// 	}

// 	// articles
// 	userGroup := app.Group("/api/v0.1/articles", func(c *fiber.Ctx) error {
// 		return c.Next()
// 	})
// 	userGroup.Get("/", articles.GetAllArticles)

// 	// uploads
// 	uploadGroup := app.Group("/api/v0.1/uploads", func(c *fiber.Ctx) error {
// 		return c.Next()
// 	})
// 	uploadGroup.Post("/", uploads.UploadFiles)

// 	// Status
// 	app.Get("/status", status.GetAppStatus)

// 	app.Use("*", func(c *fiber.Ctx) error {
// 		message := fmt.Sprintf("api route '%s' doesn't exist!", c.Path())
// 		return fiber.NewError(fiber.StatusNotFound, message)
// 	})

// 	// Initialize all event subscribers in the app
// 	subscribers.InitEventSubscribers()

// 	log.Fatal(app.Listen(":3000"))
// }

package main

import (
	"fmt"
	"log"
	_ "net/http/pprof" // Import pprof for HTTP endpoints
	"os"
	"runtime"

	"github.com/Tibz-Dankan/hackernoon-articles/internal/events/subscribers"
	"github.com/Tibz-Dankan/hackernoon-articles/internal/handlers/articles"
	"github.com/Tibz-Dankan/hackernoon-articles/internal/handlers/status"
	"github.com/Tibz-Dankan/hackernoon-articles/internal/handlers/uploads"
	"github.com/Tibz-Dankan/hackernoon-articles/internal/middlewares"
	"github.com/Tibz-Dankan/hackernoon-articles/internal/pkg"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/joho/godotenv"
)

func init() {
	// Enable all profiling types
	runtime.SetMutexProfileFraction(1) // Enable mutex profiling
	runtime.SetBlockProfileRate(1)     // Enable block profiling
}

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: pkg.DefaultErrorHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders: "Content-Length",
	}))

	app.Use(logger.New())

	app.Use(pprof.New())

	app.Use(middlewares.RateLimit)

	// Load dev .env file
	env := os.Getenv("GO_ENV")
	if env == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
		log.Println("Loaded .env var file")
	}

	// articles
	userGroup := app.Group("/api/v0.1/articles", func(c *fiber.Ctx) error {
		return c.Next()
	})
	userGroup.Get("/", articles.GetAllArticles)

	// uploads
	uploadGroup := app.Group("/api/v0.1/uploads", func(c *fiber.Ctx) error {
		return c.Next()
	})
	uploadGroup.Post("/", uploads.UploadFiles)

	// Status
	app.Get("/status", status.GetAppStatus)

	// Runtime metrics endpoint
	app.Get("/debug/runtime", func(c *fiber.Ctx) error {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		metrics := fiber.Map{
			"goroutines": runtime.NumGoroutine(),
			"gomaxprocs": runtime.GOMAXPROCS(0), // 0 means don't change, just return current value
			"memory": fiber.Map{
				"alloc":          memStats.Alloc,         // Current allocated memory
				"total_alloc":    memStats.TotalAlloc,    // Total allocated memory
				"sys":            memStats.Sys,           // System memory obtained from OS
				"num_gc":         memStats.NumGC,         // Number of GC runs
				"gc_cpu_percent": memStats.GCCPUFraction, // Fraction of CPU time used by GC
				"heap_alloc":     memStats.HeapAlloc,     // Heap allocated memory
				"heap_sys":       memStats.HeapSys,       // Heap system memory
				"heap_idle":      memStats.HeapIdle,      // Heap idle memory
				"heap_inuse":     memStats.HeapInuse,     // Heap in-use memory
				"heap_objects":   memStats.HeapObjects,   // Number of allocated heap objects
				"stack_inuse":    memStats.StackInuse,    // Stack in-use memory
				"stack_sys":      memStats.StackSys,      // Stack system memory
			},
			"timestamp": c.Context().Time(),
		}

		return c.JSON(metrics)
	})

	app.Use("*", func(c *fiber.Ctx) error {
		message := fmt.Sprintf("api route '%s' doesn't exist!", c.Path())
		return fiber.NewError(fiber.StatusNotFound, message)
	})

	// Initialize all event subscribers in the app
	subscribers.InitEventSubscribers()

	log.Println("pprof endpoints available at:")
	log.Println("  - CPU Profile: http://localhost:3000/debug/pprof/profile")
	log.Println("  - Memory Profile: http://localhost:3000/debug/pprof/heap")
	log.Println("  - Goroutines: http://localhost:3000/debug/pprof/goroutine")
	log.Println("  - Block Profile: http://localhost:3000/debug/pprof/block")
	log.Println("  - Mutex Profile: http://localhost:3000/debug/pprof/mutex")
	log.Println("  - All Profiles: http://localhost:3000/debug/pprof/")
	log.Println("  - Runtime Metrics: http://localhost:3000/debug/runtime")

	log.Fatal(app.Listen(":3000"))
}
