// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The vet package defines the suite of analyzers used by cmd/vet,
// the default analysis tool run by "go vet".
// Its behavior is equivalent to:
//
//	func main() { unitchecker.Main(vet.Suite...) }
//
// If you need a different suite, define your own tool
// and run "go vet -vettool=mytool".
package vet

import (
	"github.com/greatliontech/go-x-tools/go/analysis"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/appends"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/asmdecl"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/assign"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/atomic"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/bools"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/buildtag"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/cgocall"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/composite"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/copylock"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/defers"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/directive"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/errorsas"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/framepointer"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/hostport"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/httpresponse"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/ifaceassert"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/loopclosure"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/lostcancel"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/nilfunc"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/printf"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/scannererr"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/shift"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/sigchanyzer"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/slog"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/sqlrowserr"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/stdmethods"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/stdversion"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/stringintconv"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/structtag"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/testinggoroutine"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/tests"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/timeformat"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/unmarshal"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/unreachable"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/unsafeptr"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/unusedresult"
	"github.com/greatliontech/go-x-tools/go/analysis/passes/waitgroup"
)

// Suite is the suite of analyzers run by cmd/vet.
//
// The vet suite analyzers report diagnostics.
// (Diagnostics must describe real problems, but need not
// suggest fixes, and fixes are not necessarily safe to apply.)
var Suite = []*analysis.Analyzer{
	appends.Analyzer,
	asmdecl.Analyzer,
	assign.Analyzer,
	atomic.Analyzer,
	bools.Analyzer,
	buildtag.Analyzer,
	cgocall.Analyzer,
	composite.Analyzer,
	copylock.Analyzer,
	defers.Analyzer,
	directive.Analyzer,
	errorsas.Analyzer,
	// fieldalignment.Analyzer omitted: too noisy
	framepointer.Analyzer,
	httpresponse.Analyzer,
	hostport.Analyzer,
	ifaceassert.Analyzer,
	loopclosure.Analyzer,
	lostcancel.Analyzer,
	nilfunc.Analyzer,
	printf.Analyzer,
	scannererr.Analyzer,
	// shadow.Analyzer omitted: too noisy
	shift.Analyzer,
	sigchanyzer.Analyzer,
	slog.Analyzer,
	sqlrowserr.Analyzer,
	stdmethods.Analyzer,
	stdversion.Analyzer,
	stringintconv.Analyzer,
	structtag.Analyzer,
	tests.Analyzer,
	testinggoroutine.Analyzer,
	timeformat.Analyzer,
	unmarshal.Analyzer,
	unreachable.Analyzer,
	unsafeptr.Analyzer,
	unusedresult.Analyzer,
	waitgroup.Analyzer,
}
