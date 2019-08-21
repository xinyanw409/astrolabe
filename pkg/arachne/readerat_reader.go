package arachne

import (
	"io"
	"sync"
)

type ReaderAtReader struct {
	readerAt io.ReaderAt
	offset   *int64
	mutex    sync.Mutex	// Lock to ensure that multiple-threads do not break offset or see the same data twice
}

func NewReaderAtReader(readerAt io.ReaderAt) ReaderAtReader {
	var offset int64
	offset = 0
	retVal := ReaderAtReader{
		readerAt: readerAt,
		offset:   &offset,
	}
	return retVal
}

func (this ReaderAtReader) Read(p []byte) (n int, err error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	bytesRead, err := this.readerAt.ReadAt(p, *this.offset)
	*this.offset += int64(bytesRead)
	return bytesRead, err
}
