package folder

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/gofrs/uuid"
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
	// I'm not sure how to determine which org to select (name, dst) from here
}

type FolderTreeNode struct {
	folder   *Folder
	children []*FolderTreeNode
}

type Org struct {
	orgId   uuid.UUID
	folders []*FolderTreeNode
}

type driver struct {
	orgs []Org
}

func NewDriver(folders []Folder) IDriver {
	return &driver{
		orgs: buildOrgs(folders),
	}
}

func preProcessFolders(folders []Folder) {
	// sort folders in place
	// primary key orgID, secondary key Path
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

	hi := 0
	for hi < len(folders) {
		lo := hi
		orgId := folders[lo].OrgId
		for hi+1 < len(folders) && folders[hi+1].OrgId == orgId {
			hi++
		}
		// fmt.Printf("%d\t%s\n%d\t%s\n", lo, folders[lo].OrgId, hi, folders[hi].OrgId)
		// fmt.Println(folders)
		hi++
		orgs = append(orgs, buildOrg(folders[lo:hi]))
	}
	return orgs
}

func buildOrg(folders []Folder) Org {
	var org Org = Org{
		orgId:   folders[0].OrgId,
		folders: []*FolderTreeNode{},
	}

	for _, folder := range folders {
		org.insertFolder(&folder)
	}

	return org
}

// assumes sorted folders, makes log(len(folders)) calls to the lambda
// returns nil on miss
func lookupTreeNode(folders []*FolderTreeNode, target string) *FolderTreeNode {
	idx := sort.Search(len(folders), func(i int) bool {
		return folders[i].folder.Name >= target
	})
	if idx == len(folders) {
		return nil
	}
	return folders[idx]
}

// only for calling from buildOrgs as assumes inputs called in sorted order
func (org *Org) insertFolder(folder *Folder) {
	parts := strings.Split(folder.Paths, ".")
	if len(parts) == 0 {
		return
	}

	curr := lookupTreeNode(org.folders, parts[0])
	if curr == nil {
		if len(parts) == 1 {
			org.folders = append(org.folders, &FolderTreeNode{
				folder,
				[]*FolderTreeNode{},
			})
		}
		return
	}
	parts = append(parts[1:])

	for idx, part := range parts {
		next := lookupTreeNode(curr.children, part)

		if next == nil {
			if idx != len(parts)-1 {
				// missing folders on path
				return
			}

			// insert folder to tree
			curr.children = append(curr.children, &FolderTreeNode{
				folder:   folder,
				children: []*FolderTreeNode{},
			})
			return
		}

		curr = next
	}

	// folder already exists at this location
	return
}

func (f *driver) getOrg(orgID uuid.UUID) (*Org, error) {
	for _, org := range f.orgs {
		if org.orgId == orgID {
			return &org, nil
		}
	}
	return nil, fmt.Errorf("No Org found with orgID %s", orgID.String())
}
