package fs

import (
	"fmt"
	"io"
	"testing"
)

func TestFSDataSource(t *testing.T) {
	t.Log("TestFSDataSource called")

	fs, err := NewFSDataSource("/Users/dsmithuchida/Downloads")
	if err != nil {
		t.Fatal("Got error " + err.Error())
	}
	fsReader, err := fs.GetReader()
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
