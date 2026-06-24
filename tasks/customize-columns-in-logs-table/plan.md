# Plan: Customize Columns in Logs Table

## Problem

The logs table panel currently shows only two fixed columns: timestamp and log line. Log entries contain structured labels (e.g., `app`, `host`, `status`, `method`, `namespace`) that are only visible when expanding a row's details panel. Users need to see these labels as first-class columns to scan, sort, and compare log entries without expanding each row individually.

## Current State

| Component | File / Location | Current Behavior |
| --------- | --------------- | ---------------- |
| Panel options type | `logstable/src/model.ts:30-39` | `LogsTableOptions` has `showTime`, `allowWrap`, `enableDetails`, `showAll`, `selection`, `actions` — no column customization |
| Plugin registration | `logstable/src/LogsTable.ts:21-34` | Two editor tabs: "Settings" and "Item Actions". No column editor tab |
| Settings editor | `logstable/src/LogsTableSettingsEditor.tsx:26-47` | Legend + Thresholds only — no column controls |
| Row rendering | `logstable/src/components/LogRow/LogRow.tsx:56-368` | Fixed grid: `16px` (expand) + `minmax(160px, max-content)` (timestamp) + `1fr` (log line) + optional actions. No dynamic columns |
| Grid layout | `logstable/src/components/LogRow/LogsStyles.tsx:26-48` | `LogRowContent` styled component with hardcoded `gridTemplateColumns` |
| Data flow | `logstable/src/LogsTableComponent.tsx:19-43` | Flattens all query results, sorts by timestamp descending, passes to `LogsList` |
| Virtualized list | `logstable/src/components/VirtualizedLogsList.tsx:39-462` | Renders `LogRow` per entry via Virtuoso. No header row, no sort controls |
| Log entry type | `@perses-dev/spec` LogEntry | `{ timestamp: number; line: string; labels: Labels }` where `Labels = Record<string, string>` |
| CUE schema | `logstable/schemas/logstable.cue:21-27` | `showTime`, `allowWrap`, `enableDetails`, `selection`, `actions` — no column settings |

## Reference Implementation: Alertmanager Column Editor (PR #647)

