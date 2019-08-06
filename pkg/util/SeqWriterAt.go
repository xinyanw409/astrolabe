package util

import "io"

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
	return this.w.Write(p)
}