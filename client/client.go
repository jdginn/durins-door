package client

type Client interface {
	Read(addr int, size int) ([]byte, error)
	Write(addr int, data []byte) error
}
