package ivd

import (
	"arachne"
	vim "github.com/vmware/govmomi/vim25/types"
//	"github.com/vmware/govmomi/vslm"
)

type IVDProtectedEntity struct {
	ipetm *IVDProtectedEntityTypeManager
	id arachne.ProtectedEntityID
	info arachne.ProtectedEntityInfo
}
func newProtectedEntityID(id vim.ID) (arachne.ProtectedEntityID) {
	return arachne.NewProtectedEntityID("ivd", id.Id)
}

func newProtectedEntityIDWithSnapshotID(id vim.ID, snapshotID arachne.ProtectedEntitySnapshotID) (arachne.ProtectedEntityID) {
	return arachne.NewProtectedEntityIDWithSnapshotID("ivd", id.Id, snapshotID)
}

func newIVDProtectedEntity(ipetm *IVDProtectedEntityTypeManager, id arachne.ProtectedEntityID) (IVDProtectedEntity, error) {
	newIPE := IVDProtectedEntity {
		ipetm: ipetm,
		id: id,
	}
	return newIPE, nil
}
func (ipe *IVDProtectedEntity) GetInfo() arachne.ProtectedEntityInfo {
	return ipe.info
}
func (ipe *IVDProtectedEntity) GetCombinedInfo() [] arachne.ProtectedEntityInfo {
	return make([]arachne.ProtectedEntityInfo, 0)
}
	/*
	 * Snapshot APIs
	 */
func (ipe *IVDProtectedEntity) Snapshot() (snapshotID arachne.ProtectedEntitySnapshotID) {
	return snapshotID
}
func (ipe *IVDProtectedEntity) ListSnapshots() [] arachne.ProtectedEntitySnapshotID {
	return make([]arachne.ProtectedEntitySnapshotID, 0)
}
func (ipe *IVDProtectedEntity) DeleteSnapshot(snapshotToDelete arachne.ProtectedEntitySnapshotID) bool {
	return true
}
func (ipe *IVDProtectedEntity) GetInfoForSnapshot(snapshotID arachne.ProtectedEntitySnapshotID) (info arachne.ProtectedEntityInfo) {
	return info
}

func (ipe *IVDProtectedEntity) GetComponents() [] arachne.ProtectedEntity {
	return make([]arachne.ProtectedEntity, 0)
}

func (ipe *IVDProtectedEntity) GetID() arachne.ProtectedEntityID {
	return ipe.id
}