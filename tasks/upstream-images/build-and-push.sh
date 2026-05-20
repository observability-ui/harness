#!/bin/bash
set -uo pipefail

PROJECTS_DIR="$(cd "$(dirname "$0")/../../projects" && pwd)"
BUILD_RESULTS=()
HAS_FAILURE=0

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
    local rc=0

    echo "=========================================="
    echo "Building $repo @ $branch with VERSION=$version"
    echo "=========================================="
    cd "$PROJECTS_DIR/$repo"
    git checkout "$branch"

    if [[ "$repo" == "monitoring-plugin" ]]; then
        make podman-cross-build-push VERSION="$version" PLUGIN_NAME=monitoring-console-plugin || rc=$?
    else
        make podman-cross-build VERSION="$version" || rc=$?
    fi

    if [[ $rc -eq 0 ]]; then
        BUILD_RESULTS+=("SUCCESS|$repo|$branch|$version")
    else
        BUILD_RESULTS+=("FAILED|$repo|$branch|$version")
        HAS_FAILURE=1
        echo "WARNING: Build failed for $repo @ $branch (exit code $rc), continuing..."
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

echo ""
echo "=========================================="
echo "  BUILD SUMMARY"
echo "=========================================="
printf "%-10s %-45s %-30s %s\n" "STATUS" "REPO" "BRANCH" "VERSION"
printf "%-10s %-45s %-30s %s\n" "------" "----" "------" "-------"
for result in "${BUILD_RESULTS[@]}"; do
    IFS='|' read -r status repo branch version <<< "$result"
    printf "%-10s %-45s %-30s %s\n" "$status" "$repo" "$branch" "$version"
done
echo "=========================================="
FAIL_COUNT=0
for result in "${BUILD_RESULTS[@]}"; do
    [[ "$result" == FAILED* ]] && ((FAIL_COUNT++))
done
SUCCESS_COUNT=$(( ${#BUILD_RESULTS[@]} - FAIL_COUNT ))
echo "Total: ${#BUILD_RESULTS[@]}  |  Succeeded: $SUCCESS_COUNT  |  Failed: $FAIL_COUNT"
echo "=========================================="

exit $HAS_FAILURE
