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
	f, err := os.OpenFile(filename, os.O_RDWR, 0777)
	if err != nil {
    wd, _ := os.Getwd()
    return &FileClient{}, fmt.Errorf("Could not open file %s:\n\ncwd:  %s\n\nerror message:\n\t%s", filename, wd, err)
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
  relAddr := int64(addr) - p.offset
	n, err := p.rw.ReadAt(val, relAddr)
	if err != nil {
		return val, err
	}
	if n != size {
		return val, fmt.Errorf("Read the incorrect number of bytes\n Expected: %d bytes; Read %d", size, n)
	}
	return val, nil
}

func (p *FileClient) Write(addr int, data []byte) error {
  relAddr := int64(addr) - p.offset
	_, err := p.rw.WriteAt(data, relAddr)
	return err
}
