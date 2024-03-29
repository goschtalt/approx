// SPDX-FileCopyrightText: 2009 The Go Authors.
// SPDX-FileCopyrightText: 2023 Weston Schmidt <weston_schmidt@alumni.purdue.edu>
// SPDX-License-Identifier: BSD-3-Clause

package approx_test

import (
	"errors"
	"math"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gerifield/approx"
)

func TestParseDuration_Mine(t *testing.T) {
	unknownErr := errors.New("unknown")
	tests := []struct {
		in     string
		expect time.Duration
		err    error
	}{
		{
			in:     "1.012h",
			expect: time.Hour + 43*time.Second + 200*time.Millisecond,
		}, {
			in:     "1y0w0d0h0m0s",
			expect: approx.Year,
		}, {
			in:     "0y1M0w0d0h0m0s",
			expect: approx.Month,
		}, {
			in:     "1y0w0d0h0m0s",
			expect: 52*approx.Week + approx.Day,
		}, {
			in:     "4y0w0d0h0m0s",
			expect: 4 * approx.Year,
		}, {
			in:     "4y3w2d1h0m0s",
			expect: 4*approx.Year + 3*approx.Week + 2*approx.Day + 1*time.Hour,
		}, {
			in:     "3w0d0h0m0s",
			expect: 3 * approx.Week,
		}, {
			in:     "3w0d0h0m1s",
			expect: 3*approx.Week + time.Second,
		}, {
			in:     "-3w0d0h0m1s",
			expect: -1 * (3*approx.Week + time.Second),
		}, {
			in:     "1h0m0s",
			expect: 60 * time.Minute,
		}, {
			in:     "2M5m",
			expect: 2*approx.Month + 5*time.Minute,
		},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			got, err := approx.ParseDuration(tc.in)

			if tc.err != nil {
				if errors.Is(tc.err, unknownErr) {
					if err == nil {
						t.Logf("expected and error but it was nil")
						t.Fail()
					}
					return
				}

				if !errors.Is(err, tc.err) {
					t.Logf("expected error %s but got %s", reflect.TypeOf(tc.err), reflect.TypeOf(err))
					t.Fail()
				}

				return
			}

			if tc.expect != got {
				t.Logf("expected: %q got: %q\n", tc.expect, got)
				t.Fail()
			}
		})
	}
}

// -- the following tests are from golang stl ----------------------------------
// https://github.com/golang/go/blob/go1.20.5/src/time/time_test.go

var durationTests = []struct {
	str string
	d   time.Duration
}{
	{"0s", 0},
	{"1ns", 1 * time.Nanosecond},
	{"1.1µs", 1100 * time.Nanosecond},
	{"2.2ms", 2200 * time.Microsecond},
	{"3.3s", 3300 * time.Millisecond},
	{"4m5s", 4*time.Minute + 5*time.Second},
	{"4m5.001s", 4*time.Minute + 5001*time.Millisecond},
	{"5h6m7.001s", 5*time.Hour + 6*time.Minute + 7001*time.Millisecond},
	{"8m0.000000001s", 8*time.Minute + 1*time.Nanosecond},
	{"2562047h47m16.854775807s", 1<<63 - 1},
	{"-2562047h47m16.854775808s", -1 << 63},
}

func TestDurationString(t *testing.T) {
	for _, tt := range durationTests {
		if str := tt.d.String(); str != tt.str {
			t.Errorf("Duration(%d).String() = %s, want %s", int64(tt.d), str, tt.str)
		}
		if tt.d > 0 {
			if str := (-tt.d).String(); str != "-"+tt.str {
				t.Errorf("Duration(%d).String() = %s, want %s", int64(-tt.d), str, "-"+tt.str)
			}
		}
	}
}

