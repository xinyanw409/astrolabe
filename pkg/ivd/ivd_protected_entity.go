package ivd

import (
	"github.com/vmware/arachne/pkg/core"
	vim "github.com/vmware/govmomi/vim25/types"
	//	"github.com/vmware/govmomi/vslm"
	"context"
	"github.com/pkg/errors"
	"net/url"
	"time"
)

type IVDProtectedEntity struct {
	ipetm *IVDProtectedEntityTypeManager
	id    core.ProtectedEntityID
}

func newProtectedEntityID(id vim.ID) core.ProtectedEntityID {
	return core.NewProtectedEntityID("ivd", id.Id)
}

func newProtectedEntityIDWithSnapshotID(id vim.ID, snapshotID core.ProtectedEntitySnapshotID) core.ProtectedEntityID {
	return core.NewProtectedEntityIDWithSnapshotID("ivd", id.Id, snapshotID)
}

func newIVDProtectedEntity(ipetm *IVDProtectedEntityTypeManager, id core.ProtectedEntityID) (IVDProtectedEntity, error) {
	newIPE := IVDProtectedEntity{
		ipetm: ipetm,
		id:    id,
	}
	return newIPE, nil
}
func (this IVDProtectedEntity) GetInfo(ctx context.Context) (core.ProtectedEntityInfo, error) {
	vsoID := vim.ID{
		Id: this.id.GetID(),
	}
	vso, err := this.ipetm.vsom.Retrieve(ctx, vsoID)
	if err != nil {
		return nil, errors.Wrap(err, "Retrieve failed")
	}
	retVal := core.ProtectedEntityInfoImpl{
		Id:           this.id,
		Name:         vso.Config.Name,
		CombinedURLs: []url.URL{},
		DataURLs:     []url.URL{},
		MetadataURLs: []url.URL{},
		ComponentIDs: []core.ProtectedEntityID{},
	}
	return retVal, nil
}

func (this IVDProtectedEntity) GetCombinedInfo(ctx context.Context) ([]core.ProtectedEntityInfo, error) {
	ivdIPE, err := this.GetInfo(ctx)
	if err != nil {
		return nil, err
	}
	return []core.ProtectedEntityInfo{ivdIPE}, nil
}

/*
 * Snapshot APIs
 */
func (this IVDProtectedEntity) Snapshot(ctx context.Context) (*core.ProtectedEntitySnapshotID, error) {
	vslmTask, err := this.ipetm.vsom.CreateSnapshot(ctx, NewVimIDFromPEID(this.GetID()), "ArachneSnapshot")
	if err != nil {
		return nil, errors.Wrap(err, "Snapshot failed")
	}
	ivdSnapshotIDAny, err := vslmTask.Wait(ctx, 60*time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "Wait failed")
	}
	ivdSnapshotID := ivdSnapshotIDAny.(vim.ID)
	return core.NewProtectedEntitySnapshotID(ivdSnapshotID.Id), nil
}

func (this IVDProtectedEntity) ListSnapshots(ctx context.Context) ([]core.ProtectedEntitySnapshotID, error) {
	snapshotInfo, err := this.ipetm.vsom.RetrieveSnapshotInfo(ctx, NewVimIDFromPEID(this.GetID()))
	if err != nil {
		return nil, errors.Wrap(err, "RetrieveSnapshotInfo failed")
	}
	peSnapshotIDs := []core.ProtectedEntitySnapshotID{}
	for _, curSnapshotInfo := range snapshotInfo {
		peSnapshotIDs = append(peSnapshotIDs, *core.NewProtectedEntitySnapshotID(curSnapshotInfo.Id.Id))
	}
	return peSnapshotIDs, nil
}
func (this IVDProtectedEntity) DeleteSnapshot(ctx context.Context, snapshotToDelete core.ProtectedEntitySnapshotID) (bool, error) {
	return true, nil
}
func (this IVDProtectedEntity) GetInfoForSnapshot(ctx context.Context, snapshotID core.ProtectedEntitySnapshotID) (*core.ProtectedEntityInfo, error) {
	return nil, nil
}

func (this IVDProtectedEntity) GetComponents(ctx context.Context) ([]core.ProtectedEntity, error) {
	return make([]core.ProtectedEntity, 0), nil
}

func (this IVDProtectedEntity) GetID() core.ProtectedEntityID {
	return this.id
}

func NewIDFromString(idStr string) vim.ID {
	return vim.ID{
		Id: idStr,
	}
}

func NewVimIDFromPEID(peid core.ProtectedEntityID) vim.ID {
	return vim.ID{
		Id: peid.GetID(),
	}
}
