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
