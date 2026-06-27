#!/usr/bin/env bash
#
# deploy.sh — build and deploy the Bluejay CMS to the production server.
#
# What it does:
#   1. Cross-compiles a static linux/amd64 binary (pure Go, no CGO).
#   2. Syncs the binary + templates/ + db/migrations/ + public/ (CSS/JS) to the server.
#      -> It NEVER touches the production database (bluejay.db) or user-uploaded
#         files (public/uploads/), so your live content and images are preserved.
#   3. Restarts the systemd service and runs health checks.
#
# Usage:
#   ./deploy.sh              # build + deploy + verify
#   ./deploy.sh --no-build   # skip the build, deploy the existing ./bluejay-cms
#   ./deploy.sh --help
#
# Override any of these via environment variables if the server ever changes:
#   SSH_TARGET   (default: root@178.105.217.158)
#   REMOTE_DIR   (default: /var/www/bluejay-cms)
#   SERVICE      (default: bluejay-cms)
#   DOMAIN       (default: newsite.bluejayinnolabs.com)
#
set -euo pipefail

# ── Config ────────────────────────────────────────────────────────────────────
SSH_TARGET="${SSH_TARGET:-root@178.105.217.158}"
REMOTE_DIR="${REMOTE_DIR:-/var/www/bluejay-cms}"
SERVICE="${SERVICE:-bluejay-cms}"
DOMAIN="${DOMAIN:-newsite.bluejayinnolabs.com}"
SSH_OPTS=(-o StrictHostKeyChecking=accept-new -o BatchMode=yes)

# Always run from the repo root (the directory this script lives in).
cd "$(dirname "$0")"

DO_BUILD=1
for arg in "$@"; do
  case "$arg" in
    --no-build) DO_BUILD=0 ;;
    -h|--help)  sed -n '2,30p' "$0"; exit 0 ;;
    *) echo "Unknown option: $arg (try --help)"; exit 1 ;;
  esac
done

say() { printf '\n\033[1;36m==> %s\033[0m\n' "$1"; }
die() { printf '\n\033[1;31mDEPLOY FAILED: %s\033[0m\n' "$1" >&2; exit 1; }

# ── 1. Build ──────────────────────────────────────────────────────────────────
if [[ "$DO_BUILD" == 1 ]]; then
  say "Building static linux/amd64 binary"
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o bluejay-cms ./cmd/server \
    || die "go build failed"
fi
[[ -f bluejay-cms ]] || die "./bluejay-cms not found (run without --no-build)"

GITREF="$(git rev-parse --short HEAD 2>/dev/null || echo unknown)"
DIRTY="$(git status --porcelain 2>/dev/null | grep -q . && echo ' +uncommitted' || true)"
say "Deploying ${GITREF}${DIRTY} -> ${SSH_TARGET}:${REMOTE_DIR}  (https://${DOMAIN})"

# ── 2. Preflight: make sure rsync exists on the server ────────────────────────
ssh "${SSH_OPTS[@]}" "$SSH_TARGET" \
  'command -v rsync >/dev/null || { apt-get update -qq && apt-get install -y -qq rsync; }' \
  || die "could not reach server / install rsync"

# ── 3. Sync code + assets (preserve DB and user uploads) ──────────────────────
RSYNC=(rsync -az -e "ssh ${SSH_OPTS[*]}")

say "Uploading binary"
"${RSYNC[@]}" bluejay-cms "$SSH_TARGET:$REMOTE_DIR/bluejay-cms"        || die "binary upload failed"

say "Syncing templates/"
"${RSYNC[@]}" --delete templates/ "$SSH_TARGET:$REMOTE_DIR/templates/" || die "templates sync failed"

say "Syncing db/migrations/"
"${RSYNC[@]}" --delete db/migrations/ "$SSH_TARGET:$REMOTE_DIR/db/migrations/" || die "migrations sync failed"

say "Syncing public/ (excluding uploads/)"
"${RSYNC[@]}" --delete --exclude 'uploads/' public/ "$SSH_TARGET:$REMOTE_DIR/public/" || die "public sync failed"

# ── 4. Fix ownership + restart ────────────────────────────────────────────────
say "Setting ownership and restarting $SERVICE"
ssh "${SSH_OPTS[@]}" "$SSH_TARGET" "
  set -e
  chown -R www-data:www-data '$REMOTE_DIR'
  systemctl restart '$SERVICE'
  sleep 3
  systemctl is-active --quiet '$SERVICE' || { journalctl -u '$SERVICE' --no-pager -n 20; exit 1; }
  curl -fsS -o /dev/null http://127.0.0.1:28090/ || { echo 'app not responding on :28090'; exit 1; }
" || die "service did not come up cleanly (see logs above)"

# ── 5. Public health check over HTTPS ─────────────────────────────────────────
say "Verifying https://${DOMAIN}"
CODE="$(curl -fsS -o /dev/null -w '%{http_code}' "https://${DOMAIN}/" || true)"
[[ "$CODE" == "200" ]] || die "public site returned HTTP ${CODE:-000}"

printf '\n\033[1;32m✓ Deploy complete — https://%s is live (HTTP %s, %s)\033[0m\n' "$DOMAIN" "$CODE" "$GITREF"
