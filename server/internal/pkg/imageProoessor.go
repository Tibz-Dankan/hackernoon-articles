// package pkg

// import (
// 	"encoding/base64"
// 	"fmt"
// 	"io"
// 	"mime/multipart"
// 	"net/http"
// 	"time"
// )

// type ImageProcessor struct{}

// // FileToBytes converts  multipart.File Image to []byte
// func (ip *ImageProcessor) FileToBytes(file multipart.File) ([]byte, error) {
// 	bytes, err := io.ReadAll(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return bytes, nil
// }

// // ToBase64 converts []byte to Base64 string
// func (ip *ImageProcessor) ToBase64(data []byte) string {
// 	return base64.StdEncoding.EncodeToString(data)
// }

// // ConvertBase64ToBytes converts Base64 string to []byte
// func (ip *ImageProcessor) ConvertBase64ToBytes(base64String string) ([]byte, error) {
// 	imageBytes, err := base64.StdEncoding.DecodeString(base64String)
// 	if err != nil {
// 		return imageBytes, err
// 	}
// 	return imageBytes, nil
// }

// // FetchImageFromURL downloads an image from the given URL
// // and returns the image data as bytes
// func (ip *ImageProcessor) GetImageFromURL(url string) ([]byte, error) {
// 	client := &http.Client{
// 		Timeout: 30 * time.Second,
// 	}

// 	resp, err := client.Get(url)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch image from URL: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("failed to fetch image: HTTP %d", resp.StatusCode)
// 	}

// 	contentType := resp.Header.Get("Content-Type")
// 	if contentType == "" || !isImageContentType(contentType) {
// 		return nil, fmt.Errorf("URL does not point to an image, content-type: %s", contentType)
// 	}

// 	imgData, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read image data: %w", err)
// 	}

// 	return imgData, nil
// }

// func isImageContentType(contentType string) bool {
// 	imageTypes := []string{
// 		"image/jpeg",
// 		"image/jpg",
// 		"image/png",
// 		"image/gif",
// 		"image/webp",
// 		"image/bmp",
// 		"image/tiff",
// 		"image/svg+xml",
// 	}

// 	for _, imgType := range imageTypes {
// 		// if contentType == imgType || contentType[:len(imgType)] == imgType {
// 		if contentType == imgType {
// 			return true
// 		}
// 	}
// 	return false
// }

package pkg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type ImageProcessor struct{}

// FileToBytes converts multipart.File Image to []byte
func (ip *ImageProcessor) FileToBytes(file multipart.File) ([]byte, error) {
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// ToBase64 converts []byte to Base64 string
func (ip *ImageProcessor) ToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// ConvertBase64ToBytes converts Base64 string to []byte
func (ip *ImageProcessor) ConvertBase64ToBytes(base64String string) ([]byte, error) {
	imageBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return imageBytes, err
	}
	return imageBytes, nil
}

// FetchImageFromURL downloads an image from the given URL
// and returns the image data as bytes
func (ip *ImageProcessor) GetImageFromURL(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch image: HTTP %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" || !isImageContentType(contentType) {
		return nil, fmt.Errorf("URL does not point to an image, content-type: %s", contentType)
	}

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	return imgData, nil
}

// GetContentTypeFromBinary detects the content type of an image from its binary data
// by examining the file signature (magic bytes)
func (ip *ImageProcessor) GetContentTypeFromBinary(data []byte) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("empty data provided")
	}

	// Check for common image file signatures
	switch {
	case len(data) >= 8 && bytes.Equal(data[:8], []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}):
		return "image/png", nil
	case len(data) >= 3 && bytes.Equal(data[:3], []byte{0xFF, 0xD8, 0xFF}):
		return "image/jpeg", nil
	case len(data) >= 6 && (bytes.Equal(data[:6], []byte("GIF87a")) || bytes.Equal(data[:6], []byte("GIF89a"))):
		return "image/gif", nil
	case len(data) >= 12 && bytes.Equal(data[8:12], []byte("WEBP")):
		return "image/webp", nil
	case len(data) >= 2 && bytes.Equal(data[:2], []byte("BM")):
		return "image/bmp", nil
	case len(data) >= 4 && (bytes.Equal(data[:4], []byte("II*\x00")) || bytes.Equal(data[:4], []byte("MM\x00*"))):
		return "image/tiff", nil
	case len(data) >= 5 && bytes.Equal(data[:5], []byte("<?xml")):
		// Basic check for SVG - could be enhanced to look for <svg tag
		return "image/svg+xml", nil
	case len(data) >= 4 && bytes.Equal(data[:4], []byte("<svg")):
		return "image/svg+xml", nil
	default:
		return "", fmt.Errorf("unknown or unsupported image format")
	}
}

// BinaryToReader converts binary data ([]byte) to io.Reader
func (ip *ImageProcessor) BinaryToReader(data []byte) io.Reader {
	return bytes.NewReader(data)
}

func isImageContentType(contentType string) bool {
	imageTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/bmp",
		"image/tiff",
		"image/svg+xml",
	}

	for _, imgType := range imageTypes {
		if contentType == imgType {
			return true
		}
	}
	return false
}
