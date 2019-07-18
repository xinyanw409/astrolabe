package arachne

type ProtectedEntityManager struct {
	typeManager map[string]ProtectedEntityTypeManager
}

func NewProtectedEntityManager(petms []ProtectedEntityTypeManager) (returnPEM ProtectedEntityManager) {
	for _, curPETM := range petms {
		returnPEM.typeManager[curPETM.GetTypeName()] = curPETM
	}
	return returnPEM
}

func (pem *ProtectedEntityManager) getProtectedEntity(id ProtectedEntityID) ProtectedEntity {
      //return typeManagers.get(id.getPeType()).getProtectedEntity(id);
      return nil
   }
   
func (pem *ProtectedEntityManager) getProtectedEntityTypeManager(peType string) ProtectedEntityTypeManager {
      return pem.typeManager[peType]
   }
   
func (pem *ProtectedEntityManager) listEntityTypeManagers() []ProtectedEntityTypeManager {
	returnArr := []ProtectedEntityTypeManager{}
	for _, curPETM := range pem.typeManager { 
	 returnArr = append(returnArr, curPETM)
}
   return returnArr
}