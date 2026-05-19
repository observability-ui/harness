#!/bin/bash
set -euo pipefail

PROJECTS_DIR="$(cd "$(dirname "$0")/../../projects" && pwd)"

REPOS=(
    console-dashboards-plugin
    distributed-tracing-console-plugin
    logging-view-plugin
    monitoring-plugin
    troubleshooting-panel-console-plugin
)

echo "Fetching all branches for each repo..."
for repo in "${REPOS[@]}"; do
    echo "  Fetching $repo..."
    (cd "$PROJECTS_DIR/$repo" && git fetch --all --prune)
done
echo "Fetch complete."

build_and_push() {
    local repo="$1" branch="$2" version="$3"

    echo "=========================================="
    echo "Building $repo @ $branch with VERSION=$version"
    echo "=========================================="
    cd "$PROJECTS_DIR/$repo"
    git checkout "$branch"

    local dockerfile
    dockerfile=$(grep -A2 'podman-cross-build:' Makefile | grep -oE '\-f [^ ]+' | head -1 | sed 's/-f //')
    dockerfile="${dockerfile:-Dockerfile}"

    sed -i.bak 's/npm ci --ignore-scripts/npm ci/g' Makefile && rm -f Makefile.bak
    sed -i.bak '/nodejs/{/--platform/!s/^FROM /FROM --platform=linux\/amd64 /;}' "$dockerfile" && rm -f "${dockerfile}.bak"

    make install

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

echo "=========================================="
echo "All images built and pushed successfully."
echo "=========================================="
