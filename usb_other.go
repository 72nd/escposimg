//go:build !unix

package escposimg

import "errors"

// NewUSBOutput is not supported on this platform.
func NewUSBOutput(devicePath string) (*USBOutput, error) {
	return nil, errors.New("USB printer output is only supported on Unix-like platforms")
}
