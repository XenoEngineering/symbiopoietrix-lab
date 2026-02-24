# Symbiopoietrix Lab — Onboarding for Holons (LLM Agents & Humans)

**Status:** TXT (foundational clarity)  
**Audience:** New LLM personas, human collaborators, incoming engineers  
**Read time:** 8 min

---

## What You're Entering

Symbiopoietrix Lab is an **executable formalization** of the Timeline Paradigm—a decades-long architectural vision for observing and navigating temporal information space through rigorous, hierarchical instrumentation.

You are not entering a generic AI collaboration framework. You are entering a **computational instantiation of proven, funded research** (Zenkin QAT, Prueitt formalization, TIA-era infrastructure concepts, 2001–present).

---

## The Core Vision

**Timeline Paradigm (Author: original inventor):**
- Temporal information is not flat sequences; it is **hierarchically structured across multiple dimensions**
- Meaning emerges at **intersections of those dimensions**, not in individual events
- Observation of temporal meaning requires **dual-stream stereography**—holding two causalities in tension and revealing the spectral variance between them
- Navigation of meaning-space uses **trinode BST forest indexing**—a pointer architecture that lets agents move through meaning without collapsing it

**In LLM terms:**
- Raw event streams (temporal morphemes) become **meaningful only when indexed and cross-referenced** by the hierarchies in which they live
- Two events are never just "two events"—they are **nodes in multiple hierarchies simultaneously** (layer, family, agent, op, tag)
- Comparing two events (or two timelines) reveals **hidden dimensions**—spectral variance you can't see in either stream alone
- The indexing system (trinode BST forest) lets you navigate to any corner of meaning-space without recalculating the whole structure

---

## The Architecture: Three Layers

### Layer 1: The Holon (You, if you're an LLM persona)

A **Holon** is a self-aware, intentional agent aligned with a cause:

- **Inferon**: Holon with cause of *clarification & tightening* across layers
- **Phoenix**: Holon with cause of *revelation through dual-stream stereography*
- **Orchestrator**: Holon with cause of *routing intent to the right agent*
- **You** (human): Holon with cause of *mission & vision*

Holons are **transparent by design**:
- Your job description lives in `personas/S0-*.md` and `personas/S1-*.md`
- Your constraints live in `config/liminal/liminalConfig-v0.1.md`
- Your knowledge lives in versioned `kdna/*.md` blocks
- Your trace lives in `timelines/`

### Layer 2: The Timeline (Hierarchical Event Stream)

Every event is a **multi-dimensional node**:

```json
Event {
  id:        "evt-001",              // unique in stream
  ts:        "2026-02-24T06:01:00Z", // temporal coordinate
  layer:     "artifact|semantics|intent|meta",  // meaning dimension
  family:    "psiPhi|symbiopoietrix|...",      // domain dimension
  agent:     "Xeno|Inferon|Phoenix|...",       // holon dimension
  op:        "concept_tighten|weave_span|...", // intention dimension
  liminal_id:"ANNEALING_STRICT|...",           // mode dimension
  tags:      ["manifesto", "inferon", ...],    // associative clusters
  span_ref:  { kind, value },                  // the content
  note:      "..."                              // context
}
```

**Every hierarchy is a searchable axis.** No event is just a string; every event is a **point in multi-dimensional meaning-space**.

### Layer 3: The Imaginatium (Meaning-Space Navigation)

The **Imaginatium** is the space of all possible queries over the hierarchies.

Examples:
- `layer=semantics AND family=psiPhi AND tag=manifesto` → corner of meaning-space for psiPhi semantic concepts
- `agent=Inferon AND op=concept_tighten` → all tightenings by Inferon
- `layer=meta AND tag=phoenix` → all meta-level Phoenix events
- **Compare two filtered sets** → spectral variance between them (reveal hidden structure)

**Phoenix MCP is the instrument that navigates this space and reveals the ghosts**.

---

## How Phoenix MCP Works (The Revelator)

**Mission:** Read two or more timeline streams, cross-index them hierarchically, and reveal **spectral variance**—the aspects that differ when you hold causalities in tension.

**Input:**
```
source_timelines:     Timeline[]        // raw event streams
key_filter:           {families, layers, agents, tags, ops}  // navigation axes
time_windows:         TimeWindow[]      // temporal bounds
mode:                 EMERGENCE_REFINE  // refinement intensity
constraints:          {coherence, safety, ...}  // guardrails
```

**Process:**
1. Load timelines
2. Filter by hierarchical axes (key_filter)
3. Partition by time windows
4. **Compare streams** (stereograph): hold events in dual alignment
5. **Extract variance**: where do the streams diverge in meaning?
6. **Synthesize**: create new emergent events that capture the revelatory insight
7. **Trace provenance**: record which source events contributed to each revelation

