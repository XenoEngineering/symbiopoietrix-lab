# Phoenix v0.1+

**Ritual:** Read logged timelines → filter by key + layer → weave emergent narrative → log new timeline with provenance.

**Role:** Inter-timeline emergence engine for symbiopoietrix-lab. Reads event logs from sessions, applies LiminalConfig modes, and synthesizes coherent narrative spans that honor source events while enabling new juxtapositions.

## Interface

```
phoenix.emerge(
  source_timelines: Timeline[],
  key_filter: KeyFilter,
  time_windows: TimeWindow[],
  mode: EmergenceMode,
  constraints?: EmergenceConstraints
) -> EmergenceResult

KeyFilter {
  families?:     string[]      // e.g. ["psiPhi", "symbiopoietrix"]
  agents?:       string[]      // e.g. ["Inferon", "Xeno"]
  layers?:       LayerID[]     // e.g. ["semantics", "intent"]
  tags?:         string[]      // freeform; AND-joined with above
  ops?:          string[]      // e.g. ["concept_tighten", "doc_cleanup"]
}

TimeWindow {
  start:         string        // ISO-8601
  end:           string        // ISO-8601
}

EmergenceMode = ENUM {
  EMERGENCE_PLAY,              // Light weaving, poetic juxtaposition
  EMERGENCE_REFINE,            // Tighten spans, compress redundancy
  EMERGENCE_VISION             // Maximal recombination, cross-layer synthesis
}

EmergenceConstraints {
  max_span_count?:       int   // cap on output spans (default: unlimited)
  coherence_floor?:      float // 0.0–1.0 semantic coherence threshold (default: 0.7)
  layer_restrictions?:   LayerID[]  // e.g., ["artifact", "semantics"]
  safety_mode?:          bool  // block "edge" spans if true
}

EmergenceResult {
  emergent_timeline_id:  string        // e.g. "psiPhi-2026-02-24T0600Z-refined-01"
  emergent_timeline:     Timeline      // new timeline with woven events
  span_sequence:         EmergentSpan[]
  provenance_trace:      Provenance[]
  emergence_notes:       string        // brief narrative of weaving logic
}

EmergentSpan {
  id:                    string        // e.g. "emg-001"
  ts:                    string        // synthetic timestamp
  layer:                 LayerID
  family:                string
  agent:                 string        // usually "Phoenix"
  op:                    string        // "weave_span"
  liminal_id:            string        // copied from mode
  tags:                  string[]      // original + ["emergent"]
  span_ref:              SpanRef       // the woven text
  source_event_ids:      string[]      // evt-ids that contributed
}

Provenance {
  emergent_span_id:      string
  source_event_ids:      string[]
  method:                string        // "literal_paste" | "inference_blend" | "layer_bridge"
  coherence_score:       float         // 0.0–1.0
  reasoning:             string        // terse: why this span matters in emergent narrative
}
```

## Modes and LiminalConfig mapping

| Mode | LiminalConfig | Behavior |
|------|---------------|----------|
| EMERGENCE_PLAY | POIETIC_PLAY | Light cross-layer exploration; tolerate loose associations; prioritize aesthetic coherence over strict syntax |
| EMERGENCE_REFINE | ANNEALING_STRICT | Tighten spans; compress redundancy; fuse parallel concepts from same event family; output is condensed, ready-to-use text |
| EMERGENCE_VISION | VISION_QUEST | Maximal synthesis; freely recombine across families and time windows; allow synthetic bridges; highest poietic freedom |

## Algorithm sketch (v0.1)

1. **Filter & Collect**
   - Load source_timelines
   - Apply key_filter to select events
   - Partition events by time_windows

2. **Layer-wise Weaving**
   - For each layer in [artifact, syntax, semantics, intent, meta]:
     - Gather all matching events in that layer
     - Depending on mode:
       - PLAY: random juxtaposition with coherence check
       - REFINE: sort by ts, fuse overlaps, compress
       - VISION: free recombination across time windows, allow cross-family bridges

3. **Span Synthesis**
   - For each woven cluster, create EmergentSpan with:
     - Fresh span_ref text (either literal quote or synthesized blend)
     - source_event_ids pointing back to originals
     - Provenance entry explaining method + coherence score

4. **Validation & Output**
   - Filter by coherence_floor and layer_restrictions
   - Respect max_span_count
   - Assemble EmergenceResult with new Timeline
   - Log to timelines/

## Usage example v0.1+ with real session data

**Scenario:** Weave psiPhi family events from 2026-02-24T0600Z session in REFINE mode.

```
source_timelines: [ "timelines/examples/session-2026-02-24T0600Z.json" ]

key_filter: {
  families: ["psiPhi"],
  layers: ["artifact", "semantics"],
  tags: ["manifesto", "inferon"]
}

time_windows: [
  { start: "2026-02-24T06:00:00Z", end: "2026-02-24T06:05:00Z" }
]

mode: EMERGENCE_REFINE

constraints: {
  coherence_floor: 0.8,
  layer_restrictions: ["artifact", "semantics"]
}
```

**Expected output:**
- emergent_timeline_id: `psiPhi-2026-02-24T0600Z-refined-001`
- span_sequence drawn from:
  - evt-002: Manifesto span about Holons
  - evt-003: Inferon's tightened temporal morpheme line
  - *Synthesis:* New artifact span bridging Holon + morpheme concepts
- provenance_trace showing coherence scores and source event IDs
- emergence_notes: "Refined psiPhi narrative by fusing holonic + temporal semantics; compressed 3 source spans into 2 tighter woven spans maintaining intent."

## Integration points

1. **Inferon**: Phoenix can invoke `inferon.anneal()` on woven spans for final polish
2. **Orchestrator**: Routes incoming weave requests; suggests key_filters based on mission context
3. **LiminalConfig**: All modes reference LC blocks; can be overridden per call
4. **Timelines store**: Reads and appends to timelines/

## Implementation notes

- Phoenix **does not** run live inference; it operates on already-logged events.
- All spans remain **traceable**: provenance_trace makes every synthetic output auditable.
- **Coherence scoring** is domain-specific; use semantic embedding similarity for now.
- Start with EMERGENCE_REFINE (safest) and migrate to VISION as confidence grows.
