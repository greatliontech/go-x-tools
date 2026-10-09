// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stringintconv

import (
	_ "embed"
	"go/ast"

	"github.com/greatliontech/go-x-tools/go/analysis"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/inspect"
	"github.com/greatliontech/go-x-tools/go/ast/inspector"
	"github.com/greatliontech/go-x-tools/internal/analysis/analyzerutil"
	"github.com/greatliontech/go-x-tools/internal/analysis/intconv"
)

//go:embed doc.go
var doc string

var Analyzer = &analysis.Analyzer{
	Name:     "stringintconv",
	Doc:      analyzerutil.MustExtractDoc(doc, "stringintconv"),
	URL:      "https://pkg.go.dev/github.com/greatliontech/go-x-tools/go/analysis/passes/stringintconv",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

// The logic of this analyzer lives in the intconv package,
// which is shared with gopls.

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	for curCall := range inspect.Root().Preorder((*ast.CallExpr)(nil)) {
		if diag, ok := intconv.Check(pass.TypesInfo, curCall); ok {
			pass.Report(diag)
		}
	}
	return nil, nil
}
