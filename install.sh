#!/usr/bin/env bash
set -euo pipefail

INSTALL_DIR="$HOME/.curriculum/bin"
REPO_DIR="$(cd "$(dirname "$0")" && pwd)"

sh "$REPO_DIR/scripts/build.sh"
sh "$REPO_DIR/scripts/install.sh"

# Add to PATH if not already present
add_to_path() {
    local profile="$1"
    if [ -f "$profile" ] && grep -q 'curriculum/bin' "$profile"; then
        return
    fi
    echo '' >> "$profile"
    echo '# curriculum CLI' >> "$profile"
    echo 'export PATH="$HOME/.curriculum/bin:$PATH"' >> "$profile"
    echo "Added to $profile"
}

if echo "$PATH" | tr ':' '\n' | grep -q "$INSTALL_DIR"; then
    echo "PATH already contains $INSTALL_DIR"
else
    shell_name="$(basename "${SHELL:-/bin/bash}")"
    case "$shell_name" in
        zsh)  add_to_path "$HOME/.zshrc" ;;
        bash)
            if [ -f "$HOME/.bash_profile" ]; then
                add_to_path "$HOME/.bash_profile"
            else
                add_to_path "$HOME/.bashrc"
            fi
            ;;
        fish)
            fish_conf="$HOME/.config/fish/conf.d/curriculum.fish"
            if [ ! -f "$fish_conf" ]; then
                mkdir -p "$(dirname "$fish_conf")"
                echo 'fish_add_path $HOME/.curriculum/bin' > "$fish_conf"
                echo "Added to $fish_conf"
            fi
            ;;
        *)
            echo "Unknown shell '$shell_name'. Add $INSTALL_DIR to your PATH manually."
            ;;
    esac
    echo "Restart your shell or run: export PATH=\"$INSTALL_DIR:\$PATH\""
fi

echo "Done. Run 'cur --help' to get started."
