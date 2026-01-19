import { useState, useRef, useEffect, useCallback } from 'react';
import {
  User,
  Bot,
  Copy,
  Check,
  Folder,
  FileCode,
  GitBranch,
  CheckCircle2,
  Loader2,
  Sparkles,
  ChevronLeft,
  ChevronRight,
  ChevronDown,
  Lock,
} from 'lucide-react';

// Language colors matching the design system
const languageColors: Record<string, string> = {
  go: '#00ADD8',
  typescript: '#3178C6',
  rust: '#CE422B',
  python: '#3776AB',
  java: '#F89820',
  csharp: '#953DAC',
};

// Step configuration
const STEPS = [
  { id: 1, label: 'User Request', shortLabel: 'Request' },
  { id: 2, label: 'Analyze Repository', shortLabel: 'Analyze' },
  { id: 3, label: 'Generate Spec', shortLabel: 'Spec' },
  { id: 4, label: 'Select Languages', shortLabel: 'Languages' },
  { id: 5, label: 'Generate TypeScript', shortLabel: 'TypeScript' },
  { id: 6, label: 'Generate Go', shortLabel: 'Go' },
  { id: 7, label: 'Verify Parity', shortLabel: 'Parity' },
  { id: 8, label: 'Complete', shortLabel: 'Complete' },
];

// ConversationMessage component
interface MessageProps {
  role: 'user' | 'assistant';
  children: React.ReactNode;
}

function ConversationMessage({ role, children }: MessageProps) {
  const isUser = role === 'user';

  return (
    <div className={`flex gap-3 ${isUser ? 'flex-row-reverse' : 'flex-row'}`}>
      <div
        className={`flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center ${
          isUser
            ? 'bg-primary-500 text-white'
            : 'bg-gradient-to-br from-emerald-500 to-cyan-500 text-white'
        }`}
      >
        {isUser ? <User className="w-4 h-4" /> : <Bot className="w-4 h-4" />}
      </div>
      <div
        className={`max-w-[85%] rounded-2xl px-4 py-3 ${
          isUser
            ? 'bg-primary-500 text-white rounded-tr-sm'
            : 'bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-gray-100 rounded-tl-sm'
        }`}
      >
        {children}
      </div>
    </div>
  );
}

// CodeBlock component
interface CodeBlockProps {
  language: string;
  code: string;
  filename?: string;
  showCopy?: boolean;
}

