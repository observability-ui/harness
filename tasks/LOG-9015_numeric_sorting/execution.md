# LOG-9015: Numeric Sorting — Execution Notes

## Data Source

- File: `openshift_logs(4).csv`
- 450 rows, 47 columns
- Namespace: `log-9015-numeric-sorting`
- Pod: `log-generator-546c895b84-kbrhs`
- Time range: `2026-05-21T15:53:59.851070415Z` to `2026-05-21T15:54:08.249655178Z` (~8.4 seconds)

## Key Columns Analyzed

| Column | Format | Description |
|---|---|---|
| `time` | Epoch nanoseconds (int) | When the log line was **generated** by the container |
| `_timestamp` | ISO 8601 with nanoseconds | Human-readable form of `time` |
| `openshift_sequence` | Epoch nanoseconds (int) | When the log line was **processed/ingested** by the collector |

## Analysis: `openshift_sequence` Logic

### Finding: `openshift_sequence` is the collector ingestion timestamp

`openshift_sequence` is a nanosecond-precision epoch timestamp representing when the log was processed/ingested by the collector — not when the log was originally emitted.

Both `time` and `openshift_sequence` are epoch nanoseconds:

```
time                = 1779378839851070415 → 2026-05-21T15:53:59.851070415Z
openshift_sequence  = 1779378840017768952 → 2026-05-21T15:54:00.017768952Z
```

The difference (`openshift_sequence - time`) represents the **processing latency** — how long it took the collector to pick up and assign a sequence number to each log line.

### Processing Latency Stats

| Metric | Value |
|---|---|
| Min diff | 3.618 ms |
| Max diff | 199.024 ms |
| Mean diff | 122.993 ms |

### Latency Over Time (sampled every 50 rows)

```
Row   0: time=2026-05-21T15:53:59.851070  diff=166.699 ms
Row  50: time=2026-05-21T15:53:59.854452  diff=164.694 ms
Row 100: time=2026-05-21T15:53:59.951438  diff=68.697 ms
Row 150: time=2026-05-21T15:54:01.955715  diff=127.035 ms
Row 200: time=2026-05-21T15:54:02.148726  diff=64.462 ms
Row 250: time=2026-05-21T15:54:04.152512  diff=134.175 ms
Row 300: time=2026-05-21T15:54:04.155525  diff=132.353 ms
Row 350: time=2026-05-21T15:54:06.160527  diff=196.847 ms
Row 400: time=2026-05-21T15:54:08.164453  diff=4.689 ms
```

The latency fluctuates, suggesting batch processing behavior by the collector.

### Uniqueness and Ordering

- **All 450 `openshift_sequence` values are unique** (vs. only 293 unique `time` values).
- **102 groups of rows share the same `time`** — `openshift_sequence` serves as a tie-breaker to distinguish them.
- **`time` is monotonically increasing** in the CSV — the file is sorted by log generation time.
- **`openshift_sequence` is NOT monotonically increasing** in the file (58 violations out of 449 transitions), because the CSV is sorted by `time`, not by processing order.
- Within same-`time` groups, the sequence gaps are **~15–40 microseconds**, consistent with the collector calling `clock_gettime()` as it processes each line in a batch.

### Consecutive Gaps in Sorted Sequences

```
Min gap:    15,682 ns
Max gap:    2,071,778,987 ns  (~2.07 seconds — likely a batch boundary)
Mean gap:   18,451,304 ns
Median gap: 18,736 ns
```

### Monotonicity Violations (first 5)

These occur because the file is ordered by `time` (generation), not by `openshift_sequence` (ingestion):

```
Row  2→3:  seq 1779378840017945533 > 1779378840017924496  (diff: 21,037 ns)
Row  3→4:  seq 1779378840017924496 > 1779378840017893876  (diff: 30,620 ns)
Row  5→6:  seq 1779378840017971952 > 1779378840017831296  (diff: 140,656 ns)
Row  7→8:  seq 1779378840018005034 > 1779378840017859323  (diff: 145,711 ns)
Row 10→11: seq 1779378840018095419 > 1779378840018077149  (diff: 18,270 ns)
```

## Conclusion

`openshift_sequence` = **nanosecond wall-clock time at ingestion**, used as a globally unique, roughly-monotonic sequence ID to order and deduplicate logs — especially when multiple log lines share the same `time` value.

This is relevant to LOG-9015 because sorting by `openshift_sequence` requires **numeric comparison** (not lexicographic), as these are 19-digit integer values. Lexicographic sorting would produce incorrect ordering.
