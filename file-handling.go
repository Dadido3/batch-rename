// Copyright (c) 2022-2023 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// createEntriesFile creates a new file with all the fileEntries and returns its file path.
func createEntriesFile(rootDir string, fileEntries []fileEntry, numbering bool) (filePath string, err error) {
	// Create temporary file with paths.
	file, err := os.CreateTemp(".", "*.batch-rename")
	if err != nil {
		log.Fatal(err)
	}

	if numbering {
		decimalWidth := getDecimalWidth(uint(len(fileEntries)))

		for i, fileEntry := range fileEntries {
			fmt.Fprintf(file, "%0"+fmt.Sprint(decimalWidth)+"d\t%s\n", i+1, fileEntry.originalPath)
		}
	} else {
		for _, fileEntry := range fileEntries {
			fmt.Fprintf(file, "%s\n", fileEntry.originalPath)
		}
	}

	if err := file.Close(); err != nil {
		return "", err
	}

	return file.Name(), nil
}

// readEntriesFile reads the (edited) temporary file back and modifies the fileEntries with the new file paths.
func readEntriesFile(filePath string, fileEntries []fileEntry, numbering bool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		if numbering {
			substrings := strings.SplitN(scanner.Text(), "\t", 2)
			entry, err := strconv.ParseUint(substrings[0], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse element number: %w", err)
			}
			if entry <= 0 || entry > uint64(len(fileEntries)) {
				return fmt.Errorf("element number outside of valid range: Got %d, but there are only %d entries", entry, len(fileEntries))
			}
			fileEntries[entry-1].newPath = substrings[1]
		} else {
			if i >= len(fileEntries) {
				return fmt.Errorf("there are more lines than files: Got %d, but expected %d", i+1, len(fileEntries))
			}
			fileEntries[i].newPath = scanner.Text()
			i++
		}
	}

	// In case we don't do numbering, we should make sure that the number of lines didn't change.
	if !numbering && i != len(fileEntries) {
		return fmt.Errorf("the number of lines don't match the number of files: Got %d, but expected %d", i, len(fileEntries))
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
			// This is important, because if there is an error later on, we know which files were already moved.
			// So we can query the user to correct the list, and just call moveFiles again, without files getting messed up.
			fileEntry.originalPath = fileEntry.newPath
		}
	}

	return nil
}
