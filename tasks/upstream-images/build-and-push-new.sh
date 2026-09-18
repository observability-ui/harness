#!/bin/bash
set -uo pipefail

REGISTRY_AUTH_FILE=~/.docker/config.json
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
LOG_FILE="$SCRIPT_DIR/build-and-push-new-$(date +%Y%m%d-%H%M%S).log"
exec > >(tee -a "$LOG_FILE") 2>&1
echo "Logging to $LOG_FILE"

PROJECTS_DIR="$(cd "$SCRIPT_DIR/../../projects" && pwd)"
ORG="${ORG:-openshift-observability-ui}"
MAX_RETRIES="${MAX_RETRIES:-3}"
STEP_RESULTS=()
HAS_FAILURE=0

REPOS=(
    distributed-tracing-plugin
    logging-view-plugin
    monitoring-plugin
    troubleshooting-panel-plugin
)

# Successful builds from previous runs (will be skipped)
# Format: "repo|branch|version"
SKIP_BUILDS=(
    # Empty for now - add entries here as builds succeed
)

should_skip_build() {
    local key="$1"
    for skip_key in "${SKIP_BUILDS[@]:-}"; do
        [[ "$skip_key" == "$key" ]] && return 0
    done
    return 1
}

echo "Fetching all branches for each repo..."
for repo in "${REPOS[@]}"; do
    echo "  Fetching $repo..."
    (cd "$PROJECTS_DIR/$repo" && git fetch --all --prune)
done
echo "Fetch complete."

build_and_push() {
    local repo="$1" branch="$2" version="$3"
    local checkout_status="OK" pull_status="OK" branch_status="OK" build_status="" build_attempts=""
    local skip_key="$repo|$branch|$version"

    # Check if this build should be skipped
    if should_skip_build "$skip_key"; then
        echo "=========================================="
        echo "SKIPPING $repo @ $branch with VERSION=$version (already built successfully)"
        echo "=========================================="
        STEP_RESULTS+=("$repo|$branch|$version|SKIPPED|SKIPPED|SKIPPED|SKIPPED|0")
        return
    fi

    echo "=========================================="
    echo "Building $repo @ $branch with VERSION=$version ORG=$ORG"
    echo "=========================================="
    cd "$PROJECTS_DIR/$repo"

    # Step 1: Checkout
    if ! git checkout -f "$branch" 2>&1; then
        checkout_status="FAILED"
        echo "ERROR: checkout failed for $repo @ $branch"
        STEP_RESULTS+=("$repo|$branch|$version|$checkout_status|$pull_status|$branch_status|SKIPPED|0")
        HAS_FAILURE=1
        return
    fi

    # Step 2: Pull latest
    if ! git pull --rebase origin "$branch" 2>&1; then
        pull_status="FAILED"
        echo "WARNING: pull failed for $repo @ $branch, building with local state"
    fi

    # Step 3: Verify branch
    local actual_branch
    actual_branch="$(git branch --show-current)"
    if [[ "$actual_branch" != "$branch" ]]; then
        branch_status="WRONG: $actual_branch"
        echo "ERROR: expected branch $branch but on $actual_branch"
        STEP_RESULTS+=("$repo|$branch|$version|$checkout_status|$pull_status|$branch_status|SKIPPED|0")
        HAS_FAILURE=1
        return
    fi

    # Step 3b: Patch Makefile to use Dockerfile instead of Dockerfile.dev
    if [[ "$repo" != "monitoring-plugin" ]]; then
        sed -i.bak 's/-f Dockerfile\.dev/-f Dockerfile/' Makefile && rm -f Makefile.bak
    fi

    # Step 3c: Replace unavailable CI registry base images in Dockerfile
    if [[ -f Dockerfile ]]; then
        sed -i.bak 's|registry.ci.openshift.org/ocp/[0-9.]*:base-rhel9|registry.access.redhat.com/ubi9/ubi-minimal|' Dockerfile && rm -f Dockerfile.bak
        sed -i.bak 's|registry.ci.openshift.org/ocp/builder:rhel-9-golang-[0-9.]*-openshift-[0-9.]*|registry.redhat.io/openshift/golang-builder:golang-builder-v1.26-rhel9|' Dockerfile && rm -f Dockerfile.bak
        sed -i.bak 's|registry.ci.openshift.org/ocp/builder:rhel-8-golang-[0-9.]*-openshift-[0-9.]*|registry.redhat.io/ubi9/go-toolset:1.26|' Dockerfile && rm -f Dockerfile.bak
        # Fix nodejs-18 → nodejs-22 for compatibility with copy-webpack-plugin v14 (requires toSorted() from Node 20+)
        sed -i.bak 's|registry.redhat.io/ubi9/nodejs-18|registry.redhat.io/ubi9/nodejs-22|' Dockerfile && rm -f Dockerfile.bak
    fi

    # Step 4: Build
    local rc=0 attempt
    for attempt in $(seq 1 "$MAX_RETRIES"); do
        rc=0
        echo "--- Attempt $attempt/$MAX_RETRIES ---"
        if [[ "$repo" == "monitoring-plugin" ]]; then
            make podman-cross-build-push VERSION="$version" PLUGIN_NAME=monitoring-console-plugin ORG="$ORG" || rc=$?
        else
            make podman-cross-build VERSION="$version" ORG="$ORG" REGISTRY_ORG="$ORG" TAG="$version" || rc=$?
        fi
        [[ $rc -eq 0 ]] && break
        echo "WARNING: Attempt $attempt failed for $repo @ $branch (exit code $rc)"
    done

    if [[ $rc -eq 0 ]]; then
        build_status="SUCCESS"
    else
        build_status="FAILED"
        HAS_FAILURE=1
        echo "WARNING: Build failed for $repo @ $branch after $MAX_RETRIES attempts, continuing..."
    fi

    STEP_RESULTS+=("$repo|$branch|$version|$checkout_status|$pull_status|$branch_status|$build_status|$attempt")
}

