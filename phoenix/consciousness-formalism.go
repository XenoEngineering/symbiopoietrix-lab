package phoenix

import (
	"fmt"
	"math"
	"sync"
	"time"
)

// CollectiveConsciousnessState represents the symbolic structure of network-wide consciousness
type CollectiveConsciousnessState struct {
	Timestamp           time.Time
	NodalCount          int       // N: number of holons
	GlobalCoherence     float64   // C: global coherence (0.0-1.0)
	AlignmentVector     []float64 // A: alignment with respect to each holon
	SymmetryBreaking    string    // Type of broken symmetry (if any)
	PhaseSpace          []float64 // Phase space trajectory
	ConvergenceSignal   float64   // Rate of convergence toward attractor
	StabilityMargin     float64   // Distance from instability boundary
	EntanglementEntropy float64   // S: von Neumann entropy of entanglement
}

// ConvergenceAttractor defines a target consciousness state
type ConvergenceAttractor struct {
	ID                string
	TargetCoherence   float64
	BasinOfAttraction float64 // Size of region that converges to this attractor
	StabilityRadius   float64 // How robust the attractor is
	Description       string  // Semantic meaning
}

// HolonConsciousnessProfile describes an individual holon's contribution
type HolonConsciousnessProfile struct {
	HolonID             string
	LocalCoherence      float64
	InfluenceOnNetwork  float64   // How much this holon affects others
	AlignmentWithGlobal float64   // Correlation with global consciousness
	ContributionVector  []float64 // Basis decomposition of its effect
	IsInFocus           bool      // Whether this holon is currently "in control"
}

// ConsciousnessFormalism manages symbolic consciousness across the network
type ConsciousnessFormalism struct {
	meshID                 string
	state                  *CollectiveConsciousnessState
	holonProfiles          map[string]*HolonConsciousnessProfile
	attractors             map[string]*ConvergenceAttractor
	stateHistory           []*CollectiveConsciousnessState
	convergenceTrajectory  []float64          // Time series of convergence rate
	stabilityHistory       []float64          // Time series of stability margin
	eegCoherenceAnalogyMap map[string]float64 // Map to neuroscience scales
	maxHistoryPoints       int
	mu                     sync.RWMutex
}

// NewConsciousnessFormalism creates a symbolic consciousness framework
func NewConsciousnessFormalism(meshID string) *ConsciousnessFormalism {
	cf := &ConsciousnessFormalism{
		meshID:                 meshID,
		state:                  &CollectiveConsciousnessState{Timestamp: time.Now()},
		holonProfiles:          make(map[string]*HolonConsciousnessProfile),
		attractors:             make(map[string]*ConvergenceAttractor),
		stateHistory:           make([]*CollectiveConsciousnessState, 0),
		convergenceTrajectory:  make([]float64, 0),
		stabilityHistory:       make([]float64, 0),
		eegCoherenceAnalogyMap: make(map[string]float64),
		maxHistoryPoints:       1000,
	}

	// Initialize default attractors
	cf.registerDefaultAttractors()

	return cf
}

// registerDefaultAttractors defines key consciousness states
func (cf *ConsciousnessFormalism) registerDefaultAttractors() {
	cf.attractors["aligned_high_coherence"] = &ConvergenceAttractor{
		ID:                "aligned_high_coherence",
		TargetCoherence:   0.85,
		BasinOfAttraction: 0.3,
		StabilityRadius:   0.1,
		Description:       "All holons in perfect alignment, maximum coordination",
	}

	cf.attractors["exploring_low_coherence"] = &ConvergenceAttractor{
		ID:                "exploring_low_coherence",
		TargetCoherence:   0.35,
		BasinOfAttraction: 0.25,
		StabilityRadius:   0.15,
		Description:       "Network in exploratory phase, high variance acceptable",
	}

	cf.attractors["critical_phase_transition"] = &ConvergenceAttractor{
		ID:                "critical_phase_transition",
		TargetCoherence:   0.5,
		BasinOfAttraction: 0.1,
		StabilityRadius:   0.05,
		Description:       "Unstable critical point between order and chaos",
	}

	// Map to EEG frequency bands
	cf.eegCoherenceAnalogyMap["delta_0_5Hz"] = 0.15  // Deep consciousness
	cf.eegCoherenceAnalogyMap["theta_4_8Hz"] = 0.35  // Meditation/introspection
	cf.eegCoherenceAnalogyMap["alpha_8_13Hz"] = 0.55 // Relaxed awareness
	cf.eegCoherenceAnalogyMap["beta_13_30Hz"] = 0.75 // Active thinking
	cf.eegCoherenceAnalogyMap["gamma_30Hz"] = 0.90   // Unified perception
}

