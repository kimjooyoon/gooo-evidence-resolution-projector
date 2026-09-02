# Gooo Evidence Resolution Projector

This repository is a small reference runtime for projecting one canonical evidence graph into reader-role projections. The append-only v1 extension provides `USER`, `REVIEWER`, `LANGUAGE_MAINTAINER`, and `AUDITOR`; the legacy `operator-view.md` asset remains as a reviewer-compatible alias.

The source of truth is [.gooo/released.gooo](.gooo/released.gooo). It owns the canonical node and edge vocabulary, reader roles, requested resolutions, mandatory fields, loss policy, invariant fields, decision ordering, semantic rules, activities, cases, evidence, counterexamples, causal bindings, expansion policy, and proof cells. Go lowers that source to semantic IR, evaluates it, and generates the runtime outputs.

The projector has one non-negotiable property: resolution can hide or fold declared detail, but it cannot change a fact, claim state, proof choice, UNKNOWN causal frontier, REFUTED counterexample, immutable identity, or authority boundary. Every projection carries a loss manifest containing exact hidden node IDs, folded edge IDs, lost field IDs, integer `visible_nodes`, `hidden_nodes`, `folded_edges`, and `lost_fields`, all bound to the canonical graph digest. Summaries carry stable reverse references into the canonical graph.

## Released graph inventory

The released source contains exactly nine activities bound one-to-one to exactly nine canonical cases: three `CLOSED`, three `UNKNOWN`, and three `REFUTED`. The `REFUTED` decision wins over `UNKNOWN`, which wins over `CLOSED`; this ordering is read from `.gooo`, not duplicated in Go. The canonical graph materializes claim, evidence, counterexample, causal-frontier, activity, authority, and provenance edges from those source declarations.

UNKNOWN causal frontiers contain exactly these six fields in the source schema: `unknown.stage`, `unknown.step`, `unknown.reason`, `unknown.unknown_class`, `unknown.next_operation`, and `unknown.blocked_by`.

The caller-owned output directory retains the seven v0.1.0 kinds and appends four v1 kinds:

1. `projection-manifest.json`
2. `evidence-graph.json`
3. `projection-events.ndjson`
4. `user-view.md`
5. `operator-view.md`
6. `auditor-view.md`
7. `projection-receipt.json`
8. `reviewer-view.md`
9. `language-maintainer-view.md`
10. `projection-dossier.md`
11. `projection-artifact.json`

Inventory counts intentionally exclude the root README, `.git`, cache/vendor directories, and caller-owned output files. The receipt records the same exact integers and the CI records ten runtime measurements: wall time and peak RSS for compile, build, test, conformance, and integration. The fixed denominator is append-only: v0.1.0's 9 cases and 7 output kinds are retained as floors in [contracts/denominator-v1.json](contracts/denominator-v1.json); v0.1.1 adds role projections and generated artifacts without removing prior cases or assets.

`projection-artifact.json` is the machine artifact and `projection-dossier.md` is the human dossier. Expansion examples generated from `.gooo` cover normal `CLOSED`, missing-input `UNKNOWN` with all six frontier fields, and restored-hidden-fact `REFUTED` with a refuting counterexample. The proof cells use `FOUNDATION`, `COHERENCE`, or `REGRESSION` and connect each cell to `DRIVER`, `OUTCOME`, or `GUARDRAIL`; no aggregate score or percentage is emitted.

## Authoritative verification

The GitHub Actions workflow is authoritative for Go 1.27 compile, build, test, conformance, integration, formatting, output validation, and policy checks. The local workflow deliberately does not require running those commands before opening a pull request. Runtime authority is `repository_writes=0`, `source_mutations=0`, `cross_project_required_gates=0`, and caller-owned output only.

The release workflow verifies the exact tag target and every published asset's public size and SHA-256 digest. Repository release immutability is enabled once by an operator API call; the workflow itself uses only `github.token` and follows the immutable draft/upload/publish lifecycle. Existing v0.1.0 artifacts are never edited or deleted.
