package s3repository

import (
	"context"
	"encoding/json"
	"github.com/vmware/arachne/pkg/arachne"
	"io"
)

type ProtectedEntity struct {
	rpetm    *ProtectedEntityTypeManager
	peinfo   arachne.ProtectedEntityInfo
}

func NewProtectedEntityFromJSONBuf(rpetm *ProtectedEntityTypeManager, buf [] byte) (pe ProtectedEntity, err error) {
	peii := arachne.ProtectedEntityInfoImpl{}
	err = json.Unmarshal(buf, &peii)
	if err != nil {
		return
	}
	pe.peinfo = peii
	pe.rpetm = rpetm
	return
}

func NewProtectedEntityFromJSONReader(rpetm * ProtectedEntityTypeManager, reader io.Reader) (pe ProtectedEntity, err error) {
	decoder:= json.NewDecoder(reader)
	err = decoder.Decode(&pe)
	return
}
func (this ProtectedEntity) GetInfo(ctx context.Context) (arachne.ProtectedEntityInfo, error) {
	return this.peinfo, nil
}

func (ProtectedEntity) GetCombinedInfo(ctx context.Context) ([]arachne.ProtectedEntityInfo, error) {
	panic("implement me")
}

func (ProtectedEntity) Snapshot(ctx context.Context) (*arachne.ProtectedEntitySnapshotID, error) {
	panic("implement me")
}

func (ProtectedEntity) ListSnapshots(ctx context.Context) ([]arachne.ProtectedEntitySnapshotID, error) {
	panic("implement me")
}

func (ProtectedEntity) DeleteSnapshot(ctx context.Context, snapshotToDelete arachne.ProtectedEntitySnapshotID) (bool, error) {
	panic("implement me")
}

func (ProtectedEntity) GetInfoForSnapshot(ctx context.Context, snapshotID arachne.ProtectedEntitySnapshotID) (*arachne.ProtectedEntityInfo, error) {
	panic("implement me")
}

func (ProtectedEntity) GetComponents(ctx context.Context) ([]arachne.ProtectedEntity, error) {
	panic("implement me")
}

func (this ProtectedEntity) GetID() arachne.ProtectedEntityID {
	return this.peinfo.GetID()
}

func (ProtectedEntity) GetDataReader() (io.Reader, error) {
	panic("implement me")
}

func (ProtectedEntity) GetMetadataReader() (io.Reader, error) {
	panic("implement me")
}
