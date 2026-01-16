# Claude Desktop Configuration

This directory contains the MCP server configuration for Claude Desktop.

## Automatic Setup

Run the setup script from the repository root:

```bash
./scripts/setup-clients.sh claude-desktop
```

Or set up all clients at once:

```bash
./scripts/setup-clients.sh
```

## Manual Setup

1. Build the RPG binary:
   ```bash
   ./run.sh --build
   ```

2. Copy the configuration to your Claude Desktop config directory:

   **macOS:**
   ```bash
   cp clients/claude-desktop/config.json ~/Library/Application\ Support/Claude/claude_desktop_config.json
   ```

   **Windows:**
   ```powershell
   copy clients\claude-desktop\config.json %APPDATA%\Claude\claude_desktop_config.json
   ```

   **Linux:**
   ```bash
   cp clients/claude-desktop/config.json ~/.config/Claude/claude_desktop_config.json
   ```

3. Update the path in the config file to point to your RPG binary:
   ```json
   {
     "mcpServers": {
       "rpg": {
         "command": "/absolute/path/to/bin/rpg"
       }
     }
   }
   ```

4. Restart Claude Desktop

## Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `command` | Path to the RPG binary | Required |
| `args` | Command line arguments | `["--output", "./generated"]` |
| `alwaysAllow` | Tools that don't require confirmation | See config |

## Verification

After setup, you can verify the connection by asking Claude:
- "List available languages in RPG"
- "Validate my spec file"
