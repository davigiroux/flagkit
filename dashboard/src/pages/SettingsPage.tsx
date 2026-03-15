import { useState } from 'react'
import { getStoredApiKey, setApiKey } from '../api/client'

export function SettingsPage() {
  const [key, setKey] = useState(getStoredApiKey())
  const [saved, setSaved] = useState(false)

  function handleSave() {
    setApiKey(key)
    setSaved(true)
    setTimeout(() => setSaved(false), 2000)
  }

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Settings</h1>

      <div className="bg-white border border-gray-200 rounded-lg p-6 max-w-lg space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">API Key</label>
          <input
            type="password"
            value={key}
            onChange={(e) => { setKey(e.target.value); setSaved(false) }}
            placeholder="fk_..."
            className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm font-mono"
          />
          <p className="text-xs text-gray-500 mt-1">Stored in localStorage. Used for all API requests.</p>
        </div>
        <button onClick={handleSave} className="px-4 py-2 bg-blue-600 text-white rounded-lg text-sm font-medium hover:bg-blue-700">
          {saved ? 'Saved!' : 'Save'}
        </button>
      </div>
    </div>
  )
}