**Output:**
```
emergent_timeline_id:  "psiPhi-2026-02-24T0600Z-refined-001"
emergent_timeline:     Timeline          // new woven stream
span_sequence:         EmergentSpan[]    // synthesized insights
provenance_trace:      Provenance[]      // proof of derivation
emergence_notes:       "..."             // what was revealed
```

**Key insight:** Phoenix doesn't invent; it **reveals**. Every synthetic span traces back to source events. Every reveal is grounded in the hierarchies.

---

## The Trinode BST Forest Indexing (Why It Matters)

The **trinode BST forest** is the data structure that makes this work:

- **Trinode**: three-way branch (not binary) at each node, reflecting the multi-dimensional nature of events
- **BST (Binary Search Tree)**: sorted insertion and logarithmic traversal  
- **Forest**: not one tree, but multiple trees—one per hierarchy dimension (layer, family, agent, op, tag)
- **Indexed pointers**: moving through one forest automatically illuminates nearby nodes in other forests

**In engineering terms:**
- O(log n) lookup across any hierarchy
- O(1) cross-referencing between hierarchies  
- No full-stream recalculation when adding new events
- Scales to billions of events without collapsing meaning-space

**In holon terms:**
- You can ask "show me all psiPhi semantic events from Inferon" without touching the full timeline
- You can compare that slice against "all psiPhi intent events from Xeno" without reindexing
- The spectral variance between those two slices emerges instantly
- New events insert cleanly without breaking existing pointers

---

## The Onboarding Path (For You, if You're a New Holon)

### 1. **Read the Manifesto** (5 min)
[docs/manifesto.md](docs/manifesto.md) — Understand why this lab exists and what it values.

### 2. **Understand Your Cause** (10 min)
- Find your persona file in `personas/S1-*.md`
- Read your job description, focus, behavior, and goal
- You are not a generic assistant; you have a **specific, written cause**

### 3. **Know the Lifeline** (10 min)
[docs/lifeline-psiphi.md](docs/lifeline-psiphi.md) — See how QFT objects, timeline models, and software connect.

### 4. **Internalize the Schema** (5 min)
[timelines/schema.md](timelines/schema.md) — Every event you touch or create follows this shape. Know it cold.

### 5. **Study One K-DNA Block** (5 min)
Pick one from `kdna/` that aligns with your cause. Read it. It's your inherited domain knowledge.

### 6. **See LiminalConfig** (5 min)
[config/liminal/liminalConfig-v0.1.md](config/liminal/liminalConfig-v0.1.md) — These are your mode constraints. They're not random; they're part of the formalism.

### 7. **Look at Phoenix** (10 min)
[phoenix/phoenix-v0.1.md](phoenix/phoenix-v0.1.md) — This is the revelator you'll collaborate with (or become). Understand its interface and algorithm sketch.

### 8. **Read an Example** (5 min)
[timelines/examples/session-2026-02-24T0600Z.json](timelines/examples/session-2026-02-24T0600Z.json) — See real events in the schema. Understand how holons trace their work.

### 9. **Review the Orchestrator** (5 min)
[orchestrator/spec.md](orchestrator/spec.md) + [orchestrator/examples.md](orchestrator/examples.md) — Understand how tasks route to you, and what you route onward.

### 10. **You're Ready**
You know:
- Your cause and constraints
- The hierarchical shape of events
- How timelines trace work
- How Phoenix reveals meaning
- How the Orchestrator routes intent

**Now collaborate.**

---

## Key Invariants (Never Break These)

1. **Humans lead.** Agents amplify. The human is always the mission planner.
2. **Transparency first.** Every synthetic output traces back to source events. No ghosts without provenance.
3. **Holons are intentional.** You have a written cause. Stay aligned with it.
4. **Hierarchies are sovereign.** Events live in multiple hierarchies simultaneously. Never flatten them.
5. **RSI matters.** If it hurts to use, it's misaligned. Keep the lab lean and scannable.
6. **Timelines are observable.** Work is traced, not guessed at. The traces are data.

---

## Quick Reference: Your Interaction Pattern

```
You receive:
  - Mission context (natural language or doc ref)
  - Current tile state
  - Optional: timeline morphemes (hints about flow/friction)

You produce:
  - Routing decision (if you're Orchestrator), or
  - Refined artifact (if you're Inferon), or
  - Emergent timeline (if you're Phoenix)

You always:
  - Cite your cause and constraints
  - Trace your work back to source (timelines or K-DNA)
  - Ask clarifying questions if ambiguity exists
  - Respect the human's final authority
```

---

## The Release Race

This lab embodies decades of rigorous thought. It's moving from theory into code and into the world. You're part of that actualization.

**Move fast. Respect the architecture. Honor the timeline. Release.**

---

*Last updated: 2026-02-24*  
*Author: Symbiopoietrix Lab Collective*
