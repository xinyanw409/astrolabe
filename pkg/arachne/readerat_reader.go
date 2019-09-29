package arachne

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vmware/gvddk/gDiskLib"
	"io"
	"sync"
)

type ReaderAtReader struct {
	readerAt io.ReaderAt
	offset   *int64
	mutex    sync.Mutex // Lock to ensure that multiple-threads do not break offset or see the same data twice
	logger   logrus.FieldLogger
}


func NewReaderAtReader(readerAt io.ReaderAt, logger logrus.FieldLogger) ReaderAtReader {
	var offset int64
	offset = 0
	retVal := ReaderAtReader{
		readerAt: readerAt,
		offset:   &offset,
		logger: logger,
	}
	return retVal
}

func (this ReaderAtReader) Read(p []byte) (n int, err error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	bytesRead, err := this.readerAt.ReadAt(p, *this.offset)
	*this.offset += int64(bytesRead)
	this.logger.Infof("Read returning %d, len(p) = %d, offset=%d\n", bytesRead, len(p), *this.offset)
	return bytesRead, err
}

type DiskDataReader struct {
	readerAt io.ReaderAt
	connection gDiskLib.VixDiskLibConnection
	params gDiskLib.ConnectParams
	offset  *int64
	mutex    sync.Mutex // Lock to ensure that multiple-threads do not break offset or see the same data twice
	logger   logrus.FieldLogger
}

func NewDiskDataReaderAtReader(param DiskConnectionParam, logger logrus.FieldLogger) DiskDataReader {
	var offset int64
	offset = 0
	retVal := DiskDataReader{
		readerAt: param.DiskHandle,
		connection: param.VixDiskLibConnection,
		params: param.ConnectParams,
		offset:   &offset,
		logger: logger,
	}
	return retVal
}

func (this DiskDataReader) Read(p []byte) (n int, err error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	bytesRead, err := this.readerAt.ReadAt(p, *this.offset)
	*this.offset += int64(bytesRead)
	this.logger.Infof("Read returning %d, len(p) = %d, offset=%d\n", bytesRead, len(p), *this.offset)
	return bytesRead, err
}

func (this DiskDataReader) Close() (err error) {
	vErr := this.readerAt.(gDiskLib.DiskHandle).Close()
	if vErr != nil {
		return errors.New(fmt.Sprintf(vErr.Error() + " with error code: %d", vErr.VixErrorCode()))
	}

	vErr = gDiskLib.Disconnect(this.connection)
	if vErr != nil {
		return errors.New(fmt.Sprintf(vErr.Error() + " with error code: %d", vErr.VixErrorCode()))
	}

	vErr = gDiskLib.EndAccess(this.params)
	if vErr != nil {
		return errors.New(fmt.Sprintf(vErr.Error() + " with error code: %d", vErr.VixErrorCode()))
	}

	return nil
}