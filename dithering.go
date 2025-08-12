package escposimg

import (
	"image"
	"image/color"
	"log/slog"
)

// ApplyDithering applies the specified dithering algorithm to the image
func ApplyDithering(img image.Image, algo DitheringType) (image.Image, error) {
	slog.Debug("Applying dithering algorithm", "algorithm", algo.String())

	switch algo {
	case DitheringFloydSteinberg:
		return applyFloydSteinberg(img)
	case DitheringAtkinson:
		return applyAtkinson(img)
	case DitheringThreshold:
		return applyThreshold(img)
	case DitheringBayer:
		return applyBayer(img)
	case DitheringBurkes:
		return applyBurkes(img)
	case DitheringSierraLite:
		return applySierraLite(img)
	case DitheringJarvisJudiceNinke:
		return applyJarvisJudiceNinke(img)
	case DitheringShadura:
		return applyShadura(img)
	default:
		slog.Warn("Unknown dithering algorithm, falling back to Floyd-Steinberg", "algorithm", algo)
		return applyFloydSteinberg(img)
	}
}

// convertToGrayscale converts an image to grayscale values
func convertToGrayscale(img image.Image) [][]uint8 {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	gray := make([][]uint8, height)
	for y := 0; y < height; y++ {
		gray[y] = make([]uint8, width)
		for x := 0; x < width; x++ {
			// Get pixel color and convert to grayscale using luminance formula
			r, g, b, _ := img.At(x+bounds.Min.X, y+bounds.Min.Y).RGBA()
			// Convert from 16-bit to 8-bit and apply luminance weights
			grayValue := uint8((0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)))
			gray[y][x] = grayValue
		}
	}
	return gray
}

// createMonochromeImage creates a black and white image from a boolean matrix
func createMonochromeImage(pixels [][]bool, width, height int) image.Image {
	img := image.NewGray(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if pixels[y][x] {
				img.SetGray(x, y, color.Gray{Y: 0}) // Black
			} else {
				img.SetGray(x, y, color.Gray{Y: 255}) // White
			}
		}
	}
	return img
}

// applyFloydSteinberg implements Floyd-Steinberg dithering
func applyFloydSteinberg(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Convert to grayscale
	gray := convertToGrayscale(img)

	// Convert to float64 for error diffusion calculations
	pixels := make([][]float64, height)
	for y := 0; y < height; y++ {
		pixels[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			pixels[y][x] = float64(gray[y][x])
		}
	}

	result := make([][]bool, height)
	for y := 0; y < height; y++ {
		result[y] = make([]bool, width)
	}

	// Apply Floyd-Steinberg dithering
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			oldPixel := pixels[y][x]
			var newPixel float64
			var isBlack bool

			if oldPixel < 128 {
				newPixel = 0
				isBlack = true
			} else {
				newPixel = 255
				isBlack = false
			}

			result[y][x] = isBlack
			quantError := oldPixel - newPixel

			// Distribute error to neighboring pixels
			if x+1 < width {
				pixels[y][x+1] += quantError * 7.0 / 16.0
			}
			if y+1 < height {
				if x > 0 {
					pixels[y+1][x-1] += quantError * 3.0 / 16.0
				}
				pixels[y+1][x] += quantError * 5.0 / 16.0
				if x+1 < width {
					pixels[y+1][x+1] += quantError * 1.0 / 16.0
				}
			}
		}
	}

	return createMonochromeImage(result, width, height), nil
}

