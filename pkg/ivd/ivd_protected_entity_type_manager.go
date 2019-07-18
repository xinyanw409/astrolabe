package ivd

import (
	"context"
	"github.com/pkg/errors"
	"github.com/vmware/arachne/pkg/core"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vslm"
	"net/url"
)

type IVDProtectedEntityTypeManager struct {
	client *govmomi.Client
	vsom   *vslm.GlobalObjectManager
}

func NewIVDProtectedEntityTypeManagerFromConfig(params map[string]interface{}) (*IVDProtectedEntityTypeManager, error) {
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
	return NewIVDProtectedEntityTypeManagerFromURL(&vcURL, insecure)
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

	vsom := vslm.NewGlobalObjectManager(vslmClient)

	retVal := IVDProtectedEntityTypeManager{
		client: client,
		vsom:   vsom,
	}
	return &retVal, nil
}

func (this *IVDProtectedEntityTypeManager) GetTypeName() string {
	return "ivd"
}

func (this *IVDProtectedEntityTypeManager) GetProtectedEntity(ctx context.Context, id core.ProtectedEntityID) (core.ProtectedEntity, error) {
	retIPE, err := newIVDProtectedEntity(this, id)
	if err != nil {
		return nil, err
	}
	return retIPE, nil
}

func (this *IVDProtectedEntityTypeManager) GetProtectedEntities(ctx context.Context) ([]core.ProtectedEntity, error) {
	res, err := this.vsom.ListObjectsForSpec(ctx, nil, 1000)
	if err != nil {
		return nil, err
	}
	var retEntities []core.ProtectedEntity
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
