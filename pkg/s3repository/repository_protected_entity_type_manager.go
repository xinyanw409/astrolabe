package s3repository

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"github.com/vmware/arachne/pkg/arachne"
	"io"
	"strings"
)

/*
 * ProtectedEntityTypeManager for an S3 repository acts as a passive, generic Protected Entity Type Manager
 * Protected Entities served by the type manager do not change and are always read-only.
 */
type ProtectedEntityTypeManager struct {
	typeName                                         string
	session                                          session.Session
	s3                                               s3.S3
	bucket                                           string
	objectPrefix, peinfoPrefix, mdPrefix, dataPrefix string
	logger    logrus.FieldLogger
}

func NewS3RepositoryProtectedEntityTypeManager(typeName string, session session.Session, bucket string, logger logrus.FieldLogger) (*ProtectedEntityTypeManager, error) {
	objectPrefix := "arachne-repo/" + typeName + "/"
	peinfoPrefix := objectPrefix + "peinfo/"
	mdPrefix := objectPrefix + "md/"
	dataPrefix := objectPrefix + "data/"
	returnPETM := ProtectedEntityTypeManager{
		typeName:     typeName,
		session:      session,
		s3:           *(s3.New(&session)),
		bucket:       bucket,
		objectPrefix: objectPrefix,
		peinfoPrefix: peinfoPrefix,
		mdPrefix:     mdPrefix,
		dataPrefix:   dataPrefix,
		logger:       logger,
	}
	return &returnPETM, nil
}

/*
 * Protected Entities are stored in the S3 repo as 1-3 files.  The peinfo file contains the Protected Entity JSON,
 * the md file contains the Protected Entity metadata, if present and the data file contains the Protected Entity data,
 * if present.  The basic structure of the repository is
 *    <bucket>/arachne-repo/<type>/{peinfo, md, data}/<peid>[, .md, .data]
 * The PEID must have a snapshot component
 * For example, an IVD would be represented as three S3 objects:
 *     /arachne-repo/ivd/peinfo/ivd:e1c3cb20-db88-4c1c-9f02-5f5347e435d5:67469e1c-50a8-4f63-9a6a-ad8a2265197c
 *     /arachne-repo/ivd/md/ivd:e1c3cb20-db88-4c1c-9f02-5f5347e435d5:67469e1c-50a8-4f63-9a6a-ad8a2265197c.md
 *     /arachne-repo/ivd/data/ivd:e1c3cb20-db88-4c1c-9f02-5f5347e435d5:67469e1c-50a8-4f63-9a6a-ad8a2265197c.data
 *
 * The combined stream is not stored in S3 but could be synthesized on demand (figure out how this would actually work)
 */
const MD_SUFFIX = ".md"
const DATA_SUFFIX = ".data"

func (this *ProtectedEntityTypeManager) peinfoName(id arachne.ProtectedEntityID) (string) {
	if !id.HasSnapshot() {
		panic("Cannot store objects that do not have snapshots")
	}
	return this.peinfoPrefix + id.String()
}

func (this *ProtectedEntityTypeManager) metadataName(id arachne.ProtectedEntityID) (string) {
	if !id.HasSnapshot() {
		panic("Cannot store objects that do not have snapshots")
	}
	return this.mdPrefix + id.String() + MD_SUFFIX
}

func (this *ProtectedEntityTypeManager) metadataTransportsForID(id arachne.ProtectedEntityID) ([]arachne.DataTransport, error) {
	endpointPtr := this.session.Config.Endpoint
	var endpoint string
	if endpointPtr != nil {
		endpoint = *endpointPtr
	} else {
		endpoint = ""
	}
	mdTransport := arachne.NewDataTransportForS3(endpoint, this.bucket, this.metadataName(id))
	return []arachne.DataTransport{mdTransport}, nil
}

func (this *ProtectedEntityTypeManager) dataName(id arachne.ProtectedEntityID) (string) {
	if !id.HasSnapshot() {
		panic("Cannot store objects that do not have snapshots")
	}
	return this.dataPrefix + id.String() + DATA_SUFFIX
}

func (this *ProtectedEntityTypeManager) dataTransportsForID(id arachne.ProtectedEntityID) ([]arachne.DataTransport, error) {
	endpointPtr := this.session.Config.Endpoint
	var endpoint string
	if endpointPtr != nil {
		endpoint = *endpointPtr
	} else {
		endpoint = ""
	}
	mdTransport := arachne.NewDataTransportForS3(endpoint, this.bucket, this.dataName(id))
	return []arachne.DataTransport{mdTransport}, nil
}

func (this *ProtectedEntityTypeManager) objectPEID(key string) (arachne.ProtectedEntityID, error) {
	var idStr string
	if strings.HasPrefix(key, this.peinfoPrefix) {
		idStr = strings.TrimPrefix(key, this.peinfoPrefix)
	}
	if strings.HasPrefix(key, this.mdPrefix) {
		if !strings.HasSuffix(key, MD_SUFFIX) {
			return arachne.ProtectedEntityID{}, errors.New(key + " has md prefix, but does not have .md suffix")
		}
		idStr = strings.TrimPrefix(key, this.mdPrefix)
		idStr = strings.TrimSuffix(key, MD_SUFFIX)
	}
	if strings.HasPrefix(key, this.dataPrefix) {
		if !strings.HasSuffix(key, DATA_SUFFIX) {
			return arachne.ProtectedEntityID{}, errors.New(key + " has data prefix, but does not have .data suffix")
		}
		idStr = strings.TrimPrefix(key, this.dataPrefix)
		idStr = strings.TrimSuffix(key, DATA_SUFFIX)
	}
	retPEID, err := arachne.NewProtectedEntityIDFromString(idStr)
	if err != nil {
		return arachne.ProtectedEntityID{}, err
	}
	return retPEID, nil
}

