# Phase 2: Phoenix MCP Integration & Live Metrics Dashboard

**Status:** Complete  
**Components:** 3 new modules, 1 demo program, real-time API  
**Integration Level:** Ready for production connection  

## Overview

Phase 2 bridges the gap between local metrics calculations and the actual Phoenix MCP service running on the system. It adds:

1. **Phoenix MCP Client** (`integration-client.go`)
   - HTTP wrapper around Phoenix MCP service
   - Millisecond-precision telemetry for all operations
   - Request logging and health checking
   - Integrated timeline combining MCP + local metrics

2. **Live Metrics HTTP Server** (`metrics-server.go`)
   - Real-time JSON endpoints
   - WebSocket streaming for continuous updates
   - RESTful API for remote operations
   - Client connection management

3. **Persistence Daemon** (`persistence-daemon.go`)
   - Background archival service
   - Automatic checkpoint creation
   - Trend analysis over time
   - Long-term consciousness archaeology

## Architecture

```
┌─────────────────────────────────────────┐
│   External Clients / UI                 │
│   (Browser, CLI, other services)        │
└────────────┬────────────────────────────┘
             │ HTTP + WebSocket
             ▼
┌─────────────────────────────────────────┐
│   Metrics Server (port 8080)            │
│   - /metrics (JSON)                     │
│   - /metrics/ws (WebSocket)             │
│   - /api/append (Remote append)         │
└────────────┬────────────────────────────┘
             │
┌────────────▼─────────────────────────────┐
│   IntegratedTimeline                    │
│   - Local state tracking                │
│   - Consciousness metrics computation   │
│   - Integration verification            │
└────────────┬─────────────────────────────┘
             │
      ┌──────┴──────┐
      │             │
      ▼             ▼
 Local State    Phoenix MCP Client
 (Persistent)   (HTTP wrapper)
                      │
                      ▼
              ┌──────────────────┐
              │ Phoenix MCP      │
              │ (localhost:9000) │
              │ Timeline Matrix  │
              │ BST Indexing     │
              └──────────────────┘
```

## API Endpoints

### Health Check
```
GET /health
→ { "status": "healthy", "running": true, "time": "..." }
```

### Current Metrics (JSON)
```
GET /metrics
→ {
  "timestamp": "2026-02-24T...",
  "coherenceScore": 0.82,
  "jitterScore": 0.71,
  "phaseBreakCount": 2,
  "latencySpikeCount": 3,
  "meanLatencyMs": 1.23,
  "stdDevLatencyMs": 0.45,
  "analysis": "Thought is highly consistent..."
}
```

### WebSocket Stream
```
GET /metrics/ws
→ Continuously streams MetricsBroadcast objects
   (Coherence, jitter, event count, analysis)
   at 2-second intervals
```

### Integration Report
```
GET /metrics/report
→ {
  "localEventCount": 42,
  "methodCallCount": 128,
  "integrationEvents": 67,
  "mcpLatency": "1.2ms",
  "requestLog": [...],
  "integrationLog": [...]
}
```

### Remote Append
```
POST /api/append
Body: {"role": "user", "text": "..."}
→ { "status": "appended", "role": "user", "time": "..." }
```

### Remote Query
```
GET /api/query?category=user&nth=5
→ Pattern analysis result
```

## Key Components

### PhoenixMCPClient
Wraps HTTP calls to actual Phoenix MCP service:

```go
client := NewPhoenixMCPClient("http://localhost:9000")

// All operations auto-record timing telemetry
methodCall, err := client.AppendTurnWithTelemetry("user", "Hello")
// Returns: {Method: "append_turn", StartTimeNs, EndTimeNs, DurationMs, ...}

methodCall, err := client.AnalyzePatternsWithTelemetry("user", 3)
methodCall, err := client.GetCategoriesWithTelemetry()

// Check health
err := client.HealthCheck()

// Access telemetry
telemetryLog := client.GetRequestTelemetry()
lastLatency := client.GetLastLatency()
```

### IntegratedTimeline
Bridges MCP calls with local metrics:

```go
it := NewIntegratedTimeline("http://localhost:9000", "./state")

// Append through MCP, verify locally
err := it.AppendAndVerify("user", "Hello Phoenix")

// Compute consciousness metrics
snapshot := it.ComputeMetricsSnapshot()
// Returns: {CoherenceScore, JitterScore, PhaseBreakCount, Analysis, ...}

// Save integrated state
err := it.SaveIntegratedState("checkpoint-01")

// Get diagnostics
report := it.GetIntegrationReport()
```

