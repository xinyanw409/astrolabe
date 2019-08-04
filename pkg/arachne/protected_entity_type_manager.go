package arachne

import (
	//	"archive/zip"
	"context"
)

type ProtectedEntityTypeManager interface {
	GetTypeName() string
	GetProtectedEntity(ctx context.Context, id ProtectedEntityID) (ProtectedEntity, error)
	GetProtectedEntities(ctx context.Context) ([]ProtectedEntityID, error)
	Copy(ctx context.Context, pe ProtectedEntity) (ProtectedEntity, error)
	CopyFromInfo(ctx context.Context, info ProtectedEntityInfo) (ProtectedEntity, error)
}
