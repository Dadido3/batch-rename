// Copyright (c) 2022 David Vogel
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
