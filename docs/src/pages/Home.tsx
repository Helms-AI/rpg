import { Link } from 'react-router-dom';
import { useState, useEffect } from 'react';
import {
  ArrowRight,
  Sparkles,
  Code2,
  GitCompare,
  FileCode,
  Zap,
  Check,
  Layers,
  RefreshCw,
} from 'lucide-react';
import Header from '../components/Layout/Header';

// Language data with distinctive colors
const languages = [
  { name: 'Go', color: '#00ADD8', bg: 'bg-[#00ADD8]/10', border: 'border-[#00ADD8]/30' },
  { name: 'Rust', color: '#CE422B', bg: 'bg-[#CE422B]/10', border: 'border-[#CE422B]/30' },
  { name: 'Java', color: '#F89820', bg: 'bg-[#F89820]/10', border: 'border-[#F89820]/30' },
  { name: 'C#', color: '#953DAC', bg: 'bg-[#953DAC]/10', border: 'border-[#953DAC]/30' },
  { name: 'Python', color: '#3776AB', bg: 'bg-[#3776AB]/10', border: 'border-[#3776AB]/30' },
  { name: 'TypeScript', color: '#3178C6', bg: 'bg-[#3178C6]/10', border: 'border-[#3178C6]/30' },
];

const features = [
  {
    icon: Sparkles,
    title: 'AI-Powered Generation',
    description: 'Write natural language specs, get idiomatic code. The AI understands context, patterns, and best practices.',
  },
  {
    icon: Layers,
    title: 'Multi-Language Support',
    description: 'One spec, six languages. Generate consistent implementations in Go, Rust, Java, C#, Python, and TypeScript.',
  },
  {
    icon: GitCompare,
    title: 'Feature Parity Checking',
    description: 'Ensure all implementations match. Detect missing features and get fix instructions automatically.',
  },
  {
    icon: RefreshCw,
    title: 'Reverse Engineering',
    description: 'Import existing code to generate specs. Perfect for documenting legacy systems or porting to new languages.',
  },
];

const steps = [
  { step: '01', title: 'Write a spec', desc: 'Describe your project in markdown' },
  { step: '02', title: 'Choose languages', desc: 'Select your target languages' },
  { step: '03', title: 'Generate code', desc: 'Get idiomatic implementations' },
];

// Animated code samples for the demo
const specCode = `# slugify

Convert text to URL-friendly slugs.

## Target Languages
- go
- typescript

## Functions
### slugify
**accepts:** text: Text
**returns:** Text
**logic:**
  - lowercase the text
  - replace spaces with dashes
  - remove special characters`;

const goCode = `package slugify

import (
    "regexp"
    "strings"
)

func Slugify(text string) string {
    // Convert to lowercase
    result := strings.ToLower(text)

    // Replace spaces with dashes
    result = strings.ReplaceAll(result, " ", "-")

    // Remove special characters
    reg := regexp.MustCompile("[^a-z0-9-]")
    result = reg.ReplaceAllString(result, "")

    return result
}`;

const tsCode = `export function slugify(text: string): string {
  return text
    .toLowerCase()
    .replace(/\\s+/g, '-')
    .replace(/[^a-z0-9-]/g, '');
}`;

function TypewriterCode({ code, delay = 0 }: { code: string; delay?: number }) {
  const [displayed, setDisplayed] = useState('');
  const [started, setStarted] = useState(false);

  useEffect(() => {
    const startTimer = setTimeout(() => setStarted(true), delay);
    return () => clearTimeout(startTimer);
  }, [delay]);

  useEffect(() => {
    if (!started) return;

    let index = 0;
    const interval = setInterval(() => {
      if (index <= code.length) {
        setDisplayed(code.slice(0, index));
        index++;
      } else {
        clearInterval(interval);
      }
    }, 15);

    return () => clearInterval(interval);
  }, [code, started]);

  return (
    <pre className="text-xs leading-relaxed whitespace-pre-wrap">
      {displayed}
      {displayed.length < code.length && started && (
        <span className="inline-block w-2 h-4 bg-emerald-400 animate-pulse ml-0.5" />
      )}
    </pre>
  );
}

