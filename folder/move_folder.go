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

	fromOrg, fromFolder, err := f.nameToOrgFolder(name)
	if fromFolder == nil {
		return []Folder{}, errors.New("Source folder does not exist")
	}

	toOrg, toFolder, err := f.nameToOrgFolder(dst)
	if toFolder == nil {
		return []Folder{}, errors.New("Destination folder does not exist")
	}

	if fromOrg.orgId != toOrg.orgId {
		return []Folder{}, errors.New("Cannot move a folder to a different organization")
	}
	if slices.Contains(strings.Split(toFolder.folder.Paths, "."), fromFolder.folder.Name) {
		return []Folder{}, errors.New("Cannot move a folder to a child of itself")
	}

	_, err = fromOrg.pruneFolder(fromFolder)
	if err != nil {
		return []Folder{}, err
	}
	fixPaths(fromFolder, toFolder.folder.Paths)
	toOrg.insertFolder(fromFolder)

	allFolders, err := f.GetAllFolders()
	if err != nil {
		return []Folder{}, err
	}
	return allFolders, nil
}

func (org *Org) pruneFolder(node *FolderTreeNode) (*FolderTreeNode, error) {
	paths := strings.Split(node.folder.Paths, ".")
	if len(paths) == 0 {
		return nil, errors.New("Could not prune tree, requested path was empty")
	}

	curr, found := lookupTreeNode(org.folders, paths[0])
	if !found {
		return nil, errors.New("Could not prune tree, folder not in this organization")
	}
	paths = paths[1:]

	for i, path := range paths {
		next, found := curr.children.Get(&FolderTreeNode{folder: &Folder{Name: path}})
		if !found {
			return nil, errors.New("Could not prune tree, missing folders on path")
		}

		if i == len(paths)-1 {
			res, found := curr.children.Delete(next)
			if !found {
				return nil, errors.New("Could not prune tree, folder not in tree")
			}
			return res, nil
		}

		curr = next
	}

	// target is a top-level folder
	res, found := org.folders.Get(node)
	if !found {
		return nil, errors.New("Likely Bug: prune tree, this should be unreachable")
	}
	res, found = org.folders.Delete(res)
	if !found {
		return nil, errors.New("Likely Bug: prune tree, this should be unreachable")
	}
	return res, nil
}

func fixPaths(node *FolderTreeNode, newPrefix string) {
	paths := strings.Split(node.folder.Paths, ".")
	oldPrefix := strings.Trim(strings.Join(paths[:len(paths)-1], "."), ".")
	newPrefix = strings.Trim(newPrefix, ".")

	stack := []*FolderTreeNode{node}
	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// preprocessing for well-formed top level nodes
		curr.folder.Paths = "." + curr.folder.Paths + "."
		curr.folder.Paths = strings.Replace(curr.folder.Paths, oldPrefix, newPrefix, 1)
		curr.folder.Paths = strings.Trim(curr.folder.Paths, ".")

		curr.children.Ascend(func(child *FolderTreeNode) bool {
			stack = append(stack, child)
			return true
		})
	}
}