> Source: [perses/plugins#647](https://github.com/perses/plugins/pull/647) on branch `feat/alert-manager-plugin`

The alertmanager plugin introduces a reusable generic `ColumnsEditor<C>` component and an `AlertTableColumnsEditor` wrapper. This is the pattern we follow for the logs table — **not** the table plugin's `ColumnsEditor/ColumnEditorContainer/ColumnEditor` hierarchy, which is heavier (drag-and-drop, expand/collapse panels, embedded visualization editors).

### Architecture overview

```
ColumnsEditor<C>                  (generic, reusable)
  └── ColumnEntry<C>              (per-column card with form fields)
        ├── Name field            (via renderNameField prop — caller decides the input type)
        ├── Header text field
        ├── Enable sorting checkbox
        ├── Sort mode select
        ├── Default sort select
        └── Extra fields          (via renderExtraFields prop — caller adds domain-specific controls)

AlertTableColumnsEditor           (alert-table-specific wrapper)
  └── ColumnsEditor<ColumnDefinition>
        renderNameField → <TextField label="Label key" />

LogsTableColumnsEditor            (logs-table-specific wrapper)
  └── ColumnsEditor<LogsColumnDefinition>
        renderNameField  → <TextField label="Column name" helperText="..." />
        renderExtraFields → <FormControlLabel "Wrap content" />
```

### Key files from PR #647

**Generic column editor — `alertmanager/src/components/ColumnsEditor.tsx`:**

```typescript
export interface BaseColumnDefinition {
  name: string;
  header?: string;
  enableSorting?: boolean;
  sort?: 'asc' | 'desc';
  sortMode?: string;
}

export type ColumnUpdater<C> = (index: number, updater: (draft: C) => void) => void;

export interface ColumnsEditorProps<C extends BaseColumnDefinition> {
  columns: C[];
  description: string;                                          // help text above the column list
  sortModeLabels: Record<string, string>;                       // map of sortMode values → display labels
  defaultSortMode: string;                                      // initial sortMode for new columns
  getDisplayName: (column: C) => string;                        // title shown in each column card
  getHeaderPlaceholder: (column: C) => string;                  // placeholder for the header text field
  onAdd: () => void;
  onRemove: (index: number) => void;
  onUpdate: ColumnUpdater<C>;                                   // immer-style draft updater
  onMoveUp: (index: number) => void;
  onMoveDown: (index: number) => void;
  renderNameField: (column: C, index: number, onUpdate: ColumnUpdater<C>) => ReactElement;
  renderExtraFields?: (column: C, index: number, onUpdate: ColumnUpdater<C>) => ReactElement; // domain-specific fields (e.g., wrap toggle)
}
```

Each `ColumnEntry` renders a bordered card (`Box` with `border: 1`) containing:
1. **Header row:** Column display name + ArrowUp/ArrowDown/Delete icon buttons
2. **Name + Header row:** Side-by-side text fields (name via `renderNameField` prop, header as plain `TextField`)
3. **Enable sorting:** `FormControlLabel` with `Checkbox` — checked by default (`enableSorting !== false`)
4. **Sort mode + Default sort:** Side-by-side `Select` dropdowns

The list uses a stable `idCounterRef` + `idsRef` pattern for React keys (not array index) to avoid unmount/remount on reorder.

**Alert-table-specific wrapper — `alertmanager/src/plugins/alert-table/AlertTableColumnsEditor.tsx`:**

```typescript
const SORT_MODE_LABELS: Record<ColumnSortMode, string> = {
  alphabetical: 'Alphabetical',
  numeric: 'Numeric',
  severity: 'Severity (critical → other)',
};

export function AlertTableColumnsEditor(props: OptionsEditorProps<AlertTableOptions>): ReactElement {
  const { value, onChange } = props;
  const columns = value.columns ?? [];

  const handleAddColumn = useCallback((): void => {
    onChange(produce(value, (draft) => {
      if (!draft.columns) draft.columns = [];
      draft.columns.push({ name: 'severity' });
    }));
  }, [value, onChange]);

  // handleRemoveColumn, handleUpdateColumn, handleMoveUp, handleMoveDown
  // all use immer's `produce` for immutable state updates

  return (
    <ColumnsEditor<ColumnDefinition>
      columns={columns}
      description="Status and Alert Name are always shown. Add extra columns below."
      sortModeLabels={SORT_MODE_LABELS}
      defaultSortMode="alphabetical"
      getDisplayName={(col) => col.header || col.name || 'New column'}
      getHeaderPlaceholder={(col) => col.name || 'Column header'}
      onAdd={handleAddColumn}
      onRemove={handleRemoveColumn}
      onUpdate={handleUpdateColumn}
      onMoveUp={handleMoveUp}
      onMoveDown={handleMoveDown}
      renderNameField={(col, index, onUpdate) => (
        <TextField label="Label key" value={col.name}
          onChange={(e) => onUpdate(index, (draft) => { draft.name = e.target.value; })}
          size="small" fullWidth />
      )}
    />
  );
}
```

**Column data model — `alertmanager/src/plugins/alert-table/alert-table-model.ts`:**

```typescript
export type SortDirection = 'asc' | 'desc';
export type ColumnSortMode = 'alphabetical' | 'numeric' | 'severity';

export interface ColumnDefinition {
  name: string;
  header?: string;
  enableSorting?: boolean;
  sort?: SortDirection;
  sortMode?: ColumnSortMode;
}

export interface AlertTableOptions {
  defaultGroupBy?: string[];
  columns?: ColumnDefinition[];
  // ... other options
}
```

**Sorting — `alertmanager/src/plugins/alert-table/alert-table-sorting.ts`:**

```typescript
export interface SortState {
  columnName: string;
  direction: SortDirection;
  mode: ColumnSortMode;
}

// Comparators by mode: compareAlphabetical, compareNumeric, compareSeverity
// Main function: compareAlertsByColumn(a, b, sort) → number

export function compareAlertsByColumn(a: Alert, b: Alert, sort: SortState): number {
  const va = a.labels[sort.columnName];
  const vb = b.labels[sort.columnName];
  // switch on sort.mode, multiply by direction
}
```

**Sort state in AlertTablePanel — `alertmanager/src/plugins/alert-table/AlertTablePanel.tsx:426-440`:**

```typescript
const initialSort = useMemo<SortState | null>(() => {
  const col = columnDefs.find((c) => c.sort && c.enableSorting !== false);
  if (!col?.sort) return null;
  return { columnName: col.name, direction: col.sort, mode: col.sortMode ?? 'alphabetical' };
}, [columnDefs]);
const [sortState, setSortState] = useState<SortState | null>(initialSort);

const handleSortClick = useCallback((col: ColumnDefinition): void => {
  setSortState((prev) => {
    if (prev?.columnName === col.name) {
      return prev.direction === 'asc' ? { ...prev, direction: 'desc' } : null;
    }
    return { columnName: col.name, direction: 'asc', mode: col.sortMode ?? 'alphabetical' };
  });
}, []);
```

**Column header rendering with sort labels — `AlertTablePanel.tsx:621-638`:**

```typescript
{columnDefs.map((col) => (
  <TableCell key={col.name}
    sortDirection={sortState?.columnName === col.name ? sortState.direction : false}>
    {col.enableSorting !== false ? (
      <TableSortLabel
        active={sortState?.columnName === col.name}
        direction={sortState?.columnName === col.name ? sortState.direction : 'asc'}
        onClick={() => handleSortClick(col)}>
        {col.header ?? col.name}
      </TableSortLabel>
    ) : (col.header ?? col.name)}
  </TableCell>
))}
```

**Column cell rendering in AlertRow — `AlertTablePanel.tsx:156-168`:**

```typescript
{columnDefs.map((col) => {
  const value = alert.labels[col.name];
  // render with color mapping if configured, else plain text
  return <TableCell key={col.name}>{value ?? ''}</TableCell>;
})}
```

**CUE schema — `alertmanager/schemas/alert-table/alert-table.cue:19-25`:**

```cue
columns?: [...close({
  name:           string
  header?:        string
  enableSorting?: bool
  sort?:          "asc" | "desc"
  sortMode?:      "alphabetical" | "numeric" | "severity"
})]
```

**Plugin registration — `alertmanager/src/plugins/alert-table/AlertTable.ts:28-33`:**

```typescript
panelOptionsEditorComponents: [
  { label: 'General', content: AlertTableOptionsEditor },
  { label: 'Columns', content: AlertTableColumnsEditor },
  { label: 'Labels', content: AlertTableLabelsEditor },
  { label: 'Deduplication', content: AlertTableDeduplicationEditor },
],
```

### Key differences from table plugin's column editor

| Aspect | Table plugin | Alertmanager PR #647 |
| ------ | ------------ | -------------------- |
| Layout | Expand/collapse panels per column | Flat cards, always fully visible |
| Reorder | `DragButton` + `useDragAndDropMonitor` + `DragAndDropElement` from `@perses-dev/components` | Simple ArrowUp/ArrowDown `IconButton`s with stable ref-based key tracking |
| State updates | Direct array mutation (`onChange(updatedColumns)`) | `immer` `produce` for immutable drafts |
| Column fields | name, header, headerDescription, cellDescription, plugin, format, align, enableSorting, sort, width, hide, cellSettings, dataLink | name, header, enableSorting, sort, sortMode |
| Name field | Customizable (text input) | Customizable via `renderNameField` prop — caller controls the input |
| Generic | No (hardcoded `ColumnSettings` type) | Yes — `ColumnsEditor<C extends BaseColumnDefinition>` |
| React keys | Array index (`key={i}`) | Stable `idCounterRef`/`idsRef` pattern |
| Dependencies | `@perses-dev/components` drag utilities | Only MUI + `@perses-dev/components` `OptionsEditorGroup` |

**We follow the alertmanager pattern** because: (1) the spec says "similar to the alertmanager table", (2) it's simpler with fewer dependencies, (3) the generic `ColumnsEditor<C>` component can be reused directly, and (4) the flat card layout is appropriate for our fewer column properties.

## Changes

### Phase 1: Types, Model, and CUE Schema

**Dependency:** None
**Parallel with:** None

#### Files Modified

| File | Change |
| ---- | ------ |
| `logstable/src/model.ts` | Add `LogsColumnDefinition`, `LogsColumnSortMode`, `SortDirection` types and `columns` to `LogsTableOptions` |
| `logstable/schemas/logstable.cue` | Add column settings definition to spec |

#### Details

**TypeScript types** (`model.ts`):

Follow the alertmanager's `ColumnDefinition` pattern. The logs table needs `alphabetical` and `numeric` sort modes (no `severity` — that's alert-specific). Add a `timestamp` sort mode for the built-in timestamp column:

