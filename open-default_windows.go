// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"os/exec"
)

func openWithDefault(filePath string) error {
	cmd := exec.Command("explorer", filePath)
	return cmd.Start()
}
