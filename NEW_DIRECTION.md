# New Direction Idea

Have a config file that pins a directory PREFIX:

/Users/${USER}/src/alteryx_futurama = gitlab.com/alteryx_futurama
/Users/${USER}/dev/futurama/        = git.alteryx.com/futurama
/Users/${USER}/dev/falkor/          = git.alteryx.com/falkor

This way we no longer have to detect or GUESS, we just see if we are in one of these
directories, and if we are, then we use the defined information.

```bash
# the GROUP should/MUST be a top level group with no parent group
gitlab-tool config group --gitlab-host git.alteryx.com --group futurama --envvar GLA_TOKEN --directory ~/dev/futurama
```

## Design by example use

```bash
# while in a top level group or sub-group
gitlab-tool ls             # output a LS like listing, identifying GROUPS/PROJECTS
                           # this process could be in charge of mkdir for ALL groups at this level
ge-w bender
gc-- farnsworth
ge+- hermes
p--- test-project

-     : group or project (g/p)
 -    : exists on disk (e exists, c created)
  -   : dirty (- clean, + dirty)
   -  : has worktrees (-/w)
# tree from current location
gitlab-tool tree

.
├── bender
│   ├── docker-images
│   │   ├── create-futurama-repo (p)
│   │   └── futurama-gitlab-runner (p)
│   ├── gcp
│   │   ├── bootstrap (p)
│   │   ├── provisioning
│   │   │   ├── dev-lz (p)
│   │   │   └── dev-net (p)
│   │   └── tf-modules
│   │       └── secrets (p)
│   └── tools
├── farnsworth
│   ├── generate-temporary-project (p)
│   ├── simple-pipeline-test (p)
│   ├── terraform-net (p)
│   └── terraform-test (p)
└── hermes
    ├── automation
    │   ├── child-pipelines
    │   │   └── argocd-application-constructor (p)
    │   └── docker-images
    │       └── argocd-application-constructor (p)
    └── control-plane
        └── gcp
            └── lowers
                └── firefly-dev (p)

# clone a project
gitlab-tool clone test-project    # clones the project based on directory level (prefix/subdir depth)

# All other commands work as they do, just change the deterministic way they return details
# pay attention to worktree, as that would exist OUTSIDE of the normal directory, may have to read `.git` file to find prefix
```
