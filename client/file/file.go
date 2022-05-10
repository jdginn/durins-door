package file

import (
	"fmt"
	"os"
)

type FileClient struct {
	rw *os.File
}

func NewFileClient(f *os.File) (*FileClient, error) {
	fw := &FileClient{
		rw: f,
	}
	return fw, nil
}

func NewFileClientFromName(filename string) (*FileClient, error) {
	f, err := os.Open(filename)
	if err != nil {
		return &FileClient{}, fmt.Errorf("Could not open file %s:\n\n\t%s", filename, err)
	}
	fw := &FileClient{
		rw: f,
	}
	return fw, nil
}

func (p *FileClient) Read(addr int, size int) ([]byte, error) {
	val := make([]byte, size)
	n, err := p.rw.ReadAt(val, int64(addr))
	if err != nil {
		return val, err
	}
	if n != size {
		return val, fmt.Errorf("Read the incorrect number of bytes\n Expected: %d bytes; Read %d", size, n)
	}
	return val, nil
}

func (p *FileClient) Write(addr int, data []byte) error {
	_, err := p.rw.WriteAt(data, int64(addr))
	return err
}
