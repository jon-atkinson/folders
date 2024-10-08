package folder_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/jon-atkinson/sc-takehome-2024-25/folder"
)

func testFolderResults(t *testing.T, got []folder.Folder, want []folder.Folder, gotErr error, expErr error) {
	errString := "<nil>"
	ttErrString := "<nil>"
	if gotErr != nil {
		errString = gotErr.Error()
	}
	if expErr != nil {
		ttErrString = expErr.Error()
	}
	if errString != ttErrString {
		t.Fatalf("GetFoldersByOrgID wanted=%s. got=%s\n",
			ttErrString, errString)
	}

	if len(want) != len(got) {
		t.Fatalf("GetAllChildFolders output does not contain %d Folders. got=%d\n",
			len(want), len(got))
	}

	slices.SortFunc(want, func(a, b folder.Folder) int {
		return strings.Compare(a.Paths, b.Paths)
	})
	slices.SortFunc(got, func(a, b folder.Folder) int {
		return strings.Compare(a.Paths, b.Paths)
	})

	if diff := deep.Equal(want, got); diff != nil {
		t.Fatalf("GetAllChildFolders output folders do not match expected:\n%s",
			diff[0])
	}
}
