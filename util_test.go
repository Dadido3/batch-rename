// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

import (
	"fmt"
	"testing"
)

func Test_getDecimalWidth(t *testing.T) {
	tests := []struct {
		n          uint
		wantDigits uint
	}{
		{0, 1},
		{1, 1},
		{9, 1},
		{10, 2},
		{11, 2},
		{19, 2},
		{20, 2},
		{21, 2},
		{99, 2},
		{100, 3},
		{101, 3},
		{999, 3},
		{1000, 4},
		{1001, 4},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.n), func(t *testing.T) {
			if gotDigits := getDecimalWidth(tt.n); gotDigits != tt.wantDigits {
				t.Errorf("getDecimalWidth() = %v, want %v", gotDigits, tt.wantDigits)
			}
		})
	}
}