function CodeBlock({ language, code, filename, showCopy = true }: CodeBlockProps) {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    await navigator.clipboard.writeText(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const langColor = languageColors[language.toLowerCase()] || '#6B7280';

  return (
    <div className="rounded-lg overflow-hidden bg-gray-900 dark:bg-gray-950 border border-gray-800 my-3">
      {filename && (
        <div className="flex items-center justify-between px-4 py-2 bg-gray-800 dark:bg-gray-900 border-b border-gray-700">
          <div className="flex items-center gap-2">
            <FileCode className="w-4 h-4" style={{ color: langColor }} />
            <span className="text-sm text-gray-300 font-mono">{filename}</span>
          </div>
          {showCopy && (
            <button
              onClick={handleCopy}
              className="p-1.5 rounded hover:bg-gray-700 text-gray-400 hover:text-white transition-colors"
            >
              {copied ? <Check className="w-4 h-4 text-emerald-400" /> : <Copy className="w-4 h-4" />}
            </button>
          )}
        </div>
      )}
      <pre className="p-4 overflow-x-auto text-sm font-mono">
        <code className="text-gray-300">{code}</code>
      </pre>
    </div>
  );
}

// FileTree component
interface FileItem {
  name: string;
  type: 'file' | 'folder';
  language?: string;
  children?: FileItem[];
}

interface FileTreeProps {
  files: FileItem[];
  title?: string;
}

function FileTree({ files, title }: FileTreeProps) {
  const renderItem = (item: FileItem, depth: number = 0) => {
    const langColor = item.language ? languageColors[item.language.toLowerCase()] : undefined;

    return (
      <div key={item.name} style={{ marginLeft: depth * 16 }}>
        <div className="flex items-center gap-2 py-1">
          {item.type === 'folder' ? (
            <Folder className="w-4 h-4 text-yellow-500" />
          ) : (
            <FileCode className="w-4 h-4" style={{ color: langColor || '#9CA3AF' }} />
          )}
          <span className="text-sm font-mono" style={{ color: langColor }}>
            {item.name}
          </span>
        </div>
        {item.children?.map((child) => renderItem(child, depth + 1))}
      </div>
    );
  };

  return (
    <div className="rounded-lg bg-gray-50 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700 p-4 my-3">
      {title && (
        <div className="flex items-center gap-2 mb-3 pb-2 border-b border-gray-200 dark:border-gray-700">
          <GitBranch className="w-4 h-4 text-gray-500" />
          <span className="text-sm font-medium text-gray-700 dark:text-gray-300">{title}</span>
        </div>
      )}
      <div className="space-y-0.5">{files.map((file) => renderItem(file))}</div>
    </div>
  );
}

// ParityTable component
interface ParityFeature {
  name: string;
  typescript: boolean;
  go: boolean;
}

interface ParityTableProps {
  features: ParityFeature[];
  score: number;
}

function ParityTable({ features, score }: ParityTableProps) {
  return (
    <div className="rounded-lg bg-gray-50 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700 overflow-hidden my-3">
      <div className="px-4 py-3 bg-gray-100 dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between">
        <span className="text-sm font-medium text-gray-700 dark:text-gray-300">Feature Parity</span>
        <span className="px-2 py-1 rounded-full bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400 text-xs font-medium">
          {score}% Match
        </span>
      </div>
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-gray-200 dark:border-gray-700">
            <th className="text-left px-4 py-2 font-medium text-gray-600 dark:text-gray-400">Feature</th>
            <th className="text-center px-4 py-2 font-medium" style={{ color: languageColors.typescript }}>
              TypeScript
            </th>
            <th className="text-center px-4 py-2 font-medium" style={{ color: languageColors.go }}>
              Go
            </th>
          </tr>
        </thead>
        <tbody>
          {features.map((feature) => (
            <tr key={feature.name} className="border-b border-gray-100 dark:border-gray-700/50 last:border-0">
              <td className="px-4 py-2 text-gray-700 dark:text-gray-300">{feature.name}</td>
              <td className="text-center px-4 py-2">
                {feature.typescript ? (
                  <CheckCircle2 className="w-5 h-5 text-emerald-500 mx-auto" />
                ) : (
                  <span className="text-gray-300">-</span>
                )}
              </td>
              <td className="text-center px-4 py-2">
                {feature.go ? (
                  <CheckCircle2 className="w-5 h-5 text-emerald-500 mx-auto" />
                ) : (
                  <span className="text-gray-300">-</span>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

// Status indicator component
function StatusIndicator({ text, done = false }: { text: string; done?: boolean }) {
  return (
    <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400 my-2">
      {done ? (
        <CheckCircle2 className="w-4 h-4 text-emerald-500" />
      ) : (
        <Loader2 className="w-4 h-4 animate-spin text-primary-500" />
      )}
      <span>{text}</span>
    </div>
  );
}

// StepTabs component for desktop
interface StepTabsProps {
  steps: typeof STEPS;
  currentStep: number;
  maxUnlockedStep: number;
  onStepClick: (step: number) => void;
}

function StepTabs({ steps, currentStep, maxUnlockedStep, onStepClick }: StepTabsProps) {
  const tabsRef = useRef<HTMLDivElement>(null);

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      let newStep = currentStep;

      switch (e.key) {
        case 'ArrowLeft':
          e.preventDefault();
          newStep = Math.max(1, currentStep - 1);
          break;
        case 'ArrowRight':
          e.preventDefault();
          newStep = Math.min(maxUnlockedStep, currentStep + 1);
          break;
        case 'Home':
          e.preventDefault();
          newStep = 1;
          break;
        case 'End':
          e.preventDefault();
          newStep = maxUnlockedStep;
          break;
        default:
          return;
      }

      if (newStep !== currentStep && newStep <= maxUnlockedStep) {
        onStepClick(newStep);
      }
    },
    [currentStep, maxUnlockedStep, onStepClick]
  );

  return (
    <div
      ref={tabsRef}
      role="tablist"
      aria-label="Quick start steps"
      className="hidden md:flex items-center justify-center gap-1 mb-6 flex-wrap"
      onKeyDown={handleKeyDown}
    >
      {steps.map((step, index) => {
        const isActive = step.id === currentStep;
        const isCompleted = step.id < maxUnlockedStep;
        const isLocked = step.id > maxUnlockedStep;
        const isAccessible = step.id <= maxUnlockedStep;

        return (
          <div key={step.id} className="flex items-center">
            <button
              role="tab"
              aria-selected={isActive}
              aria-disabled={isLocked}
              aria-controls={`step-panel-${step.id}`}
              id={`step-tab-${step.id}`}
              tabIndex={isActive ? 0 : -1}
              onClick={() => isAccessible && onStepClick(step.id)}
              disabled={isLocked}
              className={`
                flex items-center gap-1.5 px-3 py-2 rounded-lg text-sm font-medium transition-all
                ${
                  isActive
                    ? 'bg-primary-500 text-white shadow-md scale-105'
                    : isCompleted
                    ? 'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-400 hover:bg-emerald-200 dark:hover:bg-emerald-900/50 cursor-pointer'
                    : 'bg-gray-100 dark:bg-gray-800 text-gray-400 dark:text-gray-500 cursor-not-allowed'
                }
              `}
            >
              {isLocked ? (
                <Lock className="w-3.5 h-3.5" />
              ) : isCompleted ? (
                <Check className="w-3.5 h-3.5" />
              ) : (
                <span className="w-4 h-4 flex items-center justify-center text-xs rounded-full bg-white/20">
                  {step.id}
                </span>
              )}
              <span className="hidden lg:inline">{step.shortLabel}</span>
              <span className="lg:hidden">{step.id}</span>
            </button>
            {index < steps.length - 1 && (
              <ChevronRight
                className={`w-4 h-4 mx-1 ${
                  step.id < maxUnlockedStep ? 'text-emerald-500' : 'text-gray-300 dark:text-gray-600'
                }`}
              />
            )}
          </div>
        );
      })}
    </div>
  );
}

// StepAccordion component for mobile
interface StepAccordionProps {
  steps: typeof STEPS;
  currentStep: number;
  maxUnlockedStep: number;
  onStepClick: (step: number) => void;
  renderContent: (step: number) => React.ReactNode;
}

function StepAccordion({ steps, currentStep, maxUnlockedStep, onStepClick, renderContent }: StepAccordionProps) {
  return (
    <div className="md:hidden space-y-2">
      {steps.map((step) => {
        const isActive = step.id === currentStep;
        const isCompleted = step.id < maxUnlockedStep;
        const isLocked = step.id > maxUnlockedStep;
        const isAccessible = step.id <= maxUnlockedStep;

        return (
          <div
            key={step.id}
            className={`rounded-lg border overflow-hidden transition-all ${
              isActive
                ? 'border-primary-500 dark:border-primary-400'
                : isCompleted
                ? 'border-emerald-300 dark:border-emerald-700'
                : 'border-gray-200 dark:border-gray-700'
            }`}
          >
            <button
              onClick={() => isAccessible && onStepClick(step.id)}
              disabled={isLocked}
              className={`w-full flex items-center justify-between px-4 py-3 text-left transition-colors ${
                isActive
                  ? 'bg-primary-50 dark:bg-primary-900/20'
                  : isCompleted
                  ? 'bg-emerald-50 dark:bg-emerald-900/10 hover:bg-emerald-100 dark:hover:bg-emerald-900/20'
                  : 'bg-gray-50 dark:bg-gray-800/50'
              }`}
            >
              <div className="flex items-center gap-3">
                <div
                  className={`w-6 h-6 rounded-full flex items-center justify-center text-xs font-medium ${
                    isActive
                      ? 'bg-primary-500 text-white'
                      : isCompleted
                      ? 'bg-emerald-500 text-white'
                      : 'bg-gray-200 dark:bg-gray-700 text-gray-500 dark:text-gray-400'
                  }`}
                >
                  {isCompleted ? <Check className="w-3.5 h-3.5" /> : isLocked ? <Lock className="w-3 h-3" /> : step.id}
                </div>
                <span
                  className={`text-sm font-medium ${
                    isActive
                      ? 'text-primary-700 dark:text-primary-300'
                      : isCompleted
                      ? 'text-emerald-700 dark:text-emerald-400'
                      : 'text-gray-500 dark:text-gray-400'
                  }`}
                >
                  {step.label}
                </span>
              </div>
              <ChevronDown
                className={`w-4 h-4 transition-transform ${
                  isActive ? 'rotate-180 text-primary-500' : 'text-gray-400'
                } ${isLocked ? 'opacity-50' : ''}`}
              />
            </button>
            {isActive && (
              <div className="px-4 py-4 bg-white dark:bg-gray-900 border-t border-gray-100 dark:border-gray-800">
                {renderContent(step.id)}
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
}

// NavigationButtons component
interface NavigationButtonsProps {
  currentStep: number;
  maxStep: number;
  maxUnlockedStep: number;
  onPrevious: () => void;
  onNext: () => void;
}

function NavigationButtons({ currentStep, maxStep, maxUnlockedStep, onPrevious, onNext }: NavigationButtonsProps) {
  const currentStepInfo = STEPS.find((s) => s.id === currentStep);

  return (
    <div className="flex items-center justify-between py-3">
      <button
        onClick={onPrevious}
        disabled={currentStep === 1}
        className="flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-medium text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        <ChevronLeft className="w-4 h-4" />
        Previous
      </button>
      <div className="text-sm text-gray-500 dark:text-gray-400 hidden sm:block">
        Step {currentStep} of {maxStep}: {currentStepInfo?.label}
      </div>
      <button
        onClick={onNext}
        disabled={currentStep >= maxUnlockedStep}
        className="flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-medium bg-primary-500 text-white hover:bg-primary-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        Next
        <ChevronRight className="w-4 h-4" />
      </button>
    </div>
  );
}

// ScrollableContent component with fade gradient
interface ScrollableContentProps {
  children: React.ReactNode;
}

function ScrollableContent({ children }: ScrollableContentProps) {
  const contentRef = useRef<HTMLDivElement>(null);
  const [showBottomFade, setShowBottomFade] = useState(false);

  useEffect(() => {
    const checkScroll = () => {
      if (contentRef.current) {
        const { scrollHeight, clientHeight, scrollTop } = contentRef.current;
        const isScrollable = scrollHeight > clientHeight;
        const isAtBottom = scrollHeight - scrollTop - clientHeight < 20;
        setShowBottomFade(isScrollable && !isAtBottom);
      }
    };

    checkScroll();
    const element = contentRef.current;
    element?.addEventListener('scroll', checkScroll);
    window.addEventListener('resize', checkScroll);

    return () => {
      element?.removeEventListener('scroll', checkScroll);
      window.removeEventListener('resize', checkScroll);
    };
  }, [children]);

  // Scroll to top when children change
  useEffect(() => {
    contentRef.current?.scrollTo({ top: 0, behavior: 'smooth' });
  }, [children]);

  return (
    <div className="relative">
      <div
        ref={contentRef}
        className="max-h-[calc(100vh-280px)] min-h-[400px] overflow-y-auto pr-2"
      >
        {children}
      </div>
      {showBottomFade && (
        <div className="absolute bottom-0 left-0 right-2 h-16 bg-gradient-to-t from-white dark:from-gray-900 to-transparent pointer-events-none" />
      )}
    </div>
  );
}

// Data definitions
const tsFiles: FileItem[] = [
  {
    name: 'task-dashboard-ts/',
    type: 'folder',
    children: [
      {
        name: 'src/',
        type: 'folder',
        children: [
          { name: 'index.ts', type: 'file', language: 'typescript' },
          { name: 'types.ts', type: 'file', language: 'typescript' },
          { name: 'task.ts', type: 'file', language: 'typescript' },
          { name: 'dashboard.ts', type: 'file', language: 'typescript' },
          { name: 'utils.ts', type: 'file', language: 'typescript' },
        ],
      },
      {
        name: 'tests/',
        type: 'folder',
        children: [
          { name: 'task.test.ts', type: 'file', language: 'typescript' },
          { name: 'dashboard.test.ts', type: 'file', language: 'typescript' },
        ],
      },
      { name: 'package.json', type: 'file' },
      { name: 'tsconfig.json', type: 'file' },
      { name: 'README.md', type: 'file' },
    ],
  },
];

const goFiles: FileItem[] = [
  {
    name: 'task-dashboard-go/',
    type: 'folder',
    children: [
      {
        name: 'cmd/',
        type: 'folder',
        children: [{ name: 'main.go', type: 'file', language: 'go' }],
      },
      {
        name: 'internal/',
        type: 'folder',
        children: [
          { name: 'task.go', type: 'file', language: 'go' },
          { name: 'dashboard.go', type: 'file', language: 'go' },
          { name: 'types.go', type: 'file', language: 'go' },
        ],
      },
      { name: 'task_test.go', type: 'file', language: 'go' },
      { name: 'dashboard_test.go', type: 'file', language: 'go' },
      { name: 'go.mod', type: 'file' },
      { name: 'README.md', type: 'file' },
    ],
  },
];

const parityFeatures: ParityFeature[] = [
  { name: 'Task CRUD operations', typescript: true, go: true },
  { name: 'Priority levels', typescript: true, go: true },
  { name: 'Due date handling', typescript: true, go: true },
  { name: 'Dashboard statistics', typescript: true, go: true },
  { name: 'Task filtering', typescript: true, go: true },
  { name: 'Data validation', typescript: true, go: true },
];

const specPreview = `# Task Dashboard

A task management application for tracking todos with priorities.

## Types

### Task
- id: UUID
- title: Text
- description: Text (optional)
- priority: Priority (low, medium, high)
- dueDate: DateTime (optional)
- completed: Boolean
- createdAt: DateTime

### Dashboard
- tasks: List<Task>
- completedCount: Integer
- pendingCount: Integer

## Functions

### createTask
Creates a new task with the given details.
**accepts:** title, priority, dueDate (optional)
**returns:** Task

### completeTask
Marks a task as completed.
**accepts:** taskId
**returns:** Task | Error

### getDashboardStats
Returns dashboard statistics.
**returns:** Dashboard`;

const tsCodePreview = `export interface Task {
  id: string;
  title: string;
  description?: string;
  priority: 'low' | 'medium' | 'high';
  dueDate?: Date;
  completed: boolean;
  createdAt: Date;
}

export function createTask(
  title: string,
  priority: Task['priority'],
  dueDate?: Date
): Task {
  return {
    id: crypto.randomUUID(),
    title,
    priority,
    dueDate,
    completed: false,
    createdAt: new Date(),
  };
}`;

const goCodePreview = `package taskdashboard

import (
    "time"
    "github.com/google/uuid"
)

type Priority string

const (
    PriorityLow    Priority = "low"
    PriorityMedium Priority = "medium"
    PriorityHigh   Priority = "high"
)

type Task struct {
    ID          string     \`json:"id"\`
    Title       string     \`json:"title"\`
    Description string     \`json:"description,omitempty"\`
    Priority    Priority   \`json:"priority"\`
    DueDate     *time.Time \`json:"dueDate,omitempty"\`
    Completed   bool       \`json:"completed"\`
    CreatedAt   time.Time  \`json:"createdAt"\`
}

func CreateTask(title string, priority Priority, dueDate *time.Time) *Task {
    return &Task{
        ID:        uuid.New().String(),
        Title:     title,
        Priority:  priority,
        DueDate:   dueDate,
        Completed: false,
        CreatedAt: time.Now(),
    }
}`;

const supportedLanguages = [
  { id: 'go', name: 'Go', color: languageColors.go },
  { id: 'rust', name: 'Rust', color: languageColors.rust },
  { id: 'java', name: 'Java', color: languageColors.java },
  { id: 'csharp', name: 'C#', color: languageColors.csharp },
  { id: 'python', name: 'Python', color: languageColors.python },
  { id: 'typescript', name: 'TypeScript', color: languageColors.typescript },
];

// Step content components - isolated content for each step
function StepContent({ step }: { step: number }) {
  switch (step) {
    case 1:
      return (
        <div className="space-y-4">
          <ConversationMessage role="user">
            <p>
              I have a task dashboard app on GitHub at{' '}
              <code className="px-1.5 py-0.5 bg-white/20 rounded text-sm">github-samples/task-dashboard</code>.
              Can you help me port it to Go and TypeScript?
            </p>
          </ConversationMessage>
        </div>
      );

    case 2:
      return (
        <div className="space-y-4">
          <ConversationMessage role="assistant">
            <p className="mb-3">
              I'll analyze that repository and generate implementations in both languages. Let me start by importing and analyzing the source code.
            </p>
            <div className="bg-white/10 dark:bg-black/20 rounded-lg p-3 space-y-1">
              <StatusIndicator text="Cloning repository..." done />
              <StatusIndicator text="Analyzing source files..." done />
              <StatusIndicator text="Detecting patterns and types..." done />
            </div>
            <div className="mt-3 p-3 bg-white/10 dark:bg-black/20 rounded-lg">
              <p className="text-sm font-medium mb-2">Analysis Complete:</p>
              <ul className="text-sm space-y-1 text-gray-700 dark:text-gray-300">
                <li>12 source files detected</li>
                <li>3 main types identified (Task, Dashboard, Priority)</li>
                <li>8 functions mapped</li>
                <li>Test coverage patterns found</li>
              </ul>
            </div>
          </ConversationMessage>
        </div>
      );

    case 3:
      return (
        <div className="space-y-4">
          <ConversationMessage role="assistant">
            <p className="mb-3">
              I've generated a specification from the source code. Here's a preview:
            </p>
            <CodeBlock language="markdown" code={specPreview} filename="task-dashboard.spec.md" />
            <p className="text-sm text-gray-600 dark:text-gray-400 mt-2">
              This spec captures all the types, functions, and behaviors from your original code.
            </p>
          </ConversationMessage>
        </div>
      );

    case 4:
      return (
        <div className="space-y-4">
          <ConversationMessage role="assistant">
            <p className="mb-3">RPG supports 6 target languages. You requested TypeScript and Go:</p>
            <div className="grid grid-cols-3 gap-2 my-3">
              {supportedLanguages.map((lang) => {
                const isSelected = lang.id === 'typescript' || lang.id === 'go';
                return (
                  <div
                    key={lang.id}
                    className={`px-3 py-2 rounded-lg text-sm font-medium text-center transition-all ${
                      isSelected ? '' : 'opacity-50'
                    }`}
                    style={{
                      backgroundColor: `${lang.color}15`,
                      color: lang.color,
                      boxShadow: isSelected ? `0 0 0 2px ${lang.color}` : undefined,
                    }}
                  >
                    {lang.name}
                    {isSelected && <Check className="w-3 h-3 inline ml-1" />}
                  </div>
                );
              })}
            </div>
          </ConversationMessage>
        </div>
      );

    case 5:
      return (
        <div className="space-y-4">
          <ConversationMessage role="assistant">
            <p className="mb-3">TypeScript implementation complete!</p>
            <FileTree files={tsFiles} title="Generated TypeScript Project" />
            <CodeBlock language="typescript" code={tsCodePreview} filename="src/task.ts" />
          </ConversationMessage>
        </div>
      );

    case 6:
      return (
        <div className="space-y-4">
          <ConversationMessage role="assistant">
            <p className="mb-3">Go implementation complete with idiomatic patterns!</p>
            <FileTree files={goFiles} title="Generated Go Project" />
            <CodeBlock language="go" code={goCodePreview} filename="internal/task.go" />
          </ConversationMessage>
        </div>
      );

    case 7:
      return (
        <div className="space-y-4">
          <ConversationMessage role="assistant">
            <p className="mb-3">Feature parity verification complete!</p>
            <ParityTable features={parityFeatures} score={100} />
            <p className="text-sm text-gray-600 dark:text-gray-400 mt-2">
              Both implementations achieve 100% feature parity, ensuring consistent behavior across languages.
            </p>
          </ConversationMessage>
        </div>
      );

    case 8:
      return (
        <div className="space-y-4">
          <ConversationMessage role="assistant">
            <div className="space-y-4">
              <div className="flex items-center gap-2 text-emerald-600 dark:text-emerald-400">
                <CheckCircle2 className="w-5 h-5" />
                <span className="font-medium">Generation Complete!</span>
              </div>
              <p>
                I've successfully ported your task dashboard to both TypeScript and Go. Both implementations:
              </p>
              <ul className="list-disc list-inside space-y-1 text-sm">
                <li>Follow language-specific idioms and conventions</li>
                <li>Include comprehensive type definitions</li>
                <li>Have matching test suites</li>
                <li>Achieve 100% feature parity</li>
              </ul>
              <div className="mt-4 p-3 bg-white/10 dark:bg-black/20 rounded-lg">
                <p className="text-sm font-medium mb-2">Generated projects:</p>
                <ul className="text-sm space-y-1">
                  <li>
                    <code className="px-1.5 py-0.5 bg-[#3178C6]/20 rounded" style={{ color: languageColors.typescript }}>
                      generated/typescript/task-dashboard-ts/
                    </code>
                  </li>
                  <li>
                    <code className="px-1.5 py-0.5 bg-[#00ADD8]/20 rounded" style={{ color: languageColors.go }}>
                      generated/go/task-dashboard-go/
                    </code>
                  </li>
                </ul>
              </div>
            </div>
          </ConversationMessage>

          {/* Next steps section */}
          <div className="mt-6 p-6 rounded-xl bg-gradient-to-br from-primary-50 to-emerald-50 dark:from-primary-900/20 dark:to-emerald-900/20 border border-primary-200 dark:border-primary-800">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">Next Steps</h3>
            <ul className="space-y-2 text-sm text-gray-700 dark:text-gray-300">
              <li>
                <a href="/getting-started/installation" className="text-primary-600 dark:text-primary-400 hover:underline">Install RPG</a> on your machine
              </li>
              <li>
                Learn to <a href="/writing-specs" className="text-primary-600 dark:text-primary-400 hover:underline">write your own specs</a> from scratch
              </li>
              <li>
                Explore <a href="/mcp-tools" className="text-primary-600 dark:text-primary-400 hover:underline">all MCP tools</a> available
              </li>
              <li>
                Try the <a href="/tutorials/multi-language" className="text-primary-600 dark:text-primary-400 hover:underline">multi-language tutorial</a> with your own project
              </li>
            </ul>
          </div>
        </div>
      );

    default:
      return null;
  }
}

export default function QuickStart() {
  const [currentStep, setCurrentStep] = useState(1);
  // For this demo, steps unlock progressively as user navigates
  const [maxUnlockedStep, setMaxUnlockedStep] = useState(8); // All unlocked for demo

  const goToStep = (step: number) => {
    if (step >= 1 && step <= maxUnlockedStep) {
      setCurrentStep(step);
    }
  };

  const handlePrevious = () => {
    if (currentStep > 1) {
      setCurrentStep(currentStep - 1);
    }
  };

  const handleNext = () => {
    if (currentStep < maxUnlockedStep) {
      const nextStep = currentStep + 1;
      setCurrentStep(nextStep);
      // Unlock next step if needed
      if (nextStep > maxUnlockedStep) {
        setMaxUnlockedStep(nextStep);
      }
    }
  };

  const renderContent = useCallback(
    (step: number) => <StepContent step={step} />,
    []
  );

  return (
    <div className="max-w-4xl mx-auto">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">Quick Start</h1>
        <p className="text-lg text-gray-600 dark:text-gray-400">
          See RPG in action: import a GitHub project and generate multi-language implementations.
        </p>
      </div>

      {/* Info banner */}
      <div className="mb-6 p-4 rounded-lg bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-200 dark:border-emerald-800">
        <div className="flex items-start gap-3">
          <Sparkles className="w-5 h-5 text-emerald-600 dark:text-emerald-400 mt-0.5" />
          <div>
            <p className="text-emerald-800 dark:text-emerald-300 text-sm">
              <strong>Interactive Demo:</strong> This shows a typical workflow using{' '}
              <code className="px-1.5 py-0.5 bg-emerald-100 dark:bg-emerald-800/50 rounded text-xs">
                github-samples/task-dashboard
              </code>{' '}
              as an example. Click tabs or use arrow keys to navigate.
            </p>
          </div>
        </div>
      </div>

      {/* Desktop: Tabs */}
      <StepTabs
        steps={STEPS}
        currentStep={currentStep}
        maxUnlockedStep={maxUnlockedStep}
        onStepClick={goToStep}
      />

      {/* Desktop: Navigation buttons (top) + Content */}
      <div className="hidden md:block">
        <NavigationButtons
          currentStep={currentStep}
          maxStep={STEPS.length}
          maxUnlockedStep={maxUnlockedStep}
          onPrevious={handlePrevious}
          onNext={handleNext}
        />

        <div
          role="tabpanel"
          id={`step-panel-${currentStep}`}
          aria-labelledby={`step-tab-${currentStep}`}
          className="border border-gray-200 dark:border-gray-700 rounded-xl p-6 bg-white dark:bg-gray-900"
        >
          <ScrollableContent>
            <StepContent step={currentStep} />
          </ScrollableContent>
        </div>

        <NavigationButtons
          currentStep={currentStep}
          maxStep={STEPS.length}
          maxUnlockedStep={maxUnlockedStep}
          onPrevious={handlePrevious}
          onNext={handleNext}
        />
      </div>

      {/* Mobile: Accordion */}
      <StepAccordion
        steps={STEPS}
        currentStep={currentStep}
        maxUnlockedStep={maxUnlockedStep}
        onStepClick={goToStep}
        renderContent={renderContent}
      />
    </div>
  );
}
