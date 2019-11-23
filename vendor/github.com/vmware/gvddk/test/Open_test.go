package main

import (
	"fmt"
	"testing"
	"github.com/sirupsen/logrus"
	"github.com/vmware/gvddk/gDiskLib"
	"github.com/vmware/gvddk/gvddk-high"
)

func Test(t *testing.T) {
	fmt.Println("Test Open")
	var majorVersion uint32 = 6
	var minorVersion uint32 = 7
	var path string = "/usr/lib/vmware-vix-disklib"
	gDiskLib.Init(majorVersion, minorVersion, path)
	fmt.Println("Open")
	params := gDiskLib.NewConnectParams("", "10.161.131.94","D7:3E:C5:99:ED:AA:74:18:B4:08:1E:40:1C:B8:D2:10:68:02:84:4F", "administrator@vsphere.local",
		"Admin!23", "ad39188b-782c-4b00-a4fb-7785378da976", "datastore-58", "", "", "vm-example", "", gDiskLib.VIXDISKLIB_FLAG_OPEN_COMPRESSION_SKIPZ,
		false, "nbd")
	//var logger logrus.FieldLogger
	diskReaderWriter, err := gvddk_high.Open(params, logrus.New())
	if err != nil {
		gDiskLib.EndAccess(params)
		t.Errorf("Open failed, got error code: %d, error message: %s.", err.VixErrorCode(), err.Error())
	}
	// ReadAt
	fmt.Printf("ReadAt test\n")
	buffer := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE)
	n, err4 := diskReaderWriter.Read(buffer)
	fmt.Printf("Read byte n = %d\n", n)
	fmt.Println(buffer)
	fmt.Println(err4)

	// WriteAt
	fmt.Println("WriteAt start")
	//buf1 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE)
	//for i,_ := range(buf1) {
	//	buf1[i] = 'B'
	//}
	//n2, err2 := diskReaderWriter.Write(buf1)
	//fmt.Printf("Write byte n = %d\n", n2)
	//fmt.Println(err2)

	buf1 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE)
	for i,_ := range(buf1) {
		buf1[i] = 'E'
	}
	n2, err2 := diskReaderWriter.WriteAt(buf1, 0)
	fmt.Printf("Write byte n = %d\n", n2)
	fmt.Println(err2)

	buffer2 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE)
	n2, err5 := diskReaderWriter.ReadAt(buffer2, 0)
	fmt.Printf("Read byte n = %d\n", n2)
	fmt.Println(buffer2)
	fmt.Println(err5)

	//// ReadAt
	//fmt.Printf("ReadAt test\n")
	//buf2 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE)
	//n1, err5 := diskReaderWriter.Read(buf2)
	//fmt.Printf("Read byte n = %d\n", n1)
	//fmt.Println(buf2)
	//fmt.Println(err5)


	diskReaderWriter.Close()
}
