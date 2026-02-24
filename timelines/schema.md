# Timeline Morpheme Schema v0.1

Goal: define a small, flexible event shape for “timeline morphemes” – time-aligned units of work such as keystroke bursts, edit sessions, file hops, or agent exchanges.

## Event envelope

Each morpheme is a JSON object with a standard envelope plus a type-specific payload.

Required fields:

- `event_id` (string, UUID or similar)
- `timestamp_start` (string, ISO 8601)
- `timestamp_end` (string, ISO 8601)
- `channel` (string; e.g. `"keystroke"`, `"edit"`, `"nav"`, `"agent"` )
- `event_type` (string; e.g. `"typing_burst"`, `"file_edit"`, `"window_switch"`, `"agent_turn"`)
- `actor` (string; `"human"` or persona ID like `"plexi-qft"`, `"orchestrator"`)
- `schema_version` (string; e.g. `"tmorph-0.1"`)
- `payload` (object; type-specific details)

## Example: keystroke burst morpheme

```json
{
  "event_id": "tm-2026-02-24T06:05:00Z-001",
  "timestamp_start": "2026-02-24T06:05:00Z",
  "timestamp_end": "2026-02-24T06:05:18Z",
  "channel": "keystroke",
  "event_type": "typing_burst",
  "actor": "human",
  "schema_version": "tmorph-0.1",
  "payload": {
    "window_id": "psiphi-lifeline.md",
    "chars_typed": 240,
    "error_rate": 0.03,
    "mean_iki_ms": 110,
    "burst_label": "high_flow"
  }
}
