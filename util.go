// Copyright (c) 2021 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package main

// getDecimalWidth returns the width in digits that is needed to represent the given unsigned integer in the decimal format.
func getDecimalWidth(n uint) (digits uint) {
	for digits = 1; n >= 10; n /= 10 {
		digits++
	}

	return
}
