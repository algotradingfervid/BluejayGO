#!/bin/bash
set -euo pipefail

PROJECT_DIR="/Users/narendhupati/Documents/ClaudeWebsiteCreator"
PLANS_DIR="$PROJECT_DIR/plans"
TRACKER="$PROJECT_DIR/automation/phase-tracker.json"
LOG_DIR="$PROJECT_DIR/automation/logs"

PHASES=(01 02 03 04 05 06 07 08 09 10 11 12 13 14 15 16 17 18 19 20 21 22)

# Derive plan filename from phase number by matching plans directory
get_plan_file() {
  local num="$1"
  basename "$(ls "$PLANS_DIR"/${num}-*.md 2>/dev/null | head -1)"
}

mkdir -p "$LOG_DIR"

for phase_num in "${PHASES[@]}"; do
  # Skip completed phases
  status=$(jq -r ".phases.\"$phase_num\".status // \"pending\"" "$TRACKER")
  if [ "$status" = "completed" ]; then
    echo "[SKIP] Phase $phase_num already completed"
    continue
  fi

  echo "=========================================="
  echo "[START] Phase $phase_num: $(get_plan_file "$phase_num")"
  echo "=========================================="

  # Update tracker: in_progress
  jq ".current_phase = \"$phase_num\" | .phases.\"$phase_num\".status = \"in_progress\" | .phases.\"$phase_num\".started_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" \
    "$TRACKER" > "${TRACKER}.tmp" && mv "${TRACKER}.tmp" "$TRACKER"

  # Read plan content
  plan_content=$(cat "$PLANS_DIR/$(get_plan_file "$phase_num")")

  # Build prompt — Claude handles build verification + error fixing internally
  prompt="You are implementing Phase $phase_num of the admin panel redesign.

## Phase Plan
$plan_content

## CRITICAL: Build Verification Loop
After implementing all changes:
1. Run \`go build ./cmd/...\` to check if the project compiles
2. If the build FAILS: read the errors, fix them, and run \`go build ./cmd/...\` again
3. Repeat until the build succeeds (max 5 attempts)
4. If SQL changes were made: run \`sqlc generate\` first, then build
5. Once the build passes, you are done with this phase
6. Do NOT move on or finish until the build passes

Begin implementation now. Follow the plan precisely."

  # Run Claude with fresh session (context is clean)
  claude -p "$prompt" \
    --system-prompt-file "$PROJECT_DIR/.claude/system-prompts/phase-runner.md" \
    --allowedTools "Edit,Write,Read,Glob,Grep,Bash(go build:*),Bash(go run:*),Bash(sqlc:*),Bash(ls:*),Bash(mkdir:*),Bash(cat:*)" \
    --max-turns 80 \
    --output-format json \
    2>&1 | tee "$LOG_DIR/phase-${phase_num}.log"

  exit_code=${PIPESTATUS[0]}

  if [ $exit_code -ne 0 ]; then
    echo "[FAIL] Phase $phase_num session failed (exit $exit_code)"
    jq ".phases.\"$phase_num\".status = \"failed\"" "$TRACKER" > "${TRACKER}.tmp" && mv "${TRACKER}.tmp" "$TRACKER"
    echo "Fix manually and re-run orchestrator."
    exit 1
  fi

  # External build check (safety net — Claude should have already fixed this)
  cd "$PROJECT_DIR"
  if ! go build ./cmd/... 2>/dev/null; then
    echo "[WARN] Build still broken after Phase $phase_num. Running fix session..."

    fix_prompt="The Go project does not compile after Phase $phase_num implementation.
Run \`go build ./cmd/...\` to see the errors, then fix them.
Keep fixing until the build passes."

    claude -p "$fix_prompt" \
      --allowedTools "Edit,Write,Read,Glob,Grep,Bash(go build:*),Bash(sqlc:*),Bash(ls:*)" \
      --max-turns 30 \
      --output-format json \
      2>&1 | tee "$LOG_DIR/phase-${phase_num}-fix.log"

    # Final check
    if ! go build ./cmd/... 2>/dev/null; then
      echo "[FAIL] Could not fix build for Phase $phase_num"
      jq ".phases.\"$phase_num\".status = \"failed\"" "$TRACKER" > "${TRACKER}.tmp" && mv "${TRACKER}.tmp" "$TRACKER"
      exit 1
    fi
  fi

  # Mark completed
  jq ".phases.\"$phase_num\".status = \"completed\" | .phases.\"$phase_num\".completed_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" \
    "$TRACKER" > "${TRACKER}.tmp" && mv "${TRACKER}.tmp" "$TRACKER"

  echo "[DONE] Phase $phase_num completed successfully"
  sleep 2
done

echo "All phases completed!"
