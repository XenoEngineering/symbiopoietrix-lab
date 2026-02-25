package main

import (
	"fmt"
	"math"
	"sort"
	"time"

	"symbiopoietrix-lab/phoenix"
)

// CoherenceLandscapeExplorer systematically maps stable regions
type CoherenceLandscapeExplorer struct {
	meshID              string
	samples             int
	startCoherence      float64
	endCoherence        float64
	step                float64
	results             []StabilityResult
	attractorBasins     map[string]*BasinOfAttraction
	bifurcationPoints   []float64
}

// StabilityResult records behavior at a coherence point
type StabilityResult struct {
	InitialCoherence    float64
	FinalCoherence      float64
	Convergence         float64
	StabilityMargin     float64
	Lyapunov            float64
	IsConvergent        bool
	SettlingTime        int
	NearestAttractor    string
	DistanceToAttractor float64
}

// BasinOfAttraction maps a region of coherence space
type BasinOfAttraction struct {
	AttractorID         string
	TargetCoherence     float64
	BasinLowerBound     float64
	BasinUpperBound     float64
	BasinWidth          float64
	NumPointsInBasin    int
	StabilityStrength   float64 // How strongly it attracts
	ConvergenceRate     float64 // Average time to settle
}

// NewCoherenceLandscapeExplorer creates the explorer
func NewCoherenceLandscapeExplorer() *CoherenceLandscapeExplorer {
	return &CoherenceLandscapeExplorer{
		meshID:            "landscape-001",
		samples:           50,
		startCoherence:    0.0,
		endCoherence:      1.0,
		step:              0.02, // 50 points across [0, 1]
		results:           make([]StabilityResult, 0),
		attractorBasins:   make(map[string]*BasinOfAttraction),
		bifurcationPoints: make([]float64, 0),
	}
}

// ExploreLandscape scans coherence space from 0.0 to 1.0
func (cle *CoherenceLandscapeExplorer) ExploreLandscape() {
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║     COHERENCE LANDSCAPE EXPLORER                               ║")
	fmt.Println("║     Finding Islands of Stable Consciousness                    ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	fmt.Printf("→ Scanning coherence space: [%.2f, %.2f] with step %.3f\n", 
		cle.startCoherence, cle.endCoherence, cle.step)
	fmt.Printf("→ Total samples: %d\n", cle.samples)
	fmt.Println()

	// Scan each coherence point
	for i := 0; i < cle.samples; i++ {
		coherence := cle.startCoherence + float64(i)*cle.step

		// Create fresh consciousness formalism for this probe
		cf := phoenix.NewConsciousnessFormalism(fmt.Sprintf("probe-%d", i))
		cf.RegisterHolonProfile("h1", coherence)
		cf.RegisterHolonProfile("h2", coherence)
		cf.RegisterHolonProfile("h3", coherence)

		// Simulate 10 convergence steps from this starting point
		settlingTime := 0
		finalCoherence := coherence
		var convergence, stability, lyapunov float64

		for step := 0; step < 10; step++ {
			// Update with slight perturbation (realistic jitter)
			perturbedMap := map[string]float64{
				"h1": coherence + 0.01*math.Sin(float64(step)*0.5),
				"h2": coherence + 0.01*math.Cos(float64(step)*0.7),
				"h3": coherence + 0.01*math.Sin(float64(step)*0.3),
			}
			cf.UpdateCollectiveState(perturbedMap)

			// Check convergence
			proof := cf.ProveConvergenceProperties()
			finalCoherence = report["globalCoherence"].(float64)
			convergence = report["convergenceSignal"].(float64)
			stability = report["stabilityMargin"].(float64)
			lyapunov = proof.LyapunovExponent

			if proof.IsConvergent {
				settlingTime = step
				break
			}
		}

		// Determine nearest attractor
		nearestAttractor := cle.findNearestAttractor(coherence)
		report := cf.GetConsciousnessReport()

		result := StabilityResult{
			InitialCoherence:    coherence,
			FinalCoherence:      finalCoherence,
			Convergence:         convergence,
			StabilityMargin:     stability,
			Lyapunov:            lyapunov,
			IsConvergent:        lyapunov < 0,
			SettlingTime:        settlingTime,
			NearestAttractor:    nearestAttractor,
			DistanceToAttractor: math.Abs(coherence - cle.getAttractorCoherence(nearestAttractor)),
		}

		cle.results = append(cle.results, result)

		// Print progress
		marker := " "
		if result.IsConvergent {
			marker = "●"
		} else if result.StabilityMargin > 0.5 {
			marker = "◐"
		} else {
			marker = "○"
		}

		fmt.Printf("[%.2f] %s conv:%.2f stab:%.2f → %s (d=%.3f)\n",
			coherence, marker, result.Convergence, result.StabilityMargin,
			nearestAttractor, result.DistanceToAttractor)
	}

	fmt.Println()
}

