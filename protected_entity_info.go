package arachne

import (
	"net/url"
)

type ProtectedEntityInfoExtern struct {
	Id           string   "id" // JSON tag
	Name         string   "name"
	DataURLs     []string "dataURIs"
	MetadataURIs []string "metadataURIs"
	CombinedURIs []string "combinedURIs"
	ComponentIDs []string "componentIDs"
}

type ProtectedEntityInfo interface {
	GetID() ProtectedEntityID
	GetName() string
	GetDataURLs() []url.URL
	GetCombinedURLs() []url.URL
	GetComponentIDs() []ProtectedEntityID
	/*
	 * Returns the external representation of the ProtectedEntityInfo, ready for
	 * JSON marshaling.  Returns the value rather than a pointer because the JSON
	 * package does not follow pointers so combining this into a larger structure will
	 * require copies anyhow
	 */
	GetExternal() ProtectedEntityInfoExtern
}

type ProtectedEntityInfoImpl struct {
	id           ProtectedEntityID
	name         string
	dataURLs     []url.URL
	metadataURLs []url.URL
	combinedURLs []url.URL
	componentIDs []ProtectedEntityID
}

func NewProtectedEntityInfo(peie ProtectedEntityInfoExtern) (ProtectedEntityInfoImpl, error) {
	var returnErr error
	var newPEI ProtectedEntityInfoImpl
	newPEI.id, returnErr = NewProtectedEntityIDFromString(peie.Id)
	if returnErr == nil {
		newPEI.name = peie.Name
		newPEI.dataURLs, returnErr = stringsToURLs(peie.DataURLs)
	}
	return newPEI, returnErr
}

func (pei *ProtectedEntityInfoImpl) GetID() ProtectedEntityID {
	return pei.id
}

func (pei *ProtectedEntityInfoImpl) GetName() string {
	return pei.name
}

func (pei *ProtectedEntityInfoImpl) GetDataURLs() []url.URL {
	return pei.dataURLs
}

func (pei *ProtectedEntityInfoImpl) GetCombinedURLs() []url.URL {
	return pei.combinedURLs
}

func (pei *ProtectedEntityInfoImpl) GetComponentIDs() []ProtectedEntityID {
	return pei.componentIDs
}

func (pei *ProtectedEntityInfoImpl) GetExternal() ProtectedEntityInfoExtern {
	var newExtern ProtectedEntityInfoExtern
	newExtern.Id = pei.id.String()
	newExtern.Name = pei.GetName()
	newExtern.DataURLs = urlsToStrings(pei.dataURLs)
	return newExtern
}

func urlsToStrings(urls []url.URL) []string {
	returnStrings := make([]string, len(urls))
	for index, curURL := range urls {
		returnStrings[index] = curURL.String()
	}
	return returnStrings
}

func stringsToURLs(strings []string) (urls []url.URL, err error) {
	returnURLs := make([]url.URL, len(strings))
	var returnError error
	for index, curString := range strings {
		curURL, urlErr := url.Parse(curString)
		if urlErr != nil {
			returnError = urlErr
			break
		}
		returnURLs[index] = *curURL
	}
	return returnURLs, returnError
}
