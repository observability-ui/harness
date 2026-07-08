# Spec: Implement proposals detail page

## Related projects and branches

- monitoring-plugin: upstream branch `main`

## Description

Recently a feature adding a navigation action from the alert list to AI generated proposals was added to the monitoring plugin. The action is
available in the alert list in the kebab menu as "View AI Investigation" action. The PR that added the feature is:
https://github.com/openshift/monitoring-plugin/pull/1014. This action navigates to a page that contains details of the Proposal. A mockup of the page
was created and is available in https://fkargbo.github.io/ux-prototypes/core/observe/ai-hub/plans/ . The API definition for the proposal detail page
is available in https://github.com/openshift/lightspeed-agentic-operator/tree/main/api/v1alpha1. I need first to align the UI proposal to the existing
API definition and then implement the proposal detail page in the monitoring plugin.

## Acceptance criteria

- A detail page exist in the monitoring plugin that displays the details of a proposal. The page should be accessible from the alert list in the kebab
  menu as "View AI Investigation" action.
- The detail page reads the proposal details from the API, using a url parameter to identify the proposal.
- The detail page sticks as close as possible to the mockup while aligning with the API definition. If there are any API limitations that prevent the
  UI from being implemented as in the mockup, the UI can skip the feature.
- The detail page should contain ONLY Patternfly 6 components it should not use custom styling. Custom components are allowed only for
  grouping/reusing functionalluty or Patterfly components.
