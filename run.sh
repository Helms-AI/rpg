#!/bin/bash

# run.sh - Start the rpg MCP server and documentation website
#
# Usage:
#   ./run.sh           - Build and run the MCP server + docs server
#   ./run.sh --build   - Just build (don't run)
#   ./run.sh --docs    - Only start the docs server (port 4000)
#   ./run.sh --mcp     - Only start the MCP server
#   ./run.sh --help    - Show help

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY_NAME="rpg"
BUILD_DIR="${SCRIPT_DIR}/bin"
BINARY_PATH="${BUILD_DIR}/${BINARY_NAME}"
DOCS_DIR="${SCRIPT_DIR}/docs"
DOCS_DIST="${DOCS_DIR}/dist"
DOCS_PORT=4000

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_docs() {
    echo -e "${CYAN}[DOCS]${NC} $1"
}

show_help() {
    echo "rpg - MCP Server & Documentation"
    echo ""
    echo "Usage: ./run.sh [options]"
    echo ""
    echo "Options:"
    echo "  --build     Build the binary and docs (don't start servers)"
    echo "  --docs      Only start the docs server on port ${DOCS_PORT}"
    echo "  --mcp       Only start the MCP server (no docs)"
    echo "  --clean     Clean build artifacts before building"
    echo "  --dev       Start docs in development mode (hot reload)"
    echo "  --help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  ./run.sh              # Build and start both MCP server and docs"
    echo "  ./run.sh --docs       # Just start the documentation server"
    echo "  ./run.sh --build      # Build everything without starting"
    echo "  ./run.sh --dev        # Start docs in development mode"
    echo ""
    echo "Servers:"
    echo "  - MCP Server: Communicates over stdio (stdin/stdout)"
    echo "  - Docs Server: http://localhost:${DOCS_PORT}"
    echo ""
    echo "MCP Configuration for Claude Code:"
    echo ""
    echo "  {"
    echo "    \"mcpServers\": {"
    echo "      \"rpg\": {"
    echo "        \"command\": \"${BINARY_PATH}\""
    echo "      }"
    echo "    }"
    echo "  }"
}

check_node() {
    if ! command -v node &> /dev/null; then
        log_error "Node.js is not installed. Please install Node.js to build/serve docs."
        exit 1
    fi
}

check_npm() {
    if ! command -v npm &> /dev/null; then
        log_error "npm is not installed. Please install npm to build/serve docs."
        exit 1
    fi
}

install_docs_deps() {
    if [[ ! -d "${DOCS_DIR}/node_modules" ]]; then
        log_docs "Installing documentation dependencies..."
        cd "${DOCS_DIR}"
        npm install
        cd "${SCRIPT_DIR}"
    fi
}

build_docs() {
    check_node
    check_npm
    install_docs_deps

    log_docs "Building documentation site..."
    cd "${DOCS_DIR}"
    npm run build
    cd "${SCRIPT_DIR}"
    log_docs "Documentation built: ${DOCS_DIST}"
}

build_binary() {
    log_info "Building rpg..."

    # Create build directory
    mkdir -p "${BUILD_DIR}"

    # Build the binary
    cd "${SCRIPT_DIR}"
    go build -o "${BINARY_PATH}" ./cmd/rpg

    log_info "Built: ${BINARY_PATH}"
}

clean() {
    log_info "Cleaning build artifacts..."
    rm -rf "${BUILD_DIR}"
    rm -rf "${DOCS_DIST}"
    log_info "Cleaned."
}

start_docs_server() {
    check_node

    if [[ ! -d "${DOCS_DIST}" ]]; then
        log_warn "Docs not built, building first..."
        build_docs
    fi

    log_docs "Starting documentation server on http://localhost:${DOCS_PORT}"
    log_docs "Press Ctrl+C to stop."

    # Use npx serve to serve the built docs
    cd "${DOCS_DIR}"
    npx serve -s dist -l ${DOCS_PORT}
}

start_docs_dev() {
    check_node
    check_npm
    install_docs_deps

    log_docs "Starting documentation in development mode..."
    log_docs "Hot reload enabled at http://localhost:5173"
    log_docs "Press Ctrl+C to stop."

    cd "${DOCS_DIR}"
    npm run dev -- --port 5173
}

start_mcp_server() {
    if [[ ! -f "${BINARY_PATH}" ]]; then
        log_warn "Binary not found, building first..."
        build_binary
    fi

    log_info "Starting MCP server..."
    log_info "Server communicates over stdio. Press Ctrl+C to stop."

    # Run the server
    exec "${BINARY_PATH}"
}

start_both_servers() {
    if [[ ! -f "${BINARY_PATH}" ]]; then
        log_warn "Binary not found, building first..."
        build_binary
    fi

    if [[ ! -d "${DOCS_DIST}" ]]; then
        log_warn "Docs not built, building first..."
        build_docs
    fi

    log_info "Starting both servers..."
    log_docs "Documentation: http://localhost:${DOCS_PORT}"
    log_info "MCP Server: stdio"
    echo ""

    # Start docs server in background
    cd "${DOCS_DIR}"
    npx serve -s dist -l ${DOCS_PORT} &
    DOCS_PID=$!
    cd "${SCRIPT_DIR}"

    # Trap to cleanup background process
    trap "kill $DOCS_PID 2>/dev/null; exit" SIGINT SIGTERM EXIT

    # Give docs server a moment to start
    sleep 1

    log_docs "Docs server running (PID: $DOCS_PID)"
    log_info "Starting MCP server..."

    # Run MCP server in foreground
    "${BINARY_PATH}"
}

# Parse arguments
BUILD_ONLY=false
CLEAN_FIRST=false
DOCS_ONLY=false
MCP_ONLY=false
DEV_MODE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --build)
            BUILD_ONLY=true
            shift
            ;;
        --clean)
            CLEAN_FIRST=true
            shift
            ;;
        --docs)
            DOCS_ONLY=true
            shift
            ;;
        --mcp)
            MCP_ONLY=true
            shift
            ;;
        --dev)
            DEV_MODE=true
            shift
            ;;
        --help|-h)
            show_help
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Main execution
if [[ "${CLEAN_FIRST}" == "true" ]]; then
    clean
fi

if [[ "${BUILD_ONLY}" == "true" ]]; then
    build_binary
    build_docs
    exit 0
fi

if [[ "${DEV_MODE}" == "true" ]]; then
    start_docs_dev
    exit 0
fi

if [[ "${DOCS_ONLY}" == "true" ]]; then
    start_docs_server
    exit 0
fi

if [[ "${MCP_ONLY}" == "true" ]]; then
    build_binary
    start_mcp_server
    exit 0
fi

# Default: build and start both
build_binary
build_docs
start_both_servers
