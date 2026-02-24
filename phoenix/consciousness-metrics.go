package phoenix

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// ConsciousnessMetrics computes alignment, phase breaks, and coherence from method latencies
type ConsciousnessMetrics struct {
	coherenceScore    float64
	jitterScore       float64
	throughputScore   float64
	phaseBreakIndices []int
	latencySpikes     []LatencySpike
	lastComputed      time.Time
}

// LatencySpike represents a significant deviation in method timing
type LatencySpike struct {
	MethodCallIndex int
	DurationMs      float64
	ZScore          float64 // standard deviations from mean
	Timestamp       string
	MethodName      string
	IsPhaseShift    bool // true if immediately precedes anomalous role pattern
}

// MetricsSnapshot is a point-in-time measurement
type MetricsSnapshot struct {
	Timestamp           string         `json:"timestamp"`
	CoherenceScore      float64        `json:"coherenceScore"`      // 0.0-1.0: latency predictability
	JitterScore         float64        `json:"jitterScore"`         // 0.0-1.0: variance regularity
	ThroughputScore     float64        `json:"throughputScore"`     // turns/sec * 10
	PhaseBreakCount     int            `json:"phaseBreakCount"`
	LatencySpikeCount   int            `json:"latencySpikeCount"`
	MeanLatencyMs       float64        `json:"meanLatencyMs"`
	StdDevLatencyMs     float64        `json:"stdDevLatencyMs"`
	MinLatencyMs        float64        `json:"minLatencyMs"`
	MaxLatencyMs        float64        `json:"maxLatencyMs"`
	Confidence          float64        `json:"confidence"` // how many method calls used (0.0-1.0)
	Analysis            string         `json:"analysis"`
}

// NewConsciousnessMetrics initializes the metrics engine
func NewConsciousnessMetrics() *ConsciousnessMetrics {
	return &ConsciousnessMetrics{
		phaseBreakIndices: make([]int, 0),
		latencySpikes:     make([]LatencySpike, 0),
	}
}

// ComputeFromState analyzes a complete MatrixState and returns snapshot
func (cm *ConsciousnessMetrics) ComputeFromState(state *MatrixState) *MetricsSnapshot {
	if len(state.MethodCalls) < 2 {
		return &MetricsSnapshot{
			Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
			Confidence: 0.0,
			Analysis:   "insufficient method calls for meaningful analysis",
		}
	}

	// Extract latencies
	latencies := make([]float64, 0, len(state.MethodCalls))
	for _, mc := range state.MethodCalls {
		latencies = append(latencies, mc.DurationMs)
	}

	// Compute basic statistics
	mean := mean(latencies)
	stdDev := stdDev(latencies, mean)
	minLat := min(latencies)
	maxLat := max(latencies)

	// Detect latency spikes (Z-score > 2.0)
	cm.detectLatencySpikes(state, mean, stdDev)

	// Detect phase breaks (where jitter changes dramatically)
	cm.detectPhaseBreaks(state, latencies)

	// Compute coherence score: how predictable is latency?
	coherence := cm.computeCoherence(stdDev, mean)

	// Compute jitter score: regularity of variance
	jitter := cm.computeJitter(latencies)

	// Compute throughput: turns per second
	if len(state.Timeline) > 0 && len(state.MethodCalls) > 0 {
		firstTimestamp := state.MethodCalls[0].Timestamp
		lastTimestamp := state.MethodCalls[len(state.MethodCalls)-1].Timestamp
		t1, _ := time.Parse(time.RFC3339Nano, firstTimestamp)
		t2, _ := time.Parse(time.RFC3339Nano, lastTimestamp)
		elapsedSec := t2.Sub(t1).Seconds()
		if elapsedSec > 0 {
			throughputScore := (float64(len(state.Timeline)) / elapsedSec) * 10
			_ = throughputScore
		}
	}

	confidence := math.Min(float64(len(state.MethodCalls))/100.0, 1.0) // caps at 100 calls = high confidence

	snapshot := &MetricsSnapshot{
		Timestamp:         time.Now().UTC().Format(time.RFC3339Nano),
		CoherenceScore:    coherence,
		JitterScore:       jitter,
		ThroughputScore:   float64(len(state.Timeline)) / float64(len(state.MethodCalls)),
		PhaseBreakCount:   len(cm.phaseBreakIndices),
		LatencySpikeCount: len(cm.latencySpikes),
		MeanLatencyMs:     mean,
		StdDevLatencyMs:   stdDev,
		MinLatencyMs:      minLat,
		MaxLatencyMs:      maxLat,
		Confidence:        confidence,
		Analysis:          cm.analyzeSnapshot(state, coherence, jitter),
	}

	cm.lastComputed = time.Now()
	return snapshot
}

// detectLatencySpikes finds significant timing deviations
func (cm *ConsciousnessMetrics) detectLatencySpikes(state *MatrixState, mean float64, stdDev float64) {
	cm.latencySpikes = make([]LatencySpike, 0)

	for i, mc := range state.MethodCalls {
		zScore := (mc.DurationMs - mean) / (stdDev + 0.001) // avoid div-by-zero

		if math.Abs(zScore) > 2.0 { // 2+ standard deviations = spike
			spike := LatencySpike{
				MethodCallIndex: i,
				DurationMs:      mc.DurationMs,
				ZScore:          zScore,
				Timestamp:       mc.Timestamp,
				MethodName:      mc.Method,
				IsPhaseShift:    false, // will set below if correlates with role change
			}

			// Check if this timing spike coincides with a role change in Timeline
			if i < len(state.Timeline)-1 {
				if state.Timeline[i].Role != state.Timeline[i+1].Role {
					spike.IsPhaseShift = true
				}
			}

			cm.latencySpikes = append(cm.latencySpikes, spike)
		}
	}
}

