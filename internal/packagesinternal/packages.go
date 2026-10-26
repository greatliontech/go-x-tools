// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package packagesinternal exposes internal-only fields from go/packages.
package packagesinternal

import "fmt"

type PackageError struct {
	ImportStack []string // shortest path from package named on command line to this one
	Pos         string   // position of error (if present, file:line:col)
	Err         string   // the error itself
}

func (err PackageError) String() string {
	return fmt.Sprintf("%s: %s (import stack: %s)", err.Pos, err.Err, err.ImportStack)
}

const (
	DepsErrors   = 1 << 10
	TypecheckCgo = 1 << 12
)