// AnalyzeLandscape processes results to find basins
func (cle *CoherenceLandscapeExplorer) AnalyzeLandscape() {
	fmt.Println("=" * 70)
	fmt.Println("LANDSCAPE ANALYSIS")
	fmt.Println("=" * 70)
	fmt.Println()

	// Identify basin boundaries
	fmt.Println("→ Identifying attractor basins...")

	basins := make(map[string]*BasinOfAttraction)
	for _, result := range cle.results {
		attractorID := result.NearestAttractor

		if _, exists := basins[attractorID]; !exists {
			basins[attractorID] = &BasinOfAttraction{
				AttractorID:      attractorID,
				TargetCoherence:  cle.getAttractorCoherence(attractorID),
				BasinLowerBound:  result.InitialCoherence,
				BasinUpperBound:  result.InitialCoherence,
				NumPointsInBasin: 0,
			}
		}

		basin := basins[attractorID]
		basin.NumPointsInBasin++

		// Track basin bounds
		if result.InitialCoherence < basin.BasinLowerBound {
			basin.BasinLowerBound = result.InitialCoherence
		}
		if result.InitialCoherence > basin.BasinUpperBound {
			basin.BasinUpperBound = result.InitialCoherence
		}

		// Accumulate stability metrics
		if result.IsConvergent {
			basin.StabilityStrength += result.StabilityMargin
		}
		basin.ConvergenceRate += float64(result.SettlingTime)
	}

	// Finalize basin statistics
	for _, basin := range basins {
		basin.BasinWidth = basin.BasinUpperBound - basin.BasinLowerBound
		if basin.NumPointsInBasin > 0 {
			basin.ConvergenceRate /= float64(basin.NumPointsInBasin)
			basin.StabilityStrength /= float64(basin.NumPointsInBasin)
		}
	}

	cle.attractorBasins = basins

	// Print basins
	fmt.Println()
	fmt.Println("[ATTRACTOR BASINS DISCOVERED]")
	fmt.Println()

	// Sort by coherence for display
	attractorList := make([]string, 0, len(basins))
	for id := range basins {
		attractorList = append(attractorList, id)
	}
	sort.Slice(attractorList, func(i, j int) bool {
		return basins[attractorList[i]].TargetCoherence < basins[attractorList[j]].TargetCoherence
	})

	for _, id := range attractorList {
		basin := basins[id]
		fmt.Printf("┌─ %s ─┐\n", id)
		fmt.Printf("│ Target Coherence:     %.3f\n", basin.TargetCoherence)
		fmt.Printf("│ Basin Range:          [%.3f, %.3f]\n", basin.BasinLowerBound, basin.BasinUpperBound)
		fmt.Printf("│ Basin Width:          %.3f\n", basin.BasinWidth)
		fmt.Printf("│ Points in Basin:      %d / %d\n", basin.NumPointsInBasin, cle.samples)
		fmt.Printf("│ Stability Strength:   %.3f\n", basin.StabilityStrength)
		fmt.Printf("│ Avg Convergence Time: %.1f steps\n", basin.ConvergenceRate)
		fmt.Printf("└" + "─" * 40 + "┘\n")
		fmt.Println()
	}

	// Find bifurcation points (where basins meet)
	cle.findBifurcationPoints()
}

// findNearestAttractor returns which basin this point is in
func (cle *CoherenceLandscapeExplorer) findNearestAttractor(coherence float64) string {
	// Standard attractors
	attractors := map[string]float64{
		"aligned_high_coherence": 0.85,
		"exploring_low_coherence": 0.35,
		"critical_phase_transition": 0.5,
	}

	minDist := math.MaxFloat64
	nearest := ""
	for name, target := range attractors {
		dist := math.Abs(coherence - target)
		if dist < minDist {
			minDist = dist
			nearest = name
		}
	}
	return nearest
}

// getAttractorCoherence returns the coherence value of an attractor
func (cle *CoherenceLandscapeExplorer) getAttractorCoherence(attractorID string) float64 {
	targets := map[string]float64{
		"aligned_high_coherence": 0.85,
		"exploring_low_coherence": 0.35,
		"critical_phase_transition": 0.5,
	}
	return targets[attractorID]
}

