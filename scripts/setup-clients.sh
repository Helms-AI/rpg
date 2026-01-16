#!/bin/bash

# setup-clients.sh - Automatically configure MCP servers for AI coding assistants
#
# Usage:
#   ./scripts/setup-clients.sh              # Setup all detected clients
#   ./scripts/setup-clients.sh claude-desktop  # Setup specific client
#   ./scripts/setup-clients.sh --list       # List available clients

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
BINARY_PATH="${PROJECT_DIR}/bin/rpg"
CLIENTS_DIR="${PROJECT_DIR}/clients"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
log_client() { echo -e "${CYAN}[$1]${NC} $2"; }

# Supported clients
CLIENTS=(
    "claude-desktop"
    "claude-code"
    "vscode-continue"
    "cursor"
    "windsurf"
    "cline"
    "zed"
)

show_help() {
    echo "Setup MCP Server Configurations for AI Coding Assistants"
    echo ""
    echo "Usage: $0 [options] [client...]"
    echo ""
    echo "Options:"
    echo "  --list        List all supported clients"
    echo "  --detect      Detect installed clients"
    echo "  --help        Show this help message"
    echo ""
    echo "Supported clients:"
    for client in "${CLIENTS[@]}"; do
        echo "  - $client"
    done
    echo ""
    echo "Examples:"
    echo "  $0                    # Setup all detected clients"
    echo "  $0 claude-desktop     # Setup only Claude Desktop"
    echo "  $0 cursor windsurf    # Setup multiple clients"
}

list_clients() {
    echo "Supported MCP Clients:"
    echo ""
    for client in "${CLIENTS[@]}"; do
        local desc=""
        case $client in
            claude-desktop) desc="Claude Desktop App" ;;
            claude-code) desc="Claude Code CLI" ;;
            vscode-continue) desc="VS Code Continue Extension" ;;
            cursor) desc="Cursor IDE" ;;
            windsurf) desc="Windsurf IDE" ;;
            cline) desc="Cline VS Code Extension" ;;
            zed) desc="Zed Editor" ;;
        esac
        printf "  %-20s %s\n" "$client" "$desc"
    done
}

detect_clients() {
    local detected=()

    # Claude Desktop
    if [[ -d "$HOME/Library/Application Support/Claude" ]] || \
       [[ -d "$HOME/.config/Claude" ]] || \
       [[ -d "$APPDATA/Claude" ]]; then
        detected+=("claude-desktop")
    fi

    # Claude Code (always available if in this repo)
    detected+=("claude-code")

    # VS Code Continue
    if command -v code &> /dev/null || [[ -d "$HOME/.vscode" ]]; then
        detected+=("vscode-continue")
    fi

    # Cursor
    if command -v cursor &> /dev/null || [[ -d "$HOME/.cursor" ]]; then
        detected+=("cursor")
    fi

    # Windsurf
    if command -v windsurf &> /dev/null || [[ -d "$HOME/.windsurf" ]]; then
        detected+=("windsurf")
    fi

    # Cline (VS Code extension)
    if command -v code &> /dev/null; then
        detected+=("cline")
    fi

    # Zed
    if command -v zed &> /dev/null || [[ -d "$HOME/.config/zed" ]]; then
        detected+=("zed")
    fi

    echo "${detected[@]}"
}

ensure_binary() {
    if [[ ! -f "$BINARY_PATH" ]]; then
        log_warn "RPG binary not found at $BINARY_PATH"
        log_info "Building binary..."
        cd "$PROJECT_DIR"
        go build -o "$BINARY_PATH" ./cmd/rpg
        log_success "Binary built: $BINARY_PATH"
    fi
}

setup_claude_desktop() {
    log_client "claude-desktop" "Setting up Claude Desktop..."

    local config_dir=""
    local config_file=""

    # Detect OS and set config path
    if [[ "$OSTYPE" == "darwin"* ]]; then
        config_dir="$HOME/Library/Application Support/Claude"
    elif [[ "$OSTYPE" == "linux"* ]]; then
        config_dir="$HOME/.config/Claude"
    elif [[ "$OSTYPE" == "msys"* ]] || [[ "$OSTYPE" == "cygwin"* ]]; then
        config_dir="$APPDATA/Claude"
    else
        log_error "Unsupported OS for Claude Desktop"
        return 1
    fi

    config_file="$config_dir/claude_desktop_config.json"

    # Create directory if needed
    mkdir -p "$config_dir"

    # Generate config with absolute path
    cat > "$config_file" << EOF
{
  "mcpServers": {
    "rpg": {
      "command": "$BINARY_PATH",
      "args": ["--output", "$PROJECT_DIR/generated"],
      "alwaysAllow": [
        "list_languages",
        "parse_spec",
        "validate_spec",
        "get_generation_context",
        "get_project_structure"
      ]
    }
  }
}
EOF

    log_success "Claude Desktop configured: $config_file"
}

setup_claude_code() {
    log_client "claude-code" "Setting up Claude Code..."

    local config_file="$PROJECT_DIR/.mcp.json"

    # Generate config with relative path (project-level)
    cat > "$config_file" << EOF
{
  "mcpServers": {
    "rpg": {
      "command": "./bin/rpg",
      "args": []
    }
  }
}
EOF

    log_success "Claude Code configured: $config_file"
}

