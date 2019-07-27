package arachne

import "context"

type ProtectedEntityManager interface {
	GetProtectedEntity(ctx context.Context, id ProtectedEntityID) ProtectedEntity
	GetProtectedEntityTypeManager(peType string) ProtectedEntityTypeManager
	ListEntityTypeManagers() []ProtectedEntityTypeManager
}
