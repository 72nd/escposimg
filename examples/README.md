# escposimg Examples

This directory contains examples demonstrating how to use the escposimg library. Each example focuses on specific features and use cases.

## Prerequisites

Make sure you have the test image available:
```bash
# The examples expect test_pattern.png to be in the examples directory
# You can use any PNG or JPEG image for testing
```

## Examples Overview

### 1. `basic_usage.go` - Getting Started

**Purpose**: Shows the simplest way to use escposimg with default settings.

**What it demonstrates**:
- Loading an image
- Using default configuration
- Outputting to stdout

**Run it**:
```bash
go run basic_usage.go
```

**Key concepts**:
- `escposimg.DefaultConfig()` - sensible defaults
- `escposimg.NewStdoutOutput()` - output to console
- `escposimg.ProcessImage()` - main processing function

### 2. `dithering_comparison.go` - Dithering Algorithms

**Purpose**: Compares all available dithering algorithms side by side.

**What it demonstrates**:
- All 8 dithering algorithms
- Debug image output for visual comparison
- Algorithm selection based on image content

**Run it**:
```bash
go run dithering_comparison.go
```

**Output**: Creates PNG files showing the result of each dithering algorithm.

**Available algorithms**:
- `floyd-steinberg` - Best overall quality for photos
- `atkinson` - Good for images with fine details  
- `threshold` - Fast, good for high-contrast images
- `bayer` - Good for textures and patterns
- `burkes` - Good detail preservation
- `sierra-lite` - Fast error diffusion
- `jarvis-judice-ninke` - High quality, slower
- `shadura` - Custom algorithm based on png2pos.c

### 3. `output_methods.go` - Output Destinations

**Purpose**: Shows how to send ESC/POS commands to different destinations.

**What it demonstrates**:
- Stdout output (for piping/redirection)
- File output (for saving commands)
- Network output (for direct printer communication)

**Run it**:
```bash
go run output_methods.go
```

**Output methods**:
- `StdoutOutput` - Command line integration
- `FileOutput` - Save for later use
- `NetworkOutput` - Direct to network printer (IP:port)

### 4. `configuration_options.go` - Printer Settings

**Purpose**: Demonstrates different printer configurations and settings.

**What it demonstrates**:
- Different paper widths (58mm, 80mm)
- Different DPI settings (180, 203, 300)
- Pixel width calculations
- Configuration customization

**Run it**:
```bash
go run configuration_options.go
```

**Key settings**:
- `PaperWidthMM` - 58mm or 80mm typically
- `DPI` - 203 for standard, 300 for high quality
- `DitheringAlgo` - Algorithm selection
- `CutPaper` - Auto-cut after printing

### 5. `debug_features.go` - Debugging and Troubleshooting

**Purpose**: Shows debugging features to help troubleshoot issues.

**What it demonstrates**:
- Debug image output
- Debug text on receipts
- Verbose logging
- Multiple configurations for comparison

**Run it**:
```bash
go run debug_features.go
```

**Debug features**:
- `DebugOutput` - Save processed images
- `DebugText` - Add text to printouts
- Verbose logging with `slog.LevelDebug`
- Visual comparison of settings

## Running the Examples

### Basic Usage
```bash
# Run any example
go run basic_usage.go

# Redirect output to file
go run basic_usage.go > output.escpos

# View ESC/POS commands in hex
go run basic_usage.go 2>/dev/null | hexdump -C
```

### With Your Own Images
```bash
# Copy your image to the examples directory
cp /path/to/your/image.jpg test_pattern.png

# Or modify the examples to use your image path
```

### Sending to Real Printers
```bash
# Network printer
go run output_methods.go  # Edit IP address first

# USB printer (Linux)
go run basic_usage.go > /dev/usb/lp0

# Serial printer
go run basic_usage.go > /dev/ttyUSB0
```

## Common Use Cases

### Receipt Printing
Use `basic_usage.go` or `configuration_options.go` with 80mm paper width.

### Label Printing  
Use `configuration_options.go` with appropriate paper width for your labels.

### Photo Printing
Use `dithering_comparison.go` to find the best algorithm for your images.

### Debugging Print Issues
Use `debug_features.go` to save debug images and troubleshoot problems.

### Batch Processing
Use `output_methods.go` with file output to process many images.

## Tips

1. **Start with `basic_usage.go`** to verify everything works
2. **Use debug output** when developing - always enable `DebugOutput: true`
3. **Test with files first** before sending to real printers
4. **Compare dithering algorithms** - different images work better with different algorithms
5. **Check your printer specs** - verify paper width and DPI settings
6. **Use appropriate paper width** - 58mm for small receipts, 80mm for standard receipts

## Troubleshooting

### Image too light or too dark?
Try different dithering algorithms with `dithering_comparison.go`.

### Wrong image size?
Check `PaperWidthMM` and `DPI` settings in `configuration_options.go`.

### Printer not responding?
Test with file output first, then check network connectivity.

### Need to see processing details?
Use verbose logging in `debug_features.go`.
