package escposimg

import (
	"fmt"
	"net"
	"os"
)

// StdoutOutput writes data to stdout
type StdoutOutput struct{}

// NewStdoutOutput creates a new stdout output method
func NewStdoutOutput() *StdoutOutput {
	return &StdoutOutput{}
}

// Write writes data to stdout
func (s *StdoutOutput) Write(data []byte) error {
	_, err := os.Stdout.Write(data)
	return err
}

// Close is a no-op for stdout
func (s *StdoutOutput) Close() error {
	return nil
}

// NetworkOutput writes data to a network connection
type NetworkOutput struct {
	conn net.Conn
}

// NewNetworkOutput creates a new network output method
func NewNetworkOutput(address string) (*NetworkOutput, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
	}
	return &NetworkOutput{conn: conn}, nil
}

// Write writes data to the network connection
func (n *NetworkOutput) Write(data []byte) error {
	_, err := n.conn.Write(data)
	return err
}

// Close closes the network connection
func (n *NetworkOutput) Close() error {
	return n.conn.Close()
}

// FileOutput writes data to a file
type FileOutput struct {
	file *os.File
}

// NewFileOutput creates a new file output method
func NewFileOutput(filePath string) (*FileOutput, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	return &FileOutput{file: file}, nil
}

// Write writes data to the file
func (f *FileOutput) Write(data []byte) error {
	_, err := f.file.Write(data)
	return err
}

// Close closes the file
func (f *FileOutput) Close() error {
	return f.file.Close()
}