```typescript
export type SortDirection = 'asc' | 'desc';

export type LogsColumnSortMode = 'alphabetical' | 'numeric' | 'timestamp';

export interface LogsColumnDefinition {
  name: string;              // 'timestamp', 'line', or a label key
  header?: string;           // Display name. Defaults to name if unset
  enableSorting?: boolean;   // Default true (same as alertmanager)
  sort?: SortDirection;      // Default sort direction for this column
  sortMode?: LogsColumnSortMode;  // How to compare values. Default: 'alphabetical'
  allowWrap?: boolean;       // When true, content wraps (pre-wrap). When false, overflow hidden + ellipsis. Default: false
}
```

The `name` field identifies the column source:
- `"timestamp"` — the log entry's timestamp (sortMode defaults to `'timestamp'`)
- `"line"` — the log entry's message/line
- Any other string — treated as a label key from `LogEntry.labels`

Add to `LogsTableOptions`:

```typescript
export interface LogsTableOptions {
  // ... existing fields ...
  columns?: LogsColumnDefinition[];
}
```

When `columns` is `undefined` or empty, the panel falls back to current default behavior (timestamp if `showTime=true`, then log line). When `columns` is defined, it determines exactly which columns are shown and in what order, overriding `showTime`.

**CUE schema** (`logstable.cue`):

