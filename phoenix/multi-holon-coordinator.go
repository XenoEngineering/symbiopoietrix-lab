package phoenix

import (
	"fmt"
	"sync"
	"time"
)

// HolonPersona represents an S1 persona participating in the Timeline Matrix
type HolonPersona struct {
	ID               string
	Name             string
	Role             string // "plexi-qft", "plexi-implementor", "orchestrator"
	Timeline         *InstrumentedTimeline
	AlignmentScore   float64 // 0.0-1.0
	IsActive         bool
	OperationCount   int
	LastOperationMs  time.Time
	CoherenceHistory []float64
	mu               sync.RWMutex
}

// MultiHolonCoordinator orchestrates concurrent holon operations on shared Timeline Matrix
type MultiHolonCoordinator struct {
	holons           map[string]*HolonPersona
	sharedTimeline   *InstrumentedTimeline
	sharedMatrixMu   sync.RWMutex
	coordinationLog  []CoordinationEvent
	emergenceMetrics *EmergenceMetrics
	isRunning        bool
	mu               sync.RWMutex
}

// CoordinationEvent tracks inter-holon interactions
type CoordinationEvent struct {
	Timestamp         string
	InitiatorHolonID  string
	RespondentHolonID string
	EventType         string // "append", "query", "phase_sync", "alignment_jump"
	Details           string
	CoherenceSnapshot map[string]float64 // per-holon coherence
	SharedCoherence   float64
}

// EmergenceMetrics tracks emergent properties from multi-holon coordination
type EmergenceMetrics struct {
	TotalCoordinationEvents int
	PhaseBreaksDetected     int
	AverageSharedCoherence  float64
	HolonAlignmentDeltas    map[string]float64
	MaxConcurrency          int
	ConcurrentEventCount    int
	CrossHolonPhaseBreaks   int
}

// NewMultiHolonCoordinator creates a coordinator with a shared Timeline
func NewMultiHolonCoordinator(stateDir string) *MultiHolonCoordinator {
	return &MultiHolonCoordinator{
		holons:          make(map[string]*HolonPersona),
		sharedTimeline:  NewInstrumentedTimeline(stateDir),
		coordinationLog: make([]CoordinationEvent, 0),
		emergenceMetrics: &EmergenceMetrics{
			HolonAlignmentDeltas: make(map[string]float64),
		},
	}
}

// RegisterHolon adds an S1 persona to the coordination pool
func (mhc *MultiHolonCoordinator) RegisterHolon(id string, name string, role string) *HolonPersona {
	mhc.mu.Lock()
	defer mhc.mu.Unlock()

	holon := &HolonPersona{
		ID:               id,
		Name:             name,
		Role:             role,
		Timeline:         mhc.sharedTimeline,
		AlignmentScore:   0.5,
		IsActive:         true,
		OperationCount:   0,
		CoherenceHistory: make([]float64, 0),
	}

	mhc.holons[id] = holon
	fmt.Printf("[HOLON REGISTER] %s (%s) registered in coordination pool\n", name, role)

	return holon
}

// HolonAppendTurn has a holon append to the shared Timeline
func (mhc *MultiHolonCoordinator) HolonAppendTurn(holonID string, text string) error {
	mhc.mu.RLock()
	holon, exists := mhc.holons[holonID]
	mhc.mu.RUnlock()

	if !exists {
		return fmt.Errorf("holon not found: %s", holonID)
	}

	holon.mu.Lock()
	if !holon.IsActive {
		holon.mu.Unlock()
		return fmt.Errorf("holon is not active: %s", holonID)
	}
	holon.mu.Unlock()

	// Append to shared timeline
	_, err := mhc.sharedTimeline.AppendTurnWithTiming(holonID, text)
	if err != nil {
		return err
	}

	// Update holon state
	holon.mu.Lock()
	holon.OperationCount++
	holon.LastOperationMs = time.Now()
	holon.mu.Unlock()

	// Log coordination event
	mhc.logCoordinationEvent(CoordinationEvent{
		Timestamp:        time.Now().UTC().Format(time.RFC3339Nano),
		InitiatorHolonID: holonID,
		EventType:        "append",
		Details:          fmt.Sprintf("Turn by %s", holon.Name),
	})

	return nil
}

