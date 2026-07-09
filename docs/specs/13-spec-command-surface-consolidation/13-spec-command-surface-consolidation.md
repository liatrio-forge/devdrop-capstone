# 13-spec-command-surface-consolidation.md

## Introduction/Overview

DevSpace has grown from a focused workspace-recovery CLI into a command surface with overlapping namespaces, repeated verbs, and release-grade commands mixed with prototypes. This feature establishes a smaller pre-1.0 CLI contract organized around user goals, makes a clean break from the old command paths, and updates help, packaging, documentation, tests, and demos so the release teaches one consistent workflow.

## Goals

- Reduce the visible root command surface from 16 product commands to no more than 14 clearly grouped commands, excluding Cobra's generated `help` and `completion` commands.
- Give manifest sync, project management, env materialization, and setup execution one canonical vocabulary with no deprecated or hidden compatibility aliases.
- Keep hosted client sync visible while placing the hosted server and FUSE mount prototypes under an explicitly experimental namespace.
- Make `devspace ui` the only visible UI entry point and ship the matching `devspace-tui` companion inside every supported release archive.
- Update help, maintained documentation, automated tests, demos, and release evidence so removed commands are neither advertised nor required.

## User Stories

- **As a new DevSpace user**, I want root help organized around common jobs so that I can identify the capture, restore, maintenance, and diagnostic workflows without reading the full README.
- **As a developer restoring a workspace**, I want manifest synchronization under `devspace sync` so that I do not need to know that the manifest used to be managed through the overloaded `workspace` namespace.
- **As a developer managing projects**, I want `track`, `untrack`, and `update` verbs so that the command names describe metadata and repository intent without implying unsafe file deletion.
- **As a developer managing encrypted env profiles**, I want an explicit `env write` command so that generating a local `.env` file is not confused with pulling remote workspace metadata.
- **As a release consumer**, I want `devspace ui` to work from the downloaded archive without a separate TUI installation command.
- **As a maintainer**, I want command-tree contract tests and synchronized examples so that later features cannot silently reintroduce duplicate or misleading command paths.

## Demoable Units of Work

### Unit 1: Release Command Taxonomy and Help

**Purpose:** Establish the release-facing command tree and make common workflows discoverable from the CLI itself.

**Functional Requirements:**
- The root help shall organize commands into labeled groups for core workflow, workspace management, diagnostics/automation, and experimental functionality.
- The visible release command set shall contain no more than 14 product commands, excluding generated `help` and `completion` commands.
- The canonical visible root commands shall include `init`, `scan`, `status`, `plan`, `apply`, `sync`, `hosted`, `project`, `env`, `setup`, `ui`, `watch`, `doctor`, and `experimental`.
- The root shall stop registering the visible `workspace`, `tui`, `mount`, and `version` commands.
- `devspace --version` shall remain the supported version interface.
- `devspace status` shall retain the workspace health view; `devspace status --verbose` shall include the saved workspace overview previously printed by `devspace workspace`.
- `devspace status <project>` shall replace `devspace project status <project>` and support plain output plus machine-readable output when `--json` is supplied.
- The canonical root and renamed commands shall provide concise `Long` or `Example` content showing the most common invocation and the next safe step.
- Invoking a command group without a required action shall display focused help rather than execute an implicit listing or mutation.
- Removed command paths shall not remain registered as aliases, deprecated commands, hidden commands, or compatibility wrappers.

**Proof Artifacts:**
- CLI capture: `devspace --help` demonstrates the grouped root surface and absence of `workspace`, `tui`, `mount`, and `version` entries.
- CLI capture: `devspace status --help`, `devspace status --verbose`, and `devspace status <project> --json` demonstrate the consolidated status workflow.
- Test output: command-tree contract tests pass and demonstrate that canonical paths resolve while removed paths return command-not-found errors.

### Unit 2: Sync and Project Vocabulary

**Purpose:** Consolidate the two most common resource workflows without changing their existing safety or domain behavior.

**Functional Requirements:**
- The system shall provide `devspace sync push`, `pull`, `diff`, `reconcile`, and `remote` as the only Git-backed manifest sync namespace.
- `devspace sync remote` shall retain the existing `set`, `get`, `create local`, and `create github` capabilities and existing commit identity flags.
- Sync operations shall continue to exchange only validated `manifest.json` metadata and shall preserve existing backup, localization, divergence, hash-guard, reconciliation, JSON, and access-advisory behavior.
- The system shall remove all `devspace workspace ...` paths, including the duplicate `workspace scan` path.
- The system shall provide explicit `devspace project list`, `track`, `untrack`, and `update` commands.
- `project track` shall preserve the current `project add` behavior and shall never modify existing project contents.
- `project untrack` shall preserve the current `project remove` behavior: remove manifest/access/state references, retain local files, and report retained encrypted profiles.
- `project update <project>` and `project update --all` shall be the only project repository update actions; they shall hydrate missing or empty Git projects and fast-forward eligible clean Git checkouts using the existing safety rules.
- The system shall remove `project add`, `remove`, `hydrate`, and `status`, and `devspace project` without a subcommand shall show project help rather than list implicitly.
- Existing human-readable and JSON data contracts shall remain stable where an equivalent canonical command still exposes them.

