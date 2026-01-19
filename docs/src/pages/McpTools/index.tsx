import { Link } from 'react-router-dom';
import { List, FileText, Sparkles, FolderTree, GitCompare, FileCode, Github, Search, Layers, RefreshCw, Zap, Languages, Files } from 'lucide-react';

// Core Generation Tools
const coreTools = [
  {
    name: 'list_languages',
    description: 'List all supported languages with their conventions and idioms',
    icon: List,
    path: '/mcp-tools/list-languages',
  },
  {
    name: 'parse_spec',
    description: 'Read a markdown specification file and return its content',
    icon: FileText,
    path: '/mcp-tools/parse-spec',
  },
  {
    name: 'get_generation_context',
    description: 'Get spec + language conventions + prompt template for code generation',
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
    name: 'generate_source_from_spec',
    description: 'Autonomous code generation with automatic parity validation',
    icon: Zap,
    path: '/mcp-tools/generate-source-from-spec',
  },
];

// Import & Analysis Tools
const importTools = [
  {
    name: 'import_spec_from_source',
    description: 'Analyze local source code for AI-powered spec generation',
    icon: FileCode,
    path: '/mcp-tools/import-spec-from-source',
  },
  {
    name: 'import_spec_from_github',
    description: 'Clone and analyze a GitHub repository for spec generation',
    icon: Github,
    path: '/mcp-tools/import-spec-from-github',
  },
  {
    name: 'deep_analyze_source',
    description: 'AST-based semantic analysis (types, functions, call graphs)',
    icon: Search,
    path: '/mcp-tools/deep-analyze-source',
  },
  {
    name: 'list_project_languages',
    description: 'Detect all programming languages in a project',
    icon: Languages,
    path: '/mcp-tools/list-project-languages',
  },
  {
    name: 'get_files_for_language',
    description: 'Get raw file contents for AI-driven analysis',
    icon: Files,
    path: '/mcp-tools/get-files-for-language',
  },
];

// Parity & Refinement Tools
const parityTools = [
  {
    name: 'ensure_parity',
    description: 'Compare implementations across languages with fix instructions',
    icon: GitCompare,
    path: '/mcp-tools/ensure-parity',
  },
  {
    name: 'semantic_parity_analysis',
    description: 'Deep semantic comparison using AST-based analysis',
    icon: Layers,
    path: '/mcp-tools/semantic-parity-analysis',
  },
  {
    name: 'iterative_refinement_loop',
    description: 'Automated refinement until parity threshold is reached',
    icon: RefreshCw,
    path: '/mcp-tools/iterative-refinement-loop',
  },
];

function ToolCard({ tool }: { tool: typeof coreTools[0] }) {
  return (
    <Link
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
  );
}

