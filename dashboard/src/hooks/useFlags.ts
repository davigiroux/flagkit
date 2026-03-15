import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { api, type Flag } from '../api/client'

export function useFlags() {
  return useQuery({ queryKey: ['flags'], queryFn: api.listFlags })
}

export function useFlag(key: string) {
  return useQuery({ queryKey: ['flags', key], queryFn: () => api.getFlag(key), enabled: !!key })
}

export function useCreateFlag() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (data: Partial<Flag>) => api.createFlag(data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['flags'] }),
  })
}

export function useUpdateFlag() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ key, data }: { key: string; data: Partial<Flag> }) => api.updateFlag(key, data),
    onSuccess: (_, { key }) => {
      qc.invalidateQueries({ queryKey: ['flags'] })
      qc.invalidateQueries({ queryKey: ['flags', key] })
    },
  })
}

export function useDeleteFlag() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (key: string) => api.deleteFlag(key),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['flags'] }),
  })
}

export function useToggleFlag() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (key: string) => api.toggleFlag(key),
    onSuccess: (flag) => {
      qc.setQueryData(['flags', flag.key], flag)
      qc.invalidateQueries({ queryKey: ['flags'] })
    },
  })
}
