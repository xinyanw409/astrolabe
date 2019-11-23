package main

import (
	"fmt"
	"testing"
	"github.com/sirupsen/logrus"
	"github.com/vmware/gvddk/gDiskLib"
	"github.com/vmware/gvddk/gvddk-high"
)

// II vs II
func TestAligned(t *testing.T) {
	//t.Parallel()
	fmt.Println("Test Multithread write for aligned case which skip lock: II vs II")
	var majorVersion uint32 = 6
	var minorVersion uint32 = 7
	var path string = "/usr/lib/vmware-vix-disklib"
	gDiskLib.Init(majorVersion, minorVersion, path)
	params := gDiskLib.NewConnectParams("", "10.161.131.94","D7:3E:C5:99:ED:AA:74:18:B4:08:1E:40:1C:B8:D2:10:68:02:84:4F", "administrator@vsphere.local",
		"Admin!23", "ad39188b-782c-4b00-a4fb-7785378da976", "datastore-58", "", "", "vm-example", "", gDiskLib.VIXDISKLIB_FLAG_OPEN_COMPRESSION_SKIPZ,
		false, "nbd")
	//var logger logrus.FieldLogger
	diskReaderWriter, err := gvddk_high.Open(params, logrus.New())
	if err != nil {
		gDiskLib.EndAccess(params)
		t.Errorf("Open failed, got error code: %d, error message: %s.", err.VixErrorCode(), err.Error())
	}
	// WriteAt
	done := make(chan bool)
	fmt.Println("---------------------WriteAt start----------------------")
	go func() {
		buf1 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE)
		for i, _ := range (buf1) {
			buf1[i] = 'A'
		}
		n2, err2 := diskReaderWriter.WriteAt(buf1, 0)
		fmt.Printf("--------Write A byte n = %d\n", n2)
		fmt.Println(err2)
		done <- true
	}()

	go func() {
		buf1 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE)
		for i, _ := range (buf1) {
			buf1[i] = 'B'
		}
		n2, err2 := diskReaderWriter.WriteAt(buf1, 0)
		fmt.Printf("--------Write B byte n = %d\n", n2)
		fmt.Println(err2)
		done <- true
	}()

	for i := 0; i < 2; i++ {
		<-done
	}
	// Verify written data by read
	fmt.Println("----------Read start to verify----------")
	buffer2 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE)
	n2, err5 := diskReaderWriter.ReadAt(buffer2, 0)
	fmt.Printf("Read byte n = %d\n", n2)
	fmt.Println(buffer2)
	fmt.Println(err5)

	diskReaderWriter.Close()
}

// I II III vs II III
func TestMiss1(t *testing.T) {
	//t.Parallel()
	fmt.Println("Test Multithread write for miss aligned case which lock: I II III vs II III")
	var majorVersion uint32 = 6
	var minorVersion uint32 = 7
	var path string = "/usr/lib/vmware-vix-disklib"
	gDiskLib.Init(majorVersion, minorVersion, path)
	params := gDiskLib.NewConnectParams("", "10.161.131.94","D7:3E:C5:99:ED:AA:74:18:B4:08:1E:40:1C:B8:D2:10:68:02:84:4F", "administrator@vsphere.local",
		"Admin!23", "ad39188b-782c-4b00-a4fb-7785378da976", "datastore-58", "", "", "vm-example", "", gDiskLib.VIXDISKLIB_FLAG_OPEN_COMPRESSION_SKIPZ,
		false, "nbd")
	//var logger logrus.FieldLogger
	diskReaderWriter, err := gvddk_high.Open(params, logrus.New())
	if err != nil {
		gDiskLib.EndAccess(params)
		t.Errorf("Open failed, got error code: %d, error message: %s.", err.VixErrorCode(), err.Error())
	}
	// WriteAt
	done := make(chan bool)
	fmt.Println("---------------------WriteAt start----------------------")
	go func() {
		buf1 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE + 14)
		for i, _ := range (buf1) {
			buf1[i] = 'C'
		}
		n2, err2 := diskReaderWriter.WriteAt(buf1, 500)
		fmt.Printf("--------Write C byte n = %d\n", n2)
		fmt.Println(err2)
		done <- true
	}()

	go func() {
		buf1 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE + 2)
		for i, _ := range (buf1) {
			buf1[i] = 'D'
		}
		n2, err2 := diskReaderWriter.WriteAt(buf1, gDiskLib.VIXDISKLIB_SECTOR_SIZE)
		fmt.Printf("--------Write D byte n = %d\n", n2)
		fmt.Println(err2)
		done <- true
	}()

	for i := 0; i < 2; i++ {
		<-done
	}
	// Verify written data by read
	fmt.Println("----------Read start to verify----------")
	buffer2 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE + 14)
	n2, err5 := diskReaderWriter.ReadAt(buffer2, 500)
	fmt.Printf("Read byte n = %d\n", n2)
	fmt.Println(buffer2)
	fmt.Println(err5)

	diskReaderWriter.Close()
}

