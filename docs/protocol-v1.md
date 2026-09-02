# Evidence-resolution projection protocol v1

`.gooo/released.gooo` is the only semantic authority. The runtime lowers it to
semantic IR, evaluates the declared decision precedence, and generates every
machine artifact and human dossier from that IR.

The canonical graph has stable IDs for the graph, authority boundary,
activities, cases, claims, evidence, counterexamples, unknown frontiers, and
provenance. Typed edges retain activity binding, claim support, refutation,
causal-frontier, evidence, and provenance relationships. A node or edge is
never recreated by a reader projection.

The append-only v1 extension declares four roles: `USER`, `REVIEWER`,
`LANGUAGE_MAINTAINER`, and `AUDITOR`. Each role declares a requested
resolution, mandatory fields, hidden node kinds, folded edge kinds, and lost
fields. Its loss manifest exposes integer `visible_nodes`, `hidden_nodes`,
`folded_edges`, and `lost_fields`, the exact IDs behind those integers, and the
canonical graph SHA-256 digest.

Only declared context can be lost. Immutable identity, authority boundary,
decision, evidence, proof cells, all six UNKNOWN frontier fields, and every
refuting counterexample remain present at every role. Conformance compares all
roles to the canonical evaluation and checks that the causal frontier is not
lost, including for a `REFUTED` case that also contains an unresolved claim.

Expansion is evaluated as an authority claim. Missing original graph, role, or
loss policy is `UNKNOWN` with the six-field frontier. Treating restored hidden
facts as canonical is `REFUTED` and points to the source counterexample.

Proof cells carry a proof choice (`FOUNDATION`, `COHERENCE`, or `REGRESSION`)
and an indicator (`DRIVER`, `OUTCOME`, or `GUARDRAIL`) per cell. The contract
does not emit aggregate scores or percentages.
