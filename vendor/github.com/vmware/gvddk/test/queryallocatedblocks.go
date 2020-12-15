package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vmware/gvddk/gDiskLib"
	"github.com/vmware/gvddk/gvddk-high"
	"os"
)

func main() {
	fmt.Println("Test OueryAllocatedBocks starts")
	var majorVersion uint32 = 7
	var minorVersion uint32 = 0
	path := os.Getenv("LIBPATH")
	gDiskLib.Init(majorVersion, minorVersion, path)
	//serverName := os.Getenv("IP")
	//thumPrint := os.Getenv("THUMBPRINT")
	////thumPrint = GetThumbPrintForServer(serverName)
	//userName := os.Getenv("USERNAME")
	//password := os.Getenv("PASSWORD")
	//fcdId := os.Getenv("FCDID")
	//ds := os.Getenv("DATASTORE")
	//identity := os.Getenv("IDENTITY")
	params := gDiskLib.NewConnectParams("", "10.185.34.44","6C:5E:43:B2:26:7F:21:12:CA:4A:02:9C:D4:FF:C4:B0:93:F6:FE:E4", "administrator@vsphere.local",
		"Admin!23", "60ad0bda-ff16-492f-8b81-8aff917872c9", "datastore-31", "", "", "vm1", "", gDiskLib.VIXDISKLIB_FLAG_OPEN_COMPRESSION_SKIPZ,
		false, gDiskLib.NBD)


	//diskReaderWriter, dli, err := Open_test(params, logrus.New())
	_, dli, err := Open_test(params, logrus.New())
	if err != nil {
		gDiskLib.EndAccess(params)
		fmt.Errorf("Open failed, got error code: %d, error message: %s.", err.VixErrorCode(), err.Error())
		return
	}

	QueryBlocks(dli, params)
}

func QueryBlocks(diskHandle gvddk_high.DiskConnectHandle, params gDiskLib.ConnectParams) {
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
		abList, err := diskHandle.QueryAllocatedBlocks(gDiskLib.VixDiskLibSectorType(offset), gDiskLib.VixDiskLibSectorType(numChunkToQuery)*gDiskLib.VixDiskLibSectorType(chunkSize), gDiskLib.VixDiskLibSectorType(chunkSize))
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

func Open_test (globalParams gDiskLib.ConnectParams, logger logrus.FieldLogger) (gvddk_high.DiskReaderWriter, gvddk_high.DiskConnectHandle, gDiskLib.VddkError) {
	err := gDiskLib.PrepareForAccess(globalParams)
	if err != nil {
		return gvddk_high.DiskReaderWriter{}, gvddk_high.DiskConnectHandle{}, err
	}
	conn, err := gDiskLib.ConnectEx(globalParams)
	if err != nil {
		gDiskLib.EndAccess(globalParams)
		return gvddk_high.DiskReaderWriter{}, gvddk_high.DiskConnectHandle{}, err
	}
	dli, err := gDiskLib.Open(conn, globalParams)
	if err != nil {
		gDiskLib.Disconnect(conn)
		gDiskLib.EndAccess(globalParams)
		return gvddk_high.DiskReaderWriter{}, gvddk_high.DiskConnectHandle{}, err
	}
	info, err := gDiskLib.GetInfo(dli)
	if err != nil {
		gDiskLib.Disconnect(conn)
		gDiskLib.EndAccess(globalParams)
		return gvddk_high.DiskReaderWriter{}, gvddk_high.DiskConnectHandle{}, err
	}
	diskHandle := gvddk_high.NewDiskHandle(dli, conn, globalParams, info)
	return gvddk_high.NewDiskReaderWriter(diskHandle, logger), diskHandle, nil
}