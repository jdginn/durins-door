//go:build linux

package plat

import (
	"debug/elf"
)

func GetReaderFromFile(f string) (*elf.File, error) {
	return elf.Open(f)
}
