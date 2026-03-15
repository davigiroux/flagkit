import { DndContext, closestCenter, type DragEndEvent } from '@dnd-kit/core'
import { SortableContext, verticalListSortingStrategy, arrayMove } from '@dnd-kit/sortable'
import type { Rule } from '../api/client'
import { RuleCard } from './RuleCard'

interface RuleBuilderProps {
  rules: Rule[]
  onChange: (rules: Rule[]) => void
}

export function RuleBuilder({ rules, onChange }: RuleBuilderProps) {
  const ids = rules.map((_, i) => `rule-${i}`)

  function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event
    if (!over || active.id === over.id) return
    const oldIndex = ids.indexOf(String(active.id))
    const newIndex = ids.indexOf(String(over.id))
    onChange(arrayMove(rules, oldIndex, newIndex))
  }

  function handleRuleChange(index: number, rule: Rule) {
    const next = [...rules]
    next[index] = rule
    onChange(next)
  }

  function handleDelete(index: number) {
    onChange(rules.filter((_, i) => i !== index))
  }

  function handleAdd() {
    onChange([...rules, { type: 'percentage', value: 50 }])
  }

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <label className="text-sm font-medium text-gray-700">Targeting Rules</label>
        <button onClick={handleAdd} className="text-sm text-blue-600 hover:text-blue-800 font-medium">
          + Add Rule
        </button>
      </div>

      {rules.length === 0 ? (
        <p className="text-sm text-gray-500 py-4 text-center border border-dashed border-gray-300 rounded-lg">
          No rules yet. Add a rule to start targeting users.
        </p>
      ) : (
        <DndContext collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
          <SortableContext items={ids} strategy={verticalListSortingStrategy}>
            <div className="space-y-2">
              {rules.map((rule, i) => (
                <RuleCard
                  key={ids[i]}
                  id={ids[i]}
                  rule={rule}
                  onChange={(r) => handleRuleChange(i, r)}
                  onDelete={() => handleDelete(i)}
                />
              ))}
            </div>
          </SortableContext>
        </DndContext>
      )}

      {rules.length > 1 && (
        <p className="text-xs text-gray-500">Rules evaluate top-to-bottom. First match wins. Drag to reorder.</p>
      )}
    </div>
  )
}
