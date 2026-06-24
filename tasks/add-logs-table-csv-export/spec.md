# Spec: Add support for CSV export to logs table

## Related projects and branches

- perses-plugins: branch `main`
- perses-shared: branch `main`

## Description

The logs table, that supports many datsources has support to export logs as JSON. We want to add support to export logs as CSV. This will allow users
to download logs in a more common format that can be used in other tools.

## Acceptance criteria

- The logs table supports exporting logs as CSV

## Hints

- The export functionallity is added into a panel through actions that the plugins can register. The logs table plugin already has an action to export
  logs as JSON, you can use that as a reference to add support for CSV.
- The format for the CSV should include the timestamp and a time field in ISO 8601 format, log line and extracted labels. The log line should be the
  raw log message, and the extracted labels should be included as additional columns in the CSV. Bear in mind the scape characters and the CSV format.
  Reuse existing libraries for CSV generation if possible.
