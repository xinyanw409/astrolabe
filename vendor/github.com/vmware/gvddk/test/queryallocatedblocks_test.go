package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vmware/gvddk/gDiskLib"
	"github.com/vmware/gvddk/gvddk-high"
	"github.com/stretchr/testify/require"
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

	// Call QueryAllocatedBlocks at T1
	abBefore, err := diskReaderWriter.QueryAllocatedBlocks(0, 2048*1024, 2048)
	if err != nil {
		gDiskLib.EndAccess(params)
		t.Errorf("QueryAllocatedBlocks failed, got error code: %d, error message: %s.", err.VixErrorCode(), err.Error())
	}
	fmt.Printf("Number of blocks: %d\n", len(abBefore))
	fmt.Printf("Offset     Length\n")
	for _, ab := range abBefore {
		fmt.Printf("0x%012x  0x%012x\n", ab.Offset(), ab.Length())
	}

	// Write to disk
	fmt.Println("WriteAt start")
	buf1 := make([]byte, 2 * gDiskLib.VIXDISKLIB_SECTOR_SIZE)
	for i,_ := range(buf1) {
		buf1[i] = 'E'
	}
	n, err1 := diskReaderWriter.WriteAt(buf1, 0)
	require.Nil(t, err1)
	fmt.Printf("Write byte n = %d\n", n)

	// Call QueryAllocatedBlocks at T2
	abLater, err := diskReaderWriter.QueryAllocatedBlocks(0, 2048*1024, 2048)
	require.Nil(t, err)
	fmt.Printf("Number of blocks: %d\n", len(abLater))
	fmt.Printf("Offset      Length\n")
	for _, ab := range abLater {
		fmt.Printf("0x%012x  0x%012x\n", ab.Offset(), ab.Length())
	}
}
