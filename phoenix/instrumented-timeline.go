package phoenix

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// InstrumentedTimeline wraps the Phoenix MCP client to add timing telemetry
type InstrumentedTimeline struct {
	mu               sync.Mutex
	currentState     *MatrixState
	persistence      *MatrixPersistence
	metrics          *ConsciousnessMetrics
	metricsSnapshots []MetricsSnapshot
	isRecording      bool
}

// NewInstrumentedTimeline creates a wrapper around the Phoenix MCP interface
func NewInstrumentedTimeline(stateDir string) *InstrumentedTimeline {
	return &InstrumentedTimeline{
		currentState: &MatrixState{
			Timeline:      make([]Event, 0),
			BSTForest:     make(map[string]*BSTNode),
			MethodCalls:   make([]MethodCall, 0),
			SchemaVersion: "0.1",
		},
		persistence:      NewMatrixPersistence(stateDir),
		metrics:          NewConsciousnessMetrics(),
		metricsSnapshots: make([]MetricsSnapshot, 0),
		isRecording:      true,
	}
}

// AppendTurnWithTiming wraps an append operation with instrumentation
// This will be called after Phoenix MCP's actual append_turn succeeds
func (it *InstrumentedTimeline) AppendTurnWithTiming(role string, text string) (*MethodCall, error) {
	it.mu.Lock()
	defer it.mu.Unlock()

	startNs := time.Now().UnixNano()

	// In real integration, this would call Phoenix MCP's actual append_turn
	// For now, we create a mock event
	tIndex := len(it.currentState.Timeline) + 1
	event := Event{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		TIndex:    tIndex,
		Role:      role,
		Layer:     "artifact",
		Op:        "append_turn",
		Tags:      []string{"telemetry"},
	}

	it.currentState.Timeline = append(it.currentState.Timeline, event)

	endNs := time.Now().UnixNano()
	durationMs := float64(endNs-startNs) / 1_000_000.0

	// Record the method call
	methodCall := MethodCall{
		Method:      "append_turn",
		StartTimeNs: startNs,
		EndTimeNs:   endNs,
		DurationMs:  durationMs,
		Timestamp:   time.Now().UTC().Format(time.RFC3339Nano),
		Input: json.RawMessage(
			[]byte(fmt.Sprintf(`{"role":"%s","text":"%s"}`, role, text)),
		),
		Result: json.RawMessage(
			[]byte(fmt.Sprintf(`{"tIndex":%d,"status":"appended"}`, tIndex)),
		),
	}

	it.currentState.MethodCalls = append(it.currentState.MethodCalls, methodCall)

	if it.isRecording {
		fmt.Printf("[TELEMETRY] append_turn: %.3fms (turn #%d)\n", durationMs, tIndex)
	}

	return &methodCall, nil
}

// AnalyzePatternsWithTiming wraps pattern analysis with instrumentation
func (it *InstrumentedTimeline) AnalyzePatternsWithTiming(category string, nth int) (*MethodCall, error) {
	it.mu.Lock()
	defer it.mu.Unlock()

	startNs := time.Now().UnixNano()

	// In real integration, this would call Phoenix MCP's actual analyze_patterns
	// For now, simple linear search simulation
	occurrenceCount := 0
	foundIndex := -1
	for i, event := range it.currentState.Timeline {
		if event.Role == category {
			occurrenceCount++
			if occurrenceCount == nth {
				foundIndex = i
				break
			}
		}
	}

	endNs := time.Now().UnixNano()
	durationMs := float64(endNs-startNs) / 1_000_000.0

	methodCall := MethodCall{
		Method:      "analyze_patterns",
		StartTimeNs: startNs,
		EndTimeNs:   endNs,
		DurationMs:  durationMs,
		Timestamp:   time.Now().UTC().Format(time.RFC3339Nano),
		CategoryKey: category,
		Input: json.RawMessage(
			[]byte(fmt.Sprintf(`{"category":"%s","nth":%d}`, category, nth)),
		),
		Result: json.RawMessage(
			[]byte(fmt.Sprintf(`{"found":%v,"tIndex":%d,"totalOccurrences":%d}`, foundIndex >= 0, foundIndex+1, occurrenceCount)),
		),
	}

	it.currentState.MethodCalls = append(it.currentState.MethodCalls, methodCall)

	if it.isRecording {
		found := "found"
		if foundIndex < 0 {
			found = "not_found"
		}
		fmt.Printf("[TELEMETRY] analyze_patterns[%s,%d]: %.3fms (%s)\n", category, nth, durationMs, found)
	}

	return &methodCall, nil
}