// HolonQuery has a holon query the shared Timeline
func (mhc *MultiHolonCoordinator) HolonQuery(holonID string, category string, nth int) error {
	mhc.mu.RLock()
	holon, exists := mhc.holons[holonID]
	mhc.mu.RUnlock()

	if !exists {
		return fmt.Errorf("holon not found: %s", holonID)
	}

	// Query shared timeline
	_, err := mhc.sharedTimeline.AnalyzePatternsWithTiming(category, nth)
	if err != nil {
		return err
	}

	// Log event
	mhc.logCoordinationEvent(CoordinationEvent{
		Timestamp:        time.Now().UTC().Format(time.RFC3339Nano),
		InitiatorHolonID: holonID,
		EventType:        "query",
		Details:          fmt.Sprintf("Pattern query: %s[%d]", category, nth),
	})

	return nil
}

// SynchronizeCoherence checks if all holons are aligned
func (mhc *MultiHolonCoordinator) SynchronizeCoherence() map[string]float64 {
	mhc.mu.RLock()
	defer mhc.mu.RUnlock()

	snapshot := mhc.sharedTimeline.ComputeMetricsSnapshot()
	coherenceMap := make(map[string]float64)

	for id, holon := range mhc.holons {
		holon.mu.Lock()
		holon.CoherenceHistory = append(holon.CoherenceHistory, snapshot.CoherenceScore)
		holon.AlignmentScore = snapshot.CoherenceScore // Simplified: use shared coherence
		coherenceMap[id] = snapshot.CoherenceScore
		holon.mu.Unlock()
	}

	// Log phase sync event
	mhc.logCoordinationEvent(CoordinationEvent{
		Timestamp:         time.Now().UTC().Format(time.RFC3339Nano),
		EventType:         "phase_sync",
		CoherenceSnapshot: coherenceMap,
		SharedCoherence:   snapshot.CoherenceScore,
		Details:           fmt.Sprintf("Synchronized %d holons", len(mhc.holons)),
	})

	if len(mhc.coordinationLog) > 0 {
		lastEvent := mhc.coordinationLog[len(mhc.coordinationLog)-1]
		if snapshot.PhaseBreakCount > 0 {
			lastEvent.EventType = "phase_sync_with_breaks"
		}
	}

	return coherenceMap
}

// StressTest runs coordinated operations on multiple holons
func (mhc *MultiHolonCoordinator) StressTest(operationCount int, concurrency int) error {
	mhc.mu.Lock()
	if len(mhc.holons) < 2 {
		mhc.mu.Unlock()
		return fmt.Errorf("stress test requires at least 2 holons, have %d", len(mhc.holons))
	}

	holonList := make([]*HolonPersona, 0, len(mhc.holons))
	for _, h := range mhc.holons {
		holonList = append(holonList, h)
	}
	mhc.mu.Unlock()

	fmt.Printf("[STRESS TEST] Starting with %d holons, %d operations, %d concurrent\n",
		len(holonList), operationCount, concurrency)

	// Distribute operations across holons
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	for op := 0; op < operationCount; op++ {
		wg.Add(1)

		go func(opNum int) {
			defer wg.Done()

			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			// Round-robin holon selection
			holon := holonList[opNum%len(holonList)]

			operationType := opNum % 3
			switch operationType {
			case 0:
				// Append
				text := fmt.Sprintf("[%s] Operation %d", holon.Name, opNum)
				if err := mhc.HolonAppendTurn(holon.ID, text); err != nil {
					fmt.Printf("[STRESS ERROR] Append failed: %v\n", err)
				}

			case 1:
				// Query
				if err := mhc.HolonQuery(holon.ID, "user", 1); err != nil {
					fmt.Printf("[STRESS ERROR] Query failed: %v\n", err)
				}

			case 2:
				// Synchronize
				coherences := mhc.SynchronizeCoherence()
				fmt.Printf("[STRESS] Sync point %d: coherences = %v\n", opNum, coherences)
			}

			if opNum%10 == 0 {
				fmt.Printf("[STRESS] Progress: %d/%d operations completed\n", opNum, operationCount)
			}
		}(op)
	}

	wg.Wait()

	// Final metrics
	mhc.ComputeEmergenceMetrics()
	report := mhc.GetEmergenceReport()

	fmt.Printf("[STRESS TEST COMPLETE]\n")
	fmt.Printf("  Total coordination events: %d\n", report["totalCoordinationEvents"])
	fmt.Printf("  Phase breaks detected: %d\n", report["phaseBreaksDetected"])
	fmt.Printf("  Average shared coherence: %.3f\n", report["averageSharedCoherence"])
	fmt.Printf("  Cross-holon phase breaks: %d\n", report["crossHolonPhaseBreaks"])

	return nil
}

