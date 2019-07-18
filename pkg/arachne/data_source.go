package arachne

import "io"

// DataSource is our internal interface representing the data for Protected Entity
// data, metadata or combined info
type DataSource interface {
	// The type of this data source, e.g. S3, VADP
	GetType() string
	// Returns a reader that can be used to access the data for this source
	GetReader() (io.Reader, error)
	// Returns the JSON representation of this DataSource
	GetJSON() []byte
}
