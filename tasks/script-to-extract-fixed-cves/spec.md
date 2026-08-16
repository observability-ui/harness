# Spec: Script to extract fixed CVEs

## Related projects and branches

- perses-operator: release-coo-1.5
- monitoring-plugin: release-coo-ocp-4.15
- monitoring-plugin: release-coo-ocp-4.19
- monitoring-plugin: release-coo-ocp-4.22
- logging-view-plugin: release-coo-ocp-4.12
- logging-view-plugin: release-coo-ocp-4.15
- logging-view-plugin: release-coo-ocp-4.22
- distributed-tracing-console-plugin: release-coo-ocp-4.12
- distributed-tracing-console-plugin: release-coo-ocp-4.15
- distributed-tracing-console-plugin: release-coo-ocp-4.19
- distributed-tracing-console-plugin: release-coo-ocp-4.22
- troubleshooting-panel-console-plugin: release-coo-ocp-4.19
- troubleshooting-panel-console-plugin: release-coo-ocp-4.22

## Description

I need to create a script that extracts the fixed CVEs from the latest commits. The script should rollback the last 3 commits and execute the
`npm audit` command to get the list of fixed CVEs. In some projects the frontend could be installed in a different directory, so the script should be
able to handle that. It should also detect the vulnerabilities that were fixed in go using govulncheck. The script should be able to handle different
package managers and should be able to extract the fixed CVEs from the output of both analysis tools

## Acceptance criteria

- The script should be able to rollback the last 3 commits and execute the `npm audit` command to get the list of fixed CVEs.
- The script should be able to rollback the last 3 commits and execute the `govulncheck` command to get the list of fixed CVEs.
- The list should be able to receive multiple projects in a loop.
