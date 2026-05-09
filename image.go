package escposimg

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	// Register Netpbm decoders (pbm, pgm, ppm, pam) with image.Decode.
	_ "github.com/spakin/netpbm"
)

// LoadImageResult holds a decoded image and hints derived from the file format.
type LoadImageResult struct {
	// Image is the decoded bitmap.
	Image image.Image

	// Format is the string returned by image.Decode (e.g. "png", "jpeg", "pbm").
	Format string

	// SkipDithering is true when the source is portable bitmap (PBM); error diffusion is not applied.
	SkipDithering bool
}

// LoadImageDetails loads an image and reports format-specific processing hints.
// Supports PNG, JPEG, and Netpbm family (pbm, pgm, ppm, pam) via github.com/spakin/netpbm.
func LoadImageDetails(imagePath string) (LoadImageResult, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return LoadImageResult{}, fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return LoadImageResult{}, fmt.Errorf("failed to decode image: %w", err)
	}

	switch format {
	case "png", "jpeg", "pbm", "pgm", "ppm", "pam":
	default:
		return LoadImageResult{}, fmt.Errorf("unsupported image format: %s (supported: PNG, JPEG, Netpbm pbm/pgm/ppm/pam)", format)
	}

	skipDither := format == "pbm"
	return LoadImageResult{Image: img, Format: format, SkipDithering: skipDither}, nil
}

// LoadImage loads an image from the specified file path.
// Supports PNG, JPEG, and Netpbm (pbm, pgm, ppm, pam).
func LoadImage(imagePath string) (image.Image, error) {
	res, err := LoadImageDetails(imagePath)
	if err != nil {
		return nil, err
	}
	return res.Image, nil
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
