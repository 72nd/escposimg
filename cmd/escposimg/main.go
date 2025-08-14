package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/72nd/escposimg"
)

func main() {
	// Define command line flags
	var (
		imagePath      = flag.String("image", "", "Path to the image file (required)")
		paperWidth     = flag.Int("paper-width", 80, "Paper width in millimeters")
		dpi            = flag.Int("dpi", 203, "Printer DPI")
		ditheringAlgo  = flag.String("dithering", "floyd-steinberg", "Dithering algorithm (floyd-steinberg, atkinson, threshold, bayer, burkes, sierra-lite, jarvis-judice-ninke, shadura)")
		printMode      = flag.String("print-mode", "raster", "ESC/POS print mode (raster, bit-image)")
		debugOutput    = flag.Bool("debug-output", false, "Save dithered image for debugging")
		debugImagePath = flag.String("debug-image", "debug_output.png", "Path to save debug image")
		debugText      = flag.String("debug-text", "", "Optional debug text to print before image")
		cutPaper       = flag.Bool("cut", false, "Send paper cut command after printing")
		outputMethod   = flag.String("output", "stdout", "Output method (stdout, network, file)")
		networkAddr    = flag.String("network-addr", "", "Network address for network output (e.g., 192.168.1.100:9100)")
		filePath       = flag.String("file-path", "", "File path for file output")
		verbose        = flag.Bool("verbose", false, "Enable verbose logging")
		version        = flag.Bool("version", false, "Show version information")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "escposimg processes images for ESC/POS thermal printers.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -image photo.jpg\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -image photo.jpg -output network -network-addr 192.168.1.100:9100\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -image photo.jpg -dithering threshold -debug-output\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -image photo.jpg -print-mode bit-image\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -image photo.jpg -print-mode raster -dithering atkinson\n", os.Args[0])
	}

	flag.Parse()

	// Show version and exit
	if *version {
		fmt.Printf("escposimg version %s\n", escposimg.Version())
		return
	}

	// Set up logging
	logLevel := slog.LevelInfo
	if *verbose {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	// Validate required arguments
	if *imagePath == "" {
		fmt.Fprintf(os.Stderr, "Error: -image is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Parse dithering algorithm
	ditheringType, err := parseDitheringAlgo(*ditheringAlgo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Parse print mode
	printModeType, err := parsePrintMode(*printMode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Create configuration
	config := &escposimg.Config{
		PaperWidthMM:   *paperWidth,
		DPI:            *dpi,
		DitheringAlgo:  ditheringType,
		PrintMode:      printModeType,
		DebugOutput:    *debugOutput,
		DebugImagePath: *debugImagePath,
		DebugText:      *debugText,
		CutPaper:       *cutPaper,
	}

	// Create output method
	output, err := createOutputMethod(*outputMethod, *networkAddr, *filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output method: %v\n", err)
		os.Exit(1)
	}

	// Process the image
	if err := escposimg.ProcessImage(*imagePath, config, output); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing image: %v\n", err)
		os.Exit(1)
	}

	slog.Info("Image processed successfully")
}

// parseDitheringAlgo converts string to DitheringType
func parseDitheringAlgo(algo string) (escposimg.DitheringType, error) {
	switch strings.ToLower(algo) {
	case "floyd-steinberg":
		return escposimg.DitheringFloydSteinberg, nil
	case "atkinson":
		return escposimg.DitheringAtkinson, nil
	case "threshold":
		return escposimg.DitheringThreshold, nil
	case "bayer":
		return escposimg.DitheringBayer, nil
	case "burkes":
		return escposimg.DitheringBurkes, nil
	case "sierra-lite":
		return escposimg.DitheringSierraLite, nil
	case "jarvis-judice-ninke":
		return escposimg.DitheringJarvisJudiceNinke, nil
	case "shadura":
		return escposimg.DitheringShadura, nil
	default:
		return 0, fmt.Errorf("unknown dithering algorithm: %s", algo)
	}
}

// parsePrintMode converts string to PrintMode
func parsePrintMode(mode string) (escposimg.PrintMode, error) {
	switch strings.ToLower(mode) {
	case "raster":
		return escposimg.PrintModeRaster, nil
	case "bit-image":
		return escposimg.PrintModeBitImage, nil
	default:
		return 0, fmt.Errorf("unknown print mode: %s (supported: raster, bit-image)", mode)
	}
}

// createOutputMethod creates the appropriate output method based on the flag
func createOutputMethod(method, networkAddr, filePath string) (escposimg.OutputMethod, error) {
	switch strings.ToLower(method) {
	case "stdout":
		return escposimg.NewStdoutOutput(), nil
	case "network":
		if networkAddr == "" {
			return nil, fmt.Errorf("network address is required for network output")
		}
		return escposimg.NewNetworkOutput(networkAddr)
	case "file":
		if filePath == "" {
			return nil, fmt.Errorf("file path is required for file output")
		}
		return escposimg.NewFileOutput(filePath)
	default:
		return nil, fmt.Errorf("unknown output method: %s", method)
	}
}
