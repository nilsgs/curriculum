#!/usr/bin/env sh
set -eu

repo_dir="$(cd "$(dirname "$0")/.." && pwd)"
src_dir="$repo_dir/src"
dist_dir="$repo_dir/dist"
version="$(tr -d '\r\n' < "$repo_dir/VERSION")"
commit="$(git -C "$repo_dir" rev-parse --short HEAD 2>/dev/null || printf 'unknown')"
ldflags="-s -w -X curriculum/cmd.version=${version} -X curriculum/cmd.commit=${commit}"

echo "Cross-building cur v${version}+${commit}..."
mkdir -p "$dist_dir"

cd "$src_dir"
for spec in \
  "linux amd64 cur-linux-amd64" \
  "linux arm64 cur-linux-arm64" \
  "darwin amd64 cur-darwin-amd64" \
  "darwin arm64 cur-darwin-arm64" \
  "windows amd64 cur-windows-amd64.exe" \
  "windows arm64 cur-windows-arm64.exe"
do
  set -- $spec
  GOOS="$1" GOARCH="$2" go build -ldflags "$ldflags" -o "$dist_dir/$3" .
done
