# Spec: Build and push upstream images for UIPlugins

## Related projects and branches

- console-dashboards-plugin: branch release-coo-ocp-4.12
- distributed-tracing-console-plugin: branches release-coo-ocp-4.12, release-coo-ocp-4.15, release-coo-ocp-4.19, release-coo-ocp-4.22
- monitoring-plugin: branches release-coo-ocp-4.15, release-coo-ocp-4.19, release-coo-ocp-4.22
- loggin-view-plugin: branches release-coo-ocp-4.12, release-coo-ocp-4.15, release-coo-ocp-4.22
- troubleshooting-panel-console-plugin: branches release-coo-ocp-4.19, release-coo-ocp-4.22

## Description

In order to generate upstream images for Observability Operator, UIPlugin images should be built and pushed to quay repository.
For each repository, we have different branches. For each branch, an image should be built with a certain tag and pushed to quay.
For reference, there is a file branches_tags.txt that helps.

The command to build the image is `make podman-cross-build` and VERSION should be the tag presented in branches_tags.txt file.

Create a shell script that, for each repository, checks out the branch, build the image and push to quay.

## Acceptance criteria

All images presented in branches_tags.txt should be pushed and updated in quay. As you can noticed, the exception of plugin name is monitoring-console-plugin instead of monitoring-plugin, that deviates from the repository name in the image name.

