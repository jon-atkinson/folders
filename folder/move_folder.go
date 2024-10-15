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

	dstFolder := f.getDstFolder(dst)
	if dstFolder == nil {
		return []Folder{}, errors.New("Destination folder does not exist")
	}
	if slices.Contains(strings.Split(dstFolder.Paths, "."), name) {
		return []Folder{}, errors.New("Cannot move a folder to a child of itself")
	}

	targetExists := false
	for nodeIdx, node := range f.folders {
		paths := strings.Split(node.Paths, ".")

		for pathIdx := range paths {
			if paths[pathIdx] == name {
				if paths[len(paths)-1] == name {
					targetExists = true
				}

				if dstFolder.OrgId != node.OrgId {
					return []Folder{}, errors.New("Cannot move a folder to a different organization")
				}

				fixPaths(&f.folders[nodeIdx], dstFolder.Paths, strings.Join(paths[:pathIdx], "."))
			}
		}
	}

	if !targetExists {
		return []Folder{}, errors.New("Source folder does not exist")
	}

	return f.folders, nil
}

func (f *driver) getDstFolder(dst string) *Folder {
	for _, node := range f.folders {
		if node.Name == dst {
			return &node
		}
	}
	return nil
}

// updates the paths for all nodes in the tree rooted at node
// node rooted at newPrefix, children are updated as required
func fixPaths(node *Folder, newPrefix, oldPrefix string) {
	newPrefix = strings.Trim(newPrefix, ".")

	// preprocessing for well-formed top level nodes in output
	node.Paths = "." + node.Paths + "."
	node.Paths = strings.Replace(node.Paths, oldPrefix, newPrefix, 1)
	node.Paths = strings.Trim(node.Paths, ".")
}