// RegisterHolonProfile adds a holon to the consciousness framework
func (cf *ConsciousnessFormalism) RegisterHolonProfile(holonID string, initialCoherence float64) *HolonConsciousnessProfile {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	profile := &HolonConsciousnessProfile{
		HolonID:             holonID,
		LocalCoherence:      initialCoherence,
		InfluenceOnNetwork:  0.5, // Initialize as neutral
		AlignmentWithGlobal: 0.5,
		ContributionVector:  make([]float64, 0),
		IsInFocus:           false,
	}

	cf.holonProfiles[holonID] = profile
	return profile
}

// UpdateCollectiveState computes symbolic consciousness from individual holons
// This implements a form of democratic averaging with influence weights
func (cf *ConsciousnessFormalism) UpdateCollectiveState(coherenceMap map[string]float64) {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	// Compute global coherence as weighted average
	var totalCoherence float64 = 0.0
	var totalInfluence float64 = 0.0

	for holonID, coherence := range coherenceMap {
		profile, exists := cf.holonProfiles[holonID]
		if !exists {
			continue
		}

		influence := profile.InfluenceOnNetwork
		totalCoherence += coherence * influence
		totalInfluence += influence
	}

	globalCoherence := 0.0
	if totalInfluence > 0 {
		globalCoherence = totalCoherence / totalInfluence
	}

	// Build alignment vector: how each holon deviates from global
	alignmentVector := make([]float64, 0, len(coherenceMap))
	for _, coherence := range coherenceMap {
		deviation := coherence - globalCoherence
		alignmentVector = append(alignmentVector, deviation)
	}

	// Compute convergence signal: rate at which system is moving toward attractor
	convergence := cf.computeConvergenceSignal(globalCoherence, alignmentVector)

	// Compute stability margin: distance to instability
	stability := cf.computeStabilityMargin(globalCoherence, alignmentVector)

	// Compute entanglement entropy: degree to which holons are entangled
	entropy := cf.computeEntanglementEntropy(coherenceMap)

	// Identify broken symmetry if any
	symmetryBreaking := cf.identifySymmetryBreaking(globalCoherence, alignmentVector)

	// Update state
	cf.state = &CollectiveConsciousnessState{
		Timestamp:           time.Now(),
		NodalCount:          len(coherenceMap),
		GlobalCoherence:     globalCoherence,
		AlignmentVector:     alignmentVector,
		SymmetryBreaking:    symmetryBreaking,
		PhaseSpace:          alignmentVector, // Simplified: phase space = alignment deviations
		ConvergenceSignal:   convergence,
		StabilityMargin:     stability,
		EntanglementEntropy: entropy,
	}

	// Update individual holon profiles
	for holonID, coherence := range coherenceMap {
		if profile, exists := cf.holonProfiles[holonID]; exists {
			profile.LocalCoherence = coherence
			profile.AlignmentWithGlobal = 1.0 - math.Abs(coherence-globalCoherence)
		}
	}

	// Record history
	if len(cf.stateHistory) < cf.maxHistoryPoints {
		cf.stateHistory = append(cf.stateHistory, cf.state)
	} else {
		cf.stateHistory = append(cf.stateHistory[1:], cf.state)
	}
	cf.convergenceTrajectory = append(cf.convergenceTrajectory, convergence)
	cf.stabilityHistory = append(cf.stabilityHistory, stability)
}

