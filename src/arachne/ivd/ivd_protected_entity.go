package ivd

import (
	"arachne"
)

type ivdProtectedEntity struct {
	id arachne.ProtectedEntityID
	info arachne.ProtectedEntityInfo
}

func (ipe *ivdProtectedEntity) GetInfo() arachne.ProtectedEntityInfo {
	return ipe.info
}
func (ipe *ivdProtectedEntity) GetCombinedInfo() [] arachne.ProtectedEntityInfo {
	return make([]arachne.ProtectedEntityInfo, 0)
}
	/*
	 * Snapshot APIs
	 */
func (ipe *ivdProtectedEntity) Snapshot() (snapshotID arachne.ProtectedEntitySnapshotID) {
	return snapshotID
}
func (ipe *ivdProtectedEntity) ListSnapshots() [] arachne.ProtectedEntitySnapshotID {
	return make([]arachne.ProtectedEntitySnapshotID, 0)
}
func (ipe *ivdProtectedEntity) DeleteSnapshot(snapshotToDelete arachne.ProtectedEntitySnapshotID) bool {
	return true
}
func (ipe *ivdProtectedEntity) GetInfoForSnapshot(snapshotID arachne.ProtectedEntitySnapshotID) (info arachne.ProtectedEntityInfo) {
	return info
}

func (ipe *ivdProtectedEntity) GetComponents() [] arachne.ProtectedEntity {
	return make([]arachne.ProtectedEntity, 0)
}

func (ipe *ivdProtectedEntity) GetID() arachne.ProtectedEntityID {
	return ipe.id
}