// I II vs I II III
func TestMiss2(t *testing.T) {
	//t.Parallel()
	fmt.Println("Test Multithread write for miss aligned case which lock: I II vs I II III")
	var majorVersion uint32 = 6
	var minorVersion uint32 = 7
	var path string = "/usr/lib/vmware-vix-disklib"
	gDiskLib.Init(majorVersion, minorVersion, path)
	params := gDiskLib.NewConnectParams("", "10.161.131.94","D7:3E:C5:99:ED:AA:74:18:B4:08:1E:40:1C:B8:D2:10:68:02:84:4F", "administrator@vsphere.local",
		"Admin!23", "ad39188b-782c-4b00-a4fb-7785378da976", "datastore-58", "", "", "vm-example", "", gDiskLib.VIXDISKLIB_FLAG_OPEN_COMPRESSION_SKIPZ,
		false, "nbd")
	//var logger logrus.FieldLogger
	diskReaderWriter, err := gvddk_high.Open(params, logrus.New())
	if err != nil {
		gDiskLib.EndAccess(params)
		t.Errorf("Open failed, got error code: %d, error message: %s.", err.VixErrorCode(), err.Error())
	}
	// WriteAt
	done := make(chan bool)
	fmt.Println("---------------------WriteAt start----------------------")
	go func() {
		buf1 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE + 12)
		for i, _ := range (buf1) {
			buf1[i] = 'E'
		}
		n2, err2 := diskReaderWriter.WriteAt(buf1, 500)
		fmt.Printf("--------Write E byte n = %d\n", n2)
		fmt.Println(err2)
		done <- true
	}()

	go func() {
		buf1 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE + 14)
		for i, _ := range (buf1) {
			buf1[i] = 'F'
		}
		n2, err2 := diskReaderWriter.WriteAt(buf1, 500)
		fmt.Printf("--------Write F byte n = %d\n", n2)
		fmt.Println(err2)
		done <- true
	}()

	for i := 0; i < 2; i++ {
		<-done
	}
	// Verify written data by read
	fmt.Println("----------Read start to verify----------")
	buffer2 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE + 14)
	n2, err5 := diskReaderWriter.ReadAt(buffer2, 500)
	fmt.Printf("Read byte n = %d\n", n2)
	fmt.Println(buffer2)
	fmt.Println(err5)

	diskReaderWriter.Close()
}

