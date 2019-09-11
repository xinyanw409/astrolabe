package s3repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/vmware/arachne/pkg/arachne"
	"github.com/vmware/arachne/pkg/util"
	"io"
	"log"
)

type ProtectedEntity struct {
	rpetm  *ProtectedEntityTypeManager
	peinfo arachne.ProtectedEntityInfo
}

func NewProtectedEntityFromJSONBuf(rpetm *ProtectedEntityTypeManager, buf [] byte) (pe ProtectedEntity, err error) {
	peii := arachne.ProtectedEntityInfoImpl{}
	err = json.Unmarshal(buf, &peii)
	if err != nil {
		return
	}
	pe.peinfo = peii
	pe.rpetm = rpetm
	return
}

func NewProtectedEntityFromJSONReader(rpetm *ProtectedEntityTypeManager, reader io.Reader) (pe ProtectedEntity, err error) {
	decoder := json.NewDecoder(reader)
	peInfo := arachne.ProtectedEntityInfoImpl{}
	err = decoder.Decode(&peInfo)
	if err == nil {
		pe.peinfo = peInfo
		pe.rpetm = rpetm
	}
	return
}
func (this ProtectedEntity) GetInfo(ctx context.Context) (arachne.ProtectedEntityInfo, error) {
	return this.peinfo, nil
}

func (ProtectedEntity) GetCombinedInfo(ctx context.Context) ([]arachne.ProtectedEntityInfo, error) {
	panic("implement me")
}

func (ProtectedEntity) Snapshot(ctx context.Context) (*arachne.ProtectedEntitySnapshotID, error) {
	return nil, errors.New("Snapshot not supported")
}

func (ProtectedEntity) ListSnapshots(ctx context.Context) ([]arachne.ProtectedEntitySnapshotID, error) {
	panic("implement me")
}

func (ProtectedEntity) DeleteSnapshot(ctx context.Context, snapshotToDelete arachne.ProtectedEntitySnapshotID) (bool, error) {
	panic("implement me")
}

func (ProtectedEntity) GetInfoForSnapshot(ctx context.Context, snapshotID arachne.ProtectedEntitySnapshotID) (*arachne.ProtectedEntityInfo, error) {
	panic("implement me")
}

func (ProtectedEntity) GetComponents(ctx context.Context) ([]arachne.ProtectedEntity, error) {
	panic("implement me")
}

func (this ProtectedEntity) GetID() arachne.ProtectedEntityID {
	return this.peinfo.GetID()
}

func (this ProtectedEntity) GetDataReader(context.Context) (io.Reader, error) {
	if len(this.peinfo.GetDataTransports()) > 0 {
		dataName := this.rpetm.dataName(this.GetID())
		return this.getReader(dataName)
	}
	return nil, nil
}

func (this ProtectedEntity) GetMetadataReader(context.Context) (io.Reader, error) {
	if len(this.peinfo.GetMetadataTransports()) > 0 {
		metadataName := this.rpetm.metadataName(this.GetID())
		return this.getReader(metadataName)
	}
	return nil, nil
}

func (this *ProtectedEntity) uploadStream(ctx context.Context, name string, reader io.Reader) error {
	uploader := s3manager.NewUploader(&this.rpetm.session)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   reader,
		Bucket: aws.String(this.rpetm.bucket),
		Key:    aws.String(name),
	})
	if err == nil {
		log.Println("Successfully uploaded to", result.Location)
	}
	return err
}

func (this *ProtectedEntity) copy(ctx context.Context, dataReader io.Reader,
	metadataReader io.Reader) error {
	peInfo := this.peinfo
	peinfoName := this.rpetm.peinfoName(peInfo.GetID())

	peInfoBuf, err := json.Marshal(peInfo)
	if err != nil {
		return err
	}
	if len(peInfoBuf) > maxPEInfoSize {
		return errors.New("JSON for pe info > 16K")
	}

	if dataReader != nil {
		dataName := this.rpetm.dataName(peInfo.GetID())
		err = this.uploadStream(ctx, dataName, dataReader)
		if err != nil {
			return err
		}
	}

	if metadataReader != nil {
		mdName := this.rpetm.metadataName(peInfo.GetID())
		err = this.uploadStream(ctx, mdName, metadataReader)
		if err != nil {
			return err
		}
	}
	jsonBytes := bytes.NewReader(peInfoBuf)

	jsonParams := &s3.PutObjectInput{
		Bucket:        aws.String(this.rpetm.bucket),
		Key:           aws.String(peinfoName),
		Body:          jsonBytes,
		ContentLength: aws.Int64(int64(len(peInfoBuf))),
		ContentType:   aws.String(peInfoFileType),
	}
	_, err = this.rpetm.s3.PutObject(jsonParams)
	if err != nil {
		return err
	}
	return err
}

func (this *ProtectedEntity) getReader(key string) (io.Reader, error) {
	downloadMgr := s3manager.NewDownloaderWithClient(&this.rpetm.s3, func(d *s3manager.Downloader) {
		d.Concurrency = 1
		//d.PartSize = 1
	})
	reader, writer := io.Pipe()
	seqWriterAt := util.NewSeqWriterAt(writer)
	go func() {
		defer writer.Close()
		downloadMgr.Download(seqWriterAt, &s3.GetObjectInput{
			Bucket: aws.String(this.rpetm.bucket),
			Key:    aws.String(key),
		})
		fmt.Printf("Download finished")
	}()

	return reader, nil
}
