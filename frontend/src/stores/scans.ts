import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Scan } from './types'

export const useScansStore = defineStore('scans', () => {
  const scans = ref<Scan[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchScans = async () => {
    loading.value = true
    error.value = null
    try {
      const res = await fetch('/api/scans', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      })
      if (!res.ok) throw new Error('Failed to fetch scans')
      scans.value = await res.json()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Unknown error'
    } finally {
      loading.value = false
    }
  }

  const getScan = (id: string) => scans.value.find(s => s.id === id)

  return { scans, loading, error, fetchScans, getScan }
})