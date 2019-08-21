package ivd

import "C"
import (
	"github.com/pkg/errors"
	"github.com/vmware/arachne/pkg/arachne"
	vim "github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/gvddk/gDiskLib"
	"io"
	"strings"

	//	"github.com/vmware/govmomi/vslm"
	"context"
	"time"
)

type IVDProtectedEntity struct {
	ipetm    *IVDProtectedEntityTypeManager
	id       arachne.ProtectedEntityID
	data     []arachne.DataTransport
	metadata []arachne.DataTransport
	combined []arachne.DataTransport
}

func (this IVDProtectedEntity) GetDataReader() (io.Reader, error) {

	url := this.ipetm.client.URL()
	serverName := url.Hostname()
	userName := this.ipetm.user
	password := this.ipetm.password
	/*
	thumbprint := this.ipetm.client.Thumbprint(serverName)
	thumbprint = "3D:62:45:37:88:36:3E:03:7A:6C:5A:63:D6:D6:AB:85:F7:DE:A3:AB"
	if thumbprint == "" {
		return nil, errors.New("Thumbprint was not set in client")
	}*/
	fcdid := this.id.GetID()

	vso, err := this.ipetm.vsom.Retrieve(context.Background(), NewVimIDFromPEID(this.id))
	if err != nil {
		return nil, err
	}
	datastore := vso.Config.Backing.GetBaseConfigInfoBackingInfo().Datastore.String()
	datastore = strings.TrimPrefix(datastore, "Datastore:")
	/*
	params := gDiskLib.ConnectParams{
		ServerName: serverName,
		UserName: userName,
		Password: password,
		ThumbPrint: thumbprint,
		FCDid: fcdid,
	}*/
	params := gDiskLib.NewConnectParams("",
		serverName,
		"3D:62:45:37:88:36:3E:03:7A:6C:5A:63:D6:D6:AB:85:F7:DE:A3:AB",
		userName,
		password,
		fcdid,
		datastore,
		"",
		"",
		"vm-example")


	conn, errno := gDiskLib.Connect(params)
	if errno != 0 {
		return nil, errors.New("Connect failed")
	}
	errno = gDiskLib.PrepareForAccess(params)
	if errno != 0 {
		return nil, errors.New("PrepareForAccess failed")
	}
	diskHandle, errno := gDiskLib.Open(conn, "", 1 /*C.VIXDISKLIB_FLAG_OPEN_UNBUFFERED*/)
	return arachne.NewReaderAtReader(diskHandle), nil
}

func (this IVDProtectedEntity) GetMetadataReader() (io.Reader, error) {
	return nil, nil
}

func newProtectedEntityID(id vim.ID) arachne.ProtectedEntityID {
	return arachne.NewProtectedEntityID("ivd", id.Id)
}

func newProtectedEntityIDWithSnapshotID(id vim.ID, snapshotID arachne.ProtectedEntitySnapshotID) arachne.ProtectedEntityID {
	return arachne.NewProtectedEntityIDWithSnapshotID("ivd", id.Id, snapshotID)
}

func newIVDProtectedEntity(ipetm *IVDProtectedEntityTypeManager, id arachne.ProtectedEntityID) (IVDProtectedEntity, error) {
	data, metadata, combined, err := ipetm.getDataTransports(id)
	if err != nil {
		return IVDProtectedEntity{}, err
	}
	newIPE := IVDProtectedEntity{
		ipetm:    ipetm,
		id:       id,
		data:     data,
		metadata: metadata,
		combined: combined,
	}
	return newIPE, nil
}
func (this IVDProtectedEntity) GetInfo(ctx context.Context) (arachne.ProtectedEntityInfo, error) {
	vsoID := vim.ID{
		Id: this.id.GetID(),
	}
	vso, err := this.ipetm.vsom.Retrieve(ctx, vsoID)
	if err != nil {
		return nil, errors.Wrap(err, "Retrieve failed")
	}

	retVal := arachne.NewProtectedEntityInfo(
		this.id,
		vso.Config.Name,
		this.data,
		this.metadata,
		this.combined,
		[]arachne.ProtectedEntityID{})
	return retVal, nil
}

