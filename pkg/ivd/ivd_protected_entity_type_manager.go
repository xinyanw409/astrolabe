package ivd

import (
	"context"
	"github.com/pkg/errors"
	"github.com/vmware/arachne/pkg/arachne"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/vslm"
	types2 "github.com/vmware/govmomi/vslm/types"
	"github.com/vmware/gvddk/gDiskLib"
	"net/url"
	"time"
)

type IVDProtectedEntityTypeManager struct {
	client    *govmomi.Client
	vsom      *vslm.GlobalObjectManager
	s3URLBase string
	user      string	// These are being kept so we can open VDDK connections, may be able to open a VDDK connection
	password  string	// in IVDProtectedEntityTypeManager instead
}

func NewIVDProtectedEntityTypeManagerFromConfig(params map[string]interface{}, s3URLBase string) (*IVDProtectedEntityTypeManager, error) {
	var vcURL url.URL
	vcHostStr, ok := params["vcHost"].(string)
	if !ok {
		return nil, errors.New("Missing vcHost param, cannot initialize IVDProtectedEntityTypeManager")
	}
	vcURL.Scheme = "https"
	vcURL.Host = vcHostStr
	insecure := false
	insecureStr, ok := params["insecureVC"].(string)
	if ok && (insecureStr == "Y" || insecureStr == "y") {
		insecure = true
	}
	vcUser, ok := params["vcUser"].(string)
	if !ok {
		return nil, errors.New("Missing vcUser param, cannot initialize IVDProtectedEntityTypeManager")
	}
	vcPassword, ok := params["vcPassword"].(string)
	if !ok {
		return nil, errors.New("Missing vcPassword param, cannot initialize IVDProtectedEntityTypeManager")
	}
	vcURL.User = url.UserPassword(vcUser, vcPassword)
	vcURL.Path = "/sdk"
	return NewIVDProtectedEntityTypeManagerFromURL(&vcURL, s3URLBase, insecure)
}

func NewIVDProtectedEntityTypeManagerFromURL(url *url.URL, s3URLBase string, insecure bool) (*IVDProtectedEntityTypeManager, error) {
	ctx := context.Background()
	client, err := govmomi.NewClient(ctx, url, insecure)
	if err != nil {
		return nil, err
	}

	vslmClient, err := vslm.NewClient(ctx, client.Client)

	if err != nil {
		return nil, err
	}

	retVal, err := newIVDProtectedEntityTypeManagerWithClient(client, s3URLBase, vslmClient)
	if err == nil {

	retVal.user = url.User.Username()
	password, hasPassword := url.User.Password()
	if !hasPassword {
		return nil, errors.New("No VC Password specified")
	}
	retVal.password = password
}
	return retVal, err
}

const vsphereMajor = 6
const vSphereMinor = 7
const disklibLib64 = "/usr/lib/vmware-vix-disklib/lib64"
func newIVDProtectedEntityTypeManagerWithClient(client *govmomi.Client, s3URLBase string, vslmClient *vslm.Client) (*IVDProtectedEntityTypeManager, error) {

	vsom := vslm.NewGlobalObjectManager(vslmClient)

	errno := gDiskLib.Init(vsphereMajor, vSphereMinor, disklibLib64)
	if errno != 0 {
		return nil, errors.New("Could not initialize VDDK errno ")
	}
	retVal := IVDProtectedEntityTypeManager{
		client:    client,
		vsom:      vsom,
		s3URLBase: s3URLBase,
	}
	return &retVal, nil
}

func (this *IVDProtectedEntityTypeManager) GetTypeName() string {
	return "ivd"
}

func (this *IVDProtectedEntityTypeManager) GetProtectedEntity(ctx context.Context, id arachne.ProtectedEntityID) (arachne.ProtectedEntity, error) {
	retIPE, err := newIVDProtectedEntity(this, id)
	if err != nil {
		return nil, err
	}
	return retIPE, nil
}

