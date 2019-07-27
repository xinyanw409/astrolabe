package server

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/vmware/arachne/pkg/arachne"
	"github.com/vmware/arachne/pkg/fs"
	"github.com/vmware/arachne/pkg/ivd"
	"github.com/vmware/arachne/pkg/kubernetes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Arachne struct {
	api_services map[string]*ServiceAPI
	s3_services map[string]*ServiceS3
	s3URLBase string
}

func NewArachne(confDirPath string, port int) *Arachne {
	api_services := make(map[string]*ServiceAPI)
	s3_services := make(map[string]*ServiceS3)
	configMap, err := readConfigFiles(confDirPath)
	if err != nil {
		log.Fatal("Could not read config files", err)
	}
	s3URLBase, err := configS3URL(port)
	if err != nil {
		log.Fatal("Could not get host IP address", err)
	}
	for serviceName, params := range configMap {
		var curService arachne.ProtectedEntityTypeManager
		switch serviceName {
		case "ivd":
			curService, err = ivd.NewIVDProtectedEntityTypeManagerFromConfig(params, s3URLBase)
		case "k8sns":
			curService, err = kubernetes.NewKubernetesNamespaceProtectedEntityTypeManagerFromConfig(params, s3URLBase)
		case "fs":
			curService, err = fs.NewFSProtectedEntityTypeManagerFromConfig(params, s3URLBase)
		default:

		}
		if err != nil {
			log.Printf("Could not start service %s err=%v", serviceName, err)
			continue
		}
		if curService != nil {
			api_services[serviceName] = NewServiceAPI(&curService)
			s3_services[serviceName] = NewServiceS3(&curService)
		}
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

func readConfigFiles(confDirPath string) (map[string]map[string]interface{}, error) {
	configMap := make(map[string]map[string]interface{})

	confDir, err := os.Stat(confDirPath)
	if err != nil {
		log.Panicln("Could not stat configuration directory " + confDirPath)
	}
	if !confDir.Mode().IsDir() {
		log.Panicln(confDirPath + " is not a directory")

	}

	files, err := ioutil.ReadDir(confDirPath)
	for _, curFile := range files {
		if !strings.HasPrefix(curFile.Name(), ".") && strings.HasSuffix(curFile.Name(), fileSuffix) {
			peTypeName := strings.TrimSuffix(curFile.Name(), fileSuffix)
			peConf, err := readConfigFile(filepath.Join(confDirPath, curFile.Name()))
			if err != nil {
				log.Panicln("Could not process conf file " + curFile.Name() + " continuing, err = " + err.Error())
			} else {
				configMap[peTypeName] = peConf
			}
		}
	}
	return configMap, nil
}
func readConfigFile(confFile string) (map[string]interface{}, error) {
	jsonFile, err := os.Open(confFile)
	if err != nil {
		return nil, errors.Wrap(err, "Could not open conf file "+confFile)
	}
	defer jsonFile.Close()
	jsonBytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, errors.Wrap(err, "Could not read conf file "+confFile)
	}
	var result map[string]interface{}
	err = json.Unmarshal([]byte(jsonBytes), &result)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal JSON from "+confFile)
	}
	return result, nil
}

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

