# Cline Configuration

This directory contains the MCP server configuration for the Cline extension in VS Code.

## Automatic Setup

Run the setup script from the repository root:

```bash
./scripts/setup-clients.sh cline
```

## Manual Setup

1. Build the RPG binary:
   ```bash
   ./run.sh --build
   ```

2. Create the VS Code config directory:
   ```bash
   mkdir -p .vscode
   ```

3. Copy the configuration:
   ```bash
   cp clients/cline/cline_mcp_settings.json .vscode/cline_mcp_settings.json
   ```

4. Reload VS Code.

## Configuration

```json
{
  "mcpServers": {
    "rpg": {
      "command": "./bin/rpg",
      "args": ["--output", "./generated"],
      "disabled": false,
      "alwaysAllow": [
        "list_languages",
        "parse_spec",
        "validate_spec"
      ]
    }
  }
}
```

| Option | Description | Default |
|--------|-------------|---------|
| `command` | Path to the RPG binary | Required |
| `args` | Command line arguments | `[]` |
| `disabled` | Disable the server | `false` |
| `alwaysAllow` | Tools that don't need approval | `[]` |

## Global Configuration

For global Cline configuration, the settings file location depends on your OS:

**macOS:**
```
~/Library/Application Support/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json
```

**Windows:**
```
%APPDATA%\Code\User\globalStorage\saoudrizwan.claude-dev\settings\cline_mcp_settings.json
```

**Linux:**
```
~/.config/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json
```

## Verification

1. Open VS Code with the Cline extension
2. Start a new Cline session
3. Try: "Use RPG to list supported languages"
