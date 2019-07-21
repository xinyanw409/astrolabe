package ivd

import (
	"context"
	"github.com/pkg/errors"
	"github.com/vmware/arachne/pkg/arachne"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vslm"
	"net/url"
)

type IVDProtectedEntityTypeManager struct {
	client    *govmomi.Client
	vsom      *vslm.GlobalObjectManager
	s3URLBase string
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

	return NewIVDProtectedEntityTypeManagerWithClient(client, s3URLBase, vslmClient)
}

func NewIVDProtectedEntityTypeManagerWithClient(client *govmomi.Client, s3URLBase string, vslmClient *vslm.Client) (*IVDProtectedEntityTypeManager, error) {

	vsom := vslm.NewGlobalObjectManager(vslmClient)

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

func (this *IVDProtectedEntityTypeManager) GetProtectedEntities(ctx context.Context) ([]arachne.ProtectedEntity, error) {
	res, err := this.vsom.ListObjectsForSpec(ctx, nil, 1000)
	if err != nil {
		return nil, err
	}
	var retEntities []arachne.ProtectedEntity
	for _, curVSOID := range res.Id {
		arachneId := newProtectedEntityID(curVSOID)
		newIPE, err := newIVDProtectedEntity(this, arachneId)
		if err != nil {
			return nil, err
		}
		retEntities = append(retEntities, &newIPE)
	}
	return retEntities, nil
}

func (this *IVDProtectedEntityTypeManager) Copy(ctx context.Context, pe arachne.ProtectedEntity) error {
	return nil
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
		arachne.NewDataTransportForS3(dataS3URL),
	}

	mdS3URL := dataS3URL + ".md"

	md := []arachne.DataTransport{
		arachne.NewDataTransportForS3(mdS3URL),
	}

	combinedS3URL := dataS3URL + ".zip"
	combined := []arachne.DataTransport{
		arachne.NewDataTransportForS3(combinedS3URL),
	}

	return data, md, combined, nil
}
