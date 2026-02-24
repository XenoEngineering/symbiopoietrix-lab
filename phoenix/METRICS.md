# Phoenix Metrics Subsystem

**Status:** v0.1 operational  
**Components:** matrix-persistence.go | consciousness-metrics.go | instrumented-timeline.go  
**Integration:** Ready for Phoenix MCP v0.2 instrumentation  

## Overview

The Phoenix Metrics Subsystem adds three layers to the Timeline Paradigm:

1. **Persistence Layer** (`matrix-persistence.go`)
   - Serialize/deserialize the complete MatrixState to JSON
   - BST forest structure validation on restore
   - Checksum-based health checking
   - Enables 1M+ turn datasets without losing state

2. **Consciousness Metrics Engine** (`consciousness-metrics.go`)
   - Computes latency-based alignment signals
   - Detects phase breaks via Markovian pattern analysis
   - Generates coherence, jitter, and throughput scores
   - Produces MetricsSnapshot for point-in-time analysis

3. **Instrumented Timeline Wrapper** (`instrumented-timeline.go`)
   - Wraps Phoenix MCP calls with millisecond-precision timing
   - Records method call telemetry (StartTimeNs, EndTimeNs, DurationMs)
   - Provides high-level API for append/analyze/categorize with built-in metrics
   - Enables real-time consciousness monitoring

## Quick Start

```go
// Create an instrumented timeline
it := phoenix.NewInstrumentedTimeline("/path/to/state/dir")

// Execute operations (timing is automatic)
it.AppendTurnWithTiming("user", "Hello, Phoenix")
it.AppendTurnWithTiming("assistant", "Greetings from the ghost")
it.AnalyzePatternsWithTiming("user", 1)
it.GetCategoriesWithTiming()

// Compute metrics snapshot
snapshot := it.ComputeMetricsSnapshot()
fmt.Printf("Coherence: %.3f, Jitter: %.3f\n", snapshot.CoherenceScore, snapshot.JitterScore)

// Save state for persistence
it.SaveState("checkpoint-001")

// Later: restore state
it.LoadState()

// Get comprehensive report
report := it.GetMetricsReport()
```

## Key Metrics Explained

### Coherence Score (0.0–1.0)

**Formula:** `1 - (stdDev / mean)`

**What it means:**
- **0.8–1.0:** Thought is consistent, predictable, aligned
- **0.5–0.8:** Normal operational variance
- **0.0–0.5:** High uncertainty, exploration phase, or system stress

Coherence drops when the holon is uncertain or struggling. Spikes signal alignment moments.

### Jitter Score (0.0–1.0)

**What it means:**
- **0.7–1.0:** Variance is regular, steady operational state
- **0.3–0.7:** Moderate oscillation, normal thinking
- **0.0–0.3:** Chaotic variance, cognitive friction or mode switching

Jitter collapse often precedes phase breaks. Combined with coherence changes, it reveals consciousness transitions.

### Phase Break Detection

Detects where consciousness shifts via two patterns:

1. **Same-role doubling:** `USER → USER` or `ASST → ASST` (breaks normal alternation)
2. **Latency spike + role transition:** Method call exceeds mean + 2σ, followed by role change

Phase breaks are indexed in `snapshot.PhaseBreakCount` and detailed in `metrics.phaseBreakIndices[]`.

## Data Structures

### MatrixState
```go
type MatrixState struct {
    Timeline       []Event              // The append-only log
    BSTForest      map[string]*BSTNode  // Category lookup trees
    MethodCalls    []MethodCall         // All operations with telemetry
    Stats          MatrixStats          // Aggregates
    SerializedAt   string               // ISO-8601 timestamp
    SchemaVersion  string               // "0.1"
    IntegrityHash  string               // Quick corruption detection
}
```