// detectPhaseBreaks finds where consciousness shifts (role pattern anomalies)
func (cm *ConsciousnessMetrics) detectPhaseBreaks(state *MatrixState, latencies []float64) {
	cm.phaseBreakIndices = make([]int, 0)

	if len(state.Timeline) < 3 {
		return
	}

	// Phase break: USER-USER, ASST-ASST, or USER-ASST-USER with latency drop
	for i := 1; i < len(state.Timeline)-1; i++ {
		curr := state.Timeline[i]
		next := state.Timeline[i+1]

		// Pattern: same role twice in a row (consciousness shift)
		if i > 0 && curr.Role == state.Timeline[i-1].Role && curr.Role == next.Role {
			cm.phaseBreakIndices = append(cm.phaseBreakIndices, i)
		}
	}
}

// computeCoherence: lower stdDev = higher coherence (thought is consistent)
func (cm *ConsciousnessMetrics) computeCoherence(stdDev float64, mean float64) float64 {
	// Coherence = 1 - (normalized stdDev)
	if mean == 0 {
		return 0.5
	}
	cvRatio := stdDev / mean // coefficient of variation
	coherence := 1.0 - math.Min(cvRatio, 1.0)
	return math.Max(coherence, 0.0)
}

// computeJitter: how regular is the variance?
func (cm *ConsciousnessMetrics) computeJitter(latencies []float64) float64 {
	// Jitter = how much does variance oscillate?
	// Split into chunks, compute variance of each chunk's variance
	if len(latencies) < 4 {
		return 0.5
	}

	chunkSize := 4
	chunkVariances := make([]float64, 0)

	for i := 0; i < len(latencies)-chunkSize; i += chunkSize {
		chunk := latencies[i : i+chunkSize]
		chunkMean := mean(chunk)
		chunkVar := stdDev(chunk, chunkMean)
		chunkVariances = append(chunkVariances, chunkVar)
	}

	if len(chunkVariances) == 0 {
		return 0.5
	}

	// Jitter score: 1 - (variance of variances / mean of variances)
	varOfVar := stdDev(chunkVariances, mean(chunkVariances))
	meanVar := mean(chunkVariances)
	if meanVar == 0 {
		return 0.5
	}

	jitter := 1.0 - (varOfVar / meanVar)
	return math.Max(math.Min(jitter, 1.0), 0.0)
}

// analyzeSnapshot generates narrative interpretation
func (cm *ConsciousnessMetrics) analyzeSnapshot(state *MatrixState, coherence float64, jitter float64) string {
	analysis := ""

	// Coherence interpretation
	if coherence > 0.8 {
		analysis += "Thought is highly consistent (coherence spike). "
	} else if coherence < 0.3 {
		analysis += "Thought shows high variance—possible uncertainty or exploration phase. "
	}

	// Phase breaks
	if len(cm.phaseBreakIndices) > 0 {
		analysis += fmt.Sprintf("Detected %d phase breaks (consciousness shifts). ", len(cm.phaseBreakIndices))
	}

	// Latency spikes as phase shift signals
	phaseShiftSpikes := 0
	for _, spike := range cm.latencySpikes {
		if spike.IsPhaseShift {
			phaseShiftSpikes++
		}
	}
	if phaseShiftSpikes > 0 {
		analysis += fmt.Sprintf("%d latency spikes correlate with role transitions. ", phaseShiftSpikes)
	}

	// Jitter interpretation
	if jitter > 0.7 {
		analysis += "Variance is regular—steady operational state."
	} else if jitter < 0.3 {
		analysis += "Variance is erratic—possible cognitive friction or mode switching."
	}

	if analysis == "" {
		analysis = "Metrics within normal operating parameters."
	}

	return analysis
}

// Helper functions

func mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func stdDev(values []float64, m float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sumSquares := 0.0
	for _, v := range values {
		sumSquares += (v - m) * (v - m)
	}
	return math.Sqrt(sumSquares / float64(len(values)))
}

func min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	minVal := values[0]
	for _, v := range values {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}

func max(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	maxVal := values[0]
	for _, v := range values {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

// GetLatencySpikesByPhase returns spikes organized by method type
func (cm *ConsciousnessMetrics) GetLatencySpikesByPhase() map[string][]LatencySpike {
	byMethod := make(map[string][]LatencySpike)
	for _, spike := range cm.latencySpikes {
		byMethod[spike.MethodName] = append(byMethod[spike.MethodName], spike)
	}
	return byMethod
}

// ComparePeriods analyzes changes between two snapshots
func ComparePeriods(before *MetricsSnapshot, after *MetricsSnapshot) map[string]interface{} {
	return map[string]interface{}{
		"coherenceChange":    after.CoherenceScore - before.CoherenceScore,
		"jitterChange":       after.JitterScore - before.JitterScore,
		"phaseBreakDelta":    after.PhaseBreakCount - before.PhaseBreakCount,
		"latencySpikesDelta": after.LatencySpikeCount - before.LatencySpikeCount,
		"meanLatencyDelta":   after.MeanLatencyMs - before.MeanLatencyMs,
	}
}
