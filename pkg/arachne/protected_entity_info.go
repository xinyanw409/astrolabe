package arachne

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
)

type ProtectedEntityInfo interface {
	GetID() ProtectedEntityID
	GetName() string
	GetDataTransports() [] DataTransport
	GetMetadataTransports() [] DataTransport
	GetCombinedTransports() [] DataTransport
	GetComponentIDs() []ProtectedEntityID
}

type ProtectedEntityInfoImpl struct {
	id                 ProtectedEntityID
	name               string
	dataTransports     []DataTransport
	metadataTransports []DataTransport
	combinedTransports []DataTransport
	componentIDs       []ProtectedEntityID
}

func NewProtectedEntityInfo(id ProtectedEntityID, name string, dataTransports []DataTransport, metadataTransports []DataTransport,
	combinedTransports []DataTransport, componentIDs []ProtectedEntityID) ProtectedEntityInfo {
	return ProtectedEntityInfoImpl{
		id:                 id,
		name:               name,
		dataTransports:     dataTransports,
		metadataTransports: metadataTransports,
		combinedTransports: combinedTransports,
		componentIDs:       componentIDs,
	}
}

type protectedEntityInfoJSON struct {
	Id                 ProtectedEntityID   `json:"id"`
	Name               string              `json:"name"`
	DataTransports     []DataTransport     `json:"dataTransports"`
	MetadataTransports []DataTransport     `json:"metadataTransports"`
	CombinedTransports []DataTransport     `json:"combinedTransports"`
	ComponentIDs       []ProtectedEntityID `json:"componentIDs"`
}

func (this ProtectedEntityInfoImpl) GetID() ProtectedEntityID {
	return this.id
}

func (this ProtectedEntityInfoImpl) GetName() string {
	return this.name
}

func stringsToURLs(urlStrs []string) ([]url.URL, error) {
	retList := []url.URL{}
	for _, curURLStr := range urlStrs {
		curURL, err := url.Parse(curURLStr)
		if err != nil {
			return nil, err
		}
		retList = append(retList, *curURL)
	}
	return retList, nil
}

func urlsToStrings(urls []url.URL) []string {
	retList := []string{}
	for _, curURL := range urls {
		curURLStr := curURL.String()
		retList = append(retList, curURLStr)
	}
	return retList
}

func (this ProtectedEntityInfoImpl) MarshalJSON() ([]byte, error) {

	jsonStruct := protectedEntityInfoJSON{
		Id:                 this.id,
		Name:               this.name,
		DataTransports:     this.dataTransports,
		MetadataTransports: this.metadataTransports,
		CombinedTransports: this.combinedTransports,
		ComponentIDs:       this.componentIDs,
	}

	return json.Marshal(jsonStruct)
}

func appendJSON(buffer *bytes.Buffer, key string, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	buffer.WriteString(fmt.Sprintf("\"%s\":%s", key, string(jsonValue)))
	return nil
}

func (this *ProtectedEntityInfoImpl) UnmarshalJSON(data []byte) error {
	jsonStruct := protectedEntityInfoJSON{}
	err := json.Unmarshal(data, &jsonStruct)
	if err != nil {
		return err
	}
	this.id = jsonStruct.Id
	this.name = jsonStruct.Name
	this.dataTransports = jsonStruct.DataTransports
	this.metadataTransports = jsonStruct.MetadataTransports
	this.combinedTransports = jsonStruct.CombinedTransports
	this.componentIDs = jsonStruct.ComponentIDs
	return nil
}

func (this ProtectedEntityInfoImpl) GetDataTransports() []DataTransport {
	return this.dataTransports
}

func (this ProtectedEntityInfoImpl) GetMetadataTransports() []DataTransport {
	return this.metadataTransports
}

func (this ProtectedEntityInfoImpl) GetCombinedTransports() []DataTransport {
	return this.dataTransports
}

func (this ProtectedEntityInfoImpl) GetComponentIDs() []ProtectedEntityID {
	return this.componentIDs
}
