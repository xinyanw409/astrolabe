package ivd

import (
	"errors"
	"io"

	vim "github.com/vmware/govmomi/vim25/types"
)

type IVDDataSource struct {
	id vim.ID
}

// The type of this data source, e.g. S3, VADP
func (this *IVDDataSource) GetType() string {
	return "vadp"
}

// Returns a reader that can be used to access the data for this source
func (this *IVDDataSource) GetReader() (io.Reader, error) {
	return nil, errors.New("not supported")
}