```cue
package model

import (
  "github.com/perses/shared/cue/common"
)

kind: "LogsTable"
spec: close({
  allowWrap?:     bool
  enableDetails?: bool
  showTime?:      bool
  columns?: [...close({
    name:           string
    header?:        string
    enableSorting?: bool
    sort?:          "asc" | "desc"
    sortMode?:      "alphabetical" | "numeric" | "timestamp"
    allowWrap?:     bool
  })]
  selection?:     common.#selection
  actions?:       common.#actions
})
```

#### Phase 1 Verification

- `cd ./projects/perses-plugins/logstable && npx tsc --noEmit` — TypeScript compiles without errors
- CUE schema validates: `cd ./projects/perses-plugins/logstable && cue vet ./schemas/logstable.cue`

---

### Phase 2: Column Editor UI

**Dependency:** Phase 1
**Parallel with:** None

#### Files Modified

| File | Change |
| ---- | ------ |
| `logstable/src/components/ColumnsEditor.tsx` | **New file.** Generic `ColumnsEditor<C>` component — copy from alertmanager's `ColumnsEditor.tsx` |
| `logstable/src/LogsTableColumnsEditor.tsx` | **New file.** Logs-specific wrapper that uses `ColumnsEditor<LogsColumnDefinition>` |
| `logstable/src/LogsTable.ts` | Add "Columns" tab to `panelOptionsEditorComponents` |

#### Details

##### ColumnsEditor (generic, reusable)

Copy `alertmanager/src/components/ColumnsEditor.tsx` from PR #647 into `logstable/src/components/ColumnsEditor.tsx`. Add one extension: an optional `renderExtraFields` prop on `ColumnsEditorProps` (and pass it through to `ColumnEntry`). When provided, `ColumnEntry` renders the extra fields after the sort mode/default sort row. This keeps the generic component reusable while allowing the logs table to add a wrap toggle. The alertmanager doesn't pass `renderExtraFields`, so its behavior is unchanged.

Exports: `BaseColumnDefinition`, `ColumnUpdater`, `ColumnsEditorProps`, `ColumnsEditor`.

##### LogsTableColumnsEditor (logs-specific wrapper)

Create `logstable/src/LogsTableColumnsEditor.tsx` following the `AlertTableColumnsEditor` pattern:

