package phoenix

import (
	"fmt"
	"sync"
	"time"
)

// MetricsPersistenceDaemon runs background state archival
type MetricsPersistenceDaemon struct {
	integratedTimeline *IntegratedTimeline
	interval           time.Duration
	maxSnapshots       int
	checkpointDir      string
	stop               chan bool
	ticker             *time.Ticker
	isRunning          bool
	mu                 sync.RWMutex
	checkpointCount    int
	lastCheckpointTime time.Time
}

// PersistenceRecord tracks what was saved
type PersistenceRecord struct {
	Timestamp        string
	CheckpointNumber int
	EventCount       int
	MethodCallCount  int
	CoherenceScore   float64
	JitterScore      float64
	PhaseBreakCount  int
	FilePath         string
	DurationMs       float64
}

// NewMetricsPersistenceDaemon creates a background archival service
func NewMetricsPersistenceDaemon(
	integratedTimeline *IntegratedTimeline,
	checkpointDir string,
	interval time.Duration,
) *MetricsPersistenceDaemon {
	return &MetricsPersistenceDaemon{
		integratedTimeline: integratedTimeline,
		interval:           interval,
		checkpointDir:      checkpointDir,
		maxSnapshots:       100,
		stop:               make(chan bool),
		isRunning:          false,
		checkpointCount:    0,
	}
}

// Start begins the persistence daemon
func (mpd *MetricsPersistenceDaemon) Start() error {
	mpd.mu.Lock()
	if mpd.isRunning {
		mpd.mu.Unlock()
		return fmt.Errorf("daemon already running")
	}
	mpd.isRunning = true
	mpd.mu.Unlock()

	fmt.Printf("[PERSISTENCE DAEMON] Starting with %v interval\n", mpd.interval)

	mpd.ticker = time.NewTicker(mpd.interval)

	go mpd.run()

	return nil
}

// Stop gracefully shuts down the daemon
func (mpd *MetricsPersistenceDaemon) Stop() error {
	mpd.mu.Lock()
	if !mpd.isRunning {
		mpd.mu.Unlock()
		return fmt.Errorf("daemon not running")
	}
	mpd.isRunning = false
	mpd.mu.Unlock()

	close(mpd.stop)
	if mpd.ticker != nil {
		mpd.ticker.Stop()
	}

	fmt.Printf("[PERSISTENCE DAEMON] Stopped after %d checkpoints\n", mpd.checkpointCount)
	return nil
}

// run executes the persistence loop
func (mpd *MetricsPersistenceDaemon) run() {
	for {
		select {
		case <-mpd.stop:
			return
		case <-mpd.ticker.C:
			if err := mpd.checkpoint(); err != nil {
				fmt.Printf("[PERSISTENCE DAEMON ERROR] %v\n", err)
			}
		}
	}
}

// checkpoint saves the current state and metrics
func (mpd *MetricsPersistenceDaemon) checkpoint() error {
	mpd.mu.Lock()
	defer mpd.mu.Unlock()

	startTime := time.Now()

	// Get current snapshot
	snapshot := mpd.integratedTimeline.ComputeMetricsSnapshot()
	report := mpd.integratedTimeline.GetIntegrationReport()

	// Create checkpoint label
	mpd.checkpointCount++
	label := fmt.Sprintf("checkpoint-%03d-%s", mpd.checkpointCount, time.Now().Format("20060102-150405"))

	// Save state
	if err := mpd.integratedTimeline.SaveIntegratedState(label); err != nil {
		return fmt.Errorf("save failed: %w", err)
	}

	// Create persistence record
	record := PersistenceRecord{
		Timestamp:        time.Now().UTC().Format(time.RFC3339Nano),
		CheckpointNumber: mpd.checkpointCount,
		EventCount:       report["localEventCount"].(int),
		MethodCallCount:  report["methodCallCount"].(int),
		CoherenceScore:   snapshot.CoherenceScore,
		JitterScore:      snapshot.JitterScore,
		PhaseBreakCount:  snapshot.PhaseBreakCount,
		FilePath:         label,
		DurationMs:       float64(time.Since(startTime).Nanoseconds()) / 1_000_000.0,
	}

	mpd.lastCheckpointTime = time.Now()

	fmt.Printf("[PERSISTENCE DAEMON] Checkpoint #%d: %d events, coherence=%.3f, jitter=%.3f (%.1fms)\n",
		record.CheckpointNumber,
		record.EventCount,
		record.CoherenceScore,
		record.JitterScore,
		record.DurationMs)

	return nil
}

