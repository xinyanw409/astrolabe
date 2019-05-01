package arachne

import (
//	"archive/zip"
)

type ProtectedEntityTypeManager interface {
   GetTypeName() string
   GetProtectedEntity(id ProtectedEntityID) ProtectedEntity
   GetProtectedEntities() [] ProtectedEntity
   //Serialize(pe ProtectedEntity, out Zip.Writer)
   //Deserialize(is ZipInputStream, ProtectedEntityInfo peInfo) ProtectedEntity
   //SerializeData(pe ProtectedEntity, out OutputStream)
}

