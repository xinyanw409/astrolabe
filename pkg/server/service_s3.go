package server

import (
	"context"
	"github.com/labstack/echo"
	"github.com/vmware/arachne/pkg/arachne"
	"io"
	"net/http"
	"strings"
)

type ServiceS3 struct {
	petm *arachne.ProtectedEntityTypeManager
}

func NewServiceS3(petm *arachne.ProtectedEntityTypeManager) *ServiceS3 {
	return &ServiceS3{
		petm: petm,
	}
}

func (this *ServiceS3) listObjects(echoContext echo.Context) error {
	/*
	 * No, this is not a correct implementation of S3 list bucket.
	 * TODO - write a proper implementation
	 */
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

func (this *ServiceS3) handleObjectRequest(echoContext echo.Context) error {
	objectKey := echoContext.Param("objectKey")
	var objectStream io.Reader
	var idStr, source, contentType string
	if (strings.HasSuffix(objectKey, ".md")) {
		idStr = strings.TrimSuffix(objectKey, ".md")
		source = "md"
		contentType = "application/octet-stream"
	} else {
		idStr = objectKey
		source = "data"
		contentType = "application/octet-stream"
	}

	_, pe, err := getProtectedEntityForIDStr(*this.petm, idStr, echoContext)
	if (err != nil) {

	}

	switch (source) {
	case "md":
		objectStream, err = pe.GetMetadataReader()
	case "data":
		objectStream, err = pe.GetDataReader()
	}
	if (err != nil) {

	}

	echoContext.Stream(http.StatusOK, contentType, objectStream)
	return nil
}
