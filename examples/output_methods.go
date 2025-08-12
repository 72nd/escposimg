// Package main demonstrates different output methods available in escposimg.
// This example shows how to send ESC/POS commands to stdout, a file, or a network printer.
package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/72nd/escposimg"
)

func main() {
	// Set up logging
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// Use the test pattern image
	imagePath := "test_pattern.png"

	// Create a basic configuration
	config := &escposimg.Config{
		PaperWidthMM:  80,
		DPI:           203,
		DitheringAlgo: escposimg.DitheringFloydSteinberg,
		DebugText:     "Output Methods Demo",
		CutPaper:      true,
	}

	fmt.Println("ESC/POS Output Methods Demo")
	fmt.Println("==========================")

	// Method 1: Output to stdout
	fmt.Println("\n1. Stdout Output:")
	fmt.Println("   Use this to pipe data to other programs or redirect to files")
	fmt.Println("   Example: go run output_methods.go > printer_data.escpos")

	stdoutOutput := escposimg.NewStdoutOutput()
	fmt.Println("   Processing image to stdout...")
	if err := escposimg.ProcessImage(imagePath, config, stdoutOutput); err != nil {
		log.Printf("Error with stdout output: %v", err)
	} else {
		fmt.Println("   ✓ ESC/POS commands sent to stdout")
	}

	// Method 2: Output to file
	fmt.Println("\n2. File Output:")
	fmt.Println("   Save ESC/POS commands to a file for later use")

	fileOutput, err := escposimg.NewFileOutput("printer_commands.escpos")
	if err != nil {
		log.Printf("Error creating file output: %v", err)
	} else {
		fmt.Println("   Processing image to file...")
		if err := escposimg.ProcessImage(imagePath, config, fileOutput); err != nil {
			log.Printf("Error with file output: %v", err)
		} else {
			fmt.Println("   ✓ ESC/POS commands saved to: printer_commands.escpos")
			fmt.Println("   You can send this file to a printer later:")
			fmt.Println("   cat printer_commands.escpos > /dev/usb/lp0  # Linux USB printer")
			fmt.Println("   nc 192.168.1.100 9100 < printer_commands.escpos  # Network printer")
		}
		fileOutput.Close()
	}

	// Method 3: Network output (example - this will likely fail unless you have a printer)
	fmt.Println("\n3. Network Output:")
	fmt.Println("   Send directly to a network printer (requires actual printer)")
	fmt.Println("   Example printer IP: 192.168.1.100:9100")

	// Note: This will likely fail in most demo environments
	printerIP := "192.168.1.100:9100"
	fmt.Printf("   Attempting to connect to %s...\n", printerIP)

	networkOutput, err := escposimg.NewNetworkOutput(printerIP)
	if err != nil {
		fmt.Printf("   ⚠ Network connection failed (expected): %v\n", err)
		fmt.Println("   This is normal if you don't have a printer at this address")
		fmt.Println("   To use network printing:")
		fmt.Println("   1. Find your printer's IP address")
		fmt.Println("   2. Ensure it accepts raw TCP connections on port 9100")
		fmt.Println("   3. Update the IP address in this example")
	} else {
		fmt.Println("   Processing image to network printer...")
		if err := escposimg.ProcessImage(imagePath, config, networkOutput); err != nil {
			log.Printf("Error with network output: %v", err)
		} else {
			fmt.Println("   ✓ ESC/POS commands sent to network printer")
		}
		networkOutput.Close()
	}

	fmt.Println("\nOutput Methods Summary:")
	fmt.Println("======================")
	fmt.Println("• Stdout: Best for piping and shell integration")
	fmt.Println("• File: Best for batch processing and debugging")
	fmt.Println("• Network: Best for direct printing to network printers")
	fmt.Println("\nTip: You can test any output by redirecting to a file and examining with a hex editor:")
	fmt.Println("go run output_methods.go 2>/dev/null | hexdump -C")
}
