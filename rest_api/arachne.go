package rest_api

import (
	"net/http"
	"github.com/labstack/echo"
	"strings"
	"context"
	"log"
	"github.com/vmware/arachne"
	"github.com/vmware/arachne/ivd"
	"net/url"
)

type Arachne struct {
	services map[string]*ServiceAPI
}

func NewArachne() (*Arachne) {
	services := make(map[string]*ServiceAPI)
	ivdService, err := InitIVDService(context.Background())
	if (err != nil) {
		log.Fatal("Could not initialize IVD service", err)
	}
	services["ivd"] = NewServiceAPI(&ivdService)
	retArachne := Arachne {
		services: services,
	}
	
	return &retArachne
}

func (this *Arachne) Get(c echo.Context) error {
	var servicesList strings.Builder
	needsComma := false
	for serviceName, _ := range this.services {
		if (needsComma) {
			servicesList.WriteString(",")
		}
		servicesList.WriteString(serviceName)
		needsComma = true
	}
	return c.String(http.StatusOK, servicesList.String())
}

func (this *Arachne) ConnectArachneAPIToEcho(echo *echo.Echo) error {
	echo.GET("/api/arachne", this.Get)

	for serviceName, service := range this.services {
		echo.GET("/api/arachne/" + serviceName, service.listObjects)
	}
	return nil
}

func InitIVDService(ctx context.Context) (arachne.ProtectedEntityTypeManager, error) {
	var vcUrl url.URL
	vcUrl.Scheme = "https"
	vcUrl.Host = "10.160.127.39"
	vcUrl.User = url.UserPassword("administrator@vsphere.local", "Admin!23")
	vcUrl.Path = "/sdk"
	
	ivdPETM, err := ivd.NewIVDProtectedEntityTypeManagerFromURL(&vcUrl, true)
	if err != nil {
		return nil, err
	}
	return ivdPETM, nil
}