setup_vscode_continue() {
    log_client "vscode-continue" "Setting up VS Code Continue..."

    local config_dir="$PROJECT_DIR/.continue"
    local config_file="$config_dir/config.json"

    mkdir -p "$config_dir"

    cat > "$config_file" << EOF
{
  "models": [],
  "mcpServers": [
    {
      "name": "rpg",
      "command": "$BINARY_PATH",
      "args": ["--output", "$PROJECT_DIR/generated"]
    }
  ],
  "allowAnonymousTelemetry": false
}
EOF

    log_success "VS Code Continue configured: $config_file"
}

setup_cursor() {
    log_client "cursor" "Setting up Cursor..."

    local config_dir="$PROJECT_DIR/.cursor"
    local config_file="$config_dir/settings.json"

    mkdir -p "$config_dir"

    cat > "$config_file" << EOF
{
  "cursor.mcpServers": {
    "rpg": {
      "command": "$BINARY_PATH",
      "args": ["--output", "$PROJECT_DIR/generated"],
      "autoStart": false
    }
  }
}
EOF

    log_success "Cursor configured: $config_file"
}

setup_windsurf() {
    log_client "windsurf" "Setting up Windsurf..."

    local config_dir="$PROJECT_DIR/.windsurf"
    local config_file="$config_dir/mcp.json"

    mkdir -p "$config_dir"

    cat > "$config_file" << EOF
{
  "mcpServers": {
    "rpg": {
      "command": "$BINARY_PATH",
      "args": ["--output", "$PROJECT_DIR/generated"],
      "transportType": "stdio"
    }
  }
}
EOF

    log_success "Windsurf configured: $config_file"
}

setup_cline() {
    log_client "cline" "Setting up Cline..."

    local config_dir="$PROJECT_DIR/.vscode"
    local config_file="$config_dir/cline_mcp_settings.json"

    mkdir -p "$config_dir"

    cat > "$config_file" << EOF
{
  "mcpServers": {
    "rpg": {
      "command": "$BINARY_PATH",
      "args": ["--output", "$PROJECT_DIR/generated"],
      "disabled": false,
      "alwaysAllow": [
        "list_languages",
        "parse_spec",
        "validate_spec"
      ]
    }
  }
}
EOF

    log_success "Cline configured: $config_file"
}

setup_zed() {
    log_client "zed" "Setting up Zed..."

    local config_dir="$PROJECT_DIR/.zed"
    local config_file="$config_dir/settings.json"

    mkdir -p "$config_dir"

    cat > "$config_file" << EOF
{
  "context_servers": {
    "rpg": {
      "command": {
        "path": "$BINARY_PATH",
        "args": ["--output", "$PROJECT_DIR/generated"]
      }
    }
  }
}
EOF

    log_success "Zed configured: $config_file"
}

setup_client() {
    local client=$1

    case $client in
        claude-desktop) setup_claude_desktop ;;
        claude-code) setup_claude_code ;;
        vscode-continue) setup_vscode_continue ;;
        cursor) setup_cursor ;;
        windsurf) setup_windsurf ;;
        cline) setup_cline ;;
        zed) setup_zed ;;
        *)
            log_error "Unknown client: $client"
            return 1
            ;;
    esac
}

# Parse arguments
SPECIFIC_CLIENTS=()

while [[ $# -gt 0 ]]; do
    case $1 in
        --help|-h)
            show_help
            exit 0
            ;;
        --list)
            list_clients
            exit 0
            ;;
        --detect)
            echo "Detected clients:"
            for client in $(detect_clients); do
                echo "  - $client"
            done
            exit 0
            ;;
        *)
            SPECIFIC_CLIENTS+=("$1")
            shift
            ;;
    esac
done

# Main execution
echo ""
echo "╔══════════════════════════════════════════════════════════╗"
echo "║     RPG MCP Server - Client Configuration Setup          ║"
echo "╚══════════════════════════════════════════════════════════╝"
echo ""

# Ensure binary exists
ensure_binary

# Determine which clients to setup
if [[ ${#SPECIFIC_CLIENTS[@]} -gt 0 ]]; then
    CLIENTS_TO_SETUP=("${SPECIFIC_CLIENTS[@]}")
else
    log_info "Detecting installed clients..."
    CLIENTS_TO_SETUP=($(detect_clients))
fi

if [[ ${#CLIENTS_TO_SETUP[@]} -eq 0 ]]; then
    log_warn "No clients detected. Specify clients manually or install some first."
    echo ""
    list_clients
    exit 1
fi

echo ""
log_info "Setting up ${#CLIENTS_TO_SETUP[@]} client(s):"
for client in "${CLIENTS_TO_SETUP[@]}"; do
    echo "  - $client"
done
echo ""

# Setup each client
SUCCESS_COUNT=0
FAIL_COUNT=0

for client in "${CLIENTS_TO_SETUP[@]}"; do
    if setup_client "$client"; then
        ((SUCCESS_COUNT++))
    else
        ((FAIL_COUNT++))
    fi
done

# Summary
echo ""
echo "════════════════════════════════════════════════════════════"
echo ""
log_info "Setup complete!"
log_success "Configured: $SUCCESS_COUNT client(s)"
if [[ $FAIL_COUNT -gt 0 ]]; then
    log_warn "Failed: $FAIL_COUNT client(s)"
fi
echo ""
echo "Next steps:"
echo "  1. Restart your AI coding assistant"
echo "  2. Open the RPG project"
echo "  3. Try: 'Use RPG to list supported languages'"
echo ""
