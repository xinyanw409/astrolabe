package s3repository

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/vmware/arachne/pkg/arachne"
	"github.com/vmware/arachne/pkg/fs"
	"log"
	"testing"
)

func TestProtectedEntityTypeManager(t *testing.T) {
	s3petm, err := setupPETM(t, "test")
	if err != nil {
		t.Fatal(err)
	}
	ids, err := s3petm.GetProtectedEntities(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("# of PEs returned = %d\n", len(ids))

	for _, id := range ids {
		t.Logf("%s\n", id)

	}
}

func setupPETM(t *testing.T, typeName string) (*ProtectedEntityTypeManager, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-1")},
	)
	if err != nil {
		return nil, err
	}
	s3petm, err := NewS3RepositoryProtectedEntityTypeManager(typeName, *sess, "dsu-velero")
	if err != nil {
		return nil, err
	}
	return s3petm, err
}

func TestCreateDeleteProtectedEntity(t *testing.T) {
	s3petm, err := setupPETM(t, "test")
	if err != nil {
		t.Fatal(err)
	}
	peID := arachne.NewProtectedEntityIDWithSnapshotID("test", "unique1", arachne.NewProtectedEntitySnapshotID("snapshot1"))
	peInfo := arachne.NewProtectedEntityInfo(peID, "testPE", nil, nil, nil, nil)
	repoPE, err := s3petm.CopyFromInfo(context.Background(), peInfo, arachne.AllocateNewObject)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Create repo PE %s\n", repoPE.GetID().String())

}

func TestCopyFSProtectedEntity(t *testing.T) {
	s3petm, err := setupPETM(t, "fs")
	if err != nil {
		t.Fatal(err)
	}

	fsParams := make(map[string]interface{})
	fsParams["root"] = "/Users/dsmithuchida/arachne_fs_root"

	fsPETM, err := fs.NewFSProtectedEntityTypeManagerFromConfig(fsParams, "notUsed")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	fsPEs, err := fsPETM.GetProtectedEntities(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, fsPEID := range  fsPEs {
		// FS doesn't have snapshots, but repository likes them, so fake one
		snapPEID := arachne.NewProtectedEntityIDWithSnapshotID(fsPEID.GetPeType(), fsPEID.GetID(),
			arachne.NewProtectedEntitySnapshotID("dummy-snap-id"))
		fsPE, err := fsPETM.GetProtectedEntity(ctx, snapPEID)
		if err != nil {
			t.Fatal(err)
		}
		s3PE, err := s3petm.Copy(ctx, fsPE, arachne.AllocateNewObject)
		if err != nil {
			t.Fatal(err)
		}

		newFSPE, err := fsPETM.Copy(ctx, s3PE, arachne.AllocateNewObject)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("Restored new FSPE %s\n", newFSPE.GetID().String())
	}
}