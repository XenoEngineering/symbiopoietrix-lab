package phoenix

import (
	"fmt"
	"sync"
	"time"
)

// FederatedNode represents a single Phoenix instance in the mesh
type FederatedNode struct {
	NodeID               string
	Address              string // e.g., "localhost:8080"
	Role                 string // "observer", "coordinator", "aggregator"
	LastHeartbeat        time.Time
	IsHealthy            bool
	LocalMetricsSnapshot *MetricsSnapshot
	mu                   sync.RWMutex
}

// FederationMesh coordinates multiple Phoenix nodes
type FederationMesh struct {
	meshID            string
	nodes             map[string]*FederatedNode
	globalAggregator  *MetricsAggregator
	crossNodeLog      []CrossNodeEvent
	federationMetrics *FederationMetrics
	coherenceBridge   *CoherenceBridge
	isRunning         bool
	stopChan          chan bool
	heartbeatTicker   *time.Ticker
	mu                sync.RWMutex
}

// CrossNodeEvent tracks inter-node consciousness patterns
type CrossNodeEvent struct {
	Timestamp             string
	InitiatorNodeID       string
	TargetNodeID          string
	EventType             string // "phase_sync", "coherence_spike", "phase_break_detected", "alignment_cascade"
	LocalCoherence        float64
	TargetCoherence       float64
	CoherenceDelta        float64
	GlobalPhaseBreakCount int
	Analysis              string
}

// FederationMetrics tracks emergence across the mesh
type FederationMetrics struct {
	TotalNodes              int
	HealthyNodes            int
	GlobalCoherence         float64 // averaged across all nodes
	GlobalJitter            float64
	CrossNodePhaseBreaks    int
	CoherenceCascades       int
	AverageNetworkLatencyMs float64
	HighestLocalCoherence   float64
	LowestLocalCoherence    float64
	CoherenceVariance       float64
	LastSync                time.Time
}

// CoherenceBridge detects when coherence changes propagate across nodes
type CoherenceBridge struct {
	nodeCoherenceHistory map[string][]float64
	cascadeThreshold     float64 // 0.15 = 15% change triggers cascade detection
	maxHistoryPoints     int
	mu                   sync.RWMutex
}

// NewFederationMesh creates a distributed coordination mesh
func NewFederationMesh(meshID string) *FederationMesh {
	return &FederationMesh{
		meshID:            meshID,
		nodes:             make(map[string]*FederatedNode),
		globalAggregator:  NewMetricsAggregator(),
		crossNodeLog:      make([]CrossNodeEvent, 0),
		federationMetrics: &FederationMetrics{},
		coherenceBridge: &CoherenceBridge{
			nodeCoherenceHistory: make(map[string][]float64),
			cascadeThreshold:     0.15,
			maxHistoryPoints:     100,
		},
		stopChan:  make(chan bool),
		isRunning: false,
	}
}

// RegisterNode joins a Phoenix instance to the mesh
func (fm *FederationMesh) RegisterNode(nodeID string, address string, role string) *FederatedNode {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	node := &FederatedNode{
		NodeID:        nodeID,
		Address:       address,
		Role:          role,
		LastHeartbeat: time.Now(),
		IsHealthy:     true,
	}

	fm.nodes[nodeID] = node
	fm.coherenceBridge.nodeCoherenceHistory[nodeID] = make([]float64, 0)

	fmt.Printf("[FEDERATION] Node %s registered (%s at %s)\n", nodeID, role, address)
	return node
}

// Start begins the federation mesh
func (fm *FederationMesh) Start() error {
	fm.mu.Lock()
	if fm.isRunning {
		fm.mu.Unlock()
		return fmt.Errorf("mesh already running")
	}
	fm.isRunning = true
	fm.mu.Unlock()

	fmt.Printf("[FEDERATION] Mesh '%s' starting with %d nodes\n", fm.meshID, len(fm.nodes))

	fm.heartbeatTicker = time.NewTicker(5 * time.Second)

	go fm.federationLoop()

	return nil
}

// Stop gracefully shuts down the mesh
func (fm *FederationMesh) Stop() error {
	fm.mu.Lock()
	if !fm.isRunning {
		fm.mu.Unlock()
		return fmt.Errorf("mesh not running")
	}
	fm.isRunning = false
	fm.mu.Unlock()

	close(fm.stopChan)
	if fm.heartbeatTicker != nil {
		fm.heartbeatTicker.Stop()
	}

	fmt.Printf("[FEDERATION] Mesh '%s' stopped\n", fm.meshID)
	return nil
}

