package ivd

import (
	"github.com/vmware/arachne"
	vim "github.com/vmware/govmomi/vim25/types"
	//	"github.com/vmware/govmomi/vslm"
	"context"
	"time"
)

type IVDProtectedEntity struct {
	ipetm *IVDProtectedEntityTypeManager
	id    arachne.ProtectedEntityID
	info  arachne.ProtectedEntityInfo
}

func newProtectedEntityID(id vim.ID) arachne.ProtectedEntityID {
	return arachne.NewProtectedEntityID("ivd", id.Id)
}

func newProtectedEntityIDWithSnapshotID(id vim.ID, snapshotID arachne.ProtectedEntitySnapshotID) arachne.ProtectedEntityID {
	return arachne.NewProtectedEntityIDWithSnapshotID("ivd", id.Id, snapshotID)
}

func newIVDProtectedEntity(ipetm *IVDProtectedEntityTypeManager, id arachne.ProtectedEntityID) (IVDProtectedEntity, error) {
	newIPE := IVDProtectedEntity{
		ipetm: ipetm,
		id:    id,
	}
	return newIPE, nil
}
func (ipe *IVDProtectedEntity) GetInfo() arachne.ProtectedEntityInfo {
	return ipe.info
}
func (ipe *IVDProtectedEntity) GetCombinedInfo(ctx context.Context) ([]arachne.ProtectedEntityInfo, error) {
	return make([]arachne.ProtectedEntityInfo, 0), nil
}

/*
 * Snapshot APIs
 */
func (ipe *IVDProtectedEntity) Snapshot(ctx context.Context) (*arachne.ProtectedEntitySnapshotID, error) {
	vslmTask, err := ipe.ipetm.vsom.CreateSnapshot(ctx, NewIDFromString(ipe.GetID().GetID()), "ArachneSnapshot")
	if err != nil {
		return nil, err
	}
	ivdSnapshotIDAny, err := vslmTask.Wait(ctx, 60*time.Second)
	if err != nil {
		return nil, err
	}
	ivdSnapshotID := ivdSnapshotIDAny.(arachne.ProtectedEntitySnapshotID)
	return arachne.NewProtectedEntitySnapshotID(ivdSnapshotID.String()), nil
}

func (ipe *IVDProtectedEntity) ListSnapshots(ctx context.Context) ([]arachne.ProtectedEntitySnapshotID, error) {
	return make([]arachne.ProtectedEntitySnapshotID, 0), nil
}
func (ipe *IVDProtectedEntity) DeleteSnapshot(ctx context.Context, snapshotToDelete arachne.ProtectedEntitySnapshotID) (bool, error) {
	return true, nil
}
func (ipe *IVDProtectedEntity) GetInfoForSnapshot(ctx context.Context, snapshotID arachne.ProtectedEntitySnapshotID) (*arachne.ProtectedEntityInfo, error) {
	return nil, nil
}

func (ipe *IVDProtectedEntity) GetComponents(ctx context.Context) []arachne.ProtectedEntity {
	return make([]arachne.ProtectedEntity, 0)
}

func (ipe *IVDProtectedEntity) GetID() arachne.ProtectedEntityID {
	return ipe.id
}

func NewIDFromString(idStr string) vim.ID {
	return vim.ID{
		Id: idStr,
	}
}