function CodeDemo() {
  const [activeTab, setActiveTab] = useState<'spec' | 'go' | 'ts'>('spec');
  const [isAnimating, setIsAnimating] = useState(false);

  const runDemo = () => {
    setIsAnimating(true);
    setActiveTab('spec');
    setTimeout(() => setActiveTab('go'), 2000);
    setTimeout(() => setActiveTab('ts'), 4000);
    setTimeout(() => setIsAnimating(false), 6000);
  };

  useEffect(() => {
    // Auto-run demo on mount
    const timer = setTimeout(runDemo, 1000);
    return () => clearTimeout(timer);
  }, []);

  return (
    <div className="relative">
      {/* Glow effect */}
      <div className="absolute -inset-4 bg-gradient-to-r from-emerald-500/20 via-cyan-500/20 to-violet-500/20 rounded-3xl blur-2xl opacity-50" />

      <div className="relative bg-gray-900 rounded-2xl border border-gray-800 overflow-hidden shadow-2xl">
        {/* Window chrome */}
        <div className="flex items-center gap-2 px-4 py-3 bg-gray-800/50 border-b border-gray-700/50">
          <div className="flex gap-1.5">
            <div className="w-3 h-3 rounded-full bg-red-500/80" />
            <div className="w-3 h-3 rounded-full bg-yellow-500/80" />
            <div className="w-3 h-3 rounded-full bg-green-500/80" />
          </div>
          <div className="flex-1 flex justify-center gap-1">
            <button
              onClick={() => setActiveTab('spec')}
              className={`px-3 py-1 text-xs rounded-md transition-all ${
                activeTab === 'spec'
                  ? 'bg-emerald-500/20 text-emerald-400'
                  : 'text-gray-500 hover:text-gray-300'
              }`}
            >
              spec.md
            </button>
            <button
              onClick={() => setActiveTab('go')}
              className={`px-3 py-1 text-xs rounded-md transition-all ${
                activeTab === 'go'
                  ? 'bg-[#00ADD8]/20 text-[#00ADD8]'
                  : 'text-gray-500 hover:text-gray-300'
              }`}
            >
              slugify.go
            </button>
            <button
              onClick={() => setActiveTab('ts')}
              className={`px-3 py-1 text-xs rounded-md transition-all ${
                activeTab === 'ts'
                  ? 'bg-[#3178C6]/20 text-[#3178C6]'
                  : 'text-gray-500 hover:text-gray-300'
              }`}
            >
              slugify.ts
            </button>
          </div>
          <button
            onClick={runDemo}
            disabled={isAnimating}
            className="p-1.5 rounded-md text-gray-500 hover:text-gray-300 hover:bg-gray-700/50 disabled:opacity-50"
          >
            <RefreshCw className={`w-4 h-4 ${isAnimating ? 'animate-spin' : ''}`} />
          </button>
        </div>

        {/* Code content */}
        <div className="p-6 h-80 overflow-hidden font-mono">
          {activeTab === 'spec' && (
            <div className="text-emerald-400/90">
              <TypewriterCode code={specCode} />
            </div>
          )}
          {activeTab === 'go' && (
            <div className="text-[#00ADD8]/90">
              <TypewriterCode code={goCode} delay={0} />
            </div>
          )}
          {activeTab === 'ts' && (
            <div className="text-[#3178C6]/90">
              <TypewriterCode code={tsCode} delay={0} />
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default function Home() {
  return (
    <div className="min-h-screen bg-gray-950 text-white">
      <Header />

      {/* Hero Section */}
      <section className="relative overflow-hidden">
        {/* Animated gradient background */}
        <div className="absolute inset-0 bg-gradient-to-br from-gray-950 via-gray-900 to-gray-950" />
        <div className="absolute inset-0">
          <div className="absolute top-0 left-1/4 w-96 h-96 bg-emerald-500/10 rounded-full blur-3xl animate-pulse" />
          <div className="absolute bottom-0 right-1/4 w-96 h-96 bg-cyan-500/10 rounded-full blur-3xl animate-pulse" style={{ animationDelay: '1s' }} />
          <div className="absolute top-1/2 left-1/2 w-96 h-96 bg-violet-500/5 rounded-full blur-3xl animate-pulse" style={{ animationDelay: '2s' }} />
        </div>

        {/* Grid pattern overlay */}
        <div
          className="absolute inset-0 opacity-[0.02]"
          style={{
            backgroundImage: `linear-gradient(rgba(255,255,255,0.1) 1px, transparent 1px),
                              linear-gradient(90deg, rgba(255,255,255,0.1) 1px, transparent 1px)`,
            backgroundSize: '64px 64px',
          }}
        />

        <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-20 pb-32">
          <div className="grid lg:grid-cols-2 gap-16 items-center">
            {/* Left: Text content */}
            <div className="space-y-8">
              {/* Badge */}
              <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-sm font-medium">
                <Sparkles className="w-4 h-4" />
                <span>AI-Powered Code Generation</span>
              </div>

              {/* Headline */}
              <h1 className="text-5xl sm:text-6xl lg:text-7xl font-bold tracking-tight">
                <span className="block text-white">Write specs in</span>
                <span className="block bg-gradient-to-r from-emerald-400 via-cyan-400 to-violet-400 bg-clip-text text-transparent">
                  markdown.
                </span>
                <span className="block text-gray-400 text-4xl sm:text-5xl lg:text-6xl mt-2">
                  Generate code in 6 languages.
                </span>
              </h1>

              {/* Description */}
              <p className="text-lg text-gray-400 max-w-xl leading-relaxed">
                RPG transforms natural language specifications into production-ready, idiomatic code.
                Describe once, deploy everywhere.
              </p>

              {/* CTAs */}
              <div className="flex flex-wrap gap-4">
                <Link
                  to="/getting-started"
                  className="group inline-flex items-center gap-2 px-6 py-3 bg-white text-gray-900 rounded-xl font-semibold hover:bg-gray-100 transition-all shadow-lg shadow-white/10"
                >
                  Get Started
                  <ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
                </Link>
              </div>

              {/* Language badges */}
              <div className="flex flex-wrap gap-2 pt-4">
                {languages.map((lang) => (
                  <span
                    key={lang.name}
                    className={`px-3 py-1.5 rounded-lg text-sm font-medium ${lang.bg} border ${lang.border}`}
                    style={{ color: lang.color }}
                  >
                    {lang.name}
                  </span>
                ))}
              </div>
            </div>

            {/* Right: Code demo */}
            <div className="lg:pl-8">
              <CodeDemo />
            </div>
          </div>
        </div>
      </section>

      {/* Quick Start Section */}
      <section className="relative py-24 bg-gray-900/50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
              Get started in minutes
            </h2>
            <p className="text-gray-400 text-lg max-w-2xl mx-auto">
              From spec to working code in three simple steps
            </p>
          </div>

          <div className="grid md:grid-cols-3 gap-8">
            {steps.map((item, index) => (
              <div
                key={item.step}
                className="relative group"
                style={{ animationDelay: `${index * 150}ms` }}
              >
                {/* Connector line */}
                {index < steps.length - 1 && (
                  <div className="hidden md:block absolute top-12 left-full w-full h-px bg-gradient-to-r from-gray-700 to-transparent z-0" />
                )}

                <div className="relative bg-gray-800/50 rounded-2xl p-8 border border-gray-700/50 hover:border-emerald-500/30 transition-all group-hover:bg-gray-800/80">
                  <span className="text-6xl font-bold text-gray-700/50 absolute top-4 right-6">
                    {item.step}
                  </span>
                  <div className="relative">
                    <div className="w-12 h-12 rounded-xl bg-emerald-500/10 flex items-center justify-center mb-4">
                      {index === 0 && <FileCode className="w-6 h-6 text-emerald-400" />}
                      {index === 1 && <Code2 className="w-6 h-6 text-emerald-400" />}
                      {index === 2 && <Zap className="w-6 h-6 text-emerald-400" />}
                    </div>
                    <h3 className="text-xl font-semibold text-white mb-2">{item.title}</h3>
                    <p className="text-gray-400">{item.desc}</p>
                  </div>
                </div>
              </div>
            ))}
          </div>

          {/* Quick start command */}
          <div className="mt-12 flex justify-center">
            <div className="inline-flex items-center gap-3 px-6 py-4 bg-gray-800 rounded-xl border border-gray-700 font-mono text-sm">
              <span className="text-gray-500">$</span>
              <span className="text-emerald-400">./scripts/setup-clients.sh</span>
              <button
                className="p-1.5 rounded-md hover:bg-gray-700 text-gray-500 hover:text-white transition-colors"
                onClick={() => navigator.clipboard.writeText('./scripts/setup-clients.sh')}
              >
                <Check className="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="relative py-24">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
              Powerful features for modern development
            </h2>
            <p className="text-gray-400 text-lg max-w-2xl mx-auto">
              Everything you need to generate, maintain, and evolve multi-language codebases
            </p>
          </div>

          <div className="grid md:grid-cols-2 gap-6">
            {features.map((feature) => (
              <div
                key={feature.title}
                className="group relative bg-gray-900/50 rounded-2xl p-8 border border-gray-800 hover:border-gray-700 transition-all"
              >
                <div className="absolute inset-0 bg-gradient-to-br from-emerald-500/5 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity" />
                <div className="relative">
                  <div className="w-14 h-14 rounded-2xl bg-gradient-to-br from-emerald-500/20 to-cyan-500/20 flex items-center justify-center mb-6 group-hover:scale-110 transition-transform">
                    <feature.icon className="w-7 h-7 text-emerald-400" />
                  </div>
                  <h3 className="text-xl font-semibold text-white mb-3">{feature.title}</h3>
                  <p className="text-gray-400 leading-relaxed">{feature.description}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Languages Grid */}
      <section className="relative py-24 bg-gray-900/50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
              Six languages, one spec
            </h2>
            <p className="text-gray-400 text-lg max-w-2xl mx-auto">
              Generate idiomatic code with proper conventions, error handling, and project structure
            </p>
          </div>

          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
            {languages.map((lang) => (
              <Link
                key={lang.name}
                to={`/languages/${lang.name.toLowerCase().replace('#', 'sharp')}`}
                className="group relative bg-gray-800/50 rounded-xl p-6 border border-gray-700/50 hover:border-gray-600 transition-all text-center"
              >
                <div
                  className="absolute inset-0 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity"
                  style={{ backgroundColor: `${lang.color}10` }}
                />
                <div className="relative">
                  <div
                    className="w-12 h-12 rounded-xl mx-auto mb-4 flex items-center justify-center font-bold text-lg"
                    style={{ backgroundColor: `${lang.color}20`, color: lang.color }}
                  >
                    {lang.name.charAt(0)}
                  </div>
                  <span className="text-white font-medium">{lang.name}</span>
                </div>
              </Link>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="relative py-24">
        <div className="absolute inset-0">
          <div className="absolute inset-0 bg-gradient-to-r from-emerald-500/10 via-cyan-500/10 to-violet-500/10" />
        </div>
        <div className="relative max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h2 className="text-4xl sm:text-5xl font-bold text-white mb-6">
            Ready to transform your workflow?
          </h2>
          <p className="text-xl text-gray-400 mb-8 max-w-2xl mx-auto">
            Start generating production-ready code from markdown specs today.
          </p>
          <div className="flex flex-wrap justify-center gap-4">
            <Link
              to="/getting-started/installation"
              className="group inline-flex items-center gap-2 px-8 py-4 bg-white text-gray-900 rounded-xl font-semibold hover:bg-gray-100 transition-all shadow-lg shadow-white/10"
            >
              Start Building
              <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
            </Link>
            <a
              href="https://github.com/kon1790/rpg"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-8 py-4 bg-gray-800 text-white rounded-xl font-semibold hover:bg-gray-700 transition-all border border-gray-700"
            >
              View on GitHub
            </a>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-gray-800 py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex flex-col md:flex-row justify-between items-center gap-4">
            <div className="flex items-center gap-2">
              <div className="h-8 w-8 rounded-lg bg-gradient-to-br from-emerald-500 to-cyan-500 flex items-center justify-center">
                <span className="text-white font-bold">R</span>
              </div>
              <span className="text-gray-400">RPG - Rosetta Project Generator</span>
            </div>
            <p className="text-gray-500 text-sm">
              Built with AI, for AI-powered development.
            </p>
          </div>
        </div>
      </footer>
    </div>
  );
}
