package kubernetes

import (
	"context"
	"github.com/vmware/arachne/pkg/core"
	v1 "k8s.io/api/core/v1"
)

type KubernetesNamespaceProtectedEntity struct {
	knpetm    *KubernetesNamespaceProtectedEntityTypeManager
	id        core.ProtectedEntityID
	namespace *v1.Namespace
}

func NewKubernetesNamespaceProtectedEntity(knpetm *KubernetesNamespaceProtectedEntityTypeManager,
	namespace *v1.Namespace) (*KubernetesNamespaceProtectedEntity, error) {
	nsPEID := core.NewProtectedEntityID("k8sns", namespace.Name)
	returnPE := KubernetesNamespaceProtectedEntity{
		knpetm:    knpetm,
		id:        nsPEID,
		namespace: namespace,
	}
	return &returnPE, nil
}

func (this *KubernetesNamespaceProtectedEntity) GetInfo(ctx context.Context) (core.ProtectedEntityInfo, error) {
	return nil, nil
}
func (this *KubernetesNamespaceProtectedEntity) GetCombinedInfo(ctx context.Context) ([]core.ProtectedEntityInfo, error) {
	return nil, nil

}

func (this *KubernetesNamespaceProtectedEntity) Snapshot(ctx context.Context) (*core.ProtectedEntitySnapshotID, error) {
	return nil, nil

}
func (this *KubernetesNamespaceProtectedEntity) ListSnapshots(ctx context.Context) ([]core.ProtectedEntitySnapshotID, error) {
	return nil, nil

}
func (this *KubernetesNamespaceProtectedEntity) DeleteSnapshot(ctx context.Context,
	snapshotToDelete core.ProtectedEntitySnapshotID) (bool, error) {
	return false, nil

}
func (this *KubernetesNamespaceProtectedEntity) GetInfoForSnapshot(ctx context.Context,
	snapshotID core.ProtectedEntitySnapshotID) (*core.ProtectedEntityInfo, error) {
	return nil, nil

}

func (this *KubernetesNamespaceProtectedEntity) GetComponents(ctx context.Context) ([]core.ProtectedEntity, error) {
	return nil, nil

}

func (this *KubernetesNamespaceProtectedEntity) GetID() core.ProtectedEntityID {
	return this.id
}