```typescript
import { Checkbox, FormControlLabel, TextField } from '@mui/material';
import { OptionsEditorProps } from '@perses-dev/plugin-system';
import { produce } from 'immer';
import { ReactElement, useCallback } from 'react';
import { ColumnsEditor } from './components/ColumnsEditor';
import { LogsTableOptions, LogsColumnDefinition, LogsColumnSortMode } from './model';

const SORT_MODE_LABELS: Record<LogsColumnSortMode, string> = {
  alphabetical: 'Alphabetical',
  numeric: 'Numeric',
  timestamp: 'Timestamp',
};

export function LogsTableColumnsEditor(props: OptionsEditorProps<LogsTableOptions>): ReactElement {
  const { value, onChange } = props;
  const columns = value.columns ?? [];

  const handleAddColumn = useCallback((): void => {
    onChange(produce(value, (draft) => {
      if (!draft.columns) draft.columns = [];
      draft.columns.push({ name: '' });
    }));
  }, [value, onChange]);

  const handleRemoveColumn = useCallback((index: number): void => {
    onChange(produce(value, (draft) => {
      draft.columns?.splice(index, 1);
    }));
  }, [value, onChange]);

  const handleUpdateColumn = useCallback(
    (index: number, updater: (draft: LogsColumnDefinition) => void): void => {
      onChange(produce(value, (draft) => {
        const column = draft.columns?.[index];
        if (column) updater(column);
      }));
    }, [value, onChange]);

  const handleMoveUp = useCallback((index: number): void => {
    if (index <= 0) return;
    onChange(produce(value, (draft) => {
      if (!draft.columns) return;
      const item = draft.columns.splice(index, 1)[0]!;
      draft.columns.splice(index - 1, 0, item);
    }));
  }, [value, onChange]);

  const handleMoveDown = useCallback((index: number): void => {
    onChange(produce(value, (draft) => {
      if (!draft.columns || index >= draft.columns.length - 1) return;
      const item = draft.columns.splice(index, 1)[0]!;
      draft.columns.splice(index + 1, 0, item);
    }));
  }, [value, onChange]);

  return (
    <ColumnsEditor<LogsColumnDefinition>
      columns={columns}
      description="Timestamp and Log line are shown by default. Add columns below to customize which columns are visible and their order."
      sortModeLabels={SORT_MODE_LABELS}
      defaultSortMode="alphabetical"
      getDisplayName={(col) => col.header || col.name || 'New column'}
      getHeaderPlaceholder={(col) => col.name || 'Column header'}
      onAdd={handleAddColumn}
      onRemove={handleRemoveColumn}
      onUpdate={handleUpdateColumn}
      onMoveUp={handleMoveUp}
      onMoveDown={handleMoveDown}
      renderNameField={(col, index, onUpdate) => (
        <TextField
          label="Column name"
          value={col.name}
          onChange={(e) => onUpdate(index, (draft) => { draft.name = e.target.value; })}
          size="small"
          fullWidth
          helperText="Use 'timestamp', 'line', or a label key"
        />
      )}
      renderExtraFields={(col, index, onUpdate) => (
        <FormControlLabel
          control={
            <Checkbox
              checked={col.allowWrap ?? false}
              onChange={(e) =>
                onUpdate(index, (draft) => {
                  draft.allowWrap = e.target.checked || undefined;
                })
              }
              size="small"
            />
          }
          label="Wrap content"
        />
      )}
    />
  );
}
```

##### Plugin registration update

In `LogsTable.ts`, add the new tab:

```typescript
panelOptionsEditorComponents: [
  { label: 'Settings', content: LogsTableSettingsEditor },
  { label: 'Columns', content: LogsTableColumnsEditor },
  { label: 'Item Actions', content: LogsTableItemSelectionActionsEditor },
],
```

##### Dependency: `immer`

Check if `immer` is already a dependency of `logstable/package.json`. The alertmanager plugin uses `immer`'s `produce` for all state updates in the columns editor. If not present, add it as a dependency (it's already used by other plugins in the monorepo).

#### Phase 2 Verification

- `cd ./projects/perses-plugins/logstable && npx tsc --noEmit` — TypeScript compiles
- Manual: open the logs table panel editor and verify the "Columns" tab appears
- Manual: add, remove, reorder columns in the editor; verify options are persisted in the panel spec JSON

---

### Phase 3: Column Rendering and Sorting

**Dependency:** Phase 1
**Parallel with:** Phase 2 (touches different files)

#### Files Modified

| File | Change |
| ---- | ------ |
| `logstable/src/components/LogRow/LogRow.tsx` | Accept resolved column definitions, render dynamic columns |
| `logstable/src/components/LogRow/LogsStyles.tsx` | Make `LogRowContent` grid template dynamic — accept `gridTemplateColumns` string prop |
| `logstable/src/components/VirtualizedLogsList.tsx` | Add header row, sort state, resolve column config from spec, apply sort |
| `logstable/src/components/LogsTableHeader.tsx` | **New file.** Header row with column names and `TableSortLabel`-style sort indicators |
| `logstable/src/components/LogRow/LogLabelCell.tsx` | **New file.** Renders a label value cell |
| `logstable/src/components/logs-table-sorting.ts` | **New file.** Sort comparators following alertmanager's `alert-table-sorting.ts` pattern |
| `logstable/src/LogsTableComponent.tsx` | Remove hardcoded timestamp sort — let VirtualizedLogsList handle sorting |

#### Details

##### Column resolution logic