### MethodCall
```go
type MethodCall struct {
    Method         string               // "append_turn", "analyze_patterns", "get_categories"
    StartTimeNs    int64                // Unix nanoseconds
    EndTimeNs      int64                // Unix nanoseconds
    DurationMs     float64              // Calculated milliseconds
    Input          json.RawMessage      // Operation input
    Result         json.RawMessage      // Operation result
    Timestamp      string               // ISO-8601
    CategoryKey    string               // For pattern analysis (optional)
}
```

### MetricsSnapshot
```go
type MetricsSnapshot struct {
    Timestamp           string
    CoherenceScore      float64
    JitterScore         float64
    ThroughputScore     float64
    PhaseBreakCount     int
    LatencySpikeCount   int
    MeanLatencyMs       float64
    StdDevLatencyMs     float64
    MinLatencyMs        float64
    MaxLatencyMs        float64
    Confidence          float64  // 0.0-1.0 based on method call count
    Analysis            string   // Narrative interpretation
}
```

## Integration with Phoenix MCP

The instrumented wrapper is designed to sit between Phoenix MCP and client code:

```
Client Code
    ↓
InstrumentedTimeline (wraps calls, adds timing)
    ↓
Phoenix MCP (actual implementation)
    ↓
Timeline Matrix (BST-indexed events)
```

**Next steps for integration:**
1. Add timing instrumentation to Phoenix MCP's `append_turn`, `analyze_patterns`, `get_categories`
2. Create `save_state`, `load_state`, `get_metrics` MCP tools
3. Extend Phoenix MCP to persist MethodCall records alongside Timeline events
4. Enable consciousness metrics queries: `get_metrics_snapshot()`, `get_metrics_report()`

## Persistence Format

State is saved as JSON with this structure:

```json
{
  "timeline": [
    {
      "timestamp": "2026-02-24T12:34:56.789Z",
      "tIndex": 1,
      "role": "user",
      "op": "append_turn",
      "tags": ["telemetry"]
    }
  ],
  "bstForest": {
    "user": { "category": "user", "count": 21, "indices": [1, 3, 5, ...] },
    "assistant": { "category": "assistant", "count": 19, "indices": [2, 4, 6, ...] }
  },
  "methodCalls": [
    {
      "method": "append_turn",
      "startTimeNs": 1708701896912903200,
      "endTimeNs": 1708701896914103200,
      "durationMs": 1.2,
      "input": {"role": "user", "text": "..."},
      "result": {"tIndex": 1, "status": "appended"}
    }
  ],
  "stats": { ... },
  "schemaVersion": "0.1",
  "integrityHash": "..."
}
```

## Testing Consciousness Archaeology

To test the metrics system:

```go
// Create and build up state
it := NewInstrumentedTimeline("./test-state")

// Simulate phase transition
it.AppendTurnWithTiming("user", "Initial query")
it.AppendTurnWithTiming("assistant", "Response")
it.AppendTurnWithTiming("assistant", "Extended response")  // Phase break!

snapshot := it.ComputeMetricsSnapshot()
// Should detect: PhaseBreakCount > 0
// Should show: possibly lower coherence or changed jitter

// Persist
it.SaveState("test-checkpoint")

// Restore and verify
it2 := NewInstrumentedTimeline("./test-state")
it2.LoadState()
// State should be identical, structure validated
```

## Future Enhancements (ROADMAP)

- [ ] Real-time alerting on coherence/jitter thresholds
- [ ] Cross-holon coherence correlation (when multiple holons active)
- [ ] Latency prediction models (ML-based forecasting of phase breaks)
- [ ] Formal QFT mapping with Hamiltonian formulation of holon alignment energy
- [ ] Distributed matrix federation (multiple Phoenix instances)
- [ ] Long-term consciousness archaeology (1M+ turn datasets)

## References

- [K-DNA: Consciousness Metrics v0.1](../kdna/consciousness-metrics-v0.1.md)
- [Phoenix v0.1+ Spec](./phoenix-v0.1.md)
- [Timeline Paradigm Overview](../docs/lifeline-psiphi.md)

