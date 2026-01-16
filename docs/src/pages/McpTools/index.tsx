import { Link } from 'react-router-dom';
import { List, FileText, CheckCircle, Sparkles, FolderTree, GitCompare, FileCode } from 'lucide-react';

const tools = [
  {
    name: 'list_languages',
    description: 'List all supported languages with their conventions and idioms',
    icon: List,
    path: '/mcp-tools/list-languages',
  },
  {
    name: 'parse_spec',
    description: 'Parse a markdown specification file and return its structured content',
    icon: FileText,
    path: '/mcp-tools/parse-spec',
  },
  {
    name: 'validate_spec',
    description: 'Check if a spec file exists and contains valid content',
    icon: CheckCircle,
    path: '/mcp-tools/validate-spec',
  },
  {
    name: 'get_generation_context',
    description: 'Get full context for code generation including spec and language conventions',
    icon: Sparkles,
    path: '/mcp-tools/get-generation-context',
  },
  {
    name: 'get_project_structure',
    description: 'Get recommended file structure for a project in the target language',
    icon: FolderTree,
    path: '/mcp-tools/get-project-structure',
  },
  {
    name: 'ensure_parity',
    description: 'Check feature parity across generated projects and provide fix instructions',
    icon: GitCompare,
    path: '/mcp-tools/ensure-parity',
  },
  {
    name: 'import_spec_from_source',
    description: 'Analyze source code for AI-powered spec generation',
    icon: FileCode,
    path: '/mcp-tools/import-spec-from-source',
  },
];

export default function McpToolsIndex() {
  return (
    <div className="prose-docs">
      <h1>MCP Tools Reference</h1>
      <p className="lead">
        RPG exposes 7 MCP tools that AI assistants can use for code generation workflows.
        Each tool is designed for a specific step in the generation pipeline.
      </p>

      <div className="not-prose mt-8 grid gap-4">
        {tools.map((tool) => (
          <Link
            key={tool.name}
            to={tool.path}
            className="group flex items-start gap-4 rounded-lg border border-gray-200 dark:border-gray-700 p-4 hover:border-primary-500 dark:hover:border-primary-400 transition-colors"
          >
            <div className="flex-shrink-0 p-2 rounded-lg bg-primary-100 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400 group-hover:bg-primary-200 dark:group-hover:bg-primary-900/50 transition-colors">
              <tool.icon className="w-5 h-5" />
            </div>
            <div>
              <h3 className="font-mono text-sm font-semibold text-gray-900 dark:text-white group-hover:text-primary-600 dark:group-hover:text-primary-400">
                {tool.name}
              </h3>
              <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
                {tool.description}
              </p>
            </div>
          </Link>
        ))}
      </div>

      <h2 className="mt-12">Tool Pipeline</h2>
      <p>
        The tools are designed to work together in a typical code generation workflow:
      </p>

      <div className="not-prose mt-6 flex flex-col gap-2">
        <div className="flex items-center gap-3 text-sm">
          <span className="flex-shrink-0 w-6 h-6 rounded-full bg-primary-500 text-white flex items-center justify-center font-bold">1</span>
          <span className="text-gray-700 dark:text-gray-300">
            <code className="text-primary-600 dark:text-primary-400">list_languages</code> - Discover available target languages
          </span>
        </div>
        <div className="flex items-center gap-3 text-sm">
          <span className="flex-shrink-0 w-6 h-6 rounded-full bg-primary-500 text-white flex items-center justify-center font-bold">2</span>
          <span className="text-gray-700 dark:text-gray-300">
            <code className="text-primary-600 dark:text-primary-400">validate_spec</code> - Verify spec file exists and is valid
          </span>
        </div>
        <div className="flex items-center gap-3 text-sm">
          <span className="flex-shrink-0 w-6 h-6 rounded-full bg-primary-500 text-white flex items-center justify-center font-bold">3</span>
          <span className="text-gray-700 dark:text-gray-300">
            <code className="text-primary-600 dark:text-primary-400">get_generation_context</code> - Get spec + conventions for generation
          </span>
        </div>
        <div className="flex items-center gap-3 text-sm">
          <span className="flex-shrink-0 w-6 h-6 rounded-full bg-primary-500 text-white flex items-center justify-center font-bold">4</span>
          <span className="text-gray-700 dark:text-gray-300">
            <code className="text-primary-600 dark:text-primary-400">get_project_structure</code> - Get recommended file layout
          </span>
        </div>
        <div className="flex items-center gap-3 text-sm">
          <span className="flex-shrink-0 w-6 h-6 rounded-full bg-primary-500 text-white flex items-center justify-center font-bold">5</span>
          <span className="text-gray-700 dark:text-gray-300">
            <code className="text-primary-600 dark:text-primary-400">ensure_parity</code> - Verify implementations match across languages
          </span>
        </div>
      </div>
    </div>
  );
}
