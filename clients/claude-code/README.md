# Claude Code Configuration

This directory contains the MCP server configuration for Claude Code (CLI).

## Automatic Setup

The `.mcp.json` file in the repository root is automatically detected by Claude Code.
No additional setup is required when working in this repository.

## Manual Setup

If you want to use RPG globally or in other projects:

1. Build the RPG binary:
   ```bash
   ./run.sh --build
   ```

2. Install globally:
   ```bash
   make install-global
   # or
   sudo cp bin/rpg /usr/local/bin/
   ```

3. Copy `.mcp.json` to your project:
   ```bash
   cp clients/claude-code/.mcp.json /path/to/your/project/
   ```

4. Update the path if needed:
   ```json
   {
     "mcpServers": {
       "rpg": {
         "command": "rpg"
       }
     }
   }
   ```

## Project-Level Configuration

For project-specific configuration, create `.mcp.json` in your project root:

```json
{
  "mcpServers": {
    "rpg": {
      "command": "/path/to/rpg",
      "args": ["--output", "./generated"]
    }
  }
}
```

## Verification

After setup, run Claude Code and try:
- "Use rpg to list supported languages"
- "Validate examples/reference-project-generator.spec.md"
