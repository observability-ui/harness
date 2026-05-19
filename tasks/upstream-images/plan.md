# Plan: Build and Push Upstream UIPlugin Images

## Context

The Observability Operator needs upstream images for its UIPlugins. Each plugin repository has specific release branches, and for each branch an image must be built with a specific version tag and pushed to `quay.io/openshift-observability-ui/`. The spec asks for a shell script to automate this across all 5 repos and 13 total image builds.

## What the script will do

Create a single shell script `build-and-push.sh` in `/Users/emurasak/workspace/harness/tasks/upstream-images/` that iterates over all repo/branch/tag combinations and for each one:

1. `cd` into the project directory (under `../../projects/<repo>`)
2. `git checkout` the release branch
3. Run `make podman-cross-build VERSION=<tag>` (which builds a multi-arch manifest and pushes it)
4. Return to the original directory

## Build matrix (from branches_tags.txt)

| Branch | Repo (directory) | Image name | Tag |
|--------|-----------------|------------|-----|
| release-coo-ocp-4.12 | console-dashboards-plugin | console-dashboards-plugin | v0.4.3 |
| release-coo-ocp-4.12 | distributed-tracing-console-plugin | distributed-tracing-console-plugin | v0.3.3 |
| release-coo-ocp-4.12 | logging-view-plugin | logging-view-plugin | v6.0.5 |
| release-coo-ocp-4.15 | distributed-tracing-console-plugin | distributed-tracing-console-plugin | v0.4.3 |
| release-coo-ocp-4.15 | logging-view-plugin | logging-view-plugin | v6.1.6 |
| release-coo-ocp-4.15 | monitoring-plugin | monitoring-console-plugin | v0.4.5 |
| release-coo-ocp-4.19 | distributed-tracing-console-plugin | distributed-tracing-console-plugin | v1.0.3 |
| release-coo-ocp-4.19 | troubleshooting-panel-console-plugin | troubleshooting-panel-console-plugin | v0.4.5 |
| release-coo-ocp-4.19 | monitoring-plugin | monitoring-console-plugin | v0.5.4 |
| release-coo-ocp-4.22 | distributed-tracing-console-plugin | distributed-tracing-console-plugin | v1.1.0 |
| release-coo-ocp-4.22 | troubleshooting-panel-console-plugin | troubleshooting-panel-console-plugin | v1.0.0 |
| release-coo-ocp-4.22 | monitoring-plugin | monitoring-console-plugin | v1.0.0 |
| release-coo-ocp-4.22 | logging-view-plugin | logging-view-plugin | v6.2.1 |

## Special handling for monitoring-plugin

Two differences from other repos:
1. **Image name differs from repo name:** The image is `monitoring-console-plugin` but the repo directory is `monitoring-plugin`. The Makefile has `PLUGIN_NAME ?= monitoring-plugin`, so we must pass `PLUGIN_NAME=monitoring-console-plugin` to override it.
2. **`podman-cross-build` does not push:** Unlike the other repos, monitoring-plugin's `podman-cross-build` target does not include `podman manifest push`. We need to use `make podman-cross-build-push` instead (which depends on `podman-cross-build` and adds the push step).

## Script structure

```bash
#!/bin/bash
set -euo pipefail

PROJECTS_DIR="$(cd "$(dirname "$0")/../../projects" && pwd)"

build_and_push() {
    local repo="$1" branch="$2" version="$3"

    echo "Building $repo @ $branch with VERSION=$version"
    cd "$PROJECTS_DIR/$repo"
    git checkout "$branch"

    if [[ "$repo" == "monitoring-plugin" ]]; then
        make podman-cross-build-push VERSION="$version" PLUGIN_NAME=monitoring-console-plugin
    else
        make podman-cross-build VERSION="$version"
    fi
}

# release-coo-ocp-4.12
build_and_push console-dashboards-plugin release-coo-ocp-4.12 v0.4.3
build_and_push distributed-tracing-console-plugin release-coo-ocp-4.12 v0.3.3
build_and_push logging-view-plugin release-coo-ocp-4.12 v6.0.5

# release-coo-ocp-4.15
build_and_push distributed-tracing-console-plugin release-coo-ocp-4.15 v0.4.3
build_and_push logging-view-plugin release-coo-ocp-4.15 v6.1.6
build_and_push monitoring-plugin release-coo-ocp-4.15 v0.4.5

# release-coo-ocp-4.19
build_and_push distributed-tracing-console-plugin release-coo-ocp-4.19 v1.0.3
build_and_push troubleshooting-panel-console-plugin release-coo-ocp-4.19 v0.4.5
build_and_push monitoring-plugin release-coo-ocp-4.19 v0.5.4

# release-coo-ocp-4.22
build_and_push distributed-tracing-console-plugin release-coo-ocp-4.22 v1.1.0
build_and_push troubleshooting-panel-console-plugin release-coo-ocp-4.22 v1.0.0
build_and_push monitoring-plugin release-coo-ocp-4.22 v1.0.0
build_and_push logging-view-plugin release-coo-ocp-4.22 v6.2.1

echo "All images built and pushed successfully."
```

## File to create

- `/Users/emurasak/workspace/harness/tasks/upstream-images/build-and-push.sh` — the shell script (make executable)

## Verification

1. Review the script for correctness against branches_tags.txt
2. Ensure podman is available and the user is logged into quay.io (`podman login quay.io`)
3. Run `./build-and-push.sh` and verify each image is pushed by checking quay.io/openshift-observability-ui/
