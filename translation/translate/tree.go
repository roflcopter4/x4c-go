package translate

import (
	XMLtypes "github.com/lestrrat-go/libxml2/types"
)

type Node struct {
	top *root
}

type root struct {
	Node
}

type documentWrapper struct {
	filename struct {
		base string
		path string
		full string
	}
	doc        XMLtypes.Document
	scriptName string
	lines      []string
}
