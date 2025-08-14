# escposimg

Go library and CLI tool for printing images with ESC/POS-compatible receipt printers.


## Background / Motivation

Many receipt printers require control via Epson's ESC/POS protocol. Unlike most modern office printers, this approach doesn't involve sending complete document layouts (such as PDF or PostScript files), but rather individual control commands like "output text" or "cut paper". The actual rendering takes place on the device itself. Consequently, the usual document creation tools (Typst, Word, LaTeX, Affinity Publisher, etc.) cannot be employed for layout purposes.

A viable workaround is to rasterise the desired document into a monochrome (black and white) image and transmit it to the printer in this format. This is precisely what this library accomplishes. Since thermal printers can only produce black output, image files containing colours or greyscale tones must first be converted to monochrome format. To this end, the tool provides a selection of different dithering algorithms. The print data can be transmitted to the printer via network connection.

Whilst printer drivers exist for all common platforms that perform the same function, experience demonstrates that these are sometimes implemented differently (e.g., employing different dithering algorithms), resulting in the same document being printed differently across various platforms.

## Usage

The `escposimg` tool provides both a command-line interface and a Go library for processing images and generating ESC/POS printer commands. Choose the approach that best suits your workflow.

### Command-Line Interface

The CLI tool offers straightforward image processing with configurable options for various printer types and use cases.

#### Basic Usage

Process an image with default settings (80mm paper, 203 DPI, Floyd-Steinberg dithering):

```bash
escposimg -image photo.jpg
```
This command converts a photograph to ESC/POS format using standard settings and outputs the printer commands to the terminal.

#### Advanced Examples

**High-quality printing with Atkinson dithering:**
```bash
escposimg -image artwork.png -dithering atkinson -dpi 300 -debug-output
```
This example processes artwork at high resolution with Atkinson dithering, which preserves fine details, and saves a debug image to verify the result.

**Legacy printer compatibility:**
```bash
escposimg -image receipt.jpg -print-mode bit-image -dpi 180 -paper-width 58
```
This command configures the tool for older thermal printers using bit-image mode, lower resolution, and narrow 58mm paper format.

#### Output Methods

The CLI supports three distinct output methods for different deployment scenarios:

**File Output (for batch processing or USB printers):**
```bash
# Generate ESC/POS file
escposimg -image logo.png -output file -file-path printer_data.escpos

# Send to USB printer (Linux/macOS)
cat printer_data.escpos > /dev/usb/lp0

# Send to USB printer (alternative method)
lp -d thermal_printer -o raw printer_data.escpos
```
These commands first generate an ESC/POS file containing the printer commands, then demonstrate two methods for sending the file to a USB-connected thermal printer.

**Network Output (direct TCP connection):**
```bash
# Print directly to network printer
escposimg -image invoice.jpg -output network -network-addr 192.168.1.100:9100
```
This command establishes a direct TCP connection to a network printer and immediately sends the processed image data for printing.

**Standard Output (for shell integration):**
```bash
# Pipe to netcat for network printing
escposimg -image receipt.png | nc 192.168.1.100 9100

# Redirect to file for later use
escposimg -image document.jpg -cut > batch_print.escpos

# Chain with other commands
escposimg -image header.png -debug-text "Order #12345" | tee order_header.escpos | nc printer.local 9100
```
These examples demonstrate shell integration: piping directly to a network printer, saving output for later use, and creating a command chain that both saves and prints simultaneously.

### Go Library

The library provides fine-grained control over image processing and printer communication within Go applications.

#### Basic Library Usage

```go
package main

import (
    "log"
    "github.com/72nd/escposimg"
)

func main() {
    // Load default configuration (80mm, 203 DPI, Floyd-Steinberg)
    config := escposimg.DefaultConfig()
    
    // Create output method
    output := escposimg.NewStdoutOutput()
    
    // Process and print
    if err := escposimg.ProcessImage("receipt.jpg", config, output); err != nil {
        log.Fatal(err)
    }
}
```
This basic example demonstrates how to process an image with default settings and output the ESC/POS commands to standard output within a Go application.

#### Advanced Configuration

```go
package main

import (
    "log"
    "github.com/72nd/escposimg"
)

func main() {
    // Custom configuration for 58mm receipt printer
    config := &escposimg.Config{
        PaperWidthMM:   escposimg.PaperWidth58mm,
        DPI:            escposimg.DPI203,
        DitheringAlgo:  escposimg.DitheringAtkinson,
        PrintMode:      escposimg.PrintModeRaster,
        DebugOutput:    true,
        DebugImagePath: "processed_image.png",
        DebugText:      "Store Receipt #1234",
        CutPaper:       true,
    }
    
    // Network printer output
    output, err := escposimg.NewNetworkOutput("192.168.1.100:9100")
    if err != nil {
        log.Fatal(err)
    }
    defer output.Close()
    
    // Process multiple images in sequence
    images := []string{"header.png", "content.jpg", "footer.png"}
    
    for _, imagePath := range images {
        if err := escposimg.ProcessImage(imagePath, config, output); err != nil {
            log.Printf("Failed to process %s: %v", imagePath, err)
        }
    }
}
```
This advanced example shows how to configure the library for a specific printer type (58mm), enable debug features, and process multiple images sequentially over a network connection.