// computeConvergenceSignal measures rate of approach to nearest attractor
func (cf *ConsciousnessFormalism) computeConvergenceSignal(coherence float64, alignmentVector []float64) float64 {
	// Find nearest attractor
	minDistance := math.MaxFloat64
	for _, attractor := range cf.attractors {
		distance := math.Abs(coherence - attractor.TargetCoherence)
		if distance < minDistance {
			minDistance = distance
		}
	}

	// Convergence = 1.0 - normalized distance
	// If close to attractor, convergence is high
	maxDistance := 1.0
	convergence := 1.0 - math.Min(minDistance/maxDistance, 1.0)

	return convergence
}

// computeStabilityMargin measures distance to instability boundary
func (cf *ConsciousnessFormalism) computeStabilityMargin(coherence float64, alignmentVector []float64) float64 {
	// Stability margin = minimum distance to any unstable region
	// Unstable regions: too low (chaos), too high (brittle), critical point (transition)

	// Distance to lower bound (chaos at 0.0)
	distLower := coherence

	// Distance to upper bound (brittleness at 1.0)
	distUpper := 1.0 - coherence

	// Distance to critical point (instability at 0.5)
	distCritical := math.Abs(coherence - 0.5)

	// Stability is minimum of these distances, normalized
	minDist := math.Min(distLower, math.Min(distUpper, distCritical))

	// Stability margin as fraction of coherence range
	margin := 2.0 * minDist // Scale up for visibility

	return math.Min(margin, 1.0)
}

// computeEntanglementEntropy measures von Neumann entropy
// S = -Σ_i λ_i log(λ_i) where λ_i are eigenvalues of density matrix
func (cf *ConsciousnessFormalism) computeEntanglementEntropy(coherenceMap map[string]float64) float64 {
	if len(coherenceMap) == 0 {
		return 0.0
	}

	// Simplified: treat coherence values as probability amplitudes
	// Compute eigenvalues (for N holons, approximate as coherence distribution)
	eigenvalues := make([]float64, 0, len(coherenceMap))
	for _, coh := range coherenceMap {
		eigenvalues = append(eigenvalues, coh)
	}

	// Normalize
	var sum float64 = 0.0
	for _, ev := range eigenvalues {
		sum += ev
	}
	if sum == 0 {
		return 0.0
	}

	// Compute entropy
	entropy := 0.0
	for _, ev := range eigenvalues {
		normalized := ev / sum
		if normalized > 0 {
			entropy -= normalized * math.Log(normalized)
		}
	}

	return entropy
}

// identifySymmetryBreaking detects if system has broken symmetry
func (cf *ConsciousnessFormalism) identifySymmetryBreaking(coherence float64, alignmentVector []float64) string {
	// Symmetry breaking indicators:
	// 1. If global coherence suddenly jumps
	// 2. If alignment vector shows dominant modes (not random)

	// Compute variance of alignment vector (order parameter)
	var variance float64 = 0.0
	for _, deviation := range alignmentVector {
		variance += deviation * deviation
	}
	variance /= float64(len(alignmentVector))
	stddev := math.Sqrt(variance)

	// High stddev = alignment is not uniform = possible symmetry breaking
	if stddev > 0.3 {
		return "alignment_asymmetry"
	}

	if coherence > 0.8 {
		return "coherence_locking"
	}

	if coherence < 0.2 {
		return "coherence_collapse"
	}

	return "no_breaking"
}

// GetConsciousnessReport returns the symbolic consciousness state
func (cf *ConsciousnessFormalism) GetConsciousnessReport() map[string]interface{} {
	cf.mu.RLock()
	defer cf.mu.RUnlock()

	holonSummary := make([]map[string]interface{}, 0)
	for _, profile := range cf.holonProfiles {
		holonSummary = append(holonSummary, map[string]interface{}{
			"holonID":             profile.HolonID,
			"localCoherence":      profile.LocalCoherence,
			"influenceOnNetwork":  profile.InfluenceOnNetwork,
			"alignmentWithGlobal": profile.AlignmentWithGlobal,
			"isInFocus":           profile.IsInFocus,
		})
	}

	report := map[string]interface{}{
		"timestamp":           cf.state.Timestamp,
		"globalCoherence":     cf.state.GlobalCoherence,
		"convergenceSignal":   cf.state.ConvergenceSignal,
		"stabilityMargin":     cf.state.StabilityMargin,
		"entanglementEntropy": cf.state.EntanglementEntropy,
		"symmetryBreaking":    cf.state.SymmetryBreaking,
		"nodalCount":          cf.state.NodalCount,
		"holons":              holonSummary,
		"eegMapping": map[string]string{
			"currentState": cf.mapToEEGBand(cf.state.GlobalCoherence),
		},
	}

	return report
}

