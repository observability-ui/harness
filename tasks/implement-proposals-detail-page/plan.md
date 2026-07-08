# Plan: Implement Proposals Detail Page

## Problem

The monitoring-plugin recently added a "View AI Investigation" action to the alert list kebab menu (PR #1014). This action navigates to
`/monitoring/v2/ai/proposals?name=<name>`, but the destination page does not exist yet. Users clicking the action see a blank page.

This plan implements the Proposal detail page that displays:

- Proposal header with name, phase status, and creation timestamp
- Root cause analysis (RCA) section from the AnalysisResult CR
- Remediation options with expandable details (diagnosis, proposed actions, risk, reversibility)
- Action buttons: Execute remediation, Stop analysis (during Analyzing phase)

The page is modeled after the [UX mockup](https://fkargbo.github.io/ux-prototypes/core/observe/ai-hub/plans/?perspective=core-platforms) but aligned
to the actual [Proposal CRD API](https://github.com/openshift/lightspeed-agentic-operator/tree/main/api/v1alpha1). Fields not present in the API
(e.g., "Tokens consumed") are omitted. The mockup's "Trigger domain" maps to the `agentic.openshift.io/source` label on the Proposal CR (e.g.,
`alertmanager`) and is displayed in the detail page header.

## Current State

| Component          | File / Location                                                 | Current Behavior                                                                              |
| ------------------ | --------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| Alert kebab action | `web/src/components/alerting/AlertList/AlertTableRow.tsx:63-80` | Navigates to `/monitoring/v2/ai/proposals?name=<name>` via query params — page does not exist |
| URL builder        | `web/src/components/hooks/usePerspective.tsx:322-332`           | `getProposalsUrl()` builds query-param URLs, no path-param variant                            |
| K8s model          | `web/src/components/console/models/index.ts:95-107`             | `ProposalModel` defined for `agentic.openshift.io/v1alpha1` — no AnalysisResult model         |
| AI proposals code  | `web/src/components/ai-proposals/`                              | `useProposalCheck.ts` fetches proposals by fingerprint for alert matching only                |
| Console extensions | `web/console-extensions.json`                                   | No route registered for `/monitoring/v2/ai/proposals/:name`                                   |
| Exposed modules    | `web/package.json:175-195`                                      | No ProposalDetailsPage module exposed                                                         |
| Translation keys   | `web/locales/en/plugin__monitoring-plugin.json`                 | Only "View AI Investigation" and "Loading investigations..." keys exist                       |

## Changes

### Phase 1: TypeScript types, K8s models, and data hooks

**Dependency:** None **Parallel with:** None

#### Files Modified

| File                                                      | Change                                                                                             |
| --------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `web/src/components/console/models/index.ts`              | Add `AnalysisResultModel`, `ExecutionResultModel`, `VerificationResultModel` K8s models            |
| `web/src/components/ai-proposals/constants.ts`             | Add `RESULT_LABEL_PROPOSAL` constant for result CR label selector                                  |
| `web/src/components/ai-proposals/api-types.ts` (new)      | Raw CRD types mirroring `lightspeed-agentic-operator/api/v1alpha1/` — changes when the API changes |
| `web/src/components/ai-proposals/types.ts` (new)          | UI-stable view model interfaces — what components consume; does NOT change when the API changes    |
| `web/src/components/ai-proposals/proposal-utils.ts` (new) | `getPhaseStatus()` — component-facing pure helper (no `api-types.ts` imports); phase labels handled by i18n |
| `web/src/components/ai-proposals/useProposal.ts` (new)    | Hook: K8s reads + mutations + **API→view-model mapping**. Single file to touch when API changes.   |

#### Details

##### Layer 1: Raw API types (`api-types.ts`)

These interfaces mirror `lightspeed-agentic-operator/api/v1alpha1/` exactly. **Only `useProposal.ts` imports this file.** Components never touch it.
When the alpha API changes fields/structure, this file and the mapping code in `useProposal.ts` change — nothing else.

```typescript
import { K8sResourceCommon } from '@openshift-console/dynamic-plugin-sdk';

// K8s condition type — matches metav1.Condition shape
interface Condition {
  type: string;
  status: 'True' | 'False' | 'Unknown';
  reason?: string;
  message?: string;
  lastTransitionTime?: string;
}

// --- Proposal CRD ---
interface ApiProposalSpec {
  request: string;
  targetNamespaces?: string[];
}
interface ApiProposalStatus {
  conditions?: Condition[];
  steps?: {
    analysis?: { conditions?: Condition[]; results?: { name: string; outcome: string }[] };
    execution?: { conditions?: Condition[]; results?: { name: string; outcome: string }[]; retryCount?: number };
    verification?: { conditions?: Condition[]; results?: { name: string; outcome: string }[] };
    escalation?: { conditions?: Condition[]; results?: { name: string; outcome: string }[] };
  };
}
interface ApiProposal extends K8sResourceCommon {
  spec: ApiProposalSpec;
  status?: ApiProposalStatus;
}

// --- AnalysisResult CRD ---
interface ApiRemediationOption {
  title: string;
  summary?: string;
  diagnosis?: { summary: string; confidence: string; rootCause: string };
  proposal?: {
    description: string;
    actions: { type: string; description: string }[];
    risk: string;
    reversible?: string;
    estimatedImpact: string;
    rollbackPlan?: { description: string; command?: string };
  };
  verification?: { description: string; steps?: { name: string; command?: string; expected?: string; type: string }[] };
}
interface ApiAnalysisResult extends K8sResourceCommon {
  spec: { proposalName: string };
  status?: { conditions?: Condition[]; options?: ApiRemediationOption[]; failureReason?: string };
}

// --- ExecutionResult CRD ---
interface ApiExecutionResult extends K8sResourceCommon {
  spec: { proposalName: string; retryIndex?: number };
  status?: {
    conditions?: Condition[];
    actionsTaken?: { type: string; description: string; outcome: string; output?: string; error?: string }[];
    verification?: { conditionOutcome: string; summary: string };
    failureReason?: string;
  };
}

// --- VerificationResult CRD ---
interface ApiVerificationResult extends K8sResourceCommon {
  spec: { proposalName: string; retryIndex?: number };
  status?: {
    conditions?: Condition[];
    checks?: { name: string; source: string; value: string; result: string }[];
    summary?: string;
    failureReason?: string;
  };
}
```

##### Layer 2: UI view models (`types.ts`)

These are what components consume. They are **UI-oriented and stable** — named for what the UI needs, not for what the API exposes. If the API renames
`diagnosis.rootCause` to `rootCauseAnalysis.cause`, only the mapping in `useProposal.ts` changes; every component still reads `rootCause.cause`.

```typescript
// --- Phase (shared between utils and components) ---
export type ProposalPhase = 'Pending' | 'Analyzing' | 'Proposed' | 'Executing'
  | 'Verifying' | 'Completed' | 'Failed' | 'Denied'
  | 'Escalating' | 'Escalated' | 'EmergencyStopped';

// --- Root cause analysis view model ---
export interface RootCauseView {
  cause: string;         // the one-line root cause statement
  detail: string;        // longer explanation / summary
  confidence?: string;   // 'Low' | 'Medium' | 'High' — rendered as-is, no enum coupling
}

// --- Remediation option view model ---
export interface RemediationOptionView {
  index: number;
  title: string;
  description: string;                 // merged from summary or proposal.description
  reversibility?: string;              // 'Reversible' | 'Irreversible' | 'Partially reversible' — rendered as-is
  risk?: string;                       // 'Low' | 'Medium' | 'High' | 'Critical' — rendered as-is
  estimatedImpact?: string;
  command?: string;                    // first actionable command for ClipboardCopy
  actions?: { type: string; description: string }[];
  rollbackDescription?: string;
  rollbackCommand?: string;
}

// --- Post-execution summary view model (Completed state) ---
export interface AuditEntry {
  label: string;       // e.g. "Applied", "System restored", "Execution time"
  value: string;       // pre-formatted string
  isLink?: boolean;    // if true, value is a URL
}

export interface PostExecutionView {
  originalRootCause: string;
  remediationDelta: string;
  outcome: string;                     // 'Improved' | 'Unchanged' | 'Degraded' — rendered as-is
  auditTrail: AuditEntry[];            // flexible key-value list — adapts to whatever the API provides
  verificationSummary?: string;
}

// --- Top-level view model returned by useProposal ---
// NOTE: This does NOT wrap trivial K8s metadata (name, createdAt). The hook returns
// the raw K8sResourceCommon for those fields — they are stable across API versions.
// Only non-trivial mapped data gets a view type.
export interface ProposalView {
  phase: ProposalPhase;
  request: string;                       // the original investigation request text (from spec.request)
  source?: string;                       // trigger domain — from label `agentic.openshift.io/source` (e.g., "alertmanager")
  failureReason?: string;                // populated when any step failed — system error explanation
  rootCause?: RootCauseView;             // populated once analysis completes with results
  options: RemediationOptionView[];      // empty during Analyzing or if analysis returned no options
  postExecution?: PostExecutionView;     // populated only in Completed state
}
```

**Key design decisions:**

- **No passthrough wrappers for stable K8s metadata.** `name` and `createdAt` come from `metadata` — standard K8s fields that won't change with the
  alpha API. The hook exposes the raw proposal as `K8sResourceCommon`, so components read `proposal.metadata.name` directly. However, `request` is
  in `spec` (not metadata) and `K8sResourceCommon` doesn't include `spec`, so `request` must go through the view model to remain accessible without
  a type assertion. `source` is extracted from a label whose key is API-specific, so it also goes through the view model.
- **String enums instead of union types for display values.** `reversibility` is `string`, not `'Reversible' | 'Irreversible' | 'Partial'`. The
  component renders it as-is in a Label. If the API adds a new value, it just appears — no type error, no component change.
- **`AuditEntry[]` instead of named fields.** The post-execution audit trail is a flexible list of `{label, value}` pairs. If the API adds "Git
  commit" or removes "Token burn", the mapping code adds/removes an entry — the component renders whatever it gets.
- **`command` is best-effort.** The API's `ProposedAction` has `type` and `description` but no dedicated command field. The mapping extracts the first
  action's description as the copyable command. This assumption may need adjustment once real API payloads are available — the mapping function is the
  only place to change.

##### K8s models for result CRDs

```typescript
export const AnalysisResultModel: K8sModel = {
  kind: 'AnalysisResult',
  label: 'AnalysisResult',
  labelKey: 'AnalysisResult',
  labelPlural: 'AnalysisResults',
  labelPluralKey: 'AnalysisResults',
  apiGroup: 'agentic.openshift.io',
  apiVersion: 'v1alpha1',
  abbr: 'AR',
  namespaced: true,
  crd: true,
  plural: 'analysisresults',
};

export const ExecutionResultModel: K8sModel = {
  kind: 'ExecutionResult',
  label: 'ExecutionResult',
  labelKey: 'ExecutionResult',
  labelPlural: 'ExecutionResults',
  labelPluralKey: 'ExecutionResults',
  apiGroup: 'agentic.openshift.io',
  apiVersion: 'v1alpha1',
  abbr: 'ER',
  namespaced: true,
  crd: true,
  plural: 'executionresults',
};

export const VerificationResultModel: K8sModel = {
  kind: 'VerificationResult',
  label: 'VerificationResult',
  labelKey: 'VerificationResult',
  labelPlural: 'VerificationResults',
  labelPluralKey: 'VerificationResults',
  apiGroup: 'agentic.openshift.io',
  apiVersion: 'v1alpha1',
  abbr: 'VR',
  namespaced: true,
  crd: true,
  plural: 'verificationresults',
};
```

##### Utilities (`proposal-utils.ts`)

This file contains **component-facing** pure functions that operate on view model types (`ProposalPhase` from `types.ts`). It does NOT import from
`api-types.ts`.

- `getPhaseStatus(phase: ProposalPhase)` — maps phase to PatternFly Label status (`danger`, `warning`, `info`, `success`, `custom`)

Note: `getPhaseLabel()` is NOT needed — display labels (e.g., `EmergencyStopped` → `"Emergency stopped"`) are handled by the i18n translation keys
via `t(phase)`. Adding a separate function would duplicate the translation system.

The `derivePhase()` function (which takes raw `Condition[]` from `api-types.ts`) lives **inside `useProposal.ts`** as a non-exported helper, since it
operates on raw API types and is only called by the hook. This keeps the boundary clean — `proposal-utils.ts` never imports `api-types.ts`.

`derivePhase` logic ported from `proposal_types.go`, priority chain:

1. EmergencyStopped=True → EmergencyStopped
2. Escalated=True → Escalated
3. Denied=True → Denied
4. Escalated exists but not True → Escalating/Failed
5. Verified exists → Completed/Verifying/Executing(retrying)/Failed
6. Executed exists → Verifying/Executing/Failed
7. Analyzed exists → Proposed/Analyzing/Failed
8. No conditions → Pending

##### Proposal API hook (`useProposal.ts`)

This hook is the **single abstraction point** for all K8s API interaction with proposals AND the **mapping boundary** between raw API types and UI
view models. It imports from `api-types.ts`, maps to `types.ts`, and exposes only view models. Components never import `api-types.ts`,
`k8sPatchResource`, `useK8sWatchResource`, or any console SDK K8s utility.

**Boundary rule:** everything above this hook (components) speaks view-model language. Everything below (K8s SDK, API types) is hidden.

```
┌─────────────────────────────────────────┐
│  Components (ProposalDetailsPage, etc.) │  ← import from types.ts only
│  Props: ProposalView, RootCauseView, etc.   │
├─────────────────────────────────────────┤
│  useProposal.ts                         │  ← THE BOUNDARY
│  - imports api-types.ts (raw CRD)       │
│  - imports types.ts (view models)       │
│  - maps ApiProposal → ProposalView        │
│  - maps ApiAnalysisResult → RootCauseView │
│  - maps ApiExecutionResult → PostExecutionView │
│  - wraps K8s SDK calls for mutations    │
├─────────────────────────────────────────┤
│  api-types.ts (raw CRD shapes)          │  ← mirrors v1alpha1, NOBODY else imports
│  K8s SDK (useK8sWatchResource, etc.)    │
└─────────────────────────────────────────┘
```

The hook exposes:

- **View model**: `ProposalView` assembled from all watched CRs
- **Mutations**: Async action callbacks with loading/error state
- **Loading state**: Split between proposal (required) and results (progressive)

```typescript
interface UseProposalReturn {
  proposal: K8sResourceCommon | undefined;  // raw metadata (name, createdAt, etc.) — stable across API versions
  view: ProposalView | undefined;           // mapped data (rootCause, options, postExecution) — API-insulated

  // Loading is split: the page renders as soon as the proposal loads.
  // Result watches load progressively — the RCA and remediation sections show
  // inline spinners while results are still loading, instead of blocking the whole page.
  proposalLoaded: boolean;
  proposalError: Error | undefined;
  resultsLoaded: boolean;                   // true when all 3 result watches have resolved
  resultsError: Error | undefined;          // set if any result watch fails (e.g., RBAC 403)

  stopAnalysis: () => Promise<void>;
  approveExecution: () => Promise<void>;
  mutationInProgress: boolean;
  mutationError: string | undefined;
}

// RESULT_LABEL_PROPOSAL imported from constants.ts (added in this phase)

// Builds the watch config for a result CRD, using server-side filtering.
// Strategy: label selector first (most efficient), field selector as fallback.
function buildResultWatch(gvk: K8sGroupVersionKind, proposalName: string) {
  return {
    groupVersionKind: gvk,
    namespace: PROPOSAL_NAMESPACE,
    isList: true,
    // Primary: label selector — server filters at the API level, only matching CRs are streamed.
    // The operator is expected to set `agentic.openshift.io/proposal: <name>` on each result CR.
    selector: { matchLabels: { [RESULT_LABEL_PROPOSAL]: proposalName } },
  };
}

const useProposal = (name: string): UseProposalReturn => {
  // --- K8s watches (raw API types, private to this hook) ---
  // Proposal: single-resource watch by name — most efficient, no list scan
  const [proposal, proposalLoaded, proposalError] = useK8sWatchResource<ApiProposal>({
    groupVersionKind: { group: 'agentic.openshift.io', version: 'v1alpha1', kind: 'Proposal' },
    name,
    namespace: PROPOSAL_NAMESPACE,
  });

  // Results: list watches with server-side label filtering per proposal
  const analysisGVK = { group: 'agentic.openshift.io', version: 'v1alpha1', kind: 'AnalysisResult' };
  const executionGVK = { group: 'agentic.openshift.io', version: 'v1alpha1', kind: 'ExecutionResult' };
  const verificationGVK = { group: 'agentic.openshift.io', version: 'v1alpha1', kind: 'VerificationResult' };

  const [analysisResults, arLoaded, arError] = useK8sWatchResource<ApiAnalysisResult[]>(
    buildResultWatch(analysisGVK, name),
  );
  const [executionResults, erLoaded, erError] = useK8sWatchResource<ApiExecutionResult[]>(
    buildResultWatch(executionGVK, name),
  );
  const [verificationResults, vrLoaded, vrError] = useK8sWatchResource<ApiVerificationResult[]>(
    buildResultWatch(verificationGVK, name),
  );

  // --- Map API → view models ---
  // filterLatest still needed: the label selector narrows to this proposal's results,
  // but there may be multiple results per step (retries). We want the latest.
  const view = useMemo(() => {
    if (!proposal) return undefined;
    const latestAnalysis = filterLatest(analysisResults);
    const latestExecution = filterLatest(executionResults);
    const latestVerification = filterLatest(verificationResults);
    return mapToProposalView(proposal, latestAnalysis, latestExecution, latestVerification);
  }, [proposal, analysisResults, executionResults, verificationResults]);

  // --- Mutations (hidden K8s SDK calls) ---
  const [mutationInProgress, setMutationInProgress] = useState(false);
  const [mutationError, setMutationError] = useState<string>();

  const stopAnalysis = useCallback(async () => { /* k8sPatchResource... */ }, [proposal]);
  const approveExecution = useCallback(async () => { /* k8sPatchResource... */ }, [proposal]);

  return {
    proposal: proposal as K8sResourceCommon,
    view,
    proposalLoaded,
    proposalError,
    resultsLoaded: arLoaded && erLoaded && vrLoaded,
    resultsError: arError || erError || vrError,
    stopAnalysis, approveExecution, mutationInProgress, mutationError,
  };
};
```

##### Mapping functions (inside `useProposal.ts`, not exported)

These are the **only place** that knows the raw API shape. When the alpha API changes, these functions absorb the diff.

```typescript
function mapToProposalView(
  proposal: ApiProposal,
  analysis?: ApiAnalysisResult,
  execution?: ApiExecutionResult,
  verification?: ApiVerificationResult,
): ProposalView {
  const phase = derivePhase(proposal.status?.conditions ?? []);
  const apiOptions = analysis?.status?.options ?? [];

  // Extract source from proposal labels — label key is API-specific, mapped here so components don't import it
  const source = proposal.metadata?.labels?.[PROPOSAL_LABEL_SOURCE];

  // Surface the most relevant failure reason from whichever step failed
  const failureReason = analysis?.status?.failureReason
    || execution?.status?.failureReason
    || verification?.status?.failureReason;

  return {
    phase,
    request: proposal.spec.request,
    source,
    failureReason,
    rootCause: mapRootCause(apiOptions),
    options: apiOptions.map(mapOption),
    postExecution: phase === 'Completed'
      ? mapPostExecution(apiOptions, execution, verification)
      : undefined,
  };
}

function mapRootCause(options: ApiRemediationOption[]): RootCauseView | undefined {
  // Extract from first option's diagnosis — if API moves this, change HERE only
  const diag = options[0]?.diagnosis;
  if (!diag) return undefined;
  return { cause: diag.rootCause, detail: diag.summary, confidence: diag.confidence };
}

function mapOption(opt: ApiRemediationOption, index: number): RemediationOptionView {
  return {
    index,
    title: opt.title,
    description: opt.summary || opt.proposal?.description || '',
    reversibility: opt.proposal?.reversible,
    risk: opt.proposal?.risk,
    estimatedImpact: opt.proposal?.estimatedImpact,
    command: opt.proposal?.actions?.[0]?.description,
    actions: opt.proposal?.actions,
    rollbackDescription: opt.proposal?.rollbackPlan?.description,
    rollbackCommand: opt.proposal?.rollbackPlan?.command,
  };
}

function mapPostExecution(
  options: ApiRemediationOption[],
  execution?: ApiExecutionResult,
  verification?: ApiVerificationResult,
): PostExecutionView {
  const auditTrail: AuditEntry[] = [];
  // Build audit trail from whatever the API provides — flexible, not hardcoded fields
  const startedCond = execution?.status?.conditions?.find(c => c.type === 'Started');
  const completedCond = execution?.status?.conditions?.find(c => c.type === 'Completed');
  if (startedCond?.lastTransitionTime)
    auditTrail.push({ label: 'Applied', value: startedCond.lastTransitionTime });
  if (completedCond?.lastTransitionTime)
    auditTrail.push({ label: 'System restored', value: completedCond.lastTransitionTime });
  // ... more entries extracted from conditions/actions as available

  return {
    originalRootCause: options[0]?.diagnosis?.summary ?? '',
    remediationDelta: options[0]?.proposal?.description ?? '',
    outcome: execution?.status?.verification?.conditionOutcome ?? '',
    auditTrail,
    verificationSummary: verification?.status?.summary,
  };
}
```

**What changes when the API changes — and what doesn't:**

| API change                                               | Files that change                                                 | Files that DON'T change                                         |
| -------------------------------------------------------- | ----------------------------------------------------------------- | --------------------------------------------------------------- |
| Field rename (e.g. `rootCause` → `cause`)                | `api-types.ts`, `mapRootCause()` in `useProposal.ts`              | All components                                                  |
| Struct restructure (e.g. flatten `diagnosis`/`proposal`) | `api-types.ts`, mapping functions in `useProposal.ts`             | All components                                                  |
| New CRD for approval                                     | `api-types.ts`, mutation callbacks in `useProposal.ts`            | All components                                                  |
| New field to display (e.g. "Git commit")                 | `mapPostExecution()` adds an `AuditEntry`                         | `PostExecutionSummary.tsx` (already renders any `AuditEntry[]`) |
| `v1alpha1` → `v1alpha2` version bump                     | `api-types.ts`, K8s models, GVK in `useProposal.ts`               | All components                                                  |
| New phase value                                          | `api-types.ts`, `derivePhase()`, `types.ts` `ProposalPhase` union | `ProposalPhaseLabel.tsx` (renders unknown phases as `custom`)   |

**Filtering strategy for result watches:**

Result CRs (AnalysisResult, ExecutionResult, VerificationResult) are fetched with server-side filtering to avoid streaming all results across all
proposals. The strategy is **label selector first, field selector fallback**:

1. **Label selector (primary).** The hook uses `selector: { matchLabels: { 'agentic.openshift.io/proposal': name } }`. This requires the operator to
   set this label on each result CR it creates — verify with the operator team. This is the most efficient path: the API server filters at watch
   time, and only matching CRs are streamed to the browser. This matches the existing pattern in `useProposalCheck.ts` which already uses label
   selectors for proposal lookup.

2. **Field selector (fallback).** If the operator does not set a proposal label on result CRs, switch to
   `fieldSelector: 'spec.proposalName=<name>'`. Field selectors work on any spec/status field without operator changes, but K8s only supports
   field selectors on indexed fields by default — `metadata.name` and `metadata.namespace` are always indexed, but `spec.proposalName` requires
   the operator to register a field index. Test against a real cluster to confirm.

3. **Client-side post-filter (last resort).** If neither server-side filter works, remove the `selector`/`fieldSelector` and post-filter in the
   `useMemo` with `results.filter(r => r.spec.proposalName === name)`. This works but streams all results namespace-wide. Acceptable for clusters
   with few proposals, but add a TODO for server-side filtering.

The `filterLatest()` helper is still needed regardless of server-side filtering — a proposal can have multiple AnalysisResults (retries). The label
selector narrows to this proposal; `filterLatest` picks the most recent one by `creationTimestamp`.

**Additional performance notes:**

- **4 concurrent K8s watches.** During the Analyzing phase, ExecutionResult and VerificationResult CRs don't exist yet — the watches return empty
  arrays, which is harmless. Consider conditional watches (only open when phase >= Executing) as a future optimization if API server load becomes a
  concern.

#### Phase 1 Verification

- `cd projects/monitoring-plugin/web && npx tsc --noEmit` — no type errors

### Phase 2: Proposal detail page component

**Dependency:** Phase 1 **Parallel with:** None

#### Files Modified

| File                                                              | Change                                                        |
| ----------------------------------------------------------------- | ------------------------------------------------------------- |
| `web/src/components/ai-proposals/ProposalDetailsPage.tsx` (new)   | Main detail page component orchestrating state-driven layout  |
| `web/src/components/ai-proposals/ProposalPhaseLabel.tsx` (new)    | Reusable phase status Label component                         |
| `web/src/components/ai-proposals/RootCauseAnalysis.tsx` (new)     | RCA section: analyzing skeleton OR detected root cause card   |
| `web/src/components/ai-proposals/RemediationOptionCard.tsx` (new) | Expandable card with radio, reversibility badge, command copy |
| `web/src/components/ai-proposals/PostExecutionSummary.tsx` (new)  | Completed-state card: contextual evidence + audit trail       |

#### Details

##### State-driven UI behavior

The mockup shows **four distinct visual states**. The page layout adapts based on the proposal's derived phase:

**Common header (all states):**

- Breadcrumb: "Alerts > Proposal details"
- "P" ResourceIcon + proposal name as h1
- ProposalPhaseLabel badge (Analyzing / Proposed / Completed / Plan aborted / etc.)
- Source label (e.g., "alertmanager") — from `view.source`, shown as outline Label when present
- "Created {timestamp}" (shown for Proposed state; sometimes omitted for Completed/Aborted)

**State: Analyzing**

- RCA section: progress bar + "Analyzing infrastructure topology to isolate root cause..." text + skeleton placeholders
- "Stop analysis" danger button (right-aligned, below RCA card) with confirmation modal
- Remediation hub section: "Remediation options will be synthesized following root cause confirmation." + skeleton placeholders
- No action button at bottom

**State: Proposed**

- RCA section: Card with "DETECTED ROOT CAUSE" label (with icon) + `rootCause.cause` (bold first line) + `rootCause.detail` (explanation)
- Remediation hub section: "{N} remediation options" badge + "Created {timestamp}"
- Expandable option cards with radio selection:
  - Option header: "Option {N}" + `option.reversibility` badge + `option.title`
  - First option expanded by default, others collapsed
  - Expanded content: `option.description` + "COMMAND EXECUTED BY AGENT" with `option.command` in ClipboardCopy + "Download plan" button
  - Radio button on right side of each option for selection
- "Execute remediation" primary button at bottom with confirmation modal

**State: Completed**

- RCA section: same as Proposed (Card with root cause)
- Remediation hub section: "Completed" success badge (replaces option count)
- PostExecutionSummary card (replaces option cards) — renders `postExecution.auditTrail` dynamically
- No action button at bottom

**State: Executing / Verifying**

- RCA section: same as Proposed (root cause available, analysis is done)
- Remediation hub section: shows all options read-only (no radio, no expand/collapse interaction) with phase badge. The API does not store which
  option the user selected — the operator picks. If the API adds a "selectedOption" field in the future, the mapping can expose it and the UI can
  highlight it. For now, all options are displayed equally in read-only mode.
- No action buttons — execution is in progress, user cannot intervene

**State: Aborted (EmergencyStopped / Failed / Denied)**

- RCA section: Card with root cause (if analysis completed before abort), otherwise empty state
- Remediation hub section: "{N} remediation option" count badge only — no expandable cards, no actions
- No action button at bottom

**State: Escalating / Escalated**

- Same layout as Aborted — show root cause if available, option count badge, no actions
- Phase badge distinguishes from Aborted

**State: Pending**

- Same as Analyzing but without the progress bar — just a "Waiting for analysis to start..." placeholder

##### Page structure (`ProposalDetailsPage.tsx`)

Following the established pattern from `SilencesDetailsPage.tsx`:

```
MonitoringProvider
└── ProposalDetailsPageWithFallback (withFallback wrapper)
    ├── DocumentTitle
    └── StatusBox (data=proposal, loaded=proposalLoaded, loadError=proposalError)
        └── PageGroup
            ├── PageBreadcrumb
            │   └── Breadcrumb: "Alerts" > "Proposal details"
            ├── PageSection (header)
            │   └── Flex
            │       ├── ResourceIcon "P" + Title h1 (proposal.metadata.name)
            │       ├── ProposalPhaseLabel (view.phase)
            │       ├── IF view.source: Label variant="outline" (view.source)
            │       └── "Created <Timestamp>" (proposal.metadata.creationTimestamp)
            ├── IF view.failureReason:
            │   └── Alert variant="danger": view.failureReason
            ├── IF resultsError:
            │   └── Alert variant="warning": "Unable to load analysis results. {resultsError.message}"
            ├── Divider
            ├── PageSection (Root cause analysis)
            │   ├── IF !resultsLoaded: inline Skeleton (not full-page loading)
            │   └── IF resultsLoaded:
            │       └── RootCauseAnalysis (rootCause, phase, failureReason, onStop, mutationInProgress)
            ├── Divider
            └── PageSection (Remediation hub)
                ├── Title h4: "Remediation hub"
                ├── IF !resultsLoaded: inline Skeleton
                ├── IF phase=Pending/Analyzing:
                │   └── placeholder text + skeleton
                ├── IF phase=Proposed AND options.length > 0:
                │   ├── "{view.options.length} remediation options" badge
                │   ├── RemediationOptionCard per option (first expanded)
                │   └── "Execute remediation" primary Button + confirmation modal
                ├── IF phase=Proposed AND options.length === 0:
                │   └── EmptyState: "Analysis completed but no remediation options were generated."
                ├── IF phase=Executing/Verifying:
                │   └── all options read-only (no radio/expand) + phase badge
                ├── IF phase=Completed:
                │   ├── "Completed" success badge
                │   ├── IF view.postExecution has content:
                │   │   └── PostExecutionSummary (view.postExecution)
                │   └── ELSE: "Execution completed. No detailed summary available."
                └── IF phase=Failed/Denied/EmergencyStopped/Escalating/Escalated:
                    └── "{view.options.length} remediation option" count badge only
```

**Loading strategy:** `StatusBox` only waits for the Proposal itself (`proposalLoaded` / `proposalError`). This means the header renders immediately.
The RCA and remediation hub sections show inline `Skeleton` placeholders while `resultsLoaded` is false. If result watches fail (e.g., RBAC 403 on
AnalysisResult), an inline `Alert` appears below the header but the proposal metadata remains visible — the page degrades gracefully instead of
showing a full-page error.

**Phase transition handling:** The `ProposalDetailsPage` component uses `view.phase` as a key signal for which section to render. Local state for
option selection and expansion (`selectedOption`, `expandedOption`) must reset when `view.phase` changes:

```typescript
const [selectedOption, setSelectedOption] = useState(0);
const [expandedOption, setExpandedOption] = useState(0);
const prevPhaseRef = useRef(view?.phase);

useEffect(() => {
  if (prevPhaseRef.current !== view?.phase) {
    setSelectedOption(0);
    setExpandedOption(0);
    prevPhaseRef.current = view?.phase;
  }
}, [view?.phase]);
```

This prevents stale local state when the phase transitions (e.g., user was viewing option 2 expanded in Proposed, phase moves to Executing — the
component resets instead of trying to render option cards in a phase that doesn't show them).

##### Phase status label (`ProposalPhaseLabel.tsx`)

Props: `{ phase: ProposalPhase }` — from `types.ts`. Renders a PatternFly `Label` with the appropriate color/status.

Uses a known-phase → status map with a fallback for unknown phases (future-proofing against new API phases):

```typescript
const ProposalPhaseLabel: FC<{ phase: ProposalPhase }> = ({ phase }) => {
  const { t } = useTranslation(process.env.I18N_NAMESPACE);
  const statusMap: Partial<Record<ProposalPhase, string>> = {
    Pending: 'info',
    Analyzing: 'info',
    Proposed: 'warning',
    Executing: 'info',
    Verifying: 'info',
    Completed: 'success',
    Failed: 'danger',
    Denied: 'danger',
    Escalating: 'warning',
    Escalated: 'warning',
    EmergencyStopped: 'danger',
  };
  return <Label status={statusMap[phase] ?? 'custom'}>{t(phase)}</Label>;
};
```

##### Root cause analysis section (`RootCauseAnalysis.tsx`)

Props: `{ rootCause?: RootCauseView; phase: ProposalPhase; failureReason?: string; onStop: () => void; mutationInProgress: boolean }` — all from
`types.ts`.

Renders differently based on phase and data availability:

```
IF phase is Pending:
  └── "Waiting for analysis to start..." placeholder text

IF phase is Analyzing:
  └── Card (dashed border):
      ├── Progress bar (indeterminate)
      ├── "Analyzing infrastructure topology to isolate root cause..." italic text
      ├── Skeleton placeholders (3 lines)
      └── "Stop analysis" danger Button (right-aligned) → calls onStop with confirmation modal

IF rootCause is present (Proposed, Completed, Executing, Verifying, Failed, etc.):
  └── Card:
      ├── Icon + "DETECTED ROOT CAUSE" (uppercase label)
      ├── rootCause.cause (bold text — the one-line root cause statement)
      └── rootCause.detail (body text — longer explanation)

IF phase is terminal AND rootCause is absent AND failureReason is present:
  └── Alert variant="danger": failureReason
      (e.g., "Analysis failed: sandbox pod timed out after 300s")

IF phase is terminal AND rootCause is absent AND no failureReason:
  └── Empty state: "Root cause analysis was not completed."
```

Key PatternFly components: `Card`, `CardBody`, `Skeleton`, `Button`.

##### Remediation option card (`RemediationOptionCard.tsx`)

Props: `{ option: RemediationOptionView; isExpanded: boolean; isSelected: boolean; onSelect: () => void }` — all from `types.ts`.

Each option is rendered as a bordered Card with an expandable body. The mockup shows:

```
Card (bordered, with radio button on right)
├── Header row:
│   ├── Expand/Collapse chevron button
│   ├── "Option {option.index + 1}" text
│   ├── option.reversibility Label badge (rendered as-is)
│   ├── option.title (bold)
│   └── Radio button (right-aligned, for option selection)
├── IF expanded:
│   ├── option.description text
│   ├── IF option.command:
│   │   ├── "COMMAND EXECUTED BY AGENT" uppercase label
│   │   └── ClipboardCopy with option.command (readonly) + Copy button
│   └── "Download plan" outlined Button with download icon
```

The first option is expanded and its radio is checked by default. Selecting a different radio collapses the current and expands the new one.

**"Download plan" button:** Generates a JSON file containing the remediation option data (`option.title`, `option.description`, `option.actions`,
`option.command`) and triggers a browser download via `URL.createObjectURL` + anchor click. The file name follows the pattern
`proposal-{name}-option-{index}.json`.

Key PatternFly components: `Card`, `CardBody`, `Radio`, `Label`, `ClipboardCopy`, `Button`, `ExpandableSection`.

##### Post-execution summary (`PostExecutionSummary.tsx`)

Props: `{ postExecution: PostExecutionView }` — from `types.ts`. Shown only when `phase === 'Completed'`. Replaces the remediation option cards:

```
Card (bordered, with success checkmark icon in header)
├── Heading: "Post-execution summary"
├── "CONTEXTUAL EVIDENCE" uppercase label
├── DescriptionList:
│   ├── "Original root cause": postExecution.originalRootCause
│   └── "Remediation delta": postExecution.remediationDelta
├── "AUDIT TRAIL" uppercase label
├── DescriptionList (dynamic, iterates postExecution.auditTrail):
│   └── for each AuditEntry: <DescriptionListTerm>{label}</DescriptionListTerm> <DescriptionListDescription>{value}</DescriptionListDescription>
└── "View post-execution logs" outlined Button (if verificationSummary present)
```

The component renders `auditTrail` as a generic list — it has no knowledge of specific field names like "Applied" or "Git commit". The mapping in
`useProposal.ts` builds the list from whatever the API provides. If the API adds a new field, the mapping adds an `AuditEntry` and the component
renders it automatically with no code change.

Key PatternFly components used:

- `ExpandableSection` — collapsible option cards
- `Label` — status badges (risk, confidence, reversibility)
- `ClipboardCopy` — copyable command text (replaces mockup's "COMMAND EXECUTED BY AGENT")
- `DescriptionList` / `DescriptionListGroup` — structured key-value display
- `Card` — wrapping each option for visual grouping
- `Button` — "Execute remediation", "Stop analysis"
- `Spinner` / progress indicator — for Analyzing state

##### Action wiring (confirmation modals)

The mutation callbacks (`stopAnalysis`, `approveExecution`) live in the `useProposal` hook (Phase 1). The page component wires them to buttons with
confirmation modals. The component never touches K8s SDK calls directly.

**Stop analysis:** Visible when `view.phase === 'Analyzing'`. Clicking opens a PatternFly `Modal` confirmation dialog. On confirm, calls
`stopAnalysis()`. Uses `mutationInProgress` to show a spinner and `mutationError` to show an inline alert on failure.

**Execute remediation:** Visible when `view.phase === 'Proposed'`. Same modal pattern — calls `approveExecution()` on confirm.

Both modals follow the `SilenceDropdown` UX pattern from `SilencesUtils.tsx`: confirmation dialog, loading spinner, error alert, auto-dismiss on
success (the K8s watch auto-updates the phase, re-rendering into the new state).

```typescript
const {
  proposal, view,
  proposalLoaded, proposalError, resultsLoaded, resultsError,
  stopAnalysis, approveExecution, mutationInProgress, mutationError,
} = useProposal(name);
```

##### Export pattern

Following the established pattern:

```typescript
const ProposalDetailsPageWithFallback = withFallback(ProposalDetailsPage_);

export const MpCmoProposalDetailsPage = () => (
  <MonitoringProvider monitoringContext={{ plugin: 'monitoring-plugin', prometheus: 'cmo' }}>
    <ProposalDetailsPageWithFallback />
  </MonitoringProvider>
);
```

#### Phase 2 Verification

- `cd projects/monitoring-plugin/web && npx tsc --noEmit` — no type errors
- `make lint-frontend` — no lint errors

### Phase 3: Route registration, URL updates, and navigation wiring

**Dependency:** Phase 2 **Parallel with:** None

#### Files Modified

| File                                                      | Change                                                            |
| --------------------------------------------------------- | ----------------------------------------------------------------- |
| `web/console-extensions.json`                             | Add route for `/monitoring/v2/ai/proposals/:name`                 |
| `web/package.json`                                        | Add `ProposalDetailsPage` to `exposedModules`                     |
| `web/src/components/hooks/usePerspective.tsx`             | Add `getProposalDetailUrl()` function, update `getProposalsUrl()` |
| `web/src/components/alerting/AlertList/AlertTableRow.tsx` | Update navigation to use path param URL                           |
| `web/locales/en/plugin__monitoring-plugin.json`           | Add i18n translation keys                                         |

#### Details

##### Route registration (`console-extensions.json`)

Add after the existing `/monitoring/alerts/:ruleID` route entry:

```json
{
  "type": "console.page/route",
  "properties": {
    "exact": false,
    "path": "/monitoring/v2/ai/proposals/:name",
    "component": {
      "$codeRef": "ProposalDetailsPage.MpCmoProposalDetailsPage"
    }
  }
}
```

##### Exposed module (`package.json`)

Add to `exposedModules`:

```json
"ProposalDetailsPage": "./components/ai-proposals/ProposalDetailsPage"
```

##### URL helpers (`usePerspective.tsx`)

Add a new function for the detail page URL:

```typescript
export const getProposalDetailUrl = (perspective: Perspective, name: string): string => {
  switch (perspective) {
    case 'admin':
      return `/monitoring/v2/ai/proposals/${encodeURIComponent(name)}`;
    default:
      return '';
  }
};
```

Update `getProposalsUrl` to keep it for the list page (future use) without query params for individual proposals.

##### Navigation update (`AlertTableRow.tsx`)

Change the navigation from query-param to path-param URL:

```typescript
// Before:
const proposalUrl = getProposalsUrl(
  perspective,
  new URLSearchParams(
    proposalName ? { name: proposalName } : { fingerprint: alertFingerprint },
  ),
);

// After:
const proposalUrl = proposalName
  ? getProposalDetailUrl(perspective, proposalName)
  : getProposalsUrl(perspective, new URLSearchParams({ fingerprint: alertFingerprint }));
```

When there's exactly one matching proposal (common case), navigate directly to the detail page using the path param. When there are multiple or only a
fingerprint, fall back to the list URL (future list page).

##### Translation keys

Add to `plugin__monitoring-plugin.json`:

```json
"Proposal details": "Proposal details",
"Root cause analysis (RCA)": "Root cause analysis (RCA)",
"Remediation options": "Remediation options",
"{{count}} remediation options": "{{count}} remediation options",
"DETECTED ROOT CAUSE": "DETECTED ROOT CAUSE",
"Analyzing": "Analyzing",
"Proposed": "Proposed",
"Executing": "Executing",
"Verifying": "Verifying",
"Completed": "Completed",
"Failed": "Failed",
"Denied": "Denied",
"Pending": "Pending",
"Escalating": "Escalating",
"Escalated": "Escalated",
"EmergencyStopped": "Emergency stopped",
"Option {{index}}": "Option {{index}}",
"Risk": "Risk",
"Confidence": "Confidence",
"Estimated impact": "Estimated impact",
"Reversibility": "Reversibility",
"Proposed actions": "Proposed actions",
"Rollback plan": "Rollback plan",
"Verification steps": "Verification steps",
"Execute remediation": "Execute remediation",
"Stop analysis": "Stop analysis",
"Root cause": "Root cause",
"Analyzing infrastructure to isolate root cause...": "Analyzing infrastructure to isolate root cause...",
"Remediation options will be synthesized following root cause confirmation.": "Remediation options will be synthesized following root cause confirmation.",
"Waiting for analysis to start...": "Waiting for analysis to start...",
"Root cause analysis was not completed.": "Root cause analysis was not completed.",
"Remediation hub": "Remediation hub",
"Post-execution summary": "Post-execution summary",
"CONTEXTUAL EVIDENCE": "CONTEXTUAL EVIDENCE",
"AUDIT TRAIL": "AUDIT TRAIL",
"Original root cause": "Original root cause",
"Remediation delta": "Remediation delta",
"View post-execution logs": "View post-execution logs",
"Download plan": "Download plan",
"Plan aborted": "Plan aborted",
"Alerts": "Alerts",
"COMMAND EXECUTED BY AGENT": "COMMAND EXECUTED BY AGENT",
"Are you sure you want to stop the analysis?": "Are you sure you want to stop the analysis?",
"Are you sure you want to execute this remediation?": "Are you sure you want to execute this remediation?",
"Analysis completed but no remediation options were generated.": "Analysis completed but no remediation options were generated.",
"Unable to load analysis results.": "Unable to load analysis results.",
"Execution completed. No detailed summary available.": "Execution completed. No detailed summary available."
```

#### Phase 3 Verification

- `make lint-frontend` — no lint errors
- `make test-translations` — all translation keys present
- `npx tsc --noEmit` — no type errors

### Phase 4: Unit tests

**Dependency:** Phase 1 **Parallel with:** Phase 2, Phase 3

#### Files Modified

| File                                                           | Change                                                                                                                          |
| -------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| `web/src/components/ai-proposals/proposal-utils.spec.ts` (new) | Unit tests for `getPhaseStatus()`                                                                                               |
| `web/src/components/ai-proposals/useProposal.spec.ts` (new)    | Unit tests for `derivePhase()` and mapping functions (`mapRootCause`, `mapOption`, `mapPostExecution`) with sample API fixtures |

#### Details

##### `derivePhase` tests (in `useProposal.spec.ts`)

Test with all condition combinations:

```typescript
describe('derivePhase', () => {
  it('returns Pending when no conditions', () => { ... });
  it('returns Analyzing when Analyzed condition is Unknown', () => { ... });
  it('returns Proposed when Analyzed condition is True', () => { ... });
  it('returns Executing when Executed condition is Unknown', () => { ... });
  it('returns Verifying when Executed condition is True', () => { ... });
  it('returns Completed when Verified condition is True', () => { ... });
  it('returns Failed when Analyzed condition is False', () => { ... });
  it('returns Failed when Verified condition is False', () => { ... });
  it('returns Executing on retry when Verified is False with RetryingExecution reason', () => { ... });
  it('returns EmergencyStopped when EmergencyStopped condition is True', () => { ... });
  it('returns Denied when Denied condition is True', () => { ... });
  it('returns Escalating when Escalated condition is Unknown', () => { ... });
  it('returns Escalated when Escalated condition is True', () => { ... });
});
```

Also test `getPhaseStatus()`.

##### Mapping function tests (`useProposal.spec.ts`)

These are the critical API-boundary tests. Use sample API payloads (fixtures) to verify the mapping produces correct view models. When the API
changes, these tests are the first to break — which is the point.

```typescript
describe('mapRootCause', () => {
  it('extracts cause and detail from first option diagnosis', () => { ... });
  it('returns undefined when no options', () => { ... });
  it('returns undefined when first option has no diagnosis', () => { ... });
});

describe('mapOption', () => {
  it('flattens nested proposal fields into flat view', () => { ... });
  it('falls back to summary when proposal.description is absent', () => { ... });
  it('extracts command from first action description', () => { ... });
  it('handles option with no proposal (analysis-only)', () => { ... });
});

describe('mapPostExecution', () => {
  it('builds audit trail from execution conditions', () => { ... });
  it('includes verification summary when available', () => { ... });
  it('produces empty audit trail when execution has no conditions', () => { ... });
});
```

#### Phase 4 Verification

- `cd projects/monitoring-plugin && make test-frontend` — all tests pass

## PR Strategy

| PR | Repository                  | Branch                      | Description                                                                                                       | Dependencies |
| -- | --------------------------- | --------------------------- | ----------------------------------------------------------------------------------------------------------------- | ------------ |
| 1  | openshift/monitoring-plugin | `feat/proposal-detail-page` | Implement proposal detail page with RCA display, remediation options, actions, route registration, and unit tests | None         |

All changes fit in a single PR to the monitoring-plugin repository. The PR should be reviewed as a whole since the route registration, component, and
navigation changes are tightly coupled.

## Verification

End-to-end verification mapped to the spec's acceptance criteria:

- **A detail page exists accessible from "View AI Investigation" action** — Click kebab menu on an alert with a proposal, click "View AI
  Investigation", verify the detail page renders with the proposal name as title and breadcrumb navigation back to Alerts.
- **The detail page reads proposal details from the API using a URL parameter** — Verify the page URL is `/monitoring/v2/ai/proposals/<name>`, the
  `useParams` hook extracts the name, and `useK8sWatchResource` fetches the Proposal and AnalysisResult from `openshift-lightspeed` namespace.
- **The detail page sticks to the mockup while aligning with the API** — Verify: breadcrumb ("Alerts > Proposal details"), proposal name title with P
  icon, phase status label, created timestamp, RCA section with root cause and summary, remediation options as expandable cards with
  risk/reversibility badges, action buttons.
- **State-driven rendering** — Verify the page renders correctly for all four states: Analyzing (skeleton + stop button), Proposed (RCA card + option
  cards + execute button), Completed (RCA card + post-execution summary), Aborted (RCA card + minimal remediation hub).
- **Only PatternFly 6 components used** — Verify no custom CSS files are created. Grep for `.css` or `styled` imports in the new files. All styling
  comes from PatternFly component props.

## Risks

| Risk                                                                                                       | Impact                                                                  | Mitigation                                                                                                                                                                            |
| ---------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| ProposalApproval CRD mechanism is not fully defined                                                        | Execute button may not work end-to-end                                  | Use k8sPatchResource with annotation as interim approach; update when operator API stabilizes                                                                                         |
| RBAC: user can read Proposals but not result CRDs                                                          | Result sections show empty/error while header renders                   | Split loading: `StatusBox` gates on proposal only; result errors show as inline `Alert` below header, not full-page block. User still sees proposal name, phase, and timestamp.        |
| The AnalysisResult plural name `analysisresults` may differ from the actual CRD                            | API calls fail with 404                                                 | Verify against the operator's CRD definition; adjust `AnalysisResultModel.plural` if needed                                                                                           |
| Phase derivation logic diverges from operator's Go implementation                                          | Phase badges show wrong status                                          | Port the Go logic faithfully with comprehensive unit tests covering all branches                                                                                                      |
| Route `/monitoring/v2/ai/proposals/:name` conflicts with future list page at `/monitoring/v2/ai/proposals` | Navigation confusion                                                    | Register the detail route first (more specific path); list page can be added later at `/monitoring/v2/ai/proposals` without conflict since Console matches more-specific routes first |
| `command` field mapping assumes `ProposedAction.description` contains a shell command                      | ClipboardCopy shows natural-language text instead of a runnable command | Mapping is best-effort; verify against real API payloads once available. `mapOption` is the only place to adjust.                                                                     |
| Operator may not set `agentic.openshift.io/proposal` label on result CRs                                  | Label selector returns empty results; detail page appears empty          | Verify label exists on real cluster. If absent, fall back to `fieldSelector: 'spec.proposalName=<name>'`. If field selector also fails (not indexed), fall back to client-side post-filter with a TODO. |
| Proposal CRD not installed (operator not deployed)                                                         | Watch fails with opaque API group 404, not a clear "not found" message  | `StatusBox` shows the error generically. Consider detecting "resource not found" vs "CRD not found" errors in the hook and surfacing a user-friendly message: "The AI investigation feature requires the Lightspeed Agentic Operator." |
| Label mismatch indistinguishable from "results not yet created"                                            | User can't tell if results are still pending or if the filter is broken | Both states show the same Pending/Analyzing skeleton. No fix for MVP — once a real cluster is available, verify that the label selector returns results. If it doesn't, the filtering fallback chain kicks in. |
| API doesn't store which remediation option the user selected                                               | Executing/Verifying state can't highlight the chosen option             | Show all options read-only for now. If the API adds a `selectedOption` field, add it to `ProposalView` and highlight in the UI. |
