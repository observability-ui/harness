# Plan: Add CSV Export to Logs Table

## Problem

The logs table panel currently supports exporting logs as JSON only. Users need the ability to export logs as CSV, a more universally compatible format for use in spreadsheets, data analysis tools, and other systems. The CSV export should include timestamps in ISO 8601 format, the raw log message, and extracted labels as additional columns.

## Current State

| Component | File / Location | Current Behavior |
| --------- | --------------- | ---------------- |
| JSON export action | `projects/perses-plugins/logstable/src/LogsTableExportAction.tsx` | Exports all log entries as a JSON file. Extracts entries from `queryResults`, stringifies with 2-space indent, downloads as `{title}_data.json` |
| Plugin definition | `projects/perses-plugins/logstable/src/LogsTable.ts:33` | Registers one action: `LogsTableExportAction` at `location: 'header'` |
| Data model | `@perses-dev/core` → `LogEntry` | `{ timestamp: number, line: string, labels: Labels }` where `timestamp` is seconds since epoch and `Labels = Record<string, string>` |
| CSV utilities | `@perses-dev/plugin-system` (from `perses-shared/plugin-system/src/utils/csv-export.ts`) | Provides `escapeCsvValue()`, `formatTimestampISO()`, and `sanitizeFilename()` — all already used by other export actions (table, bar chart, time series) |
| ANSI stripping | `projects/perses-plugins/logstable/src/utils/ansi.ts` | Provides `stripAnsi()` to remove ANSI escape codes from log lines |
| Table CSV export | `projects/perses-plugins/table/src/TableExportAction.tsx` | Reference implementation: builds header row + data rows using `escapeCsvValue`, creates `text/csv` blob, downloads |

## Changes

### Phase 1: Create CSV Export Action Component

**Dependency:** None
**Parallel with:** None

#### Files Modified

| File | Change |
| ---- | ------ |
| `projects/perses-plugins/logstable/src/LogsTableCsvExportAction.tsx` | **New file.** CSV export action component |

#### Details

Create `LogsTableCsvExportAction.tsx` following the same pattern as `LogsTableExportAction.tsx` (JSON) and `table/src/TableExportAction.tsx` (CSV reference).

