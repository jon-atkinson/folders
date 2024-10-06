package folder

import (
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

func GetAllFolders() []Folder {
	return GetSampleData()
}

func (f *driver) GetFoldersByOrgID(orgID uuid.UUID) ([]Folder, error) {
	org, err := f.getOrg(orgID)
	if err != nil {
		return []Folder{}, err
	}

	// I chose in-order traversal here, this function could be extended to
	// support different output orderings as required
	return org.collectFoldersInOrder(), nil
}

func (org Org) collectFoldersInOrder() []Folder {
	if org.folders == nil || len(org.folders) == 0 {
		return []Folder{}
	}

	var folders []Folder
	for _, tree := range org.folders {
		stack := []FolderTreeNode{*tree}

		for len(stack) > 0 {
			curr := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if curr.folder != nil {
				folders = append(folders, *curr.folder)
			}

			for i := len(curr.children) - 1; i >= 0; i-- {
				stack = append(stack, *curr.children[i])
			}
		}
	}

	return folders
}

func (f *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error) {
	org, err := f.getOrg(orgID)
	if err != nil {
		return []Folder{}, err
	}

	target, err := org.GetNamedFolder(name)
	if err != nil {
		for _, org := range f.orgs {
			// check if folder belongs to another Org
			if res, _ := org.GetNamedFolder(name); res != nil {
				return []Folder{}, errors.New("Error: Folder does not exist in the specified organization")
			}
		}
		return []Folder{}, err
	}

	var folders []Folder
	stack := []FolderTreeNode{*target}

	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if curr.folder != nil {
			folders = append(folders, *curr.folder)
		}

		for i := len(curr.children) - 1; i >= 0; i-- {
			stack = append(stack, *curr.children[i])
		}
	}

	return folders, nil
}

func (org Org) GetNamedFolder(name string) (*FolderTreeNode, error) {
	if org.folders == nil || len(org.folders) == 0 {
		return nil, fmt.Errorf("Org %s has no folders", org.orgId.String())
	}

	for _, tree := range org.folders {
		stack := []FolderTreeNode{*tree}

		for len(stack) > 0 {
			curr := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if curr.folder != nil && curr.folder.Name == name {
				return &curr, nil
			}

			for i := len(curr.children) - 1; i >= 0; i-- {
				stack = append(stack, *curr.children[i])
			}
		}
	}

	return nil, errors.New("Error: Folder does not exist")
}

func (f *driver) GetAllFolders() ([]Folder, error) {
	var folders []Folder
	for _, org := range f.orgs {
		orgFolders, err := f.GetFoldersByOrgID(org.orgId)
		if err != nil {
			return nil, err
		}
		folders = append(folders, orgFolders...)
	}
	return folders, nil
}
