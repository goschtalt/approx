# approx

[![Build Status](https://github.com/goschtalt/approx/actions/workflows/ci.yml/badge.svg)](https://github.com/goschtalt/approx/actions/workflows/ci.yml)
[![codecov.io](http://codecov.io/github/goschtalt/approx/coverage.svg?branch=main)](http://codecov.io/github/goschtalt/approx?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/goschtalt/approx)](https://goreportcard.com/report/github.com/goschtalt/approx)
[![GitHub Release](https://img.shields.io/github/release/goschtalt/approx.svg)](https://github.com/goschtalt/approx/releases)
[![GoDoc](https://pkg.go.dev/badge/github.com/goschtalt/approx)](https://pkg.go.dev/github.com/goschtalt/approx)

Package approx adds support for durations of days, weeks, months and years.  The code
used is from the go standard library.  Only very minor adjustments were made
to enable parsing to support extra units of time.

# Usage

This really simple library allows you to normally use the `time.Duration` object
from the go standard library, but suppliments two functions for handling
approximate time durations.

```golang
package main

import (
	"fmt"

	"github.com/gerifield/approx"
)

func main() {
	d, err := approx.ParseDuration("1w4d")
	if err != nil {
		panic(err)
	}

	fmt.Println(approx.String(d))
	fmt.Println(d)
}
```

[Go Playground](https://go.dev/play/p/zX2FeTrC8Qb)
