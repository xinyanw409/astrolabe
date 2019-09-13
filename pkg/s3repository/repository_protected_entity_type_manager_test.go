package s3repository

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sirupsen/logrus"
	"github.com/vmware/arachne/pkg/arachne"
	"github.com/vmware/arachne/pkg/fs"
	"github.com/vmware/arachne/pkg/ivd"
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
	s3petm, err := NewS3RepositoryProtectedEntityTypeManager(typeName, *sess, "velero-plugin-s3-repo" /*"dsu-velero"*/)
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

	fsPETM, err := fs.NewFSProtectedEntityTypeManagerFromConfig(fsParams, "notUsed", logrus.New())
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	fsPEs, err := fsPETM.GetProtectedEntities(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, fsPEID := range fsPEs {
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

func TestRetrieveEntity(t *testing.T) {
	s3petm, err := setupPETM(t, "ivd")
	if err != nil {
		t.Fatal(err)
	}

	ivdParams := make(map[string]interface{})
	ivdParams["vcHost"] = "10.161.99.58"
	ivdParams["insecureVC"] = "Y"
	ivdParams["vcUser"] = "administrator@vsphere.local"
	ivdParams["vcPassword"] = "Admin!23"

	ivdPETM, err := ivd.NewIVDProtectedEntityTypeManagerFromConfig(ivdParams, "notUsed", logrus.New())
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	peid := arachne.NewProtectedEntityIDWithSnapshotID("ivd", "ff9ac770-7ecd-405d-841b-6232857520d4",
		arachne.NewProtectedEntitySnapshotID("bde7e96d-8065-4bd5-a82a-edd7b2f540de"))
	s3PE, err := s3petm.GetProtectedEntity(ctx, peid)
	if err != nil {
		t.Fatal(err)
	}

	//mdr, err := s3PE.GetMetadataReader(ctx)
	/*
		mdr, err := s3PE.GetDataReader(ctx)

		if err != nil {
			t.Fatal(err)
		}
		op, err := os.Create("/home/dsmithuchida/tmp/xyzzy")
		if err != nil {
			t.Fatal(err)
		}

		bytesCopied, err := io.Copy(op, mdr)
		if err != nil {
			t.Fatal(err)
		}


		fmt.Printf("%d bytes copied\n", bytesCopied)
	*/
	newIVDPE, err := ivdPETM.Copy(ctx, s3PE, arachne.AllocateNewObject)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Restored new IVDPE %s\n", newIVDPE.GetID().String())

	/*
		dataReader, err := pe.GetDataReader(ctx)
		if err != nil {
			t.Fatal(err)
		}
		buf := make([]byte, 1024*1024)
		keepReading := true
		for keepReading {
			read, err := dataReader.Read(buf)
			if err != nil {
				keepReading = false
				if err != io.EOF {
					t.Fatal(err)
				}
			}
			fmt.Printf("Read %d bytes\n", read)
		}

	*/
}
func TestCopyIVDProtectedEntity(t *testing.T) {
	s3petm, err := setupPETM(t, "ivd")
	if err != nil {
		t.Fatal(err)
	}

	ivdParams := make(map[string]interface{})
	ivdParams["vcHost"] = "10.161.99.58"
	ivdParams["insecureVC"] = "Y"
	ivdParams["vcUser"] = "administrator@vsphere.local"
	ivdParams["vcPassword"] = "Admin!23"

	ivdPETM, err := ivd.NewIVDProtectedEntityTypeManagerFromConfig(ivdParams, "notUsed", logrus.New())
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	/*ivdPEs, err := ivdPETM.GetProtectedEntities(ctx)
	if err != nil {
		t.Fatal(err)
	}
	*/
	ctx = context.Background()

	//PESSID := arachne.NewProtectedEntitySnapshotID("ecb7fa78-cef9-4459-b898-17a39f582d9b")
	//ivdPEID := arachne.NewProtectedEntityIDWithSnapshotID("ivd", "cf29221a-381b-4036-825a-56bf8294ed38", ivdPESSID)
	ivdPEID := arachne.NewProtectedEntityID("ivd", "9d886896-f7f4-46d4-b6ab-f50b30013467")
	ivdPE, err := ivdPETM.GetProtectedEntity(ctx, ivdPEID)

	snapID, err := ivdPE.Snapshot(ctx)
	if err != nil {
		t.Fatal(err)
	}
	snapPEID := arachne.NewProtectedEntityIDWithSnapshotID("ivd", ivdPEID.GetID(), *snapID)
	snapPE, err := ivdPETM.GetProtectedEntity(ctx, snapPEID)
	if err != nil {
		t.Fatal(err)
	}
	s3PE, err := s3petm.Copy(ctx, snapPE, arachne.AllocateNewObject)
	if err != nil {
		t.Fatal(err)
	}

	newIVDPE, err := ivdPETM.Copy(ctx, s3PE, arachne.AllocateNewObject)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Restored new IVDPE %s\n", newIVDPE.GetID().String())

	/*
		for _, ivdPEID := range ivdPEs {
			ivdPE, err := ivdPETM.GetProtectedEntity(ctx, ivdPEID)
			if err != nil {
				t.Fatal(err)
			}
			snapID, err := ivdPE.Snapshot(ctx)
			if err == nil {

				snapPEID := arachne.NewProtectedEntityIDWithSnapshotID("ivd", ivdPEID.GetID(), *snapID)
				snapPE, err := ivdPETM.GetProtectedEntity(ctx, snapPEID)
				if err != nil {
					t.Fatal(err)
				}
				s3PE, err := s3petm.Copy(ctx, snapPE, arachne.AllocateNewObject)
				if err != nil {
					t.Fatal(err)
				}

				newIVDPE, err := ivdPETM.Copy(ctx, s3PE, arachne.AllocateNewObject)
				if err != nil {
					t.Fatal(err)
				}
				log.Printf("Restored new IVDPE %s\n", newIVDPE.GetID().String())
				status, err := ivdPE.DeleteSnapshot(ctx, *snapID)
				if err != nil {
					t.Fatal(err)
				}
				if !status {
					t.Fatal("Snapshot delete returned false")
				}
			} else {
				log.Printf("Snapshot failed for %s, skipping\n", ivdPEID.String())
			}
		}
	*/
}
