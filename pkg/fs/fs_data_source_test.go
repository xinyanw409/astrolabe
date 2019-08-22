package fs

import (
	"fmt"
	"github.com/vmware/arachne/pkg/arachne"
	"io"
	"testing"
)

func TestFSDataSource(t *testing.T) {
	t.Log("TestFSDataSource called")

	fs, err := newFSProtectedEntity(nil, arachne.ProtectedEntityID{}, "test", "/Users/dsmithuchida/Downloads")
	if err != nil {
		t.Fatal("Got error " + err.Error())
	}
	fsReader, err := fs.GetDataReader(nil)
	if err != nil {
		t.Fatal("Got error " + err.Error())
	}

	buf := make([]byte, 1024*1024)
	keepReading := true
	for keepReading {
		bytesRead, err := fsReader.Read(buf)
		fmt.Printf("Read %d bytes\n", bytesRead)
		if err != nil && err != io.EOF {
			t.Fatal("Got error " + err.Error())
		}
		if bytesRead == 0 && err == io.EOF {
			keepReading = false
		}
	}
	fmt.Printf("Finished reading\n")
}
