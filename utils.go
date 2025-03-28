package main

import (
	"os"
)

func createFolder(folderName string) error {
	err := os.MkdirAll(folderName, 0755)
	if err != nil {
		return err
	}
	return nil
}
