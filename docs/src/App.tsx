import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { ThemeProvider } from './context/ThemeContext';
import Layout from './components/Layout/Layout';
import Home from './pages/Home';

// Lazy load pages for code splitting
import { lazy, Suspense } from 'react';

// Getting Started
const GettingStartedIndex = lazy(() => import('./pages/GettingStarted'));
const Installation = lazy(() => import('./pages/GettingStarted/Installation'));
const QuickStart = lazy(() => import('./pages/GettingStarted/QuickStart'));
const Configuration = lazy(() => import('./pages/GettingStarted/Configuration'));

// Writing Specs
const WritingSpecsIndex = lazy(() => import('./pages/WritingSpecs'));
const WritingSpecsExamples = lazy(() => import('./pages/WritingSpecs/Examples'));
const WritingSpecsBestPractices = lazy(() => import('./pages/WritingSpecs/BestPractices'));

// Importing Specs
const ImportingSpecsIndex = lazy(() => import('./pages/ImportingSpecs'));
const ImportingSpecsExamples = lazy(() => import('./pages/ImportingSpecs/Examples'));
const ImportingSpecsBestPractices = lazy(() => import('./pages/ImportingSpecs/BestPractices'));

// Languages
const LanguagesIndex = lazy(() => import('./pages/Languages'));
const GoLang = lazy(() => import('./pages/Languages/Go'));
const RustLang = lazy(() => import('./pages/Languages/Rust'));
const JavaLang = lazy(() => import('./pages/Languages/Java'));
const CSharpLang = lazy(() => import('./pages/Languages/CSharp'));
const PythonLang = lazy(() => import('./pages/Languages/Python'));
const TypeScriptLang = lazy(() => import('./pages/Languages/TypeScript'));

// MCP Tools
const McpToolsIndex = lazy(() => import('./pages/McpTools'));
const ListLanguages = lazy(() => import('./pages/McpTools/ListLanguages'));
const ParseSpec = lazy(() => import('./pages/McpTools/ParseSpec'));
const GetGenerationContext = lazy(() => import('./pages/McpTools/GetGenerationContext'));
const GetProjectStructure = lazy(() => import('./pages/McpTools/GetProjectStructure'));
const GenerateSourceFromSpec = lazy(() => import('./pages/McpTools/GenerateSourceFromSpec'));
const ImportSpecFromSource = lazy(() => import('./pages/McpTools/ImportSpecFromSource'));
const ImportSpecFromGithub = lazy(() => import('./pages/McpTools/ImportSpecFromGithub'));
const DeepAnalyzeSource = lazy(() => import('./pages/McpTools/DeepAnalyzeSource'));
const ListProjectLanguages = lazy(() => import('./pages/McpTools/ListProjectLanguages'));
const GetFilesForLanguage = lazy(() => import('./pages/McpTools/GetFilesForLanguage'));
const EnsureParity = lazy(() => import('./pages/McpTools/EnsureParity'));
const SemanticParityAnalysis = lazy(() => import('./pages/McpTools/SemanticParityAnalysis'));
const IterativeRefinementLoop = lazy(() => import('./pages/McpTools/IterativeRefinementLoop'));

// Setup
const SetupIndex = lazy(() => import('./pages/Setup'));
const ClaudeDesktop = lazy(() => import('./pages/Setup/ClaudeDesktop'));
const ClaudeCode = lazy(() => import('./pages/Setup/ClaudeCode'));
const VSCodeContinue = lazy(() => import('./pages/Setup/VSCodeContinue'));
const VSCodeCopilot = lazy(() => import('./pages/Setup/VSCodeCopilot'));
const Cursor = lazy(() => import('./pages/Setup/Cursor'));
const Windsurf = lazy(() => import('./pages/Setup/Windsurf'));
const Cline = lazy(() => import('./pages/Setup/Cline'));

// Tutorials
const TutorialsIndex = lazy(() => import('./pages/Tutorials'));
const FirstSpec = lazy(() => import('./pages/Tutorials/FirstSpec'));
const MultiLanguage = lazy(() => import('./pages/Tutorials/MultiLanguage'));
const RestApi = lazy(() => import('./pages/Tutorials/RestApi'));
const ParityChecking = lazy(() => import('./pages/Tutorials/ParityChecking'));

// API
const ApiIndex = lazy(() => import('./pages/Api'));
const Schemas = lazy(() => import('./pages/Api/Schemas'));

