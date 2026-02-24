package phoenix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// MetricsServer exposes real-time consciousness metrics via HTTP + WebSocket
type MetricsServer struct {
	port               int
	integratedTimeline *IntegratedTimeline
	server             *http.Server
	upgrader           websocket.Upgrader
	clients            map[*websocket.Conn]bool
	broadcast          chan interface{}
	stop               chan bool
	mu                 sync.RWMutex
	isRunning          bool
}

// MetricsBroadcast is sent to all WebSocket clients
type MetricsBroadcast struct {
	Timestamp         string  `json:"timestamp"`
	CoherenceScore    float64 `json:"coherenceScore"`
	JitterScore       float64 `json:"jitterScore"`
	PhaseBreakCount   int     `json:"phaseBreakCount"`
	LatencySpikeCount int     `json:"latencySpikeCount"`
	MeanLatencyMs     float64 `json:"meanLatencyMs"`
	EventCount        int     `json:"eventCount"`
	MethodCallCount   int     `json:"methodCallCount"`
	LastLatencyMs     float64 `json:"lastLatencyMs"`
	Analysis          string  `json:"analysis"`
}

// NewMetricsServer creates an HTTP server for real-time metrics
func NewMetricsServer(port int, integratedTimeline *IntegratedTimeline) *MetricsServer {
	ms := &MetricsServer{
		port:               port,
		integratedTimeline: integratedTimeline,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins for now (can restrict in production)
				return true
			},
		},
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan interface{}, 10),
		stop:      make(chan bool),
	}

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", ms.healthHandler)
	mux.HandleFunc("/metrics", ms.metricsJSONHandler)
	mux.HandleFunc("/metrics/ws", ms.metricsWebSocketHandler)
	mux.HandleFunc("/metrics/report", ms.metricsReportHandler)
	mux.HandleFunc("/metrics/history", ms.metricsHistoryHandler)
	mux.HandleFunc("/api/append", ms.apiAppendHandler)
	mux.HandleFunc("/api/query", ms.apiQueryHandler)

	ms.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return ms
}

// Start runs the metrics server
func (ms *MetricsServer) Start() error {
	ms.mu.Lock()
	ms.isRunning = true
	ms.mu.Unlock()

	fmt.Printf("[METRICS SERVER] Starting on port %d\n", ms.port)

	// Start broadcaster goroutine
	go ms.broadcastLoop()

	// Start periodic metrics collector
	go ms.metricsCollectorLoop()

	// Start HTTP server
	go func() {
		if err := ms.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("[METRICS SERVER ERROR] %v\n", err)
		}
	}()

	return nil
}

// Stop gracefully shuts down the metrics server
func (ms *MetricsServer) Stop() error {
	ms.mu.Lock()
	ms.isRunning = false
	ms.mu.Unlock()

	close(ms.stop)
	close(ms.broadcast)

	// Close all WebSocket connections
	ms.mu.Lock()
	for client := range ms.clients {
		client.Close()
		delete(ms.clients, client)
	}
	ms.mu.Unlock()

	return ms.server.Close()
}

// broadcastLoop sends metrics to all connected WebSocket clients
func (ms *MetricsServer) broadcastLoop() {
	for {
		select {
		case <-ms.stop:
			return
		case message := <-ms.broadcast:
			ms.mu.RLock()
			for client := range ms.clients {
				ms.mu.RUnlock()
				err := client.WriteJSON(message)
				ms.mu.RLock()
				if err != nil {
					client.Close()
					delete(ms.clients, client)
				}
			}
			ms.mu.RUnlock()
		}
	}
}

