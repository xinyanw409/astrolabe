package server

import (
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

func (this OpenAPIArachneHandler) Handle(params operations.ListServicesParams) middleware.Responder {
	etms := this.pem.ListEntityTypeManagers()
	var serviceNames = make(models.ServiceList, len(etms))
	for curETMNum, curETM := range etms {
		serviceNames[curETMNum] = curETM.GetTypeName()
	}
	return operations.NewListServicesOK().WithPayload(serviceNames)
}