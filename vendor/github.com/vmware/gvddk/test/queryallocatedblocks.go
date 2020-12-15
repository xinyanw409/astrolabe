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
	diskReaderWriter, err := gvddk_high.Open(params, logrus.New())
	if err != nil {
		gDiskLib.EndAccess(params)
		fmt.Errorf("Open failed, got error code: %d, error message: %s.", err.VixErrorCode(), err.Error())
		return
	}

	//// Call QueryAllocatedBlocks at T1
	//abBefore, err := diskReaderWriter.QueryAllocatedBlocks(0, 2 * gDiskLib.VIXDISKLIB_SECTOR_SIZE, gDiskLib.VIXDISKLIB_SECTOR_SIZE)
	//if err != nil {
	//	gDiskLib.EndAccess(params)
	//	fmt.Errorf("QueryAllocatedBlocks failed, got error code: %d, error message: %s.", err.VixErrorCode(), err.Error())
	//	return
	//}
	//fmt.Printf("Number of blocks: %d\n", len(abBefore))
	//fmt.Printf("Offset Length\n")
	//for _, ab := range abBefore {
	//	fmt.Printf("0x%012x  0x%012x\n", ab.Offset(), ab.Length())
	//}

	// Write to disk
	fmt.Println("WriteAt start")
	buf1 := make([]byte, 20 * gDiskLib.VIXDISKLIB_SECTOR_SIZE)
	for i,_ := range(buf1) {
		buf1[i] = 'E'
	}
	n, err1 := diskReaderWriter.WriteAt(buf1, 0)
	if err1 != nil {
		gDiskLib.EndAccess(params)
		fmt.Errorf("Write failed, got error %v.", err1)
		return
	}
	fmt.Printf("Write byte n = %d\n", n)

	buffer2 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE)
	n2, err5 := diskReaderWriter.ReadAt(buffer2, 0)
	fmt.Printf("Read byte n = %d\n", n2)
	fmt.Println(buffer2)
	fmt.Println(err5)

	//// Call QueryAllocatedBlocks at T2
	//abLater, err := diskReaderWriter.QueryAllocatedBlocks(0, 2 * gDiskLib.VIXDISKLIB_SECTOR_SIZE, gDiskLib.VIXDISKLIB_SECTOR_SIZE)
	//if err != nil {
	//	gDiskLib.EndAccess(params)
	//	fmt.Errorf("QueryAllocatedBlocks failed, got error code: %d, error message: %s.", err.VixErrorCode(), err.Error())
	//	return
	//}
	//fmt.Printf("Number of blocks: %d\n", len(abLater))
	//fmt.Printf("Offset      Length\n")
	//for _, ab := range abLater {
	//	fmt.Printf("0x%012x  0x%012x\n", ab.Offset(), ab.Length())
	//}

	offset := 0
	capacity := 10240
	chunkSize := 2048
	numChunk := capacity /chunkSize
	var numChunkToQuery int
	for numChunk > 0 {
		if numChunk > gDiskLib.VIXDISKLIB_MAX_CHUNK_NUMBER {
			numChunkToQuery = gDiskLib.VIXDISKLIB_MAX_CHUNK_NUMBER
		} else {
			numChunkToQuery = numChunk
		}
		abList, err := diskReaderWriter.QueryAllocatedBlocks(gDiskLib.VixDiskLibSectorType(offset), gDiskLib.VixDiskLibSectorType(numChunkToQuery) * gDiskLib.VixDiskLibSectorType(chunkSize), gDiskLib.VixDiskLibSectorType(chunkSize))
		if err != nil {
			gDiskLib.EndAccess(params)
			fmt.Errorf("QueryAllocatedBlocks failed, got error code: %d, error message: %s.", err.VixErrorCode(), err.Error())
			return
		}
		fmt.Printf("Number of blocks: %d\n", len(abList))
		fmt.Printf("Offset      Length\n")
		for _, ab := range abList {
			fmt.Printf("0x%012x  0x%012x\n", ab.Offset(), ab.Length())
		}

		numChunk = numChunk - numChunkToQuery
		offset = offset + numChunkToQuery * chunkSize
	}
}
