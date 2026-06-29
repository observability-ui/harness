# Execution: Add CSV Export to Logs Table

> Results are annotated inline: `-- **value**` for discovered values, `-- **passes/FAILED**` for verification.

## Phase 1: Create CSV Export Action Component
Depends on: nothing | Parallel with: none | Type: implementation | Projects: perses-plugins

### 1a. Extract testable CSV generation logic
- [x] Create pure function `collectLabelKeys(entries: LogEntry[]): string[]` to gather and sort unique label keys
- [x] Create pure function `buildLogsCsvString(entries: LogEntry[]): string` to generate CSV content
- [x] Write failing tests for `collectLabelKeys` and `buildLogsCsvString` - `logstable/src/LogsTableCsvExportAction.test.ts` -- **11 tests**
- [x] Implement functions to pass tests - `logstable/src/LogsTableCsvExportAction.tsx`

### 1b. Build the React action component
- [x] Create `LogsTableCsvExportAction` React component using extracted functions - `logstable/src/LogsTableCsvExportAction.tsx`
- [x] Uses `FileDelimitedOutline` icon from `mdi-material-ui/FileDelimitedOutline`
- [x] Uses `escapeCsvValue`, `sanitizeFilename`, `formatTimestampISO` from `@perses-dev/plugin-system`
- [x] Uses `stripAnsi` from `./utils/ansi`
- [x] Tooltip: "Export as CSV", aria-label: "Export Logs Table Data as CSV"
- [x] Downloads as `{sanitizeFilename(title)}_data.csv` with MIME type `text/csv;charset=utf-8`

### Phase 1 Verification
- [x] `npx tsx node_modules/.bin/jest --config logstable/jest.config.ts` — **55 tests pass, 11 new** (pre-existing: `npm test` fails due to Node 25/Jest 30 config resolution, `LogsTablePanel.test.tsx` fails due to echarts init in jsdom)
- [x] `cd logstable && npm run type-check` — **pre-existing errors only** (`Cannot find module '@perses-dev/spec'` across all project files, not specific to new code)
- [x] `cd logstable && npm run lint` — **passes, no errors**

## Phase 2: Register CSV Export Action
Depends on: Phase 1 | Parallel with: none | Type: configuration | Projects: perses-plugins

- [x] Import `LogsTableCsvExportAction` in `logstable/src/LogsTable.ts`
- [x] Add `{ component: LogsTableCsvExportAction, location: 'header' }` to `actions` array

### Phase 2 Verification
- [x] `cd logstable && npm run lint` — **passes, no errors**
- [x] All 11 CSV tests still pass after registration

---

## Summary

**Status:** Complete (2 of 2 phases done)

### Files changed

| File | Change |
| ---- | ------ |
| `logstable/src/LogsTableCsvExportAction.tsx` | New file. CSV export component with `collectLabelKeys()` and `buildLogsCsvString()` pure functions + `LogsTableCsvExportAction` React component |
| `logstable/src/LogsTableCsvExportAction.test.ts` | New file. 11 unit tests covering label key collection and CSV string generation |
| `logstable/src/LogsTable.ts` | Added import and registered CSV export action alongside existing JSON export |

### Outstanding items

- [ ] Commit changes on `feat/logs-table-csv-export` branch
- [ ] Manual verification: load a dashboard with logs table data, verify CSV download works correctly
- [ ] Push branch and create PR to `perses/plugins`

### Notes

- Pre-existing issue: `npm test` does not work in logstable (Node 25 / Jest 30 config resolution for `jest.shared.ts`). Workaround: `npx tsx node_modules/.bin/jest --config logstable/jest.config.ts`
- Pre-existing issue: `npm run type-check` shows errors for `@perses-dev/spec` across all files in the project — not related to new code
- The codebase uses `@perses-dev/spec` for `LogEntry` (not `@perses-dev/core` as initially assumed from the plan)
