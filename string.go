// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: BSD-3-Clause

// Package approx provides a set of approximate duration units that extend the
// standard go package time.Duration.
package approx

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	Day   = time.Hour * 24
	Week  = Day * 7
	Month = Week * 4
	Year  = Day * 365
)

var (
	ErrInvalidDuration = errors.New("time: invalid duration")
)

// String converts a time.Duration into the larger approximate units for output.
// The default form is to prefix days, but weeks and years may be prepended, too.
// Passing a string with the options "d", "w", and "y" will enable each.  Order
// and repeats are ignored.
func String(d time.Duration, format ...string) string {
	var b strings.Builder

	// Ensure the buffer doesn't have to grow more than once.
	b.Grow(40)

	format = append(format, "d")
	withYears := strings.Contains(format[0], "y")
	withMonths := strings.Contains(format[0], "M")
	withWeeks := strings.Contains(format[0], "w")
	withDays := strings.Contains(format[0], "d")

	ns := int64(d)

	if ns < 0 {
		b.WriteString("-")
		ns = -1 * ns
	}

	setSmaller := false

	if withYears && ns >= int64(Year) {
		years := ns / int64(Year)
		ns -= years * int64(Year)
		b.WriteString(fmt.Sprintf("%dy", years))
		setSmaller = true
	}

	if withMonths && ns >= int64(Month) {
		months := ns / int64(Month)
		ns -= months * int64(Month)
		b.WriteString(fmt.Sprintf("%dM", months))
		setSmaller = true
	} else {
		if withMonths && setSmaller {
			b.WriteString("0M")
		}
	}
	if withWeeks && ns >= int64(Week) {
		weeks := ns / int64(Week)
		ns -= weeks * int64(Week)
		b.WriteString(fmt.Sprintf("%dw", weeks))
		setSmaller = true
	} else {
		if withWeeks && setSmaller {
			b.WriteString("0w")
		}
	}
	if withDays && ns >= int64(Day) {
		days := ns / int64(Day)
		ns -= days * int64(Day)
		b.WriteString(fmt.Sprintf("%dd", days))
		setSmaller = true
	} else {
		if withDays && setSmaller {
			b.WriteString("0d")
		}
	}
	if setSmaller {
		if ns < int64(time.Hour) {
			b.WriteString("0h")
		}
		if ns < int64(time.Minute) {
			b.WriteString("0m")
		}
	}

	b.WriteString(time.Duration(ns).String())

	return b.String()
}