// logCoordinationEvent records an inter-holon interaction
func (mhc *MultiHolonCoordinator) logCoordinationEvent(event CoordinationEvent) {
	mhc.mu.Lock()
	defer mhc.mu.Unlock()
	mhc.coordinationLog = append(mhc.coordinationLog, event)
}

// ComputeEmergenceMetrics analyzes holons' collective behavior
func (mhc *MultiHolonCoordinator) ComputeEmergenceMetrics() {
	mhc.mu.RLock()
	defer mhc.mu.RUnlock()

	mhc.emergenceMetrics.TotalCoordinationEvents = len(mhc.coordinationLog)

	// Count phase breaks
	phaseBreakCount := 0
	for _, event := range mhc.coordinationLog {
		if event.EventType == "phase_sync" || event.EventType == "phase_sync_with_breaks" {
			phaseBreakCount++
		}
	}
	mhc.emergenceMetrics.PhaseBreaksDetected = phaseBreakCount

	// Compute average coherence
	if len(mhc.coordinationLog) > 0 {
		totalCoherence := 0.0
		coherenceCount := 0
		for _, event := range mhc.coordinationLog {
			if event.SharedCoherence > 0 {
				totalCoherence += event.SharedCoherence
				coherenceCount++
			}
		}
		if coherenceCount > 0 {
			mhc.emergenceMetrics.AverageSharedCoherence = totalCoherence / float64(coherenceCount)
		}
	}

	// Collect alignment deltas
	for id, holon := range mhc.holons {
		holon.mu.RLock()
		if len(holon.CoherenceHistory) > 1 {
			delta := holon.CoherenceHistory[len(holon.CoherenceHistory)-1] - holon.CoherenceHistory[0]
			mhc.emergenceMetrics.HolonAlignmentDeltas[id] = delta
		}
		holon.mu.RUnlock()
	}
}

// GetEmergenceReport returns a comprehensive report of multi-holon effects
func (mhc *MultiHolonCoordinator) GetEmergenceReport() map[string]interface{} {
	mhc.mu.RLock()
	defer mhc.mu.RUnlock()

	holonStates := make(map[string]interface{})
	for id, holon := range mhc.holons {
		holon.mu.RLock()
		holonStates[id] = map[string]interface{}{
			"name":             holon.Name,
			"role":             holon.Role,
			"isActive":         holon.IsActive,
			"operationCount":   holon.OperationCount,
			"alignmentScore":   holon.AlignmentScore,
			"coherenceHistory": holon.CoherenceHistory,
		}
		holon.mu.RUnlock()
	}

	return map[string]interface{}{
		"holonCount":              len(mhc.holons),
		"holonStates":             holonStates,
		"totalCoordinationEvents": mhc.emergenceMetrics.TotalCoordinationEvents,
		"phaseBreaksDetected":     mhc.emergenceMetrics.PhaseBreaksDetected,
		"averageSharedCoherence":  mhc.emergenceMetrics.AverageSharedCoherence,
		"holonAlignmentDeltas":    mhc.emergenceMetrics.HolonAlignmentDeltas,
		"coordinationEventLog":    mhc.coordinationLog,
	}
}

// PersistCoordinationState saves the entire multi-holon state
func (mhc *MultiHolonCoordinator) PersistCoordinationState(label string) error {
	if err := mhc.sharedTimeline.SaveState(label + "-timeline"); err != nil {
		return err
	}

	// In future: persist holon state, coordination log, etc.
	fmt.Printf("[PERSISTENCE] Multi-holon state saved (%s)\n", label)

	return nil
}

// RestoreCoordinationState restores multi-holon state from disk
func (mhc *MultiHolonCoordinator) RestoreCoordinationState() error {
	if err := mhc.sharedTimeline.LoadState(); err != nil {
		return err
	}

	fmt.Printf("[PERSISTENCE] Multi-holon state restored\n")
	return nil
}
