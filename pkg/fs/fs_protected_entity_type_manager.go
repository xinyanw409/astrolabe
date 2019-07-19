package fs

import (
	"context"
	"github.com/vmware/arachne/pkg/arachne"
)

type FSProtectedEntityTypeManager struct {
	root string
}

func NewFSProtectedEntityTypeManagerFromConfig(params map[string]interface{}) (*FSProtectedEntityTypeManager, error) {
	root := params["root"].(string)

	returnTypeManager := FSProtectedEntityTypeManager{
		root: root,
	}
	return &returnTypeManager, nil
}

func (this *FSProtectedEntityTypeManager) GetTypeName() string {
	return "fs"
}

func (this *FSProtectedEntityTypeManager) GetProtectedEntity(ctx context.Context, id arachne.ProtectedEntityID) (
	arachne.ProtectedEntity, error) {
	return nil, nil
}

func (this *FSProtectedEntityTypeManager) GetProtectedEntities(ctx context.Context) ([]arachne.ProtectedEntity, error) {
	return nil, nil
}

func (this *FSProtectedEntityTypeManager) Copy(ctx context.Context, pe arachne.ProtectedEntity) error {
	return nil
}
