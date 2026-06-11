package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/72nd/escposimg"
)

// PrintModeComparison demonstrates the differences between all three printing modes:
// - raster: Modern GS v 0 command (efficient, single command)
// - graphics: Advanced GS ( L command (extended features, modern printers)
// - column: Legacy ESC * command (maximum compatibility, band-by-band)
//
// This example processes the same image with all three modes and compares:
// - Output size (command overhead)
// - Processing time
// - Command structure differences
func main() {
	// Set up logging for detailed output
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// Use a test image - replace with your own image path
	imagePath := "benchmark/benchmark.png"
	if len(os.Args) > 1 {
		imagePath = os.Args[1]
	}

	// Check if image exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Fatalf("Image file not found: %s", imagePath)
	}

	fmt.Printf("🖨️  ESC/POS Print Mode Comparison\n")
	fmt.Printf("════════════════════════════════════\n")
	fmt.Printf("Image: %s\n\n", imagePath)

	// Test configuration for consistent comparison
	baseConfig := &escposimg.Config{
		PaperWidthMM:  80,
		DPI:           203,
		DitheringAlgo: escposimg.DitheringFloydSteinberg,
		DebugOutput:   false,
		CutPaper:      true,
	}

	// Define the three print modes to test
	modes := []struct {
		name        string
		mode        escposimg.PrintMode
		description string
		useCase     string
	}{
		{
			name:        "Raster",
			mode:        escposimg.PrintModeRaster,
			description: "Modern GS v 0 command - single command block",
			useCase:     "Modern printers, network printing, efficiency",
		},
		{
			name:        "Graphics",
			mode:        escposimg.PrintModeGraphics,
			description: "Advanced GS ( L command - extended features",
			useCase:     "Modern printers with graphics support, advanced features",
		},
		{
			name:        "Column",
			mode:        escposimg.PrintModeColumn,
			description: "Legacy ESC * command - band-by-band processing",
			useCase:     "Legacy printers, maximum compatibility, troubleshooting",
		},
	}

	// Load and process image once to get dimensions
	img, err := escposimg.LoadImage(imagePath)
	if err != nil {
		log.Fatalf("Failed to load image: %v", err)
	}

	// Calculate target dimensions
	config := *baseConfig
	config.PrintMode = escposimg.PrintModeRaster // Just for dimension calculation
	targetWidth := config.CalculatePixelWidth()
	scaledImg, err := escposimg.ScaleImage(img, targetWidth)
	if err != nil {
		log.Fatalf("Failed to scale image: %v", err)
	}
	bounds := scaledImg.Bounds()

	fmt.Printf("📐 Image Dimensions:\n")
	fmt.Printf("   Original: %dx%d pixels\n", img.Bounds().Dx(), img.Bounds().Dy())
	fmt.Printf("   Scaled:   %dx%d pixels (for %dmm paper at %d DPI)\n\n",
		bounds.Dx(), bounds.Dy(), config.PaperWidthMM, config.DPI)

	// Test each print mode
	results := make([]struct {
		name     string
		size     int
		duration time.Duration
		error    error
	}, len(modes))

	for i, mode := range modes {
		fmt.Printf("🔄 Testing %s Mode...\n", mode.name)

		// Configure for this mode
		config := *baseConfig
		config.PrintMode = mode.mode

		// Create file output for this mode
		outputFile := fmt.Sprintf("/tmp/escposimg_comparison_%s.bin", mode.name)
		output, err := escposimg.NewFileOutput(outputFile)
		if err != nil {
			results[i].error = fmt.Errorf("failed to create output: %w", err)
			continue
		}

		// Measure processing time
		startTime := time.Now()

		// Process the image
		err = escposimg.ProcessImage(imagePath, &config, output)

		duration := time.Since(startTime)
		results[i].duration = duration

		if err != nil {
			results[i].error = fmt.Errorf("processing failed: %w", err)
			output.Close()
			continue
		}

		output.Close()

		// Get file size
		if stat, err := os.Stat(outputFile); err == nil {
			results[i].size = int(stat.Size())
		} else {
			results[i].error = fmt.Errorf("failed to get file size: %w", err)
			continue
		}

		results[i].name = mode.name

		fmt.Printf("   ✅ Complete: %d bytes in %v\n", results[i].size, duration)
	}

	// Print comparison table
	fmt.Printf("\n📊 Comparison Results:\n")
	fmt.Printf("══════════════════════════════════════════════════════════════════\n")
	fmt.Printf("| %-10s | %-12s | %-12s | %-25s |\n", "Mode", "Size (bytes)", "Time", "Best Use Case")
	fmt.Printf("══════════════════════════════════════════════════════════════════\n")

	for i, mode := range modes {
		result := results[i]
		if result.error != nil {
			fmt.Printf("| %-10s | %-12s | %-12s | %-25s |\n",
				mode.name, "ERROR", "N/A", "See error below")
		} else {
			fmt.Printf("| %-10s | %-12d | %-12v | %-25s |\n",
				mode.name, result.size, result.duration.Truncate(time.Millisecond), mode.useCase)
		}
	}
	fmt.Printf("══════════════════════════════════════════════════════════════════\n\n")

	// Print detailed mode descriptions
	fmt.Printf("📋 Mode Details:\n")
	for i, mode := range modes {
		fmt.Printf("• %s: %s\n", mode.name, mode.description)
		if results[i].error != nil {
			fmt.Printf("  ❌ Error: %v\n", results[i].error)
		}
	}

	// Print technical comparison
	fmt.Printf("\n🔧 Technical Comparison:\n")
	fmt.Printf("• Raster mode uses a single GS v 0 command with all image data\n")
	fmt.Printf("• Graphics mode uses GS ( L commands for definition + print operations\n")
	fmt.Printf("• Column mode uses multiple ESC * commands (one per 8-pixel band)\n")

	// Calculate overhead
	if results[0].size > 0 && results[1].size > 0 && results[2].size > 0 {
		rasterSize := results[0].size
		graphicsSize := results[1].size
		columnSize := results[2].size

		fmt.Printf("\n📈 Command Overhead Analysis:\n")
		fmt.Printf("• Graphics vs Raster: +%d bytes (+%.1f%%)\n",
			graphicsSize-rasterSize, float64(graphicsSize-rasterSize)*100.0/float64(rasterSize))
		fmt.Printf("• Column vs Raster: +%d bytes (+%.1f%%)\n",
			columnSize-rasterSize, float64(columnSize-rasterSize)*100.0/float64(rasterSize))
	}

	fmt.Printf("\n✨ Recommendation:\n")
	fmt.Printf("• Use 'raster' for most modern thermal printers (best efficiency)\n")
	fmt.Printf("• Use 'graphics' for advanced features or specific printer requirements\n")
	fmt.Printf("• Use 'column' for legacy printers or compatibility troubleshooting\n")
}
