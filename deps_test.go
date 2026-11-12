// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This test enforces invariants of the import dependency graph.

package tools_test

import (
	"os/exec"
	"slices"
	"strings"
	"testing"

	"github.com/greatliontech/go-x-tools/internal/testenv"
)

// rules lists "allow" or "deny" constraints on the transitive
// dependencies among non-test packages in x/tools.
var rules = [...]struct {
	rationale string
	from      []string // import paths or patterns (e.g. "github.com/greatliontech/go-x-tools/...") of constrained packages
	except    []string // import paths or patterns of packages matched by from that are not constrained
	deny      []string // if set, packages that none of them may depend on
	allow     []string // if set, the only packages of this module that they may depend on
}{
	{
		rationale: "syntax packages must not depend on the type checker",
		from: []string{
			"github.com/greatliontech/go-x-tools/go/ast/astutil",
			"github.com/greatliontech/go-x-tools/go/ast/edge",
			"github.com/greatliontech/go-x-tools/go/ast/inspector",
			"github.com/greatliontech/go-x-tools/internal/astutil",
			"github.com/greatliontech/go-x-tools/internal/cursorutil",
		},
		deny: []string{"go/types"},
	},
	{
		rationale: "cursorutil must have no dependencies other than those of cursor itself",
		from:      []string{"github.com/greatliontech/go-x-tools/internal/cursorutil"},
		allow: []string{
			"github.com/greatliontech/go-x-tools/go/ast/edge",
			"github.com/greatliontech/go-x-tools/go/ast/inspector",
		},
	},
	{
		rationale: "gcexportdata is vendored into std for use by go/importer, so must have minimal dependencies",
		from:      []string{"github.com/greatliontech/go-x-tools/go/gcexportdata"},
		allow: []string{
			"github.com/greatliontech/go-x-tools/go/types/objectpath",
			"github.com/greatliontech/go-x-tools/internal/gcimporter",
			"github.com/greatliontech/go-x-tools/internal/pkgbits",
		},
	},
	{
		rationale: "packages vendored into cmd must have minimal dependencies",
		// These are the packages imported by cmd (see
		// $GOROOT/src/cmd/vendor/modules.txt). A violation would be
		// discovered only when cmd next updates its vendored copy.
		from: []string{
			"github.com/greatliontech/go-x-tools/cmd/bisect",
			"github.com/greatliontech/go-x-tools/cover",
			"github.com/greatliontech/go-x-tools/go/analysis/suite/fix",
			"github.com/greatliontech/go-x-tools/go/analysis/suite/vet",
			"github.com/greatliontech/go-x-tools/go/analysis/unitchecker",
		},
		deny: []string{
			"github.com/greatliontech/go-x-tools/go/packages",
			"github.com/greatliontech/go-x-tools/go/ssa",
			"github.com/greatliontech/go-x-tools/internal/event",
			"github.com/greatliontech/go-x-tools/internal/gocommand",
			"testing",
		},
	},
	{
		rationale: "the analysis packages vendored into cmd (and, prospectively, std) must not depend on x/sync, so that std need not vendor it",
		from: []string{
			"github.com/greatliontech/go-x-tools/go/analysis/suite/fix",
			"github.com/greatliontech/go-x-tools/go/analysis/suite/vet",
			"github.com/greatliontech/go-x-tools/go/analysis/unitchecker",
		},
		deny: []string{
			"golang.org/x/sync/errgroup",
			"golang.org/x/sync/semaphore",
			"golang.org/x/sync/singleflight",
			"golang.org/x/sync/syncmap",
		},
	},
	{
		rationale: "only go/ssa/ssautil may connect go/ssa and go/packages",
		from:      []string{"github.com/greatliontech/go-x-tools/go/ssa"},
		deny:      []string{"github.com/greatliontech/go-x-tools/go/packages"},
	},
	{
		rationale: "low-level type utilities must not depend on higher layers",
		// (Mostly implied by the cmd rule above, since these
		// packages are all vendored into cmd.)
		from: []string{
			"github.com/greatliontech/go-x-tools/go/types/objectpath",
			"github.com/greatliontech/go-x-tools/go/types/typeutil",
			"github.com/greatliontech/go-x-tools/internal/typeparams",
			"github.com/greatliontech/go-x-tools/internal/typesinternal",
		},
		deny: []string{
			"github.com/greatliontech/go-x-tools/go/analysis",
			"github.com/greatliontech/go-x-tools/go/packages",
			"github.com/greatliontech/go-x-tools/go/ssa",
		},
	},
	{
		rationale: "only test helpers may depend on test helpers",
		from:      []string{"github.com/greatliontech/go-x-tools/..."},
		except: []string{
			"github.com/greatliontech/go-x-tools/go/analysis/analysistest",
			"github.com/greatliontech/go-x-tools/go/packages/packagestest",
			"github.com/greatliontech/go-x-tools/internal/packagestest",
			"github.com/greatliontech/go-x-tools/internal/testfiles",
		},
		deny: []string{
			"github.com/greatliontech/go-x-tools/go/analysis/analysistest",
			"github.com/greatliontech/go-x-tools/go/packages/packagestest",
			"github.com/greatliontech/go-x-tools/internal/diff/difftest",
			"github.com/greatliontech/go-x-tools/internal/event/export/eventtest",
			"github.com/greatliontech/go-x-tools/internal/jsonrpc2/stack/stacktest",
			"github.com/greatliontech/go-x-tools/internal/packagestest",
			"github.com/greatliontech/go-x-tools/internal/testenv",
			"github.com/greatliontech/go-x-tools/internal/testfiles",
		},
	},
}

