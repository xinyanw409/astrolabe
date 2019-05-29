package arachne

import (
	"testing"
	"gotest.tools/assert"
)

func TestProtectedEntityIDFromString(t *testing.T) {
	t.Log("TestProtectedEntityIDFromString called")
	const test1Str = "k8s:nginx-example"
	test1ID, test1Err := NewProtectedEntityIDFromString(test1Str)
	if (test1Err != nil) {
		t.Error("Got error " + test1Err.Error())
	}
	assert.Equal(t, test1Str, test1ID.String())
	t.Log("test1ID = " + test1ID.String())
	
	// Test with ivd with snapshot
	const test2Str = "ivd:e1c3cb20-db88-4c1c-9f02-5f5347e435d5:67469e1c-50a8-4f63-9a6a-ad8a2265197c"
	test2ID, test2Err := NewProtectedEntityIDFromString(test2Str)
	if (test2Err != nil) {
				t.Error("Got error " + test2Err.Error())
	}
	assert.Equal(t, test2Str, test2ID.String())
	t.Log("test2ID = " + test2ID.String())
}

