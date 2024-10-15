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

	for _, folder := range f.folders {
		if folder.OrgId == orgID {
			res = append(res, folder)
		}
	}
	return res
}

func (f *driver) GetAllChildFolders(orgID uuid.UUID, name string) []Folder {
	res := []Folder{}
	for _, node := range f.folders {
		paths := strings.Split(node.Paths, ".")
		if node.OrgId == orgID && slices.Contains(paths, name) {
			res = append(res, node)
		}
	}

	return res
}
