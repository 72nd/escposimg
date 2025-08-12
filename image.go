package escposimg

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

// LoadImage loads an image from the specified file path.
// Supports PNG and JPEG formats.
func LoadImage(imagePath string) (image.Image, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	// Decode the image
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Log the detected format for debugging
	switch format {
	case "png", "jpeg":
		// Supported formats
	default:
		return nil, fmt.Errorf("unsupported image format: %s (supported: PNG, JPEG)", format)
	}

	return img, nil
}

// SaveDebugImage saves an image to the specified path for debugging purposes
func SaveDebugImage(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create debug image file: %w", err)
	}
	defer file.Close()

	// Save as PNG for debugging
	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode debug image: %w", err)
	}

	return nil
}

func init() {
	// Register image formats
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}
