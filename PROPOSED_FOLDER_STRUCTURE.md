## Monorepo

```
observability-console-plugin/
|
+-- cmd/
|   +-- plugin-backend.go
|
+-- pkg/
|   +-- server/
|   |   +-- server.go           // config, lifecycle, TLS
|   |   +-- server_test.go
|   |   +-- routes.go           // Route registration
|   |   +-- handlers.go         // health, non-patch features, config, files, CORS
|   |   +-- manifest.go         // JSON Patch manifest handler
|   +-- monitoring/
|   |   +-- acm_proxy.go        // ACM Alertmanager/Thanos reverse proxy
|   |   +-- acm_proxy_test.go
|   +-- tracing/
|   |   +-- api.go              // Tempo resource discovery
|   |   +-- api_test.go
|   |   +-- proxy.go            // Tempo reverse proxy
|   |   +-- proxy_test.go
|   +-- alertmanagement/        // move all management files in here
|   |   +-- management.go
|   |   +-- types.go
|   |   +-- ...
|   +-- k8s/                    // k8s utils
|   |   +-- client.go
|   |   +-- types.go
|   |   +-- ...
|
+-- api/
|   +-- alertmanagement-openapi.yaml
|   +-- oapi-codegen.yaml
|   +-- korrel8r-openapi.yaml
|
+-- internal/
|   +-- managementrouter/       // Generated API handlers
|       +-- ...
|
+-- web/
|   +-- src/
|   |   +-- features/
|   |   |   +-- alerting/
|   |   |   |   +-- pages/
|   |   |   |   |   ...
|   |   |   |   +-- components/
|   |   |   |   |   ...
|   |   |   |   +-- hooks/
|   |   |   |   |   +-- useAlerts.ts
|   |   |   |   |   +-- useSortAlerts.ts
|   |   |   |   |   +-- useSelectedFilters.ts
|   |   |   |   +-- types.ts
|   |   |   |
|   |   |   +-- metrics/
|   |   |   |   +-- pages/
|   |   |   |   |   +-- MetricsPage.tsx
|   |   |   |   |   +-- TargetsPage.tsx
|   |   |   |   |   +-- TargetDetail.tsx
|   |   |   |   +-- components/
|   |   |   |   |   +-- PromQLExpressionInput.tsx
|   |   |   |   +-- hooks/
|   |   |   |   +-- types.ts
|   |   |   |
|   |   |   +-- legacy-dashboards/
|   |   |   |   +-- pages/
|   |   |   |   |   +-- LegacyDashboardPage.tsx
|   |   |   |   +-- components/
|   |   |   |   |   ...
|   |   |   |   +-- hooks/
|   |   |   |   |   +-- useLegacyDashboards.ts
|   |   |   |   +-- types.ts
|   |   |   |
|   |   |   +-- perses-dashboards/
|   |   |   |   +-- pages/
|   |   |   |   |   +-- DashboardListPage.tsx
|   |   |   |   |   +-- DashboardPage.tsx
|   |   |   |   +-- components/
|   |   |   |   |   +-- DashboardActions.tsx
|   |   |   |   |   +-- ProjectBar.tsx
|   |   |   |   +-- hooks/
|   |   |   |   |   +-- useDashboardsData.ts
|   |   |   |   |   +-- usePerses.ts
|   |   |   |   +-- types.ts
|   |   |   |
|   |   |   +-- incidents/
|   |   |   |   +-- pages/
|   |   |   |   |   +-- IncidentsPage.tsx
|   |   |   |   +-- components/
|   |   |   |   |   ...
|   |   |   |   +-- hooks/
|   |   |   |   |   +-- useIncidents.ts
|   |   |   |   +-- utils/
|   |   |   |   |   +-- processAlerts.ts
|   |   |   |   |   +-- processIncidents.ts
|   |   |   |   +-- types.ts
|   |   |   |
|   |   |   +-- logs/                               // to be replaced with the logs explore view. Will keep code in the logging-view-plugin for ELS versions
|   |   |   |
|   |   |   +-- traces/
|   |   |   |   +-- pages/
|   |   |   |   |   ...
|   |   |   |   +-- components/
|   |   |   |   |   ...
|   |   |   |   +-- hooks/
|   |   |   |   |   ...
|   |   |   |   +-- utils/
|   |   |   |   |   +-- transformTrace.ts
|   |   |   |   |   +-- links.ts                    // should we co-locate all of the link utils? There is this, monitorings usePerspective stuff and also the troubleshooting panel
|   |   |   |   |   +-- filter.ts
|   |   |   |   +-- types.ts
|   |   |   |
|   |   |   +-- troubleshooting/
|   |   |       +-- components/
|   |   |       |   ...
|   |   |       +-- hooks/
|   |   |       |   ...
|   |   |       +-- korrel8r/
|   |   |       |   +-- gen-client/
|   |   |       |   ... domain converters
|   |   |       +-- types.ts
|   |   |
|   |   +-- shared/
|   |   |
|   |   |   +-- components/
|   |   |   |   +-- ErrorAlert.tsx
|   |   |   |   +-- LoadingState.tsx
|   |   |   |   +-- EmptyState.tsx
|   |   |   |   +-- TypeaheadSelect.tsx
|   |   |   |   +-- TimeRangeDropdown.tsx
|   |   |   |   +-- TimeRangeSelectModal.tsx
|   |   |   |   +-- DateTimePicker.tsx
|   |   |   |   +-- RefreshIntervalDropdown.tsx
|   |   |   |   +-- Labels.tsx
|   |   |   |   +-- KebabDropdown.tsx
|   |   |   |   +-- TablePagination.tsx
|   |   |   |   +-- ToggleButton.tsx
|   |   |   |   +-- PersesWrapper.tsx
|   |   |   |
|   |   |   +-- hooks/
|   |   |   |   +-- useBoolean.ts
|   |   |   |   +-- usePatternFlyTheme.ts
|   |   |   |   +-- usePerspective.ts
|   |   |   |   +-- useFeatures.ts
|   |   |   |   +-- useDeepMemo.ts
|   |   |   |   +-- useDebounce.ts
|   |   |   |   +-- useIsVisible.ts
|   |   |   |   +-- useRefWidth.ts
|   |   |   |   +-- useMonitoring.ts
|   |   |   |
|   |   |   +-- utils/
|   |   |   |   +-- cancellableFetch.ts     // Should just use react-query where possible
|   |   |   |   +-- dateTime.ts
|   |   |   |   +-- units.ts
|   |   |   |   +-- sort.ts
|   |   |   |   +-- severity.ts
|   |   |   |   +-- queryParams.ts
|   |   |   |   +-- dataTest.ts
|   |   |   |
|   |   |   +-- store/
|   |   |   |   +-- store.ts
|   |   |   |   +-- actions.ts
|   |   |   |   +-- reducers.ts
|   |   |   |   +-- monitoring/            // Monitoring slice
|   |   |   |   |   +-- actions.ts
|   |   |   |   |   +-- reducers.ts
|   |   |   |   |   +-- thunks.ts
|   |   |   |   +-- troubleshooting/       // Troubleshooting panel slice
|   |   |   |       +-- actions.ts
|   |   |   |       +-- reducers.ts
|   |   |   |
|   |   |   +-- types/
|   |   |   |   +-- monitoring.ts
|   |   |   |   +-- k8s.ts
|   |   |   |
|   |   |   +-- console/                   // Vendored console utilities
|   |   |       +-- models/
|   |   |       +-- utils/
|   |   |       +-- graphs/
|   |   |
|   |   +-- i18n.ts
|   |   +-- index.d.ts
|   |
|   +-- cypress/
|   |   +-- e2e/
|   |   |   +-- alerting/
|   |   |   +-- metrics/
|   |   |   +-- dashboards/
|   |   |   +-- incidents/
|   |   |   +-- logs/
|   |   |   +-- traces/
|   |   |   +-- troubleshooting/
|   |   +-- fixtures/
|   |   +-- support/
|   |   |   +-- commands/
|   |   +-- views/
|   |
|   +-- locales/
|   |   +-- en/
|   |       +-- plugin__observability-console-plugin.json
|   |
|   +-- console-extensions.json             // should be empty, all features are optional
|   +-- package.json
|   +-- tsconfig.json
|   +-- webpack.config.ts
|   +-- jest.config.js
|   +-- cypress.config.ts
|   +-- eslint.config.ts
|   +-- .prettierrc.yml
|   +-- .swcrc
+-- config/
|   +-- alerting.patch.json                     (admin + virt + dev)
|   +-- metrics.patch.json                      (admin + virt + dev)
|   +-- legacy-dashboards.patch.json            (admin + virt + dev)
|   +-- perses-dashboards.patch.json            (admin + virt + ACM) + OLS tool-ui
|   +-- acm-alerting.patch.json
|   +-- cluster-health-analyzer.patch.json
|   |
|   +-- logs.patch.json
|   +-- logs-dev-console.patch.json
|   +-- logs-alerts.patch.json
|   +-- logs-alerts-charts.patch.json
|   |
|   +-- traces.patch.json
|   |
|   +-- troubleshooting.patch.json
|   +-- troubleshooting-agent-navigation.patch.json
|
+-- scripts/
|   +-- start-console.sh        // This is gonna be realy unfun to rewrite with all the feature options
|   ...
|
+-- hack/
|   +-- docker-compose/
|   +-- deploy-tempostack.sh
|   +-- deploy-otel.sh
|
+-- Makefile
+-- Dockerfile
+-- Dockerfile.dev
+-- Dockerfile.test
+-- Dockerfile.devspace
+-- .ci-operator.yaml
+-- .tool-versions
+-- ct.yaml
+-- go.mod
+-- go.sum
+-- OWNERS
+-- LICENSE
+-- README.md
```

### Notes:

- All tests targeted at a single file / component should be co-located with the file it is testing (ie. `server_test.go` or `AlertsPage.spec.tsx`)
-

### Implementation Steps:

1. Start adjusting structure in monitoring-plugin, don’t move plugins
2. Reach out with other teams to align with the new structure
3. Update CMO with new "send all feature flags" structure to have it send the alerting, metrics, legacy-dashboards, targets features
4. Create a new version of the UIPlugin CR (v1beta1 or v1), as we would need only one plugin moving forward. Create a version upgrade hook to merge
   all old CR's into the single one. This can be completed before all items are in a single repo, we can just perform the adjustments in the operator
5. Define owners / ownership codification
6. Define test running optimizations (ie. only run X test when Y folder changes) to prevent lots of tests from running all the time
7.
