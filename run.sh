#!/bin/bash

# run.sh - Start the rpg MCP server
#
# Usage:
#   ./run.sh           - Build and run the MCP server
#   ./run.sh --build   - Just build (don't run)
#   ./run.sh --help    - Show help

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY_NAME="rpg"
BUILD_DIR="${SCRIPT_DIR}/bin"
BINARY_PATH="${BUILD_DIR}/${BINARY_NAME}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
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

show_help() {
    echo "rpg - MCP Server"
    echo ""
    echo "Usage: ./run.sh [options]"
    echo ""
    echo "Options:"
    echo "  --build     Build the binary only (don't start server)"
    echo "  --clean     Clean build artifacts before building"
    echo "  --help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  ./run.sh              # Build and start the MCP server"
    echo "  ./run.sh --build      # Just build the binary"
    echo "  ./run.sh --clean      # Clean and rebuild"
    echo ""
    echo "The MCP server communicates over stdio (stdin/stdout)."
    echo "Configure Claude Code to use this server by adding to your MCP settings:"
    echo ""
    echo "  {"
    echo "    \"mcpServers\": {"
    echo "      \"rpg\": {"
    echo "        \"command\": \"${BINARY_PATH}\""
    echo "      }"
    echo "    }"
    echo "  }"
}

build() {
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
    log_info "Cleaned."
}

run_server() {
    if [[ ! -f "${BINARY_PATH}" ]]; then
        log_warn "Binary not found, building first..."
        build
    fi

    log_info "Starting MCP server..."
    log_info "Server communicates over stdio. Press Ctrl+C to stop."

    # Run the server
    exec "${BINARY_PATH}"
}

# Parse arguments
BUILD_ONLY=false
CLEAN_FIRST=false

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
    build
    exit 0
fi

build
run_server
