Implement a custom `workspace` command that provides utilities around setting up projects
via git worktrees.

Set it up initially with the `new` subcommand, ie. `mdcli workspace new ...`

`new` subcommand takes a single argument, the git URL to clone from.

In addition, it supports the following flags:

* `--name $NAME`: the name of the workspace/project, defaults to the name of the repo
* `--project-dir $PROJECT_DIR`: the directory to create the workspace in
* `--git-branch $BRANCH`: the git ref to checkout, defaults to `main`
* `--initialize-worktree`: whether to create a worktree for the project, defaults to true
* `--default-worktree-name $DEFAULT_WORKTREE_NAME`: the name of the default worktree, defaults to $BRANCH

When run, the command will:

1. Initialize the workspace directory, ie. `mkdir -p $PROJECT_DIR/$NAME`
2. Perform a bare clone of the repo to the project directory, ie. `git clone --bare $URL $PROJECT_DIR/$NAME/.git`
3. If `--initialize-worktree` is set, create a worktree for the project in the workspace directory, ie. `git worktree add $PROJECT_DIR/$NAME/worktrees/$DEFAULT_WORKTREE_NAME $BRANCH`

Worktrees will always be in the following layout:

$PROJECT_DIR/
    $NAME/
        worktrees/
            $DEFAULT_WORKTREE_NAME/
            # other worktrees ...
