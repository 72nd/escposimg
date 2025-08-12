// Package main demonstrates various configuration options in escposimg.
// This example shows how to configure paper width, DPI, and other settings
// for different printer types and use cases.
package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/irvin/escposimg"
)

func main() {
	// Set up logging
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// Use the test pattern image
	imagePath := "test_pattern.png"

	fmt.Println("ESC/POS Configuration Options Demo")
	fmt.Println("==================================")

	// Example 1: 58mm paper width with standard DPI
	fmt.Println("\n1. Small Receipt Printer (58mm paper):")
	config58mm := &escposimg.Config{
		PaperWidthMM:   escposimg.PaperWidth58mm, // 58mm
		DPI:            escposimg.DPI203,         // 203 DPI
		DitheringAlgo:  escposimg.DitheringFloydSteinberg,
		DebugOutput:    true,
		DebugImagePath: "config_58mm.png",
		DebugText:      "58mm paper, 203 DPI",
		CutPaper:       true,
	}

	pixelWidth58 := config58mm.CalculatePixelWidth()
	fmt.Printf("   Paper width: %dmm → %d pixels\n", config58mm.PaperWidthMM, pixelWidth58)

	fileOutput58, err := escposimg.NewFileOutput("config_58mm.escpos")
	if err != nil {
		log.Printf("Error creating 58mm output: %v", err)
	} else {
		if err := escposimg.ProcessImage(imagePath, config58mm, fileOutput58); err != nil {
			log.Printf("Error processing 58mm config: %v", err)
		} else {
			fmt.Println("   ✓ Processed for 58mm printer → config_58mm.escpos")
		}
		fileOutput58.Close()
	}

	// Example 2: 80mm paper width with high DPI
	fmt.Println("\n2. Standard Receipt Printer (80mm paper, high DPI):")
	config80mmHiDPI := &escposimg.Config{
		PaperWidthMM:   escposimg.PaperWidth80mm, // 80mm
		DPI:            escposimg.DPI300,         // 300 DPI for higher quality
		DitheringAlgo:  escposimg.DitheringAtkinson,
		DebugOutput:    true,
		DebugImagePath: "config_80mm_300dpi.png",
		DebugText:      "80mm paper, 300 DPI, Atkinson dithering",
		CutPaper:       true,
	}

	pixelWidth80HiDPI := config80mmHiDPI.CalculatePixelWidth()
	fmt.Printf("   Paper width: %dmm → %d pixels (high resolution)\n", config80mmHiDPI.PaperWidthMM, pixelWidth80HiDPI)

	fileOutput80HiDPI, err := escposimg.NewFileOutput("config_80mm_300dpi.escpos")
	if err != nil {
		log.Printf("Error creating 80mm high DPI output: %v", err)
	} else {
		if err := escposimg.ProcessImage(imagePath, config80mmHiDPI, fileOutput80HiDPI); err != nil {
			log.Printf("Error processing 80mm high DPI config: %v", err)
		} else {
			fmt.Println("   ✓ Processed for high-DPI 80mm printer → config_80mm_300dpi.escpos")
		}
		fileOutput80HiDPI.Close()
	}

	// Example 3: Custom configuration for older printers
	fmt.Println("\n3. Older Thermal Printer (80mm paper, 180 DPI):")
	configOlder := &escposimg.Config{
		PaperWidthMM:   80,
		DPI:            escposimg.DPI180,             // Lower DPI for older printers
		DitheringAlgo:  escposimg.DitheringThreshold, // Simpler algorithm for speed
		DebugOutput:    true,
		DebugImagePath: "config_older_printer.png",
		DebugText:      "Older printer: 180 DPI, Threshold dithering",
		CutPaper:       false, // Some older printers don't support auto-cut
	}

	pixelWidthOlder := configOlder.CalculatePixelWidth()
	fmt.Printf("   Paper width: %dmm → %d pixels (older printer)\n", configOlder.PaperWidthMM, pixelWidthOlder)

	fileOutputOlder, err := escposimg.NewFileOutput("config_older_printer.escpos")
	if err != nil {
		log.Printf("Error creating older printer output: %v", err)
	} else {
		if err := escposimg.ProcessImage(imagePath, configOlder, fileOutputOlder); err != nil {
			log.Printf("Error processing older printer config: %v", err)
		} else {
			fmt.Println("   ✓ Processed for older printer → config_older_printer.escpos")
		}
		fileOutputOlder.Close()
	}

	// Example 4: Using DefaultConfig() and customizing
	fmt.Println("\n4. Using DefaultConfig() as a starting point:")
	defaultConfig := escposimg.DefaultConfig()
	fmt.Printf("   Default settings: %dmm paper, %d DPI, %s dithering\n",
		defaultConfig.PaperWidthMM, defaultConfig.DPI, defaultConfig.DitheringAlgo.String())

	// Customize the default config
	defaultConfig.DitheringAlgo = escposimg.DitheringBayer
	defaultConfig.DebugOutput = true
	defaultConfig.DebugImagePath = "config_customized_default.png"
	defaultConfig.DebugText = "Customized default config"
	defaultConfig.CutPaper = true

	fileOutputDefault, err := escposimg.NewFileOutput("config_customized_default.escpos")
	if err != nil {
		log.Printf("Error creating default config output: %v", err)
	} else {
		if err := escposimg.ProcessImage(imagePath, defaultConfig, fileOutputDefault); err != nil {
			log.Printf("Error processing default config: %v", err)
		} else {
			fmt.Println("   ✓ Processed with customized defaults → config_customized_default.escpos")
		}
		fileOutputDefault.Close()
	}

	fmt.Println("\nConfiguration Summary:")
	fmt.Println("=====================")
	fmt.Printf("58mm printer: %d pixels wide\n", pixelWidth58)
	fmt.Printf("80mm printer (203 DPI): %d pixels wide\n", escposimg.DefaultConfig().CalculatePixelWidth())
	fmt.Printf("80mm printer (300 DPI): %d pixels wide\n", pixelWidth80HiDPI)
	fmt.Printf("80mm printer (180 DPI): %d pixels wide\n", pixelWidthOlder)

	fmt.Println("\nGenerated files:")
	fmt.Println("• config_58mm.png & config_58mm.escpos")
	fmt.Println("• config_80mm_300dpi.png & config_80mm_300dpi.escpos")
	fmt.Println("• config_older_printer.png & config_older_printer.escpos")
	fmt.Println("• config_customized_default.png & config_customized_default.escpos")

	fmt.Println("\nTips:")
	fmt.Println("• Use 58mm for small receipt printers")
	fmt.Println("• Use 80mm for standard receipt printers")
	fmt.Println("• Use 203 DPI for most modern thermal printers")
	fmt.Println("• Use 300 DPI for high-quality photo printing")
	fmt.Println("• Use 180 DPI for older or slower printers")
	fmt.Println("• Choose dithering algorithm based on image content and speed requirements")
}
