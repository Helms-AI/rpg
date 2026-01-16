# VS Code Continue Configuration

This directory contains the MCP server configuration for the Continue extension in VS Code.

## Automatic Setup

Run the setup script from the repository root:

```bash
./scripts/setup-clients.sh vscode-continue
```

## Manual Setup

1. Build the RPG binary:
   ```bash
   ./run.sh --build
   ```

2. Create the Continue config directory:
   ```bash
   mkdir -p .continue
   ```

3. Copy the configuration:
   ```bash
   cp clients/vscode-continue/config.json .continue/config.json
   ```

4. Open VS Code and the Continue extension will detect the MCP server.

## Configuration

The configuration uses VS Code variables:
- `${workspaceFolder}` - resolves to the current workspace root

```json
{
  "mcpServers": [
    {
      "name": "rpg",
      "command": "${workspaceFolder}/bin/rpg",
      "args": ["--output", "${workspaceFolder}/generated"]
    }
  ]
}
```

## Global Configuration

For global Continue configuration:

**macOS/Linux:**
```bash
cp clients/vscode-continue/config.json ~/.continue/config.json
# Update paths to absolute paths
```

**Windows:**
```powershell
copy clients\vscode-continue\config.json %USERPROFILE%\.continue\config.json
```

## Verification

1. Open VS Code with the Continue extension
2. Open a project with the RPG configuration
3. In the Continue chat, try: "List RPG supported languages"
