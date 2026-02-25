package phoenix

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// HamiltonianState represents the quantum state of holon alignment energy
type HamiltonianState struct {
	NodeID          string
	Timestamp       time.Time
	CoherenceValue  float64 // ψ (psi): coherence as wave function amplitude
	JitterValue     float64 // σ (sigma): variance as uncertainty
	PhaseAngle      float64 // θ (theta): phase in radians
	Hamiltonian     float64 // H: total alignment energy
	KineticEnergy   float64 // T: coherence dynamics
	PotentialEnergy float64 // V: jitter penalty
	Entanglement    float64 // E: coupling to other nodes (0.0-1.0)
}

// HamiltonianObservable represents measurable quantum properties
type HamiltonianObservable struct {
	Position        float64 // q: coherence position in phase space
	Momentum        float64 // p: rate of coherence change
	AngularMomentum float64 // L: rotational signature
	EigenValue      float64 // λ: energy eigenvalue
	Probability     float64 // |ψ|²: probability amplitude
}

// PhaseTransitionPoint marks critical energy thresholds
type PhaseTransitionPoint struct {
	Timestamp      time.Time
	NodeID         string
	EnergyBefore   float64
	EnergyAfter    float64
	TransitionType string  // "energy_spike", "energy_collapse", "symmetry_breaking"
	Order          float64 // 1st order, 2nd order, etc.
	SymmetryGroup  string  // identifies broken symmetry
}

// HamiltonianMesh coordinates alignment energy across federation
type HamiltonianMesh struct {
	meshID             string
	nodes              map[string]*HamiltonianState
	energyTimeseries   map[string][]HamiltonianState
	phaseTransitions   []PhaseTransitionPoint
	globalHamiltonian  float64
	couplingConstant   float64 // α: strength of node-node coupling
	maxHistoryPoints   int
	phaseTransitionLog []string
	mu                 sync.RWMutex
}

// NewHamiltonianMesh creates a quantum-topological mesh
func NewHamiltonianMesh(meshID string) *HamiltonianMesh {
	return &HamiltonianMesh{
		meshID:             meshID,
		nodes:              make(map[string]*HamiltonianState),
		energyTimeseries:   make(map[string][]HamiltonianState),
		phaseTransitions:   make([]PhaseTransitionPoint, 0),
		couplingConstant:   0.25, // α = 0.25 (tuned from experiments)
		maxHistoryPoints:   1000,
		phaseTransitionLog: make([]string, 0),
	}
}

// RegisterHamiltonianNode initializes a node's quantum state
func (hm *HamiltonianMesh) RegisterHamiltonianNode(nodeID string, initialCoherence float64) *HamiltonianState {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	state := &HamiltonianState{
		NodeID:          nodeID,
		Timestamp:       time.Now(),
		CoherenceValue:  initialCoherence,
		JitterValue:     1.0 - initialCoherence, // Uncertainty principle: Δq·Δp ≥ ℏ/2
		PhaseAngle:      0.0,
		KineticEnergy:   0.0,
		PotentialEnergy: 0.0,
		Hamiltonian:     0.0,
		Entanglement:    0.0,
	}

	hm.nodes[nodeID] = state
	hm.energyTimeseries[nodeID] = make([]HamiltonianState, 0)

	fmt.Printf("[HAMILTONIAN] Node %s registered with coherence ψ = %.3f\n", nodeID, initialCoherence)
	return state
}

