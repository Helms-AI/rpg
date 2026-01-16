import { Link } from 'react-router-dom';
import { ArrowRight, Download, Zap, Settings } from 'lucide-react';

const sections = [
  {
    title: 'Installation',
    description: 'Get RPG set up on your machine in minutes',
    href: '/getting-started/installation',
    icon: Download,
  },
  {
    title: 'Quick Start',
    description: 'Write your first spec and generate code',
    href: '/getting-started/quick-start',
    icon: Zap,
  },
  {
    title: 'Configuration',
    description: 'Configure RPG for your workflow',
    href: '/getting-started/configuration',
    icon: Settings,
  },
];

export default function GettingStarted() {
  return (
    <div className="prose-docs">
      <h1>Getting Started</h1>
      <p className="lead text-xl text-gray-600 dark:text-gray-400">
        Learn how to install RPG, write your first spec, and generate code in multiple languages.
      </p>

      <div className="grid gap-4 mt-8 not-prose">
        {sections.map((section) => (
          <Link
            key={section.href}
            to={section.href}
            className="group flex items-center gap-4 p-4 rounded-xl border border-gray-200 dark:border-gray-800 hover:border-primary-500/50 dark:hover:border-primary-500/50 transition-colors"
          >
            <div className="p-3 rounded-lg bg-primary-50 dark:bg-primary-900/20">
              <section.icon className="w-6 h-6 text-primary-600 dark:text-primary-400" />
            </div>
            <div className="flex-1">
              <h3 className="font-semibold text-gray-900 dark:text-white">{section.title}</h3>
              <p className="text-sm text-gray-600 dark:text-gray-400">{section.description}</p>
            </div>
            <ArrowRight className="w-5 h-5 text-gray-400 group-hover:text-primary-500 group-hover:translate-x-1 transition-all" />
          </Link>
        ))}
      </div>
    </div>
  );
}
