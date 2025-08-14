package escposimg

// DitheringType represents the available dithering algorithms
type DitheringType int

const (
	DitheringFloydSteinberg DitheringType = iota
	DitheringAtkinson
	DitheringThreshold
	DitheringBayer
	DitheringBurkes
	DitheringSierraLite
	DitheringJarvisJudiceNinke
	DitheringShadura
)

// PrintMode defines the ESC/POS printing mode for images.
//
// ESC/POS supports two main approaches for printing bitmap images:
// - Raster mode (GS v 0): Modern, efficient, single command per image
// - Bit Image mode (ESC *): Legacy compatible, line-by-line processing
//
// Both modes produce identical visual output when using the same dithering
// algorithm, but differ in command structure and printer compatibility.
type PrintMode int

const (
	// PrintModeRaster uses the GS v 0 command (Print Raster Bit Image).
	//
	// This is the modern standard for thermal printers and is recommended
	// for most use cases. The entire image is sent as a single command,
	// making it more efficient for large images and network printing.
	//
	// Command format: GS v 0 m xL xH yL yH [data]
	//
	// Best for:
	// - Modern thermal printers (post-2010)
	// - Network/Ethernet connections
	// - High-volume printing
	// - Large or high-resolution images
	//
	// Compatibility: Most modern ESC/POS printers support this mode.
	PrintModeRaster PrintMode = iota

	// PrintModeBitImage uses the ESC * command (Bit Image Mode).
	//
	// This is the traditional approach that processes images in 8-pixel
	// height bands, sending one command per band. While less efficient
	// than raster mode, it offers better compatibility with older printers.
	//
	// Command format: ESC * m nL nH [data] (repeated for each 8-pixel band)
	//
	// Best for:
	// - Legacy thermal printers (pre-2010)
	// - Serial/RS232 connections
	// - Compatibility troubleshooting
	// - Printers that don't support GS v 0
	//
	// Compatibility: Supported by virtually all ESC/POS printers,
	// including very old models.
	PrintModeBitImage
)

// String returns the string representation of the print mode.
// Returns "raster" for PrintModeRaster, "bit-image" for PrintModeBitImage,
// or "unknown" for invalid values.
func (p PrintMode) String() string {
	switch p {
	case PrintModeRaster:
		return "raster"
	case PrintModeBitImage:
		return "bit-image"
	default:
		return "unknown"
	}
}

// String returns the string representation of the dithering type
func (d DitheringType) String() string {
	switch d {
	case DitheringFloydSteinberg:
		return "floyd-steinberg"
	case DitheringAtkinson:
		return "atkinson"
	case DitheringThreshold:
		return "threshold"
	case DitheringBayer:
		return "bayer"
	case DitheringBurkes:
		return "burkes"
	case DitheringSierraLite:
		return "sierra-lite"
	case DitheringJarvisJudiceNinke:
		return "jarvis-judice-ninke"
	case DitheringShadura:
		return "shadura"
	default:
		return "unknown"
	}
}

// Config holds the configuration for image processing and printing
type Config struct {
	// Paper width in millimeters (default: 80mm)
	PaperWidthMM int

	// Printer DPI (default: 203 DPI)
	DPI int

	// Dithering algorithm to use
	DitheringAlgo DitheringType

	// ESC/POS printing mode for images (default: PrintModeRaster).
	//
	// Determines which ESC/POS command sequence to use for image printing:
	// - PrintModeRaster: Modern GS v 0 command, efficient, single command
	// - PrintModeBitImage: Legacy ESC * command, compatible, line-by-line
	//
	// Use PrintModeRaster for modern printers, PrintModeBitImage for legacy
	// compatibility or when experiencing printer communication issues.
	PrintMode PrintMode

	// Save dithered image for debugging
	DebugOutput bool

	// Path to save debug image (if DebugOutput is true)
	DebugImagePath string

	// Optional debug text to print before image
	DebugText string

	// Send paper cut command after printing
	CutPaper bool
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		PaperWidthMM:   80,
		DPI:            203,
		DitheringAlgo:  DitheringFloydSteinberg,
		PrintMode:      PrintModeRaster, // Default to modern raster mode
		DebugOutput:    false,
		DebugImagePath: "debug_output.png",
		DebugText:      "",
		CutPaper:       false,
	}
}

// OutputMethod defines the interface for different output methods
type OutputMethod interface {
	Write(data []byte) error
	Close() error
}

// Common DPI values for thermal printers
const (
	DPI203 = 203 // Standard thermal printer DPI
	DPI300 = 300 // High quality thermal printer DPI
	DPI180 = 180 // Older thermal printer models
)

// Common paper widths in millimeters
const (
	PaperWidth58mm = 58
	PaperWidth80mm = 80
)

// CalculatePixelWidth calculates the pixel width based on paper width and DPI
func (c *Config) CalculatePixelWidth() int {
	// Convert mm to inches, then multiply by DPI
	inches := float64(c.PaperWidthMM) / 25.4
	return int(inches * float64(c.DPI))
}
