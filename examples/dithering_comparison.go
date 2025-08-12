// Package main demonstrates different dithering algorithms available in escposimg.
// This example processes the same image with all available dithering algorithms
// and saves the dithered results as debug images for comparison.
package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/irvin/escposimg"
)

func main() {
	// Set up logging to see the processing details
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// Use the test pattern image
	imagePath := "test_pattern.png"

	// All available dithering algorithms
	algorithms := []escposimg.DitheringType{
		escposimg.DitheringFloydSteinberg,
		escposimg.DitheringAtkinson,
		escposimg.DitheringThreshold,
		escposimg.DitheringBayer,
		escposimg.DitheringBurkes,
		escposimg.DitheringSierraLite,
		escposimg.DitheringJarvisJudiceNinke,
		escposimg.DitheringShadura,
	}

	// Create a file output to save the ESC/POS commands (we'll use the last one)
	output, err := escposimg.NewFileOutput("dithering_comparison.escpos")
	if err != nil {
		log.Fatalf("Failed to create file output: %v", err)
	}
	defer output.Close()

	fmt.Println("Comparing dithering algorithms:")
	fmt.Println("==============================")

	// Process the image with each dithering algorithm
	for _, algo := range algorithms {
		fmt.Printf("Processing with %s dithering...\n", algo.String())

		// Create configuration for this algorithm
		config := &escposimg.Config{
			PaperWidthMM:   80,
			DPI:            203,
			DitheringAlgo:  algo,
			DebugOutput:    true,
			DebugImagePath: fmt.Sprintf("dithered_%s.png", algo.String()),
			DebugText:      fmt.Sprintf("Dithering: %s", algo.String()),
			CutPaper:       false,
		}

		// For the comparison, we'll save debug images but only process ESC/POS for the last one
		if algo == escposimg.DitheringShadura {
			// Process the last one with actual ESC/POS output
			if err := escposimg.ProcessImage(imagePath, config, output); err != nil {
				log.Printf("Error processing with %s: %v", algo.String(), err)
				continue
			}
		} else {
			// For others, just create a dummy output to generate the debug images
			dummyOutput := escposimg.NewStdoutOutput()
			if err := escposimg.ProcessImage(imagePath, config, dummyOutput); err != nil {
				log.Printf("Error processing with %s: %v", algo.String(), err)
				continue
			}
		}

		fmt.Printf("  → Debug image saved: dithered_%s.png\n", algo.String())
	}

	fmt.Println("\nComparison complete!")
	fmt.Println("Check the generated PNG files to see the differences between algorithms:")
	for _, algo := range algorithms {
		fmt.Printf("  - dithered_%s.png\n", algo.String())
	}
	fmt.Println("ESC/POS commands saved to: dithering_comparison.escpos")
	fmt.Println("\nRecommendations:")
	fmt.Println("  - Floyd-Steinberg: Best overall quality for photos")
	fmt.Println("  - Atkinson: Good for images with fine details")
	fmt.Println("  - Threshold: Fast, good for high-contrast images")
	fmt.Println("  - Bayer: Good for textures and patterns")
}