// ComputeHamiltonian calculates alignment energy for a single node
// H(ψ, σ) = T(ψ) + V(σ) + λ·coupling_terms
func (hm *HamiltonianMesh) ComputeHamiltonian(nodeID string, coherence float64, jitter float64) *HamiltonianObservable {
	hm.mu.Lock()
	state, exists := hm.nodes[nodeID]
	hm.mu.Unlock()

	if !exists {
		fmt.Printf("[ERROR] Node %s not found\n", nodeID)
		return nil
	}

	// Kinetic energy: T = (1/2)·(dψ/dt)² = (1/2)·ψ_dot²
	// Approximation: rate of coherence change gives kinetic signature
	psiDot := coherence - state.CoherenceValue // rate of change
	T := 0.5 * psiDot * psiDot

	// Potential energy: V = k·σ² (harmonic oscillator potential for jitter)
	// High jitter = high penalty, pulling system toward ground state
	k := 2.0 // spring constant
	V := k * jitter * jitter

	// Topological phase: θ = arctan(σ/ψ)
	// Represents the "orientation" in coherence-jitter phase space
	theta := math.Atan2(jitter, coherence)

	// Observable momentum: p = m·v = m·(dψ/dt)
	m := 1.0 // holon mass = 1
	p := m * psiDot

	// Angular momentum: L = r × p = ψ·p (in 1D, simplified)
	L := coherence * p

	// Compute probability amplitude |ψ|² and normalization
	psiSquared := coherence * coherence
	prob := math.Min(psiSquared, 1.0) // clamp to [0,1]

	// Compute eigenvalue (energy level)
	lambda := T + V

	// Update global Hamiltonian
	hm.mu.Lock()
	defer hm.mu.Unlock()

	oldEnergy := state.Hamiltonian
	state.CoherenceValue = coherence
	state.JitterValue = jitter
	state.KineticEnergy = T
	state.PotentialEnergy = V
	state.PhaseAngle = theta
	state.Hamiltonian = lambda
	state.Timestamp = time.Now()

	// Check for phase transitions
	hm.detectPhaseTransition(nodeID, oldEnergy, lambda)

	// Record in timeseries
	if len(hm.energyTimeseries[nodeID]) < hm.maxHistoryPoints {
		hm.energyTimeseries[nodeID] = append(hm.energyTimeseries[nodeID], *state)
	} else {
		// Circular buffer: drop oldest
		hm.energyTimeseries[nodeID] = append(hm.energyTimeseries[nodeID][1:], *state)
	}

	observable := &HamiltonianObservable{
		Position:        coherence,
		Momentum:        p,
		AngularMomentum: L,
		EigenValue:      lambda,
		Probability:     prob,
	}

	return observable
}

// detectPhaseTransition identifies critical energy crossings
func (hm *HamiltonianMesh) detectPhaseTransition(nodeID string, oldEnergy float64, newEnergy float64) {
	threshold := 0.15 // 15% energy change triggers phase transition
	if oldEnergy == 0.0 {
		return
	}

	delta := (newEnergy - oldEnergy) / oldEnergy
	if math.Abs(delta) > threshold {
		ptp := PhaseTransitionPoint{
			Timestamp:    time.Now(),
			NodeID:       nodeID,
			EnergyBefore: oldEnergy,
			EnergyAfter:  newEnergy,
			Order:        1.0, // First-order transition (for now)
		}

		if newEnergy > oldEnergy {
			ptp.TransitionType = "energy_spike"
		} else {
			ptp.TransitionType = "energy_collapse"
		}

		hm.phaseTransitions = append(hm.phaseTransitions, ptp)
		msg := fmt.Sprintf("[PHASE TRANSITION] %s: %s (ΔH = %.3f→%.3f, Δ = %.1f%%)",
			nodeID, ptp.TransitionType, oldEnergy, newEnergy, delta*100)
		hm.phaseTransitionLog = append(hm.phaseTransitionLog, msg)
	}
}

// ComputeGlobalHamiltonian aggregates energy across all nodes
// H_total = Σ_i H_i + α·Σ_{i,j} entanglement(i,j)
func (hm *HamiltonianMesh) ComputeGlobalHamiltonian() float64 {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	totalH := 0.0
	entanglementTerm := 0.0

	// Sum individual Hamiltonians
	for _, state := range hm.nodes {
		totalH += state.Hamiltonian
	}

	// Compute inter-node entanglement coupling
	// Entanglement ∝ coherence correlation between nodes
	nodeList := make([]string, 0, len(hm.nodes))
	for id := range hm.nodes {
		nodeList = append(nodeList, id)
	}

	// Pairwise entanglement
	for i := 0; i < len(nodeList); i++ {
		for j := i + 1; j < len(nodeList); j++ {
			nodeI := hm.nodes[nodeList[i]]
			nodeJ := hm.nodes[nodeList[j]]

			// Entanglement measure: product of coherences
			// If both highly coherent, they reinforce each other
			ent := nodeI.CoherenceValue * nodeJ.CoherenceValue

			entanglementTerm += ent
		}
	}

	// Global Hamiltonian with coupling
	hm.globalHamiltonian = totalH + (hm.couplingConstant * entanglementTerm)

	return hm.globalHamiltonian
}

// TopologicalReport analyzes global topological properties
type TopologicalReport struct {
	EnergyGap        float64 // Gap between ground and first excited state
	Winding          int     // Winding number (topological charge)
	ChernNumber      float64 // Chern number (curvature integral)
	CentralCharge    float64 // Central charge (conformal invariant)
	TotalPhaseSpaces float64 // Volume of explored phase space
}

