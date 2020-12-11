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
	fmt.Println("Test Open")
	var majorVersion uint32 = 7
	var minorVersion uint32 = 0
	path := os.Getenv("LIBPATH")
	if path == "" {
		t.Skip("Skipping testing if environment variables are not set.")
	}
	gDiskLib.Init(majorVersion, minorVersion, path)
	serverName := os.Getenv("IP")
	thumPrint := os.Getenv("THUMBPRINT")
	//thumPrint = GetThumbPrintForServer(serverName)
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
	// QAB (assume at least 1GiB volume and 1MiB block size)
	abInitial, err := diskReaderWriter.QueryAllocatedBlocks(0, 2048*1024, 2048)
	if err != nil {
		t.Errorf("QueryAllocatedBlocks failed: %d, error message: %s", err.VixErrorCode(), err.Error())
	} else {
		fmt.Printf("Number of blocks: %d\n", len(abInitial))
		fmt.Printf("Offset      Length\n")
		for _, ab := range abInitial {
			fmt.Printf("0x%012x  0x%012x\n", ab.Offset(), ab.Length())
		}
	}
}
