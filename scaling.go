package escposimg

import (
	"image"
	"log/slog"

	"github.com/nfnt/resize"
)

// ScaleImage scales an image to the specified width while maintaining aspect ratio.
// Uses Lanczos3 interpolation for high quality scaling.
func ScaleImage(img image.Image, targetWidth int) (image.Image, error) {
	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	// If the image is already the target width, return as-is
	if originalWidth == targetWidth {
		slog.Debug("Image already at target width, no scaling needed", "width", targetWidth)
		return img, nil
	}

	slog.Debug("Scaling image",
		"original_width", originalWidth,
		"original_height", originalHeight,
		"target_width", targetWidth)

	// Use Lanczos3 for high-quality scaling
	// Height is set to 0 to preserve aspect ratio automatically
	scaledImg := resize.Resize(uint(targetWidth), 0, img, resize.Lanczos3)

	newBounds := scaledImg.Bounds()
	slog.Debug("Image scaled successfully",
		"new_width", newBounds.Dx(),
		"new_height", newBounds.Dy())

	return scaledImg, nil
}
