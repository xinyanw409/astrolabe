package core

import (
	"encoding/json"
	"io"
	"net/http"
)

type S3DataSource struct {
	url string
}

type s3DataSourceJSON struct {
	URL string "json:url"
}
func NewS3DataSource(url string) (S3DataSource, error) {
	return S3DataSource{
		url:url,
	}, nil
}

// The type of this data source, e.g. S3, VADP
func (this *S3DataSource) GetType() string {
	return "s3"
}
// Returns a reader that can be used to access the data for this source
func (this *S3DataSource) GetReader() (io.Reader, error) {
	resp, err := http.Get(this.url)
	if (err != nil) {
		return nil, err
	}
	return resp.Body, nil
}

func (this S3DataSource) MarshalJSON() ([]byte, error) {

	jsonStruct := s3DataSourceJSON{
		URL:           this.url,
	}

	return json.Marshal(jsonStruct)
}

func (this *S3DataSource) UnmarshalJSON(data []byte) error {
	jsonStruct := s3DataSourceJSON{}
	err := json.Unmarshal(data, &jsonStruct)
	if err != nil {
		return err
	}
	this.url = jsonStruct.URL
	return nil
}