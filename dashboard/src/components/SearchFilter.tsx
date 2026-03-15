interface SearchFilterProps {
  search: string
  onSearchChange: (v: string) => void
  environment: string
  onEnvironmentChange: (v: string) => void
}

export function SearchFilter({ search, onSearchChange, environment, onEnvironmentChange }: SearchFilterProps) {
  return (
    <div className="flex gap-3">
      <input
        type="text"
        placeholder="Search flags..."
        value={search}
        onChange={(e) => onSearchChange(e.target.value)}
        className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 w-64"
      />
      <select
        value={environment}
        onChange={(e) => onEnvironmentChange(e.target.value)}
        className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="">All environments</option>
        <option value="production">Production</option>
        <option value="staging">Staging</option>
        <option value="development">Development</option>
      </select>
    </div>
  )
}
