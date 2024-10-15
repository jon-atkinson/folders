package folder_test

import (
	"errors"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/georgechieng-sc/interns-2022/folder"
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
		{
			"deeper trees, more organizations",
			"echo",
			"sierra",
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"charlie", firstOrgId, "alpha.charlie"},
				{"delta", firstOrgId, "alpha.bravo.delta"},
				{"echo", firstOrgId, "alpha.bravo.echo"},
				{"foxtrot", firstOrgId, "alpha.charlie.foxtrot"},
				{"golf", firstOrgId, "alpha.charlie.golf"},
				{"hotel", firstOrgId, "alpha.bravo.delta.hotel"},
				{"india", firstOrgId, "alpha.bravo.delta.india"},
				{"juliet", firstOrgId, "alpha.bravo.delta.juliet"},
				{"kilo", firstOrgId, "alpha.bravo.echo.kilo"},
				{"lima", firstOrgId, "alpha.bravo.echo.lima"},
				{"mike", firstOrgId, "alpha.bravo.echo.mike"},
				{"november", firstOrgId, "alpha.charlie.foxtrot.november"},
				{"oscar", firstOrgId, "alpha.charlie.foxtrot.oscar"},
				{"papa", firstOrgId, "alpha.charlie.foxtrot.papa"},
				{"quebec", firstOrgId, "alpha.charlie.golf.quebec"},
				{"romeo", firstOrgId, "alpha.charlie.golf.romeo"},
				{"sierra", firstOrgId, "alpha.charlie.golf.sierra"},
				{"tango", firstOrgId, "alpha.bravo.delta.hotel.tango"},
				{"uniform", firstOrgId, "alpha.bravo.delta.hotel.uniform"},
				{"victor", firstOrgId, "alpha.bravo.delta.india.victor"},
				{"whiskey", firstOrgId, "alpha.bravo.delta.india.whiskey"},
				{"x-ray", firstOrgId, "alpha.bravo.delta.juliet.x-ray"},
				{"yankee", firstOrgId, "alpha.bravo.delta.juliet.yankee"},
				{"zulu", firstOrgId, "alpha.bravo.echo.kilo.zulu"},
				{"alpha-2", secondOrgId, "alpha-2"},
				{"bravo-2", secondOrgId, "alpha-2.bravo-2"},
				{"charlie-2", secondOrgId, "alpha-2.charlie-2"},
				{"delta-2", secondOrgId, "alpha-2.bravo-2.delta-2"},
				{"echo-2", secondOrgId, "alpha-2.bravo-2.echo-2"},
				{"foxtrot-2", secondOrgId, "alpha-2.charlie-2.foxtrot-2"},
				{"golf-2", secondOrgId, "alpha-2.charlie-2.golf-2"},
				{"hotel-2", secondOrgId, "alpha-2.bravo-2.delta-2.hotel-2"},
				{"india-2", secondOrgId, "alpha-2.bravo-2.delta-2.india-2"},
				{"juliet-2", secondOrgId, "alpha-2.bravo-2.delta-2.juliet-2"},
				{"kilo-2", secondOrgId, "alpha-2.bravo-2.echo-2.kilo-2"},
				{"lima-2", secondOrgId, "alpha-2.bravo-2.echo-2.lima-2"},
				{"mike-2", secondOrgId, "alpha-2.bravo-2.echo-2.mike-2"},
				{"november-2", secondOrgId, "alpha-2.charlie-2.foxtrot-2.november-2"},
				{"oscar-2", secondOrgId, "alpha-2.charlie-2.foxtrot-2.oscar-2"},
				{"papa-2", secondOrgId, "alpha-2.charlie-2.foxtrot-2.papa-2"},
				{"quebec-2", secondOrgId, "alpha-2.charlie-2.golf-2.quebec-2"},
				{"romeo-2", secondOrgId, "alpha-2.charlie-2.golf-2.romeo-2"},
				{"sierra-2", secondOrgId, "alpha-2.charlie-2.golf-2.sierra-2"},
				{"tango-2", secondOrgId, "alpha-2.bravo-2.delta-2.hotel-2.tango-2"},
				{"uniform-2", secondOrgId, "alpha-2.bravo-2.delta-2.hotel-2.uniform-2"},
				{"victor-2", secondOrgId, "alpha-2.bravo-2.delta-2.india-2.victor-2"},
				{"whiskey-2", secondOrgId, "alpha-2.bravo-2.delta-2.india-2.whiskey-2"},
				{"x-ray-2", secondOrgId, "alpha-2.bravo-2.delta-2.juliet-2.x-ray-2"},
				{"yankee-2", secondOrgId, "alpha-2.bravo-2.delta-2.juliet-2.yankee-2"},
				{"zulu-2", secondOrgId, "alpha-2.bravo-2.echo-2.kilo-2.zulu-2"},
			},
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"charlie", firstOrgId, "alpha.charlie"},
				{"delta", firstOrgId, "alpha.bravo.delta"},
				{"echo", firstOrgId, "alpha.charlie.golf.sierra.echo"},
				{"foxtrot", firstOrgId, "alpha.charlie.foxtrot"},
				{"golf", firstOrgId, "alpha.charlie.golf"},
				{"hotel", firstOrgId, "alpha.bravo.delta.hotel"},
				{"india", firstOrgId, "alpha.bravo.delta.india"},
				{"juliet", firstOrgId, "alpha.bravo.delta.juliet"},
				{"kilo", firstOrgId, "alpha.charlie.golf.sierra.echo.kilo"},
				{"lima", firstOrgId, "alpha.charlie.golf.sierra.echo.lima"},
				{"mike", firstOrgId, "alpha.charlie.golf.sierra.echo.mike"},
				{"november", firstOrgId, "alpha.charlie.foxtrot.november"},
				{"oscar", firstOrgId, "alpha.charlie.foxtrot.oscar"},
				{"papa", firstOrgId, "alpha.charlie.foxtrot.papa"},
				{"quebec", firstOrgId, "alpha.charlie.golf.quebec"},
				{"romeo", firstOrgId, "alpha.charlie.golf.romeo"},
				{"sierra", firstOrgId, "alpha.charlie.golf.sierra"},
				{"tango", firstOrgId, "alpha.bravo.delta.hotel.tango"},
				{"uniform", firstOrgId, "alpha.bravo.delta.hotel.uniform"},
				{"victor", firstOrgId, "alpha.bravo.delta.india.victor"},
				{"whiskey", firstOrgId, "alpha.bravo.delta.india.whiskey"},
				{"x-ray", firstOrgId, "alpha.bravo.delta.juliet.x-ray"},
				{"yankee", firstOrgId, "alpha.bravo.delta.juliet.yankee"},
				{"zulu", firstOrgId, "alpha.charlie.golf.sierra.echo.kilo.zulu"},
				{"alpha-2", secondOrgId, "alpha-2"},
				{"bravo-2", secondOrgId, "alpha-2.bravo-2"},
				{"charlie-2", secondOrgId, "alpha-2.charlie-2"},
				{"delta-2", secondOrgId, "alpha-2.bravo-2.delta-2"},
				{"echo-2", secondOrgId, "alpha-2.bravo-2.echo-2"},
				{"foxtrot-2", secondOrgId, "alpha-2.charlie-2.foxtrot-2"},
				{"golf-2", secondOrgId, "alpha-2.charlie-2.golf-2"},
				{"hotel-2", secondOrgId, "alpha-2.bravo-2.delta-2.hotel-2"},
				{"india-2", secondOrgId, "alpha-2.bravo-2.delta-2.india-2"},
				{"juliet-2", secondOrgId, "alpha-2.bravo-2.delta-2.juliet-2"},
				{"kilo-2", secondOrgId, "alpha-2.bravo-2.echo-2.kilo-2"},
				{"lima-2", secondOrgId, "alpha-2.bravo-2.echo-2.lima-2"},
				{"mike-2", secondOrgId, "alpha-2.bravo-2.echo-2.mike-2"},
				{"november-2", secondOrgId, "alpha-2.charlie-2.foxtrot-2.november-2"},
				{"oscar-2", secondOrgId, "alpha-2.charlie-2.foxtrot-2.oscar-2"},
				{"papa-2", secondOrgId, "alpha-2.charlie-2.foxtrot-2.papa-2"},
				{"quebec-2", secondOrgId, "alpha-2.charlie-2.golf-2.quebec-2"},
				{"romeo-2", secondOrgId, "alpha-2.charlie-2.golf-2.romeo-2"},
				{"sierra-2", secondOrgId, "alpha-2.charlie-2.golf-2.sierra-2"},
				{"tango-2", secondOrgId, "alpha-2.bravo-2.delta-2.hotel-2.tango-2"},
				{"uniform-2", secondOrgId, "alpha-2.bravo-2.delta-2.hotel-2.uniform-2"},
				{"victor-2", secondOrgId, "alpha-2.bravo-2.delta-2.india-2.victor-2"},
				{"whiskey-2", secondOrgId, "alpha-2.bravo-2.delta-2.india-2.whiskey-2"},
				{"x-ray-2", secondOrgId, "alpha-2.bravo-2.delta-2.juliet-2.x-ray-2"},
				{"yankee-2", secondOrgId, "alpha-2.bravo-2.delta-2.juliet-2.yankee-2"},
				{"zulu-2", secondOrgId, "alpha-2.bravo-2.echo-2.kilo-2.zulu-2"},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)
			got, err := f.MoveFolder(tt.target, tt.dst)

			testFolderResults(t, got, tt.want)
			testFolderError(t, err, tt.err)
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
				{"golf", "foxtrot"},
				{"foxtrot", "echo"},
			},
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"charlie", firstOrgId, "alpha.charlie"},
				{"delta", firstOrgId, "alpha.bravo.delta"},
				{"echo", firstOrgId, "alpha.bravo.echo"},
				{"foxtrot", firstOrgId, "alpha.charlie.foxtrot"},
				{"golf", firstOrgId, "alpha.charlie.golf"},
			},
			[]folder.Folder{
				{"alpha", firstOrgId, "alpha"},
				{"bravo", firstOrgId, "alpha.bravo"},
				{"charlie", firstOrgId, "alpha.bravo.charlie"},
				{"delta", firstOrgId, "alpha.bravo.charlie.delta"},
				{"echo", firstOrgId, "alpha.bravo.charlie.delta.echo"},
				{"foxtrot", firstOrgId, "alpha.bravo.charlie.delta.echo.foxtrot"},
				{"golf", firstOrgId, "alpha.bravo.charlie.delta.echo.foxtrot.golf"},
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

			testFolderResults(t, got, tt.want)
			testFolderError(t, err, tt.err)
		})
	}
}
