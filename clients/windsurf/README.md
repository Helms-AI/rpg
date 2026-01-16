# Windsurf Configuration

This directory contains the MCP server configuration for Windsurf IDE.

## Automatic Setup

Run the setup script from the repository root:

```bash
./scripts/setup-clients.sh windsurf
```

## Manual Setup

1. Build the RPG binary:
   ```bash
   ./run.sh --build
   ```

2. Create the Windsurf config directory:
   ```bash
   mkdir -p .windsurf
   ```

3. Copy the configuration:
   ```bash
   cp clients/windsurf/mcp.json .windsurf/mcp.json
   ```

4. Restart Windsurf.

## Configuration

```json
{
  "mcpServers": {
    "rpg": {
      "command": "./bin/rpg",
      "args": ["--output", "./generated"],
      "transportType": "stdio"
    }
  }
}
```

## Global Configuration

For global Windsurf configuration:

**macOS:**
```bash
cp clients/windsurf/mcp.json ~/.windsurf/mcp.json
# Update paths to absolute paths
```

## Verification

1. Open Windsurf in the RPG repository
2. Open the Cascade AI panel
3. Try: "List available RPG languages"
