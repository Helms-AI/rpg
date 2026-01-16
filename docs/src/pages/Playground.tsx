import { useState } from 'react';
import { Play, Copy, Check, Download } from 'lucide-react';

const defaultSpec = `# Hello World

A simple greeting function to demonstrate RPG.

## Target Languages

- go
- typescript
- python

## Functions

### greet

Generates a personalized greeting message.

**Accepts:**
- \`name\` (string): The name to greet

**Returns:**
- string: The greeting message

**Logic:**
1. If name is empty, use "World"
2. Return "Hello, {name}!"

## Tests

### Default greeting
- **Given**: greet("")
- **Expect**: "Hello, World!"

### Named greeting
- **Given**: greet("Alice")
- **Expect**: "Hello, Alice!"
`;

const languages = [
  { id: 'go', name: 'Go', color: 'bg-cyan-500' },
  { id: 'rust', name: 'Rust', color: 'bg-orange-500' },
  { id: 'typescript', name: 'TypeScript', color: 'bg-blue-500' },
  { id: 'python', name: 'Python', color: 'bg-yellow-500' },
  { id: 'java', name: 'Java', color: 'bg-red-500' },
  { id: 'csharp', name: 'C#', color: 'bg-purple-500' },
];

export default function Playground() {
  const [spec, setSpec] = useState(defaultSpec);
  const [selectedLanguage, setSelectedLanguage] = useState('go');
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    await navigator.clipboard.writeText(spec);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="min-h-[calc(100vh-4rem)] flex flex-col">
      <div className="border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 px-4 py-3">
        <div className="flex items-center justify-between">
          <h1 className="text-xl font-bold text-gray-900 dark:text-white">
            Playground
          </h1>
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <span className="text-sm text-gray-600 dark:text-gray-400">Target:</span>
              <select
                value={selectedLanguage}
                onChange={(e) => setSelectedLanguage(e.target.value)}
                className="px-3 py-1.5 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              >
                {languages.map((lang) => (
                  <option key={lang.id} value={lang.id}>
                    {lang.name}
                  </option>
                ))}
              </select>
            </div>
            <button
              onClick={handleCopy}
              className="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg border border-gray-300 dark:border-gray-600 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
            >
              {copied ? <Check className="w-4 h-4" /> : <Copy className="w-4 h-4" />}
              {copied ? 'Copied!' : 'Copy Spec'}
            </button>
            <button className="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg border border-gray-300 dark:border-gray-600 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors">
              <Download className="w-4 h-4" />
              Download
            </button>
          </div>
        </div>
      </div>

      <div className="flex-1 grid md:grid-cols-2 divide-x divide-gray-200 dark:divide-gray-700">
        {/* Editor Panel */}
        <div className="flex flex-col">
          <div className="px-4 py-2 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50">
            <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
              Specification
            </span>
          </div>
          <div className="flex-1 relative">
            <textarea
              value={spec}
              onChange={(e) => setSpec(e.target.value)}
              className="absolute inset-0 w-full h-full p-4 font-mono text-sm bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 resize-none focus:outline-none"
              placeholder="Write your spec here..."
              spellCheck={false}
            />
          </div>
        </div>

        {/* Preview Panel */}
        <div className="flex flex-col bg-gray-50 dark:bg-gray-800/30">
          <div className="px-4 py-2 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50 flex items-center justify-between">
            <div className="flex items-center gap-2">
              <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                Generated Code
              </span>
              <span className={`px-2 py-0.5 rounded text-xs font-medium text-white ${languages.find(l => l.id === selectedLanguage)?.color}`}>
                {languages.find(l => l.id === selectedLanguage)?.name}
              </span>
            </div>
            <button className="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg bg-primary-600 text-white text-sm font-medium hover:bg-primary-700 transition-colors">
              <Play className="w-4 h-4" />
              Generate
            </button>
          </div>
          <div className="flex-1 p-4">
            <div className="h-full flex items-center justify-center text-gray-500 dark:text-gray-400">
              <div className="text-center">
                <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-gray-200 dark:bg-gray-700 flex items-center justify-center">
                  <Play className="w-8 h-8" />
                </div>
                <p className="text-lg font-medium mb-2">Ready to Generate</p>
                <p className="text-sm max-w-xs">
                  Click "Generate" to create {languages.find(l => l.id === selectedLanguage)?.name} code
                  from your specification.
                </p>
                <p className="text-xs mt-4 text-gray-400 dark:text-gray-500">
                  Note: In this demo, generation requires an MCP-connected AI assistant.
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Language Tabs */}
      <div className="border-t border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 px-4 py-2">
        <div className="flex items-center gap-2">
          <span className="text-xs text-gray-500 dark:text-gray-400 mr-2">Quick switch:</span>
          {languages.map((lang) => (
            <button
              key={lang.id}
              onClick={() => setSelectedLanguage(lang.id)}
              className={`px-3 py-1 rounded text-xs font-medium transition-colors ${
                selectedLanguage === lang.id
                  ? `${lang.color} text-white`
                  : 'bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-700'
              }`}
            >
              {lang.name}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}
