#!/usr/bin/env sh
set -eu

repo_dir="$(cd "$(dirname "$0")/.." && pwd)"
install_dir="$HOME/.curriculum/bin"
artifact="$repo_dir/dist/cur"

if [ ! -f "$artifact" ]; then
  echo "Build artifact not found: $artifact. Run scripts/build.sh first." >&2
  exit 1
fi

echo "Installing to $install_dir..."
mkdir -p "$install_dir"
cp "$artifact" "$install_dir/cur"
chmod +x "$install_dir/cur"
