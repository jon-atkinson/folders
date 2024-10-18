package folder

import (
	"errors"
	"slices"
	"strings"
)

// moves Folder name to be a child of Folder dst
// runs in O(n) in length of folders due to path updates
func (f *driver) MoveFolder(name string, dst string) ([]Folder, error) {
	if name == dst {
		return []Folder{}, errors.New("Cannot move a folder to itself")
	}

	fromFolder, found := f.folderMap[name]
	if !found {
		return []Folder{}, errors.New("Source folder does not exist")
	}

	toFolder, found := f.folderMap[dst]
	if !found {
		return []Folder{}, errors.New("Destination folder does not exist")
	}

	if fromFolder.folder.OrgId != toFolder.folder.OrgId {
		return []Folder{}, errors.New("Cannot move a folder to a different organization")
	}
	if slices.Contains(strings.Split(toFolder.folder.Paths, "."), fromFolder.folder.Name) {
		return []Folder{}, errors.New("Cannot move a folder to a child of itself")
	}

	// update position
	toFolder.children[fromFolder.folder.Name] = fromFolder
	if fromFolder.parent != nil {
		delete(fromFolder.parent.children, fromFolder.folder.Name)
		toFolder.children[fromFolder.folder.Name] = fromFolder
		fromFolder.parent = toFolder
	}

	// update paths
	fixPaths(fromFolder, toFolder.folder.Paths)

	return f.GetAllFolders(), nil
}

// updates the paths for all nodes in the tree rooted at node
// node rooted at newPrefix, children are updated as required
func fixPaths(node *FolderTreeNode, newPrefix string) {
	paths := strings.Split(node.folder.Paths, ".")
	oldPrefix := strings.Trim(strings.Join(paths[:len(paths)-1], "."), ".")
	newPrefix = strings.Trim(newPrefix, ".")

	stack := []*FolderTreeNode{node}
	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// preprocessing for well-formed top level nodes in output
		curr.folder.Paths = "." + curr.folder.Paths + "."
		curr.folder.Paths = strings.Replace(curr.folder.Paths, oldPrefix, newPrefix, 1)
		curr.folder.Paths = strings.Trim(curr.folder.Paths, ".")

		for _, child := range curr.children {
			stack = append(stack, child)
		}
	}
}
