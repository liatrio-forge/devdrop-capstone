# 13 Questions Round 1 - Command Surface Consolidation

Please answer each question below (select one or more options, or add your own notes). Feel free to add additional context under any question.

## 1. Supported Release Surface

How should the release CLI expose the hosted-sync server and FUSE mount prototypes?

- [ ] (A) Move them under a hidden `devspace experimental` namespace; keep them executable for capstone proofs but omit them from normal root help
- [ ] (B) Remove them from the release binary and keep them available only through development builds or separate binaries
- [ ] (C) Keep `devspace hosted` and `devspace mount` visible at the root and improve their prototype labels
- [x] (D) Move only `hosted serve` and `mount` under `experimental`, while keeping hosted client sync visible
- [ ] (E) Other (describe)

**Current best-practice context:** Cobra supports hidden commands for administrative or transitional functionality, and its current guidance recommends grouping commands by user-facing domain rather than implementation detail.

**Recommended answer(s):** [(D)]

**Why these are recommended:**

- `(D)` preserves hosted client testing and the existing hosted roadmap without presenting a prototype server or FUSE mount as release-grade core behavior.
- `(D)` keeps capstone proof paths executable while removing two advanced implementation surfaces from the beginner workflow.
- `(A)` is simpler but would hide all hosted client sync even though it mirrors a real user workflow; `(B)` introduces packaging work; `(C)` preserves the reported command overload.

## 2. Canonical Manifest Sync Namespace

Which command family should be canonical for Git-backed manifest synchronization?

- [x] (A) Introduce `devspace sync push|pull|diff|reconcile|remote`; retain `workspace ...` as hidden deprecated compatibility paths for one release
- [ ] (B) Keep `devspace workspace ...` canonical and improve help grouping only
- [ ] (C) Introduce one generic `devspace sync --backend git|hosted` surface for both Git and hosted synchronization
- [ ] (D) Keep Git under `workspace` and hosted under `hosted`, but add a new task-oriented `restore` command
- [ ] (E) Other (describe)

**Current best-practice context:** Current Cobra guidance recommends shallow command trees organized by user-facing domains. It also recommends avoiding excessive aliases and using compatibility paths only during migration.

**Recommended answer(s):** [(A)]

**Why these are recommended:**

- `(A)` names the user goal directly and removes the ambiguity between the executable `workspace` overview and its sync subcommands.
- `(A)` avoids the new backend abstraction and repeated flags required by `(C)` while keeping hosted sync independently experimental.
- `(B)` does not address the reported discoverability problem, and `(D)` adds another command without removing the existing complexity.

## 3. Project Command Vocabulary

How aggressively should project commands adopt intent-revealing names?

- [x] (A) Make `project list|track|untrack|update` canonical; make `hydrate`, `add`, `remove`, and `project status` hidden deprecated compatibility paths
- [ ] (B) Keep `add|remove|hydrate|update`, but improve descriptions and examples
- [ ] (C) Rename only `remove` to `untrack`; keep the other project verbs unchanged
- [ ] (D) Keep all existing commands and add aliases for `track` and `untrack`
- [ ] (E) Other (describe)

**Recommended answer(s):** [(A)]

**Why these are recommended:**

- `(A)` makes filesystem safety explicit: tracking and untracking describe manifest operations without implying file creation or deletion.
- The pending `project update` already hydrates missing projects, so keeping both as equally visible canonical actions would preserve overlapping behavior.
- `(D)` increases the command vocabulary, while `(B)` and `(C)` leave most of the reported ambiguity intact.

## 4. Secondary Verb Collisions

Should the release rename the env and setup commands whose verbs conflict with workspace operations?

- [x] (A) Rename `env pull` to `env write`, `setup plan` to `setup show`, and `setup apply` to `setup run --all`; retain hidden deprecated compatibility paths for one release
- [ ] (B) Rename only `env pull` because `pull` currently has three different meanings
- [ ] (C) Keep all existing names and rely on command-group help text
- [ ] (D) Replace env and setup commands with interactive prompts in `devspace ui`
- [ ] (E) Other (describe)

**Recommended answer(s):** [(A)]

**Why these are recommended:**

- `(A)` gives each common verb one primary meaning and removes the second plan/apply lifecycle from the same CLI.
- `(A)` preserves automation through explicit compatibility paths and does not require new interactive behavior.
- `(C)` improves documentation but not predictability; `(D)` expands this feature into a larger UI redesign.

## 5. Compatibility and Companion TUI Policy

What migration policy should apply to renamed commands and the `devspace-tui` installer?

- [ ] (A) Keep old commands functional but hidden and deprecated for one minor release; ship the companion in release archives so `devspace ui` is the only visible UI command
- [x] (B) Make a clean breaking change before 1.0 and remove all old command paths immediately
- [ ] (C) Preserve old aliases indefinitely and keep `devspace tui install` visible
- [ ] (D) Keep command aliases for one release, but continue distributing the companion through the visible installer
- [ ] (E) Other (describe)

**Current best-practice context:** Cobra's current flag guidance recommends keeping deprecated interfaces active for at least one minor release with clear replacement messages. Its command guidance recommends limiting aliases to one or two obvious compatibility names so help and completion remain predictable.

**Recommended answer(s):** [(A)]

**Why these are recommended:**

- `(A)` protects existing scripts and checked-in demos while presenting one stable UI entry point to new users.
- Bundling the companion removes installation mechanics from the product vocabulary and aligns release behavior with `devspace ui`.
- `(B)` is cheaper but unnecessarily breaks existing workflows; `(C)` preserves permanent clutter; `(D)` leaves the `ui` versus `tui` distinction visible.