func (this IVDProtectedEntity) GetCombinedInfo(ctx context.Context) ([]arachne.ProtectedEntityInfo, error) {
	ivdIPE, err := this.GetInfo(ctx)
	if err != nil {
		return nil, err
	}
	return []arachne.ProtectedEntityInfo{ivdIPE}, nil
}

const waitTime = 3600 * time.Second
/*
 * Snapshot APIs
 */
func (this IVDProtectedEntity) Snapshot(ctx context.Context) (*arachne.ProtectedEntitySnapshotID, error) {
	vslmTask, err := this.ipetm.vsom.CreateSnapshot(ctx, NewVimIDFromPEID(this.GetID()), "ArachneSnapshot")
	if err != nil {
		return nil, errors.Wrap(err, "Snapshot failed")
	}
	ivdSnapshotIDAny, err := vslmTask.Wait(ctx, waitTime)
	if err != nil {
		return nil, errors.Wrap(err, "Wait failed")
	}
	ivdSnapshotID := ivdSnapshotIDAny.(vim.ID)
	/*
		ivdSnapshotStr := ivdSnapshotIDAny.(string)
		ivdSnapshotID := vim.ID{
			id: ivdSnapshotStr,
		}
	*/
	retVal := arachne.NewProtectedEntitySnapshotID(ivdSnapshotID.Id)
	return &retVal, nil
}

func (this IVDProtectedEntity) ListSnapshots(ctx context.Context) ([]arachne.ProtectedEntitySnapshotID, error) {
	snapshotInfo, err := this.ipetm.vsom.RetrieveSnapshotInfo(ctx, NewVimIDFromPEID(this.GetID()))
	if err != nil {
		return nil, errors.Wrap(err, "RetrieveSnapshotInfo failed")
	}
	peSnapshotIDs := []arachne.ProtectedEntitySnapshotID{}
	for _, curSnapshotInfo := range snapshotInfo {
		peSnapshotIDs = append(peSnapshotIDs, arachne.NewProtectedEntitySnapshotID(curSnapshotInfo.Id.Id))
	}
	return peSnapshotIDs, nil
}
func (this IVDProtectedEntity) DeleteSnapshot(ctx context.Context, snapshotToDelete arachne.ProtectedEntitySnapshotID) (bool, error) {
	vslmTask, err := this.ipetm.vsom.DeleteSnapshot(ctx, NewVimIDFromPEID(this.id), NewVimSnapshotIDFromPEID(this.id))
	if err != nil {
		return false, errors.Wrap(err, "DeleteSnapshot failed")
	}
	_, err = vslmTask.Wait(ctx, waitTime)
	if err != nil {
		return false, errors.Wrap(err, "Wait failed")
	}
	return true, nil
}

func (this IVDProtectedEntity) GetInfoForSnapshot(ctx context.Context, snapshotID arachne.ProtectedEntitySnapshotID) (*arachne.ProtectedEntityInfo, error) {
	return nil, nil
}

func (this IVDProtectedEntity) GetComponents(ctx context.Context) ([]arachne.ProtectedEntity, error) {
	return make([]arachne.ProtectedEntity, 0), nil
}

func (this IVDProtectedEntity) GetID() arachne.ProtectedEntityID {
	return this.id
}

func NewIDFromString(idStr string) vim.ID {
	return vim.ID{
		Id: idStr,
	}
}

func NewVimIDFromPEID(peid arachne.ProtectedEntityID) vim.ID {
	if peid.GetPeType() == "ivd" {
		return vim.ID{
			Id: peid.GetID(),
		}
	} else {
		return vim.ID{}
	}
}

func NewVimSnapshotIDFromPEID(peid arachne.ProtectedEntityID) vim.ID {
	if peid.HasSnapshot() {
		return vim.ID{
			Id: peid.GetSnapshotID().GetID(),
		}
	} else {
		return vim.ID{}
	}
}
