package arachne

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vmware/gvddk/gDiskLib"
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

type DiskDataWriter struct {
	writerAt io.WriterAt
	connection gDiskLib.VixDiskLibConnection
	params gDiskLib.ConnectParams
	offset   *int64
	mutex    sync.Mutex // Lock to ensure that multiple-threads do not break offset or see the same data twice
	logger   logrus.FieldLogger
}

func NewDiskDataWriter(param DiskConnectionParam, logger logrus.FieldLogger) DiskDataWriter {
	var offset int64
	offset = 0
	retVal := DiskDataWriter{
		writerAt: param.DiskHandle,
		connection: param.VixDiskLibConnection,
		params: param.ConnectParams,
		offset:   &offset,
		logger:   logger,
	}
	return retVal
}

func (this DiskDataWriter) Write(p []byte) (n int, err error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	bytesWritten, err := this.writerAt.WriteAt(p, *this.offset)
	*this.offset += int64(bytesWritten)
	this.logger.Infof("Write returning %d, len(p) = %d, offset=%d\n", bytesWritten, len(p), *this.offset)
	return bytesWritten, err
}

func (this DiskDataWriter) Close() (err error) {
	vErr := this.writerAt.(gDiskLib.DiskHandle).Close()
	if vErr != nil {
		return errors.New(fmt.Sprintf(vErr.Error() + " with error code: %v", vErr))
	}

	vErr = gDiskLib.Disconnect(this.connection)
	if vErr != nil {
		return errors.New(fmt.Sprintf(vErr.Error() + " with error code: %v", vErr))
	}

	vErr = gDiskLib.EndAccess(this.params)
	if vErr != nil {
		return errors.New(fmt.Sprintf(vErr.Error() + " with error code: %v", vErr))
	}

	return nil
}