**Imports:**
- `escapeCsvValue`, `sanitizeFilename`, `formatTimestampISO` from `@perses-dev/plugin-system`
- `InfoTooltip` from `@perses-dev/components`
- `IconButton` from `@mui/material`
- `FileDelimitedOutline` from `mdi-material-ui/FileDelimitedOutline` (CSV-specific icon, visually distinct from the JSON export's `Download` icon)
- `LogEntry` from `@perses-dev/core`
- `stripAnsi` from `./utils/ansi`
- `LogsTableProps` from `./model`

**Component logic:**

1. **Extract entries** — same as JSON export: `queryResults.flatMap((q) => q.data?.logs?.entries ?? [])`

2. **Collect all label keys** — iterate all entries, gather the union of all label keys, sort alphabetically. This handles entries with different label sets gracefully.

3. **Build CSV string:**
   - **Header row:** `timestamp,body,{label1},{label2},...` — fixed columns `timestamp` and `body` first, then sorted label columns
   - **Data rows:** For each entry:
     - `formatTimestampISO(entry.timestamp)` for the timestamp column (ISO 8601)
     - `stripAnsi(entry.line)` for the body column (raw message without ANSI codes)
     - `entry.labels[key] ?? ''` for each label column
     - All values escaped with `escapeCsvValue()`
   - Join rows with `\n`, add trailing newline

4. **Download** — create `Blob` with `text/csv;charset=utf-8` MIME type, use same download pattern as JSON export. Filename: `{sanitizeFilename(title)}_data.csv`

5. **UI** — same `InfoTooltip` + `IconButton` pattern. Tooltip: `"Export as CSV"`. Aria label: `"Export Logs Table Data as CSV"`. Disabled when no data.

**Code snippet for CSV generation core:**

```typescript
const allLabelKeys = useMemo(() => {
  const keys = new Set<string>();
  for (const entry of entries) {
    if (entry.labels) {
      for (const key of Object.keys(entry.labels)) {
        keys.add(key);
      }
    }
  }
  return Array.from(keys).sort();
}, [entries]);

const handleDownload = useCallback((): void => {
  if (isDisabled) return;
  try {
    const headerRow = ['timestamp', 'body', ...allLabelKeys].map(escapeCsvValue).join(',');
    const dataRows = entries.map((entry) => {
      const timestamp = escapeCsvValue(formatTimestampISO(entry.timestamp));
      const body = escapeCsvValue(stripAnsi(entry.line));
      const labels = allLabelKeys.map((key) => escapeCsvValue(entry.labels?.[key] ?? ''));
      return [timestamp, body, ...labels].join(',');
    });
    const csvString = [headerRow, ...dataRows].join('\n') + '\n';
    // ... blob creation and download (same pattern as JSON export)
  } catch (error) {
    console.error('Logs table CSV export failed:', error);
  }
}, [entries, allLabelKeys, isDisabled, definition]);
```

#### Phase 1 Verification

- File compiles without TypeScript errors: `cd ./projects/perses-plugins && npx tsc --noEmit --project logstable/tsconfig.json` (or equivalent)
- Manual review: CSV output for test data has correct header, ISO 8601 timestamps, stripped ANSI codes, properly escaped values

### Phase 2: Register CSV Export Action

**Dependency:** Phase 1
**Parallel with:** None

#### Files Modified

| File | Change |
| ---- | ------ |
| `projects/perses-plugins/logstable/src/LogsTable.ts` | Add `LogsTableCsvExportAction` import and register as second action |

#### Details

Add the CSV export action to the `actions` array in `LogsTable.ts:33`:

```typescript
import { LogsTableCsvExportAction } from './LogsTableCsvExportAction';

// ...
actions: [
  { component: LogsTableExportAction, location: 'header' },
  { component: LogsTableCsvExportAction, location: 'header' },
],
```

Both actions render as separate icon buttons in the panel header. The CSV action uses the `FileDelimitedOutline` icon (a document with delimiter lines) while the JSON action keeps its existing `Download` icon, making them visually distinct at a glance. Tooltips further clarify each button's function.

#### Phase 2 Verification

- File compiles without TypeScript errors
- Both actions are registered in the plugin definition
- Lint passes: `cd ./projects/perses-plugins && npm run lint -- --filter logstable` (or equivalent)

## PR Strategy

| PR | Repository | Branch | Description | Dependencies |
| -- | ---------- | ------ | ----------- | ------------ |
| 1  | perses/plugins | `feat/logs-table-csv-export` from `main` | Add CSV export action to logs table panel | None |

Single PR since all changes are in one repo (perses-plugins) and the two phases are tightly coupled.

## Verification

- **The logs table supports exporting logs as CSV** — Load a dashboard with a logs table panel that has data. Verify a second download button appears in the panel header with tooltip "Export as CSV". Click it. Verify the downloaded `.csv` file:
  - Has a header row with `timestamp,body,{sorted label columns}`
  - Timestamps are in ISO 8601 format (e.g., `2026-01-21T15:32:31.000Z`)
  - Log body is the raw message text with ANSI codes stripped
  - Labels are correctly extracted into individual columns
  - Values containing commas, quotes, or newlines are properly escaped
  - File is named `{panelName}_data.csv`
- **JSON export still works** — Verify the existing JSON export button still functions correctly and is unaffected
- **Empty state** — When no log data is available, both export buttons should be disabled with appropriate tooltips
- **Build** — `npm run build` in the logstable plugin directory succeeds
- **Lint** — `npm run lint` passes

## Risks

| Risk | Impact | Mitigation |
| ---- | ------ | ---------- |
| Entries with highly heterogeneous label sets produce many sparse columns | CSV file has many empty cells, making it harder to read | This is inherent to the data shape; sorted columns and empty-string defaults keep it parseable. Could add a future option to select which labels to export. |
| Log lines containing CSV-special characters (commas, quotes, newlines) | Malformed CSV if not properly escaped | Using `escapeCsvValue()` from `@perses-dev/plugin-system`, which handles all RFC 4180 escaping (wraps in quotes, doubles internal quotes) |
| Log lines containing ANSI escape codes | Raw ANSI codes pollute CSV data | Using `stripAnsi()` to clean log lines before export |
| Header button clutter as more actions are added | Panel header becomes crowded | Each action has a distinct icon (`Download` for JSON, `FileDelimitedOutline` for CSV) and tooltip. This is consistent with how action buttons work across panels. |
| Large log datasets produce large CSV files | Browser may lag or run out of memory | Same risk exists for JSON export; no regression. Could add streaming/pagination in a future iteration. |