export default function McpToolsIndex() {
  return (
    <div className="prose-docs">
      <h1>MCP Tools Reference</h1>
      <p className="lead">
        RPG exposes 13 MCP tools organized into three categories for code generation workflows.
      </p>

      <h2 className="mt-8">Core Generation</h2>
      <p>Essential tools for spec-based code generation.</p>
      <div className="not-prose mt-4 grid gap-4">
        {coreTools.map((tool) => (
          <ToolCard key={tool.name} tool={tool} />
        ))}
      </div>

      <h2 className="mt-10">Import & Analysis</h2>
      <p>Tools for analyzing existing codebases and generating specs.</p>
      <div className="not-prose mt-4 grid gap-4">
        {importTools.map((tool) => (
          <ToolCard key={tool.name} tool={tool} />
        ))}
      </div>

      <h2 className="mt-10">Parity & Refinement</h2>
      <p>Tools for comparing and refining multi-language implementations.</p>
      <div className="not-prose mt-4 grid gap-4">
        {parityTools.map((tool) => (
          <ToolCard key={tool.name} tool={tool} />
        ))}
      </div>

      <h2 className="mt-12">Common Workflows</h2>

      <h3 className="mt-6">Spec-First Generation</h3>
      <div className="not-prose mt-4 flex flex-col gap-2">
        <div className="flex items-center gap-3 text-sm">
          <span className="flex-shrink-0 w-6 h-6 rounded-full bg-primary-500 text-white flex items-center justify-center font-bold">1</span>
          <span className="text-gray-700 dark:text-gray-300">
            <code className="text-primary-600 dark:text-primary-400">list_languages</code> - Discover available target languages
          </span>
        </div>
        <div className="flex items-center gap-3 text-sm">
          <span className="flex-shrink-0 w-6 h-6 rounded-full bg-primary-500 text-white flex items-center justify-center font-bold">2</span>
          <span className="text-gray-700 dark:text-gray-300">
            <code className="text-primary-600 dark:text-primary-400">get_generation_context</code> - Get spec + conventions for generation
          </span>
        </div>
        <div className="flex items-center gap-3 text-sm">
          <span className="flex-shrink-0 w-6 h-6 rounded-full bg-primary-500 text-white flex items-center justify-center font-bold">3</span>
          <span className="text-gray-700 dark:text-gray-300">
            <code className="text-primary-600 dark:text-primary-400">generate_source_from_spec</code> - Autonomous code generation with validation
          </span>
        </div>
      </div>

      <h3 className="mt-6">Code-First Porting</h3>
      <div className="not-prose mt-4 flex flex-col gap-2">
        <div className="flex items-center gap-3 text-sm">
          <span className="flex-shrink-0 w-6 h-6 rounded-full bg-primary-500 text-white flex items-center justify-center font-bold">1</span>
          <span className="text-gray-700 dark:text-gray-300">
            <code className="text-primary-600 dark:text-primary-400">import_spec_from_github</code> - Analyze existing repository
          </span>
        </div>
        <div className="flex items-center gap-3 text-sm">
          <span className="flex-shrink-0 w-6 h-6 rounded-full bg-primary-500 text-white flex items-center justify-center font-bold">2</span>
          <span className="text-gray-700 dark:text-gray-300">
            <code className="text-primary-600 dark:text-primary-400">deep_analyze_source</code> - Extract types and functions via AST
          </span>
        </div>
        <div className="flex items-center gap-3 text-sm">
          <span className="flex-shrink-0 w-6 h-6 rounded-full bg-primary-500 text-white flex items-center justify-center font-bold">3</span>
          <span className="text-gray-700 dark:text-gray-300">
            <code className="text-primary-600 dark:text-primary-400">iterative_refinement_loop</code> - Generate and refine until target parity
          </span>
        </div>
      </div>

      <h3 className="mt-6">Multi-Language Parity</h3>
      <div className="not-prose mt-4 flex flex-col gap-2">
        <div className="flex items-center gap-3 text-sm">
          <span className="flex-shrink-0 w-6 h-6 rounded-full bg-primary-500 text-white flex items-center justify-center font-bold">1</span>
          <span className="text-gray-700 dark:text-gray-300">
            <code className="text-primary-600 dark:text-primary-400">ensure_parity</code> - Quick comparison across implementations
          </span>
        </div>
        <div className="flex items-center gap-3 text-sm">
          <span className="flex-shrink-0 w-6 h-6 rounded-full bg-primary-500 text-white flex items-center justify-center font-bold">2</span>
          <span className="text-gray-700 dark:text-gray-300">
            <code className="text-primary-600 dark:text-primary-400">semantic_parity_analysis</code> - Deep comparison with detailed scores
          </span>
        </div>
        <div className="flex items-center gap-3 text-sm">
          <span className="flex-shrink-0 w-6 h-6 rounded-full bg-primary-500 text-white flex items-center justify-center font-bold">3</span>
          <span className="text-gray-700 dark:text-gray-300">
            Apply fix instructions to achieve full parity
          </span>
        </div>
      </div>
    </div>
  );
}
