package server

import (
	"context"
	"github.com/go-openapi/runtime/middleware"
	"github.com/vmware/arachne/gen/models"
	"github.com/vmware/arachne/gen/restapi/operations"
	"github.com/vmware/arachne/pkg/arachne"
	"github.com/vmware/govmomi/task"
	"time"
)

type OpenAPIArachneHandler struct {
	pem arachne.ProtectedEntityManager
	tm TaskManager
}


func NewOpenAPIArachneHandler(pem arachne.ProtectedEntityManager, tm TaskManager) OpenAPIArachneHandler {
	return OpenAPIArachneHandler{
		pem: pem,
		tm: tm,
	}
}
func (this OpenAPIArachneHandler) AttachHandlers(api *operations.ArachneAPI) {
	api.ListServicesHandler = operations.ListServicesHandlerFunc(this.ListServices)
	api.ListProtectedEntitiesHandler = operations.ListProtectedEntitiesHandlerFunc(this.ListProtectedEntities)
	api.GetProtectedEntityInfoHandler = operations.GetProtectedEntityInfoHandlerFunc(this.GetProtectedEntityInfo)
	api.CreateSnapshotHandler = operations.CreateSnapshotHandlerFunc(this.CreateSnapshot)
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
	peid, err := arachne.NewProtectedEntityIDFromString(params.ProtectedEntityID)
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
	peid, err := arachne.NewProtectedEntityIDFromString(params.ProtectedEntityID)
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

func (this OpenAPIArachneHandler) CopyProtectedEntity(params operations.CopyProtectedEntityParams) middleware.Responder {
	petm := this.pem.GetProtectedEntityTypeManager(params.Service)
	if petm == nil {

	}
	pei, err := arachne.NewProtectedEntityInfoFromModel(*params.Body)
	if err != nil {

	}
	startedTime := time.Now()
	newPE, err := petm.CopyFromInfo(context.Background(), pei, arachne.AllocateNewObject)
	var taskStatus arachne.TaskStatus
	if err != nil {
		taskStatus = arachne.TaskStatus{

		}
	}
	// Fake a task for now
	task := arachne.GenericTask{
		ID:           arachne.TaskID{},
		Completed:    true,
		TaskStatus:   0,
		Details:      "",
		StartedTime:  startedTime,
		FinishedTime: time.Now(),
		Progress:     100,
	}
}