In `VirtualizedLogsList`, resolve the effective columns from `spec.columns`:

```typescript
interface ResolvedColumn {
  name: string;
  header: string;
  type: 'timestamp' | 'line' | 'label';
  enableSorting: boolean;
  sortMode: LogsColumnSortMode;
  allowWrap: boolean;
  width?: number;
}
```

When `spec.columns` is undefined or empty, produce default columns:
- If `spec.showTime !== false`: `{ name: 'timestamp', header: 'Timestamp', type: 'timestamp', enableSorting: true, sortMode: 'timestamp' }`
- Always: `{ name: 'line', header: 'Log line', type: 'line', enableSorting: false, sortMode: 'alphabetical' }`

When `spec.columns` is defined, map each entry to a `ResolvedColumn`, determining `type` from the `name` field. The `sortMode` defaults to `'timestamp'` for the timestamp column and `'alphabetical'` for all others, unless explicitly set.

##### Grid layout changes

`LogRowContent` currently has a hardcoded grid. The grid template needs to become dynamic:

- Expand button column: `16px` (always present when `isExpandable`)
- For each resolved column:
  - `timestamp` type: `minmax(160px, max-content)`
  - `label` type: `minmax(80px, max-content)`
  - `line` type: `1fr` (fills remaining space)
- Copy/action area: `min-content`

`LogsStyles.tsx:LogRowContent` will accept a `gridTemplateColumns` string prop instead of computing it internally:

```typescript
export const LogRowContent = styled(Box, {
  shouldForwardProp: (prop) =>
    prop !== 'gridTemplateColumns' && prop !== 'isHighlighted' && prop !== 'isSelected',
})<{ gridTemplateColumns: string; isHighlighted?: boolean; isSelected?: boolean }>(
  ({ theme, gridTemplateColumns, isHighlighted, isSelected }) => ({
    display: 'grid',
    gridTemplateColumns,
    // ... rest stays the same
  })
);
```

The parent (`VirtualizedLogsList`) computes the grid template string from the resolved columns and passes it down.

##### LogRow changes

`LogRow` currently renders timestamp and log line in fixed positions. With dynamic columns:

1. Accept `resolvedColumns: ResolvedColumn[]` and `gridTemplateColumns: string` props
2. After the optional expand button, iterate over `resolvedColumns` and render each:
   - `timestamp` → `<LogTimestamp timestamp={log.timestamp} />`
   - `line` → `<LogText>` with ANSI rendering (existing code, wrapped in a Box with copy menu and action buttons)
   - `label` → `<LogLabelCell value={log.labels[column.name]} allowWrap={allowWrap} />` (new component)
3. The copy menu and action buttons remain associated with the log line area

The `LogLabelCell` component renders a monospace text value styled consistently with `LogText`. For missing labels, it renders `—` (em-dash).

**Wrap behavior per column** (`allowWrap` on `LogsColumnDefinition`):
- `allowWrap: false` (default): `overflow: hidden; text-overflow: ellipsis; white-space: nowrap` — long values are truncated with an ellipsis. The full value is shown in a tooltip on hover.
- `allowWrap: true`: `word-break: break-word; white-space: pre-wrap; overflow: visible` — content wraps to new lines (same CSS as the existing `LogText` wrap mode in `LogsStyles.tsx:67-74`).

This mirrors the existing global `allowWrap` option on `LogsTableOptions` but at the per-column level. The `line` column inherits the global `spec.allowWrap` when no per-column `allowWrap` is set; label columns default to `false` (no wrap) unless explicitly toggled.

##### Expanded details panel (enableDetails)

Currently, when a row is expanded, the `Collapse` section in `LogRow.tsx:343-363` renders `LogDetailsTable` (all labels as key-value pairs). This `Collapse` is a sibling of `LogRowContent` inside `LogRowContainer` — it sits below the grid row. The current code uses an inner alignment grid that matches the old fixed columns:

```typescript
// Current code (LogRow.tsx:345-361) — broken with dynamic columns
<Box sx={{
  display: 'grid',
  gridTemplateColumns: !showTime ? '1fr' : '8px minmax(160px, max-content) 1fr',
  gap: '12px',
}}>
  {showTime && (<><Box /><Box /></>)}  // spacer boxes to align under log line
  <Box><LogDetailsTable log={log.labels} /></Box>
</Box>
```

This hardcoded alignment grid must be replaced. With dynamic columns, the details panel should **span the full width** of all columns rather than trying to align under a specific column:

