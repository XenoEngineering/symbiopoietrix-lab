# LiminalConfig v0.1
ID: liminalConfig-v0.1

Purpose
- Define named liminality modes reused by Inferon and Phoenix.
- Control how wild a run is allowed to be across layers and timelines.

## Structure

LiminalConfig {
    ID:              string        // e.g. "VISION_QUEST", "ANNEALING_STRICT"
    LayerFlex:       LayerFlex     // how much to loosen each layer
    ExplorationBand: Band          // narrative band: utility ↔ poiesis
    TimelineFreedom: TimeMode      // how far we can roam across branches
    Sponsor:         HolonRef      // who authorizes this liminality
}

LayerFlex {
    Artifact: FlexLevel   // NONE | LOW | MEDIUM | HIGH
    Syntax:   FlexLevel
    Semantics:FlexLevel
    Intent:   FlexLevel
}

Band = ENUM {
    STRICT_UTILITY,   // minimal play, focus on precise outcomes
    REFINEMENT,       // clarify, compress, formalize
    PLAY,             // light exploration, keep coherence
    VISION_QUEST      // maximal poiesis, tolerate wild forms
}

TimeMode = ENUM {
    LOCAL_ONLY,           // stay in current timeline window
    NEAR_BRANCHES,        // siblings/parents/children only
    CROSS_BRANCH_PARALLEL // free to mix parallel ranges
}
