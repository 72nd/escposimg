// Package escposimg provides functionality to process images for ESC/POS thermal printers.
// It supports loading PNG, JPEG, and Netpbm images; applying dithering (skipped for PBM);
// scaling to paper width; and generating ESC/POS commands for printing.
package escposimg

import (
	"fmt"
	"image"
	"log/slog"

	"github.com/nfnt/resize"
)

// ProcessImage is the main function that processes an image and sends it to the specified output.
// It performs the complete pipeline: load → scale → dither → generate ESC/POS → output.
// PBM inputs skip dithering and use nearest-neighbor scaling to preserve bitmap edges.
func ProcessImage(imagePath string, config *Config, output OutputMethod) error {
	slog.Debug("Starting image processing", "path", imagePath, "config", config)

	// Step 1: Load the image
	loadRes, err := LoadImageDetails(imagePath)
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}
	img := loadRes.Image
	slog.Debug("Image loaded successfully", "width", img.Bounds().Dx(), "height", img.Bounds().Dy(),
		"format", loadRes.Format, "skip_dither", loadRes.SkipDithering)

	warnWeakBinarization(config, loadRes)

	// Step 2: Calculate target pixel width based on paper width and DPI
	targetWidth := config.CalculatePixelWidth()
	slog.Debug("Target width calculated", "width_pixels", targetWidth, "paper_mm", config.PaperWidthMM, "dpi", config.DPI)

	// Step 3: Scale the image to fit the paper width
	interp := resize.Lanczos3
	if loadRes.SkipDithering {
		interp = resize.NearestNeighbor
	}
	scaledImg, err := ScaleImageInterp(img, targetWidth, interp)
	if err != nil {
		return fmt.Errorf("failed to scale image: %w", err)
	}
	slog.Debug("Image scaled successfully", "new_width", scaledImg.Bounds().Dx(), "new_height", scaledImg.Bounds().Dy())

	var printImg image.Image
	if loadRes.SkipDithering {
		printImg = scaledImg
		slog.Debug("Skipping dithering for PBM source", "algorithm", config.DitheringAlgo.String())
	} else {
		printImg, err = ApplyDithering(scaledImg, config.DitheringAlgo)
		if err != nil {
			return fmt.Errorf("failed to apply dithering: %w", err)
		}
		slog.Debug("Dithering applied successfully", "algorithm", config.DitheringAlgo.String())
	}

	// Step 5: Save debug image if requested
	if config.DebugOutput {
		if err := SaveDebugImage(printImg, config.DebugImagePath); err != nil {
			slog.Warn("Failed to save debug image", "error", err)
		} else {
			slog.Debug("Debug image saved", "path", config.DebugImagePath)
		}
	}

	// Step 6: Generate ESC/POS commands
	escposData, err := GenerateESCPOS(printImg, config)
	if err != nil {
		return fmt.Errorf("failed to generate ESC/POS commands: %w", err)
	}
	slog.Debug("ESC/POS commands generated", "data_size", len(escposData))

	// Step 7: Send to output
	if err := output.Write(escposData); err != nil {
		return fmt.Errorf("failed to write to output: %w", err)
	}
	slog.Debug("Data sent to output successfully")

	// Step 8: Close output
	if err := output.Close(); err != nil {
		return fmt.Errorf("failed to close output: %w", err)
	}

	slog.Info("Image processing completed successfully")
	return nil
}

func warnWeakBinarization(config *Config, loadRes LoadImageResult) {
	if loadRes.SkipDithering {
		return
	}
	if config.DitheringAlgo != DitheringNone && config.DitheringAlgo != DitheringThreshold {
		return
	}
	slog.Warn("fixed threshold binarization is often weak for continuous-tone images; consider error-diffusion or ordered dither unless the asset is already high-contrast",
		"format", loadRes.Format,
		"dithering", config.DitheringAlgo.String())
}

// Version returns the current version of the escposimg library
func Version() string {
	return "0.1.0"
}
