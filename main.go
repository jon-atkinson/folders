package main

import (
	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
)

func main() {
	orgID := uuid.FromStringOrNil(folder.DefaultOrgID)
	uuid.FromStringOrNil(folder.DefaultOrgID)
	res := folder.GetAllFolders()

	// example usage
	folderDriver := folder.NewDriver(res)

	folders, _ := folderDriver.GetAllChildFolders(orgID, "safe-infragirl")
	for _, f := range folders {
		folder.PrettyPrint(f)
	}
}
