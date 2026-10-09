// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stdversion_test

import (
	"path/filepath"
	"testing"

	"github.com/greatliontech/go-x-tools/go/analysis/analysistest"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/stdversion"
	"github.com/greatliontech/go-x-tools/internal/testfiles"
)

func Test(t *testing.T) {
	dir := testfiles.ExtractTxtarFileToTmp(t, filepath.Join(analysistest.TestData(), "test.txtar"))
	analysistest.Run(t, dir, stdversion.Analyzer,
		"example.com/basic",
		"example.com/despite",
		"example.com/mod20",
		"example.com/mod21",
		"example.com/mod22",
		"example.com/old",
	)
}