// Loading component
const PageLoader = () => (
  <div className="flex items-center justify-center min-h-[50vh]">
    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"></div>
  </div>
);

function App() {
  return (
    <ThemeProvider>
      <BrowserRouter>
        <Suspense fallback={<PageLoader />}>
          <Routes>
            {/* Home - no layout */}
            <Route path="/" element={<Home />} />

            {/* Documentation routes with layout */}
            <Route element={<Layout />}>
              {/* Getting Started */}
              <Route path="/getting-started" element={<GettingStartedIndex />} />
              <Route path="/getting-started/installation" element={<Installation />} />
              <Route path="/getting-started/quick-start" element={<QuickStart />} />
              <Route path="/getting-started/configuration" element={<Configuration />} />

              {/* Writing Specs */}
              <Route path="/writing-specs" element={<WritingSpecsIndex />} />
              <Route path="/writing-specs/examples" element={<WritingSpecsExamples />} />
              <Route path="/writing-specs/best-practices" element={<WritingSpecsBestPractices />} />

              {/* Importing Specs */}
              <Route path="/importing-specs" element={<ImportingSpecsIndex />} />
              <Route path="/importing-specs/examples" element={<ImportingSpecsExamples />} />
              <Route path="/importing-specs/best-practices" element={<ImportingSpecsBestPractices />} />

              {/* Languages */}
              <Route path="/languages" element={<LanguagesIndex />} />
              <Route path="/languages/go" element={<GoLang />} />
              <Route path="/languages/rust" element={<RustLang />} />
              <Route path="/languages/java" element={<JavaLang />} />
              <Route path="/languages/csharp" element={<CSharpLang />} />
              <Route path="/languages/python" element={<PythonLang />} />
              <Route path="/languages/typescript" element={<TypeScriptLang />} />

              {/* MCP Tools */}
              <Route path="/mcp-tools" element={<McpToolsIndex />} />
              <Route path="/mcp-tools/list-languages" element={<ListLanguages />} />
              <Route path="/mcp-tools/parse-spec" element={<ParseSpec />} />
              <Route path="/mcp-tools/get-generation-context" element={<GetGenerationContext />} />
              <Route path="/mcp-tools/get-project-structure" element={<GetProjectStructure />} />
              <Route path="/mcp-tools/generate-source-from-spec" element={<GenerateSourceFromSpec />} />
              <Route path="/mcp-tools/import-spec-from-source" element={<ImportSpecFromSource />} />
              <Route path="/mcp-tools/import-spec-from-github" element={<ImportSpecFromGithub />} />
              <Route path="/mcp-tools/deep-analyze-source" element={<DeepAnalyzeSource />} />
              <Route path="/mcp-tools/list-project-languages" element={<ListProjectLanguages />} />
              <Route path="/mcp-tools/get-files-for-language" element={<GetFilesForLanguage />} />
              <Route path="/mcp-tools/ensure-parity" element={<EnsureParity />} />
              <Route path="/mcp-tools/semantic-parity-analysis" element={<SemanticParityAnalysis />} />
              <Route path="/mcp-tools/iterative-refinement-loop" element={<IterativeRefinementLoop />} />

              {/* Setup */}
              <Route path="/setup" element={<SetupIndex />} />
              <Route path="/setup/claude-desktop" element={<ClaudeDesktop />} />
              <Route path="/setup/claude-code" element={<ClaudeCode />} />
              <Route path="/setup/vscode-continue" element={<VSCodeContinue />} />
              <Route path="/setup/vscode-copilot" element={<VSCodeCopilot />} />
              <Route path="/setup/cursor" element={<Cursor />} />
              <Route path="/setup/windsurf" element={<Windsurf />} />
              <Route path="/setup/cline" element={<Cline />} />

              {/* Tutorials */}
              <Route path="/tutorials" element={<TutorialsIndex />} />
              <Route path="/tutorials/first-spec" element={<FirstSpec />} />
              <Route path="/tutorials/multi-language" element={<MultiLanguage />} />
              <Route path="/tutorials/rest-api" element={<RestApi />} />
              <Route path="/tutorials/parity-checking" element={<ParityChecking />} />

              {/* API */}
              <Route path="/api" element={<ApiIndex />} />
              <Route path="/api/schemas" element={<Schemas />} />
            </Route>
          </Routes>
        </Suspense>
      </BrowserRouter>
    </ThemeProvider>
  );
}

export default App;
