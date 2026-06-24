# Execution: Customize Columns in Logs Table

> Results are annotated inline: `-- **value**` for discovered values, `-- **passes/FAILED**` for verification.

## Phase 1: Types, Model, and CUE Schema
Depends on: nothing | Parallel with: none | Type: implementation | Projects: perses-plugins (logstable)

### 1a. TypeScript types
- [x] Add `SortDirection`, `LogsColumnSortMode`, `LogsColumnDefinition` types to model - `logstable/src/model.ts`
- [x] Add `columns?: LogsColumnDefinition[]` to `LogsTableOptions` - `logstable/src/model.ts`

### 1b. CUE schema
- [x] Add `columns` field with column definition to CUE spec - `logstable/schemas/logstable.cue`

### Phase 1 Verification
- [x] `cd logstable && npm run type-check` — **passes**
- [ ] CUE schema validates — skipped (cue CLI not available locally, schema follows table.cue pattern)

---

## Phase 2: Column Editor UI
Depends on: Phase 1 | Parallel with: none | Type: implementation | Projects: perses-plugins (logstable)

### 2a. Generic ColumnsEditor component
- [x] Create generic `ColumnsEditor<C>` component (adapted from alertmanager PR #647 with `renderExtraFields` extension) - `logstable/src/components/ColumnsEditor.tsx` -- **20 tests**
- [x] Tests for ColumnsEditor - `logstable/src/components/ColumnsEditor.test.tsx`

### 2b. LogsTableColumnsEditor wrapper
- [x] Write tests for LogsTableColumnsEditor (add/remove/reorder/update columns, wrap toggle) - `logstable/src/LogsTableColumnsEditor.test.tsx` -- **13 tests**
- [x] Create `LogsTableColumnsEditor` wrapper with sort mode labels and wrap content checkbox - `logstable/src/LogsTableColumnsEditor.tsx`

### 2c. Plugin registration and dependency
- [x] Add `immer` to logstable `dependencies` in package.json - `logstable/package.json`
- [x] Add "Columns" tab to `panelOptionsEditorComponents` in `LogsTable.ts` - `logstable/src/LogsTable.ts`

### Phase 2 Verification
- [x] `cd logstable && npm test` — **33 tests pass** (20 ColumnsEditor + 13 LogsTableColumnsEditor)
- [x] `cd logstable && npm run type-check` — **passes** (6 errors are pre-existing in LogRow.test.tsx from Phase 3 parallel work)

---
## Phases 2 and 3 touch different files and can run in parallel after Phase 1
---

## Phase 3: Column Rendering and Sorting
Depends on: Phase 1 | Parallel with: Phase 2 (different files) | Type: implementation | Projects: perses-plugins (logstable)

### 3a. Sort comparators
- [x] Write tests for sort comparators (alphabetical, numeric, timestamp modes; asc/desc directions) - `logstable/src/components/logs-table-sorting.test.ts` -- **14 tests**
- [x] Implement `SortState`, `compareLogsByColumn`, and mode-specific comparators - `logstable/src/components/logs-table-sorting.ts`

### 3b. Column resolution and grid template
- [x] Write tests for `resolveColumns` (default columns, custom columns, showTime fallback, hidden columns) - `logstable/src/components/column-resolution.test.ts` -- **17 tests**
- [x] Implement `ResolvedColumn` interface and `resolveColumns` function - `logstable/src/components/column-resolution.ts`
- [x] Implement `buildGridTemplate` function from resolved columns - `logstable/src/components/column-resolution.ts`

### 3c. LogLabelCell component
- [x] Write tests for LogLabelCell (renders value, renders em-dash for missing, wrap vs no-wrap) - `logstable/src/components/LogRow/LogLabelCell.test.tsx` -- **5 tests**
- [x] Implement LogLabelCell component with wrap/ellipsis styling - `logstable/src/components/LogRow/LogLabelCell.tsx`

### 3d. LogsTableHeader component
- [x] Write tests for LogsTableHeader (renders column headers, sort indicators, click-to-sort) - `logstable/src/components/LogsTableHeader.test.tsx` -- **7 tests**
- [x] Implement LogsTableHeader with grid layout and MUI TableSortLabel - `logstable/src/components/LogsTableHeader.tsx`

### 3e. LogRow dynamic columns and details panel
- [x] Update LogRow tests for dynamic columns (resolvedColumns prop, label columns, expanded details spanning full width) - `logstable/src/components/LogRow/LogRow.test.tsx` -- **3 new tests**
- [x] Make `LogRowContent` accept `gridTemplateColumns` string prop instead of computing it - `logstable/src/components/LogRow/LogsStyles.tsx`
- [x] Update LogRow to accept `resolvedColumns` and `gridTemplateColumns`, render columns dynamically - `logstable/src/components/LogRow/LogRow.tsx`
- [x] Replace hardcoded details alignment grid with full-width span (paddingLeft indent) - `logstable/src/components/LogRow/LogRow.tsx`

### 3f. VirtualizedLogsList integration
- [x] Add header row, sort state management, column resolution, and sorted logs to VirtualizedLogsList - `logstable/src/components/VirtualizedLogsList.tsx`
- [x] Remove hardcoded timestamp sort from LogsTableComponent (sorting now in VirtualizedLogsList) - `logstable/src/LogsTableComponent.tsx`

### Phase 3 Verification
- [x] `cd logstable && npm run type-check` — **passes**
- [x] `cd logstable && npm test` — **128 tests pass** (10 suites)

---

## Phase 4: Integration Testing and Polish
Depends on: Phase 2, Phase 3 | Parallel with: none | Type: configuration | Projects: perses-plugins (logstable)

- [x] Verify `createInitialOptions` still works with no `columns` field - `logstable/src/LogsTable.ts` -- **confirmed, no `columns` in initial options**
- [x] Verify copy functionality still works with custom columns - `logstable/src/utils/copyHelpers.ts` -- **confirmed, uses `labels` directly (0 refs to `columns`)**
- [x] Run npm install at repo root to update lock file for immer dependency -- **done**

### Phase 4 Verification
- [x] `npm run build --workspace=logstable` — **passes** (33 files compiled, types emitted)
- [x] `cd logstable && npm test` — **128 tests pass** (10 suites)
- [x] `cd logstable && npm run lint` — **no lint errors**

---

## Summary

**Status:** Complete (4 of 4 phases done)

### What was built

- **Types & schema:** `LogsColumnDefinition` with name, header, enableSorting, sort, sortMode, allowWrap fields added to `LogsTableOptions.columns`. CUE schema updated.
- **Column editor UI:** Generic `ColumnsEditor<C>` component (adapted from alertmanager PR #647) + `LogsTableColumnsEditor` wrapper with wrap content toggle. Registered as "Columns" tab.
- **Column rendering:** Dynamic grid-based columns in LogRow, LogLabelCell for label values (with wrap/ellipsis per column), LogsTableHeader with sticky sort indicators, full-width details panel on expand.
- **Sorting:** `compareLogsByColumn` with alphabetical/numeric/timestamp modes. Sort state in VirtualizedLogsList with click-to-toggle (asc → desc → none).
- **Backward compat:** No `columns` in initial options = default behavior (timestamp + log line). `showTime` respected when no columns configured.

### Test coverage

- 128 tests total (79 new across 6 test files)
- Sort comparators: 14 tests
- Column resolution: 17 tests
- LogLabelCell: 5 tests
- LogsTableHeader: 7 tests
- ColumnsEditor: 20 tests
- LogsTableColumnsEditor: 13 tests
- LogRow: 3 new tests (existing tests updated)

### Git state

```
Branch: feat/logstable-column-settings (2 commits ahead of main)
  5406f6d feat: add column rendering, sorting, and header for logs table plugin
  15ca7ac feat: add column editor UI for logs table plugin (Phase 2)
```

### Outstanding items

- [ ] Push branch and create PR on perses/plugins
- [ ] Manual testing with a live Perses instance (dev server + real log data)

### Notes

- `immer` added as a dependency to logstable/package.json (was already installed at workspace root)
- `jest.shared.ts` was fixed for Node 25 ESM compatibility (pre-existing issue with `__dirname`)
- Phase 1 types were committed as part of Phase 2/3 agent commits (direct execution, no separate commit)
