# General gitlab interaction notes

## Get Project Info

```bash
GL_PAT=$(get-gitlab-api-pat)
curl --request GET --header "PRIVATE-TOKEN: ${GL_PAT}" "https://git.alteryx.com/api/v4/projects/5784"
```

## Delete MR

```bash
GL_PAT=$(get-gitlab-api-pat)
curl --request DELETE --header "PRIVATE-TOKEN: ${GL_PAT}" "https://git.alteryx.com/api/v4/projects/5784/releases/v0.0.6"
```

## Local Testing

```mod
replace github.com/maahsome/gitlab-go v0.1.5 => /Users/christopher.maahs/Worktrees/GT-5/gitlab-go/get-groups
```
