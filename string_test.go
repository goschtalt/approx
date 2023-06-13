// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: BSD-3-Clause

package approx_test

import (
	"testing"
	"time"

	"github.com/goschtalt/approx"
)

func TestString(t *testing.T) {
	tests := []struct {
		alt    string
		d      time.Duration
		format string
		expect string
	}{
		{
			d:      approx.Year,
			expect: "1y0w0d0h0m0s",
			format: "ywd",
		}, {
			d:      52*approx.Week + approx.Day,
			expect: "1y0w0d0h0m0s",
			format: "ywd",
			alt:    "weeks+days",
		}, {
			d:      4 * approx.Year,
			expect: "4y0w0d0h0m0s",
			format: "ywd",
		}, {
			d:      4*approx.Year + 3*approx.Week + 2*approx.Day + 1*time.Hour,
			expect: "4y3w2d1h0m0s",
			format: "ywd",
		}, {
			d:      3 * approx.Week,
			expect: "3w0d0h0m0s",
			format: "ywd",
		}, {
			d:      3*approx.Week + time.Second,
			expect: "3w0d0h0m1s",
			format: "ywd",
		}, {
			d:      -1 * (3*approx.Week + time.Second),
			expect: "-3w0d0h0m1s",
			format: "ywd",
		}, {
			d:      60 * time.Minute,
			expect: "1h0m0s",
			format: "ywd",
		},
	}

	for _, tc := range tests {
		t.Run(tc.expect+tc.alt, func(t *testing.T) {
			got := approx.String(tc.d, tc.format)
			if tc.expect != got {
				t.Logf("expected: %q got: %q\n", tc.expect, got)
				t.Fail()
			}
		})
	}
}
