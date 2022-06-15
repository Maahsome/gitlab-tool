# Local Build

Use this set of commands to perform a local build for tesing.

```bash
SEMVER=v0.0.999; echo ${SEMVER}
BUILD_DATE=$(gdate --utc +%FT%T.%3NZ); echo ${BUILD_DATE}
GIT_COMMIT=$(git rev-parse HEAD); echo ${GIT_COMMIT}

go build -ldflags "-X github.com/maahsome/gitlab-tool/cmd.semVer=${SEMVER} -X github.com/maahsome/gitlab-tool/cmd.buildDate=${BUILD_DATE} -X github.com/maahsome/gitlab-tool/cmd.gitCommit=${GIT_COMMIT} -X github.com/maahsome/gitlab-tool/cmd.gitRef=/refs/tags/${SEMVER}" && \
  cp ./gitlab-tool ~/tbin && \
  ./gitlab-tool version | jq .
```
