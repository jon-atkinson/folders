package folder

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/google/btree"
)

type IDriver interface {
	// I took the liberty of adding an error to the return of these functions
	// as the spec requests specific error case handling

	// This is just nice to have access to, if having it here is an issue, it
	// can be made private instead
	GetAllFolders() ([]Folder, error)

	// GetFoldersByOrgID returns all folders that belong to a specific orgID.
	GetFoldersByOrgID(orgID uuid.UUID) ([]Folder, error)

	// component 1
	// Implement the following methods:
	// GetAllChildFolders returns all child folders of a specific folder.
	GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error)

	// component 2
	// Implement the following methods:
	// MoveFolder moves a folder to a new destination.
	MoveFolder(name string, dst string) ([]Folder, error)
	// From context I'm assuming that organizations share a file system and
	// cannot have identical names? I've assumed this to be true since otherwise
	// I'm not sure how to resolve namespace collisions
}

type FolderTreeNode struct {
	folder   *Folder
	children *btree.BTreeG[*FolderTreeNode]
}

type Org struct {
	orgId   uuid.UUID
	folders *btree.BTreeG[*FolderTreeNode]
}

type driver struct {
	orgs []Org
}

func NewDriver(folders []Folder) IDriver {
	return &driver{
		orgs: buildOrgs(folders),
	}
}

func NewFolderTreeNode(folder *Folder) *FolderTreeNode {
	return &FolderTreeNode{
		folder:   folder,
		children: btree.NewG(3, folderTreeLess),
	}
}

func NewOrg(orgID uuid.UUID) Org {
	return Org{
		orgId:   orgID,
		folders: btree.NewG(3, folderTreeLess),
	}
}

func folderTreeLess(a, b *FolderTreeNode) bool {
	return a.folder.Name < b.folder.Name
}

func preProcessFolders(folders []Folder) {
	// sort folders in place
	// primary key OrgId, secondary key paths
	slices.SortFunc(folders, func(a, b Folder) int {
		return strings.Compare(a.Paths, b.Paths)
	})
	slices.SortStableFunc(folders, func(a, b Folder) int {
		if string(a.OrgId.Bytes()) < string(b.OrgId.Bytes()) {
			return -1
		} else if string(a.OrgId.Bytes()) > string(b.OrgId.Bytes()) {
			return 1
		}
		return 0
	})
	return
}

func buildOrgs(folders []Folder) []Org {
	var orgs []Org
	preProcessFolders(folders)

	var orgMu sync.Mutex
	var orgWg sync.WaitGroup

	hi := 0
	for hi < len(folders) {
		lo := hi
		orgId := folders[lo].OrgId
		hi = sort.Search(len(folders), func(i int) bool {
			return folders[i].OrgId != orgId
		})

		orgWg.Add(1)
		go func(lo, hi int) {
			defer orgWg.Done()
			org := buildOrg(folders[lo:hi])

			orgMu.Lock()
			orgs = append(orgs, org)
			orgMu.Unlock()
		}(lo, hi)
	}

	orgWg.Wait()
	return orgs
}

func buildOrg(folders []Folder) Org {
	org := NewOrg(folders[0].OrgId)

	for _, folder := range folders {
		org.insertFolder(NewFolderTreeNode(&folder))
	}

	return org
}

func lookupTreeNode(folders *btree.BTreeG[*FolderTreeNode], target string) (*FolderTreeNode, bool) {
	return folders.Get(&FolderTreeNode{
		folder: &Folder{Name: target},
	})
}

func (org *Org) insertFolder(node *FolderTreeNode) error {
	parts := strings.Split(node.folder.Paths, ".")
	if len(parts) == 0 {
		return errors.New("Cannot insert folder with empty path")
	}

	curr, found := lookupTreeNode(org.folders, parts[0])
	if found == false {
		if len(parts) == 1 {
			org.folders.ReplaceOrInsert(node)
		}
		return nil
	}
	parts = append(parts[1:])

	for idx, part := range parts {
		next, found := lookupTreeNode(curr.children, part)

		if found == false {
			if idx != len(parts)-1 {
				// missing folders on path
				return errors.New("Could not insert, missing folders on path")
			}

			// insert folder to tree
			curr.children.ReplaceOrInsert(node)
			return nil
		}

		curr = next
	}

	// folder already exists at this location
	return errors.New("Could not insert, folder already exists at this location")
}

func (f *driver) getOrg(orgID uuid.UUID) (*Org, error) {
	for _, org := range f.orgs {
		if org.orgId == orgID {
			return &org, nil
		}
	}
	return nil, fmt.Errorf("No Org found with orgID %s", orgID.String())
}

func (f *driver) nameToOrgFolder(name string) (*Org, *FolderTreeNode, error) {
	for _, org := range f.orgs {
		// ignore err as folder may be in a different Org
		folder, _ := org.GetNamedFolder(name)
		if folder != nil {
			return &org, folder, nil
		}
	}
	return nil, nil, errors.New("Could not locate folder")
}
