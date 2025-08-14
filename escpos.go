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
// Supports both raster mode (GS v 0) and bit image mode (ESC *)
func GenerateESCPOS(img image.Image, config *Config) ([]byte, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	slog.Debug("Generating ESC/POS commands",
		"width", width,
		"height", height,
		"print_mode", config.PrintMode.String())

	// Dispatch to appropriate mode-specific function
	switch config.PrintMode {
	case PrintModeRaster:
		return generateRasterMode(img, config)
	case PrintModeBitImage:
		return generateBitImageMode(img, config)
	default:
		return nil, fmt.Errorf("unsupported print mode: %v", config.PrintMode)
	}
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

// convertToBitImageFormat converts a monochrome image to bit image format for ESC *.
//
// The ESC * command processes images in horizontal bands of 8 pixels height.
// Each column in a band is represented by a single byte, where each bit
// corresponds to a vertical pixel (bit 0 = top, bit 7 = bottom).
//
// This format is compatible with legacy thermal printers and provides
// line-by-line processing for better compatibility with older hardware.
//
// Parameters:
//   - img: Source image (should be monochrome/dithered)
//
// Returns:
//   - []byte: Formatted data ready for ESC * commands
//   - error: If image processing fails
func convertToBitImageFormat(img image.Image) ([]byte, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// ESC * mode 0: 8-dot single-density
	// Each band is 8 pixels high, each column takes 1 byte
	bands := (height + 7) / 8
	bytesPerBand := width
	bitImageData := make([]byte, bands*bytesPerBand)

	for band := 0; band < bands; band++ {
		for x := 0; x < width; x++ {
			var columnByte byte

			// Process 8 pixels vertically for this column
			for bit := 0; bit < 8; bit++ {
				y := band*8 + bit
				if y < height {
					// Get pixel color
					pixel := img.At(x+bounds.Min.X, y+bounds.Min.Y)
					grayColor := color.GrayModel.Convert(pixel).(color.Gray)

					// Black pixels (Y=0) should print
					if grayColor.Y < 128 {
						// Set bit (bit 0 = top pixel, bit 7 = bottom pixel)
						columnByte |= 1 << uint(bit)
					}
				}
			}

			// Store the column byte
			bitImageData[band*bytesPerBand+x] = columnByte
		}
	}

	return bitImageData, nil
}

// writeBitImageCommand writes ESC * commands for bit image printing.
//
// Generates a series of ESC * commands to print the image data band by band.
// Each band represents 8 pixels of height, and the entire image width is
// sent with each command. After each band, a line feed advances the paper.
//
// Command format for each band: ESC * m nL nH [data]
// Where:
//   - ESC * = Start of bit image command
//   - m = Mode (0 = 8-dot single-density)
//   - nL, nH = Width in dots (little-endian 16-bit)
//   - [data] = Column data for this band
//
// Parameters:
//   - buf: Buffer to write commands to
//   - width: Image width in pixels
//   - height: Image height in pixels
//   - bitImageData: Pre-formatted bit image data from convertToBitImageFormat
//
// Returns:
//   - error: If command generation fails
func writeBitImageCommand(buf *bytes.Buffer, width, height int, bitImageData []byte) error {
	bands := (height + 7) / 8
	bytesPerBand := width

	slog.Debug("Writing bit image command",
		"width", width,
		"height", height,
		"bands", bands,
		"bytes_per_band", bytesPerBand)

	for band := 0; band < bands; band++ {
		// ESC * m nL nH [data]
		buf.WriteByte(ESC) // ESC
		buf.WriteByte('*') // *
		buf.WriteByte(0)   // m (mode 0: 8-dot single-density)

		// Width in dots (nL + nH * 256)
		buf.WriteByte(byte(width & 0xFF))        // nL
		buf.WriteByte(byte((width >> 8) & 0xFF)) // nH

		// Write band data
		bandStart := band * bytesPerBand
		bandEnd := bandStart + bytesPerBand
		buf.Write(bitImageData[bandStart:bandEnd])

		// Line feed after each band
		buf.WriteByte(LF)

		slog.Debug("Wrote bit image band",
			"band", band,
			"data_size", bytesPerBand)
	}

	return nil
}

// generateRasterMode generates ESC/POS commands using GS v 0 (raster mode).
//
// This function implements the modern raster image printing approach using
// the GS v 0 command. The entire image is sent as a single command block,
// making it efficient for large images and network printing.
//
// Process:
//  1. Initialize printer (ESC @)
//  2. Add optional debug text
//  3. Convert image to raster format (horizontal bit packing)
//  4. Send single GS v 0 command with all image data
//  5. Add paper feeds and optional cut command
//
// Parameters:
//   - img: Source image (should be monochrome/dithered)
//   - config: Configuration including paper settings and options
//
// Returns:
//   - []byte: Complete ESC/POS command sequence
//   - error: If generation fails
func generateRasterMode(img image.Image, config *Config) ([]byte, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	slog.Debug("Generating raster mode commands", "width", width, "height", height)

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

	slog.Debug("Raster mode command generation completed", "total_bytes", buf.Len())
	return buf.Bytes(), nil
}

// generateBitImageMode generates ESC/POS commands using ESC * (bit image mode).
//
// This function implements the traditional bit image printing approach using
// ESC * commands. The image is processed in 8-pixel height bands, with each
// band sent as a separate command. This provides better compatibility with
// legacy thermal printers at the cost of increased command overhead.
//
// Process:
//  1. Initialize printer (ESC @)
//  2. Add optional debug text
//  3. Convert image to bit image format (vertical column packing)
//  4. Send series of ESC * commands, one per 8-pixel band
//  5. Add paper feeds and optional cut command
//
// Parameters:
//   - img: Source image (should be monochrome/dithered)
//   - config: Configuration including paper settings and options
//
// Returns:
//   - []byte: Complete ESC/POS command sequence
//   - error: If generation fails
func generateBitImageMode(img image.Image, config *Config) ([]byte, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	slog.Debug("Generating bit image mode commands", "width", width, "height", height)

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

	// Step 3: Convert image to bit image format and generate print commands
	bitImageData, err := convertToBitImageFormat(img)
	if err != nil {
		return nil, fmt.Errorf("failed to convert image to bit image format: %w", err)
	}

	// Step 4: Generate bit image commands (ESC *)
	err = writeBitImageCommand(&buf, width, height, bitImageData)
	if err != nil {
		return nil, fmt.Errorf("failed to write bit image command: %w", err)
	}

	// Step 5: Feed paper and cut if requested
	buf.WriteByte(LF)
	buf.WriteByte(LF)

	if config.CutPaper {
		// Partial cut command (GS V 1)
		buf.WriteByte(GS)
		buf.WriteByte('V')
		buf.WriteByte(1)
		slog.Debug("Added paper cut command")
	}

	slog.Debug("Bit image mode command generation completed", "total_bytes", buf.Len())
	return buf.Bytes(), nil
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
