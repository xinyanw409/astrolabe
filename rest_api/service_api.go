package rest_api

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/vmware/arachne"
	"context"
)

type ServiceAPI struct {
	petm * arachne.ProtectedEntityTypeManager
}

func NewServiceAPI(petm * arachne.ProtectedEntityTypeManager) (*ServiceAPI) {
	return &ServiceAPI {
		petm: petm,
	}
}

func (this *ServiceAPI) listObjects(echoContext echo.Context) error {
	
	pes, err := (*this.petm).GetProtectedEntities(context.Background());
	if (err != nil) {
		return err
	}
	var pesList []string
	for _, curPes := range pes {
		pesList = append(pesList, curPes.GetID().String())
	}
	echoContext.JSON(http.StatusOK, pesList)
	return nil
}