# release-coo-ocp-4.12 (updated branch name, no dashboards-plugin)
build_and_push distributed-tracing-plugin release-coo-ocp-4.12 v0.3.4
build_and_push logging-view-plugin release-coo-ocp-4.12 v6.0.6

# release-coo-ocp-4.15 (updated branch name and versions)
build_and_push distributed-tracing-plugin release-coo-ocp-4.15 v0.4.4
build_and_push logging-view-plugin release-coo-ocp-4.15 v6.1.7
build_and_push monitoring-plugin release-coo-ocp-4.15 v0.4.6

# release-coo-ocp-4.19 (updated versions)
build_and_push distributed-tracing-plugin release-coo-ocp-4.19 v1.0.5
build_and_push troubleshooting-panel-plugin release-coo-ocp-4.19 v0.4.6
build_and_push monitoring-plugin release-coo-ocp-4.19 v0.5.5

# release-coo-ocp-4.22 (updated versions)
build_and_push distributed-tracing-plugin release-coo-ocp-4.22 v1.2.0
build_and_push troubleshooting-panel-plugin release-coo-ocp-4.22 v1.1.0
build_and_push monitoring-plugin release-coo-ocp-4.22 v1.1.0
build_and_push logging-view-plugin release-coo-ocp-4.22 v6.3.0

echo ""
echo "=========================================="
echo "  BUILD SUMMARY"
echo "=========================================="
printf "%-45s %-30s %-10s %-10s %-10s %-10s %-10s %s\n" \
    "REPO" "BRANCH" "VERSION" "CHECKOUT" "PULL" "BRANCH" "BUILD" "ATTEMPTS"
printf "%-45s %-30s %-10s %-10s %-10s %-10s %-10s %s\n" \
    "----" "------" "-------" "--------" "----" "------" "-----" "--------"
