// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: BSD-3-Clause

package approx_test

import (
	"fmt"

	"github.com/gerifield/approx"
)

func Example_test() {

	d, err := approx.ParseDuration("1w4d")
	if err != nil {
		panic(err)
	}

	fmt.Println(approx.String(d))
	fmt.Println(d)

	// Output:
	// 11d0h0m0s
	// 264h0m0s
}
