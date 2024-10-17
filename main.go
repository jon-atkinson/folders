package main

import (
	"fmt"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
)

const (
	FirstOrgID  = "c1556e17-b7c0-45a3-a6ae-9546248fb17a"
	SecondOrgID = "38b9879b-f73b-4b0e-b9d9-4fc4c23643a7"
)

func main() {
	orgID := uuid.FromStringOrNil(folder.DefaultOrgID)

	// res := folder.GetAllFolders()

	// example usage
	folderDriver := folder.NewDriver([]folder.Folder{
		{"bravo", orgID, "bravo"},
		{"alpha", orgID, "alpha"},
		{"charlie", orgID, "charlie"},
	})
	// orgFolder := folderDriver.GetFoldersByOrgID(orgID)
	folders, _ := folderDriver.MoveFolder("bravo", "alpha")
	folders, _ = folderDriver.MoveFolder("charlie", "bravo")

	fmt.Printf("\n Folders for orgID: %s", orgID)
	folder.PrettyPrint(folders)
}
