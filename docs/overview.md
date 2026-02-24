# Symbiopoietrix Lab — Overview

Symbiopoietrix Lab is an experiment in working with AI as a **research partner** instead of a black box. It gives one human and a small family of AI personas a shared workspace, shared language, and simple rules for how to collaborate.

At a high level:

- The **human** decides what matters, chooses the direction, and names the important objects.
- **AI personas** (like Plexi-QFT and Plexi-Implementor) each have a small, written “job description”.
- Short text snippets called **K-DNA blocks** capture important domain ideas (like what “psiPhi” or “UPT” means).
- Simple **timeline logs** record how the work moves over time (typing bursts, file switches, AI replies).
- An **Orchestrator** persona helps route tasks to the right specialist persona, instead of one agent trying to do everything.

The goal is not to automate the human away, but to:
- Make AI behavior more predictable and inspectable.
- Turn long, messy projects into small, traceable steps.
- Keep the structure of the work visible over time.

## How to use this lab as a new human

1. Read `docs/README.md` and this `overview.md` to get the basic idea.
2. Skim `docs/lifeline-psiphi.md` to see what psiPhi connects (physics, timelines, tools).
3. Look at `personas/` to understand what each AI persona is “allowed” to do.
4. Glance at `kdna/` to see the small K-DNA blocks that define key concepts.
5. When you start a session with an AI model, pick:
   - A persona spec (from `personas/`),
   - One or more K-DNA blocks (from `kdna/`),
   - A small, concrete task.
6. As you work, you may log timeline events under `timelines/` or extend the K-DNA and personas, but you never have to touch everything at once.

If you only remember one sentence: this lab is a way to keep human intent, AI personas, and the history of work aligned in one place.
