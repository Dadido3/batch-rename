// Copyright (c) 2021-2022 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"bufio"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/skratchdot/open-golang/open"
)

func main() {
	log.Printf("Started batch-rename %v", Version)

	rootDir := "."

	fileEntries := []fileEntry{}

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		switch d.IsDir() {
		case false:
			fileEntries = append(fileEntries, fileEntry{
				name:         d.Name(),
				originalPath: path,
				newPath:      path,
			})
		case true:
		}

		return nil
	})
	if err != nil {
		log.Panicf("Failed to read file tree: %v", err)
	}

	reader := bufio.NewReader(os.Stdin)

	if len(fileEntries) == 0 {
		log.Print("There are no files to be renamed.")
		reader.ReadString('\n')
		return
	}

	log.Printf("Got %d files to rename", len(fileEntries))

	// Create temporary file.
	filePath, err := createEntriesFile(rootDir, fileEntries)
	if err != nil {
		log.Panicf("Failed to create temporary text file: %v", err)
	}

	/*if err := openWithDefault(filePath); err != nil {
		log.Panicf("Failed to open temporary text file: %v", err)
	}*/

	if err := open.Run(filePath); err != nil {
		log.Panicf("Failed to open temporary text file: %v", err)
	}

	log.Print("Editor opened, edit file paths and press any key to continue...")
	reader.ReadString('\n')

	// Read temporary file back and rename files. Retry on fail.
	for {
		if err := readEntriesFile(filePath, fileEntries); err != nil {
			log.Printf("Failed to read temporary text file: %v", err)
			log.Print("Update the temporary text file and press any key to retry...")
			reader.ReadString('\n')
			continue
		}

		if err := moveFiles(fileEntries); err != nil {
			log.Printf("Failed to move file: %v", err)
			log.Print("Update the temporary file and press any key to retry...")
			reader.ReadString('\n')
			continue
		}

		break
	}

	log.Print("Deleting temporary file")

	// Try to delete the temporary file. Retry on fail.
	for {
		if err := os.Remove(filePath); err != nil {
			log.Printf("Failed to remove temporary text file: %v", err)
			log.Printf("Close the temporary text file and press any key to try again...")
			reader.ReadString('\n')
			continue
		}

		break
	}

}
