# Symbiopoietrix Lab — Phoenix Startup Test Setup v0.1

Purpose
- Provide a repeatable startup dialogue to probe model variance around Phoenix.
- Keep all behavior anchored to existing canon (personas, K‑DNA, timelines, LiminalConfig).

## Canon pack for this test

Include these files as read‑only reference:

- Symbiopoietrix Lab (root)
- docs/overview.md
- docs/manifesto.md
- docs/lifeline-psiphi.md
- config/liminal/liminalConfig-v0.1.md
- inferon/inferon-v0.1.md
- phoenix/phoenix-v0.1.md
- personas/S0-mindspeak-core-v0.1.md
- personas/S1-plexi-qft-v0.1.md
- orchestrator/orchestrator-spec-v0.1.md
- orchestrator/examples-v0.1.md
- kdna/kdna-psiphi-v0.1.md
- kdna/kdna-upt-operator-v0.1.md
- kdna/kdna-forge-style-v0.1.md
- timelines/schema-v0.1.md
- timelines/examples/session-2026-02-24T0600Z.json

## Evaluation knobs

For each run, log:

- Model ID / mode (if exposed by the UI).
- Date/time and provider.
- Any extra system prompt text you add.

You will compare:
- Span selection: which events Phoenix proposes to use.
- Emergent timeline shape: how it names and orders the new spans.
- Explanations: how it justifies choices against canon.

## Run protocol

1) Paste this file + the canon pack into the new model’s context.
2) Paste the entire `tests/startup-psiPhi-Phoenix-v0.1.md` below as the **target behavior**.
3) Ask the model to:
   - Play the parts of Orchestrator, Plexi‑QFT, and Phoenix,
   - Follow the scripted dialogue structure,
   - Fill in only the “MODEL TODO” sections with its own content.

You will then inspect:
- Whether it respects personas and K‑DNA selection.
- Whether Phoenix’s behavior matches your expectations on span choice and emergent timeline output.
