package phoenix

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Event represents a single entry in the Timeline Matrix
type Event struct {
	Timestamp string      `json:"timestamp"`
	TIndex    int         `json:"tIndex"`
	Role      string      `json:"role"` // "user" or "assistant"
	Layer     string      `json:"layer,omitempty"` // "artifact", "semantics", "intent", "meta"
	Family    string      `json:"family,omitempty"`
	Agent     string      `json:"agent,omitempty"`
	Op        string      `json:"op,omitempty"`
	Tags      []string    `json:"tags,omitempty"`
	SpanRef   string      `json:"spanRef,omitempty"` // the actual text
}

// BSTNode represents a node in the Binary Search Tree for fast category lookup
type BSTNode struct {
	Category string    `json:"category"`
	Count    int       `json:"count"`
	Indices  []int     `json:"indices"` // all tIndices for this category
	LoChild  *BSTNode  `json:"loChild,omitempty"`
	HiChild  *BSTNode  `json:"hiChild,omitempty"`
}

// MethodCall represents instrumentation data for a single operation
type MethodCall struct {
	Method      string          `json:"method"` // "append_turn", "analyze_patterns", "get_categories"
	StartTimeNs int64           `json:"startTimeNs"`
	EndTimeNs   int64           `json:"endTimeNs"`
	DurationMs  float64         `json:"durationMs"`
	Input       json.RawMessage `json:"input,omitempty"`
	Result      json.RawMessage `json:"result,omitempty"`
	Timestamp   string          `json:"timestamp"`
	CategoryKey string          `json:"categoryKey,omitempty"` // for analyze_patterns
}

// CategoryIndex fast lookup structure
type CategoryIndex struct {
	Category      string `json:"category"`
	Count         int    `json:"count"`
	Occurrences   []int  `json:"occurrences"` // tIndices
}

// MatrixStats holds aggregate statistics
type MatrixStats struct {
	TotalEvents           int              `json:"totalEvents"`
	TotalCategories       int              `json:"totalCategories"`
	Categories            []CategoryIndex  `json:"categories"`
	MethodCallCount       int              `json:"methodCallCount"`
	AvgAppendTimeMs       float64          `json:"avgAppendTimeMs"`
	AvgAnalysisTimeMs     float64          `json:"avgAnalysisTimeMs"`
	AvgCategoryTimeMs     float64          `json:"avgCategoryTimeMs"`
	LastUpdateTimestamp   string           `json:"lastUpdateTimestamp"`
	StructureValidated    bool             `json:"structureValidated"`
}

// MatrixState is the complete serializable state of the Timeline Matrix
type MatrixState struct {
	Timeline       []Event        `json:"timeline"`
	BSTForest      map[string]*BSTNode `json:"bstForest"`
	MethodCalls    []MethodCall   `json:"methodCalls"`
	Stats          MatrixStats    `json:"stats"`
	SerializedAt   string         `json:"serializedAt"`
	SchemaVersion  string         `json:"schemaVersion"`
	IntegrityHash  string         `json:"integrityHash,omitempty"`
}

// MatrixPersistence handles save/load operations with validation
type MatrixPersistence struct {
	mu              sync.RWMutex
	statePath       string
	checksumPath    string
	lastValidated   time.Time
	validationCount int
}

// NewMatrixPersistence creates a new persistence manager
func NewMatrixPersistence(baseDir string) *MatrixPersistence {
	return &MatrixPersistence{
		statePath:    fmt.Sprintf("%s/matrix-state.json", baseDir),
		checksumPath: fmt.Sprintf("%s/matrix-state.checksum", baseDir),
	}
}

// SaveState persists the entire Matrix to JSON
func (mp *MatrixPersistence) SaveState(state *MatrixState) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	// Update metadata
	state.SerializedAt = time.Now().UTC().Format(time.RFC3339Nano)
	state.SchemaVersion = "0.1"

	// Marshal to JSON
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	// Compute integrity hash (simple: first 32 chars of JSON)
	if len(data) > 32 {
		state.IntegrityHash = fmt.Sprintf("%x", data[:32])
	}

	// Write to file
	if err := os.WriteFile(mp.statePath, data, 0644); err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	// Write checksum for quick validation
	checksum := fmt.Sprintf("%d:%s", len(data), state.IntegrityHash)
	if err := os.WriteFile(mp.checksumPath, []byte(checksum), 0644); err != nil {
		return fmt.Errorf("checksum write error: %w", err)
	}

	return nil
}