```typescript
// New code — details span all columns
<Collapse in={isExpanded} timeout={200}>
  <Box sx={{ padding: '8px', paddingLeft: isExpandable ? '36px' : '8px' }}>
    <LogDetailsTable log={log.labels} />
  </Box>
</Collapse>
```

The `paddingLeft` indents the details to visually nest under the row content (past the expand chevron). The details table itself takes the full container width, displaying all extracted labels regardless of which columns are configured. This is the correct behavior because:

1. Custom columns show a **subset** of labels inline — the details panel shows the **complete** set
2. Trying to column-align the details under dynamic columns adds complexity for no UX benefit
3. The full-width layout is consistent with how details/expansion panels work in other Perses table plugins

##### Header row

`LogsTableHeader` renders a fixed header row above the Virtuoso list, using the same grid template. Following the alertmanager's pattern (`AlertTablePanel.tsx:619-638`), each header cell shows the column's `header` text. If `enableSorting` is true, uses MUI's `TableSortLabel` for sort direction indicators and click handling.

The header row lives outside the Virtuoso component (above it in the flex column) so it doesn't scroll.

```typescript
interface LogsTableHeaderProps {
  resolvedColumns: ResolvedColumn[];
  gridTemplateColumns: string;
  isExpandable: boolean;
  sortState: SortState | null;
  onSortClick: (column: ResolvedColumn) => void;
}
```

##### Sorting

Follow the alertmanager's sorting pattern exactly.

**Sort state type** (`logs-table-sorting.ts`):

```typescript
export interface SortState {
  columnName: string;
  direction: SortDirection;
  mode: LogsColumnSortMode;
}
```

**Sort comparators** (`logs-table-sorting.ts`):

```typescript
export function compareLogsByColumn(a: LogEntry, b: LogEntry, sort: SortState): number {
  let va: string | number | undefined;
  let vb: string | number | undefined;

  if (sort.columnName === 'timestamp') {
    va = a.timestamp;
    vb = b.timestamp;
  } else if (sort.columnName === 'line') {
    va = a.line;
    vb = b.line;
  } else {
    va = a.labels[sort.columnName];
    vb = b.labels[sort.columnName];
  }

  // switch on sort.mode: 'timestamp' (numeric comparison on timestamp),
  // 'numeric' (parseFloat), 'alphabetical' (localeCompare)
  // multiply result by direction multiplier
}
```

**Sort state management in VirtualizedLogsList** (following `AlertTablePanel.tsx:426-440`):

```typescript
const initialSort = useMemo<SortState | null>(() => {
  if (!resolvedColumns.length) return { columnName: 'timestamp', direction: 'desc', mode: 'timestamp' };
  const col = resolvedColumns.find((c) => c.enableSorting && /* has default sort from spec */);
  if (col) return { columnName: col.name, direction: col.sort!, mode: col.sortMode };
  return { columnName: 'timestamp', direction: 'desc', mode: 'timestamp' };
}, [resolvedColumns]);

const [sortState, setSortState] = useState<SortState | null>(initialSort);

const handleSortClick = useCallback((col: ResolvedColumn): void => {
  setSortState((prev) => {
    if (prev?.columnName === col.name) {
      return prev.direction === 'asc' ? { ...prev, direction: 'desc' } : null;
    }
    return { columnName: col.name, direction: 'asc', mode: col.sortMode };
  });
}, []);

// Apply sort to logs
const sortedLogs = useMemo(() => {
  if (!sortState) return logs;
  return [...logs].sort((a, b) => compareLogsByColumn(a, b, sortState));
}, [logs, sortState]);
```

**Move sorting out of LogsTableComponent:** Currently `LogsTableComponent.tsx:23-25` sorts by timestamp. Remove that sort and let `VirtualizedLogsList` handle it via `sortState`, so sorting is unified and controllable by the user.

##### Component prop threading

Update the component chain to pass columns through:

1. `LogsTableComponent` → `LogsList`: already passes `spec`
2. `LogsList` → `VirtualizedLogsList`: already passes `spec`
3. `VirtualizedLogsList`: resolves columns from `spec`, manages sort state, computes grid template, renders header + list
4. `VirtualizedLogsList` → `LogRow`: add `resolvedColumns` and `gridTemplateColumns` props
5. `LogRow`: renders cells dynamically

#### Phase 3 Verification