**Proof Artifacts:**
- CLI workflow: an isolated two-machine run uses `sync remote`, `sync push`, `sync pull`, `plan`, `apply`, and `project update --all` to recreate and hydrate workspace structure.
- Test output: sync regression tests pass under the canonical namespace and confirm that metadata-only, backup, validation, and reconciliation safeguards remain unchanged.
- Test output: project command tests demonstrate list, track, untrack, single-project update, all-project update, dirty-repo skip, and non-destructive file retention behavior.

### Unit 3: Unambiguous Env, Setup, Hosted, Experimental, and UI Paths

**Purpose:** Remove secondary verb collisions and clearly distinguish supported user workflows from prototypes and packaging mechanics.

**Functional Requirements:**
- The system shall provide `devspace env write <project>` as the only command that decrypts a selected env profile and writes the project's local `.env` file.
- `env write` shall preserve profile selection, atomic write behavior, state refresh, and `0600` file permissions from the current env materialization path.
- The system shall remove `devspace env pull`.
- The system shall provide `devspace setup show` for read-only setup review and `devspace setup run <project>` or `devspace setup run --all` for explicit execution.
- `setup run` shall reject simultaneous use of a project argument and `--all`; it shall preserve confirmation, dry-run, unknown-command, and global-install safeguards.
- The system shall remove `devspace setup plan` and `devspace setup apply`.
- The visible `devspace hosted` group shall retain only client operations: `config`, `push`, `pull`, and `reconcile`.
- The visible `devspace experimental` group shall contain `mount` and `hosted serve`, preserving all existing flags, security guards, preview behavior, and prototype labels.
- The system shall remove the visible `devspace tui install` path.
- Every supported release archive shall contain the matching `devspace` and `devspace-tui` executables for its operating system and architecture.
- `devspace ui` shall remain the only visible UI command, prefer the adjacent bundled companion, and retain the built-in legacy dashboard fallback for source builds or incomplete local installations.

**Proof Artifacts:**
- CLI workflow: `env write` materializes a masked test profile into a `0600` `.env` file without exposing the value in captured output.
- CLI workflow: `setup show`, `setup run <project> --dry-run`, and `setup run --all --dry-run` demonstrate the consolidated review and execution model.
- CLI capture: `hosted --help` omits `serve`, while `experimental --help` exposes the prototype paths with explicit labels.
- Release artifact listing: each supported archive contains both platform-matched executables, and a smoke run demonstrates that `devspace ui` locates the adjacent companion.

### Unit 4: Documentation, Demo, and Release Contract

**Purpose:** Make the new command vocabulary the single maintained release contract and prevent historical examples from driving new users toward removed commands.

**Functional Requirements:**
- The README shall lead with task-oriented capture, restore, maintain, and troubleshoot workflows before the complete command reference.
- Maintained architecture, operations, capstone, and demo documentation shall use only the canonical command paths from this specification.
- `scripts/demo-check.sh` and maintained VHS/demo tapes shall exercise the canonical command surface using isolated temporary `DEVSPACE_HOME` and workspace directories.
- Historical completed spec, audit, proof, and validation artifacts shall remain unchanged unless they are executable inputs to a maintained release or demo gate.
- Release notes shall identify the command redesign as an intentional pre-1.0 breaking change and provide a compact old-to-new migration table.
- Command help and maintained docs shall explicitly preserve the product boundary that manifest sync never transfers source code, dependency folders, plaintext `.env` files, or encrypted secret payloads.
- The full Go, TUI, demo, and release-configuration gates shall pass after the command redesign.

**Proof Artifacts:**
- Documentation check: maintained docs and scripts contain the canonical commands and no removed command paths, excluding historical SDD evidence.
- Demo output: `scripts/demo-check.sh` passes with the new capture-to-restore workflow.
- Verification output: `make verify`, `make tui-verify`, and `goreleaser check` pass.
- Release proof: a snapshot or CI release dry-run validates archive construction and companion inclusion for supported targets.

## Non-Goals (Out of Scope)

1. **New restore orchestration:** This feature will not add a `restore` wizard or automatically chain pull, plan, apply, project update, env writing, or setup execution.
2. **Unified Git/hosted backend abstraction:** Hosted client sync remains under `hosted`; `sync --backend` or a shared backend configuration model is not part of this feature.
3. **Domain behavior changes:** Manifest merge rules, Git safety checks, plan/apply hash guards, project update eligibility, env encryption, setup command validation, and access-role enforcement will not be redesigned.
4. **TUI feature redesign:** The OpenTUI and legacy dashboard layouts, key bindings, RPC protocol, and action set will not change except where packaging or help text must reflect the single `ui` entry point.
5. **Productionizing prototypes:** Moving FUSE mount and hosted serve under `experimental` does not make either capability production-supported.
6. **Permanent compatibility aliases:** Removed commands will not remain functional after this pre-1.0 breaking change.
7. **Historical artifact rewriting:** Completed SDD specs and proof records remain historical evidence even when they mention the old command surface.
8. **Broad command-layer refactoring:** Splitting `commands.go` may occur only when necessary to implement and test the new command tree; unrelated domain or package restructuring is excluded.

