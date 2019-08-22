package util

import (
	"fmt"
	"io"
)

type SeqWriterAt struct {
	w io.Writer
}

func NewSeqWriterAt(w io.Writer) SeqWriterAt {
	return SeqWriterAt{
		w: w,
	}
}
func (this SeqWriterAt) WriteAt(p []byte, offset int64) (n int, err error) {
	// ignore 'offset' because we forced sequential downloads
	fmt.Printf("SeqWriterAt WriteAt %d at %d\n", len(p), offset)
	return this.w.Write(p)
}