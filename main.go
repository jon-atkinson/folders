package main

import (
	"github.com/jon-atkinson/sc-takehome-2024-25/folder"
)

func main() {
	res := folder.GenerateData()

	folder.PrettyPrint(res)

	folder.WriteSampleData(res)
}
