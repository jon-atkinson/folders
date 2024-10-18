package folder

import (
	"slices"
	"strings"

	"github.com/gofrs/uuid"
)

type IDriver interface {
	// GetFoldersByOrgID returns all folders that belong to a specific orgID.
	GetFoldersByOrgID(orgID uuid.UUID) []Folder
	// component 1
	// Implement the following methods:
	// GetAllChildFolders returns all child folders of a specific folder.
	GetAllChildFolders(orgID uuid.UUID, name string) []Folder

	// component 2
	// Implement the following methods:
	// MoveFolder moves a folder to a new destination.
	MoveFolder(name string, dst string) ([]Folder, error)
}

type driver struct {
	folderMap   map[string]*FolderTreeNode
	folderTree  map[string]*FolderTreeNode
	folderSlice *[]Folder
}

type FolderTreeNode struct {
	folder   *Folder
	children map[string]*FolderTreeNode
	parent   *FolderTreeNode
}

func NewDriver(folders []Folder) IDriver {
	f := &driver{
		folderMap:   make(map[string]*FolderTreeNode, len(folders)),
		folderTree:  make(map[string]*FolderTreeNode, len(folders)),
		folderSlice: &folders,
	}
	buildFolderTree(&folders, &f.folderTree, &f.folderMap)
	return f
}

func NewFolderTreeNode(folder *Folder) *FolderTreeNode {
	return &FolderTreeNode{
		folder:   folder,
		children: make(map[string]*FolderTreeNode),
	}
}

// sort folders in place on Paths
func preProcessFolders(folders *[]Folder) {
	slices.SortFunc(*folders, func(a, b Folder) int {
		return strings.Compare(a.Paths, b.Paths)
	})
}

// Builds the folderTree, inserting each node into the global name lookup map
// Assumes well-formed folder trees in input which are properly seperated by OrgId
func buildFolderTree(folders *[]Folder, folderTree, folderMap *map[string]*FolderTreeNode) {
	if len(*folders) == 0 {
		return
	}

	preProcessFolders(folders)

	// assumes folders sorted by path
	for i := range *folders {
		node := NewFolderTreeNode(&(*folders)[i])

		// assumes all folders have a valid path
		paths := strings.Split(node.folder.Paths, ".")
		if len(paths) == 1 {
			(*folderTree)[(*folders)[i].Name] = node
		} else {
			(*folderMap)[paths[len(paths)-2]].children[(*folders)[i].Name] = node
			node.parent = (*folderMap)[paths[len(paths)-2]]
		}

		(*folderMap)[(*folders)[i].Name] = node
	}

	return
}

// used to ensure unordered slices are ordered in the output to match tests that
// request it
func SortFoldersByPath(folders []Folder) []Folder {
	slices.SortFunc(folders, func(a, b Folder) int {
		return strings.Compare(a.Paths, b.Paths)
	})
	return folders
}
