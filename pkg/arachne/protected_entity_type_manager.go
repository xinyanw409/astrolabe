package arachne

import (
	//	"archive/zip"
	"context"
)

type CopyCreateOptions int

const (
	AllocateNewObject    CopyCreateOptions = 1
	UpdateExistingObject CopyCreateOptions = 2
	AllocateObjectWithID CopyCreateOptions = 3
)

type ProtectedEntityTypeManager interface {
	GetTypeName() string
	GetProtectedEntity(ctx context.Context, id ProtectedEntityID) (ProtectedEntity, error)
	GetProtectedEntities(ctx context.Context) ([]ProtectedEntityID, error)
	Copy(ctx context.Context, pe ProtectedEntity, options CopyCreateOptions) (ProtectedEntity, error)
	CopyFromInfo(ctx context.Context, info ProtectedEntityInfo, options CopyCreateOptions) (ProtectedEntity, error)
}
