# 13 Questions Round 2 - Command Surface Consolidation

Please answer the question below (select one option, or add your own notes). Feel free to add additional context under the question.

## 1. Breaking-Change Policy

Round 1 selected both `sync` compatibility paths for one release in Question 2 and a clean pre-1.0 break with immediate removal in Question 5. Which policy should govern all renamed commands?

- [x] (A) Make a clean break everywhere: remove `workspace ...`, `project add|remove|hydrate|status`, `env pull`, `setup plan|apply`, and the visible `tui install` path when their replacements ship
- [ ] (B) Keep all old paths hidden and deprecated for one minor release, then remove them
- [ ] (C) Keep one-release compatibility only for `workspace ...`; make a clean break for project, env, setup, and TUI commands
- [ ] (D) Keep one-release compatibility only for script-sensitive mutating commands; remove read-only and installer paths immediately
- [ ] (E) Other (describe)

**Current best-practice context:** Cobra supports gradual deprecation, but this project is still pre-1.0 and can intentionally choose a clean contract break. The specification needs one explicit policy so implementation and proof artifacts do not disagree about which old commands must remain executable.

**Recommended answer(s):** [(A)]

**Why these are recommended:**

- `(A)` treats the explicit clean-break selection from Round 1 Question 5 as authoritative and produces the smallest, clearest release surface.
- `(A)` avoids carrying compatibility-only command wiring and tests into a pre-1.0 capstone release.
- `(B)` is safer for existing scripts but contradicts the selected clean-break policy; `(C)` and `(D)` create exceptions that are harder to explain and validate.