#### Production Integration Example

```go
package main

import (
    "fmt"
    "log"
    "github.com/72nd/escposimg"
)

func PrintOrderReceipt(orderID string, logoPath string, printerIP string) error {
    // Configure for standard receipt printer
    config := &escposimg.Config{
        PaperWidthMM:  escposimg.PaperWidth80mm,
        DPI:           escposimg.DPI203,
        DitheringAlgo: escposimg.DitheringFloydSteinberg,
        PrintMode:     escposimg.PrintModeRaster,
        DebugText:     fmt.Sprintf("Order: %s", orderID),
        CutPaper:      true,
    }
    
    // Establish printer connection
    printer, err := escposimg.NewNetworkOutput(printerIP + ":9100")
    if err != nil {
        return fmt.Errorf("printer connection failed: %w", err)
    }
    defer printer.Close()
    
    // Process company logo
    if err := escposimg.ProcessImage(logoPath, config, printer); err != nil {
        return fmt.Errorf("logo printing failed: %w", err)
    }
    
    return nil
}
```
This production-ready function demonstrates error handling, configuration management, and network printer integration for printing order receipts with company logos in a commercial application.

## Available Options

### Command-Line Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `-image` | string | *required* | Path to the input image file |
| `-paper-width` | int | `80` | Paper width in millimetres (58, 80, etc.) |
| `-dpi` | int | `203` | Printer resolution in dots per inch |
| `-dithering` | string | `floyd-steinberg` | Dithering algorithm (see table below) |
| `-print-mode` | string | `raster` | ESC/POS printing mode (`raster`, `bit-image`) |
| `-debug-output` | bool | `false` | Save processed image for debugging |
| `-debug-image` | string | `debug_output.png` | Path for debug image output |
| `-debug-text` | string | `` | Optional text printed before image |
| `-cut` | bool | `false` | Send paper cut command after printing |
| `-output` | string | `stdout` | Output method (`stdout`, `network`, `file`) |
| `-network-addr` | string | `` | Network address for network output |
| `-file-path` | string | `` | File path for file output |
| `-verbose` | bool | `false` | Enable detailed logging |
| `-version` | bool | `false` | Display version information |

### Library Configuration Options

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `PaperWidthMM` | int | `80` | Paper width in millimetres |
| `DPI` | int | `203` | Printer dots per inch |
| `DitheringAlgo` | DitheringType | `DitheringFloydSteinberg` | Algorithm for monochrome conversion |
| `PrintMode` | PrintMode | `PrintModeRaster` | ESC/POS command structure |
| `DebugOutput` | bool | `false` | Generate debug image files |
| `DebugImagePath` | string | `debug_output.png` | Debug image save location |
| `DebugText` | string | `` | Text printed before image |
| `CutPaper` | bool | `false` | Automatic paper cutting |

### Dithering Algorithms

| Algorithm | CLI Value | Description | Best For |
|-----------|-----------|-------------|----------|
| Floyd-Steinberg | `floyd-steinberg` | Error diffusion with excellent quality | Photographs, general use |
| Atkinson | `atkinson` | Lighter error diffusion, preserves highlights | Fine details, line art |
| Threshold | `threshold` | Simple binary conversion, fastest processing | High-contrast images, speed |
| Bayer | `bayer` | Ordered dithering with regular patterns | Textures, consistent patterns |
| Burkes | `burkes` | Error diffusion with wider distribution | Complex images, varied tones |
| Sierra Lite | `sierra-lite` | Balanced error diffusion | General purpose, moderate quality |
| Jarvis-Judice-Ninke | `jarvis-judice-ninke` | Comprehensive error diffusion | High-quality output, detailed images |
| Shadura | `shadura` | Optimised for thermal printer characteristics | Thermal printing, bitmap graphics |

### Print Modes

| Mode | CLI Value | Description | Compatibility |
|------|-----------|-------------|---------------|
| Raster | `raster` | Modern GS v 0 command, efficient single-command printing | Modern thermal printers (post-2010) |
| Bit Image | `bit-image` | Legacy ESC * command, line-by-line processing | All ESC/POS printers, including vintage models |

### Common DPI Values

| DPI | Description | Use Case |
|-----|-------------|----------|
| 180 | Lower resolution | Older printers, faster processing |
| 203 | Standard resolution | Most thermal printers, balanced quality |
| 300 | High resolution | Photo-quality printing, detailed graphics |

### Paper Widths

| Width (mm) | Description | Typical Use |
|------------|-------------|-------------|
| 58 | Narrow format | Small receipt printers, mobile devices |
| 80 | Standard format | Retail receipts, standard thermal printers |