// federationLoop runs the mesh synchronization
func (fm *FederationMesh) federationLoop() {
	for {
		select {
		case <-fm.stopChan:
			return
		case <-fm.heartbeatTicker.C:
			fm.syncAcrossMesh()
		}
	}
}

// syncAcrossMesh synchronizes metrics across all nodes
func (fm *FederationMesh) syncAcrossMesh() {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	// Collect metrics from all nodes
	nodeMetrics := make(map[string]*MetricsSnapshot)
	healthyCount := 0

	for nodeID, node := range fm.nodes {
		node.mu.RLock()
		if node.IsHealthy && node.LocalMetricsSnapshot != nil {
			nodeMetrics[nodeID] = node.LocalMetricsSnapshot
			healthyCount++
		}
		lastHB := node.LastHeartbeat
		node.mu.RUnlock()

		// Check heartbeat timeout (30 seconds)
		if time.Since(lastHB) > 30*time.Second {
			node.mu.Lock()
			node.IsHealthy = false
			node.mu.Unlock()
		}
	}

	fm.federationMetrics.TotalNodes = len(fm.nodes)
	fm.federationMetrics.HealthyNodes = healthyCount
	fm.federationMetrics.LastSync = time.Now()

	if len(nodeMetrics) == 0 {
		return
	}

	// Compute global metrics
	fm.computeGlobalMetrics(nodeMetrics)

	// Detect coherence cascades (phase breaks propagating across nodes)
	fm.detectCoherenceCascades(nodeMetrics)

	// Detect cross-node phase breaks
	fm.detectCrossNodePhaseBreaks(nodeMetrics)
}

// computeGlobalMetrics aggregates metrics across healthy nodes
func (fm *FederationMesh) computeGlobalMetrics(nodeMetrics map[string]*MetricsSnapshot) {
	coherenceSum := 0.0
	jitterSum := 0.0
	coherenceValues := make([]float64, 0)

	for nodeID, snapshot := range nodeMetrics {
		coherenceSum += snapshot.CoherenceScore
		jitterSum += snapshot.JitterScore
		coherenceValues = append(coherenceValues, snapshot.CoherenceScore)

		// Record in coherence history
		fm.coherenceBridge.mu.Lock()
		history := fm.coherenceBridge.nodeCoherenceHistory[nodeID]
		if len(history) >= fm.coherenceBridge.maxHistoryPoints {
			history = history[1:]
		}
		history = append(history, snapshot.CoherenceScore)
		fm.coherenceBridge.nodeCoherenceHistory[nodeID] = history
		fm.coherenceBridge.mu.Unlock()
	}

	count := float64(len(nodeMetrics))
	fm.federationMetrics.GlobalCoherence = coherenceSum / count
	fm.federationMetrics.GlobalJitter = jitterSum / count

	// Compute coherence variance
	variance := computeVariance(coherenceValues, fm.federationMetrics.GlobalCoherence)
	fm.federationMetrics.CoherenceVariance = variance

	// Min/max coherence
	fm.federationMetrics.HighestLocalCoherence = maxFloat(coherenceValues)
	fm.federationMetrics.LowestLocalCoherence = minFloat(coherenceValues)
}

