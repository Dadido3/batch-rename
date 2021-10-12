# batch-rename

A simple command line tool to rename a bunch of files.

When run from inside some directory, this program will:

1. Create a temporary file with a list of all files contained in the given directory and its sub directories
2. Opens the temporary file with the default editor (For files with the file extension `.batch-rename`)
3. Reads the file back, and renames/moves the files

## Usage

There are multiple ways to use this tool:

- Download the compiled binary, then move and execute it into the directory you want to rename files in.
- Download the compiled binary, and move it to any directory that is in your `path` environment variable. Afterwards you can run `batch-rename` from inside any directory.
- Use `go install github.com/Dadido3/batch-rename`. Afterwards you can run `batch-rename` from inside any directory.
