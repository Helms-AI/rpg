import { Link } from 'react-router-dom';
import { Monitor, Terminal, Code, Sparkles, Wind, Bot, Zap } from 'lucide-react';

const clients = [
  {
    name: 'Claude Desktop',
    description: 'Anthropic\'s official desktop application for Claude',
    icon: Monitor,
    path: '/setup/claude-desktop',
    configFile: 'claude_desktop_config.json',
  },
  {
    name: 'Claude Code',
    description: 'Claude\'s CLI tool for software engineering',
    icon: Terminal,
    path: '/setup/claude-code',
    configFile: '.mcp.json',
  },
  {
    name: 'VS Code Continue',
    description: 'Open-source AI code assistant extension',
    icon: Code,
    path: '/setup/vscode-continue',
    configFile: '.continue/config.json',
  },
  {
    name: 'Cursor',
    description: 'AI-first code editor built on VS Code',
    icon: Sparkles,
    path: '/setup/cursor',
    configFile: '.cursor/mcp.json',
  },
  {
    name: 'Windsurf',
    description: 'Codeium\'s AI-powered IDE',
    icon: Wind,
    path: '/setup/windsurf',
    configFile: '.windsurf/mcp.json',
  },
  {
    name: 'Cline',
    description: 'Autonomous coding agent for VS Code',
    icon: Bot,
    path: '/setup/cline',
    configFile: 'cline_mcp_settings.json',
  },
  {
    name: 'Zed',
    description: 'High-performance, multiplayer code editor',
    icon: Zap,
    path: '/setup/zed',
    configFile: '.zed/settings.json',
  },
];

export default function SetupIndex() {
  return (
    <div className="prose-docs">
      <h1>Client Setup</h1>
      <p className="lead">
        RPG works with all major AI coding assistants that support the Model Context Protocol (MCP).
        Choose your client below for setup instructions.
      </p>

      <div className="not-prose my-8 p-4 rounded-lg bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800">
        <h3 className="text-green-800 dark:text-green-300 font-semibold mb-2">
          Zero-Configuration Setup
        </h3>
        <p className="text-green-700 dark:text-green-400 text-sm">
          RPG includes pre-configured templates for all supported clients.
          Run the setup script to automatically configure your AI assistant:
        </p>
        <pre className="mt-3 p-3 bg-green-100 dark:bg-green-900/40 rounded text-sm font-mono text-green-800 dark:text-green-300">
          ./scripts/setup-clients.sh
        </pre>
      </div>

      <h2>Supported Clients</h2>
      <div className="not-prose mt-6 grid gap-4">
        {clients.map((client) => (
          <Link
            key={client.name}
            to={client.path}
            className="group flex items-start gap-4 rounded-lg border border-gray-200 dark:border-gray-700 p-4 hover:border-primary-500 dark:hover:border-primary-400 transition-colors"
          >
            <div className="flex-shrink-0 p-2 rounded-lg bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400 group-hover:bg-primary-100 dark:group-hover:bg-primary-900/30 group-hover:text-primary-600 dark:group-hover:text-primary-400 transition-colors">
              <client.icon className="w-5 h-5" />
            </div>
            <div className="flex-1">
              <div className="flex items-center justify-between">
                <h3 className="font-semibold text-gray-900 dark:text-white group-hover:text-primary-600 dark:group-hover:text-primary-400">
                  {client.name}
                </h3>
                <span className="text-xs font-mono text-gray-500 dark:text-gray-500">
                  {client.configFile}
                </span>
              </div>
              <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
                {client.description}
              </p>
            </div>
          </Link>
        ))}
      </div>

      <h2 className="mt-12">Manual Setup</h2>
      <p>
        If you prefer manual configuration, each client page includes the full configuration
        file content and step-by-step instructions for your platform.
      </p>

      <h2>Troubleshooting</h2>
      <p>Common issues and solutions:</p>
      <ul>
        <li>
          <strong>Binary not found</strong> - Run <code>make build</code> first to compile the RPG binary
        </li>
        <li>
          <strong>Permission denied</strong> - Ensure the binary is executable: <code>chmod +x bin/rpg</code>
        </li>
        <li>
          <strong>MCP not connecting</strong> - Check that the path in your config is absolute
        </li>
      </ul>
    </div>
  );
}
