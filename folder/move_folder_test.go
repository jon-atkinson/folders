package folder_test

import (
	"errors"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jon-atkinson/sc-takehome-2024-25/folder"
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

			testFolderResults(t, got, tt.want, err, tt.err)
		})
	}
}

func Test_folder_MoveFolder_Complex(t *testing.T) {
	firstOrgId := uuid.FromStringOrNil(FirstOrgID)

	t.Parallel()
	tests := [...]struct {
		name  string
		moves []struct {
			target string
			dst    string
		}
		folders []folder.Folder
		want    []folder.Folder
		err     error
	}{
		{
			"two moves, top-level",
			[]struct {
				target string
				dst    string
			}{
				{"bravo", "alpha"},
				{"charlie", "bravo"},
			},
			[]folder.Folder{
				{"bravo", firstOrgId, "bravo"},
				{"alpha", firstOrgId, "alpha"},
				{"charlie", firstOrgId, "charlie"},
			},
			[]folder.Folder{
				{"bravo", firstOrgId, "alpha.bravo"},
				{"alpha", firstOrgId, "alpha"},
				{"charlie", firstOrgId, "alpha.bravo.charlie"},
			},
			nil,
		},
		{
			"moves are commutative",
			[]struct {
				target string
				dst    string
			}{
				{"charlie", "bravo"},
				{"bravo", "alpha"},
			},
			[]folder.Folder{
				{"bravo", firstOrgId, "bravo"},
				{"alpha", firstOrgId, "alpha"},
				{"charlie", firstOrgId, "charlie"},
			},
			[]folder.Folder{
				{"bravo", firstOrgId, "alpha.bravo"},
				{"alpha", firstOrgId, "alpha"},
				{"charlie", firstOrgId, "alpha.bravo.charlie"},
			},
			nil,
		},
		{
			"moving several folders at once",
			[]struct {
				target string
				dst    string
			}{
				{"charlie", "bravo"},
				{"bravo", "alpha"},
			},
			[]folder.Folder{
				{"bravo", firstOrgId, "bravo"},
				{"alpha", firstOrgId, "alpha"},
				{"charlie", firstOrgId, "charlie"},
				{"delta", firstOrgId, "charlie.delta"},
			},
			[]folder.Folder{
				{"bravo", firstOrgId, "alpha.bravo"},
				{"alpha", firstOrgId, "alpha"},
				{"charlie", firstOrgId, "alpha.bravo.charlie"},
				{"delta", firstOrgId, "alpha.bravo.charlie.delta"},
			},
			nil,
		},
		{
			"flatten btree",
			[]struct {
				target string
				dst    string
			}{
				{"charlie", "bravo"},
				{"charlie", "bravo"},
				{"echo", "delta"},
				{"delta", "charlie"},
				{"gamma", "foxtrot"},
				{"foxtrot", "echo"},
			},
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"charlie", firstOrgId, "alpha.charlie"},
				{"delta", firstOrgId, "alpha.bravo.delta"},
				{"echo", firstOrgId, "alpha.bravo.echo"},
				{"foxtrot", firstOrgId, "alpha.charlie.foxtrot"},
				{"gamma", firstOrgId, "alpha.charlie.gamma"},
			},
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"charlie", firstOrgId, "alpha.bravo.charlie"},
				{"delta", firstOrgId, "alpha.bravo.charlie.delta"},
				{"echo", firstOrgId, "alpha.bravo.charlie.delta.echo"},
				{"foxtrot", firstOrgId, "alpha.bravo.charlie.delta.echo.foxtrot"},
				{"gamma", firstOrgId, "alpha.bravo.charlie.delta.echo.foxtrot.gamma"},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)
			var got []folder.Folder
			var err error
			for _, move := range tt.moves {
				got, err = f.MoveFolder(move.target, move.dst)
			}

			testFolderResults(t, got, tt.want, err, tt.err)
		})
	}
}
