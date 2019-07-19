package fs

import (
	"github.com/vmware/arachne/pkg/arachne"
	vim "github.com/vmware/govmomi/vim25/types"
	//	"github.com/vmware/govmomi/vslm"
	"context"
	"net/url"
)

type FSProtectedEntity struct {
	fspetm *FSProtectedEntityTypeManager
	id     arachne.ProtectedEntityID
	name   string
}

func newProtectedEntityID(id vim.ID) arachne.ProtectedEntityID {
	return arachne.NewProtectedEntityID("fs", id.Id)
}

func newFSProtectedEntity(fspetm *FSProtectedEntityTypeManager, id arachne.ProtectedEntityID) (FSProtectedEntity, error) {
	newFSPE := FSProtectedEntity{
		fspetm: fspetm,
		id:     id,
	}
	return newFSPE, nil
}
func (this FSProtectedEntity) GetInfo(ctx context.Context) (arachne.ProtectedEntityInfo, error) {
	retVal := arachne.ProtectedEntityInfoImpl{
		Id:           this.id,
		Name:         this.name,
		CombinedURLs: []url.URL{},
		DataURLs:     []url.URL{},
		MetadataURLs: []url.URL{},
		ComponentIDs: []arachne.ProtectedEntityID{},
	}
	return retVal, nil
}

func (this FSProtectedEntity) GetCombinedInfo(ctx context.Context) ([]arachne.ProtectedEntityInfo, error) {
	fsIPE, err := this.GetInfo(ctx)
	if err != nil {
		return nil, err
	}
	return []arachne.ProtectedEntityInfo{fsIPE}, nil
}

/*
 * Snapshot APIs
 */
func (this FSProtectedEntity) Snapshot(ctx context.Context) (*arachne.ProtectedEntitySnapshotID, error) {
	return nil, nil
}

func (this FSProtectedEntity) ListSnapshots(ctx context.Context) ([]arachne.ProtectedEntitySnapshotID, error) {
	return nil, nil
}
func (this FSProtectedEntity) DeleteSnapshot(ctx context.Context, snapshotToDelete arachne.ProtectedEntitySnapshotID) (bool, error) {
	return true, nil
}
func (this FSProtectedEntity) GetInfoForSnapshot(ctx context.Context, snapshotID arachne.ProtectedEntitySnapshotID) (*arachne.ProtectedEntityInfo, error) {
	return nil, nil
}

func (this FSProtectedEntity) GetComponents(ctx context.Context) ([]arachne.ProtectedEntity, error) {
	return make([]arachne.ProtectedEntity, 0), nil
}

func (this FSProtectedEntity) GetID() arachne.ProtectedEntityID {
	return this.id
}

func NewIDFromString(idStr string) vim.ID {
	return vim.ID{
		Id: idStr,
	}
}

func NewVimIDFromPEID(peid arachne.ProtectedEntityID) vim.ID {
	return vim.ID{
		Id: peid.GetID(),
	}
}
