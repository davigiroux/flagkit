import { useSortable } from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import type { Rule } from '../api/client'

interface RuleCardProps {
  id: string
  rule: Rule
  onChange: (rule: Rule) => void
  onDelete: () => void
}

export function RuleCard({ id, rule, onChange, onDelete }: RuleCardProps) {
  const { attributes, listeners, setNodeRef, transform, transition } = useSortable({ id })

  const style = { transform: CSS.Transform.toString(transform), transition }

  return (
    <div ref={setNodeRef} style={style} className="bg-white border border-gray-200 rounded-lg p-4 flex gap-4 items-start">
      <button {...attributes} {...listeners} className="mt-1 cursor-grab text-gray-400 hover:text-gray-600" title="Drag to reorder">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <circle cx="5" cy="3" r="1.5" />
          <circle cx="11" cy="3" r="1.5" />
          <circle cx="5" cy="8" r="1.5" />
          <circle cx="11" cy="8" r="1.5" />
          <circle cx="5" cy="13" r="1.5" />
          <circle cx="11" cy="13" r="1.5" />
        </svg>
      </button>

      <div className="flex-1 space-y-3">
        <select
          value={rule.type}
          onChange={(e) => onChange({ ...rule, type: e.target.value as Rule['type'] })}
          className="px-3 py-1.5 border border-gray-300 rounded-md text-sm"
        >
          <option value="percentage">Percentage Rollout</option>
          <option value="allowlist">User Allowlist</option>
        </select>

        {rule.type === 'percentage' ? (
          <div className="flex items-center gap-3">
            <input
              type="range"
              min={0}
              max={100}
              value={rule.value ?? 0}
              onChange={(e) => onChange({ ...rule, value: Number(e.target.value) })}
              className="flex-1"
            />
            <span className="text-sm font-mono w-12 text-right">{rule.value ?? 0}%</span>
          </div>
        ) : (
          <div>
            <textarea
              value={(rule.userIds ?? []).join('\n')}
              onChange={(e) => onChange({ ...rule, userIds: e.target.value.split('\n').filter(Boolean) })}
              placeholder="One user ID per line"
              rows={3}
              className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm font-mono"
            />
            <span className="text-xs text-gray-500">{(rule.userIds ?? []).length} user(s)</span>
          </div>
        )}
      </div>

      <button onClick={onDelete} className="text-gray-400 hover:text-red-500 mt-1" title="Remove rule">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path d="M4.646 4.646a.5.5 0 0 1 .708 0L8 7.293l2.646-2.647a.5.5 0 0 1 .708.708L8.707 8l2.647 2.646a.5.5 0 0 1-.708.708L8 8.707l-2.646 2.647a.5.5 0 0 1-.708-.708L7.293 8 4.646 5.354a.5.5 0 0 1 0-.708z" />
        </svg>
      </button>
    </div>
  )
}
