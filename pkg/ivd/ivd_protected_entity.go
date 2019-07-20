package ivd

import (
	"github.com/pkg/errors"
	"github.com/vmware/arachne/pkg/arachne"
	vim "github.com/vmware/govmomi/vim25/types"
	"io"

	//	"github.com/vmware/govmomi/vslm"
	"context"
	"time"
)

type IVDProtectedEntity struct {
	ipetm *IVDProtectedEntityTypeManager
	id    arachne.ProtectedEntityID
}

func (this IVDProtectedEntity) GetDataReader() (io.Reader, error) {
	return nil, nil
}

func (this IVDProtectedEntity) GetMetadataReader() (io.Reader, error) {
	return nil, nil
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
func (this IVDProtectedEntity) GetInfo(ctx context.Context) (arachne.ProtectedEntityInfo, error) {
	vsoID := vim.ID{
		Id: this.id.GetID(),
	}
	vso, err := this.ipetm.vsom.Retrieve(ctx, vsoID)
	if err != nil {
		return nil, errors.Wrap(err, "Retrieve failed")
	}

	retVal := arachne.NewProtectedEntityInfo(
		this.id,
		vso.Config.Name,
		[]arachne.DataTransport{},
		[]arachne.DataTransport{},
		[]arachne.DataTransport{},
		[]arachne.ProtectedEntityID{})
	return retVal, nil
}

func (this IVDProtectedEntity) GetCombinedInfo(ctx context.Context) ([]arachne.ProtectedEntityInfo, error) {
	ivdIPE, err := this.GetInfo(ctx)
	if err != nil {
		return nil, err
	}
	return []arachne.ProtectedEntityInfo{ivdIPE}, nil
}

/*
 * Snapshot APIs
 */
func (this IVDProtectedEntity) Snapshot(ctx context.Context) (*arachne.ProtectedEntitySnapshotID, error) {
	vslmTask, err := this.ipetm.vsom.CreateSnapshot(ctx, NewVimIDFromPEID(this.GetID()), "ArachneSnapshot")
	if err != nil {
		return nil, errors.Wrap(err, "Snapshot failed")
	}
	ivdSnapshotIDAny, err := vslmTask.Wait(ctx, 60*time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "Wait failed")
	}
	ivdSnapshotID := ivdSnapshotIDAny.(vim.ID)
	/*
		ivdSnapshotStr := ivdSnapshotIDAny.(string)
		ivdSnapshotID := vim.ID{
			id: ivdSnapshotStr,
		}
	*/
	return arachne.NewProtectedEntitySnapshotID(ivdSnapshotID.Id), nil
}

func (this IVDProtectedEntity) ListSnapshots(ctx context.Context) ([]arachne.ProtectedEntitySnapshotID, error) {
	snapshotInfo, err := this.ipetm.vsom.RetrieveSnapshotInfo(ctx, NewVimIDFromPEID(this.GetID()))
	if err != nil {
		return nil, errors.Wrap(err, "RetrieveSnapshotInfo failed")
	}
	peSnapshotIDs := []arachne.ProtectedEntitySnapshotID{}
	for _, curSnapshotInfo := range snapshotInfo {
		peSnapshotIDs = append(peSnapshotIDs, *arachne.NewProtectedEntitySnapshotID(curSnapshotInfo.Id.Id))
	}
	return peSnapshotIDs, nil
}
func (this IVDProtectedEntity) DeleteSnapshot(ctx context.Context, snapshotToDelete arachne.ProtectedEntitySnapshotID) (bool, error) {
	return true, nil
}
func (this IVDProtectedEntity) GetInfoForSnapshot(ctx context.Context, snapshotID arachne.ProtectedEntitySnapshotID) (*arachne.ProtectedEntityInfo, error) {
	return nil, nil
}

func (this IVDProtectedEntity) GetComponents(ctx context.Context) ([]arachne.ProtectedEntity, error) {
	return make([]arachne.ProtectedEntity, 0), nil
}

func (this IVDProtectedEntity) GetID() arachne.ProtectedEntityID {
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
