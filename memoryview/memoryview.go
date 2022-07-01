package memoryview

type MemoryView interface {
	Read(addr int, size int) ([]byte, error)
	Write(addr int, data []byte) error
	SetOffset(offset int64)
}
