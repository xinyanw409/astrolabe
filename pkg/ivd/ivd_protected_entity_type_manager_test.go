package ivd

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/url"
	"testing"
)

func TestProtectedEntityTypeManager(t *testing.T) {
	var vcUrl url.URL
	vcUrl.Scheme = "https"
	vcUrl.Host = "10.160.127.39"
	vcUrl.User = url.UserPassword("administrator@vsphere.local", "Admin!23")
	vcUrl.Path = "/sdk"

	t.Logf("%s\n", vcUrl.String())

	ivdPETM, err := NewIVDProtectedEntityTypeManagerFromURL(&vcUrl, "/ivd", true, logrus.New())
	ctx := context.Background()

	pes, err := ivdPETM.GetProtectedEntities(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("# of PEs returned = %d\n", len(pes))
}
