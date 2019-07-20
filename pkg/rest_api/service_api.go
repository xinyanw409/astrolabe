package rest_api

import (
	"context"
	"github.com/labstack/echo"
	"github.com/vmware/arachne/pkg/arachne"
	"net/http"
)

type ServiceAPI struct {
	petm *arachne.ProtectedEntityTypeManager
}

func NewServiceAPI(petm *arachne.ProtectedEntityTypeManager) *ServiceAPI {
	return &ServiceAPI{
		petm: petm,
	}
}

func (this *ServiceAPI) listObjects(echoContext echo.Context) error {

	pes, err := (*this.petm).GetProtectedEntities(context.Background())
	if err != nil {
		return err
	}
	var pesList []string
	for _, curPes := range pes {
		pesList = append(pesList, curPes.GetID().String())
	}
	echoContext.JSON(http.StatusOK, pesList)
	return nil
}

func (this *ServiceAPI) handleObjectRequest(echoContext echo.Context) error {
	idStr := echoContext.Param("id")
	id, pe, err := this.getProtectedEntityForIDStr(idStr, echoContext)
	if err != nil {
		return nil
	}

	if action, ok := echoContext.Request().URL.Query()["action"]; ok {
		switch action[0] {
		case "snapshot":
			this.snapshot(echoContext, pe)
		case "deleteSnapshot":
			this.deleteSnapshot(echoContext, pe)
		default:
			echoContext.String(http.StatusBadRequest, "Action "+action[0]+" not understood")
		}
		return nil
	}
	info, err := pe.GetInfo(context.Background())
	if err != nil {
		echoContext.String(http.StatusNotFound, "Could not retrieve info for id "+id.String()+" error = "+err.Error())
		return nil
	}
	echoContext.JSON(http.StatusOK, info)

	return nil
}

func (this *ServiceAPI) snapshot(echoContext echo.Context, pe arachne.ProtectedEntity) {
	snapshotID, err := pe.Snapshot(context.Background())
	if err != nil {
		echoContext.String(http.StatusNotFound, "Snapshot failed for id "+pe.GetID().String()+" error = "+err.Error())
		return
	}
	if snapshotID == nil {
		echoContext.String(http.StatusInternalServerError, "snapshotID was nil for "+pe.GetID().String())
		return
	}
	echoContext.String(http.StatusOK, snapshotID.String())
}

func (this *ServiceAPI) deleteSnapshot(echoContext echo.Context, pe arachne.ProtectedEntity) {
	snapshotID := pe.GetID().GetSnapshotID()
	if snapshotID.GetID() == "" {
		echoContext.String(http.StatusBadRequest, "No snapshot ID specified in id "+pe.GetID().String()+" for delete")
		return
	}
	deleted, err := pe.DeleteSnapshot(context.Background(), snapshotID)
	if err != nil {
		echoContext.String(http.StatusNotFound, "Snapshot delete failed for id "+pe.GetID().String()+" error = "+err.Error())
		return
	}
	if deleted == false {
		echoContext.String(http.StatusInternalServerError, "Could not delete snapshot "+pe.GetID().String())
		return
	}
	echoContext.String(http.StatusOK, snapshotID.String())
}
func (this *ServiceAPI) getProtectedEntityForIDStr(idStr string, echoContext echo.Context) (arachne.ProtectedEntityID, arachne.ProtectedEntity, error) {
	var id arachne.ProtectedEntityID
	var pe arachne.ProtectedEntity
	var err error

	id, err = arachne.NewProtectedEntityIDFromString(idStr)
	if err != nil {
		echoContext.String(http.StatusBadRequest, "id = "+idStr+" is invalid "+err.Error())
		return id, pe, err
	}
	if id.GetPeType() != (*this.petm).GetTypeName() {
		echoContext.String(http.StatusBadRequest, "id = "+idStr+" is not type "+(*this.petm).GetTypeName())
		return id, pe, err
	}
	pe, err = (*this.petm).GetProtectedEntity(context.Background(), id)
	if err != nil {
		echoContext.String(http.StatusNotFound, "Could not retrieve id "+id.String()+" error = "+err.Error())
		return id, pe, err
	}
	if pe == nil {
		echoContext.String(http.StatusInternalServerError, "pe was nil for "+id.String())
		return id, pe, err
	}
	return id, pe, nil
}

func (this *ServiceAPI) handleSnapshotListRequest(echoContext echo.Context) error {
	idStr := echoContext.Param("id")
	id, pe, err := this.getProtectedEntityForIDStr(idStr, echoContext)
	if err != nil {
		return nil
	}
	snapshotIDs, err := pe.ListSnapshots(context.Background())
	if pe == nil {
		echoContext.String(http.StatusInternalServerError, "Could not retrieve snapshots "+id.String())
		return nil
	}
	snapshotIDStrs := []string{}
	for _, curSnapshotID := range snapshotIDs {
		snapshotIDStrs = append(snapshotIDStrs, curSnapshotID.String())
	}
	echoContext.JSON(http.StatusOK, snapshotIDStrs)
	return nil
}
