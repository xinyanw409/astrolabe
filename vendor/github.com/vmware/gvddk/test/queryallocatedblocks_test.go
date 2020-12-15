package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vmware/gvddk/gDiskLib"
	"github.com/vmware/gvddk/gvddk-high"
	"os"
	"testing"
)

func TestQueryAllocatedBlocks(t *testing.T) {
	fmt.Println("Test OueryAllocatedBocks starts")
	var majorVersion uint32 = 7
	var minorVersion uint32 = 0
	path := os.Getenv("LIBPATH")
	if path == "" {
		t.Skip("Skipping testing if environment variables are not set.")
	}
	gDiskLib.Init(majorVersion, minorVersion, path)
	serverName := os.Getenv("IP")
	thumPrint := os.Getenv("THUMBPRINT")
	userName := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	fcdId := os.Getenv("FCDID")
	ds := os.Getenv("DATASTORE")
	identity := os.Getenv("IDENTITY")
	params := gDiskLib.NewConnectParams("", serverName,thumPrint, userName,
		password, fcdId, ds, "", "", identity, "", gDiskLib.VIXDISKLIB_FLAG_OPEN_COMPRESSION_SKIPZ,
		false, gDiskLib.NBD)
	diskReaderWriter, err := gvddk_high.Open(params, logrus.New())
	if err != nil {
		gDiskLib.EndAccess(params)
		t.Errorf("Open failed, got error code: %d, error message: %s.", err.VixErrorCode(), err.Error())
	}

	queryAllocatedBlocks(diskReaderWriter, params)
}

func queryAllocatedBlocks(diskReaderWriter gvddk_high.DiskReaderWriter, params gDiskLib.ConnectParams) {
	var offset int64
	var capacity int64 = 10240 // Disk size in MB

	chunkSize := gDiskLib.VIXDISKLIB_MIN_CHUNK_SIZE
	numChunk := capacity / int64(chunkSize)
	var numChunkToQuery int64
	var abFinal []gDiskLib.VixDiskLibBlock
	for numChunk > 0 {
		if numChunk > gDiskLib.VIXDISKLIB_MAX_CHUNK_NUMBER {
			numChunkToQuery = gDiskLib.VIXDISKLIB_MAX_CHUNK_NUMBER
		} else {
			numChunkToQuery = int64(numChunk)
		}
		abList, err := diskReaderWriter.QueryAllocatedBlocks(gDiskLib.VixDiskLibSectorType(offset), gDiskLib.VixDiskLibSectorType(numChunkToQuery)*gDiskLib.VixDiskLibSectorType(chunkSize), gDiskLib.VixDiskLibSectorType(chunkSize))
		if err != nil {
			gDiskLib.EndAccess(params)
			fmt.Errorf("QueryAllocatedBlocks failed, got error code: %d, error message: %s.", err.VixErrorCode(), err.Error())
			return
		}
		fmt.Printf("Number of blocks: %d\n", len(abList))
		fmt.Printf("Offset      Length\n")
		for _, ab := range abList {
			fmt.Printf("0x%012x  0x%012x\n", ab.Offset(), ab.Length())
			abFinal = append(abFinal, ab)
		}

		numChunk = numChunk - numChunkToQuery
		offset = offset + numChunkToQuery*int64(chunkSize)
	}

	allocatedSize := 0
	for _, ab := range abFinal {
		fmt.Printf("0x%012x  0x%012x\n", ab.Offset(), ab.Length())
		allocatedSize = allocatedSize + int(ab.Length())
	}
	fmt.Printf("Allocated size is %d / capacity %d", allocatedSize, capacity)
}