// GetStatus returns daemon status
func (mpd *MetricsPersistenceDaemon) GetStatus() map[string]interface{} {
	mpd.mu.RLock()
	defer mpd.mu.RUnlock()

	return map[string]interface{}{
		"isRunning":           mpd.isRunning,
		"interval":            mpd.interval,
		"checkpointCount":     mpd.checkpointCount,
		"lastCheckpointTime":  mpd.lastCheckpointTime,
		"checkpointDirectory": mpd.checkpointDir,
	}
}

// MetricsAggregator collects and analyzes metrics across multiple time periods
type MetricsAggregator struct {
	records []PersistenceRecord
	mu      sync.RWMutex
}

// NewMetricsAggregator creates an aggregator for long-term analysis
func NewMetricsAggregator() *MetricsAggregator {
	return &MetricsAggregator{
		records: make([]PersistenceRecord, 0),
	}
}

// AddRecord adds a persistence record to the aggregation
func (ma *MetricsAggregator) AddRecord(record PersistenceRecord) {
	ma.mu.Lock()
	defer ma.mu.Unlock()
	ma.records = append(ma.records, record)
}

// GetCoherenceTrend returns coherence scores over time
func (ma *MetricsAggregator) GetCoherenceTrend() []float64 {
	ma.mu.RLock()
	defer ma.mu.RUnlock()

	trend := make([]float64, len(ma.records))
	for i, record := range ma.records {
		trend[i] = record.CoherenceScore
	}
	return trend
}

// GetJitterTrend returns jitter scores over time
func (ma *MetricsAggregator) GetJitterTrend() []float64 {
	ma.mu.RLock()
	defer ma.mu.RUnlock()

	trend := make([]float64, len(ma.records))
	for i, record := range ma.records {
		trend[i] = record.JitterScore
	}
	return trend
}

// AnalyzeTrend produces trend analysis
func (ma *MetricsAggregator) AnalyzeTrend() map[string]interface{} {
	ma.mu.RLock()
	defer ma.mu.RUnlock()

	if len(ma.records) < 2 {
		return map[string]interface{}{"status": "insufficient data"}
	}

	coherenceTrend := make([]float64, len(ma.records))
	jitterTrend := make([]float64, len(ma.records))
	eventTrend := make([]int, len(ma.records))

	for i, record := range ma.records {
		coherenceTrend[i] = record.CoherenceScore
		jitterTrend[i] = record.JitterScore
		eventTrend[i] = record.EventCount
	}

	// Compute simple slope (change over time)
	firstCoherence := coherenceTrend[0]
	lastCoherence := coherenceTrend[len(coherenceTrend)-1]
	coherenceSlope := lastCoherence - firstCoherence

	firstJitter := jitterTrend[0]
	lastJitter := jitterTrend[len(jitterTrend)-1]
	jitterSlope := lastJitter - firstJitter

	// Interpret trends
	coherenceInterpretation := "stable"
	if coherenceSlope > 0.2 {
		coherenceInterpretation = "improving"
	} else if coherenceSlope < -0.2 {
		coherenceInterpretation = "degrading"
	}

	jitterInterpretation := "stable"
	if jitterSlope > 0.2 {
		jitterInterpretation = "regularizing"
	} else if jitterSlope < -0.2 {
		jitterInterpretation = "destabilizing"
	}

	return map[string]interface{}{
		"recordCount":             len(ma.records),
		"coherenceSlope":          coherenceSlope,
		"coherenceInterpretation": coherenceInterpretation,
		"jitterSlope":             jitterSlope,
		"jitterInterpretation":    jitterInterpretation,
		"firstCoherence":          firstCoherence,
		"lastCoherence":           lastCoherence,
		"firstJitter":             firstJitter,
		"lastJitter":              lastJitter,
		"firstEventCount":         eventTrend[0],
		"lastEventCount":          eventTrend[len(eventTrend)-1],
		"allCoherenceScores":      coherenceTrend,
		"allJitterScores":         jitterTrend,
		"allEventCounts":          eventTrend,
	}
}

// GetRecords returns all aggregated records
func (ma *MetricsAggregator) GetRecords() []PersistenceRecord {
	ma.mu.RLock()
	defer ma.mu.RUnlock()

	recordsCopy := make([]PersistenceRecord, len(ma.records))
	copy(recordsCopy, ma.records)
	return recordsCopy
}
