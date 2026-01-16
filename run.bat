@echo off
setlocal EnableDelayedExpansion

:: run.bat - Start the rpg MCP server and documentation website (Windows)
::
:: Usage:
::   run.bat           - Build and run the MCP server + docs server
::   run.bat --build   - Just build (don't run)
::   run.bat --docs    - Only start the docs server (port 4000)
::   run.bat --mcp     - Only start the MCP server
::   run.bat --help    - Show help

set "SCRIPT_DIR=%~dp0"
set "SCRIPT_DIR=%SCRIPT_DIR:~0,-1%"
set "BINARY_NAME=rpg.exe"
set "BUILD_DIR=%SCRIPT_DIR%\bin"
set "BINARY_PATH=%BUILD_DIR%\%BINARY_NAME%"
set "DOCS_DIR=%SCRIPT_DIR%\docs"
set "DOCS_DIST=%DOCS_DIR%\dist"
set "DOCS_PORT=4000"

:: Parse arguments
set "BUILD_ONLY=false"
set "CLEAN_FIRST=false"
set "DOCS_ONLY=false"
set "MCP_ONLY=false"
set "DEV_MODE=false"

:parse_args
if "%~1"=="" goto :main
if /i "%~1"=="--build" (
    set "BUILD_ONLY=true"
    shift
    goto :parse_args
)
if /i "%~1"=="--clean" (
    set "CLEAN_FIRST=true"
    shift
    goto :parse_args
)
if /i "%~1"=="--docs" (
    set "DOCS_ONLY=true"
    shift
    goto :parse_args
)
if /i "%~1"=="--mcp" (
    set "MCP_ONLY=true"
    shift
    goto :parse_args
)
if /i "%~1"=="--dev" (
    set "DEV_MODE=true"
    shift
    goto :parse_args
)
if /i "%~1"=="--help" goto :show_help
if /i "%~1"=="-h" goto :show_help
echo [ERROR] Unknown option: %~1
goto :show_help

:main
:: Main execution
if "%CLEAN_FIRST%"=="true" call :clean

if "%BUILD_ONLY%"=="true" (
    call :build_binary
    call :build_docs
    goto :eof
)

if "%DEV_MODE%"=="true" (
    call :start_docs_dev
    goto :eof
)

if "%DOCS_ONLY%"=="true" (
    call :start_docs_server
    goto :eof
)

if "%MCP_ONLY%"=="true" (
    call :build_binary
    call :start_mcp_server
    goto :eof
)

:: Default: build and start both
call :build_binary
call :build_docs
call :start_both_servers
goto :eof

:: ============================================================================
:: Functions
:: ============================================================================

:log_info
echo [INFO] %~1
goto :eof

:log_warn
echo [WARN] %~1
goto :eof

:log_error
echo [ERROR] %~1
goto :eof

:log_docs
echo [DOCS] %~1
goto :eof

:show_help
echo rpg - MCP Server ^& Documentation
echo.
echo Usage: run.bat [options]
echo.
echo Options:
echo   --build     Build the binary and docs (don't start servers)
echo   --docs      Only start the docs server on port %DOCS_PORT%
echo   --mcp       Only start the MCP server (no docs)
echo   --clean     Clean build artifacts before building
echo   --dev       Start docs in development mode (hot reload)
echo   --help      Show this help message
echo.
echo Examples:
echo   run.bat              # Build and start both MCP server and docs
echo   run.bat --docs       # Just start the documentation server
echo   run.bat --build      # Build everything without starting
echo   run.bat --dev        # Start docs in development mode
echo.
echo Servers:
echo   - MCP Server: Communicates over stdio (stdin/stdout)
echo   - Docs Server: http://localhost:%DOCS_PORT%
echo.
echo MCP Configuration for Claude Code:
echo.
echo   {
echo     "mcpServers": {
echo       "rpg": {
echo         "command": "%BINARY_PATH:\=/%"
echo       }
echo     }
echo   }
goto :eof

:check_node
where node >nul 2>nul
if %errorlevel% neq 0 (
    call :log_error "Node.js is not installed. Please install Node.js to build/serve docs."
    exit /b 1
)
goto :eof

:check_npm
where npm >nul 2>nul
if %errorlevel% neq 0 (
    call :log_error "npm is not installed. Please install npm to build/serve docs."
    exit /b 1
)
goto :eof

:check_go
where go >nul 2>nul
if %errorlevel% neq 0 (
    call :log_error "Go is not installed. Please install Go to build the binary."
    exit /b 1
)
goto :eof

:install_docs_deps
if not exist "%DOCS_DIR%\node_modules" (
    call :log_docs "Installing documentation dependencies..."
    pushd "%DOCS_DIR%"
    call npm install
    popd
)
goto :eof

:build_docs
call :check_node
call :check_npm
call :install_docs_deps

call :log_docs "Building documentation site..."
pushd "%DOCS_DIR%"
call npm run build
popd
call :log_docs "Documentation built: %DOCS_DIST%"
goto :eof

:build_binary
call :check_go
call :log_info "Building rpg..."

:: Create build directory
if not exist "%BUILD_DIR%" mkdir "%BUILD_DIR%"

:: Build the binary
pushd "%SCRIPT_DIR%"
go build -o "%BINARY_PATH%" ./cmd/rpg
popd

call :log_info "Built: %BINARY_PATH%"
goto :eof

:clean
call :log_info "Cleaning build artifacts..."
if exist "%BUILD_DIR%" rmdir /s /q "%BUILD_DIR%"
if exist "%DOCS_DIST%" rmdir /s /q "%DOCS_DIST%"
call :log_info "Cleaned."
goto :eof

:start_docs_server
call :check_node

if not exist "%DOCS_DIST%" (
    call :log_warn "Docs not built, building first..."
    call :build_docs
)

call :log_docs "Starting documentation server on http://localhost:%DOCS_PORT%"
call :log_docs "Press Ctrl+C to stop."

:: Use npx serve to serve the built docs
pushd "%DOCS_DIR%"
call npx serve -s dist -l %DOCS_PORT%
popd
goto :eof

:start_docs_dev
call :check_node
call :check_npm
call :install_docs_deps

call :log_docs "Starting documentation in development mode..."
call :log_docs "Hot reload enabled at http://localhost:5173"
call :log_docs "Press Ctrl+C to stop."

pushd "%DOCS_DIR%"
call npm run dev -- --port 5173
popd
goto :eof

:start_mcp_server
if not exist "%BINARY_PATH%" (
    call :log_warn "Binary not found, building first..."
    call :build_binary
)

call :log_info "Starting MCP server..."
call :log_info "Server communicates over stdio. Press Ctrl+C to stop."

:: Run the server
"%BINARY_PATH%"
goto :eof

:start_both_servers
if not exist "%BINARY_PATH%" (
    call :log_warn "Binary not found, building first..."
    call :build_binary
)

if not exist "%DOCS_DIST%" (
    call :log_warn "Docs not built, building first..."
    call :build_docs
)

call :log_info "Starting both servers..."
call :log_docs "Documentation: http://localhost:%DOCS_PORT%"
call :log_info "MCP Server: stdio"
echo.

:: Start docs server in a new window
start "RPG Docs Server" cmd /c "cd /d "%DOCS_DIR%" && npx serve -s dist -l %DOCS_PORT%"

:: Give docs server a moment to start
timeout /t 2 /nobreak >nul

call :log_docs "Docs server running in separate window"
call :log_info "Starting MCP server..."
call :log_info "Close the docs server window manually when done."

:: Run MCP server in foreground
"%BINARY_PATH%"
goto :eof

:eof
endlocal
