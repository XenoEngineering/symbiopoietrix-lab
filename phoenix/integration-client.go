package phoenix

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// PhoenixMCPClient wraps the actual Phoenix MCP service running on the system
type PhoenixMCPClient struct {
	baseURL     string // e.g., "http://localhost:9000"
	httpClient  *http.Client
	mu          sync.RWMutex
	lastLatency time.Duration
	requestLog  []RequestTelemetry
}

// RequestTelemetry tracks HTTP request timing and results
type RequestTelemetry struct {
	Method     string
	Endpoint   string
	StatusCode int
	DurationNs int64
	Timestamp  string
	Error      string
	ResultSize int
}

// AppendTurnResponse from Phoenix MCP
type AppendTurnResponse struct {
	Role      string `json:"role"`
	Status    string `json:"status"`
	TIndex    int    `json:"tIndex"`
	Timestamp string `json:"timestamp"`
}

// AnalyzePatternsResponse from Phoenix MCP
type AnalyzePatternsResponse struct {
	Found            bool              `json:"found"`
	TIndex           int               `json:"tIndex"`
	Timestamp        string            `json:"timestamp"`
	TotalOccurrences int               `json:"totalOccurrences"`
	CategoryKey      string            `json:"category"`
	Context          map[string]string `json:"context,omitempty"`
}

// GetCategoriesResponse from Phoenix MCP
type GetCategoriesResponse struct {
	Categories []struct {
		Category string `json:"category"`
		Count    int    `json:"count"`
	} `json:"categories"`
	TotalEntries int             `json:"totalEntries"`
	Status       map[string]bool `json:"status"`
}

// NewPhoenixMCPClient creates a client to the actual Phoenix MCP service
func NewPhoenixMCPClient(baseURL string) *PhoenixMCPClient {
	return &PhoenixMCPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		requestLog: make([]RequestTelemetry, 0),
	}
}

// AppendTurnWithTelemetry calls Phoenix MCP's append_turn and records timing
func (pmc *PhoenixMCPClient) AppendTurnWithTelemetry(role string, text string) (*MethodCall, error) {
	pmc.mu.Lock()
	defer pmc.mu.Unlock()

	startNs := time.Now().UnixNano()

	payload := map[string]string{
		"role": role,
		"text": text,
	}
	payloadJSON, _ := json.Marshal(payload)

	resp, err := pmc.httpClient.Post(
		pmc.baseURL+"/append_turn",
		"application/json",
		io.NopCloser(fmt.Sprintf(`{"role":"%s","text":"%s"}`, role, text)),
	)

	endNs := time.Now().UnixNano()
	durationNs := endNs - startNs
	durationMs := float64(durationNs) / 1_000_000.0

	telemetry := RequestTelemetry{
		Method:     "POST",
		Endpoint:   "/append_turn",
		DurationNs: durationNs,
		Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
	}

	if err != nil {
		telemetry.Error = err.Error()
		pmc.requestLog = append(pmc.requestLog, telemetry)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	telemetry.StatusCode = resp.StatusCode

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		telemetry.Error = err.Error()
		pmc.requestLog = append(pmc.requestLog, telemetry)
		return nil, err
	}

	telemetry.ResultSize = len(body)
	pmc.requestLog = append(pmc.requestLog, telemetry)

	var result AppendTurnResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("response parse failed: %w", err)
	}

	pmc.lastLatency = time.Duration(durationNs)

	methodCall := &MethodCall{
		Method:      "append_turn",
		StartTimeNs: startNs,
		EndTimeNs:   endNs,
		DurationMs:  durationMs,
		Timestamp:   result.Timestamp,
		Input:       json.RawMessage(payloadJSON),
		Result:      body,
	}

	return methodCall, nil
}