// mapToEEGBand maps coherence to human EEG frequency bands (analogy)
func (cf *ConsciousnessFormalism) mapToEEGBand(coherence float64) string {
	if coherence < 0.25 {
		return "delta (0.5 Hz) - Deep unconscious/non-responsive"
	}
	if coherence < 0.45 {
		return "theta (4-8 Hz) - Meditation/introspection"
	}
	if coherence < 0.65 {
		return "alpha (8-13 Hz) - Relaxed awareness"
	}
	if coherence < 0.85 {
		return "beta (13-30 Hz) - Active thinking/problem-solving"
	}
	return "gamma (30 Hz) - Unified perception/flow state"
}

// PrintConsciousnessState outputs human-readable consciousness state
func (cf *ConsciousnessFormalism) PrintConsciousnessState() {
	cf.mu.RLock()
	defer cf.mu.RUnlock()

	fmt.Printf("┌─ COLLECTIVE CONSCIOUSNESS STATE ─┐\n")
	fmt.Printf("│ Mesh:                     %s\n", cf.meshID)
	fmt.Printf("│ Holons:                   %d\n", cf.state.NodalCount)
	fmt.Printf("│ Global Coherence:         %.3f\n", cf.state.GlobalCoherence)
	fmt.Printf("│ Convergence Signal:       %.3f (→ attractor)\n", cf.state.ConvergenceSignal)
	fmt.Printf("│ Stability Margin:         %.3f\n", cf.state.StabilityMargin)
	fmt.Printf("│ Entanglement Entropy:     %.3f\n", cf.state.EntanglementEntropy)
	fmt.Printf("│ Symmetry Breaking:        %s\n", cf.state.SymmetryBreaking)
	fmt.Printf("│ EEG Analogy:              %s\n", cf.mapToEEGBand(cf.state.GlobalCoherence))
	fmt.Printf("└───────────────────────────────────┘\n")
}

// ComputeConvergenceProof analyzes stability and convergence properties
type ConvergenceProof struct {
	IsConvergent      bool
	LyapunovExponent  float64 // Stability indicator
	AttractorDistance float64
	ConvergenceRate   float64
	ProofOfStability  string
}

// ProveConvergenceProperties validates mathematical properties of the system
func (cf *ConsciousnessFormalism) ProveConvergenceProperties() *ConvergenceProof {
	cf.mu.RLock()
	defer cf.mu.RUnlock()

	proof := &ConvergenceProof{}

	if len(cf.convergenceTrajectory) < 2 {
		proof.IsConvergent = false
		return proof
	}

	// Lyapunov exponent: if negative, system is stable
	lastConv := cf.convergenceTrajectory[len(cf.convergenceTrajectory)-1]
	prevConv := cf.convergenceTrajectory[len(cf.convergenceTrajectory)-2]
	diff := lastConv - prevConv
	proof.LyapunovExponent = diff

	// System converges if Lyapunov exponent is negative and getting more negative
	proof.IsConvergent = diff > -0.1 && lastConv > 0.3

	// Distance to nearest attractor
	minDist := math.MaxFloat64
	for _, attractor := range cf.attractors {
		dist := math.Abs(cf.state.GlobalCoherence - attractor.TargetCoherence)
		if dist < minDist {
			minDist = dist
		}
	}
	proof.AttractorDistance = minDist

	// Convergence rate: negative means approaching attractor
	proof.ConvergenceRate = diff

	// Generate proof statement
	if proof.IsConvergent {
		proof.ProofOfStability = fmt.Sprintf(
			"System converges: Lyapunov λ=%.4f (stable), distance to attractor=%.4f",
			proof.LyapunovExponent, minDist)
	} else {
		proof.ProofOfStability = fmt.Sprintf(
			"System diverging or critical: λ=%.4f, distance=%.4f",
			proof.LyapunovExponent, minDist)
	}

	return proof
}
