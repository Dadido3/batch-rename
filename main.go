// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// createEntriesFile creates a new file with all the fileEntries and returns its file path.
func createEntriesFile(rootDir string, fileEntries []fileEntry) (filePath string, err error) {
	// Create temporary file with paths.
	file, err := os.CreateTemp(".", "*.batch-rename")
	if err != nil {
		log.Fatal(err)
	}

	decimalWidth := getDecimalWidth(uint(len(fileEntries)))

	for i, fileEntry := range fileEntries {
		fmt.Fprintf(file, "%0"+fmt.Sprint(decimalWidth)+"d\t%s\n", i+1, fileEntry.originalPath)
	}

	if err := file.Close(); err != nil {
		return "", err
	}

	return file.Name(), nil
}

// readEntriesFile reads the (edited) temporary file back and modifies the fileEntries with the new file paths.
func readEntriesFile(filePath string, fileEntries []fileEntry) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		substrings := strings.SplitN(scanner.Text(), "\t", 2)
		entry, err := strconv.ParseUint(substrings[0], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse element number: %w", err)
		}
		if entry <= 0 || entry > uint64(len(fileEntries)) {
			return fmt.Errorf("element number outside of valid range: Given %v, there are %v entries", entry, len(fileEntries))
		}
		fileEntries[entry-1].newPath = substrings[1]
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to scan file: %w", err)
	}

	return nil
}

// moveFiles renames the files according to the list of file entries.
func moveFiles(fileEntries []fileEntry) error {
	for _, fileEntry := range fileEntries {
		if fileEntry.originalPath != fileEntry.newPath {
			if err := os.Rename(fileEntry.originalPath, fileEntry.newPath); err != nil {
				return err
			}
			log.Printf("Renamed %q to %q", fileEntry.originalPath, fileEntry.newPath)
			fileEntry.originalPath = fileEntry.newPath
		}
	}

	return nil
}

func main() {
	log.Printf("Started batch-rename v%v", version)

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

	if err := openWithDefault(filePath); err != nil {
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
