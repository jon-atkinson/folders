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
	if org.folders == nil {
		return []Folder{}
	}

	var folders []Folder

	org.folders.Ascend(func(node *FolderTreeNode) bool {
		stack := []*FolderTreeNode{node}

		for len(stack) > 0 {
			curr := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if curr.folder != nil {
				folders = append(folders, *curr.folder)
			}

			curr.children.Descend(func(child *FolderTreeNode) bool {
				stack = append(stack, child)
				return true
			})
		}
		return true
	})

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
				return []Folder{}, errors.New("Folder does not exist in the specified organization")
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

		curr.children.Descend(func(child *FolderTreeNode) bool {
			stack = append(stack, *child)
			return true
		})
	}

	return folders, nil
}

func (org Org) GetNamedFolder(name string) (*FolderTreeNode, error) {
	if org.folders == nil {
		return nil, fmt.Errorf("Org %s has no folders", org.orgId.String())
	}

	var targetNode *FolderTreeNode

	org.folders.Ascend(func(node *FolderTreeNode) bool {
		stack := []*FolderTreeNode{node}

		for len(stack) > 0 {
			curr := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if curr.folder != nil && curr.folder.Name == name {
				targetNode = curr
				return false
			}

			curr.children.Descend(func(child *FolderTreeNode) bool {
				stack = append(stack, child)
				return true
			})
		}
		return true
	})

	if targetNode == nil {
		return nil, errors.New("Folder does not exist")
	}
	return targetNode, nil
}

func (f *driver) GetAllFolders() ([]Folder, error) {
	var folders []Folder
	for _, org := range f.orgs {
		orgFolders, err := f.GetFoldersByOrgID(org.orgId)
		if err != nil {
			return []Folder{}, err
		}
		folders = append(folders, orgFolders...)
	}
	return folders, nil
}
