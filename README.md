# Gooo Evidence Resolution Projector

This repository is a small reference runtime for projecting one canonical evidence graph at three reader resolutions: `user`, `operator`, and `auditor`.

The source of truth is [.gooo/released.gooo](/Users/alice/Documents/Codex/2026-09-01/gooo-evidence-resolution-projector/.gooo/released.gooo). It owns the projection schema, allowed omissions, invariant fields, decision ordering, semantic rules, activities, cases, evidence, and causal bindings. Go only evaluates that source and generates the runtime outputs.

The projector has one non-negotiable property: resolution can remove declared context, but it cannot change a fact, claim state, proof choice, UNKNOWN causal frontier, or REFUTED counterexample. Every projection carries a loss declaration containing exact omitted field IDs and a count. Summaries carry stable reverse references into the canonical graph.

## Released graph inventory

The released source contains exactly nine activities bound one-to-one to exactly nine canonical cases: three `CLOSED`, three `UNKNOWN`, and three `REFUTED`. The `REFUTED` decision wins over `UNKNOWN`, which wins over `CLOSED`; this ordering is read from `.gooo`, not duplicated in Go.

UNKNOWN causal frontiers contain exactly these six fields in the source schema: `unknown.stage`, `unknown.step`, `unknown.reason`, `unknown.unknown_class`, `unknown.next_operation`, and `unknown.blocked_by`.

The caller-owned output directory contains exactly these seven kinds:

1. `projection-manifest.json`
2. `evidence-graph.json`
3. `projection-events.ndjson`
4. `user-view.md`
5. `operator-view.md`
6. `auditor-view.md`
7. `projection-receipt.json`

Inventory counts intentionally exclude the root README, `.git`, cache/vendor directories, and caller-owned output files. The receipt records the same exact integers and the CI records ten runtime measurements: wall time and peak RSS for compile, build, test, conformance, and integration.

## Authoritative verification

The GitHub Actions workflow is authoritative for Go 1.27 compile, build, test, conformance, integration, formatting, output validation, and policy checks. The local workflow deliberately does not require running those commands before opening a pull request.

The release workflow verifies the exact tag target and every published asset's public size and SHA-256 digest. Repository release immutability is enabled once by an operator API call; the workflow itself uses only `github.token` and follows the immutable draft/upload/publish lifecycle.

