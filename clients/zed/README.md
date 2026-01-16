# Zed Configuration

This directory contains the MCP server configuration for Zed Editor.

## Automatic Setup

Run the setup script from the repository root:

```bash
./scripts/setup-clients.sh zed
```

## Manual Setup

1. Build the RPG binary:
   ```bash
   ./run.sh --build
   ```

2. Create the Zed config directory:
   ```bash
   mkdir -p .zed
   ```

3. Copy the configuration:
   ```bash
   cp clients/zed/settings.json .zed/settings.json
   ```

4. Restart Zed.

## Configuration

Zed uses "context_servers" for MCP integration:

```json
{
  "context_servers": {
    "rpg": {
      "command": {
        "path": "./bin/rpg",
        "args": ["--output", "./generated"]
      }
    }
  }
}
```

## Global Configuration

For global Zed configuration:

**macOS:**
```bash
# Edit ~/.config/zed/settings.json
# Add the context_servers section with absolute paths
```

**Linux:**
```bash
# Edit ~/.config/zed/settings.json
```

## Verification

1. Open Zed in the RPG repository
2. Open the AI assistant panel (Cmd/Ctrl + ?)
3. The RPG context server should be available
