# MindChat Prompts v0.1

This file collects ready-made prompt templates for using Groq/MindChat with the Symbiopoietrix Lab personas and K-DNA blocks.

## 1. Plexi-QFT — psiPhi summary (U₁)

System (S):
- Use `personas/S0-mindspeak-core.md` + `personas/S1-plexi-qft.md` concatenated as the **system** message in MindChat.

User (U₁):

```text
TASK:
In this turn, your concrete task is:
- Draft a 6–8 line summary of the psiPhi framework,
- Written for an expert engineer with RSI,
- Avoid marketing or productivity language.

CONTEXT (short):
- psiPhi is a QFT-flavored, category-tinged lifeline framework.
- It links: QFT objects, timeline/Markov models (UPTs, morphemes), and software artifacts (code, tools, data structures).

OUTPUT:
- Step 1: 1-sentence restatement of the task.
- Step 2: The 6–8 line summary (compact, technical, easy to scan).
- Step 3: One focused question if any part of psiPhi’s role or structure is ambiguous.
