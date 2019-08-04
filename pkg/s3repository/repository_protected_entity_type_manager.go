package s3repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/vmware/arachne/pkg/arachne"
	"io"
	"log"
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
}

func NewS3RepositoryProtectedEntityTypeManager(typeName string, session session.Session, bucket string) (*ProtectedEntityTypeManager, error) {
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

func (this *ProtectedEntityTypeManager) peinfoName(id arachne.ProtectedEntityID) (string, error) {
	if !id.HasSnapshot() {
		return "", errors.New("Cannot store objects that do not have snapshots")
	}
	return this.peinfoPrefix + id.String(), nil
}

func (this *ProtectedEntityTypeManager) mdName(id arachne.ProtectedEntityID) (string, error) {
	if !id.HasSnapshot() {
		return "", errors.New("Cannot store objects that do not have snapshots")
	}
	return this.mdPrefix + id.String() + MD_SUFFIX, nil
}

func (this *ProtectedEntityTypeManager) dataName(id arachne.ProtectedEntityID) (string, error) {
	if !id.HasSnapshot() {
		return "", errors.New("Cannot store objects that do not have snapshots")
	}
	return this.dataPrefix + id.String() + DATA_SUFFIX, nil
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
	peKey, err := this.peinfoName(id)
	if err != nil {
		return nil, err
	}
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

func (this *ProtectedEntityTypeManager) Copy(ctx context.Context, pe arachne.ProtectedEntity) (arachne.ProtectedEntity, error) {
	if pe.GetID().GetPeType() != this.typeName {
		return nil, errors.New(pe.GetID().GetPeType() + " is not of type " + this.typeName)
	}
	peInfo, err := pe.GetInfo(ctx)
	if err != nil {
		return nil, err
	}
	dataReader, err := pe.GetDataReader()
	if err != nil {
		return nil, err
	}

	mdReader, err := pe.GetMetadataReader()
	if err != nil {
		return nil, err
	}
	return this.copy(ctx, peInfo, dataReader, mdReader)
}

func (this *ProtectedEntityTypeManager) CopyFromInfo(ctx context.Context, info arachne.ProtectedEntityInfo) (arachne.ProtectedEntity, error) {

	return this.copy(ctx, info, nil, nil)
}

func (this *ProtectedEntityTypeManager)uploadStream(ctx context.Context, name string, reader io.Reader) error {
	uploader := s3manager.NewUploader(&this.session)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   reader,
		Bucket: aws.String(this.bucket),
		Key:    aws.String(name),
	})
	if err == nil {
		log.Println("Successfully uploaded to", result.Location)
	}
	return err
}
func (this *ProtectedEntityTypeManager) copy(ctx context.Context, peInfo arachne.ProtectedEntityInfo, dataReader io.Reader,
	metadataReader io.Reader) (arachne.ProtectedEntity, error) {
	peinfoName, err := this.peinfoName(peInfo.GetID())
	if err != nil {
		return nil, err
	}
	buf, err := json.Marshal(peInfo)
	if err != nil {
		return nil, err
	}
	if len(buf) > maxPEInfoSize {
		return nil, errors.New("JSON for pe info > 16K")
	}

	if dataReader != nil {
		dataName, err := this.dataName(peInfo.GetID())
		if err != nil {
			return nil, err
		}
		err = this.uploadStream(ctx, dataName, dataReader)
		if err != nil {
			return nil, err
		}
	}

	if metadataReader != nil {
		mdName, err := this.mdName(peInfo.GetID())
		if err != nil {
			return nil, err
		}
		err = this.uploadStream(ctx, mdName, metadataReader)
		if err != nil {
			return nil, err
		}
	}
	jsonBytes := bytes.NewReader(buf)

	jsonParams := &s3.PutObjectInput{
		Bucket:        aws.String(this.bucket),
		Key:           aws.String(peinfoName),
		Body:          jsonBytes,
		ContentLength: aws.Int64(int64(len(buf))),
		ContentType:   aws.String(peInfoFileType),
	}
	_, err = this.s3.PutObject(jsonParams)
	if err != nil {
		return nil, err
	}

	returnPEInfo := arachne.NewProtectedEntityInfo(peInfo.GetID(),
		peInfo.GetName(),
		nil, nil, nil, nil)
	returnPE := ProtectedEntity{
		rpetm:  this,
		peinfo: returnPEInfo,
	}
	return returnPE, err
}
