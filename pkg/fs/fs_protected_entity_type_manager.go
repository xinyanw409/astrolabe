package fs

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/vmware/arachne/pkg/arachne"
	"io"
	"io/ioutil"
	"path/filepath"
)

type FSProtectedEntityTypeManager struct {
	root      string
	s3URLBase string
}

const kTYPE_NAME = "fs"

func NewFSProtectedEntityTypeManagerFromConfig(params map[string]interface{}, s3URLBase string) (*FSProtectedEntityTypeManager, error) {
	root := params["root"].(string)

	returnTypeManager := FSProtectedEntityTypeManager{
		root:      root,
		s3URLBase: s3URLBase,
	}
	return &returnTypeManager, nil
}

func (this *FSProtectedEntityTypeManager) GetTypeName() string {
	return kTYPE_NAME
}

func (this *FSProtectedEntityTypeManager) GetProtectedEntity(ctx context.Context, id arachne.ProtectedEntityID) (
	arachne.ProtectedEntity, error) {
	return newFSProtectedEntity(this, id, id.GetID(), filepath.Join(this.root, id.GetID()))
}

func (this *FSProtectedEntityTypeManager) GetProtectedEntities(ctx context.Context) ([]arachne.ProtectedEntityID, error) {
	files, err := ioutil.ReadDir(this.root)
	if err != nil {
		return nil, err
	}

	var retVal = make([]arachne.ProtectedEntityID, len(files))
	for index, curFile := range files {
		peid := arachne.NewProtectedEntityID("fs", curFile.Name())
		retVal[index] = peid
	}
	return retVal, nil
}

func (this *FSProtectedEntityTypeManager) Copy(ctx context.Context, pe arachne.ProtectedEntity,
	options arachne.CopyCreateOptions) (arachne.ProtectedEntity, error) {

	sourcePEInfo, err := pe.GetInfo(ctx)
	if err != nil {
		return nil, err
	}
	dataReader, err := pe.GetDataReader()
	if err != nil {
		return nil, err
	}
	metadataReader, err := pe.GetMetadataReader()
	if err != nil {
		return nil, err
	}
	return this.copyInt(ctx, sourcePEInfo, options, dataReader, metadataReader)
}

func (this *FSProtectedEntityTypeManager) CopyFromInfo(ctx context.Context, pe arachne.ProtectedEntityInfo,
	options arachne.CopyCreateOptions) (arachne.ProtectedEntity, error) {
	return nil, nil
}

func (this *FSProtectedEntityTypeManager) copyInt(ctx context.Context, sourcePEInfo arachne.ProtectedEntityInfo,
	options arachne.CopyCreateOptions, dataReader io.Reader, metadataReader io.Reader) (arachne.ProtectedEntity, error) {
	id := sourcePEInfo.GetID()
	if id.GetPeType() != kTYPE_NAME {
		return nil, errors.New(id.GetPeType() + " is not of type fs")
	}
	if options == arachne.AllocateObjectWithID {
		return nil, errors.New("AllocateObjectWithID not supported")
	}

	if options == arachne.UpdateExistingObject {
		return nil, errors.New("UpdateExistingObject not supported")
	}

	fsUUID, err := uuid.NewRandom()
	if err != nil {
		panic("uuid.NewRandom return err ")
	}
	newPEID := arachne.NewProtectedEntityID(kTYPE_NAME, fsUUID.String())
	newPE, err := newFSProtectedEntity(this, newPEID, sourcePEInfo.GetName(), filepath.Join(this.root, newPEID.GetID()))
	if err != nil {
		return nil, err
	}
	err = newPE.createDir()
	if err != nil {
		return nil, err
	}
	err = newPE.copy(ctx, dataReader, metadataReader)
	if err != nil {
		return nil, err
	}
	return newPE, nil
}

func (this *FSProtectedEntityTypeManager) getDataTransports(id arachne.ProtectedEntityID) ([]arachne.DataTransport,
	[]arachne.DataTransport,
	[]arachne.DataTransport, error) {
	dataS3URL := this.s3URLBase + "fs/" + id.String()
	data := []arachne.DataTransport{
		arachne.NewDataTransportForS3URL(dataS3URL),
	}

	mdS3URL := dataS3URL + ".md"

	md := []arachne.DataTransport{
		arachne.NewDataTransportForS3URL(mdS3URL),
	}

	combinedS3URL := dataS3URL + ".zip"
	combined := []arachne.DataTransport{
		arachne.NewDataTransportForS3URL(combinedS3URL),
	}

	return data, md, combined, nil
}
