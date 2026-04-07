#!/usr/bin/env bash
set -euo pipefail

VERSION=$(cat VERSION)
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS="-s -w -X curriculum/cmd.version=${VERSION} -X curriculum/cmd.commit=${COMMIT}"

echo "Building cur ${VERSION}+${COMMIT}..."
cd src
go install -ldflags "${LDFLAGS}" .
echo "Installed: $(which cur)"
