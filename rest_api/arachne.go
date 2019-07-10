package rest_api

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/vmware/arachne"
	"github.com/vmware/arachne/ivd"
	"github.com/vmware/arachne/kubernetes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Arachne struct {
	services map[string]*ServiceAPI
}

func NewArachne(confDirPath string) *Arachne {
	services := make(map[string]*ServiceAPI)

	configMap, err := readConfigFiles(confDirPath)
	if err != nil {
		log.Fatal("Could not read config files", err)
	}

	for serviceName, params := range configMap {
		var curService arachne.ProtectedEntityTypeManager
		switch serviceName {
		case "ivd":
			curService, err = ivd.NewIVDProtectedEntityTypeManagerFromConfig(params)
		case "k8sns":
			curService, err = kubernetes.NewKubernetesNamespaceProtectedEntityTypeManagerFromConfig(params)
		default:

		}
		if err != nil {
			log.Printf("Could not start service %s err=%v", serviceName, err)
			continue
		}
		if curService != nil {
			services[serviceName] = NewServiceAPI(&curService)
		}
	}
	retArachne := Arachne{
		services: services,
	}

	return &retArachne
}

func (this *Arachne) Get(c echo.Context) error {
	var servicesList strings.Builder
	needsComma := false
	for serviceName := range this.services {
		if needsComma {
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
		echo.GET("/api/arachne/"+serviceName, service.listObjects)
		echo.GET("/api/arachne/"+serviceName+"/:id", service.handleObjectRequest)
		echo.GET("/api/arachne/"+serviceName+"/:id/snapshots", service.handleSnapshotListRequest)

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