## Design Considerations

- Root help must be scannable without reading every command. Use Cobra command groups and short descriptions that begin with the user's goal rather than an implementation detail.
- Keep the common command depth at three levels or fewer (`devspace resource action`). Deeper nesting is acceptable only for infrequent configuration or experimental paths such as `sync remote create github` and `experimental hosted serve`.
- Help examples must show safe sequencing. In particular, `sync pull` should point to `plan`, `apply`, and then `project update --all` without implying that source repositories or secrets were synchronized.
- Use consistent verbs: `pull` for remote manifest retrieval, `write` for local env materialization, `show` for read-only setup review, `run` for setup execution, `track`/`untrack` for manifest membership, and `update` for Git hydration/fast-forward behavior.
- Human output may retain the established Fang/Lipgloss presentation. JSON output must remain free of ANSI styling and stable for automation.
- The `experimental` help text must state that its commands are prototypes and are not part of the supported recovery workflow.

## Repository Standards

- Follow `AGENTS.md`: Go code lives under `internal/devspace/`, uses standard formatting, returns explicit user-facing errors, and remains local-first unless a command explicitly opts into Git or hosted sync.
- Keep Cobra command construction separate from domain behavior. Canonical and experimental commands must reuse the existing domain functions, output helpers, and `withAppLock` boundaries rather than duplicate business logic.
- Use Go's built-in `testing` package with command tests beside the implementation. Test with isolated temporary application homes and workspaces; never use the developer's real `~/.devspace` state.
- Preserve JSON field names and clean-output guarantees unless this specification explicitly defines a new report shape.
- Run `gofmt` on changed Go files and use Conventional Commit-style subjects.
- Required verification is `make verify` for Go, `make tui-verify` for the companion, the maintained demo check, and release configuration validation.
- CLI output changes require captured terminal or text proof; release packaging changes require archive or workflow evidence.

## Technical Considerations

- The implementation shall use the existing Cobra v1.10.2 command tree and Fang v2 execution layer; no new CLI framework or configuration dependency is needed.
- Current Cobra guidance supports command groups, aliases, hidden commands, deprecation, and examples. This feature intentionally deviates from gradual deprecation guidance because the user selected a clean pre-1.0 break; old paths must be removed rather than retained through Cobra aliases or `Deprecated` fields. See [Working with Commands](https://cobra.dev/docs/how-to-guides/working-with-commands/) and [Working with Flags](https://cobra.dev/docs/how-to-guides/working-with-flags/).
- Prefer small command constructors that share existing `RunE` helpers or domain calls. Do not create parallel implementations for old and new paths.
- Root and group commands should use Cobra `AddGroup`/`GroupID`, `Example`, argument validators, and mutually exclusive flag validation where appropriate.
- `status <project> --json` requires one documented project-status JSON shape. Prefer reusing the existing project list row representation rather than introducing a second overlapping project report.
- Release packaging must associate each Go archive target with the corresponding prebuilt `devspace-tui` artifact instead of publishing the companion only as a separate release attachment.
- Source and development workflows may continue to use `make tui-install-local`; it is a maintainer command and not part of the release CLI contract.
- Maintained documentation includes the root README, architecture and operations docs, current capstone/demo guidance, executable scripts, and active release evidence. Completed historical SDD artifacts are excluded from command-reference migration scans.

## Security Considerations

- Renaming or moving commands must not weaken path validation, atomic writes, application locking, Git remote validation, manifest conflict protection, or access-role advisory behavior.
- `experimental hosted serve` must retain the existing loopback default, public-cleartext bind refusal, trusted-proxy validation, token handling, and TLS-termination warnings.
- `env write` must never print decrypted values and must preserve `0600` output permissions and encrypted-at-rest storage.
- Sync help, docs, and output must continue to state that only manifest metadata is transferred; source code, dependency directories, `.env` files, identities, and secret payloads are excluded.
- Proof artifacts must use placeholders or temporary credentials and must not commit real hosted tokens, age identities, `.env` values, developer workspace state, or authenticated endpoint details.
- Release-archive validation must verify artifact identity and checksums without embedding credentials in logs or proof files.

## Success Metrics

1. **Root discoverability:** Root help exposes no more than 14 product commands, groups them by user goal, and contains none of the removed root commands.
2. **Single vocabulary:** All maintained docs, scripts, and help use `sync`, `project list|track|untrack|update`, `env write`, `setup show|run`, `ui`, and `experimental` consistently; migration scans find zero removed paths outside historical evidence.
3. **Clean break:** Automated contract tests prove every canonical path resolves and every removed path fails as an unknown command.
4. **Behavior preservation:** Existing safety, JSON, sync, project, env, setup, hosted, and TUI regression tests pass under the new command wiring.
5. **Release completeness:** Each supported Linux/macOS amd64/arm64 archive contains both matching executables and passes the release smoke check.
6. **Quality gates:** `make verify`, `make tui-verify`, the maintained demo check, and release configuration validation complete successfully.

## Open Questions

No open questions at this time.