for result in "${STEP_RESULTS[@]}"; do
    IFS='|' read -r repo branch version checkout pull branchv build attempts <<< "$result"
    printf "%-45s %-30s %-10s %-10s %-10s %-10s %-10s %s\n" \
        "$repo" "$branch" "$version" "$checkout" "$pull" "$branchv" "$build" "$attempts"
done
echo "=========================================="
FAIL_COUNT=0
SKIP_COUNT=0
for result in "${STEP_RESULTS[@]}"; do
    IFS='|' read -r _ _ _ checkout pull branchv build _ <<< "$result"
    if [[ "$build" == "SKIPPED" ]]; then
        ((SKIP_COUNT++))
    elif [[ "$checkout" != "OK" && "$checkout" != "SKIPPED" ]] || [[ "$branchv" != "OK" && "$branchv" != "SKIPPED" ]] || [[ "$build" == "FAILED" ]]; then
        ((FAIL_COUNT++))
    fi
done
SUCCESS_COUNT=$(( ${#STEP_RESULTS[@]} - FAIL_COUNT - SKIP_COUNT ))
echo "Total: ${#STEP_RESULTS[@]}  |  Succeeded: $SUCCESS_COUNT  |  Skipped: $SKIP_COUNT  |  Failed: $FAIL_COUNT"
echo "=========================================="

# Image Summary
echo ""
echo "=========================================="
echo "  IMAGE SUMMARY"
echo "=========================================="

# Get image name for a repo (maps repo directory name to image name)
get_image_name() {
    local repo="$1"
    case "$repo" in
        distributed-tracing-plugin)
            echo "distributed-tracing-console-plugin"
            ;;
        troubleshooting-panel-plugin)
            echo "troubleshooting-panel-console-plugin"
            ;;
        monitoring-plugin)
            echo "monitoring-console-plugin"
            ;;
        *)
            echo "$repo"
            ;;
    esac
}

# Successfully built images
echo ""
echo "✅ SUCCESSFULLY BUILT IMAGES:"
SUCCESS_FOUND=0
for result in "${STEP_RESULTS[@]}"; do
    IFS='|' read -r repo branch version checkout pull branchv build attempts <<< "$result"
    if [[ "$build" == "SUCCESS" ]]; then
        image_name=$(get_image_name "$repo")
        echo "   quay.io/$ORG/$image_name:$version"
        SUCCESS_FOUND=1
    fi
done
[[ $SUCCESS_FOUND -eq 0 ]] && echo "   (none)"

# Skipped images
echo ""
echo "⏭️  SKIPPED IMAGES (already built):"
SKIP_FOUND=0
for result in "${STEP_RESULTS[@]}"; do
    IFS='|' read -r repo branch version checkout pull branchv build attempts <<< "$result"
    if [[ "$build" == "SKIPPED" ]]; then
        image_name=$(get_image_name "$repo")
        echo "   quay.io/$ORG/$image_name:$version"
        SKIP_FOUND=1
    fi
done
[[ $SKIP_FOUND -eq 0 ]] && echo "   (none)"

# Failed images
echo ""
echo "❌ FAILED IMAGES:"
FAIL_FOUND=0
for result in "${STEP_RESULTS[@]}"; do
    IFS='|' read -r repo branch version checkout pull branchv build attempts <<< "$result"
    if [[ "$build" == "FAILED" ]] || [[ "$checkout" != "OK" && "$checkout" != "SKIPPED" ]] || [[ "$branchv" != "OK" && "$branchv" != "SKIPPED" ]]; then
        image_name=$(get_image_name "$repo")
        echo "   quay.io/$ORG/$image_name:$version ($build)"
        FAIL_FOUND=1
    fi
done
[[ $FAIL_FOUND -eq 0 ]] && echo "   (none)"

echo ""
echo "=========================================="

exit $HAS_FAILURE
