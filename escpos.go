package escposimg

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log/slog"
)

// ESC/POS command constants
const (
	ESC = 0x1B // Escape character
	GS  = 0x1D // Group separator
	LF  = 0x0A // Line feed
	CR  = 0x0D // Carriage return
)

// GenerateESCPOS generates ESC/POS commands from a dithered image
func GenerateESCPOS(img image.Image, config *Config) ([]byte, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	slog.Debug("Generating ESC/POS commands", "width", width, "height", height)

	var buf bytes.Buffer

	// Step 1: Initialize printer (ESC @)
	buf.WriteByte(ESC)
	buf.WriteByte('@')
	slog.Debug("Added printer initialization command")

	// Step 2: Optional debug text
	if config.DebugText != "" {
		buf.WriteString(config.DebugText)
		buf.WriteByte(LF)
		slog.Debug("Added debug text", "text", config.DebugText)
	}

	// Step 3: Convert image to raster format and generate print commands
	rasterData, err := convertToRasterFormat(img)
	if err != nil {
		return nil, fmt.Errorf("failed to convert image to raster format: %w", err)
	}

	// Step 4: Generate raster image command (GS v 0)
	err = writeRasterImageCommand(&buf, width, height, rasterData)
	if err != nil {
		return nil, fmt.Errorf("failed to write raster image command: %w", err)
	}

	// Step 5: Feed paper and cut if requested
	buf.WriteByte(LF)
	buf.WriteByte(LF)
	buf.WriteByte(LF)

	if config.CutPaper {
		// Partial cut command (GS V 1)
		buf.WriteByte(GS)
		buf.WriteByte('V')
		buf.WriteByte(1)
		slog.Debug("Added paper cut command")
	}

	slog.Debug("ESC/POS command generation completed", "total_bytes", buf.Len())
	return buf.Bytes(), nil
}

// convertToRasterFormat converts a monochrome image to raster format for ESC/POS
func convertToRasterFormat(img image.Image) ([]byte, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Calculate bytes per line (width rounded up to nearest byte boundary)
	bytesPerLine := (width + 7) / 8

	rasterData := make([]byte, height*bytesPerLine)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Get pixel color
			pixel := img.At(x+bounds.Min.X, y+bounds.Min.Y)
			grayColor := color.GrayModel.Convert(pixel).(color.Gray)

			// Black pixels (Y=0) should print, white pixels (Y=255) should not
			if grayColor.Y < 128 {
				// Set bit for black pixel
				byteIndex := y*bytesPerLine + x/8
				bitIndex := uint(7 - (x % 8))
				rasterData[byteIndex] |= 1 << bitIndex
			}
		}
	}

	return rasterData, nil
}

// writeRasterImageCommand writes the GS v 0 command for raster image printing
func writeRasterImageCommand(buf *bytes.Buffer, width, height int, rasterData []byte) error {
	// Calculate bytes per line
	bytesPerLine := (width + 7) / 8

	// GS v 0 m xL xH yL yH [data]
	buf.WriteByte(GS)  // GS
	buf.WriteByte('v') // v
	buf.WriteByte('0') // 0
	buf.WriteByte(0)   // m (normal mode)

	// Width in bytes (xL + xH * 256)
	buf.WriteByte(byte(bytesPerLine & 0xFF))        // xL
	buf.WriteByte(byte((bytesPerLine >> 8) & 0xFF)) // xH

	// Height in dots (yL + yH * 256)
	buf.WriteByte(byte(height & 0xFF))        // yL
	buf.WriteByte(byte((height >> 8) & 0xFF)) // yH

	// Write raster data
	buf.Write(rasterData)

	slog.Debug("Wrote raster image command",
		"width_bytes", bytesPerLine,
		"height", height,
		"data_size", len(rasterData))

	return nil
}

// GenerateTestPattern generates a simple test pattern for debugging
func GenerateTestPattern(width, height int) []byte {
	var buf bytes.Buffer

	// Initialize printer
	buf.WriteByte(ESC)
	buf.WriteByte('@')

	// Add test text
	buf.WriteString("ESC/POS Test Pattern")
	buf.WriteByte(LF)
	buf.WriteByte(LF)

	// Generate simple pattern data
	bytesPerLine := (width + 7) / 8
	rasterData := make([]byte, height*bytesPerLine)

	// Create checkerboard pattern
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if (x/8+y/8)%2 == 0 {
				byteIndex := y*bytesPerLine + x/8
				bitIndex := uint(7 - (x % 8))
				rasterData[byteIndex] |= 1 << bitIndex
			}
		}
	}

	// Write raster command
	buf.WriteByte(GS)  // GS
	buf.WriteByte('v') // v
	buf.WriteByte('0') // 0
	buf.WriteByte(0)   // m
	buf.WriteByte(byte(bytesPerLine & 0xFF))
	buf.WriteByte(byte((bytesPerLine >> 8) & 0xFF))
	buf.WriteByte(byte(height & 0xFF))
	buf.WriteByte(byte((height >> 8) & 0xFF))
	buf.Write(rasterData)

	// Feed and cut
	buf.WriteByte(LF)
	buf.WriteByte(LF)
	buf.WriteByte(GS)
	buf.WriteByte('V')
	buf.WriteByte(1)

	return buf.Bytes()
}
