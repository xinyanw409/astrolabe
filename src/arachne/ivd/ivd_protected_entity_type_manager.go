package ivd

import (
	"arachne"
	"context"
	govmomi "github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vslm"
	"net/url"
)

type IVDProtectedEntityTypeManager struct {
	client *govmomi.Client
	vsom   *vslm.VslmObjectManager
}

func NewIVDProtectedEntityTypeManagerFromURL(url *url.URL, insecure bool) (*IVDProtectedEntityTypeManager, error) {
	ctx := context.Background()
	client, err := govmomi.NewClient(ctx, url, insecure)
	if err != nil {
		return nil, err
	}

	vslmClient, err := vslm.NewClient(ctx, client.Client)

	if err != nil {
		return nil, err
	}
	
	return NewIVDProtectedEntityTypeManagerWithClient(client, vslmClient)
}

func NewIVDProtectedEntityTypeManagerWithClient(client *govmomi.Client, vslmClient *vslm.Client) (*IVDProtectedEntityTypeManager, error) {

	vsom := vslm.NewVslmObjectManager(vslmClient)

	retVal := IVDProtectedEntityTypeManager{
		client: client,
		vsom:   vsom,
	}
	return &retVal, nil
}

func (ipetm *IVDProtectedEntityTypeManager) GetTypeName() string {
	return "ivd"
}

func (ipetm *IVDProtectedEntityTypeManager) GetProtectedEntity(id arachne.ProtectedEntityID) (arachne.ProtectedEntity, error) {
	return nil, nil
}

func (ipetm *IVDProtectedEntityTypeManager) GetProtectedEntities(ctx context.Context) ([]arachne.ProtectedEntity, error) {
	res, err := ipetm.vsom.ListVStorageObjectForSpec(ctx, nil, 1000)
	if (err != nil) {
		return nil, err
	}
	var retEntities []arachne.ProtectedEntity
	for _, curVSOID := range res.Id {
		arachneId := newProtectedEntityID(curVSOID)
		newIPE, err := newIVDProtectedEntity(ipetm, arachneId)
		if (err != nil) {
			return nil, err
		}
		retEntities = append(retEntities, &newIPE)
	}
	return retEntities, nil
}
