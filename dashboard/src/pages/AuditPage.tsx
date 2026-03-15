import { useState } from 'react'
import { useAuditLogs } from '../hooks/useAudit'
import { formatDate } from '../lib/utils'

export function AuditPage() {
  const [page, setPage] = useState(1)
  const { data, isLoading } = useAuditLogs({ page, perPage: 20 })

  if (isLoading) return <p className="text-gray-500">Loading...</p>

  const logs = data?.data ?? []
  const total = data?.total ?? 0
  const totalPages = Math.ceil(total / 20)

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Audit Log</h1>

      {logs.length === 0 ? (
        <p className="text-gray-500 text-center py-8">No audit entries yet.</p>
      ) : (
        <div className="bg-white border border-gray-200 rounded-lg overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Action</th>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Changes</th>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Time</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100">
              {logs.map((log) => (
                <tr key={log.id} className="hover:bg-gray-50">
                  <td className="px-4 py-3">
                    <span className="bg-gray-100 px-2 py-0.5 rounded text-xs font-medium">{log.action}</span>
                  </td>
                  <td className="px-4 py-3">
                    {log.diff && Object.keys(log.diff).length > 0 ? (
                      <div className="space-y-1">
                        {Object.entries(log.diff).map(([field, change]) => (
                          <div key={field} className="font-mono text-xs">
                            <span className="text-gray-600">{field}:</span>{' '}
                            {change.from == null ? (
                              <span className="text-green-600">{JSON.stringify(change.to)}</span>
                            ) : change.to == null ? (
                              <span className="text-red-500 line-through">{JSON.stringify(change.from)}</span>
                            ) : (
                              <>
                                <span className="text-red-500">{JSON.stringify(change.from)}</span>
                                {' → '}
                                <span className="text-green-600">{JSON.stringify(change.to)}</span>
                              </>
                            )}
                          </div>
                        ))}
                      </div>
                    ) : (
                      <span className="text-gray-400 text-xs">—</span>
                    )}
                  </td>
                  <td className="px-4 py-3 text-gray-500 text-xs">{formatDate(log.createdAt)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {totalPages > 1 && (
        <div className="flex justify-center gap-2">
          <button onClick={() => setPage((p) => Math.max(1, p - 1))} disabled={page === 1} className="px-3 py-1 text-sm border rounded disabled:opacity-50">
            Previous
          </button>
          <span className="px-3 py-1 text-sm text-gray-600">
            Page {page} of {totalPages}
          </span>
          <button onClick={() => setPage((p) => p + 1)} disabled={page >= totalPages} className="px-3 py-1 text-sm border rounded disabled:opacity-50">
            Next
          </button>
        </div>
      )}
    </div>
  )
}
