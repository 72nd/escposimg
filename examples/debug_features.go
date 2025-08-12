// Package main demonstrates debug features in escposimg.
// This example shows how to use debug output, debug text, and logging
// to troubleshoot image processing and printer issues.
package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/irvin/escposimg"
)

func main() {
	fmt.Println("ESC/POS Debug Features Demo")
	fmt.Println("===========================")

	// Use the test pattern image
	imagePath := "test_pattern.png"

	// Example 1: Basic debug output
	fmt.Println("\n1. Basic Debug Output:")
	fmt.Println("   Saves the dithered image to see processing results")

	// Set up INFO level logging to see processing steps
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	basicDebugConfig := &escposimg.Config{
		PaperWidthMM:   80,
		DPI:            203,
		DitheringAlgo:  escposimg.DitheringFloydSteinberg,
		DebugOutput:    true,
		DebugImagePath: "debug_basic.png",
		DebugText:      "",
		CutPaper:       false,
	}

	fileOutput1, err := escposimg.NewFileOutput("debug_basic.escpos")
	if err != nil {
		log.Printf("Error creating basic debug output: %v", err)
	} else {
		if err := escposimg.ProcessImage(imagePath, basicDebugConfig, fileOutput1); err != nil {
			log.Printf("Error processing basic debug: %v", err)
		} else {
			fmt.Println("   ✓ Debug image saved to: debug_basic.png")
			fmt.Println("   ✓ ESC/POS commands saved to: debug_basic.escpos")
		}
		fileOutput1.Close()
	}

	// Example 2: Debug with text output on receipt
	fmt.Println("\n2. Debug Text on Receipt:")
	fmt.Println("   Adds text to the printed output for identification")

	debugTextConfig := &escposimg.Config{
		PaperWidthMM:   80,
		DPI:            203,
		DitheringAlgo:  escposimg.DitheringAtkinson,
		DebugOutput:    true,
		DebugImagePath: "debug_with_text.png",
		DebugText:      "DEBUG: Test Print - Atkinson Dithering",
		CutPaper:       true,
	}

	fileOutput2, err := escposimg.NewFileOutput("debug_with_text.escpos")
	if err != nil {
		log.Printf("Error creating debug text output: %v", err)
	} else {
		if err := escposimg.ProcessImage(imagePath, debugTextConfig, fileOutput2); err != nil {
			log.Printf("Error processing debug text: %v", err)
		} else {
			fmt.Println("   ✓ Debug image saved to: debug_with_text.png")
			fmt.Println("   ✓ ESC/POS with debug text saved to: debug_with_text.escpos")
			fmt.Println("   The printed receipt will include the debug text above the image")
		}
		fileOutput2.Close()
	}

	// Example 3: Verbose logging for troubleshooting
	fmt.Println("\n3. Verbose Debug Logging:")
	fmt.Println("   Enables detailed logging to see every processing step")

	// Switch to DEBUG level logging for detailed output
	debugLogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(debugLogger)

	verboseConfig := &escposimg.Config{
		PaperWidthMM:   58,  // Use smaller paper to see scaling effects
		DPI:            300, // Use higher DPI to see calculation details
		DitheringAlgo:  escposimg.DitheringBayer,
		DebugOutput:    true,
		DebugImagePath: "debug_verbose.png",
		DebugText:      "VERBOSE DEBUG: 58mm, 300 DPI, Bayer",
		CutPaper:       true,
	}

	fmt.Println("   Watch the detailed log output below:")
	fmt.Println("   " + strings.Repeat("─", 50))

	fileOutput3, err := escposimg.NewFileOutput("debug_verbose.escpos")
	if err != nil {
		log.Printf("Error creating verbose debug output: %v", err)
	} else {
		if err := escposimg.ProcessImage(imagePath, verboseConfig, fileOutput3); err != nil {
			log.Printf("Error processing verbose debug: %v", err)
		} else {
			fmt.Println("   " + strings.Repeat("─", 50))
			fmt.Println("   ✓ Verbose debug completed")
			fmt.Println("   ✓ Debug image saved to: debug_verbose.png")
			fmt.Println("   ✓ ESC/POS commands saved to: debug_verbose.escpos")
		}
		fileOutput3.Close()
	}

	// Example 4: Multiple debug images for comparison
	fmt.Println("\n4. Multiple Debug Images for Algorithm Comparison:")
	fmt.Println("   Generate debug images with different settings for side-by-side comparison")

	// Reset to INFO level for cleaner output
	infoLogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(infoLogger)

	algorithms := []escposimg.DitheringType{
		escposimg.DitheringThreshold,
		escposimg.DitheringFloydSteinberg,
		escposimg.DitheringBayer,
	}

	for i, algo := range algorithms {
		config := &escposimg.Config{
			PaperWidthMM:   80,
			DPI:            203,
			DitheringAlgo:  algo,
			DebugOutput:    true,
			DebugImagePath: fmt.Sprintf("debug_comparison_%d_%s.png", i+1, algo.String()),
			DebugText:      fmt.Sprintf("Comparison %d: %s", i+1, algo.String()),
			CutPaper:       false,
		}

		// Use stdout for these (we just want the debug images)
		stdoutOutput := escposimg.NewStdoutOutput()
		if err := escposimg.ProcessImage(imagePath, config, stdoutOutput); err != nil {
			log.Printf("Error processing comparison %d: %v", i+1, err)
		} else {
			fmt.Printf("   ✓ Comparison image %d saved: debug_comparison_%d_%s.png\n", i+1, i+1, algo.String())
		}
	}

	fmt.Println("\nDebug Features Summary:")
	fmt.Println("======================")
	fmt.Println("Generated debug files:")
	fmt.Println("• debug_basic.png - Basic dithered output")
	fmt.Println("• debug_with_text.png - With debug text")
	fmt.Println("• debug_verbose.png - With verbose logging")
	fmt.Println("• debug_comparison_*.png - Algorithm comparison")

	fmt.Println("\nESC/POS command files:")
	fmt.Println("• debug_basic.escpos")
	fmt.Println("• debug_with_text.escpos")
	fmt.Println("• debug_verbose.escpos")

	fmt.Println("\nDebugging Tips:")
	fmt.Println("==============")
	fmt.Println("1. Always enable DebugOutput when developing/testing")
	fmt.Println("2. Use DebugText to identify prints when testing multiple configurations")
	fmt.Println("3. Enable verbose logging (slog.LevelDebug) to see detailed processing steps")
	fmt.Println("4. Compare debug images to choose the best dithering algorithm")
	fmt.Println("5. Check debug images if printed output doesn't match expectations")
	fmt.Println("6. Use different paper widths and DPIs in debug mode to test scaling")
	fmt.Println("\nTroubleshooting:")
	fmt.Println("• If image is too light: Try threshold or bayer dithering")
	fmt.Println("• If image is too dark: Try atkinson or sierra-lite dithering")
	fmt.Println("• If image is wrong size: Check paper width and DPI settings")
	fmt.Println("• If printer doesn't respond: Test with file output first")
}
