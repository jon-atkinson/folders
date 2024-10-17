package folder

import (
	"github.com/gofrs/uuid"
)

func GetAllFolders() []Folder {
	return GetSampleData()
}

func (f *driver) GetFoldersByOrgID(orgID uuid.UUID) []Folder {
	var folders []Folder
	for _, folder := range f.folderTree {
		if folder.folder.OrgId == orgID {
			folders = append(folders, folder.collectFoldersInOrder()...)
		}
	}

	// I chose in-order traversal here, this function could be extended to
	// support different output orderings as required
	return folders
}

// returns folders in Org in-order
func (fol *FolderTreeNode) collectFoldersInOrder() []Folder {
	var folders []Folder
	stack := []*FolderTreeNode{fol}

	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if curr.folder != nil {
			folders = append(folders, *curr.folder)
		}

		for _, child := range curr.children {
			stack = append(stack, child)
		}
	}

	return folders
}

func (f *driver) GetAllChildFolders(orgID uuid.UUID, name string) []Folder {
	namedFolder, found := f.folderMap[name]
	if !found || namedFolder.folder.OrgId != orgID {
		return nil
	}

	var folders []Folder
	stack := []*FolderTreeNode{namedFolder}

	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if curr.folder != nil && curr.folder.Name != name {
			folders = append(folders, *curr.folder)
		}

		for _, child := range curr.children {
			stack = append(stack, child)
		}
	}

	return folders
}

// returns all folders on f
// folders are collected for earch Org by seperate goroutines
func (f *driver) GetAllFolders() []Folder {
	return *f.folderSlice
}
