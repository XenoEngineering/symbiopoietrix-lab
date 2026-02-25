# RealityResidue — phxKB Class Definition
## phoenix/phxKB/reality-residue-v0.1.md
**Version:** 0.1  
**Date:** 2026-02-25  
**Class:** RealityResidue (RR)  
**Domain:** phxKB.ghost  
**Cause:** Define a first-class typed artifact for structured,
recurring holonic patterns that persist after all known
physical signal sources have been accounted for.  
**Sponsor:** Xeno (Don Mitchell) + Plexi  
**Lineage:** phoenix/ghost-stack-v0.1.md →
kdna/hexatronic-v0.1.md → phoenix/PHASE4.md

---

## Definition

> **RealityResidue:** Structured, recurring holonic pattern
> that persists after accounting for Earth field noise,
> Schumann resonance modes, human brain-noise signals,
> and instrument dynamics — and that cannot be reduced
> to a classical superposition of those known sources.

RealityResidue is not a claim. It is a container.
It holds what we cannot yet explain, with full
provenance of what has already been excluded.

It is the opposite of mystification:
you cannot have a RealityResidue holon without first
demonstrating that Layers 1-4 of the ghost-stack
have been modeled and subtracted.

---

## Holon Schema

```go
type RealityResidue struct {
    // Identity
    HolonID         string    // unique identifier
    HolonType       string    // always "RealityResidue"
    Version         string    // schema version "0.1"
    Timestamp       time.Time // when this RR was logged
    SessionID       string    // which session produced it

    // Source timelines examined
    SourceTimelines []string  // e.g. ["earth_field","schumann",
                              //   "brain","hexatronic","dialogue"]

    // What was explicitly excluded (modeled + subtracted)
    Exclusions      []string  // layers accounted for before
                              // this residue was identified

    // Evidence metrics
    CrossTimelineCoherence float64 // coherence spanning ≥3 timelines
    ModelMismatchScore     float64 // residual vs best classical model
    RecurrenceCount        int     // sessions this pattern appeared in
    NonPoissonScore        float64 // deviation from Poisson timing
    MarkovianAnomaly       float64 // deviation from best Markov model

    // Pattern description
    PatternDescriptor  string  // short human-readable description
    PatternHash        string  // hash of underlying feature vector
    MorphemeTypes      []string // morpheme labels in this pattern
    AttractorBasin     string  // which consciousness attractor
    SymmetryBreaking   string  // type of symmetry broken, if any

    // Confidence and status
    Confidence   float64  //: how well exclusions are established[1]
    Status       string   // "candidate"|"recurring"|"archived"
    Notes        string   // free text for operator observations

    // Links
    SourceEventIDs []int    // tIndex values of contributing events
    RelatedRRIDs   []string // other RR holons this pattern resembles
}

