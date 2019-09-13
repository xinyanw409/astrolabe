package arachne

import (
	"encoding/json"
)

// DataTransport is our internal interface representing the data transport for Protected Entity
// data, metadata or combined info
// DataTransport contains parameters for the transport but does not actually move data
// DataTransport is used in two ways:
//		Each ProtectedEntity exports a set of DataTransports for accessing its data, metadata, and combined streams
//		These exported DataTransports are used to form the JSON
//
//		When we copy from a ProtectedEntity, the DataTransports of the source may be used by the ProtectedEntity to
//		return a stream.  This is most useful for remote ProtectedEntities.

type DataTransport struct {
	// The type of this data source, e.g. S3, VADP
	transportType string
	params        map[string]string
}

type dataTransportJSON struct {
	TransportType string            `json:"transportType"`
	Params        map[string]string `json:"params"`
}

func NewDataTransport(transportType string, params map[string]string) DataTransport {
	return DataTransport{
		transportType: transportType,
		params:        params,
	}
}

func NewDataTransportForS3URL(url string) DataTransport {
	return DataTransport{
		transportType: "s3",
		params: map[string]string{
			"url": url,
		},
	}
}

func NewDataTransportForS3(host string, bucket string, key string) DataTransport {
	url := "http://" + host + "/" + bucket + "/" + key
	return DataTransport{
		transportType: "s3",
		params: map[string]string{
			"url":    url,
			"host":   host,
			"bucket": bucket,
			"key":    key,
		},
	}
}
func (this DataTransport) GetTransportType() string {
	return this.transportType
}

func (this DataTransport) GetParam(key string) (string, bool) {
	val, ok := this.params[key]
	return val, ok
}

func (this DataTransport) MarshalJSON() ([]byte, error) {

	jsonStruct := dataTransportJSON{
		TransportType: this.transportType,
		Params:        this.params,
	}

	return json.Marshal(jsonStruct)
}

func (this *DataTransport) UnmarshalJSON(data []byte) error {
	jsonStruct := dataTransportJSON{}
	err := json.Unmarshal(data, &jsonStruct)
	if err != nil {
		return err
	}
	this.transportType = jsonStruct.TransportType
	this.params = jsonStruct.Params
	return nil
}
