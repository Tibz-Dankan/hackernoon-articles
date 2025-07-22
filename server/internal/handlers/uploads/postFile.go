package uploads

import (
	"context"
	"log"
	"time"

	"github.com/Tibz-Dankan/hackernoon-articles/internal/pkg"
	"github.com/gofiber/fiber/v2"
)

func UploadFiles(c *fiber.Ctx) error {
	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid multipart form",
			"message": err.Error(),
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "No files provided",
			"message": "Please upload at least one file with field name 'files'",
		})
	}

	var uploadResponses []pkg.UploadResponse

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	for _, file := range files {
		// Open the file
		src, err := file.Open()
		if err != nil {
			log.Printf("Error opening file %s: %v", file.Filename, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to open file",
				"message": err.Error(),
			})
		}

		// Upload to S3
		s3Client := pkg.S3Client{}

		newS3Client, err := s3Client.NewS3Client(ctx)
		if err != nil {
			log.Printf("Error creating newS3Client: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to upload file to S3",
				"message": err.Error(),
			})
		}

		uploadResp, err := newS3Client.UploadFile(
			ctx,
			src,
			file.Filename,
			file.Header.Get("Content-Type"),
			file.Size,
		)
		src.Close()

		if err != nil {
			log.Printf("Error uploading file %s to S3: %v", file.Filename, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to upload file to S3",
				"message": err.Error(),
			})
		}

		uploadResponses = append(uploadResponses, *uploadResp)
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Files uploaded successfully",
		"files":   uploadResponses,
		"count":   len(uploadResponses),
	})
}