- `cd ./projects/perses-plugins/logstable && npx tsc --noEmit` — TypeScript compiles
- Manual: with no `columns`, verify default behavior is unchanged (timestamp + log line, sorted by timestamp desc)
- Manual: add columns (e.g., `app`, `status`, `method`) and verify they appear with label values
- Manual: click a sortable column header and verify sort order toggles (asc → desc → none)
- Manual: verify expanded details panel spans full width across all columns and shows all labels (not just configured columns)

---

### Phase 4: Integration Testing and Polish

**Dependency:** Phase 2, Phase 3
**Parallel with:** None

#### Files Modified

| File | Change |
| ---- | ------ |
| `logstable/src/LogsTable.ts` | Verify `createInitialOptions` still works, ensure backward compat |
| `logstable/src/utils/copyHelpers.ts` | May need updates if copy format should include custom column values |

#### Details

- Verify backward compatibility: panels with no `columns` in their spec should render identically to before
- Verify the column editor persists settings correctly in the panel spec JSON
- Verify `showTime` is respected when no `columns` is defined
- Verify copy functionality still works with custom columns (log copy should include all labels regardless of visible columns)
- Test with different datasource results that have varying label sets (some entries missing labels that are configured as columns)
- Verify row expand/collapse and details table still works alongside custom columns

#### Phase 4 Verification

- `cd ./projects/perses-plugins && npm run build` — full build passes
- Run the dev server with a dashboard containing a logs table panel and test all interactions

## PR Strategy

| PR | Repository | Branch | Description | Dependencies |
| -- | ---------- | ------ | ----------- | ------------ |
| 1  | perses/plugins | `feat/logstable-column-settings` | Add customizable column support to logs table panel | None |

All changes are in a single repository (perses-plugins) and can be delivered in a single PR. The generic `ColumnsEditor` component is copied into the logstable plugin (same pattern as alertmanager keeping its own copy). No upstream `@perses-dev/*` package changes are needed.

## Verification

End-to-end verification mapped to acceptance criteria:

- **"The logs table supports customizing the columns to show in the table, including the log line, timestamp and extracted labels"**
  - Create a logs table panel, open editor, go to "Columns" tab
  - Add columns: `timestamp`, `line`, `app`, `status`
  - Verify all four columns render with correct values from log entries
  - Remove `timestamp` column, verify it disappears
  - Verify labels with missing values show `—`

- **"The edit should be similar to the alertmanager table"**
  - Compare the column editor UI side-by-side with the alertmanager table's column editor
  - Verify: add/delete columns, ArrowUp/ArrowDown reorder, header text, enable sorting checkbox, sort mode and default sort selects all work identically
  - Verify the same card-based layout with bordered cards and dividers between entries

- **"The user should be able to save the configuration for the panel"**
  - Add column settings, save the dashboard
  - Reload the page, verify column settings persist
  - Inspect the dashboard JSON, verify `columns` array is in the panel spec

- **"The column editor should be added as a new tab in the panel editor"**
  - Open the logs table panel editor
  - Verify three tabs: "Settings", "Columns", "Item Actions"

- **Backward compatibility**
  - Load an existing dashboard with a logs table panel that has no `columns`
  - Verify it renders with default timestamp + log line (unchanged)

- **Column sort**
  - Configure a column with `enableSorting: true`
  - Click the column header, verify rows sort ascending
  - Click again, verify descending
  - Click again, verify returns to default sort (timestamp desc)

## Risks

| Risk | Impact | Mitigation |
| ---- | ------ | ---------- |
| Performance with many label columns | Grid with many columns + virtualized list could degrade scroll performance | Keep column cells lightweight (plain text), avoid heavy re-renders. Test with 10+ columns and 1000+ log entries |
| Variable label sets across log entries | Different log entries may have different labels, causing empty cells | Render `—` for missing labels. Document this behavior |
| Grid layout breaks with very long label values | Long label values could overflow or push columns off-screen | Use `overflow: hidden; text-overflow: ellipsis` on label cells. Tooltip on hover for full value |
| `showTime` and `columns` interaction confusion | Users may set `showTime: false` but also add a `timestamp` column in settings | When `columns` is defined, it takes precedence over `showTime`. Document this |
| CUE schema backward compatibility | Existing dashboards with logs table panels must validate against updated schema | `columns` is optional, so existing specs without it remain valid |
| `immer` dependency | logstable may not currently depend on `immer` | Check `package.json`; add if missing. `immer` is already used by other plugins in the monorepo (alertmanager, table) |
