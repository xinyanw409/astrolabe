package fs

import (
	"context"
	"github.com/vmware/arachne/pkg/arachne"
	"io/ioutil"
	"path/filepath"
)

type FSProtectedEntityTypeManager struct {
	root      string
	s3URLBase string
}

func NewFSProtectedEntityTypeManagerFromConfig(params map[string]interface{}, s3URLBase string) (*FSProtectedEntityTypeManager, error) {
	root := params["root"].(string)

	returnTypeManager := FSProtectedEntityTypeManager{
		root:      root,
		s3URLBase: s3URLBase,
	}
	return &returnTypeManager, nil
}

func (this *FSProtectedEntityTypeManager) GetTypeName() string {
	return "fs"
}

func (this *FSProtectedEntityTypeManager) GetProtectedEntity(ctx context.Context, id arachne.ProtectedEntityID) (
	arachne.ProtectedEntity, error) {
	return newFSProtectedEntity(this, id, id.GetID(), filepath.Join(this.root, id.GetID()))
}

func (this *FSProtectedEntityTypeManager) GetProtectedEntities(ctx context.Context) ([]arachne.ProtectedEntity, error) {
	files, err := ioutil.ReadDir(this.root)
	if err != nil {
		return nil, err
	}

	var retVal = make([]arachne.ProtectedEntity, len(files))
	for index, curFile := range files {
		peid := arachne.NewProtectedEntityID("fs", curFile.Name())
		retVal[index], err = newFSProtectedEntity(this, peid, curFile.Name(), filepath.Join(this.root, curFile.Name()))
		if err != nil {
			return nil, err
		}
	}
	return retVal, nil
}

func (this *FSProtectedEntityTypeManager) Copy(ctx context.Context, pe arachne.ProtectedEntity) error {
	return nil
}
