# Cursor Configuration

This directory contains the MCP server configuration for Cursor IDE.

## Automatic Setup

Run the setup script from the repository root:

```bash
./scripts/setup-clients.sh cursor
```

## Manual Setup

1. Build the RPG binary:
   ```bash
   ./run.sh --build
   ```

2. Create the Cursor config directory:
   ```bash
   mkdir -p .cursor
   ```

3. Copy the configuration:
   ```bash
   cp clients/cursor/settings.json .cursor/settings.json
   ```

4. Restart Cursor IDE.

## Configuration Options

```json
{
  "cursor.mcpServers": {
    "rpg": {
      "command": "./bin/rpg",
      "args": ["--output", "./generated"],
      "autoStart": false
    }
  }
}
```

| Option | Description | Default |
|--------|-------------|---------|
| `command` | Path to the RPG binary | Required |
| `args` | Command line arguments | `[]` |
| `autoStart` | Auto-start server on IDE launch | `false` |

## Global Configuration

For global Cursor MCP configuration, add to your Cursor settings:

1. Open Cursor Settings (Cmd/Ctrl + ,)
2. Search for "MCP Servers"
3. Add the RPG configuration

Or edit `~/.cursor/settings.json` directly.

## Verification

1. Open Cursor in the RPG repository
2. Open the AI chat panel
3. Try: "Use RPG to list supported languages"
