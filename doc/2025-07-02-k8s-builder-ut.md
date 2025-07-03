This document outlines the unit tests for functions in `builder.go` in the `k8s` module.

# Notes

* Implement as table-based test cases

# Test Cases for `Substitute()`

1. Basic alias substitution
2. Multiple substitutions
3. No matching substitutions
4. First substitution wins -- when multiple substitutions share the same alias, the first
   one wins
5. Multiple substitutions in the same argument

# Test Cases for `BuildKubectlArgs()`

1.  **Basic command**: A simple command with no special flags.
2.  **With context and namespace**: The command includes context and namespace flags.
3.  **All namespaces**: The command includes the `--all-namespaces` flag when `allNamespaces` is true.
4.  **Assume cluster admin**: When `assumeClusterAdmin` is true, the `--as=compute:cluster-admin` flag is added to mutating commands (e.g., `create`, `edit`, `delete`).
5.  **Confirmation for mutating commands**: The user is prompted for confirmation when running a mutating command (e.g., `create`, `delete`, `edit`, `apply`, `patch`, `replace`, `scale`).
6.  **Rewrites resource for relevant command and resource**: For commands that operate on a single resource (e.g., `logs`, `exec`, `describe`, `edit`), rewrite the resource name (e.g., `logs job foo` becomes `logs job/foo`).
6.  **Doesn't modify qualified / resource**: If the resource is already qualified (e.g., `job/foo`), it is not modified.
7.  **Empty exec assumes `-it -- bash`**: An `exec` command with no specified command defaults to `exec -it <pod> -- bash`.
8.  **Empty exec for rewritable resource**: If empty exec is passed for a rewritable resource (ex. `job foo`), then the builder should also assume `-it -- bash`.
9.  **Empty exec for rewritable resource**: An `exec` command for a rewritable resource (ex. `job/foo`) defaults to `exec -it <pod> -- bash`.
10.  **Exec with args**
11.  **Empty arguments**: No arguments are provided.

# Test Cases for `BuildK9sArgs()`

1.  Basic command
2.  With context and namespace
3.  All namespaces
4.  Empty args defaults to pods
5.  Too many args
6.  Get with resource
7.  Resource only
