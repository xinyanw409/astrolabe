package core

import (
	"encoding/json"
	"gotest.tools/assert"
	"net/url"
	"reflect"
	"testing"
)

func TestProtectedEntityInfoJSON(t *testing.T) {
	t.Log("TestProtectedEntityIDJSON called")
	const test1Str = "k8s:nginx-example"
	test1ID, test1Err := NewProtectedEntityIDFromString(test1Str)
	if test1Err != nil {
		t.Fatal("Got error " + test1Err.Error())
	}
	dataURL, test1Err := url.Parse("https://data1")
	if test1Err != nil {
		t.Fatal("Got error " + test1Err.Error())
	}
	dataURLs := []url.URL{*dataURL}

	metadataURL, test1Err := url.Parse("https://meta1")
	if test1Err != nil {
		t.Fatal("Got error " + test1Err.Error())
	}
	metadataURLs := []url.URL{*metadataURL}

	combinedURL, test1Err := url.Parse("https://combined1")
	if test1Err != nil {
		t.Fatal("Got error " + test1Err.Error())
	}
	combinedURLs := []url.URL{*combinedURL}

	component1ID, test1Err := NewProtectedEntityIDFromString("ivd:aa-bbb-cc")
	if test1Err != nil {
		t.Fatal("Got error " + test1Err.Error())
	}
	peii := ProtectedEntityInfoImpl{
		Id:           test1ID,
		Name:         "peiiTestJSON",
		DataURLs:     dataURLs,
		MetadataURLs: metadataURLs,
		CombinedURLs: combinedURLs,
		ComponentIDs: []ProtectedEntityID{component1ID},
	}

	jsonBuffer, test1Err := json.Marshal(peii)
	if test1Err != nil {
		t.Fatal("Got error " + test1Err.Error())
	}
	t.Log("jsonStr = " + string(jsonBuffer))
	unmarshalled := ProtectedEntityInfoImpl{}
	test1Err = json.Unmarshal(jsonBuffer, &unmarshalled)
	if test1Err != nil {
		t.Fatal("Got error " + test1Err.Error())
	}
	json2Buffer, test1Err := json.Marshal(unmarshalled)
	if test1Err != nil {
		t.Fatal("Got error " + test1Err.Error())
	}
	t.Log("unmarshalled = " + string(json2Buffer))
	assert.Assert(t, reflect.DeepEqual(peii, unmarshalled), "peii  != unmarshalled")
	//assert.Equal(t, peii, unmarshalled, "peii  != unmarshalled")
}
