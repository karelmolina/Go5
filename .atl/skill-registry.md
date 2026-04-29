# Skill Registry

**Delegator use only.** Any agent that launches sub-agents reads this registry to resolve compact rules, then injects them directly into sub-agent prompts. Sub-agents do NOT read this registry or individual SKILL.md files.

See `_shared/skill-resolver.md` for the full resolution protocol.

## User Skills

| Trigger | Skill | Path |
|---------|-------|------|
| When writing Go tests, using teatest, or adding test coverage. | go-testing | /home/karel/.config/opencode/skills/go-testing/SKILL.md |
| When user says "caveman mode", "talk like caveman", "use caveman", "less tokens", "be brief", or invokes /caveman. Also auto-triggers when token efficiency is requested. | caveman | /home/karel/.agents/skills/caveman/SKILL.md |
| When user says "write a commit", "commit message", "generate commit", "/commit", or invokes /caveman-commit. Auto-triggers when staging changes. | caveman-commit | /home/karel/.agents/skills/caveman-commit/SKILL.md |
| When user says "caveman help", "what caveman commands", "how do I use caveman". | caveman-help | /home/karel/.agents/skills/caveman-help/SKILL.md |

## Compact Rules

Pre-digested rules per skill. Delegators copy matching blocks into sub-agent prompts as `## Project Standards (auto-resolved)`.

### go-testing
- Use table-driven tests for multiple cases: `tests := []struct{ name, input, expected string; wantErr bool }`
- Test Bubbletea models directly: `newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})`
- Use teatest for full TUI flows: `tm := teatest.NewTestModel(t, m); tm.Send(...); tm.WaitFinished(t, ...)`
- Golden file testing: compare `m.View()` against `testdata/TestName.golden`
- Mock system info by assigning `m.SystemInfo = &system.SystemInfo{...}`
- Run: `go test ./...`, `go test -cover ./...`, `go test -update ./...`

### caveman
- Default intensity: **full**. Switch with `/caveman lite|full|ultra|wenyan-*`
- Drop articles, filler, pleasantries, hedging. Fragments OK. Short synonyms.
- Pattern: `[thing] [action] [reason]. [next step].`
- Auto-clarity for: security warnings, irreversible actions, multi-step sequences, user asks to clarify
- Code/commits/PRs: write normal. "stop caveman" or "normal mode" to revert.

### caveman-commit
- Conventional Commits format. Subject ≤50 chars, body only when "why" isn't obvious.
- Subject: `<type>(<scope>): <imperative summary>` — types: feat, fix, refactor, perf, docs, test, chore, build, ci, style, revert
- Body: skip when self-explanatory. Add for non-obvious why, breaking changes, migration notes.
- Never include: "This commit does X", "I", "we", "Generated with...", emoji (unless project convention).

### caveman-help
- One-shot display of all caveman modes, skills, and commands.
- Not a persistent mode — just prints the reference card.

## Project Conventions

| File | Path | Notes |
|------|------|-------|
| (none found) | — | No agents.md, CLAUDE.md, .cursorrules, GEMINI.md, or copilot-instructions.md detected |

Read the convention files listed above for project-specific patterns and rules. All referenced paths have been extracted — no need to read index files to discover more.