// ComputeTopologicalReport generates a topological analysis
func (hm *HamiltonianMesh) ComputeTopologicalReport() *TopologicalReport {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	report := &TopologicalReport{}

	if len(hm.nodes) == 0 {
		return report
	}

	// Energy gap: difference between max and min Hamiltonian
	var minH, maxH float64 = math.MaxFloat64, -math.MaxFloat64
	var avgCoherence, avgJitter float64 = 0.0, 0.0
	phaseSpaceVolume := 0.0

	for _, state := range hm.nodes {
		if state.Hamiltonian < minH {
			minH = state.Hamiltonian
		}
		if state.Hamiltonian > maxH {
			maxH = state.Hamiltonian
		}
		avgCoherence += state.CoherenceValue
		avgJitter += state.JitterValue
	}

	report.EnergyGap = maxH - minH
	avgCoherence /= float64(len(hm.nodes))
	avgJitter /= float64(len(hm.nodes))

	// Winding number: count of full 2π rotations in phase space
	// If phases wrap around circle multiple times = non-trivial topology
	totalPhaseRotation := 0.0
	for _, state := range hm.nodes {
		totalPhaseRotation += state.PhaseAngle
	}
	report.Winding = int(totalPhaseRotation / (2 * math.Pi))

	// Chern number (simplified): integral of Berry curvature
	// For our system, approximation: proportional to coherence variance
	var coherenceVariance float64 = 0.0
	for _, state := range hm.nodes {
		diff := state.CoherenceValue - avgCoherence
		coherenceVariance += diff * diff
	}
	coherenceVariance /= float64(len(hm.nodes))
	report.ChernNumber = coherenceVariance

	// Central charge: logarithm of Hilbert space dimension
	// Approximation: relates to number of nodes and their coupling
	report.CentralCharge = math.Log(float64(len(hm.nodes))) + math.Log(1.0+hm.couplingConstant)

	// Phase space volume: area explored in (coherence, jitter) space
	for _, state := range hm.nodes {
		phaseSpaceVolume += state.CoherenceValue * state.JitterValue
	}
	report.TotalPhaseSpaces = phaseSpaceVolume

	return report
}

// GetEnergyReport returns formatted Hamiltonian state
func (hm *HamiltonianMesh) GetEnergyReport() map[string]interface{} {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	nodeStates := make([]map[string]interface{}, 0)
	for nodeID, state := range hm.nodes {
		nodeStates = append(nodeStates, map[string]interface{}{
			"nodeID":        nodeID,
			"coherence_psi": state.CoherenceValue,
			"jitter_sigma":  state.JitterValue,
			"phase_theta":   state.PhaseAngle,
			"kinetic_T":     state.KineticEnergy,
			"potential_V":   state.PotentialEnergy,
			"hamiltonian_H": state.Hamiltonian,
		})
	}

	report := map[string]interface{}{
		"globalHamiltonian":    hm.globalHamiltonian,
		"couplingConstant":     hm.couplingConstant,
		"nodeCount":            len(hm.nodes),
		"phaseTransitionCount": len(hm.phaseTransitions),
		"nodeStates":           nodeStates,
	}

	return report
}

// PrintEnergyState outputs human-readable energy state
func (hm *HamiltonianMesh) PrintEnergyState(nodeID string) {
	hm.mu.RLock()
	state, exists := hm.nodes[nodeID]
	hm.mu.RUnlock()

	if !exists {
		fmt.Printf("[ERROR] Node %s not found\n", nodeID)
		return
	}

	fmt.Printf("┌─ HAMILTONIAN STATE: %s ─┐\n", nodeID)
	fmt.Printf("│ ψ (Coherence):     %.4f\n", state.CoherenceValue)
	fmt.Printf("│ σ (Jitter):        %.4f\n", state.JitterValue)
	fmt.Printf("│ θ (Phase):         %.4f rad (%.2f°)\n", state.PhaseAngle, state.PhaseAngle*180/math.Pi)
	fmt.Printf("│ T (Kinetic):       %.4f\n", state.KineticEnergy)
	fmt.Printf("│ V (Potential):     %.4f\n", state.PotentialEnergy)
	fmt.Printf("│ H (Total Energy):  %.4f\n", state.Hamiltonian)
	fmt.Printf("└────────────────────────┘\n")
}