// I II vs II III
func TestMiss3(t *testing.T) {
	//t.Parallel()
	fmt.Println("Test Multithread write for miss aligned case which lock: I II vs II III")
	var majorVersion uint32 = 6
	var minorVersion uint32 = 7
	var path string = "/usr/lib/vmware-vix-disklib"
	gDiskLib.Init(majorVersion, minorVersion, path)
	params := gDiskLib.NewConnectParams("", "10.161.131.94","D7:3E:C5:99:ED:AA:74:18:B4:08:1E:40:1C:B8:D2:10:68:02:84:4F", "administrator@vsphere.local",
		"Admin!23", "ad39188b-782c-4b00-a4fb-7785378da976", "datastore-58", "", "", "vm-example", "", gDiskLib.VIXDISKLIB_FLAG_OPEN_COMPRESSION_SKIPZ,
		false, "nbd")
	//var logger logrus.FieldLogger
	diskReaderWriter, err := gvddk_high.Open(params, logrus.New())
	if err != nil {
		gDiskLib.EndAccess(params)
		t.Errorf("Open failed, got error code: %d, error message: %s.", err.VixErrorCode(), err.Error())
	}
	// WriteAt
	done := make(chan bool)
	fmt.Println("---------------------WriteAt start----------------------")
	go func() {
		buf1 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE + 12)
		for i, _ := range (buf1) {
			buf1[i] = 'G'
		}
		n2, err2 := diskReaderWriter.WriteAt(buf1, 500)
		fmt.Printf("--------Write G byte n = %d\n", n2)
		fmt.Println(err2)
		done <- true
	}()

	go func() {
		buf1 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE + 2)
		for i, _ := range (buf1) {
			buf1[i] = 'H'
		}
		n2, err2 := diskReaderWriter.WriteAt(buf1, gDiskLib.VIXDISKLIB_SECTOR_SIZE)
		fmt.Printf("--------Write H byte n = %d\n", n2)
		fmt.Println(err2)
		done <- true
	}()

	for i := 0; i < 2; i++ {
		<-done
	}
	// Verify written data by read
	fmt.Println("----------Read start to verify----------")
	buffer2 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE + 14)
	n2, err5 := diskReaderWriter.ReadAt(buffer2, 500)
	fmt.Printf("Read byte n = %d\n", n2)
	fmt.Println(buffer2)
	fmt.Println(err5)

	diskReaderWriter.Close()
}

// I II III vs II
func TestMissAlign(t *testing.T) {
	//t.Parallel()
	fmt.Println("Test Multithread write for case which lock: I II III vs II")
	var majorVersion uint32 = 6
	var minorVersion uint32 = 7
	var path string = "/usr/lib/vmware-vix-disklib"
	gDiskLib.Init(majorVersion, minorVersion, path)
	params := gDiskLib.NewConnectParams("", "10.161.131.94","D7:3E:C5:99:ED:AA:74:18:B4:08:1E:40:1C:B8:D2:10:68:02:84:4F", "administrator@vsphere.local",
		"Admin!23", "ad39188b-782c-4b00-a4fb-7785378da976", "datastore-58", "", "", "vm-example", "", gDiskLib.VIXDISKLIB_FLAG_OPEN_COMPRESSION_SKIPZ,
		false, "nbd")
	//var logger logrus.FieldLogger
	diskReaderWriter, err := gvddk_high.Open(params, logrus.New())
	if err != nil {
		gDiskLib.EndAccess(params)
		t.Errorf("Open failed, got error code: %d, error message: %s.", err.VixErrorCode(), err.Error())
	}
	// WriteAt
	done := make(chan bool)
	fmt.Println("---------------------WriteAt start----------------------")
	go func() {
		buf1 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE + 14)
		for i, _ := range (buf1) {
			buf1[i] = 'A'
		}
		n2, err2 := diskReaderWriter.WriteAt(buf1, 500)
		fmt.Printf("--------Write A byte n = %d\n", n2)
		fmt.Println(err2)
		done <- true
	}()

	go func() {
		buf1 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE)
		for i, _ := range (buf1) {
			buf1[i] = 'B'
		}
		n2, err2 := diskReaderWriter.WriteAt(buf1, gDiskLib.VIXDISKLIB_SECTOR_SIZE)
		fmt.Printf("--------Write B byte n = %d\n", n2)
		fmt.Println(err2)
		done <- true
	}()

	for i := 0; i < 2; i++ {
		<-done
	}
	// Verify written data by read
	fmt.Println("----------Read start to verify----------")
	buffer2 := make([]byte, gDiskLib.VIXDISKLIB_SECTOR_SIZE + 14)
	n2, err5 := diskReaderWriter.ReadAt(buffer2, 500)
	fmt.Printf("Read byte n = %d\n", n2)
	fmt.Println(buffer2)
	fmt.Println(err5)

	diskReaderWriter.Close()
}
