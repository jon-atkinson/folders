package folder_test

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/go-test/deep"
	"github.com/gofrs/uuid"
)

const (
	FirstOrgID  = "c1556e17-b7c0-45a3-a6ae-9546248fb17a"
	SecondOrgID = "38b9879b-f73b-4b0e-b9d9-4fc4c23643a7"
)

func Test_folder_GetFoldersByOrgID(t *testing.T) {
	firstOrgId := uuid.FromStringOrNil(FirstOrgID)
	secondOrgID := uuid.FromStringOrNil(SecondOrgID)

	t.Parallel()
	tests := [...]struct {
		name    string
		orgID   uuid.UUID
		folders []folder.Folder
		want    []folder.Folder
		err     error
	}{
		{
			"no folders",
			firstOrgId,
			[]folder.Folder{},
			[]folder.Folder{},
			fmt.Errorf("No Org found with orgID %s", firstOrgId.String()),
		},
		{
			"single folder",
			firstOrgId,
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
			},
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
			},
			nil,
		},
		{
			"child folders",
			firstOrgId,
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"charlie", firstOrgId, "alpha.charlie"},
			},
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"charlie", firstOrgId, "alpha.charlie"},
			},
			nil,
		},
		{
			"multiple orgs",
			firstOrgId,
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", secondOrgID, "bravo"},
			},
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
			},
			nil,
		},
		{
			"multiple orgs, same folder name",
			firstOrgId,
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"alpha", secondOrgID, "alpha"},
			},
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
			},
			nil,
		},
		{
			"deeper folder trees",
			firstOrgId,
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"charlie", firstOrgId, "alpha.bravo.charlie"},
				{"delta", firstOrgId, "alpha.delta"},
				{"echo", firstOrgId, "alpha.delta.echo"},
				{"foxtrot", firstOrgId, "alpha.foxtrot"},
				{"golf", firstOrgId, "alpha.golf"},
				{"hotel", firstOrgId, "alpha.foxtrot.hotel"},
				{"india", firstOrgId, "alpha.foxtrot.india"},
				{"juliet", firstOrgId, "alpha.foxtrot.india.juliet"},
				{"kilo", firstOrgId, "alpha.foxtrot.india.juliet.kilo"},
			},
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"charlie", firstOrgId, "alpha.bravo.charlie"},
				{"delta", firstOrgId, "alpha.delta"},
				{"echo", firstOrgId, "alpha.delta.echo"},
				{"foxtrot", firstOrgId, "alpha.foxtrot"},
				{"golf", firstOrgId, "alpha.golf"},
				{"hotel", firstOrgId, "alpha.foxtrot.hotel"},
				{"india", firstOrgId, "alpha.foxtrot.india"},
				{"juliet", firstOrgId, "alpha.foxtrot.india.juliet"},
				{"kilo", firstOrgId, "alpha.foxtrot.india.juliet.kilo"},
			},
			nil,
		},
		{
			"unsorted input folders",
			firstOrgId,
			[]folder.Folder{
				{"foxtrot", firstOrgId, "alpha.foxtrot"},
				{"alpha", firstOrgId, "alpha"},
				{"charlie", firstOrgId, "alpha.bravo.charlie"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"kilo", firstOrgId, "alpha.foxtrot.india.juliet.kilo"},
				{"delta", firstOrgId, "alpha.delta"},
				{"golf", firstOrgId, "alpha.golf"},
				{"india", firstOrgId, "alpha.foxtrot.india"},
				{"hotel", firstOrgId, "alpha.foxtrot.hotel"},
				{"echo", firstOrgId, "alpha.delta.echo"},
				{"juliet", firstOrgId, "alpha.foxtrot.india.juliet"},
				{"kilo", firstOrgId, "alpha.bravo.kilo"},
			},
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"charlie", firstOrgId, "alpha.bravo.charlie"},
				{"delta", firstOrgId, "alpha.delta"},
				{"echo", firstOrgId, "alpha.delta.echo"},
				{"kilo", firstOrgId, "alpha.bravo.kilo"},
				{"foxtrot", firstOrgId, "alpha.foxtrot"},
				{"golf", firstOrgId, "alpha.golf"},
				{"hotel", firstOrgId, "alpha.foxtrot.hotel"},
				{"india", firstOrgId, "alpha.foxtrot.india"},
				{"juliet", firstOrgId, "alpha.foxtrot.india.juliet"},
				{"kilo", firstOrgId, "alpha.foxtrot.india.juliet.kilo"},
			},
			nil,
		},
		{
			"incorrect OrgID",
			secondOrgID,
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
			},
			[]folder.Folder{},
			fmt.Errorf("No Org found with orgID %s", secondOrgID.String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)
			got, err := f.GetFoldersByOrgID(tt.orgID)

			errString := "<nil>"
			ttErrString := "<nil>"
			if err != nil {
				errString = err.Error()
			}
			if tt.err != nil {
				ttErrString = tt.err.Error()
			}
			if errString != ttErrString {
				t.Fatalf("GetFoldersByOrgID wanted=%s. got=%s\n",
					errString, ttErrString)
			}

			if len(tt.want) != len(got) {
				t.Fatalf("GetFoldersByOrgID output does not contain %d Folders. got=%d\n",
					len(tt.want), len(got))
			}

			slices.SortFunc(tt.want, func(a, b folder.Folder) int {
				return strings.Compare(a.Paths, b.Paths)
			})
			slices.SortFunc(got, func(a, b folder.Folder) int {
				return strings.Compare(a.Paths, b.Paths)
			})

			if diff := deep.Equal(tt.want, got); diff != nil {
				t.Fatalf("GetFoldersByOrgID output folders do not match expected:\n%s",
					diff)
			}
		})
	}
}