// GetCategoriesWithTiming wraps category lookup with instrumentation
func (it *InstrumentedTimeline) GetCategoriesWithTiming() (*MethodCall, error) {
	it.mu.Lock()
	defer it.mu.Unlock()

	startNs := time.Now().UnixNano()

	// Rebuild category census
	categoryMap := make(map[string]int)
	for _, event := range it.currentState.Timeline {
		categoryMap[event.Role]++
	}

	endNs := time.Now().UnixNano()
	durationMs := float64(endNs-startNs) / 1_000_000.0

	// Build result JSON
	resultBytes, _ := json.Marshal(categoryMap)

	methodCall := MethodCall{
		Method:      "get_categories",
		StartTimeNs: startNs,
		EndTimeNs:   endNs,
		DurationMs:  durationMs,
		Timestamp:   time.Now().UTC().Format(time.RFC3339Nano),
		Result:      json.RawMessage(resultBytes),
	}

	it.currentState.MethodCalls = append(it.currentState.MethodCalls, methodCall)

	if it.isRecording {
		fmt.Printf("[TELEMETRY] get_categories: %.3fms (categories: %d)\n", durationMs, len(categoryMap))
	}

	return &methodCall, nil
}

// ComputeMetricsSnapshot generates consciousness metrics at current point
func (it *InstrumentedTimeline) ComputeMetricsSnapshot() *MetricsSnapshot {
	it.mu.Lock()
	defer it.mu.Unlock()

	snapshot := it.metrics.ComputeFromState(it.currentState)
	it.metricsSnapshots = append(it.metricsSnapshots, *snapshot)

	if it.isRecording {
		fmt.Printf("[METRICS] Coherence: %.3f | Jitter: %.3f | Spikes: %d | Breaks: %d\n",
			snapshot.CoherenceScore, snapshot.JitterScore,
			snapshot.LatencySpikeCount, snapshot.PhaseBreakCount)
		fmt.Printf("          Mean latency: %.3fms | Turns: %d | Calls: %d\n",
			snapshot.MeanLatencyMs, it.currentState.Stats.TotalEvents,
			it.currentState.Stats.MethodCallCount)
		if snapshot.Analysis != "" {
			fmt.Printf("          Analysis: %s\n", snapshot.Analysis)
		}
	}

	return snapshot
}

// SaveState persists the entire Matrix to disk
func (it *InstrumentedTimeline) SaveState(label string) error {
	it.mu.Lock()
	defer it.mu.Unlock()

	it.currentState.Stats.TotalEvents = len(it.currentState.Timeline)
	it.currentState.Stats.MethodCallCount = len(it.currentState.MethodCalls)
	it.currentState.Stats.LastUpdateTimestamp = time.Now().UTC().Format(time.RFC3339Nano)

	if err := it.persistence.SaveState(it.currentState); err != nil {
		return fmt.Errorf("save failed: %w", err)
	}

	if it.isRecording {
		fmt.Printf("[PERSISTENCE] Saved state (%s): %d events, %d method calls\n",
			label, len(it.currentState.Timeline), len(it.currentState.MethodCalls))
	}

	return nil
}

// LoadState restores state from disk
func (it *InstrumentedTimeline) LoadState() error {
	it.mu.Lock()
	defer it.mu.Unlock()

	state, err := it.persistence.LoadState()
	if err != nil {
		return fmt.Errorf("load failed: %w", err)
	}

	it.currentState = state
	if it.isRecording {
		fmt.Printf("[PERSISTENCE] Loaded state: %d events, %d method calls, structure validated\n",
			len(it.currentState.Timeline), len(it.currentState.MethodCalls))
	}

	return nil
}

// GetMetricsReport generates a comprehensive report of consciousness metrics over time
func (it *InstrumentedTimeline) GetMetricsReport() map[string]interface{} {
	it.mu.Lock()
	defer it.mu.Unlock()

	if len(it.metricsSnapshots) == 0 {
		return map[string]interface{}{"status": "no metrics snapshots recorded"}
	}

	first := it.metricsSnapshots[0]
	last := it.metricsSnapshots[len(it.metricsSnapshots)-1]

	return map[string]interface{}{
		"snapshotCount":      len(it.metricsSnapshots),
		"initialCoherence":   first.CoherenceScore,
		"currentCoherence":   last.CoherenceScore,
		"coherenceChange":    last.CoherenceScore - first.CoherenceScore,
		"initialJitter":      first.JitterScore,
		"currentJitter":      last.JitterScore,
		"jitterChange":       last.JitterScore - first.JitterScore,
		"totalPhaseBreaks":   last.PhaseBreakCount,
		"totalLatencySpikes": last.LatencySpikeCount,
		"meanLatencyMs":      last.MeanLatencyMs,
		"stdDevLatencyMs":    last.StdDevLatencyMs,
		"minLatencyMs":       last.MinLatencyMs,
		"maxLatencyMs":       last.MaxLatencyMs,
		"confidence":         last.Confidence,
		"currentAnalysis":    last.Analysis,
		"allSnapshots":       it.metricsSnapshots,
	}
}

// SetTelemetryRecording enables/disables console output
func (it *InstrumentedTimeline) SetTelemetryRecording(enabled bool) {
	it.mu.Lock()
	defer it.mu.Unlock()
	it.isRecording = enabled
}

// GetCurrentState returns a read-only snapshot of the matrix state
func (it *InstrumentedTimeline) GetCurrentState() *MatrixState {
	it.mu.Lock()
	defer it.mu.Unlock()
	// Return a shallow copy to prevent external mutation
	stateCopy := *it.currentState
	return &stateCopy
}
