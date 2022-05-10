//go:build darwin

package plat

import (
	"debug/macho"
)

func GetReaderFromFile(f string) (*macho.File, error) {
	return macho.Open(f)
}