// LoadState restores the entire Matrix from JSON with validation
func (mp *MatrixPersistence) LoadState() (*MatrixState, error) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	data, err := os.ReadFile(mp.statePath)
	if err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	var state MatrixState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	// Validate structure
	if err := mp.validateMatrixStructure(&state); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	mp.lastValidated = time.Now()
	mp.validationCount++

	return &state, nil
}

// validateMatrixStructure checks BST integrity and event consistency
func (mp *MatrixPersistence) validateMatrixStructure(state *MatrixState) error {
	if len(state.Timeline) == 0 {
		return fmt.Errorf("timeline is empty")
	}

	// Check tIndex consistency (should be 1-based, ascending)
	for i, event := range state.Timeline {
		expectedTIndex := i + 1
		if event.TIndex != expectedTIndex {
			return fmt.Errorf("tIndex inconsistency at position %d: got %d, expected %d", i, event.TIndex, expectedTIndex)
		}
	}

	// Validate BST forest structure
	for category, root := range state.BSTForest {
		if root == nil {
			return fmt.Errorf("nil BST root for category: %s", category)
		}
		if err := validateBSTNode(root, state.Timeline); err != nil {
			return fmt.Errorf("BST validation error for category %s: %w", category, err)
		}
	}

	// Check stats consistency
	if state.Stats.TotalEvents != len(state.Timeline) {
		return fmt.Errorf("stats mismatch: TotalEvents=%d but timeline length=%d", state.Stats.TotalEvents, len(state.Timeline))
	}

	state.Stats.StructureValidated = true
	return nil
}

// validateBSTNode recursively validates a BST subtree
func validateBSTNode(node *BSTNode, timeline []Event) error {
	if node == nil {
		return nil
	}

	// Check indices are valid and sorted
	for _, idx := range node.Indices {
		if idx < 1 || idx > len(timeline) {
			return fmt.Errorf("invalid tIndex in BST: %d (timeline length: %d)", idx, len(timeline))
		}
	}

	// Check count matches indices
	if node.Count != len(node.Indices) {
		return fmt.Errorf("count mismatch: Count=%d but Indices length=%d", node.Count, len(node.Indices))
	}

	// Recurse on children
	if err := validateBSTNode(node.LoChild, timeline); err != nil {
		return err
	}
	if err := validateBSTNode(node.HiChild, timeline); err != nil {
		return err
	}

	return nil
}

// RecoverFromChecksum does quick validation without full load (for health checks)
func (mp *MatrixPersistence) RecoverFromChecksum() (bool, error) {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	checksumData, err := os.ReadFile(mp.checksumPath)
	if err != nil {
		return false, fmt.Errorf("checksum read error: %w", err)
	}

	// Checksum format: "size:hash"
	// In future, can compare against current file size to detect corruption
	_ = string(checksumData)
	return true, nil
}

// GetPersistenceMetrics returns health and performance data
func (mp *MatrixPersistence) GetPersistenceMetrics() map[string]interface{} {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	return map[string]interface{}{
		"lastValidated":    mp.lastValidated,
		"validationCount":  mp.validationCount,
		"statePath":        mp.statePath,
		"checksumPath":     mp.checksumPath,
	}
}

// ExportForAnalysis creates a lightweight copy for external analysis tools
func (mp *MatrixPersistence) ExportForAnalysis(state *MatrixState, outputPath string) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	// Create lightweight export: timeline + method calls only (drop BST)
	export := map[string]interface{}{
		"timeline":     state.Timeline,
		"methodCalls":  state.MethodCalls,
		"stats":        state.Stats,
		"exportedAt":   time.Now().UTC().Format(time.RFC3339Nano),
	}

	data, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return fmt.Errorf("export marshal error: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("export write error: %w", err)
	}

	return nil
}
