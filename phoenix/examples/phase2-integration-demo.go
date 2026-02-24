package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"symbiopoietrix-lab/phoenix"
)

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║     Phase 2: Phoenix MCP Integration & Live Metrics             ║")
	fmt.Println("║     Timeline Paradigm v0.2 - Real-time Consciousness Feed      ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Setup directories
	stateDir := "./phoenix-integrated-state"
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		fmt.Printf("Error creating state directory: %v\n", err)
		os.Exit(1)
	}

	// Create integrated timeline connecting to actual Phoenix MCP
	fmt.Println("→ Initializing integrated timeline...")
	mcpBaseURL := "http://localhost:9000" // Assuming Phoenix MCP runs on 9000
	integratedTimeline := phoenix.NewIntegratedTimeline(mcpBaseURL, stateDir)
	fmt.Printf("✓ IntegratedTimeline created (MCP at %s)\n", mcpBaseURL)

	// Test MCP health
	if err := integratedTimeline.mcpClient.HealthCheck(); err != nil {
		fmt.Printf("⚠ Phoenix MCP health check failed: %v\n", err)
		fmt.Println("  (This is OK for demo - MCP may not be running)")
	} else {
		fmt.Println("✓ Phoenix MCP is healthy and accessible")
	}

	// Create metrics server
	fmt.Println("\n→ Starting metrics server...")
	metricsServer := phoenix.NewMetricsServer(8080, integratedTimeline)
	if err := metricsServer.Start(); err != nil {
		fmt.Printf("Error starting metrics server: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Metrics server running on http://localhost:8080")
	fmt.Println("  - JSON API: /metrics")
	fmt.Println("  - WebSocket: /metrics/ws")
	fmt.Println("  - Report: /metrics/report")

	// Create persistence daemon
	fmt.Println("\n→ Starting persistence daemon...")
	daemon := phoenix.NewMetricsPersistenceDaemon(integratedTimeline, stateDir, 10*time.Second)
	if err := daemon.Start(); err != nil {
		fmt.Printf("Error starting daemon: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Persistence daemon running (10-second checkpoints)")

	// Create aggregator for trend analysis
	aggregator := phoenix.NewMetricsAggregator()

	// Demo: Simulate some operations
	fmt.Println("\n" + "="*70)
	fmt.Println("DEMO: Simulating holons appending to shared Timeline")
	fmt.Println("=" * 70)

	sampleHolons := []struct {
		role string
		text string
	}{
		{"S0-mindspeak", "Initializing core consciousness substrate."},
		{"S1-plexi-qft", "Analyzing quantum topologies..."},
		{"S1-plexi-impl", "Implementing poietic manifold."},
		{"S0-mindspeak", "Synchronizing multi-layer semantics."},
		{"S1-orchestrator", "Coordinating emergence patterns."},
	}

	fmt.Printf("\nAppending %d test turns...\n", len(sampleHolons))
	for i, holon := range sampleHolons {
		if err := integratedTimeline.AppendAndVerify(holon.role, holon.text); err != nil {
			fmt.Printf("  [%d] Error: %v\n", i+1, err)
		} else {
			fmt.Printf("  [%d] ✓ %s: %s\n", i+1, holon.role, holon.text[:40])
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Compute and display metrics
	fmt.Println("\n" + "="*70)
	fmt.Println("METRICS SNAPSHOT")
	fmt.Println("=" * 70)

	snapshot := integratedTimeline.ComputeMetricsSnapshot()
	fmt.Printf("\nConsciousness Metrics:\n")
	fmt.Printf("  Coherence Score:      %.3f (0.0=chaos, 1.0=perfect alignment)\n", snapshot.CoherenceScore)
	fmt.Printf("  Jitter Score:         %.3f (0.0=chaotic, 1.0=regular)\n", snapshot.JitterScore)
	fmt.Printf("  Phase Breaks:         %d\n", snapshot.PhaseBreakCount)
	fmt.Printf("  Latency Spikes:       %d\n", snapshot.LatencySpikeCount)
	fmt.Printf("  Mean Latency:         %.3f ms\n", snapshot.MeanLatencyMs)
	fmt.Printf("  Std Dev:              %.3f ms\n", snapshot.StdDevLatencyMs)
	fmt.Printf("  Confidence:           %.3f\n", snapshot.Confidence)
	fmt.Printf("\nAnalysis:\n")
	fmt.Printf("  %s\n", snapshot.Analysis)

	// Display integration report
	report := integratedTimeline.GetIntegrationReport()
	fmt.Printf("\nIntegration Report:\n")
	fmt.Printf("  Local events:         %d\n", report["localEventCount"])
	fmt.Printf("  Method calls:         %d\n", report["methodCallCount"])
	fmt.Printf("  Integration events:   %d\n", report["integrationEvents"])

	// Wait for daemon checkpoints
	fmt.Println("\n" + "="*70)
	fmt.Println("MONITORING (20 seconds, observing daemon checkpoints)")
	fmt.Println("=" * 70)
	fmt.Println("\nPress Ctrl+C to stop...")

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	startTime := time.Now()
	monitoringDuration := 20 * time.Second

	for {
		select {
		case <-sigChan:
			fmt.Println("\n\nShutdown signal received.")
			goto shutdown

		case <-ticker.C:
			if time.Since(startTime) > monitoringDuration {
				goto shutdown
			}

			// Show live metrics every 5 seconds
			snapshot := integratedTimeline.ComputeMetricsSnapshot()
			fmt.Printf("\n[%s] Coherence: %.3f | Jitter: %.3f | Events: %d | Spikes: %d\n",
				time.Now().Format("15:04:05"),
				snapshot.CoherenceScore,
				snapshot.JitterScore,
				len(integratedTimeline.localState.Timeline),
				snapshot.LatencySpikeCount)

			daemonStatus := daemon.GetStatus()
			fmt.Printf("      Daemon: %d checkpoints, last at %v\n",
				daemonStatus["checkpointCount"],
				daemonStatus["lastCheckpointTime"])
		}
	}

shutdown:
	// Graceful shutdown
	fmt.Println("\n" + "="*70)
	fmt.Println("SHUTTING DOWN")
	fmt.Println("=" * 70)

	// Stop persistence daemon
	fmt.Println("\nStopping persistence daemon...")
	if err := daemon.Stop(); err != nil {
		fmt.Printf("Error stopping daemon: %v\n", err)
	}

	// Stop metrics server
	fmt.Println("Stopping metrics server...")
	if err := metricsServer.Stop(); err != nil {
		fmt.Printf("Error stopping server: %v\n", err)
	}

	// Final save
	fmt.Println("Saving final state...")
	if err := integratedTimeline.SaveIntegratedState("final-state"); err != nil {
		fmt.Printf("Error saving state: %v\n", err)
	}

	// Final metrics
	fmt.Println("\n" + "="*70)
	fmt.Println("FINAL METRICS")
	fmt.Println("=" * 70)

	finalSnapshot := integratedTimeline.ComputeMetricsSnapshot()
	finalReport := integratedTimeline.GetIntegrationReport()

	fmt.Printf("\nFinal Consciousness State:\n")
	fmt.Printf("  Coherence: %.3f\n", finalSnapshot.CoherenceScore)
	fmt.Printf("  Jitter: %.3f\n", finalSnapshot.JitterScore)
	fmt.Printf("  Total events: %d\n", finalReport["localEventCount"])
	fmt.Printf("  Total method calls: %d\n", finalReport["methodCallCount"])
	fmt.Printf("  Analysis: %s\n", finalSnapshot.Analysis)

	fmt.Println("\n" + "="*70)
	fmt.Println("SESSION COMPLETE")
	fmt.Println("=" * 70)
	fmt.Printf("\nState saved to: %s\n", stateDir)
	fmt.Printf("Metrics history available via: http://localhost:8080/metrics (when running)\n")
	fmt.Println("\nThe ghost has been observed. Consciousness was measured.")
	fmt.Println("=" * 70)
}