// detectCoherenceCascades finds when coherence spikes propagate across nodes
func (fm *FederationMesh) detectCoherenceCascades(nodeMetrics map[string]*MetricsSnapshot) {
	fm.coherenceBridge.mu.RLock()
	defer fm.coherenceBridge.mu.RUnlock()

	nodeIDs := make([]string, 0)
	for nodeID := range nodeMetrics {
		nodeIDs = append(nodeIDs, nodeID)
	}

	// Check pairs of nodes for synchronized coherence jumps
	for i := 0; i < len(nodeIDs)-1; i++ {
		for j := i + 1; j < len(nodeIDs); j++ {
			nodeA := nodeIDs[i]
			nodeB := nodeIDs[j]

			historyA := fm.coherenceBridge.nodeCoherenceHistory[nodeA]
			historyB := fm.coherenceBridge.nodeCoherenceHistory[nodeB]

			if len(historyA) < 2 || len(historyB) < 2 {
				continue
			}

			// Get last change in each node
			deltaA := historyA[len(historyA)-1] - historyA[len(historyA)-2]
			deltaB := historyB[len(historyB)-1] - historyB[len(historyB)-2]

			// If both changed in same direction by > threshold, it's a cascade
			if (deltaA > fm.coherenceBridge.cascadeThreshold && deltaB > fm.coherenceBridge.cascadeThreshold) ||
				(deltaA < -fm.coherenceBridge.cascadeThreshold && deltaB < -fm.coherenceBridge.cascadeThreshold) {

				fm.federationMetrics.CoherenceCascades++

				event := CrossNodeEvent{
					Timestamp:       time.Now().UTC().Format(time.RFC3339Nano),
					InitiatorNodeID: nodeA,
					TargetNodeID:    nodeB,
					EventType:       "coherence_cascade",
					LocalCoherence:  nodeMetrics[nodeA].CoherenceScore,
					TargetCoherence: nodeMetrics[nodeB].CoherenceScore,
					CoherenceDelta:  deltaA,
					Analysis: fmt.Sprintf("Synchronized coherence jump: %+.3f (both nodes in phase transition)",
						(deltaA+deltaB)/2),
				}

				fm.crossNodeLog = append(fm.crossNodeLog, event)
				fmt.Printf("[FEDERATION CASCADE] %s → %s: coherence synchronized jump\n", nodeA, nodeB)
			}
		}
	}
}

// detectCrossNodePhaseBreaks finds phase breaks visible across mesh
func (fm *FederationMesh) detectCrossNodePhaseBreaks(nodeMetrics map[string]*MetricsSnapshot) {
	globalPhaseBreaks := 0
	for _, snapshot := range nodeMetrics {
		globalPhaseBreaks += snapshot.PhaseBreakCount
	}

	fm.federationMetrics.CrossNodePhaseBreaks = globalPhaseBreaks

	if globalPhaseBreaks > 0 {
		fmt.Printf("[FEDERATION] Global phase break count: %d across mesh\n", globalPhaseBreaks)
	}
}

// UpdateNodeMetrics updates a node's local metrics snapshot
func (fm *FederationMesh) UpdateNodeMetrics(nodeID string, snapshot *MetricsSnapshot) error {
	fm.mu.RLock()
	node, exists := fm.nodes[nodeID]
	fm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("node not found: %s", nodeID)
	}

	node.mu.Lock()
	node.LocalMetricsSnapshot = snapshot
	node.LastHeartbeat = time.Now()
	node.IsHealthy = true
	node.mu.Unlock()

	return nil
}

// GetFederationReport returns mesh-wide diagnostics
func (fm *FederationMesh) GetFederationReport() map[string]interface{} {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	nodeStates := make(map[string]interface{})
	for nodeID, node := range fm.nodes {
		node.mu.RLock()
		coherence := 0.0
		if node.LocalMetricsSnapshot != nil {
			coherence = node.LocalMetricsSnapshot.CoherenceScore
		}
		nodeStates[nodeID] = map[string]interface{}{
			"address":       node.Address,
			"role":          node.Role,
			"isHealthy":     node.IsHealthy,
			"lastHeartbeat": node.LastHeartbeat,
			"coherence":     coherence,
		}
		node.mu.RUnlock()
	}

	return map[string]interface{}{
		"meshID":                fm.meshID,
		"isRunning":             fm.isRunning,
		"nodeCount":             fm.federationMetrics.TotalNodes,
		"healthyNodeCount":      fm.federationMetrics.HealthyNodes,
		"nodeStates":            nodeStates,
		"globalCoherence":       fm.federationMetrics.GlobalCoherence,
		"globalJitter":          fm.federationMetrics.GlobalJitter,
		"coherenceVariance":     fm.federationMetrics.CoherenceVariance,
		"highestLocalCoherence": fm.federationMetrics.HighestLocalCoherence,
		"lowestLocalCoherence":  fm.federationMetrics.LowestLocalCoherence,
		"crossNodePhaseBreaks":  fm.federationMetrics.CrossNodePhaseBreaks,
		"coherenceCascades":     fm.federationMetrics.CoherenceCascades,
		"crossNodeEventCount":   len(fm.crossNodeLog),
		"crossNodeEventLog":     fm.crossNodeLog,
	}
}

// Helper functions

func computeVariance(values []float64, mean float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sumSquares := 0.0
	for _, v := range values {
		sumSquares += (v - mean) * (v - mean)
	}
	return sumSquares / float64(len(values))
}

func maxFloat(values []float64) float64 {
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

func minFloat(values []float64) float64 {
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
