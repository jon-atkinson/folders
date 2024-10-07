package main

import (
	"github.com/georgechieng-sc/interns-2022/folder"
)

func main() {
	res := folder.GenerateData()

	folder.PrettyPrint(res)

	folder.WriteSampleData(res)
}
