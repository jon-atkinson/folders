package folder

import (
	"errors"
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
}

func NewDriver(folders []Folder) IDriver {
	folderMap := make(map[string]*FolderTreeNode)

	return &driver{
		folderMap:   folderMap,
		folderTree:  buildFolderTree(&folders, &folderMap),
		folderSlice: &folders,
	}
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
func buildFolderTree(folders *[]Folder, folderMap *map[string]*FolderTreeNode) map[string]*FolderTreeNode {
	if len(*folders) == 0 {
		return make(map[string]*FolderTreeNode)
	}

	preProcessFolders(folders)
	folderTree := make(map[string]*FolderTreeNode)

	for i := range *folders {
		node := NewFolderTreeNode(&(*folders)[i])
		insertFolder(node, folderTree)
		(*folderMap)[(*folders)[i].Name] = node
	}

	return folderTree
}

// inserts folder into correct position folderTree
// navigates tree based on node.folder.Paths
func insertFolder(node *FolderTreeNode, parent map[string]*FolderTreeNode) error {
	parts := strings.Split(node.folder.Paths, ".")
	if len(parts) == 0 {
		return errors.New("Cannot insert folder with empty path")
	}

	curr, found := parent[parts[0]]

	if !found {
		if len(parts) == 1 {
			parent[node.folder.Name] = node
		}
		return nil
	}
	parts = parts[1:]

	for idx, part := range parts {
		next, found := curr.children[part]

		if !found {
			if idx != len(parts)-1 {
				// missing folders on path
				return errors.New("Could not insert, missing folders on path")
			}

			// insert folder to tree
			curr.children[node.folder.Name] = node
			return nil
		}

		curr = next
	}

	// folder already exists at this location
	return errors.New("Could not insert, folder already exists at this location")
}

// used to ensure unordered slices are ordered in the output to match tests that
// request it
func SortFoldersByPath(folders []Folder) []Folder {
	slices.SortFunc(folders, func(a, b Folder) int {
		return strings.Compare(a.Paths, b.Paths)
	})
	return folders
}
