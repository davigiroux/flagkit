import { cn } from '../lib/utils'

const colors = {
  production: 'bg-red-100 text-red-700',
  staging: 'bg-yellow-100 text-yellow-700',
  development: 'bg-blue-100 text-blue-700',
}

export function EnvBadge({ env }: { env: string }) {
  return (
    <span className={cn('px-2 py-0.5 rounded-full text-xs font-medium', colors[env as keyof typeof colors] || 'bg-gray-100 text-gray-700')}>
      {env}
    </span>
  )
}
