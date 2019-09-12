package arachne

import (
	"github.com/sirupsen/logrus"
	"io"
	"sync"
)

type WriterAtWriter struct {
	writerAt io.WriterAt
	offset   *int64
	mutex    sync.Mutex // Lock to ensure that multiple-threads do not break offset or see the same data twice
	logger   logrus.FieldLogger
}

func NewWriterAtWriter(writerAt io.WriterAt, logger logrus.FieldLogger) WriterAtWriter {
	var offset int64
	offset = 0
	retVal := WriterAtWriter{
		writerAt: writerAt,
		offset:   &offset,
		logger:   logger,
	}
	return retVal
}

func (this WriterAtWriter) Write(p []byte) (n int, err error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	bytesWritten, err := this.writerAt.WriteAt(p, *this.offset)
	*this.offset += int64(bytesWritten)
	this.logger.Infof("Write returning %d, len(p) = %d, offset=%d\n", bytesWritten, len(p), *this.offset)
	return bytesWritten, err
}
