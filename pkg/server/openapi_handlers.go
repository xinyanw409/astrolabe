package server

import (
	"context"
	"github.com/go-openapi/runtime/middleware"
	"github.com/vmware-tanzu/astrolabe/gen/models"
	"github.com/vmware-tanzu/astrolabe/gen/restapi/operations"
	"github.com/vmware-tanzu/astrolabe/pkg/astrolabe"
	"time"
)

type OpenAPIArachneHandler struct {
	pem astrolabe.ProtectedEntityManager
	tm  TaskManager
}


func NewOpenAPIArachneHandler(pem astrolabe.ProtectedEntityManager, tm TaskManager) OpenAPIArachneHandler {
	return OpenAPIArachneHandler{
		pem: pem,
		tm: tm,
	}
}
func (this OpenAPIArachneHandler) AttachHandlers(api *operations.AstrolabeAPI) {
	api.ListServicesHandler = operations.ListServicesHandlerFunc(this.ListServices)
	api.ListProtectedEntitiesHandler = operations.ListProtectedEntitiesHandlerFunc(this.ListProtectedEntities)
	api.GetProtectedEntityInfoHandler = operations.GetProtectedEntityInfoHandlerFunc(this.GetProtectedEntityInfo)
	api.CreateSnapshotHandler = operations.CreateSnapshotHandlerFunc(this.CreateSnapshot)
	api.ListSnapshotsHandler = operations.ListSnapshotsHandlerFunc(this.ListSnapshots)
	api.CopyProtectedEntityHandler = operations.CopyProtectedEntityHandlerFunc(this.CopyProtectedEntity)
}

func (this OpenAPIArachneHandler) ListServices(params operations.ListServicesParams) middleware.Responder {
	etms := this.pem.ListEntityTypeManagers()
	var serviceNames = make(models.ServiceList, len(etms))
	for curETMNum, curETM := range etms {
		serviceNames[curETMNum] = curETM.GetTypeName()
	}
	return operations.NewListServicesOK().WithPayload(serviceNames)
}

func (this OpenAPIArachneHandler) ListProtectedEntities(params operations.ListProtectedEntitiesParams) middleware.Responder {
	petm := this.pem.GetProtectedEntityTypeManager(params.Service)
	if petm == nil {
		return operations.NewListProtectedEntitiesNotFound()
	}
	peids, err := petm.GetProtectedEntities(context.Background())
	if err != nil {

	}
	mpeids := make([]models.ProtectedEntityID, len(peids))
	for peidNum, peid := range peids {
		mpeids[peidNum] = models.ProtectedEntityID(peid.String())
	}
	peList := models.ProtectedEntityList{
		List:      mpeids,
		Truncated: false,
	}
	return operations.NewListProtectedEntitiesOK().WithPayload(&peList)
}

func (this OpenAPIArachneHandler) GetProtectedEntityInfo(params operations.GetProtectedEntityInfoParams) middleware.Responder {

	petm := this.pem.GetProtectedEntityTypeManager(params.Service)
	if petm == nil {

	}
	peid, err := astrolabe.NewProtectedEntityIDFromString(params.ProtectedEntityID)
	if err != nil {

	}
	pe, err := petm.GetProtectedEntity(context.Background(), peid)
	if err != nil {

	}
	peInfo, err := pe.GetInfo(context.Background())
	peInfoResponse := peInfo.GetModelProtectedEntityInfo()
	return operations.NewGetProtectedEntityInfoOK().WithPayload(&peInfoResponse);
}

func (this OpenAPIArachneHandler) CreateSnapshot(params operations.CreateSnapshotParams) middleware.Responder {
	petm := this.pem.GetProtectedEntityTypeManager(params.Service)
	if petm == nil {

	}
	peid, err := astrolabe.NewProtectedEntityIDFromString(params.ProtectedEntityID)
	if err != nil {

	}
	pe, err := petm.GetProtectedEntity(context.Background(), peid)
	if err != nil {

	}
	snapshotID, err := pe.Snapshot(context.Background())
	if err != nil {

	}

	return operations.NewCreateSnapshotOK().WithPayload(snapshotID.GetModelProtectedEntitySnapshotID())
}

func (this OpenAPIArachneHandler) ListSnapshots(params operations.ListSnapshotsParams) middleware.Responder {
	return nil
}


func (this OpenAPIArachneHandler) CopyProtectedEntity(params operations.CopyProtectedEntityParams) middleware.Responder {
	petm := this.pem.GetProtectedEntityTypeManager(params.Service)
	if petm == nil {

	}
	pei, err := astrolabe.NewProtectedEntityInfoFromModel(*params.Body)
	if err != nil {

	}
	startedTime := time.Now()
	newPE, err := petm.CopyFromInfo(context.Background(), pei, astrolabe.AllocateNewObject)
	var taskStatus astrolabe.TaskStatus
	if err != nil {
		taskStatus = astrolabe.Failed
	} else {
		taskStatus = astrolabe.Success
	}
	// Fake a task for now
	task := astrolabe.NewGenericTask()
	task.Completed = true
	task.StartedTime = startedTime
	task.FinishedTime = time.Now()
	task.Progress = 100
	task.TaskStatus = taskStatus
	task.Result = newPE.GetID().GetModelProtectedEntityID()
	return operations.NewCopyProtectedEntityAccepted()
}

