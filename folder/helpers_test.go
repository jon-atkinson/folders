package folder_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/go-test/deep"
)

func testFolderResults(t *testing.T, got []folder.Folder, want []folder.Folder) {
	if len(want) != len(got) {
		fmt.Println("got", got, "want", want)
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

// helper function for testing the errors returned by folder IDriver interface
// functions that return errors
func testFolderError(t *testing.T, gotErr error, expErr error) {
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
}
