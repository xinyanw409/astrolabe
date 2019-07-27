package server

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/vmware/arachne/pkg/arachne"
	"github.com/vmware/arachne/pkg/fs"
	"github.com/vmware/arachne/pkg/ivd"
	"github.com/vmware/arachne/pkg/kubernetes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type DirectProtectedEntityManager struct {
	typeManager map[string]arachne.ProtectedEntityTypeManager
	s3URLBase   string
}

func NewDirectProtectedEntityManager(petms []arachne.ProtectedEntityTypeManager, s3URLBase string) (returnPEM *DirectProtectedEntityManager) {
	returnPEM = &DirectProtectedEntityManager{
		typeManager: make(map[string]arachne.ProtectedEntityTypeManager),
	}
	for _, curPETM := range petms {
		returnPEM.typeManager[curPETM.GetTypeName()] = curPETM
	}
	returnPEM.s3URLBase = s3URLBase
	return
}

func NewDirectProtectedEntityManagerFromConfigDir(confDirPath string, s3URLBase string) (*DirectProtectedEntityManager) {
	configMap, err := readConfigFiles(confDirPath)
	if err != nil {
		log.Fatal("Could not read config files", err)
	}
	return NewDirectProtectedEntityManagerFromParamMap(configMap, s3URLBase)
}

func NewDirectProtectedEntityManagerFromParamMap(configMap map[string]map[string]interface{}, s3URLBase string) (
*DirectProtectedEntityManager) {
	petms := make([]arachne.ProtectedEntityTypeManager, 0) // No guarantee all configs will be valid, so don't preallocate
	var err error

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
			petms = append(petms, curService)
		}
	}
	return NewDirectProtectedEntityManager(petms, s3URLBase)
}

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

func (this *DirectProtectedEntityManager) GetProtectedEntity(ctx context.Context, id arachne.ProtectedEntityID) (arachne.ProtectedEntity, error) {
	return this.typeManager[id.GetPeType()].GetProtectedEntity(ctx, id);
}

func (this *DirectProtectedEntityManager) GetProtectedEntityTypeManager(peType string) arachne.ProtectedEntityTypeManager {
	return this.typeManager[peType]
}

func (this *DirectProtectedEntityManager) ListEntityTypeManagers() []arachne.ProtectedEntityTypeManager {
	returnArr := []arachne.ProtectedEntityTypeManager{}
	for _, curPETM := range this.typeManager {
		returnArr = append(returnArr, curPETM)
	}
	return returnArr
}
