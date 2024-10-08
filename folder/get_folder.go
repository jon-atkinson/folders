package folder

import (
	"errors"
	"fmt"
	"sync"

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
	org.mux.Lock()
	defer org.mux.Unlock()

	// I chose in-order traversal here, this function could be extended to
	// support different output orderings as required
	return org.collectFoldersInOrder(), nil
}

func (org *Org) collectFoldersInOrder() []Folder {
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

	otherOrg, folder, err := f.nameToOrgFolder(name)
	if err != nil {
		return []Folder{}, err
	}

	if otherOrg.orgId != orgID {
		return []Folder{}, errors.New("Folder does not exist in the specified organization")
	}

	var folders []Folder
	stack := []*FolderTreeNode{folder}

	org.mux.Lock()
	defer org.mux.Unlock()

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

	return folders, nil
}

func (org *Org) GetNamedFolder(name string) (*FolderTreeNode, error) {
	org.mux.RLock()
	defer org.mux.RUnlock()

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
	resultChan := make(chan []Folder)
	errChan := make(chan error)
	var wg sync.WaitGroup

	for i := range f.orgs {
		concOrg := f.orgs[i]
		wg.Add(1)
		go func(org *Org) {
			org.mux.RLock()
			defer org.mux.RUnlock()
			defer wg.Done()

			folders, err := f.GetFoldersByOrgID(org.orgId)
			if err != nil {
				select {
				case errChan <- err:
				}
				return
			}
			resultChan <- folders
		}(concOrg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var folders []Folder

	for {
		select {
		case orgFolders, ok := <-resultChan:
			if ok {
				folders = append(folders, orgFolders...)
			} else {
				return folders, nil
			}
		case err := <-errChan:
			return nil, err
		}
	}
}
