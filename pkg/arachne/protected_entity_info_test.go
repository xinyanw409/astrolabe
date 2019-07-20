package arachne

import (
	"encoding/json"
	"gotest.tools/assert"
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

	component1ID, test1Err := NewProtectedEntityIDFromString("ivd:aa-bbb-cc")
	if test1Err != nil {
		t.Fatal("Got error " + test1Err.Error())
	}
	peii := ProtectedEntityInfoImpl{
		id:           test1ID,
		name:         "peiiTestJSON",
		dataTransports:     []DataTransport {
			NewDataTransportForS3("http://localhost/s3/data1"),
		},
		metadataTransports: []DataTransport {
			NewDataTransportForS3("http://localhost/s3/metadata1"),
		},
		combinedTransports: []DataTransport {
			NewDataTransportForS3("http://localhost/s3/combined1"),
		},
		componentIDs: []ProtectedEntityID{component1ID},
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
