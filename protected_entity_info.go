package arachne

import (
	"bytes"
	"encoding/json"
	"fmt"
	//	"log"
	"net/url"
)

type ProtectedEntityInfo interface {
	GetID() ProtectedEntityID
	GetName() string
	GetDataURLs() []url.URL
	GetCombinedURLs() []url.URL
	GetComponentIDs() []ProtectedEntityID
}

type ProtectedEntityInfoImpl struct {
	Id           ProtectedEntityID   `json:"id"`
	Name         string              `json:"name"`
	DataURLs     []url.URL           `json:"dataURLs"`
	MetadataURLs []url.URL           `json:"metadataURLs"`
	CombinedURLs []url.URL           `json:"combinedURLs"`
	ComponentIDs []ProtectedEntityID `json:"componentIDs"`
}

type protectedEntityInfoJSON struct {
	Id           ProtectedEntityID   `json:"id"`
	Name         string              `json:"name"`
	DataURLs     []string            `json:"dataURLs"`
	MetadataURLs []string            `json:"metadataURLs"`
	CombinedURLs []string            `json:"combinedURLs"`
	ComponentIDs []ProtectedEntityID `json:"componentIDs"`
}

func (this *ProtectedEntityInfoImpl) GetID() ProtectedEntityID {
	return this.Id
}

func (this *ProtectedEntityInfoImpl) GetName() string {
	return this.Name
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
		Id:           this.Id,
		Name:         this.Name,
		DataURLs:     urlsToStrings(this.DataURLs),
		CombinedURLs: urlsToStrings(this.CombinedURLs),
		MetadataURLs: urlsToStrings(this.MetadataURLs),
		ComponentIDs: this.ComponentIDs,
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
	this.Id = jsonStruct.Id
	this.Name = jsonStruct.Name
	this.DataURLs, err = stringsToURLs(jsonStruct.DataURLs)
	if err != nil {
		return err
	}

	this.CombinedURLs, err = stringsToURLs(jsonStruct.CombinedURLs)
	if err != nil {
		return err
	}
	
	this.MetadataURLs, err = stringsToURLs(jsonStruct.MetadataURLs)
	this.ComponentIDs = jsonStruct.ComponentIDs
	return nil
}

func (this *ProtectedEntityInfoImpl) GetDataURLs() []url.URL {
	return this.DataURLs
}

func (this *ProtectedEntityInfoImpl) GetCombinedURLs() []url.URL {
	return this.CombinedURLs
}

func (this *ProtectedEntityInfoImpl) GetComponentIDs() []ProtectedEntityID {
	return this.ComponentIDs
}