// applyAtkinson implements Atkinson dithering
func applyAtkinson(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	gray := convertToGrayscale(img)

	pixels := make([][]float64, height)
	for y := 0; y < height; y++ {
		pixels[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			pixels[y][x] = float64(gray[y][x])
		}
	}

	result := make([][]bool, height)
	for y := 0; y < height; y++ {
		result[y] = make([]bool, width)
	}

	// Apply Atkinson dithering
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			oldPixel := pixels[y][x]
			var newPixel float64
			var isBlack bool

			if oldPixel < 128 {
				newPixel = 0
				isBlack = true
			} else {
				newPixel = 255
				isBlack = false
			}

			result[y][x] = isBlack
			quantError := oldPixel - newPixel

			// Atkinson dithering pattern (error distributed to 6 neighbors)
			if x+1 < width {
				pixels[y][x+1] += quantError / 8.0
			}
			if x+2 < width {
				pixels[y][x+2] += quantError / 8.0
			}
			if y+1 < height {
				if x > 0 {
					pixels[y+1][x-1] += quantError / 8.0
				}
				pixels[y+1][x] += quantError / 8.0
				if x+1 < width {
					pixels[y+1][x+1] += quantError / 8.0
				}
			}
			if y+2 < height {
				pixels[y+2][x] += quantError / 8.0
			}
		}
	}

	return createMonochromeImage(result, width, height), nil
}

// applyThreshold implements simple threshold dithering
func applyThreshold(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	gray := convertToGrayscale(img)
	result := make([][]bool, height)

	for y := 0; y < height; y++ {
		result[y] = make([]bool, width)
		for x := 0; x < width; x++ {
			// Simple threshold at 128
			result[y][x] = gray[y][x] < 128
		}
	}

	return createMonochromeImage(result, width, height), nil
}

// applyBayer implements Bayer matrix dithering (4x4)
func applyBayer(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 4x4 Bayer matrix
	bayerMatrix := [][]int{
		{0, 8, 2, 10},
		{12, 4, 14, 6},
		{3, 11, 1, 9},
		{15, 7, 13, 5},
	}

	gray := convertToGrayscale(img)
	result := make([][]bool, height)

	for y := 0; y < height; y++ {
		result[y] = make([]bool, width)
		for x := 0; x < width; x++ {
			threshold := bayerMatrix[y%4][x%4] * 16
			result[y][x] = int(gray[y][x]) < threshold
		}
	}

	return createMonochromeImage(result, width, height), nil
}

// applyBurkes implements Burkes dithering
func applyBurkes(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	gray := convertToGrayscale(img)

	pixels := make([][]float64, height)
	for y := 0; y < height; y++ {
		pixels[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			pixels[y][x] = float64(gray[y][x])
		}
	}

	result := make([][]bool, height)
	for y := 0; y < height; y++ {
		result[y] = make([]bool, width)
	}

	// Apply Burkes dithering
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			oldPixel := pixels[y][x]
			var newPixel float64
			var isBlack bool

			if oldPixel < 128 {
				newPixel = 0
				isBlack = true
			} else {
				newPixel = 255
				isBlack = false
			}

			result[y][x] = isBlack
			quantError := oldPixel - newPixel

			// Burkes dithering pattern
			if x+1 < width {
				pixels[y][x+1] += quantError * 8.0 / 32.0
			}
			if x+2 < width {
				pixels[y][x+2] += quantError * 4.0 / 32.0
			}
			if y+1 < height {
				if x-2 >= 0 {
					pixels[y+1][x-2] += quantError * 2.0 / 32.0
				}
				if x-1 >= 0 {
					pixels[y+1][x-1] += quantError * 4.0 / 32.0
				}
				pixels[y+1][x] += quantError * 8.0 / 32.0
				if x+1 < width {
					pixels[y+1][x+1] += quantError * 4.0 / 32.0
				}
				if x+2 < width {
					pixels[y+1][x+2] += quantError * 2.0 / 32.0
				}
			}
		}
	}

	return createMonochromeImage(result, width, height), nil
}

