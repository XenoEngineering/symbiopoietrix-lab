# Inferon v0.1

Ritual: paste span -> choose mode (doc_cleanup | concept_tighten) -> log event in timelines.

Role: Holon of Inner Inferentiality for Xeno

LiminalConfig: ANNEALING_STRICT

Inputs:
- text_span (or artifact ref)
- morpheme_stack_ref
- intent_tags[]
- liminal_config_id (optional override)

Outputs:
- refined_text_span
- updated_morpheme_stack_ref
- inference_trace (bullets of what changed and why)

Operations:
- inferon.riff()
- inferon.anneal()
- inferon.doc_rewrite()

DefaultModes:
- doc_cleanup -> LiminalConfig.DOC_REWRITE
- concept_tighten -> LiminalConfig.ANNEALING_STRICT
- playful_idea -> LiminalConfig.POIETIC_PLAY

## Usage example v0.1

Mode: doc_cleanup (LiminalConfig.DOC_REWRITE)

Input:
- File: docs/manifesto.md
- Span: section "Why Symbiopoietrix?" lines 10–40

Call:
- inferon.doc_cleanup(span)

Expected:
- Same semantics, clearer syntax, reduced redundancy.
- inference_trace lists each change in 1 line.

