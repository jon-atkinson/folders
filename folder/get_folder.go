package folder

import (
	"slices"
	"strings"

	"github.com/gofrs/uuid"
)

func GetAllFolders() []Folder {
	return GetSampleData()
}

func (f *driver) GetFoldersByOrgID(orgID uuid.UUID) []Folder {
	res := []Folder{}
	orgFound := false

	for _, folder := range f.folders {
		if folder.OrgId == orgID {
			orgFound = true
			res = append(res, folder)
		}
	}
	if !orgFound {
		return nil
	}
	return res
}

func (f *driver) GetAllChildFolders(orgID uuid.UUID, name string) []Folder {
	res := []Folder{}
	orgFound := false

	for _, node := range(f.folders) {
		paths := strings.Split(node.Paths, ".")
		if node.OrgId == orgID && slices.Contains(paths, name) && node.Name != name{
			orgFound = true
			res = append(res, node)
		}
	}

	if !orgFound {
		return nil
	}
	return res
}