// applySierraLite implements Sierra Lite dithering (Sierra-2-4A)
func applySierraLite(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	gray := convertToGrayscale(img)

	pixels := make([][]float64, height)
	for y := 0; y < height; y++ {
		pixels[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			pixels[y][x] = float64(gray[y][x])
		}
	}

	result := make([][]bool, height)
	for y := 0; y < height; y++ {
		result[y] = make([]bool, width)
	}

	// Apply Sierra Lite dithering
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			oldPixel := pixels[y][x]
			var newPixel float64
			var isBlack bool

			if oldPixel < 128 {
				newPixel = 0
				isBlack = true
			} else {
				newPixel = 255
				isBlack = false
			}

			result[y][x] = isBlack
			quantError := oldPixel - newPixel

			// Sierra Lite dithering pattern
			if x+1 < width {
				pixels[y][x+1] += quantError * 2.0 / 4.0
			}
			if y+1 < height {
				if x-1 >= 0 {
					pixels[y+1][x-1] += quantError * 1.0 / 4.0
				}
				pixels[y+1][x] += quantError * 1.0 / 4.0
			}
		}
	}

	return createMonochromeImage(result, width, height), nil
}

// applyJarvisJudiceNinke implements Jarvis-Judice-Ninke dithering
func applyJarvisJudiceNinke(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	gray := convertToGrayscale(img)

	pixels := make([][]float64, height)
	for y := 0; y < height; y++ {
		pixels[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			pixels[y][x] = float64(gray[y][x])
		}
	}

	result := make([][]bool, height)
	for y := 0; y < height; y++ {
		result[y] = make([]bool, width)
	}

	// Apply Jarvis-Judice-Ninke dithering
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			oldPixel := pixels[y][x]
			var newPixel float64
			var isBlack bool

			if oldPixel < 128 {
				newPixel = 0
				isBlack = true
			} else {
				newPixel = 255
				isBlack = false
			}

			result[y][x] = isBlack
			quantError := oldPixel - newPixel

			// Jarvis-Judice-Ninke dithering pattern
			if x+1 < width {
				pixels[y][x+1] += quantError * 7.0 / 48.0
			}
			if x+2 < width {
				pixels[y][x+2] += quantError * 5.0 / 48.0
			}
			if y+1 < height {
				if x-2 >= 0 {
					pixels[y+1][x-2] += quantError * 3.0 / 48.0
				}
				if x-1 >= 0 {
					pixels[y+1][x-1] += quantError * 5.0 / 48.0
				}
				pixels[y+1][x] += quantError * 7.0 / 48.0
				if x+1 < width {
					pixels[y+1][x+1] += quantError * 5.0 / 48.0
				}
				if x+2 < width {
					pixels[y+1][x+2] += quantError * 3.0 / 48.0
				}
			}
			if y+2 < height {
				if x-2 >= 0 {
					pixels[y+2][x-2] += quantError * 1.0 / 48.0
				}
				if x-1 >= 0 {
					pixels[y+2][x-1] += quantError * 3.0 / 48.0
				}
				pixels[y+2][x] += quantError * 5.0 / 48.0
				if x+1 < width {
					pixels[y+2][x+1] += quantError * 3.0 / 48.0
				}
				if x+2 < width {
					pixels[y+2][x+2] += quantError * 1.0 / 48.0
				}
			}
		}
	}

	return createMonochromeImage(result, width, height), nil
}

// applyShadura implements a simplified version of the Shadura algorithm
// Based on the png2pos.c implementation approach
func applyShadura(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	gray := convertToGrayscale(img)

	pixels := make([][]float64, height)
	for y := 0; y < height; y++ {
		pixels[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			pixels[y][x] = float64(gray[y][x])
		}
	}

	result := make([][]bool, height)
	for y := 0; y < height; y++ {
		result[y] = make([]bool, width)
	}

	// Apply Shadura-style dithering (simplified error diffusion)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			oldPixel := pixels[y][x]
			var newPixel float64
			var isBlack bool

			if oldPixel < 128 {
				newPixel = 0
				isBlack = true
			} else {
				newPixel = 255
				isBlack = false
			}

			result[y][x] = isBlack
			quantError := oldPixel - newPixel

			// Shadura-style error distribution (simplified pattern)
			if x+1 < width {
				pixels[y][x+1] += quantError * 0.5
			}
			if y+1 < height {
				pixels[y+1][x] += quantError * 0.5
			}
		}
	}

	return createMonochromeImage(result, width, height), nil
}
