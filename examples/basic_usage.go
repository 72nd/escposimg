// Package main demonstrates basic usage of the escposimg library.
// This example shows how to process an image with default settings
// and output the ESC/POS commands to stdout.
package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/irvin/escposimg"
)

func main() {
	// Set up basic logging to see what's happening
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// Use the test pattern image that comes with the library
	imagePath := "test_pattern.png"

	// Create a configuration with default settings:
	// - 80mm paper width
	// - 203 DPI
	// - Floyd-Steinberg dithering
	// - No debug output
	// - No paper cut
	config := escposimg.DefaultConfig()

	// Create stdout output to see the ESC/POS commands
	output := escposimg.NewStdoutOutput()

	// Process the image and send ESC/POS commands to stdout
	if err := escposimg.ProcessImage(imagePath, config, output); err != nil {
		log.Fatalf("Error processing image: %v", err)
	}

	slog.Info("Basic image processing completed successfully")
	slog.Info("ESC/POS commands have been sent to stdout")
	slog.Info("You can redirect this output to a printer or file: go run basic_usage.go > output.escpos")
}
