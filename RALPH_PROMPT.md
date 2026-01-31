Execute ApertureStack Consolidation Phase 5: Composition Layer (PRD-150 + PRD-151)

## Context
- Working directory: /Users/jraymond/Documents/Projects/ApertureStack/toolcompose
- Source repos: ../toolset, ../toolskill
- Target: toolcompose with set/ and skill/ subpackages
- Dependencies already consolidated:
  - toolfoundation/model (from toolmodel)
  - tooldiscovery/index (from toolindex)
  - toolexec/run (from toolrun)

## PHASE 1: PRD-150 - Migrate toolset → toolcompose/set

1. Create set/ directory in toolcompose
2. Copy all .go files from ../toolset/ to set/
3. Update package name: `package toolset` → `package set`
4. Update imports:
   - github.com/jonwraymond/toolset → github.com/jonwraymond/toolcompose/set
   - github.com/jonwraymond/toolmodel → github.com/jonwraymond/toolfoundation/model
   - github.com/jonwraymond/toolindex → github.com/jonwraymond/tooldiscovery/index
5. Update type references: toolmodel. → model., toolindex. → index.
6. Update go.mod with dependencies and replace directives for ../toolfoundation and ../tooldiscovery
7. Run: GOWORK=off go mod tidy
8. Run: GOWORK=off go build ./set/...
9. Run: GOWORK=off go test ./set/... -cover
10. Commit with message: feat(set): migrate toolset to toolcompose/set package (PRD-150)

## PHASE 2: PRD-151 - Migrate toolskill → toolcompose/skill

1. Create skill/ directory in toolcompose
2. Copy all .go files from ../toolskill/ to skill/
3. Update package name: `package toolskill` → `package skill`
4. Update imports:
   - github.com/jonwraymond/toolskill → github.com/jonwraymond/toolcompose/skill
   - github.com/jonwraymond/toolmodel → github.com/jonwraymond/toolfoundation/model
   - github.com/jonwraymond/toolrun → github.com/jonwraymond/toolexec/run
5. Update type references accordingly
6. Add toolexec replace directive to go.mod
7. Run: GOWORK=off go mod tidy
8. Run: GOWORK=off go build ./skill/...
9. Run: GOWORK=off go test ./skill/... -cover
10. Commit with message: feat(skill): migrate toolskill to toolcompose/skill package (PRD-151)

## PHASE 3: Verify Complete Module

1. Run: GOWORK=off go build ./...
2. Run: GOWORK=off go test ./... -cover
3. Push all commits: git push

## Self-Correction
- If tests fail, read error carefully, identify root cause, fix, and re-run
- If import errors occur, check that all type references are updated
- If build fails, ensure go.mod has all required dependencies and replace directives

## Rules
- Use GOWORK=off for all go commands
- All commits must include: Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
- Commit after each successful phase

## Completion
When ALL phases complete and ALL tests pass:
<promise>COMPOSITION_LAYER_COMPLETE</promise>