func Test_folder_GetAllChildFolders(t *testing.T) {
	firstOrgId := uuid.FromStringOrNil(FirstOrgID)
	secondOrgID := uuid.FromStringOrNil(SecondOrgID)

	t.Parallel()
	tests := [...]struct {
		name         string
		orgID        uuid.UUID
		targetFolder string
		folders      []folder.Folder
		want         []folder.Folder
		err          error
	}{
		{
			"no folders",
			firstOrgId,
			"alpha",
			[]folder.Folder{},
			[]folder.Folder{},
			fmt.Errorf("No Org found with orgID %s", firstOrgId.String()),
		},
		{
			"single folder",
			firstOrgId,
			"alpha",
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
			},
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
			},
			nil,
		},
		{
			"child folders",
			firstOrgId,
			"alpha",
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"charlie", firstOrgId, "alpha.charlie"},
			},
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"charlie", firstOrgId, "alpha.charlie"},
			},
			nil,
		},
		{
			"child folders exclude some",
			firstOrgId,
			"alpha",
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"charlie", firstOrgId, "charlie"},
			},
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "alpha.bravo"},
			},
			nil,
		},
		{
			"target not root",
			firstOrgId,
			"bravo",
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"charlie", firstOrgId, "alpha.charlie"},
			},
			[]folder.Folder{
				{"bravo", firstOrgId, "alpha.bravo"},
			},
			nil,
		},
		{
			"multiple orgs",
			firstOrgId,
			"alpha",
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", secondOrgID, "bravo"},
			},
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
			},
			nil,
		},
		{
			"multiple orgs, same folder name",
			firstOrgId,
			"alpha",
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"alpha", secondOrgID, "alpha"},
			},
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
			},
			nil,
		},
		{
			"deeper folder trees",
			firstOrgId,
			"india",
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"charlie", firstOrgId, "alpha.bravo.charlie"},
				{"delta", firstOrgId, "alpha.delta"},
				{"echo", firstOrgId, "alpha.delta.echo"},
				{"foxtrot", firstOrgId, "alpha.foxtrot"},
				{"golf", firstOrgId, "alpha.golf"},
				{"hotel", firstOrgId, "alpha.foxtrot.hotel"},
				{"india", firstOrgId, "alpha.foxtrot.india"},
				{"juliet", firstOrgId, "alpha.foxtrot.india.juliet"},
				{"kilo", firstOrgId, "alpha.foxtrot.india.juliet.kilo"},
			},
			[]folder.Folder{
				{"india", firstOrgId, "alpha.foxtrot.india"},
				{"juliet", firstOrgId, "alpha.foxtrot.india.juliet"},
				{"kilo", firstOrgId, "alpha.foxtrot.india.juliet.kilo"},
			},
			nil,
		},
		{
			"unsorted input folders",
			firstOrgId,
			"india",
			[]folder.Folder{
				{"foxtrot", firstOrgId, "alpha.foxtrot"},
				{"alpha", firstOrgId, "alpha"},
				{"charlie", firstOrgId, "alpha.bravo.charlie"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"kilo", firstOrgId, "alpha.foxtrot.india.juliet.kilo"},
				{"delta", firstOrgId, "alpha.delta"},
				{"golf", firstOrgId, "alpha.golf"},
				{"india", firstOrgId, "alpha.foxtrot.india"},
				{"hotel", firstOrgId, "alpha.foxtrot.hotel"},
				{"echo", firstOrgId, "alpha.delta.echo"},
				{"juliet", firstOrgId, "alpha.foxtrot.india.juliet"},
				{"kilo", firstOrgId, "alpha.bravo.kilo"},
			},
			[]folder.Folder{
				{"india", firstOrgId, "alpha.foxtrot.india"},
				{"juliet", firstOrgId, "alpha.foxtrot.india.juliet"},
				{"kilo", firstOrgId, "alpha.foxtrot.india.juliet.kilo"},
			},
			nil,
		},
		{
			"incorrect OrgID",
			secondOrgID,
			"alpha",
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
			},
			[]folder.Folder{},
			fmt.Errorf("No Org found with orgID %s", secondOrgID.String()),
		},
		{
			"folder does not exist",
			firstOrgId,
			"delta",
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"charlie", firstOrgId, "alpha.bravo.charlie"},
				{"sierra", firstOrgId, "alpha.sierra"},
				{"kilo", firstOrgId, "alpha.bravo.charlie.kilo"},
				{"uniform", firstOrgId, "alpha.sierra.uniform"},
				{"zulu", firstOrgId, "alpha.sierra.zulu"},
				{"mike", firstOrgId, "alpha.sierra.mike"},
			},
			[]folder.Folder{},
			errors.New("Folder does not exist"),
		},
		{
			"folder belongs to a different organization",
			firstOrgId,
			"delta",
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"charlie", firstOrgId, "alpha.bravo.charlie"},
				{"sierra", firstOrgId, "alpha.sierra"},
				{"kilo", firstOrgId, "alpha.bravo.charlie.kilo"},
				{"delta", secondOrgID, "delta"},
			},
			[]folder.Folder{},
			errors.New("Folder does not exist in the specified organization"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)
			got, err := f.GetAllChildFolders(tt.orgID, tt.targetFolder)

			if err != tt.err && err.Error() != tt.err.Error() {
				t.Fatalf("GetAllChildFolders wanted=%s. got=%s\n", err.Error(),
					tt.err.Error())
			}

			if len(tt.want) != len(got) {
				t.Fatalf("GetAllChildFolders output does not contain %d Folders. got=%d\n",
					len(tt.want), len(got))
			}

			slices.SortFunc(tt.want, func(a, b folder.Folder) int {
				return strings.Compare(a.Paths, b.Paths)
			})
			slices.SortFunc(got, func(a, b folder.Folder) int {
				return strings.Compare(a.Paths, b.Paths)
			})

			if diff := deep.Equal(tt.want, got); diff != nil {
				t.Fatalf("GetAllChildFolders output folders do not match expected:\n%s",
					diff)
			}
		})
	}
}
