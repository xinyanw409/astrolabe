package main

import (
	"github.com/vmware/gvddk/gDiskLib"
	"strings"
	"testing"
)

func TestGetThumbPrintForServer(t *testing.T) {
	host := "www.vmware.com";
	port := "443";
	vmwareThumbprint := "E8:F5:BC:3E:11:6B:C1:80:4D:17:E7:45:D1:47:8E:0E:4C:DF:98:F7";
	thumbprint, err := gDiskLib.GetThumbPrintForServer(host, port);
	if err != nil {
		t.Errorf("Thumbprint for %s:%s failed, err = %s\n", host, port, err);
	}
	t.Logf("Thumbprint for %s:%s is %s\n", host, port, thumbprint);
	if strings.Compare(vmwareThumbprint, thumbprint) != 0 {
		t.Errorf("Thumbprint %s does not match expected thumbprint %s for %s - check to see if cert has been updated at %s\n",
			thumbprint, vmwareThumbprint, host, host);
	}
}
