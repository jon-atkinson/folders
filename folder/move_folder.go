package folder

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
)

func (f *driver) MoveFolder(name string, dst string) ([]Folder, error) {
	if name == dst {
		return nil, errors.New("Error: Cannot move a folder to itself")
	}

	var from *FolderTreeNode
	var fromOrg *Org
	for _, checkOrg := range f.orgs {
		// ignore err as folder may be in a different Org
		from, _ = checkOrg.GetNamedFolder(name)
		if from != nil {
			fromOrg = &checkOrg
			break
		}
	}

	if from == nil {
		return nil, errors.New("Error: Source folder does not exist")
	}

	var to *FolderTreeNode
	var toOrg *Org
	for _, checkOrg := range f.orgs {
		// ignore err as folder may be in a different Org
		to, _ := checkOrg.GetNamedFolder(dst)
		if to != nil {
			toOrg = &checkOrg
			break
		}
	}

	if to == nil {
		return nil, errors.New("Error: Destination folder does not exist")
	}
	if fromOrg != toOrg {
		return nil, errors.New("Error: Cannot move a folder to a different organization")
	}
	if slices.Contains(strings.Split(to.folder.Paths, "."), from.folder.Name) {
		return nil, errors.New("Error: Cannot move a folder to a child of itself")
	}

	// todo, make sure ordering is preserved here
	to.children = append(to.children, from)
	// fmt.Println(*from, "\n\n\n")
	// fmt.Println(*to.children[0], "\n\n\n")
	fromOrg.PruneFolder(from)
	fixPaths(from, to.folder.Paths)

	allFolders, err := f.GetAllFolders()
	if err != nil {
		return nil, err
	}
	return allFolders, nil
}

func (org Org) PruneFolder(node *FolderTreeNode) {
	paths := strings.Split(node.folder.Paths, ".")
	if len(paths) == 0 {
		// bad path, we've got real problems if we're hitting this
		return
	}
	curr := lookupTreeNode(org.folders, paths[0])
	paths = paths[1:]

	// navigate branches
	for i, path := range paths {
		idx := sort.Search(len(curr.children), func(i int) bool {
			return curr.children[i].folder.Name == path
		})

		// for _, child := range curr.children {
		// 	fmt.Println(child.folder)
		// }

		// prune the correct branch
		if i == len(paths)-1 {
			fmt.Println(curr.folder.Name)
			fmt.Println(curr.children[idx].folder.Name)
			if len(curr.children) > idx+1 {
				curr.children = append(curr.children[:idx], curr.children[idx+1:]...)
			} else {
				curr.children = append(curr.children[:idx])
			}
			return
		}

		curr = curr.children[idx]
	}
	return
}

func fixPaths(node *FolderTreeNode, newPrefix string) {
	oldPrefix := node.folder.Paths
	stack := []FolderTreeNode{*node}

	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		strings.Replace(curr.folder.Paths, oldPrefix, newPrefix, 1)

		for i := len(curr.children) - 1; i >= 0; i-- {
			stack = append(stack, *curr.children[i])
		}
	}
}
