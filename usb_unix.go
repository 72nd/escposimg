//go:build unix

package escposimg

import (
	"fmt"
	"os"
	"syscall"
)

// NewUSBOutput opens a USB printer device for writing.
// devicePath is the path to the device file, e.g. "/dev/usb/lp0".
//
// On Linux the usblp character device is pollable, so the Go runtime registers
// the file descriptor with the netpoller and switches it to non-blocking mode.
// In that mode a write returns before the data has been fully transmitted, and
// the subsequent Close aborts the in-flight transfer — so nothing prints even
// though every call reports success. Forcing the descriptor back to blocking
// mode makes write(2) wait for the transfer to complete, matching the behaviour
// of `cat file > /dev/usb/lp0`.
func NewUSBOutput(devicePath string) (*USBOutput, error) {
	file, err := os.OpenFile(devicePath, os.O_WRONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to open USB device %s: %w", devicePath, err)
	}
	if err := syscall.SetNonblock(int(file.Fd()), false); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("failed to set blocking mode on USB device %s: %w", devicePath, err)
	}
	return &USBOutput{file: file}, nil
}
