package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
)

const (
	FirstOrgID  = "c1556e17-b7c0-45a3-a6ae-9546248fb17a"
	SecondOrgID = "38b9879b-f73b-4b0e-b9d9-4fc4c23643a7"
)

func main() {
	fmt.Println("Starting Virtual File System REPL...")
	fmt.Println("Available commands:")
	fmt.Println("  - list: List all folders")
	fmt.Println("  - get <orgID>: Get folders by organization ID")
	fmt.Println("  - children <name>: Get children by name")
	fmt.Println("  - exit: Exit the REPL")
	fmt.Println()

	res := folder.GetAllFolders()
	folderDriver := folder.NewDriver(res)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		tokens := strings.Fields(line)

		if len(tokens) == 0 {
			continue
		}

		command := tokens[0]
		switch command {
		case "list":
			fmt.Println("Listing all folders:")
			folder.PrettyPrint(res)

		case "get":
			if len(tokens) < 2 {
				fmt.Println("Error: Missing orgID argument. Usage: get <orgID>")
				continue
			}
			orgIDStr := tokens[1]
			orgID := uuid.FromStringOrNil(orgIDStr)
			orgFolders := folderDriver.GetFoldersByOrgID(orgID)

			if len(orgFolders) == 0 {
				fmt.Printf("No folders found for orgID: %s\n", orgID)
			} else {
				fmt.Printf("Folders for orgID: %s\n", orgID)
				folder.PrettyPrint(orgFolders)
			}

		case "children":
			if len(tokens) < 3 {
				fmt.Println("Error: Missing argument. Usage: children <orgID,name>")
				continue
			}
			orgID := uuid.FromStringOrNil(tokens[1])
			nameStr := tokens[2]
			childFolders := folderDriver.GetAllChildFolders(orgID, nameStr)
			if len(childFolders) == 0 {
				fmt.Printf("No folders found for <orgID,name>: %s,%s\n", orgID, nameStr)
			} else {
				fmt.Printf("Folders for <orgID,name>: %s,%s\n", orgID, nameStr)
				folder.PrettyPrint(childFolders)
			}

		case "move":
			if len(tokens) < 3 {
				fmt.Println("Error: Missing argument. Usage: move <src,dst>")
				continue
			}
			src := tokens[1]
			dst := tokens[2]
			resultFolders, err := folderDriver.MoveFolder(src, dst)
			if err != nil {
				fmt.Printf("Error encountered. %s\n", err.Error())
			} else {
				fmt.Printf("Folders for <src,dst>: %s,%s\n", src, dst)
				folder.PrettyPrint(resultFolders)
			}

		case "q":
			fmt.Println("Exiting...")
			return
		case "quit":
			fmt.Println("Exiting...")
			return
		case "exit":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Printf("Unknown command: %s\n", command)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}
