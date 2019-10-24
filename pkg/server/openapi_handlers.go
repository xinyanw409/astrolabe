package server

import (
	"context"
	"github.com/go-openapi/runtime/middleware"
	"github.com/vmware/arachne/gen/models"
	"github.com/vmware/arachne/gen/restapi/operations"
	"github.com/vmware/arachne/pkg/arachne"
)

type OpenAPIArachneHandler struct {
	pem arachne.ProtectedEntityManager
}

func NewOpenAPIArachneHandler(pem arachne.ProtectedEntityManager) OpenAPIArachneHandler {
	return OpenAPIArachneHandler{
		pem: pem,
	}
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