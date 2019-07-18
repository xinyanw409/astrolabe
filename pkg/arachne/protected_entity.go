package arachne

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
)

type ProtectedEntityID struct {
	peType     string
	id         string
	snapshotID ProtectedEntitySnapshotID
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
	/*
		components := strings.Split(peiString, ":")
		if len(components) > 1 {
			returnPEI.peType = components[0]
			returnPEI.id = components[1]
			if len(components) == 3 {
				returnPEI.snapshotID = *NewProtectedEntitySnapshotID(components[2])
			}
		} else {
			returnError = errors.New("arachne: '" + peiString + "' is not a valid protected entity ID")
		}
		return returnPEI, returnError
	*/
	returnError = fillInProtectedEntityIDFromString(&returnPEI, peiString)
	return returnPEI, returnError
}

func fillInProtectedEntityIDFromString(pei *ProtectedEntityID, peiString string) error {
	components := strings.Split(peiString, ":")
	if len(components) > 1 {
		pei.peType = components[0]
		pei.id = components[1]
		if len(components) == 3 {
			pei.snapshotID = *NewProtectedEntitySnapshotID(components[2])
		}
		log.Print("pei = " + pei.String())
	} else {
		return errors.New("arachne: '" + peiString + "' is not a valid protected entity ID")
	}
	return nil
}
func (peid ProtectedEntityID) GetID() string {
	return peid.id
}

func (peid ProtectedEntityID) GetPeType() string {
	return peid.peType
}

func (peid ProtectedEntityID) GetSnapshotID() ProtectedEntitySnapshotID {
	return peid.snapshotID

}

func (peid ProtectedEntityID) String() string {
	var returnString string
	returnString = peid.peType + ":" + peid.id
	if (peid.snapshotID) != (ProtectedEntitySnapshotID{}) {
		returnString += ":" + peid.snapshotID.String()
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

func NewProtectedEntitySnapshotID(pesiString string) *ProtectedEntitySnapshotID {
	returnPESI := ProtectedEntitySnapshotID{
		id: pesiString,
	}
	return &returnPESI
}

func (pesid *ProtectedEntitySnapshotID) GetID() string {
	return pesid.id
}

func (pesid *ProtectedEntitySnapshotID) String() string {
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
}
