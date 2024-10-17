package folder

import (
	"errors"
	"slices"
	"strings"
)

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

	_, err := f.pruneFolder(fromFolder)
	if err != nil {
		return []Folder{}, err
	}
	fixPaths(fromFolder, toFolder.folder.Paths)
	toFolder.children[fromFolder.folder.Name] = fromFolder

	return f.GetAllFolders(), nil
}

// removes target node from Org
// errors on: no/incorrect path, folder not in the org
func (f *driver) pruneFolder(node *FolderTreeNode) (*FolderTreeNode, error) {
	paths := strings.Split(node.folder.Paths, ".")
	if len(paths) == 0 {
		return nil, errors.New("Could not prune tree, requested path was empty")
	}

	curr, found := f.folderTree[paths[0]]
	if !found {
		return nil, errors.New("Could not prune tree, folder does not exist")
	}
	paths = paths[1:]

	for i, path := range paths {
		next, found := curr.children[path]
		if !found {
			return nil, errors.New("Could not prune tree, missing folders on path")
		}

		if i == len(paths)-1 {
			delete(curr.children, path)
			return next, nil
		}

		curr = next
	}

	// target is a top-level folder
	if _, found := f.folderTree[node.folder.Name]; found {
		res := f.folderTree[node.folder.Name]
		delete(f.folderTree, node.folder.Name)
		return res, nil
	}

	return nil, errors.New("Likely Bug: prune tree, this should be unreachable")
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
