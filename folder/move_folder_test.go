package folder_test

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/go-test/deep"
	"github.com/gofrs/uuid"
)

func Test_folder_MoveFolder(t *testing.T) {
	firstOrgId := uuid.FromStringOrNil(FirstOrgID)
	secondOrgId := uuid.FromStringOrNil(SecondOrgID)

	t.Parallel()
	tests := [...]struct {
		name    string
		target  string
		dst     string
		folders []folder.Folder
		want    []folder.Folder
		err     error
	}{
		{
			"move top-level to top-level",
			"alpha",
			"bravo",
			[]folder.Folder{
				{"bravo", firstOrgId, "bravo"},
				{"alpha", firstOrgId, "alpha"},
			},
			[]folder.Folder{
				{"alpha", firstOrgId, "bravo.alpha"},
				{"bravo", firstOrgId, "bravo"},
			},
			nil,
		},
		{
			"move non-top-level to top-level",
			"charlie",
			"bravo",
			[]folder.Folder{
				{"bravo", firstOrgId, "bravo"},
				{"alpha", firstOrgId, "alpha"},
				{"charlie", firstOrgId, "alpha.charlie"},
			},
			[]folder.Folder{
				{"bravo", firstOrgId, "bravo"},
				{"alpha", firstOrgId, "alpha"},
				{"charlie", firstOrgId, "bravo.charlie"},
			},
			nil,
		},
		{
			"move top-level to non-top-level",
			"bravo",
			"charlie",
			[]folder.Folder{
				{"bravo", firstOrgId, "bravo"},
				{"alpha", firstOrgId, "alpha"},
				{"charlie", firstOrgId, "alpha.charlie"},
			},
			[]folder.Folder{
				{"bravo", firstOrgId, "alpha.charlie.bravo"},
				{"alpha", firstOrgId, "alpha"},
				{"charlie", firstOrgId, "alpha.charlie"},
			},
			nil,
		},
		{
			"move non-top-level to non-top-level",
			"delta",
			"charlie",
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "bravo"},
				{"charlie", firstOrgId, "alpha.charlie"},
				{"delta", firstOrgId, "bravo.delta"},
			},
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "bravo"},
				{"charlie", firstOrgId, "alpha.charlie"},
				{"delta", firstOrgId, "alpha.charlie.delta"},
			},
			nil,
		},
		{
			"deeper general case",
			"india",
			"kilo",
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"charlie", firstOrgId, "alpha.bravo.charlie"},
				{"kilo", firstOrgId, "alpha.bravo.kilo"},
				{"delta", firstOrgId, "alpha.delta"},
				{"echo", firstOrgId, "alpha.delta.echo"},
				{"foxtrot", firstOrgId, "alpha.foxtrot"},
				{"hotel", firstOrgId, "alpha.foxtrot.hotel"},
				{"india", firstOrgId, "alpha.foxtrot.india"},
				{"juliet", firstOrgId, "alpha.foxtrot.india.juliet"},
			},
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"charlie", firstOrgId, "alpha.bravo.charlie"},
				{"kilo", firstOrgId, "alpha.bravo.kilo"},
				{"delta", firstOrgId, "alpha.delta"},
				{"echo", firstOrgId, "alpha.delta.echo"},
				{"foxtrot", firstOrgId, "alpha.foxtrot"},
				{"hotel", firstOrgId, "alpha.foxtrot.hotel"},
				{"india", firstOrgId, "alpha.bravo.kilo.india"},
				{"juliet", firstOrgId, "alpha.bravo.kilo.india.juliet"},
			},
			nil,
		},
		{
			"attempt move to own child",
			"alpha",
			"charlie",
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "bravo"},
				{"charlie", firstOrgId, "alpha.charlie"},
			},
			[]folder.Folder{},
			errors.New("Cannot move a folder to a child of itself"),
		},
		{
			"attempt move to self",
			"alpha",
			"alpha",
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
			},
			[]folder.Folder{},
			errors.New("Cannot move a folder to itself"),
		},
		{
			"attempt move to different organization",
			"alpha",
			"bravo",
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", secondOrgId, "bravo"},
			},
			[]folder.Folder{},
			errors.New("Cannot move a folder to a different organization"),
		},
		{
			"attempt move non-existant source folder",
			"invalid",
			"bravo",
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "bravo"},
			},
			[]folder.Folder{},
			errors.New("Source folder does not exist"),
		},
		{
			"attempt move to non-existant destination folder",
			"bravo",
			"invalid",
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "bravo"},
			},
			[]folder.Folder{},
			errors.New("Destination folder does not exist"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)
			got, err := f.MoveFolder(tt.target, tt.dst)

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
					ttErrString, errString)
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
