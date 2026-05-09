package escposimg

import (
	"image"
	"log/slog"

	"github.com/nfnt/resize"
)

// ScaleImage scales an image to the specified width while maintaining aspect ratio.
// Uses Lanczos3 interpolation for high quality scaling.
func ScaleImage(img image.Image, targetWidth int) (image.Image, error) {
	return ScaleImageInterp(img, targetWidth, resize.Lanczos3)
}

// ScaleImageInterp scales an image to the target width while maintaining aspect ratio,
// using the given nfnt/resize interpolation (e.g. Lanczos3 or NearestNeighbor).
func ScaleImageInterp(img image.Image, targetWidth int, interp resize.InterpolationFunction) (image.Image, error) {
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

	scaledImg := resize.Resize(uint(targetWidth), 0, img, interp)

	newBounds := scaledImg.Bounds()
	slog.Debug("Image scaled successfully",
		"new_width", newBounds.Dx(),
		"new_height", newBounds.Dy())

	return scaledImg, nil
}