func (this *ProtectedEntityTypeManager) GetTypeName() string {
	return this.typeName
}

const maxPEInfoSize int = 16 * 1024

func (this *ProtectedEntityTypeManager) GetProtectedEntity(ctx context.Context, id arachne.ProtectedEntityID) (arachne.ProtectedEntity, error) {
	peKey := this.peinfoName(id)
	oi := s3.GetObjectInput{
		Bucket: &this.bucket,
		Key:    &peKey,
	}

	oo, err := this.s3.GetObject(&oi)
	if err != nil {
		return nil, err
	}
	returnPE, err := NewProtectedEntityFromJSONReader(this, oo.Body)
	if err != nil {
		return nil, err
	}
	return returnPE, nil
}

const maxS3ObjectsToFetch int64 = 1000

func (this *ProtectedEntityTypeManager) GetProtectedEntities(ctx context.Context) ([]arachne.ProtectedEntityID, error) {
	hasMore := true
	var continuationToken *string = nil
	prefix := this.peinfoPrefix
	retPEIDs := make([]arachne.ProtectedEntityID, 0)
	for hasMore {
		maxKeys := maxS3ObjectsToFetch
		listParams := s3.ListObjectsV2Input{
			Bucket:            aws.String(this.bucket),
			Prefix:            &prefix,
			ContinuationToken: continuationToken,
			MaxKeys:           &maxKeys,
		}

		results, err := this.s3.ListObjectsV2(&listParams)

		if err != nil {
			return nil, err
		}

		for _, item := range results.Contents {
			s3Key := *item.Key
			retPEID, err := this.objectPEID(s3Key)
			if err == nil {
				retPEIDs = append(retPEIDs, retPEID)
			} else {

			}
		}
		if !*results.IsTruncated {
			hasMore = false
		} else {
			continuationToken = results.ContinuationToken
		}
	}
	return retPEIDs, nil
}

const peInfoFileType = "application/json"

func (this *ProtectedEntityTypeManager) Copy(ctx context.Context, sourcePE arachne.ProtectedEntity,
	options arachne.CopyCreateOptions) (arachne.ProtectedEntity, error) {
	sourcePEInfo, err := sourcePE.GetInfo(ctx)
	if err != nil {
		return nil, err
	}
	dataReader, err := sourcePE.GetDataReader(ctx)
	if dataReader != nil {
		defer func() {
			if err := dataReader.Close(); err != nil {
				this.logger.Errorf("The deferred data reader is closed with error, %v", err)
			}
		}()
	}

	if err != nil {
		return nil, err
	}

	metadataReader, err := sourcePE.GetMetadataReader(ctx)
	if err != nil {
		return nil, err
	}
	return this.copyInt(ctx, sourcePEInfo, options, dataReader, metadataReader)
}

func (this *ProtectedEntityTypeManager) CopyFromInfo(ctx context.Context, sourcePEInfo arachne.ProtectedEntityInfo,
	options arachne.CopyCreateOptions) (arachne.ProtectedEntity, error) {
	return nil, nil
}

func (this *ProtectedEntityTypeManager) copyInt(ctx context.Context, sourcePEInfo arachne.ProtectedEntityInfo,
	options arachne.CopyCreateOptions, dataReader io.Reader, metadataReader io.Reader) (arachne.ProtectedEntity, error) {
	id := sourcePEInfo.GetID()
	if id.GetPeType() != this.typeName {
		return nil, errors.New(id.GetPeType() + " is not of type " + this.typeName)
	}
	if options == arachne.AllocateObjectWithID {
		return nil, errors.New("AllocateObjectWithID not supported")
	}

	if options == arachne.UpdateExistingObject {
		return nil, errors.New("UpdateExistingObject not supported")
	}

	_, err := this.GetProtectedEntity(ctx, id)
	if err == nil {
		return nil, errors.New("id " + id.String() + " already exists")
	}

	var dataTransports []arachne.DataTransport
	if len(sourcePEInfo.GetDataTransports()) > 0 {
		dataTransports, err = this.dataTransportsForID(id)
		if err != nil {
			return nil, err
		}
	} else {
		dataTransports = []arachne.DataTransport{}
	}

	var metadataTransports []arachne.DataTransport
	if len(sourcePEInfo.GetMetadataTransports()) > 0 {
		metadataTransports, err = this.metadataTransportsForID(id)
		if err != nil {
			return nil, err
		}
	} else {
		metadataTransports = []arachne.DataTransport{}
	}

	combinedTransports := []arachne.DataTransport{}

	rPEInfo := arachne.NewProtectedEntityInfo(sourcePEInfo.GetID(), sourcePEInfo.GetName(),
		dataTransports, metadataTransports, combinedTransports, sourcePEInfo.GetComponentIDs())

	rpe := ProtectedEntity{
		rpetm:  this,
		peinfo: rPEInfo,
	}

	err = rpe.copy(ctx, dataReader, metadataReader)
	if err != nil {
		return nil, err
	}
	return rpe, nil
}

func (this *ProtectedEntityTypeManager) getDataTransports(id arachne.ProtectedEntityID) ([]arachne.DataTransport,
	[]arachne.DataTransport,
	[]arachne.DataTransport, error) {
	dataS3URL := this.dataName(id)
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
