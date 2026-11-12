// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tools_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestForkModuleIdentity(t *testing.T) {
	const (
		upstreamMod = "golang.org/" + "x/tools"
		forkMod     = "github.com/greatliontech/go-x-tools"
	)

	goModBytes, err := os.ReadFile("go.mod")
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	hasModLine := false
	for line := range strings.SplitSeq(string(goModBytes), "\n") {
		if rest, ok := strings.CutPrefix(strings.TrimSpace(line), "module "); ok {
			hasModLine = strings.TrimSpace(rest) == forkMod
			break
		}
	}
	if !hasModLine {
		t.Fatalf("go.mod does not declare 'module %s'", forkMod)
	}

	err = filepath.WalkDir(".", func(fullPath string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel := filepath.ToSlash(fullPath)
		if rel == "." {
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") || rel == "gopls" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(rel, ".go") || strings.HasPrefix(rel, "testdata/") || strings.Contains(rel, "/testdata/") {
			return nil
		}

		src, readErr := os.ReadFile(fullPath)
		if readErr != nil {
			return readErr
		}
		fset := token.NewFileSet()
		f, parseErr := parser.ParseFile(fset, fullPath, src, parser.ImportsOnly|parser.ParseComments)
		if parseErr != nil {
			t.Errorf("parse %s: %v", rel, parseErr)
			return nil
		}
		for _, imp := range f.Imports {
			impPath, unqErr := strconv.Unquote(imp.Path.Value)
			if unqErr != nil {
				continue
			}
			if impPath == upstreamMod || (strings.HasPrefix(impPath, upstreamMod+"/") && !strings.HasPrefix(impPath, upstreamMod+"/gopls")) {
				t.Errorf("%s: imports %q; want %s", rel, impPath, forkMod)
			}
		}
		for line := range strings.SplitSeq(string(src), "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed == "//go:linkname" || strings.HasPrefix(trimmed, "//go:linkname ") || strings.HasPrefix(trimmed, "//go:linkname\t") {
				for _, field := range strings.Fields(trimmed)[1:] {
					if field == upstreamMod ||
						strings.HasPrefix(field, upstreamMod+".") ||
						(strings.HasPrefix(field, upstreamMod+"/") && !strings.HasPrefix(field, upstreamMod+"/gopls")) {
						t.Errorf("%s: unrewritten //go:linkname target: %s", rel, trimmed)
						break
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk root module: %v", err)
	}
}