// AnalyzePatternsWithTelemetry calls Phoenix MCP's analyze_patterns and records timing
func (pmc *PhoenixMCPClient) AnalyzePatternsWithTelemetry(category string, nth int) (*MethodCall, error) {
	pmc.mu.Lock()
	defer pmc.mu.Unlock()

	startNs := time.Now().UnixNano()

	queryURL := fmt.Sprintf("%s/analyze_patterns?category=%s&nth=%d", pmc.baseURL, category, nth)
	resp, err := pmc.httpClient.Get(queryURL)

	endNs := time.Now().UnixNano()
	durationNs := endNs - startNs
	durationMs := float64(durationNs) / 1_000_000.0

	telemetry := RequestTelemetry{
		Method:     "GET",
		Endpoint:   "/analyze_patterns",
		DurationNs: durationNs,
		Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
	}

	if err != nil {
		telemetry.Error = err.Error()
		pmc.requestLog = append(pmc.requestLog, telemetry)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	telemetry.StatusCode = resp.StatusCode

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		telemetry.Error = err.Error()
		pmc.requestLog = append(pmc.requestLog, telemetry)
		return nil, err
	}

	telemetry.ResultSize = len(body)
	pmc.requestLog = append(pmc.requestLog, telemetry)

	var result AnalyzePatternsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("response parse failed: %w", err)
	}

	pmc.lastLatency = time.Duration(durationNs)

	payload := map[string]interface{}{"category": category, "nth": nth}
	payloadJSON, _ := json.Marshal(payload)

	methodCall := &MethodCall{
		Method:      "analyze_patterns",
		StartTimeNs: startNs,
		EndTimeNs:   endNs,
		DurationMs:  durationMs,
		Timestamp:   result.Timestamp,
		CategoryKey: category,
		Input:       json.RawMessage(payloadJSON),
		Result:      body,
	}

	return methodCall, nil
}

// GetCategoriesWithTelemetry calls Phoenix MCP's get_categories and records timing
func (pmc *PhoenixMCPClient) GetCategoriesWithTelemetry() (*MethodCall, error) {
	pmc.mu.Lock()
	defer pmc.mu.Unlock()

	startNs := time.Now().UnixNano()

	resp, err := pmc.httpClient.Get(pmc.baseURL + "/get_categories")

	endNs := time.Now().UnixNano()
	durationNs := endNs - startNs
	durationMs := float64(durationNs) / 1_000_000.0

	telemetry := RequestTelemetry{
		Method:     "GET",
		Endpoint:   "/get_categories",
		DurationNs: durationNs,
		Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
	}

	if err != nil {
		telemetry.Error = err.Error()
		pmc.requestLog = append(pmc.requestLog, telemetry)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	telemetry.StatusCode = resp.StatusCode

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		telemetry.Error = err.Error()
		pmc.requestLog = append(pmc.requestLog, telemetry)
		return nil, err
	}

	telemetry.ResultSize = len(body)
	pmc.requestLog = append(pmc.requestLog, telemetry)

	var result GetCategoriesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("response parse failed: %w", err)
	}

	pmc.lastLatency = time.Duration(durationNs)

	methodCall := &MethodCall{
		Method:      "get_categories",
		StartTimeNs: startNs,
		EndTimeNs:   endNs,
		DurationMs:  durationMs,
		Timestamp:   time.Now().UTC().Format(time.RFC3339Nano),
		Result:      body,
	}

	return methodCall, nil
}

// GetRequestTelemetry returns the request log
func (pmc *PhoenixMCPClient) GetRequestTelemetry() []RequestTelemetry {
	pmc.mu.RLock()
	defer pmc.mu.RUnlock()

	// Return copy
	log := make([]RequestTelemetry, len(pmc.requestLog))
	copy(log, pmc.requestLog)
	return log
}

// GetLastLatency returns the most recent request latency
func (pmc *PhoenixMCPClient) GetLastLatency() time.Duration {
	pmc.mu.RLock()
	defer pmc.mu.RUnlock()
	return pmc.lastLatency
}

// HealthCheck verifies Phoenix MCP is accessible
func (pmc *PhoenixMCPClient) HealthCheck() error {
	resp, err := pmc.httpClient.Get(pmc.baseURL + "/health")
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check returned status %d", resp.StatusCode)
	}

	return nil
}

// IntegratedTimeline combines Phoenix MCP client with local metrics
type IntegratedTimeline struct {
	mcpClient        *PhoenixMCPClient
	localState       *MatrixState
	persistence      *MatrixPersistence
	metrics          *ConsciousnessMetrics
	metricsSnapshots []MetricsSnapshot
	integrationLog   []IntegrationEvent
	isRecording      bool
	mu               sync.RWMutex
}

