package server

import (
	"context"
	"github.com/labstack/echo"
	"github.com/vmware/arachne/pkg/arachne"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type Arachne struct {
	petm *DirectProtectedEntityManager
	api_services map[string]*ServiceAPI
	s3_services map[string]*ServiceS3
	s3URLBase string
}

func NewArachne(confDirPath string, port int) *Arachne {
	api_services := make(map[string]*ServiceAPI)
	s3_services := make(map[string]*ServiceS3)
	s3URLBase, err := configS3URL(port)
	if err != nil {
		log.Fatal("Could not get host IP address", err)
	}
	petm := NewDirectProtectedEntityManagerFromConfigDir(confDirPath, s3URLBase)
	for _, curService := range petm.ListEntityTypeManagers() {
		serviceName := curService.GetTypeName()
		api_services[serviceName] = NewServiceAPI(curService)
		s3_services[serviceName] = NewServiceS3(curService)
	}

	retArachne := Arachne{
		api_services: api_services,
		s3_services: s3_services,
		s3URLBase: s3URLBase,
	}

	return &retArachne
}

func NewArachneRepository() *Arachne {
	return nil
}

func (this *Arachne) Get(c echo.Context) error {
	var servicesList strings.Builder
	needsComma := false
	for serviceName := range this.api_services {
		if needsComma {
			servicesList.WriteString(",")
		}
		servicesList.WriteString(serviceName)
		needsComma = true
	}
	return c.String(http.StatusOK, servicesList.String())
}

func (this *Arachne) ConnectArachneAPIToEcho(echo *echo.Echo) error {
	echo.GET("/arachne", this.Get)

	for serviceName, service := range this.api_services {
		echo.GET("/arachne/"+serviceName, service.listObjects)
		echo.POST("/arachne/"+serviceName, service.handleCopyObject)
		echo.GET("/arachne/"+serviceName+"/:id", service.handleObjectRequest)
		echo.GET("/arachne/"+serviceName+"/:id/snapshots", service.handleSnapshotListRequest)

	}
	return nil
}

func configS3URL(port int) (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return "http://" + ipnet.IP.String() + ":" + strconv.Itoa(port) + "/s3/", nil
			}
		}
	}
	return "", nil
}

func (this *Arachne) ConnectMiniS3ToEcho(echo *echo.Echo) error {
	echo.GET("/s3", this.Get)

	for serviceName, service := range this.s3_services {
		echo.GET("/s3/"+serviceName, service.listObjects)
		echo.GET("/s3/"+serviceName+"/:objectKey", service.handleObjectRequest)
	}
	return nil
}
const fileSuffix = ".pe.json"

func getProtectedEntityForIDStr(petm arachne.ProtectedEntityTypeManager, idStr string,
	echoContext echo.Context) (arachne.ProtectedEntityID, arachne.ProtectedEntity, error) {
	var id arachne.ProtectedEntityID
	var pe arachne.ProtectedEntity
	var err error

	id, err = arachne.NewProtectedEntityIDFromString(idStr)
	if err != nil {
		echoContext.String(http.StatusBadRequest, "id = "+idStr+" is invalid "+err.Error())
		return id, pe, err
	}
	if id.GetPeType() != (petm).GetTypeName() {
		echoContext.String(http.StatusBadRequest, "id = "+idStr+" is not type "+petm.GetTypeName())
		return id, pe, err
	}
	pe, err = (petm).GetProtectedEntity(context.Background(), id)
	if err != nil {
		echoContext.String(http.StatusNotFound, "Could not retrieve id "+id.String()+" error = "+err.Error())
		return id, pe, err
	}
	if pe == nil {
		echoContext.String(http.StatusInternalServerError, "pe was nil for "+id.String())
		return id, pe, err
	}
	return id, pe, nil
}

