// Copyright (c) 2021-2023 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"bufio"
	"flag"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/skratchdot/open-golang/open"
)

var flagNoNumbers = flag.Bool("no-numbers", false, "If set, batch-rename will not prepend numbers to every line. If you enable this option you have to make sure that you don't add or remove lines, as otherwise it will mess up your filenames!")
var flagFilterRegex = flag.String("regex", "", "Filters entries with the specified regular expression by their filepath relative to the working directory. Example: batch-rename --regex \"\\.png$|\\.jpg$\", which will only include png and jpg files.")

func main() {
	flag.Parse()

	log.Printf("Started batch-rename %v", Version)

	rootDir := "."
	numbering := !*flagNoNumbers

	var filterRegex *regexp.Regexp
	if *flagFilterRegex != "" {
		var err error
		if filterRegex, err = regexp.Compile(*flagFilterRegex); err != nil {
			log.Printf("Invalid regular expression: %v", err)
			return
		}
	}

	fileEntries := []fileEntry{}

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		switch d.IsDir() {
		case false:

			if filterRegex != nil && !filterRegex.Match([]byte(path)) {
				return nil
			}

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

	log.Printf("Found %d files", len(fileEntries))

	// Create temporary file.
	filePath, err := createEntriesFile(rootDir, fileEntries, numbering)
	if err != nil {
		log.Panicf("Failed to create file with filepaths: %v", err)
	}

	if err := open.Run(filePath); err != nil {
		log.Panicf("Failed to open file with filepaths: %v", err)
	}

	log.Printf("Opening %s in the default editor", filePath)
	log.Print("You can now edit the filepaths and press any key to continue...")
	reader.ReadString('\n')

	// Read temporary file back and rename files. Retry on fail.
	for {
		if err := readEntriesFile(filePath, fileEntries, numbering); err != nil {
			log.Printf("Failed to read filepaths: %v", err)
			log.Print("Please edit your filepaths, save them and press any key to retry...")
			reader.ReadString('\n')
			continue
		}

		if err := moveFiles(fileEntries); err != nil {
			log.Printf("Failed to move file: %v", err)
			log.Print("Please edit your filepaths, save them and press any key to retry...")
			reader.ReadString('\n')
			continue
		}

		break
	}

	log.Print("Deleting filepaths file")

	// Try to delete the temporary file. Retry on fail.
	for {
		if err := os.Remove(filePath); err != nil {
			log.Printf("Failed to remove file with filepaths: %v", err)
			log.Printf("Please close your editor and press any key to try again...")
			reader.ReadString('\n')
			continue
		}

		break
	}

}
