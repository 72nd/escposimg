package escposimg

import (
	"image"
)

// GenerateESCPOS generates ESC/POS commands from a dithered image
func GenerateESCPOS(img image.Image, config *Config) ([]byte, error) {
	// TODO: Implement ESC/POS command generation
	// For now, return empty byte slice
	return []byte{}, nil
}