var parseDurationTests = []struct {
	in   string
	want time.Duration
}{
	// simple
	{"0", 0},
	{"5s", 5 * time.Second},
	{"30s", 30 * time.Second},
	{"1478s", 1478 * time.Second},
	// sign
	{"-5s", -5 * time.Second},
	{"+5s", 5 * time.Second},
	{"-0", 0},
	{"+0", 0},
	// decimal
	{"5.0s", 5 * time.Second},
	{"5.6s", 5*time.Second + 600*time.Millisecond},
	{"5.s", 5 * time.Second},
	{".5s", 500 * time.Millisecond},
	{"1.0s", 1 * time.Second},
	{"1.00s", 1 * time.Second},
	{"1.004s", 1*time.Second + 4*time.Millisecond},
	{"1.0040s", 1*time.Second + 4*time.Millisecond},
	{"100.00100s", 100*time.Second + 1*time.Millisecond},
	// different units
	{"10ns", 10 * time.Nanosecond},
	{"11us", 11 * time.Microsecond},
	{"12µs", 12 * time.Microsecond}, // U+00B5
	{"12μs", 12 * time.Microsecond}, // U+03BC
	{"13ms", 13 * time.Millisecond},
	{"14s", 14 * time.Second},
	{"15m", 15 * time.Minute},
	{"16h", 16 * time.Hour},
	// composite durations
	{"3h30m", 3*time.Hour + 30*time.Minute},
	{"10.5s4m", 4*time.Minute + 10*time.Second + 500*time.Millisecond},
	{"-2m3.4s", -(2*time.Minute + 3*time.Second + 400*time.Millisecond)},
	{"1h2m3s4ms5us6ns", 1*time.Hour + 2*time.Minute + 3*time.Second + 4*time.Millisecond + 5*time.Microsecond + 6*time.Nanosecond},
	{"39h9m14.425s", 39*time.Hour + 9*time.Minute + 14*time.Second + 425*time.Millisecond},
	// large value
	{"52763797000ns", 52763797000 * time.Nanosecond},
	// more than 9 digits after decimal point, see https://golang.org/issue/6617
	{"0.3333333333333333333h", 20 * time.Minute},
	// 9007199254740993 = 1<<53+1 cannot be stored precisely in a float64
	{"9007199254740993ns", (1<<53 + 1) * time.Nanosecond},
	// largest duration that can be represented by int64 in nanoseconds
	{"9223372036854775807ns", (1<<63 - 1) * time.Nanosecond},
	{"9223372036854775.807us", (1<<63 - 1) * time.Nanosecond},
	{"9223372036s854ms775us807ns", (1<<63 - 1) * time.Nanosecond},
	{"-9223372036854775808ns", -1 << 63 * time.Nanosecond},
	{"-9223372036854775.808us", -1 << 63 * time.Nanosecond},
	{"-9223372036s854ms775us808ns", -1 << 63 * time.Nanosecond},
	// largest negative value
	{"-9223372036854775808ns", -1 << 63 * time.Nanosecond},
	// largest negative round trip value, see https://golang.org/issue/48629
	{"-2562047h47m16.854775808s", -1 << 63 * time.Nanosecond},
	// huge string; issue 15011.
	{"0.100000000000000000000h", 6 * time.Minute},
	// This value tests the first overflow check in leadingFraction.
	{"0.830103483285477580700h", 49*time.Minute + 48*time.Second + 372539827*time.Nanosecond},
}

func TestParseDuration(t *testing.T) {
	for _, tc := range parseDurationTests {
		d, err := approx.ParseDuration(tc.in)
		if err != nil || d != tc.want {
			t.Errorf("ParseDuration(%q) = %v, %v, want %v, nil", tc.in, d, err, tc.want)
		}
	}
}

var parseDurationErrorTests = []struct {
	in     string
	expect string
}{
	// invalid
	{"", `""`},
	{"3", `"3"`},
	{"-", `"-"`},
	{"s", `"s"`},
	{".", `"."`},
	{"-.", `"-."`},
	{".s", `".s"`},
	{"+.s", `"+.s"`},
	// {"1d", `"1d"`}, WTS This is valid now
	{"\x85\x85", `"\x85\x85"`},
	{"\xffff", `"\xffff"`},
	{"hello \xffff world", `"hello \xffff world"`},
	{"\uFFFD", `"\xef\xbf\xbd"`},                                             // utf8.RuneError
	{"\uFFFD hello \uFFFD world", `"\xef\xbf\xbd hello \xef\xbf\xbd world"`}, // utf8.RuneError
	// overflow
	{"9223372036854775810ns", `"9223372036854775810ns"`},
	{"9223372036854775808ns", `"9223372036854775808ns"`},
	{"-9223372036854775809ns", `"-9223372036854775809ns"`},
	{"9223372036854776us", `"9223372036854776us"`},
	{"3000000h", `"3000000h"`},
	{"9223372036854775.808us", `"9223372036854775.808us"`},
	{"9223372036854ms775us808ns", `"9223372036854ms775us808ns"`},
}

func TestParseDurationErrors(t *testing.T) {
	for _, tc := range parseDurationErrorTests {
		_, err := approx.ParseDuration(tc.in)
		if err == nil {
			t.Errorf("ParseDuration(%q) = _, nil, want _, non-nil", tc.in)
		} else if !strings.Contains(err.Error(), tc.expect) {
			t.Errorf("ParseDuration(%q) = _, %q, error does not contain %q", tc.in, err, tc.expect)
		}
	}
}

func TestParseDurationRoundTrip(t *testing.T) {
	// https://golang.org/issue/48629
	max0 := time.Duration(math.MaxInt64)
	max1, err := approx.ParseDuration(max0.String())
	if err != nil || max0 != max1 {
		t.Errorf("round-trip failed: %d => %q => %d, %v", max0, max0.String(), max1, err)
	}

	min0 := time.Duration(math.MinInt64)
	min1, err := approx.ParseDuration(min0.String())
	if err != nil || min0 != min1 {
		t.Errorf("round-trip failed: %d => %q => %d, %v", min0, min0.String(), min1, err)
	}

	for i := 0; i < 100; i++ {
		// Resolutions finer than milliseconds will result in
		// imprecise round-trips.
		d0 := time.Duration(rand.Int31()) * time.Millisecond
		s := d0.String()
		d1, err := approx.ParseDuration(s)
		if err != nil || d0 != d1 {
			t.Errorf("round-trip failed: %d => %q => %d, %v", d0, s, d1, err)
		}
	}
}
