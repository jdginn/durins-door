package parser

import ()

type Proxy interface {
	Name() string
	ListChildren() []string
	GetChild(childName string) (*TypeDefProxy, error)
}
