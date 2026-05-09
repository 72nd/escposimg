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
// Supports raster mode (GS v 0), graphics mode (GS ( L), and column mode (ESC *)
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
	case PrintModeGraphics:
		return generateGraphicsMode(img, config)
	case PrintModeColumn:
		return generateColumnMode(img, config)
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
// Uses ESC * mode 33 (24-dot double-density), which is the standard for 203 DPI
// thermal printers. Each band is 24 pixels tall; each column is represented by
// 3 bytes. Within each byte, bit 7 (MSB) is the topmost dot.
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

	// ESC * mode 33: 24-dot double-density
	// Each band is 24 pixels high, each column takes 3 bytes (one per 8-dot sub-row).
	bands := (height + 23) / 24
	bytesPerBand := width * 3
	bitImageData := make([]byte, bands*bytesPerBand)

	for band := 0; band < bands; band++ {
		for x := 0; x < width; x++ {
			// Process 24 pixels vertically for this column (3 bytes)
			for dot := 0; dot < 24; dot++ {
				y := band*24 + dot
				if y >= height {
					break
				}
				pixel := img.At(x+bounds.Min.X, y+bounds.Min.Y)
				grayColor := color.GrayModel.Convert(pixel).(color.Gray)
				if grayColor.Y < 128 {
					// Byte 0 covers dots 0-7, byte 1 covers 8-15, byte 2 covers 16-23.
					// Within each byte, bit 7 (MSB) = topmost dot per ESC/POS spec.
					byteOffset := dot / 8
					bitIndex := uint(7 - (dot % 8))
					bitImageData[band*bytesPerBand+x*3+byteOffset] |= 1 << bitIndex
				}
			}
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
	// Mode 33: 24-dot double-density. Each band is 24 rows; each column = 3 bytes.
	bands := (height + 23) / 24
	bytesPerBand := width * 3

	slog.Debug("Writing bit image command",
		"width", width,
		"height", height,
		"bands", bands,
		"bytes_per_band", bytesPerBand)

	// Set line spacing to 24 dots so each LF advances exactly one band height.
	// Default line spacing (~30 dots) would cause excess paper between bands.
	buf.WriteByte(ESC)
	buf.WriteByte('3')
	buf.WriteByte(24)

	for band := 0; band < bands; band++ {
		// ESC * m nL nH [data]
		buf.WriteByte(ESC)  // ESC
		buf.WriteByte('*')  // *
		buf.WriteByte(0x21) // m=33: 24-dot double-density

		// Width in dots (nL + nH * 256)
		buf.WriteByte(byte(width & 0xFF))        // nL
		buf.WriteByte(byte((width >> 8) & 0xFF)) // nH

		// Write band data (3 bytes per column)
		bandStart := band * bytesPerBand
		bandEnd := bandStart + bytesPerBand
		buf.Write(bitImageData[bandStart:bandEnd])

		// Line feed advances paper by the 24-dot line spacing set above
		buf.WriteByte(LF)

		slog.Debug("Wrote bit image band",
			"band", band,
			"data_size", bytesPerBand)
	}

	// Restore default line spacing
	buf.WriteByte(ESC)
	buf.WriteByte('2')

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

// generateColumnMode generates ESC/POS commands using ESC * (column mode).
//
// This function implements the traditional column-based printing approach using
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
func generateColumnMode(img image.Image, config *Config) ([]byte, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	slog.Debug("Generating column mode commands", "width", width, "height", height)

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

	slog.Debug("Column mode command generation completed", "total_bytes", buf.Len())
	return buf.Bytes(), nil
}

// writeGraphicsCommand writes the GS 8 L store command followed by a GS ( L print command.
//
// GS 8 L uses a 32-bit length field (p1..p4), which avoids the 16-bit overflow that
// would occur with GS ( L for images larger than ~900 rows at full paper width.
//
// Store command: GS 8 L p1 p2 p3 p4  m  fn  a bx by  c xL xH yL yH [data]
//                          (32-bit)  30  70 30  1  1 31
//   - fn=0x70 (112): store graphics data in the print buffer (raster format)
//   - a=0x30 (48): tone (mono)
//   - bx=by=0x01: 1x scale
//   - c=0x31 (49): black
//
// Print command: GS ( L 02 00 30 32
//   - fn=0x32 (50): print graphics data from the print buffer
func writeGraphicsCommand(buf *bytes.Buffer, width, height int, rasterData []byte) error {
	// Fixed header: m(1) fn(1) a(1) bx(1) by(1) c(1) xL(1) xH(1) yL(1) yH(1) = 10 bytes
	paramLength := uint32(10 + len(rasterData))

	slog.Debug("Writing graphics command (GS 8 L)",
		"width", width,
		"height", height,
		"data_size", len(rasterData),
		"param_length", paramLength)

	// GS 8 L — store graphics data in the print buffer (32-bit length)
	buf.WriteByte(GS)
	buf.WriteByte('8')
	buf.WriteByte('L')
	buf.WriteByte(byte(paramLength & 0xFF))         // p1
	buf.WriteByte(byte((paramLength >> 8) & 0xFF))  // p2
	buf.WriteByte(byte((paramLength >> 16) & 0xFF)) // p3
	buf.WriteByte(byte((paramLength >> 24) & 0xFF)) // p4
	buf.WriteByte(0x30)                             // m=48
	buf.WriteByte(0x70)                             // fn=112: store in print buffer (raster)
	buf.WriteByte(0x30)                             // a=48: tone (mono)
	buf.WriteByte(0x01)                             // bx=1: horizontal scale 1x
	buf.WriteByte(0x01)                             // by=1: vertical scale 1x
	buf.WriteByte(0x31)                             // c=49: black
	buf.WriteByte(byte(width & 0xFF))               // xL
	buf.WriteByte(byte((width >> 8) & 0xFF))        // xH
	buf.WriteByte(byte(height & 0xFF))              // yL
	buf.WriteByte(byte((height >> 8) & 0xFF))       // yH
	buf.Write(rasterData)

	// GS ( L — print graphics data from the print buffer
	buf.WriteByte(GS)
	buf.WriteByte('(')
	buf.WriteByte('L')
	buf.WriteByte(0x02) // pL: 2 parameter bytes (m + fn)
	buf.WriteByte(0x00) // pH
	buf.WriteByte(0x30) // m=48
	buf.WriteByte(0x32) // fn=50: print graphics from print buffer

	return nil
}

// generateGraphicsMode generates ESC/POS commands using GS 8 L / GS ( L (graphics mode).
//
// Uses GS 8 L (fn=112) to store the full image in the print buffer in one command —
// its 32-bit length field handles images of any size without overflow. Then GS ( L
// (fn=50) triggers printing.
//
// Process:
//  1. Initialize printer (ESC @)
//  2. Add optional debug text
//  3. Convert image to raster format
//  4. GS 8 L: store graphics data in the print buffer
//  5. GS ( L fn=50: print from the print buffer
//  6. Add paper feeds and optional cut command
//
// Parameters:
//   - img: Source image (should be monochrome/dithered)
//   - config: Configuration including paper settings and options
//
// Returns:
//   - []byte: Complete ESC/POS command sequence
//   - error: If generation fails
func generateGraphicsMode(img image.Image, config *Config) ([]byte, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	slog.Debug("Generating graphics mode commands", "width", width, "height", height)

	var buf bytes.Buffer

	// Step 1: Initialize printer (ESC @)
	buf.WriteByte(ESC)
	buf.WriteByte('@')

	// Step 2: Optional debug text
	if config.DebugText != "" {
		buf.WriteString(config.DebugText)
		buf.WriteByte(LF)
		slog.Debug("Added debug text", "text", config.DebugText)
	}

	// Step 3: Convert image to raster format
	rasterData, err := convertToRasterFormat(img)
	if err != nil {
		return nil, fmt.Errorf("failed to convert image to raster format: %w", err)
	}

	// Steps 4+5: Store and print graphics
	if err := writeGraphicsCommand(&buf, width, height, rasterData); err != nil {
		return nil, fmt.Errorf("failed to write graphics command: %w", err)
	}

	// Step 6: Feed paper and cut if requested
	buf.WriteByte(LF)
	buf.WriteByte(LF)
	buf.WriteByte(LF)

	if config.CutPaper {
		buf.WriteByte(GS)
		buf.WriteByte('V')
		buf.WriteByte(1)
		slog.Debug("Added paper cut command")
	}

	slog.Debug("Graphics mode command generation completed", "total_bytes", buf.Len())
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