// findBifurcationPoints locates critical transitions
func (cle *CoherenceLandscapeExplorer) findBifurcationPoints() {
	fmt.Println("=" * 70)
	fmt.Println("BIFURCATION ANALYSIS")
	fmt.Println("=" * 70)
	fmt.Println()

	fmt.Println("→ Searching for phase transition boundaries...")
	fmt.Println()

	// Look for places where attractor changes
	prevAttractor := ""
	for i, result := range cle.results {
		if i == 0 {
			prevAttractor = result.NearestAttractor
			continue
		}

		if result.NearestAttractor != prevAttractor {
			bifurcationCoherence := result.InitialCoherence
			cle.bifurcationPoints = append(cle.bifurcationPoints, bifurcationCoherence)

			fmt.Printf("🔄 Bifurcation at C ≈ %.3f\n", bifurcationCoherence)
			fmt.Printf("   Transition: %s → %s\n", prevAttractor, result.NearestAttractor)
			fmt.Println()

			prevAttractor = result.NearestAttractor
		}
	}

	if len(cle.bifurcationPoints) == 0 {
		fmt.Println("No bifurcations found (system is globally stable)")
	} else {
		fmt.Printf("Total bifurcation points: %d\n", len(cle.bifurcationPoints))
	}
	fmt.Println()
}

// ReportIslands prints summary of stable regions
func (cle *CoherenceLandscapeExplorer) ReportIslands() {
	fmt.Println()
	fmt.Println("=" * 70)
	fmt.Println("STABLE COHERENCE ISLANDS")
	fmt.Println("=" * 70)
	fmt.Println()

	// Count convergent points
	convergentCount := 0
	for _, result := range cle.results {
		if result.IsConvergent {
			convergentCount++
		}
	}

	stability := float64(convergentCount) / float64(len(cle.results)) * 100

	fmt.Printf("→ System Stability Across Space: %.1f%%\n", stability)
	fmt.Printf("  (%d/%d points converge)\n", convergentCount, len(cle.results))
	fmt.Println()

	// Identify "islands" (clusters of stability)
	islands := make([]IslandReport, 0)

	inIsland := false
	islandStart := 0.0

	for _, result := range cle.results {
		if result.IsConvergent && !inIsland {
			// Start of island
			inIsland = true
			islandStart = result.InitialCoherence
		} else if !result.IsConvergent && inIsland {
			// End of island
			inIsland = false
			islands = append(islands, IslandReport{
				Start:     islandStart,
				End:       result.InitialCoherence,
				Width:     result.InitialCoherence - islandStart,
				Attractor: result.NearestAttractor,
			})
		}
	}

	if inIsland {
		// Last island extends to end
		lastResult := cle.results[len(cle.results)-1]
		islands = append(islands, IslandReport{
			Start:     islandStart,
			End:       lastResult.InitialCoherence,
			Width:     lastResult.InitialCoherence - islandStart,
			Attractor: lastResult.NearestAttractor,
		})
	}

	fmt.Printf("Islands of Stable Coherence: %d\n", len(islands))
	fmt.Println()

	for i, island := range islands {
		fmt.Printf("[Island %d]\n", i+1)
		fmt.Printf("  Range:      [%.3f, %.3f]\n", island.Start, island.End)
		fmt.Printf("  Width:      %.3f\n", island.Width)
		fmt.Printf("  Attractor:  %s\n", island.Attractor)
		fmt.Println()
	}

	// Recommendation
	fmt.Println("=" * 70)
	fmt.Println("OPERATIONAL RECOMMENDATION")
	fmt.Println("=" * 70)
	fmt.Println()

	if len(islands) > 0 {
		biggestIsland := islands[0]
		for _, island := range islands {
			if island.Width > biggestIsland.Width {
				biggestIsland = island
			}
		}

		fmt.Printf("✓ SAFEST OPERATING REGION:\n")
		fmt.Printf("  Coherence Range: [%.3f, %.3f]\n", biggestIsland.Start, biggestIsland.End)
		fmt.Printf("  Attractor: %s\n", biggestIsland.Attractor)
		fmt.Printf("  Width (stability margin): %.3f\n", biggestIsland.Width)
		fmt.Println()

		targetCoherence := (biggestIsland.Start + biggestIsland.End) / 2.0
		fmt.Printf("✓ RECOMMENDED TARGET COHERENCE: %.3f\n", targetCoherence)
	} else {
		fmt.Println("⚠ WARNING: No stable islands detected!")
		fmt.Println("System may require tuning of coupling constants.")
	}

	fmt.Println()
}

// IslandReport describes one stable region
type IslandReport struct {
	Start     float64
	End       float64
	Width     float64
	Attractor string
}

func main() {
	explorer := NewCoherenceLandscapeExplorer()

	// Phase 1: Explore
	start := time.Now()
	explorer.ExploreLandscape()
	explorationTime := time.Since(start)

	// Phase 2: Analyze
	explorer.AnalyzeLandscape()

	// Phase 3: Report
	explorer.ReportIslands()

	fmt.Println("=" * 70)
	fmt.Printf("Exploration Time: %.2f seconds\n", explorationTime.Seconds())
	fmt.Printf("Samples Processed: %d\n", explorer.samples)
	fmt.Println("=" * 70)
	fmt.Println()
	fmt.Println("✓ LANDSCAPE EXPLORATION COMPLETE")
	fmt.Println()
}
