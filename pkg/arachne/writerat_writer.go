package arachne

import (
	"fmt"
	"io"
	"sync"
)

type WriterAtWriter struct {
	writerAt io.WriterAt
	offset   *int64
	mutex    sync.Mutex	// Lock to ensure that multiple-threads do not break offset or see the same data twice
}

func NewWriterAtWriter(writerAt io.WriterAt) WriterAtWriter {
	var offset int64
	offset = 0
	retVal := WriterAtWriter{
		writerAt: writerAt,
		offset:   &offset,
	}
	return retVal
}

func (this WriterAtWriter) Write(p []byte) (n int, err error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	bytesWritten, err := this.writerAt.WriteAt(p, *this.offset)
	*this.offset += int64(bytesWritten)
	fmt.Printf("Write returning %d, len(p) = %d\n", bytesWritten, len(p))
	return bytesWritten, err
}