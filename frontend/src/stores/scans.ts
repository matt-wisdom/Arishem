import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useAuth } from '@clerk/vue'
import type { Scan } from './types'

export const useScansStore = defineStore('scans', () => {
  const scans = ref<Scan[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  
  const { getToken } = useAuth()

  const getHeaders = async () => {
    const tokenFn = typeof getToken.value === 'function' ? getToken.value : getToken
    const token = (await (tokenFn as any)()) || localStorage.getItem('token') || ''
    return {
      'Authorization': `Bearer ${token}`
    }
  }

  const fetchScans = async () => {
    loading.value = true
    error.value = null
    try {
      const headers = await getHeaders()
      const res = await fetch('/api/scans', { headers })
      if (!res.ok) throw new Error('Failed to fetch scans')
      scans.value = await res.json()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Unknown error'
    } finally {
      loading.value = false
    }
  }

  const createScan = async (type: 'code' | 'webapp', target: string, branch?: string) => {
    loading.value = true
    error.value = null
    try {
      const url = type === 'code' ? '/api/scans/code' : '/api/scans/webapp'
      const body = type === 'code' ? { target, branch: branch || 'main' } : { target }
      const authHeaders = await getHeaders()
      const res = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...authHeaders
        },
        body: JSON.stringify(body)
      })
      if (!res.ok) throw new Error('Failed to create scan')
      await fetchScans()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Unknown error'
      throw e
    } finally {
      loading.value = false
    }
  }

  const cancelScan = async (id: string) => {
    loading.value = true
    error.value = null
    try {
      const authHeaders = await getHeaders()
      const res = await fetch(`/api/scans/${id}`, {
        method: 'DELETE',
        headers: authHeaders
      })
      if (!res.ok) throw new Error('Failed to cancel scan')
      await fetchScans()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Unknown error'
      throw e
    } finally {
      loading.value = false
    }
  }

  const getScan = (id: string) => scans.value.find(s => s.id === id)

  return { scans, loading, error, fetchScans, createScan, cancelScan, getScan }
})