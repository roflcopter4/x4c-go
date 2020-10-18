package util

import (
	"bytes"
)

type bslice [][]byte

//========================================================================================

func (slc *bslice) Len() int {
	return len(*slc)
}

func (slc *bslice) Less(i, x int) bool {
	return bytes.Compare((*slc)[i], (*slc)[x]) < 0
}

func (slc *bslice) Swap(i, x int) {
	tmp := (*slc)[i]
	(*slc)[i] = (*slc)[x]
	(*slc)[x] = tmp
}