func TestDeps(t *testing.T) {
	testenv.NeedsTool(t, "go")

	goList := func(args ...string) []string {
		args = append([]string{"list"}, args...)
		out, err := testenv.Command(t, "go", args...).Output()
		if err != nil {
			if err, ok := err.(*exec.ExitError); ok {
				t.Fatalf("go %s: %v\n%s", strings.Join(args, " "), err, err.Stderr)
			}
			t.Fatalf("go %s: %v", strings.Join(args, " "), err)
		}
		return strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	}

	// expand expands dir/... import path patterns within this module.
	expand := func(patterns []string) []string {
		if slices.ContainsFunc(patterns, func(s string) bool { return strings.Contains(s, "...") }) {
			const filter = `{{if and .Module (eq .Module.Path "github.com/greatliontech/go-x-tools")}}{{.ImportPath}}{{end}}`
			pkgs := goList(append([]string{"-f", filter}, patterns...)...)
			return slices.DeleteFunc(pkgs, func(pkg string) bool { return pkg == "" })
		}
		return patterns
	}

	// Expand patterns and exceptions in each rule's from list.
	// In doing, build the set of roots for the graph query.
	froms := make([][]string, len(rules)) // froms[i] is the expansion of depRules[i].from
	var roots []string                    // may contain dups
	for i, rule := range rules {
		if (rule.deny == nil) == (rule.allow == nil) {
			t.Fatalf("invalid rule for %v: must set exactly one of deny or allow", rule.from)
		}
		except := expand(rule.except)
		from := slices.DeleteFunc(slices.Clone(expand(rule.from)), func(pkg string) bool {
			return slices.Contains(except, pkg)
		})
		froms[i] = from
		roots = append(roots, from...)
	}

	// Compute the import graph of all constrained packages.
	imports := make(map[string][]string)
	for _, line := range goList(append([]string{"-deps", "-f", "{{.ImportPath}}{{range .Imports}} {{.}}{{end}}"}, roots...)...) {
		// line has form "from to1 ... toN"
		if fields := strings.Fields(line); len(fields) > 0 {
			imports[fields[0]] = fields[1:]
		}
	}

	// reaches returns an arbitrary path from x to y (inclusive)
	// through the import graph, or nil if there is none.
	// The path contains at least one import edge, so a package
	// does not trivially reach itself.
	reaches := func(x, y string) []string {
		var (
			seen = make(map[string]bool)
			path []string // reversed: y first
		)
		// importsReachY reports whether some import of pkg is, or reaches, y.
		var importsReachY func(pkg string) bool
		importsReachY = func(pkg string) bool {
			if !seen[pkg] {
				seen[pkg] = true
				if slices.ContainsFunc(imports[pkg], func(imp string) bool {
					if imp == y {
						path = append(path, y)
						return true
					}
					return importsReachY(imp)
				}) {
					path = append(path, pkg)
					return true
				}
			}
			return false
		}
		if importsReachY(x) {
			slices.Reverse(path)
		}
		return path
	}

	for i, rule := range rules {
		for _, from := range froms[i] {
			for _, deny := range rule.deny {
				if path := reaches(from, deny); path != nil {
					t.Errorf("%s must not depend on %s (%s), but does:\n\t%s",
						from, deny, rule.rationale, strings.Join(path, " ->\n\t"))
				}
			}
			if rule.allow != nil {
				for dep := range imports {
					if strings.HasPrefix(dep, "github.com/greatliontech/go-x-tools/") && !slices.Contains(rule.allow, dep) {
						if path := reaches(from, dep); path != nil {
							t.Errorf("%s has disallowed dependency on %s (%s):\n\t%s",
								from, dep, rule.rationale, strings.Join(path, " ->\n\t"))
						}
					}
				}
			}
		}
	}
}
