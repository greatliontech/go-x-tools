// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gofix_test

import (
	"testing"

	"github.com/greatliontech/go-x-tools/go/analysis/analysistest"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/gofix"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, gofix.Analyzer, "a")
}
