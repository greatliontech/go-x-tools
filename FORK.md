# Great Lion Technologies analysis fork

This fork originates from [Go Tools](https://go.googlesource.com/tools),
based on upstream tag **v0.50.0**. The fork branch is `glt-v0.50`.
Upstream copyright and license notices remain in force.

The root module is **`github.com/greatliontech/go-x-tools`**. Consumers
import that path and require an explicit published `glt` revision, allowing
the patched dependency to propagate transitively. A downstream `replace`
directive is not a substitute: replacements in dependencies are ignored.
For example, after the fork revision is published:

```sh
go get github.com/greatliontech/go-x-tools@<glt-commit>
```

Keep the resulting exact version in `go.mod`. Upstream tags still declare
`golang.org/x/tools`; do not select them through the fork path or use
`@latest` to choose the patched dependency.

The fork's scope is V5-compatible export-data analyses and the associated
alias-traversal correctness correction. Its stack parser also accepts hyphens
in import paths, including the fork's own module path. It is not a general tools or gopls
distribution. Upstream documentation URLs, standalone-tool installation
examples, and contribution instructions are retained as source references.

The separate `gopls` and `gopls/doc/assets` modules retain their upstream
identities, imports, and dependency declarations. They are not migrated or
supported by this fork. In particular, `gopls/go.mod`'s upstream root-module
replacement is not a supported way to build gopls against this renamed
root; gopls would need its own coordinated module migration to import the
fork's internal packages. Root-module test fixtures retain their synthetic
module identities while their references to root packages use the fork.
