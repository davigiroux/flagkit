import { useParams, useNavigate } from 'react-router'
import { useState, useEffect } from 'react'
import { useFlag, useUpdateFlag, useToggleFlag } from '../hooks/useFlags'
import { useAuditLogs } from '../hooks/useAudit'
import { RuleBuilder } from '../components/RuleBuilder'
import { Toggle } from '../components/Toggle'
import { EnvBadge } from '../components/EnvBadge'
import { formatDate } from '../lib/utils'
import type { Rule } from '../api/client'

export function FlagDetailPage() {
  const { key } = useParams<{ key: string }>()
  const navigate = useNavigate()
  const { data: flag, isLoading } = useFlag(key!)
  const updateFlag = useUpdateFlag()
  const toggleFlag = useToggleFlag()
  const { data: auditData } = useAuditLogs({ flagId: flag?.id, perPage: 5 })

  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [environment, setEnvironment] = useState('development')
  const [rules, setRules] = useState<Rule[]>([])
  const [dirty, setDirty] = useState(false)

  useEffect(() => {
    if (flag) {
      setName(flag.name)
      setDescription(flag.description)
      setEnvironment(flag.environment)
      setRules(flag.rules ?? [])
      setDirty(false)
    }
  }, [flag])

  function handleSave() {
    if (!key) return
    updateFlag.mutate(
      { key, data: { name, description, environment: environment as 'production' | 'staging' | 'development', rules } },
      { onSuccess: () => setDirty(false) }
    )
  }

  if (isLoading) return <p className="text-gray-500">Loading...</p>
  if (!flag) return <p className="text-gray-500">Flag not found.</p>

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <button onClick={() => navigate('/')} className="text-gray-500 hover:text-gray-700">
            &larr; Back
          </button>
          <h1 className="text-2xl font-bold font-mono">{flag.key}</h1>
          <EnvBadge env={flag.environment} />
        </div>
        <div className="flex items-center gap-4">
          <span className="text-sm text-gray-500">{flag.enabled ? 'Enabled' : 'Disabled'}</span>
          <Toggle enabled={flag.enabled} onChange={() => toggleFlag.mutate(flag.key)} />
        </div>
      </div>

      <div className="bg-white border border-gray-200 rounded-lg p-6 space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Name</label>
            <input
              value={name}
              onChange={(e) => { setName(e.target.value); setDirty(true) }}
              className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Environment</label>
            <select
              value={environment}
              onChange={(e) => { setEnvironment(e.target.value); setDirty(true) }}
              className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm"
            >
              <option value="development">Development</option>
              <option value="staging">Staging</option>
              <option value="production">Production</option>
            </select>
          </div>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
          <input
            value={description}
            onChange={(e) => { setDescription(e.target.value); setDirty(true) }}
            className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm"
          />
        </div>

        <RuleBuilder rules={rules} onChange={(r) => { setRules(r); setDirty(true) }} />

        <div className="flex justify-end pt-2">
          <button
            onClick={handleSave}
            disabled={!dirty || updateFlag.isPending}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm font-medium hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {updateFlag.isPending ? 'Saving...' : 'Save Changes'}
          </button>
        </div>
      </div>

      {auditData && auditData.data.length > 0 && (
        <div className="bg-white border border-gray-200 rounded-lg p-6">
          <h2 className="text-lg font-semibold mb-4">Recent Activity</h2>
          <div className="space-y-3">
            {auditData.data.map((log) => (
              <div key={log.id} className="flex items-start gap-3 text-sm">
                <span className="text-xs bg-gray-100 px-2 py-0.5 rounded font-medium">{log.action}</span>
                <span className="text-gray-500">{formatDate(log.createdAt)}</span>
                {log.diff && Object.keys(log.diff).length > 0 && (
                  <span className="text-gray-400 font-mono text-xs">
                    {Object.keys(log.diff).join(', ')} changed
                  </span>
                )}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