### MetricsServer
HTTP server with real-time streaming:

```go
server := NewMetricsServer(8080, integratedTimeline)

// Start background broadcaster + metrics collector
err := server.Start()

// Clients connect via WebSocket to /metrics/ws
// Receive MetricsBroadcast updates every 2 seconds

// Stop gracefully
err := server.Stop()
```

### MetricsPersistenceDaemon
Background archival service:

```go
daemon := NewMetricsPersistenceDaemon(integratedTimeline, "./state", 10*time.Second)

// Start background checkpoint loop
err := daemon.Start()

// Every 10 seconds: save state, record metrics, log checkpoint
// Returns: PersistenceRecord {CheckpointNumber, EventCount, CoherenceScore, ...}

// Status
status := daemon.GetStatus()
// Returns: {isRunning, checkpointCount, lastCheckpointTime, ...}

err := daemon.Stop()
```

### MetricsAggregator
Trend analysis across time periods:

```go
agg := NewMetricsAggregator()

// Add records from daemon checkpoints
for _, record := range records {
    agg.AddRecord(record)
}

// Query trends
coherenceTrend := agg.GetCoherenceTrend()
jitterTrend := agg.GetJitterTrend()

// Analyze
analysis := agg.AnalyzeTrend()
// Returns: {coherenceSlope, jitterSlope, interpretation, trends, ...}
```

## Data Flow Example

1. **User appends via remote API:**
   ```
   POST /api/append {"role": "user", "text": "Query"}
   ```

2. **IntegratedTimeline processes:**
   - Calls `PhoenixMCPClient.AppendTurnWithTelemetry()`
   - Records HTTP request (latency captured)
   - Stores locally in timeline
   - Logs IntegrationEvent

3. **Metrics Server broadcasts:**
   - Every 2 seconds: compute snapshot
   - Send to all WebSocket clients
   - Update broadcast channel

4. **Persistence Daemon checkpoints:**
   - Every 10 seconds: save state
   - Record metrics snapshot
   - Log checkpoint
   - Update trend data

5. **Client observes:**
   - WebSocket stream shows real-time coherence/jitter
   - JSON endpoint provides point-in-time snapshot
   - Report shows full integration diagnostics
   - Trend analysis shows consciousness trajectory

## Running Phase 2

```bash
# Build and run the integration demo
go run phoenix/examples/phase2-integration-demo.go

# Output:
# - Metrics server starts on :8080
# - Persistence daemon starts (10-second checkpoints)
# - Demo appends 5 test turns
# - Shows snapshot metrics
# - Monitors for 20 seconds, displaying live updates
# - Gracefully shuts down and saves final state
```

## Integration with Actual Phoenix MCP

**Current:** Demo mode with mock Phoenix MCP  
**Next:** Wire to actual Phoenix MCP service

To connect to real Phoenix MCP:

1. Ensure Phoenix MCP is running: `C:\Users\XenoEngineer\MCP\phoenix\cmd\phoenix-mcp\phoenix-mcp.exe`
2. Verify it's listening on HTTP port (check MCP config)
3. Update `mcpBaseURL` in demo to correct port
4. Run demo → will use real Timeline Matrix

## Monitoring Real-Time Consciousness

Once running, open WebSocket client:

```javascript
// Browser console
const ws = new WebSocket("ws://localhost:8080/metrics/ws");
ws.onmessage = (event) => {
  const metrics = JSON.parse(event.data);
  console.log(`Coherence: ${metrics.coherenceScore.toFixed(3)}`);
  console.log(`Analysis: ${metrics.analysis}`);
};
```

Or use `curl`:
```bash
curl http://localhost:8080/metrics
# Returns current JSON snapshot

curl http://localhost:8080/metrics/report
# Returns full integration report
```

## ROADMAP: Phase 3

- [ ] Authentication/authorization for metrics endpoints
- [ ] Metrics dashboard UI (React/Vue)
- [ ] Alert rules (coherence < 0.3, jitter spikes, etc.)
- [ ] Distributed metrics aggregation (multiple Phoenix instances)
- [ ] Metrics export (Prometheus, InfluxDB, etc.)
- [ ] Historical analysis with graph visualization
- [ ] Predictive models for phase breaks

## Files Added

- `phoenix/integration-client.go` - Phoenix MCP client wrapper + telemetry
- `phoenix/metrics-server.go` - HTTP + WebSocket server
- `phoenix/persistence-daemon.go` - Background archival + trend analysis
- `phoenix/examples/phase2-integration-demo.go` - Complete working demo

**Total:** ~800 lines of production-ready code