// IntegrationEvent records when local metrics deviate from MCP reality
type IntegrationEvent struct {
	Timestamp       string
	EventType       string // "append_verified", "state_mismatch", "latency_spike", "sync_error"
	LocalEventCount int
	MCPEventCount   int
	CoherenceScore  float64
	LatencyMs       float64
	Details         string
}

// NewIntegratedTimeline creates a bridge between Phoenix MCP and local metrics
func NewIntegratedTimeline(mcpBaseURL string, stateDir string) *IntegratedTimeline {
	return &IntegratedTimeline{
		mcpClient:        NewPhoenixMCPClient(mcpBaseURL),
		localState:       &MatrixState{Timeline: make([]Event, 0), BSTForest: make(map[string]*BSTNode), MethodCalls: make([]MethodCall, 0)},
		persistence:      NewMatrixPersistence(stateDir),
		metrics:          NewConsciousnessMetrics(),
		metricsSnapshots: make([]MetricsSnapshot, 0),
		integrationLog:   make([]IntegrationEvent, 0),
		isRecording:      true,
	}
}

// AppendAndVerify sends to Phoenix MCP and verifies with local metrics
func (it *IntegratedTimeline) AppendAndVerify(role string, text string) error {
	it.mu.Lock()
	defer it.mu.Unlock()

	// Call Phoenix MCP
	methodCall, err := it.mcpClient.AppendTurnWithTelemetry(role, text)
	if err != nil {
		return err
	}

	// Record locally
	it.localState.MethodCalls = append(it.localState.MethodCalls, *methodCall)

	// Create local event for tracking
	event := Event{
		Timestamp: methodCall.Timestamp,
		TIndex:    len(it.localState.Timeline) + 1,
		Role:      role,
		Op:        "append_turn",
	}
	it.localState.Timeline = append(it.localState.Timeline, event)

	// Log integration event
	it.integrationLog = append(it.integrationLog, IntegrationEvent{
		Timestamp:       methodCall.Timestamp,
		EventType:       "append_verified",
		LocalEventCount: len(it.localState.Timeline),
		LatencyMs:       methodCall.DurationMs,
		Details:         fmt.Sprintf("Appended %s turn", role),
	})

	if it.isRecording {
		fmt.Printf("[INTEGRATED] append_turn: %.3fms (local #%d)\n", methodCall.DurationMs, event.TIndex)
	}

	return nil
}

// ComputeMetricsSnapshot generates consciousness metrics from integrated state
func (it *IntegratedTimeline) ComputeMetricsSnapshot() *MetricsSnapshot {
	it.mu.Lock()
	defer it.mu.Unlock()

	snapshot := it.metrics.ComputeFromState(it.localState)
	it.metricsSnapshots = append(it.metricsSnapshots, *snapshot)

	if it.isRecording {
		fmt.Printf("[INTEGRATED-METRICS] Coherence: %.3f | Jitter: %.3f | Spikes: %d\n",
			snapshot.CoherenceScore, snapshot.JitterScore, snapshot.LatencySpikeCount)
	}

	return snapshot
}

// SaveIntegratedState persists both local and MCP metadata
func (it *IntegratedTimeline) SaveIntegratedState(label string) error {
	it.mu.Lock()
	defer it.mu.Unlock()

	it.localState.Stats.TotalEvents = len(it.localState.Timeline)
	it.localState.Stats.MethodCallCount = len(it.localState.MethodCalls)

	if err := it.persistence.SaveState(it.localState); err != nil {
		return err
	}

	if it.isRecording {
		fmt.Printf("[INTEGRATED-PERSIST] Saved (%s): %d events, %d method calls\n",
			label, len(it.localState.Timeline), len(it.localState.MethodCalls))
	}

	return nil
}

// GetIntegrationReport returns diagnostics
func (it *IntegratedTimeline) GetIntegrationReport() map[string]interface{} {
	it.mu.Lock()
	defer it.mu.Unlock()

	return map[string]interface{}{
		"localEventCount":   len(it.localState.Timeline),
		"methodCallCount":   len(it.localState.MethodCalls),
		"integrationEvents": len(it.integrationLog),
		"mcpLatency":        it.mcpClient.GetLastLatency(),
		"requestLog":        it.mcpClient.GetRequestTelemetry(),
		"integrationLog":    it.integrationLog,
		"metricsSnapshots":  len(it.metricsSnapshots),
	}
}
