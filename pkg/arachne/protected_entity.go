package arachne

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/vmware/gvddk/gDiskLib"
	"io"
	"log"
	"strings"
)

type ProtectedEntityID struct {
	peType     string
	id         string
	snapshotID ProtectedEntitySnapshotID
}

type DiskConnectionParam struct {
	gDiskLib.DiskHandle
	gDiskLib.VixDiskLibConnection
	gDiskLib.ConnectParams
}

func NewProtectedEntityID(peType string, id string) ProtectedEntityID {
	return NewProtectedEntityIDWithSnapshotID(peType, id, ProtectedEntitySnapshotID{})
}

func NewProtectedEntityIDWithSnapshotID(peType string, id string, snapshotID ProtectedEntitySnapshotID) ProtectedEntityID {
	newID := ProtectedEntityID{
		peType:     peType,
		id:         id,
		snapshotID: snapshotID,
	}
	return newID
}

func NewProtectedEntityIDFromString(peiString string) (returnPEI ProtectedEntityID, returnError error) {
	returnError = fillInProtectedEntityIDFromString(&returnPEI, peiString)
	return returnPEI, returnError
}

func fillInProtectedEntityIDFromString(pei *ProtectedEntityID, peiString string) error {
	components := strings.Split(peiString, ":")
	if len(components) > 1 {
		pei.peType = components[0]
		pei.id = components[1]
		if len(components) == 3 {
			pei.snapshotID = NewProtectedEntitySnapshotID(components[2])
		}
		log.Print("pei = " + pei.String())
	} else {
		return errors.New("arachne: '" + peiString + "' is not a valid protected entity ID")
	}
	return nil
}
func (this ProtectedEntityID) GetID() string {
	return this.id
}

func (this ProtectedEntityID) GetPeType() string {
	return this.peType
}

func (this ProtectedEntityID) GetSnapshotID() ProtectedEntitySnapshotID {
	return this.snapshotID

}

func (this ProtectedEntityID) HasSnapshot() bool {
	return this.snapshotID.id != ""
}

func (this ProtectedEntityID) String() string {
	var returnString string
	returnString = this.peType + ":" + this.id
	if (this.snapshotID) != (ProtectedEntitySnapshotID{}) {
		returnString += ":" + this.snapshotID.String()
	}
	return returnString
}

func (this ProtectedEntityID) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.String()) // Use marshal to make sure encoding happens
}

func (this *ProtectedEntityID) UnmarshalJSON(b []byte) error {
	var idStr string
	json.Unmarshal(b, &idStr) // Use unmarshall to make sure decoding happens
	log.Print("UnmarshalJSON idStr = " + idStr)
	return fillInProtectedEntityIDFromString(this, idStr)
}

type ProtectedEntitySnapshotID struct {
	// We should move this to actually being a UUID internally
	id string
}

func NewProtectedEntitySnapshotID(pesiString string) ProtectedEntitySnapshotID {
	returnPESI := ProtectedEntitySnapshotID{
		id: pesiString,
	}
	return returnPESI
}

func (pesid ProtectedEntitySnapshotID) GetID() string {
	return pesid.id
}

func (pesid ProtectedEntitySnapshotID) String() string {
	return pesid.id
}

type ProtectedEntity interface {
	GetInfo(ctx context.Context) (ProtectedEntityInfo, error)
	GetCombinedInfo(ctx context.Context) ([]ProtectedEntityInfo, error)
	/*
	 * Snapshot APIs
	 */
	Snapshot(ctx context.Context) (*ProtectedEntitySnapshotID, error)
	ListSnapshots(ctx context.Context) ([]ProtectedEntitySnapshotID, error)
	DeleteSnapshot(ctx context.Context, snapshotToDelete ProtectedEntitySnapshotID) (bool, error)
	GetInfoForSnapshot(ctx context.Context, snapshotID ProtectedEntitySnapshotID) (*ProtectedEntityInfo, error)

	GetComponents(ctx context.Context) ([]ProtectedEntity, error)
	GetID() ProtectedEntityID

	// GetDataReader returns a reader for the data of the ProtectedEntity.  The ProtectedEntity will pick the
	// best data path to provide the Reader stream.  If the ProtectedEntity does not have any data, nil will be
	// returned
	GetDataReader(ctx context.Context) (io.Reader, error)

	// GetMetadataReader returns a reader for the metadata of the ProtectedEntity.  The ProtectedEntity will pick the
	// best data path to provide the Reader stream.  If the ProtectedEntity does not have any metadata, nil will be
	// returned
	GetMetadataReader(ctx context.Context) (io.Reader, error)
}