func (this *IVDProtectedEntityTypeManager) GetProtectedEntities(ctx context.Context) ([]arachne.ProtectedEntityID, error) {
	// Kludge because of PR
	spec := types2.VslmVsoVStorageObjectQuerySpec{
		QueryField:    "createTime",
		QueryOperator: "greaterThan",
		QueryValue:    []string{"0"},
	}
	res, err := this.vsom.ListObjectsForSpec(ctx, []types2.VslmVsoVStorageObjectQuerySpec{spec}, 1000)
	if err != nil {
		return nil, err
	}
	retIDs := make([]arachne.ProtectedEntityID, len(res.Id))
	for idNum, curVSOID := range res.Id {
		arachneId := newProtectedEntityID(curVSOID)
		retIDs[idNum] = arachneId
	}
	return retIDs, nil
}

func (this *IVDProtectedEntityTypeManager) Copy(ctx context.Context, pe arachne.ProtectedEntity, options arachne.CopyCreateOptions) (arachne.ProtectedEntity, error) {
	info, err := pe.GetInfo(ctx)
	if err != nil {
		return nil, err
	}
	return this.CopyFromInfo(ctx, info, options)
}

func (this *IVDProtectedEntityTypeManager) CopyFromInfo(ctx context.Context, peInfo arachne.ProtectedEntityInfo, options arachne.CopyCreateOptions) (arachne.ProtectedEntity, error) {

	if (peInfo.GetID().GetPeType() != "ivd") {
		return nil, errors.New("Copy source must be an ivd")
	}
	ourVC := false
	existsInOurVC := false
	for _, checkData := range peInfo.GetDataTransports() {
		vcenterURL, ok := checkData.GetParam("vcenter")

		if checkData.GetTransportType() == "vadp" && ok && vcenterURL == this.client.URL().Host {
			ourVC = true
			existsInOurVC = true
			break
		}
	}

	if (ourVC) {
	_, err := this.vsom.Retrieve(ctx, NewVimIDFromPEID(peInfo.GetID()))
	if err != nil {
		if soap.IsSoapFault(err) {
			fault := soap.ToSoapFault(err).Detail.Fault
			if _, ok := fault.(types.NotFound); ok {
				// Doesn't exist in our local system, we can't just clone it
				existsInOurVC = false
			} else {
				return nil, err
			}
		}
	}
	}
	var retPE arachne.ProtectedEntity
	retPE = nil
	if ourVC && existsInOurVC {
		var createTask *vslm.Task
		var err error
		if (peInfo.GetID().HasSnapshot()) {
			createTask, err = this.vsom.CreateDiskFromSnapshot(ctx, NewVimIDFromPEID(peInfo.GetID()), NewVimSnapshotIDFromPEID(peInfo.GetID()),
				peInfo.GetName(),  nil, nil, "")
		} else {
			keepAfterDeleteVm := true
			cloneSpec := types.VslmCloneSpec{
				Name: "",
				KeepAfterDeleteVm: &keepAfterDeleteVm,
			}
			createTask, err = this.vsom.Clone(ctx, NewVimIDFromPEID(peInfo.GetID()), cloneSpec)
		}
		retVal, err := createTask.WaitNonDefault(ctx, time.Hour * 24, time.Second * 10, true, time.Second * 30);
		if err != nil {
			return nil, err
		}
		newVSO := retVal.(types.VStorageObject)
		retPE, err = newIVDProtectedEntity(this, newProtectedEntityID(newVSO.Config.Id))
		if err != nil {
			return nil, err
		}
	} else {

	}
	return retPE, nil
}

func (this *IVDProtectedEntityTypeManager) getDataTransports(id arachne.ProtectedEntityID) ([]arachne.DataTransport,
	[]arachne.DataTransport,
	[]arachne.DataTransport, error) {
	vadpParams := make(map[string]string)
	vadpParams["id"] = id.GetID()
	if id.GetSnapshotID().String() != "" {
		vadpParams["snapshotID"] = id.GetSnapshotID().String()
	}
	vadpParams["vcenter"] = this.client.URL().Host

	dataS3URL := this.s3URLBase + "ivd/" + id.String()
	data := []arachne.DataTransport{
		arachne.NewDataTransport("vadp", vadpParams),
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
