import { Link } from 'react-router-dom';
import { FileText, Languages, Globe, GitCompare } from 'lucide-react';

const tutorials = [
  {
    title: 'Your First Spec',
    description: 'Learn the basics of writing specification files for RPG',
    icon: FileText,
    path: '/tutorials/first-spec',
    difficulty: 'Beginner',
    time: '10 min',
  },
  {
    title: 'Multi-Language Generation',
    description: 'Generate the same project in multiple programming languages',
    icon: Languages,
    path: '/tutorials/multi-language',
    difficulty: 'Intermediate',
    time: '15 min',
  },
  {
    title: 'Building a REST API',
    description: 'Create a complete REST API from spec to implementation',
    icon: Globe,
    path: '/tutorials/rest-api',
    difficulty: 'Intermediate',
    time: '20 min',
  },
  {
    title: 'Feature Parity Checking',
    description: 'Ensure consistent behavior across language implementations',
    icon: GitCompare,
    path: '/tutorials/parity-checking',
    difficulty: 'Advanced',
    time: '15 min',
  },
];

const difficultyColors = {
  Beginner: 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400',
  Intermediate: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400',
  Advanced: 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400',
};

export default function TutorialsIndex() {
  return (
    <div className="prose-docs">
      <h1>Tutorials</h1>
      <p className="lead">
        Step-by-step guides to help you master RPG, from writing your first spec
        to maintaining multi-language codebases.
      </p>

      <div className="not-prose mt-8 grid gap-6">
        {tutorials.map((tutorial) => (
          <Link
            key={tutorial.title}
            to={tutorial.path}
            className="group block rounded-xl border border-gray-200 dark:border-gray-700 p-6 hover:border-primary-500 dark:hover:border-primary-400 hover:shadow-lg transition-all"
          >
            <div className="flex items-start gap-4">
              <div className="flex-shrink-0 p-3 rounded-lg bg-primary-100 dark:bg-primary-900/30 text-primary-600 dark:text-primary-400 group-hover:bg-primary-200 dark:group-hover:bg-primary-900/50 transition-colors">
                <tutorial.icon className="w-6 h-6" />
              </div>
              <div className="flex-1">
                <div className="flex items-center gap-3 mb-2">
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-white group-hover:text-primary-600 dark:group-hover:text-primary-400">
                    {tutorial.title}
                  </h3>
                  <span className={`px-2 py-0.5 rounded text-xs font-medium ${difficultyColors[tutorial.difficulty as keyof typeof difficultyColors]}`}>
                    {tutorial.difficulty}
                  </span>
                </div>
                <p className="text-gray-600 dark:text-gray-400">
                  {tutorial.description}
                </p>
                <p className="mt-2 text-sm text-gray-500 dark:text-gray-500">
                  Estimated time: {tutorial.time}
                </p>
              </div>
            </div>
          </Link>
        ))}
      </div>

      <h2 className="mt-12">Learning Path</h2>
      <p>
        We recommend following the tutorials in order if you're new to RPG:
      </p>
      <ol>
        <li>
          <strong>First Spec</strong> - Understand the spec format and basic generation
        </li>
        <li>
          <strong>Multi-Language</strong> - Learn to target multiple languages from one spec
        </li>
        <li>
          <strong>REST API</strong> - Build a complete, real-world project
        </li>
        <li>
          <strong>Parity Checking</strong> - Maintain consistency across implementations
        </li>
      </ol>

      <h2>Prerequisites</h2>
      <p>Before starting the tutorials, make sure you have:</p>
      <ul>
        <li>RPG installed and configured (<Link to="/getting-started/installation">Installation Guide</Link>)</li>
        <li>An AI assistant with MCP support (<Link to="/setup">Client Setup</Link>)</li>
        <li>Basic familiarity with at least one supported language</li>
      </ul>
    </div>
  );
}
