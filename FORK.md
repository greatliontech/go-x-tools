# Great Lion Technologies analysis fork

This fork originates from [Go Tools](https://go.googlesource.com/tools)
(`https://github.com/golang/tools.git`). Upstream copyright and license
notices remain in force.

- **Current upstream anchor:** `d4453b555` (`v0.52.0-90-gd4453b555`)
- **Active fork branch:** `main`
- **Upstream tracking branch:** `master` (tracks upstream `master` unchanged)
- **Archived fork branches:** `glt-v0.50` (`v0.50.0`-based), `glt-v0.51` (`v0.51.0`-based)

## Module identity, toolchain target, and consumer pinning

The root module is **`github.com/greatliontech/go-x-tools`**. Consumers
import that path and require an explicit published `glt` revision, allowing
the patched dependency to propagate transitively under Go Minimum Version
Selection. A downstream `replace` directive is not a substitute: replacements
in dependencies are ignored by `go`.

After a fork revision is published on `main`:

```sh
go get github.com/greatliontech/go-x-tools@<glt-commit>
```

Keep the resulting exact pseudo-version in `go.mod`. Upstream tags still
declare `golang.org/x/tools`; do not select them through the fork path or use
`@latest` to choose the patched dependency.

The fork's scope is **Go 1.27+ V5 unified export-data analyses**, concurrency
safety, and startup/dynamic-state purity. CI tests the root module on Go
`1.27.1` and `1.27.x`; root `go.mod` retains upstream's `go 1.26.0` directive
unchanged. It is not a general standalone-tools or `gopls` distribution.

The separate `gopls` and `gopls/doc/assets` modules retain their upstream
identities, imports, and dependency declarations. They are not migrated or
supported by this fork. Upstream documentation URLs, standalone-tool
installation examples (`@latest`), and synthetic `testdata` module identities
are retained unchanged.

## Fork patch inventory (`<anchor>..main`)

`main` maintains a strictly ordered linear stack on top of the upstream
anchor: **all semantic commits sit directly on `<anchor>` using stock
`"golang.org/x/tools/..."` imports, and the single top commit (`HEAD`)
applies the mechanical `github.com/greatliontech/go-x-tools` module identity
and fork files.** Keeping stock imports on the semantic commits eliminates
import-block conflicts during upstream rebases and keeps each semantic commit
independently droppable once upstream fixes the underlying issue.

Current semantic commits (bottom to top, beneath `HEAD`):

1. **`internal/typesinternal`, `internal/jsonrpc2/stack`** — Normalize
   `types.Unalias(T)` before `ForEachElement` invokes its caller's visited
   callback (preventing identity-keyed visited sets in RTA and SSA from
   pruning structural descendants of aliased types), fix negative expectation
   lookup in `element_test.go`, and allow hyphens in `stack/parse.go` import
   paths. *Drop individual parts only if upstream normalizes aliases before
   `f(T, access)` in `ForEachElement` and admits `-` in `reCall`.*
2. **`go/ssa`, `go/packages`, `go/callgraph/rta`, `internal/gcimporter`,
   `internal/gocommand`, `internal/packagesinternal`, `internal/typeparams`** —
   Eliminate package-level dynamic carriers (`allTermlist`, `tRune`,
   predeclared globals, event label keys, `packagesinternal` function vars)
   and thread loader context. *Required so downstream analyses observe zero
   undischarged package-level dynamic state.*
3. **`go/packages`** — Eagerly unpack `*types.Named` (`Underlying()`) and
   `*types.Alias` (`types.Unalias()`) immediately after shallow export-data
   import, before releasing concurrent dependent package type-checkers. *Drop
   only if upstream eliminates concurrent lazy `Named`/`Alias` unpacking
   races across dependent packages.*
4. **`go/packages`, `go/ssa`, `go/types/typeutil`** — Evaluate
   `GOPACKAGESDEBUG` on demand instead of package `init`, replace concrete
   `reflect.TypeFor` calls in `go/ssa/ssa.go` with `reflect.TypeOf` on zero
   values, and replace `maphash.MakeSeed()` in `typeutil.Hasher` with
   deterministic `TypeName` hashing. *Required for startup purity and
   deterministic hashing.*
5. **`internal/gcimporter`** — Store value-only `gorootExportEntry` structs
   instead of closures in `exportMap` so `sync.Map` memoization carries no
   dynamic function state.
6. **`go/ast/edge`** — Replace the package-level `fieldInfos` table (whose
   `reflect.Type` interface fields look writable to static escape analysis)
   and `reflect.TypeFor` calls with a pure `Kind.info()` switch.

*(Previously carried `internal/gocommand` process-group cancellation was
absorbed upstream in `1332450b8` / `v0.52.0` and dropped.)*

## Upstream sync workflow (LLM / maintainer runbook)

Never `git merge` upstream into `main`: merging collides on rewritten import
blocks while silently admitting newly added upstream `.go` files that still
import `"golang.org/x/tools/..."`. Always archive the outgoing revision if
published, rebase the stock-import semantic commits onto the new upstream
anchor, and regenerate the single top module-identity commit.

### 1. Fetch upstream and audit changes

Ensure `main` is checked out and the `upstream` remote points to
`https://github.com/golang/tools.git`:

```sh
git checkout main
git fetch --tags upstream master
git branch -f master upstream/master
ANCHOR=$(git merge-base main upstream/master)
```

Inspect both directions before rebasing:

1. **Have any of our semantic patches been fixed upstream?** Query upstream
   commits against the exact file set touched by our semantic commits
   (`$ANCHOR..main~1`) — do not rely solely on `git patch-id`, since upstream
   may fix the same bug with different wording:
   ```sh
   git log -p "$ANCHOR..upstream/master" -- $(git diff --name-only "$ANCHOR..main~1")
   ```
2. **What new upstream changes affect consumers or startup/state purity?**
   Review root-module commits (`git log --oneline "$ANCHOR..upstream/master" -- . ':!gopls'`)
   and check non-test root `.go` diffs under `go/`, `internal/`, `imports/`,
   `cover/`, and `refactor/` for any newly added package-level `var`,
   `func init()`, `os.Getenv`, `maphash.MakeSeed`, or `reflect.TypeFor`
   that would violate downstream dynamic-state / startup-purity gates.

### 2. Archive the outgoing `main` tip

Before rewriting `main`, preserve any published revision pinned by downstream
`go.mod` files on an immutable `glt-v0.<minor>` branch (matching the outgoing
anchor's `v0.<minor>` tag series) so GitHub never garbage-collects those
commits:

```sh
# Example when moving off a v0.52-series revision:
git branch glt-v0.52 main
```

Do not overwrite an existing `glt-v0.*` branch once pushed unless intentionally
updating within that exact series before consumers pin a newer base.

### 3. Rebase the stock-import semantic commits

Select `<TARGET>` (either the latest upstream `v0.*` tag or `upstream/master`).
Rebase only the semantic commits (`$ANCHOR..$OLD_MAIN~1`) in detached `HEAD`
first, moving `main` only after the rebase completes cleanly (so `git rebase
--abort` or a failed command leaves `main` and `$OLD_MAIN` intact):

```sh
OLD_MAIN=$(git rev-parse main)
git rebase --onto <TARGET> "$ANCHOR" "$OLD_MAIN~1"
git checkout -B main HEAD
```

- Drop any semantic commit whose fix landed upstream.
- Resolve any code conflicts while keeping `"golang.org/x/tools/..."` import
  paths in all semantic commits (never introduce `github.com/greatliontech/go-x-tools`
  below the top commit).

### 4. Regenerate the single top module-identity commit

On top of the rebased semantic stack:

1. Restore fork-owned files from `$OLD_MAIN`:
   ```sh
   git checkout "$OLD_MAIN" -- .github/workflows/fork.yml FORK.md fork_identity_test.go
   ```
2. Update `- **Current upstream anchor:**` at the top of `FORK.md` (and the
   patch inventory if a commit was dropped or added), and re-apply the 4-line
   fork notice blockquote at the top of upstream's current `README.md`.
3. Deterministically rewrite root-module self-references from
   `golang.org/x/tools` to `github.com/greatliontech/go-x-tools`:
   - Rewrite `module golang.org/x/tools` in root `go.mod` only (never nested
     `gopls/go.mod` or `**/testdata/**`).
   - Rewrite `cmd/bundle/testdata/out.golden`'s generated-by header.
   - In all root-module `.go` files outside `gopls/`, `**/testdata/**`, and
     `fork_identity_test.go`, rewrite `golang.org/x/tools` to
     `github.com/greatliontech/go-x-tools` **except**:
     - occurrences preceded by `[A-Za-z0-9_.-]` or followed immediately by
       `[A-Za-z0-9_-]`,
     - `http://` and `https://` URLs (e.g. `https://pkg.go.dev/golang.org/x/tools/...`),
     - `golang.org/x/tools/gopls` paths (where `/gopls` is not followed by
       `[A-Za-z0-9_-]`),
     - `@latest` standalone-tool examples (`go install` / `go run ...@latest`).
   - In `go/buildutil/util_test.go`, also rewrite
     `"golang.org", "x", "tools"` to `"github.com", "greatliontech", "go-x-tools"`.
   - Run `ast.SortImports` + `gofmt` across all modified `.go` files so
     `"github.com/..."` imports sort canonically above `"golang.org/x/..."`.
4. Commit all module-identity and fork-doc changes as the single top commit:
   `chore(fork): apply github.com/greatliontech/go-x-tools module identity`.

### 5. Verify before pushing

1. **Root test suite (includes `TestForkModuleIdentity`, `TestDeps`, and
   `TestToolsCopyright`):**
   ```sh
   go test -p 4 -count=1 -timeout=20m ./...
   ```
   `TestForkModuleIdentity` (`fork_identity_test.go`) verifies that root
   `go.mod` declares `github.com/greatliontech/go-x-tools`, zero root-module
   `.go` files import `golang.org/x/tools`, and `//go:linkname` directives
   use the fork path.
2. **Downstream `gofresh` dynamic-state check:** When `gofresh` is checked out
   locally alongside this repository, verify
   `TestModulePackagesHaveNoSharedDynamicStateUnderPackageExecution` against
   the local `go-x-tools` tree inside a subshell-scoped temporary worktree
   (leaving `gofresh`'s checkout and outer shell traps untouched):
   ```sh
   (
     set -eu
     XTOOLS_DIR=$(pwd)
     GOFRESH_DIR=$(cd ../gofresh && pwd)
     TMPW=$(mktemp -d)
     trap 'git -C "$GOFRESH_DIR" worktree remove --force "$TMPW" 2>/dev/null || true; rm -rf "$TMPW"' EXIT
     git -C "$GOFRESH_DIR" worktree add --detach "$TMPW" HEAD
     cd "$TMPW"
     go mod edit -replace="github.com/greatliontech/go-x-tools=$XTOOLS_DIR"
     go mod tidy
     go test -count=1 -timeout=10m -run='^TestModulePackagesHaveNoSharedDynamicStateUnderPackageExecution$' .
   )
   ```
3. **Publish:** Push `master`, any new `glt-v0.*` archive branch, and
   `--force-with-lease` on `main`, then update downstream `go.mod` pins with
   `go get github.com/greatliontech/go-x-tools@<new-main-sha>`.
