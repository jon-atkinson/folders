package folder

import (
	"context"
	"errors"
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

type driver struct {
	orgs *btree.BTreeG[*Org]
}

type Org struct {
	orgId   uuid.UUID
	folders *btree.BTreeG[*FolderTreeNode]
}

type FolderTreeNode struct {
	folder   *Folder
	children *btree.BTreeG[*FolderTreeNode]
}

func NewDriver(folders []Folder) IDriver {
	return &driver{
		orgs: buildOrgs(folders),
	}
}

func NewOrg(orgID uuid.UUID) *Org {
	return &Org{
		orgId:   orgID,
		folders: btree.NewG(3, folderTreeLess),
	}
}

func NewFolderTreeNode(folder *Folder) *FolderTreeNode {
	return &FolderTreeNode{
		folder:   folder,
		children: btree.NewG(3, folderTreeLess),
	}
}

func orgTreeLess(a, b *Org) bool {
	return a.orgId.String() < b.orgId.String()
}

func folderTreeLess(a, b *FolderTreeNode) bool {
	return a.folder.Name < b.folder.Name
}

// sort folders in place
// primary key OrgId, secondary key paths
func preProcessFolders(folders []Folder) {
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
}

// Builds all Orgs, Org construction is managed by one goroutine / Org
// All goroutines return before returning Orgs btree
func buildOrgs(folders []Folder) *btree.BTreeG[*Org] {
	preProcessFolders(folders)
	var orgs *btree.BTreeG[*Org] = btree.NewG(3, orgTreeLess)

	orgChan := make(chan *Org)
	var orgWg sync.WaitGroup
	var mu sync.Mutex

	hi := 0
	for hi < len(folders) {
		orgId := folders[hi].OrgId
		lo := hi

		hi = sort.Search(len(folders), func(i int) bool {
			return folders[i].OrgId != orgId
		})

		orgWg.Add(1)
		go func(folderSlice []Folder) {
			defer orgWg.Done()
			orgChan <- buildOrg(folderSlice)
		}(folders[lo:hi])
	}

	go func() {
		orgWg.Wait()
		close(orgChan)
	}()

	for org := range orgChan {
		mu.Lock()
		orgs.ReplaceOrInsert(org)
		mu.Unlock()
	}

	return orgs
}

func buildOrg(folders []Folder) *Org {
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

// inserts folder into correct position in org
// assumes org is correct
// navigates org based on node.folder.Paths
func (org *Org) insertFolder(node *FolderTreeNode) error {
	parts := strings.Split(node.folder.Paths, ".")
	if len(parts) == 0 {
		return errors.New("Cannot insert folder with empty path")
	}

	curr, found := lookupTreeNode(org.folders, parts[0])

	if !found {
		if len(parts) == 1 {
			org.folders.ReplaceOrInsert(node)
		}
		return nil
	}
	parts = parts[1:]

	for idx, part := range parts {
		next, found := lookupTreeNode(curr.children, part)

		if !found {
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
	node, found := f.orgs.Get(&Org{orgId: orgID})
	if !found {
		return nil, errors.New("Organization does not exist")
	}
	return node, nil
}

// searches for folder called name in all orgs concurrently
// returns pointer to the tree node folder is stored in and the owning Org
// errors on miss
func (f *driver) nameToOrgFolder(name string) (*Org, *FolderTreeNode, error) {
	resultChan := make(chan struct {
		org    *Org
		folder *FolderTreeNode
	}, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(f.orgs.Len())

	f.orgs.Ascend(func(testOrg *Org) bool {

		go func(routineOrg *Org) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
				folder, err := routineOrg.GetNamedFolder(name)
				if folder != nil && err == nil {
					select {
					case resultChan <- struct {
						org    *Org
						folder *FolderTreeNode
					}{routineOrg, folder}:
						cancel()
					default:
						return
					}
				}
			}
		}(testOrg)
		return true
	})

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	result, ok := <-resultChan
	if ok {
		return result.org, result.folder, nil
	}
	return nil, nil, errors.New("Folder does not exist")
}
