package file

import (
	"fmt"
	"os"
)

type FileClient struct {
	rw     *os.File
	offset int64
}

func New(f *os.File) (*FileClient, error) {
	fw := &FileClient{
		rw: f,
	}
	return fw, nil
}

func NewFromPath(filename string) (*FileClient, error) {
	f, err := os.Open(filename)
	if err != nil {
		return &FileClient{}, fmt.Errorf("Could not open file %s:\n\n\t%s", filename, err)
	}
	fw := &FileClient{
		rw:     f,
		offset: 0,
	}
	return fw, nil
}

// TODO: setting offset in client seems like a bad idea
func (p *FileClient) SetOffset(offset int64) {
	p.offset = offset
}

func (p *FileClient) Read(addr int, size int) ([]byte, error) {
	val := make([]byte, size)
	n, err := p.rw.ReadAt(val, int64(addr)-p.offset)
	if err != nil {
		return val, err
	}
	if n != size {
		return val, fmt.Errorf("Read the incorrect number of bytes\n Expected: %d bytes; Read %d", size, n)
	}
	return val, nil
}

func (p *FileClient) Write(addr int, data []byte) error {
	_, err := p.rw.WriteAt(data, int64(addr)-p.offset)
	return err
}