// metricsCollectorLoop periodically computes and broadcasts metrics
func (ms *MetricsServer) metricsCollectorLoop() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ms.stop:
			return
		case <-ticker.C:
			snapshot := ms.integratedTimeline.ComputeMetricsSnapshot()
			report := ms.integratedTimeline.GetIntegrationReport()

			broadcast := MetricsBroadcast{
				Timestamp:         time.Now().UTC().Format(time.RFC3339Nano),
				CoherenceScore:    snapshot.CoherenceScore,
				JitterScore:       snapshot.JitterScore,
				PhaseBreakCount:   snapshot.PhaseBreakCount,
				LatencySpikeCount: snapshot.LatencySpikeCount,
				MeanLatencyMs:     snapshot.MeanLatencyMs,
				EventCount:        report["localEventCount"].(int),
				MethodCallCount:   report["methodCallCount"].(int),
				Analysis:          snapshot.Analysis,
			}

			if mcpLatency, ok := report["mcpLatency"].(time.Duration); ok {
				broadcast.LastLatencyMs = float64(mcpLatency.Nanoseconds()) / 1_000_000.0
			}

			select {
			case ms.broadcast <- broadcast:
			case <-ms.stop:
				return
			default:
				// Broadcast channel full, skip this cycle
			}
		}
	}
}

// healthHandler returns server health status
func (ms *MetricsServer) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "healthy",
		"running": ms.isRunning,
		"time":    time.Now().UTC(),
	})
}

// metricsJSONHandler returns current metrics as JSON
func (ms *MetricsServer) metricsJSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	snapshot := ms.integratedTimeline.ComputeMetricsSnapshot()
	json.NewEncoder(w).Encode(snapshot)
}

// metricsReportHandler returns comprehensive integration report
func (ms *MetricsServer) metricsReportHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	report := ms.integratedTimeline.GetIntegrationReport()
	json.NewEncoder(w).Encode(report)
}

// metricsHistoryHandler returns all metrics snapshots
func (ms *MetricsServer) metricsHistoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Note: would need to expose metricsSnapshots in IntegratedTimeline
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "History API coming in v0.2",
	})
}

// metricsWebSocketHandler upgrades connection to WebSocket for streaming metrics
func (ms *MetricsServer) metricsWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := ms.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("[METRICS WS ERROR] Upgrade failed: %v\n", err)
		return
	}

	ms.mu.Lock()
	ms.clients[conn] = true
	ms.mu.Unlock()

	fmt.Printf("[METRICS WS] Client connected (total: %d)\n", len(ms.clients))

	// Send initial snapshot
	snapshot := ms.integratedTimeline.ComputeMetricsSnapshot()
	conn.WriteJSON(MetricsBroadcast{
		Timestamp:         time.Now().UTC().Format(time.RFC3339Nano),
		CoherenceScore:    snapshot.CoherenceScore,
		JitterScore:       snapshot.JitterScore,
		PhaseBreakCount:   snapshot.PhaseBreakCount,
		LatencySpikeCount: snapshot.LatencySpikeCount,
		MeanLatencyMs:     snapshot.MeanLatencyMs,
		Analysis:          "Initial connection",
	})

	// Read from client (keep connection alive)
	for {
		var message map[string]interface{}
		if err := conn.ReadJSON(&message); err != nil {
			ms.mu.Lock()
			delete(ms.clients, conn)
			ms.mu.Unlock()
			fmt.Printf("[METRICS WS] Client disconnected (total: %d)\n", len(ms.clients))
			conn.Close()
			break
		}
	}
}

// apiAppendHandler allows remote clients to append turns
func (ms *MetricsServer) apiAppendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Role string `json:"role"`
		Text string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Delegate to Phoenix MCP via integrated client
	if err := ms.integratedTimeline.AppendAndVerify(req.Role, req.Text); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "appended",
		"role":   req.Role,
		"time":   time.Now().UTC(),
	})
}

// apiQueryHandler allows remote clients to query patterns
func (ms *MetricsServer) apiQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "GET required", http.StatusMethodNotAllowed)
		return
	}

	category := r.URL.Query().Get("category")
	nth := r.URL.Query().Get("nth")

	if category == "" || nth == "" {
		http.Error(w, "category and nth required", http.StatusBadRequest)
		return
	}

	// Query via MCP client
	// (implementation simplified; full parsing omitted)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Query API coming in v0.2",
	})
}
