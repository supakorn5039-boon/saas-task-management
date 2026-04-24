#!/usr/bin/env bash
# ════════════════════════════════════════════════════════════════
#  Full pre-commit audit — checks the WHOLE project, not just
#  staged files. Use before pushing or opening a PR.
#  Run: bash scripts/pre-commit-check.sh
# ════════════════════════════════════════════════════════════════
set -e

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'
FAIL=0
pass() { echo -e "${GREEN}✓ PASS${NC}: $1"; }
fail() { echo -e "${RED}✗ FAIL${NC}: $1"; FAIL=1; }
warn() { echo -e "${YELLOW}⚠ WARN${NC}: $1"; }

REPO_ROOT="$(git rev-parse --show-toplevel)"
cd "$REPO_ROOT"

echo "══════════════════════════════════════════"
echo "  saas-task-management — Full Audit"
echo "══════════════════════════════════════════"

# ── 1. Backend ───────────────────────────────
echo ""
echo "▸ 1. Backend (Go)"
if [ -d backend ]; then
  pushd backend >/dev/null
  if go build ./src/... 2>&1; then pass "go build"; else fail "go build"; fi
  if go vet ./src/...   2>&1; then pass "go vet";   else fail "go vet";   fi
  if find ./src -name '*_test.go' -type f | grep -q .; then
    if go test ./src/... -count=1 2>&1 | tee /tmp/gotest.out | grep -qE '^FAIL|--- FAIL'; then
      fail "go test"
    else
      pass "go test"
    fi
  else
    warn "no *_test.go files — skipping go test"
  fi
  popd >/dev/null
else
  warn "backend/ not found — skipping"
fi

# ── 2. Frontend ──────────────────────────────
echo ""
echo "▸ 2. Frontend (React + TS)"
if [ -d frontend ]; then
  pushd frontend >/dev/null
  if [ ! -d node_modules ]; then
    warn "node_modules missing — running npm ci"
    npm ci --silent
  fi
  if npx --no-install eslint ./src --max-warnings=0; then pass "eslint"; else fail "eslint"; fi
  if npx --no-install tsc -b --noEmit;                  then pass "tsc type-check"; else fail "tsc type-check"; fi
  if npx --no-install prettier --check ./src;            then pass "prettier formatted"; else fail "prettier — run 'npm run format'"; fi
  if npm run build --silent;                             then pass "vite build"; else fail "vite build"; fi
  popd >/dev/null
else
  warn "frontend/ not found — skipping"
fi

# ── 3. Repo hygiene ──────────────────────────
echo ""
echo "▸ 3. Repo hygiene"
TODO_COUNT=$(grep -rEn 'TODO|FIXME' --include='*.go' --include='*.ts' --include='*.tsx' \
  backend/src frontend/src 2>/dev/null \
  | grep -v '_test.go' | grep -v '\.test\.' | grep -v 'PLANNED:' | wc -l | tr -d ' ')
if [ "$TODO_COUNT" = "0" ]; then pass "no TODO/FIXME"; else warn "$TODO_COUNT TODO/FIXME (use PLANNED: if intentional)"; fi

SECRETS=$(grep -rEn '(api[_-]?key|secret|password)\s*[:=]\s*["'\''`][^"'\''`]{8,}' \
  --include='*.go' --include='*.ts' --include='*.tsx' \
  backend/src frontend/src 2>/dev/null | wc -l | tr -d ' ')
if [ "$SECRETS" = "0" ]; then pass "no hardcoded-secret patterns"; else fail "$SECRETS suspicious secret literals"; fi

# ── Summary ──────────────────────────────────
echo ""
echo "══════════════════════════════════════════"
if [ "$FAIL" = "0" ]; then
  echo -e "${GREEN}  ALL CHECKS PASSED — OK to commit/push${NC}"
else
  echo -e "${RED}  SOME CHECKS FAILED — fix before pushing${NC}"
  exit 1
fi
