import { useState } from 'react'
import { Link } from 'react-router'
import { useFlags, useToggleFlag, useDeleteFlag, useCreateFlag } from '../hooks/useFlags'
import type { Flag } from '../api/client'
import { Toggle } from '../components/Toggle'
import { EnvBadge } from '../components/EnvBadge'
import { SearchFilter } from '../components/SearchFilter'
import { ConfirmDialog } from '../components/ConfirmDialog'
import { formatDate } from '../lib/utils'

export function FlagsPage() {
  const { data: flags, isLoading } = useFlags()
  const toggleFlag = useToggleFlag()
  const deleteFlag = useDeleteFlag()
  const createFlag = useCreateFlag()

  const [search, setSearch] = useState('')
  const [envFilter, setEnvFilter] = useState('')
  const [deleteKey, setDeleteKey] = useState<string | null>(null)
  const [showCreate, setShowCreate] = useState(false)
  const [newFlag, setNewFlag] = useState({ key: '', name: '', description: '', environment: 'development' })

  const filtered = (flags ?? []).filter((f) => {
    const matchSearch = f.key.includes(search) || f.name.toLowerCase().includes(search.toLowerCase())
    const matchEnv = !envFilter || f.environment === envFilter
    return matchSearch && matchEnv
  })

  function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    createFlag.mutate({ ...newFlag, environment: newFlag.environment as Flag['environment'] }, {
      onSuccess: () => {
        setShowCreate(false)
        setNewFlag({ key: '', name: '', description: '', environment: 'development' })
      },
    })
  }

  if (isLoading) return <p className="text-gray-500">Loading...</p>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Flags</h1>
        <button onClick={() => setShowCreate(true)} className="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm font-medium hover:bg-blue-700">
          Create Flag
        </button>
      </div>

      <SearchFilter search={search} onSearchChange={setSearch} environment={envFilter} onEnvironmentChange={setEnvFilter} />

      {showCreate && (
        <form onSubmit={handleCreate} className="bg-white border border-gray-200 rounded-lg p-4 space-y-3">
          <div className="grid grid-cols-2 gap-3">
            <input
              placeholder="Flag key (e.g. new-checkout)"
              value={newFlag.key}
              onChange={(e) => setNewFlag({ ...newFlag, key: e.target.value })}
              className="px-3 py-2 border border-gray-300 rounded-md text-sm"
              required
            />
            <input
              placeholder="Display name"
              value={newFlag.name}
              onChange={(e) => setNewFlag({ ...newFlag, name: e.target.value })}
              className="px-3 py-2 border border-gray-300 rounded-md text-sm"
              required
            />
          </div>
          <input
            placeholder="Description"
            value={newFlag.description}
            onChange={(e) => setNewFlag({ ...newFlag, description: e.target.value })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm"
          />
          <select
            value={newFlag.environment}
            onChange={(e) => setNewFlag({ ...newFlag, environment: e.target.value })}
            className="px-3 py-2 border border-gray-300 rounded-md text-sm"
          >
            <option value="development">Development</option>
            <option value="staging">Staging</option>
            <option value="production">Production</option>
          </select>
          <div className="flex gap-2">
            <button type="submit" className="px-4 py-2 bg-blue-600 text-white rounded-md text-sm" disabled={createFlag.isPending}>
              Create
            </button>
            <button type="button" onClick={() => setShowCreate(false)} className="px-4 py-2 text-sm text-gray-600">
              Cancel
            </button>
          </div>
        </form>
      )}

      <div className="bg-white border border-gray-200 rounded-lg overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="text-left px-4 py-3 font-medium text-gray-600">Flag</th>
              <th className="text-left px-4 py-3 font-medium text-gray-600">Environment</th>
              <th className="text-left px-4 py-3 font-medium text-gray-600">Status</th>
              <th className="text-left px-4 py-3 font-medium text-gray-600">Updated</th>
              <th className="text-right px-4 py-3 font-medium text-gray-600">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-100">
            {filtered.map((flag) => (
              <tr key={flag.id} className="hover:bg-gray-50">
                <td className="px-4 py-3">
                  <Link to={`/flags/${flag.key}`} className="font-medium text-blue-600 hover:underline">
                    {flag.key}
                  </Link>
                  <p className="text-gray-500 text-xs">{flag.name}</p>
                </td>
                <td className="px-4 py-3">
                  <EnvBadge env={flag.environment} />
                </td>
                <td className="px-4 py-3">
                  <Toggle enabled={flag.enabled} onChange={() => toggleFlag.mutate(flag.key)} />
                </td>
                <td className="px-4 py-3 text-gray-500 text-xs">{formatDate(flag.updatedAt)}</td>
                <td className="px-4 py-3 text-right">
                  <button onClick={() => setDeleteKey(flag.key)} className="text-red-500 hover:text-red-700 text-xs">
                    Delete
                  </button>
                </td>
              </tr>
            ))}
            {filtered.length === 0 && (
              <tr>
                <td colSpan={5} className="px-4 py-8 text-center text-gray-500">
                  No flags found.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <ConfirmDialog
        open={!!deleteKey}
        title="Delete flag"
        message={`Are you sure you want to delete "${deleteKey}"? This action cannot be undone.`}
        onConfirm={() => {
          if (deleteKey) deleteFlag.mutate(deleteKey)
          setDeleteKey(null)
        }}
        onCancel={() => setDeleteKey(null)}
      />
    </div>
  )
}
