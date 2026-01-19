import { Link, useLocation } from 'react-router-dom';
import { ChevronRight } from 'lucide-react';
import { useState } from 'react';

interface NavItem {
  title: string;
  href?: string;
  items?: NavItem[];
}

const navigation: NavItem[] = [
  {
    title: 'Getting Started',
    items: [
      { title: 'Overview', href: '/getting-started' },
      { title: 'Installation', href: '/getting-started/installation' },
      { title: 'Quick Start', href: '/getting-started/quick-start' },
      { title: 'Configuration', href: '/getting-started/configuration' },
    ],
  },
  {
    title: 'Writing Specs',
    items: [
      { title: 'Overview', href: '/writing-specs' },
      { title: 'Examples', href: '/writing-specs/examples' },
      { title: 'Best Practices', href: '/writing-specs/best-practices' },
    ],
  },
  {
    title: 'Importing Specs',
    items: [
      { title: 'Overview', href: '/importing-specs' },
      { title: 'Examples', href: '/importing-specs/examples' },
      { title: 'Best Practices', href: '/importing-specs/best-practices' },
    ],
  },
  {
    title: 'Languages',
    items: [
      { title: 'Overview', href: '/languages' },
      { title: 'Go', href: '/languages/go' },
      { title: 'Rust', href: '/languages/rust' },
      { title: 'Java', href: '/languages/java' },
      { title: 'C#', href: '/languages/csharp' },
      { title: 'Python', href: '/languages/python' },
      { title: 'TypeScript', href: '/languages/typescript' },
    ],
  },
  {
    title: 'MCP Tools',
    items: [
      { title: 'Overview', href: '/mcp-tools' },
      { title: 'list_languages', href: '/mcp-tools/list-languages' },
      { title: 'parse_spec', href: '/mcp-tools/parse-spec' },
      { title: 'get_generation_context', href: '/mcp-tools/get-generation-context' },
      { title: 'get_project_structure', href: '/mcp-tools/get-project-structure' },
      { title: 'generate_source_from_spec', href: '/mcp-tools/generate-source-from-spec' },
      { title: 'import_spec_from_source', href: '/mcp-tools/import-spec-from-source' },
      { title: 'import_spec_from_github', href: '/mcp-tools/import-spec-from-github' },
      { title: 'deep_analyze_source', href: '/mcp-tools/deep-analyze-source' },
      { title: 'list_project_languages', href: '/mcp-tools/list-project-languages' },
      { title: 'get_files_for_language', href: '/mcp-tools/get-files-for-language' },
      { title: 'ensure_parity', href: '/mcp-tools/ensure-parity' },
      { title: 'semantic_parity_analysis', href: '/mcp-tools/semantic-parity-analysis' },
      { title: 'iterative_refinement_loop', href: '/mcp-tools/iterative-refinement-loop' },
    ],
  },
  {
    title: 'Client Setup',
    items: [
      { title: 'Overview', href: '/setup' },
      { title: 'Claude Desktop', href: '/setup/claude-desktop' },
      { title: 'Claude Code', href: '/setup/claude-code' },
      { title: 'VS Code Continue', href: '/setup/vscode-continue' },
      { title: 'VS Code Copilot', href: '/setup/vscode-copilot' },
      { title: 'Cursor', href: '/setup/cursor' },
      { title: 'Windsurf', href: '/setup/windsurf' },
      { title: 'Cline', href: '/setup/cline' },
    ],
  },
  {
    title: 'Tutorials',
    items: [
      { title: 'Overview', href: '/tutorials' },
      { title: 'Your First Spec', href: '/tutorials/first-spec' },
      { title: 'Multi-Language', href: '/tutorials/multi-language' },
      { title: 'REST API', href: '/tutorials/rest-api' },
      { title: 'Parity Checking', href: '/tutorials/parity-checking' },
    ],
  },
  {
    title: 'API Reference',
    items: [
      { title: 'Overview', href: '/api' },
      { title: 'Schemas', href: '/api/schemas' },
    ],
  },
];

function NavSection({ section }: { section: NavItem }) {
  const location = useLocation();
  const [isOpen, setIsOpen] = useState(() => {
    // Auto-open section if current path is in this section
    return section.items?.some((item) => location.pathname === item.href) ?? false;
  });

  const isActive = (href: string) => location.pathname === href;

  return (
    <div className="mb-4">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center justify-between w-full text-left px-3 py-2 text-sm font-semibold text-gray-900 dark:text-white hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg"
      >
        {section.title}
        <ChevronRight
          className={`h-4 w-4 text-gray-500 transition-transform ${isOpen ? 'rotate-90' : ''}`}
        />
      </button>
      {isOpen && section.items && (
        <div className="mt-1 ml-2 space-y-1">
          {section.items.map((item) => (
            <Link
              key={item.href}
              to={item.href!}
              className={`sidebar-link ${isActive(item.href!) ? 'active' : ''}`}
            >
              {item.title}
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}

export default function Sidebar() {
  return (
    <aside className="hidden lg:block w-64 flex-shrink-0 border-r border-gray-200 dark:border-gray-800">
      <div className="sticky top-16 h-[calc(100vh-4rem)] overflow-y-auto p-4">
        <nav>
          {navigation.map((section) => (
            <NavSection key={section.title} section={section} />
          ))}
        </nav>
      </div>
    </aside>
  );
}
