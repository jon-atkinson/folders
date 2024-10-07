package main

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jon-atkinson/sc-takehome-2024-25/folder"
)

const (
	FirstOrgID  = "c1556e17-b7c0-45a3-a6ae-9546248fb17a"
	SecondOrgID = "38b9879b-f73b-4b0e-b9d9-4fc4c23643a7"
)

func main() {
	firstOrgId := uuid.FromStringOrNil(FirstOrgID)

	f := folder.NewDriver([]folder.Folder{
		{"alpha", firstOrgId, "alpha"},
		{"a", firstOrgId, "a"},
		{"b", firstOrgId, "b"},
		{"c", firstOrgId, "c"},
		{"d", firstOrgId, "d"},
		{"e", firstOrgId, "e"},
		{"f", firstOrgId, "f"},
		{"g", firstOrgId, "g"},
		{"h", firstOrgId, "h"},
		{"i", firstOrgId, "i"},
		{"j", firstOrgId, "j"},
		{"k", firstOrgId, "k"},
		{"l", firstOrgId, "l"},
		{"m", firstOrgId, "m"},
		{"n", firstOrgId, "n"},
		{"o", firstOrgId, "o"},
		{"p", firstOrgId, "p"},
		{"q", firstOrgId, "q"},
		{"r", firstOrgId, "r"},
		{"s", firstOrgId, "s"},
		{"t", firstOrgId, "t"},
		{"bravo", firstOrgId, "alpha.bravo"},
		{"charlie", firstOrgId, "alpha.charlie"},
		{"delta", firstOrgId, "alpha.bravo.delta"},
		{"echo", firstOrgId, "alpha.bravo.echo"},
		{"foxtrot", firstOrgId, "alpha.charlie.foxtrot"},
		{"gamma", firstOrgId, "alpha.charlie.gamma"},
	})

	res, err := f.GetAllFolders()
	if err != nil {
		fmt.Println(err)
	} else {
		for _, ff := range res {
			folder.PrettyPrint(ff)
		}
	}
}
