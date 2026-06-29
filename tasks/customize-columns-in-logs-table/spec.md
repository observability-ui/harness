# Spec: Customize columns in logs table

## Related projects and branches

- perses-plugins: branch `main`

## Description

The logs table, that supports many datsources includes a column for the log line and a column for the timestamp. We want to add support to add custom
columns from labels or remove the default columns. This will allow users to see the extracted labels in a more structured way and be able to sort and
filter by them.

## Acceptance criteria

- The logs table supports customizing the columns to show in the table, including the log line, timestamp and extracted labels.
- The edit should be similar to the alertmanager table, where the user can select which columns to show and in which order. The user should be able to
  save the configuration for the panel.
- The column editor should be added as a new tab in the panel editor.

## Hints

- The alert manager column editor is in
  `https://github.com/perses/plugins/pull/647/changes#diff-e428cf7c445b42c208064b4b228d820f7329e1594d21b206c72f1e525707c087`
