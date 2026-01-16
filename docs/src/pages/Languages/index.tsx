import { Link } from 'react-router-dom';

const languages = [
  { name: 'Go', href: '/languages/go', color: '#00ADD8' },
  { name: 'Rust', href: '/languages/rust', color: '#CE422B' },
  { name: 'Java', href: '/languages/java', color: '#F89820' },
  { name: 'C#', href: '/languages/csharp', color: '#953DAC' },
  { name: 'Python', href: '/languages/python', color: '#3776AB' },
  { name: 'TypeScript', href: '/languages/typescript', color: '#3178C6' },
];

export default function Languages() {
  return (
    <div className="prose-docs">
      <h1>Supported Languages</h1>
      <p className="lead">RPG generates idiomatic code for six programming languages.</p>

      <div className="grid grid-cols-2 md:grid-cols-3 gap-4 mt-8 not-prose">
        {languages.map((lang) => (
          <Link
            key={lang.name}
            to={lang.href}
            className="p-6 rounded-xl border border-gray-200 dark:border-gray-800 hover:border-gray-300 dark:hover:border-gray-700 transition-colors text-center"
          >
            <div
              className="w-12 h-12 rounded-xl mx-auto mb-3 flex items-center justify-center font-bold text-lg"
              style={{ backgroundColor: `${lang.color}20`, color: lang.color }}
            >
              {lang.name.charAt(0)}
            </div>
            <span className="font-medium text-gray-900 dark:text-white">{lang.name}</span>
          </Link>
        ))}
      </div>
    </div>
  );
}
