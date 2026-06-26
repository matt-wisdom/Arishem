<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { Scan, Finding } from '@/stores/types'

const route = useRoute()
const router = useRouter()
const scan = ref<Scan | null>(null)
const findings = ref<Finding[]>([])
const loading = ref(true)

const fetchScan = async () => {
  try {
    const res = await fetch(`/api/scans/${route.params.id}`, {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
    if (!res.ok) throw new Error('Failed to fetch scan')
    const data = await res.json()
    scan.value = data.scan
    findings.value = data.findings || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const cancelScan = async () => {
  if (!confirm('Are you sure you want to cancel this scan?')) return
  try {
    const res = await fetch(`/api/scans/${route.params.id}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
    if (res.ok) await fetchScan()
  } catch (e) {
    console.error(e)
  }
}

onMounted(fetchScan)
</script>

<template>
  <div class="scan-detail">
    <header class="header">
      <button @click="router.push('/scans')">&larr; Back to Scans</button>
      <h1>Scan Details</h1>
    </header>

    <div v-if="loading" class="loading">Loading...</div>
    <div v-else-if="scan" class="content">
      <div class="scan-info">
        <p><strong>Type:</strong> {{ scan.type }}</p>
        <p><strong>Target:</strong> {{ scan.target }}</p>
        <p><strong>Status:</strong> <span :class="['status', scan.status]">{{ scan.status }}</span></p>
        <p><strong>Created:</strong> {{ new Date(scan.created_at).toLocaleString() }}</p>
        <p v-if="scan.completed_at"><strong>Completed:</strong> {{ new Date(scan.completed_at).toLocaleString() }}</p>
        <button v-if="scan.status === 'queued' || scan.status === 'running'" @click="cancelScan" class="danger-btn">Cancel Scan</button>
      </div>

      <section class="findings">
        <h2>Findings ({{ findings.length }})</h2>
        <div v-for="finding in findings" :key="finding.id" class="finding-card">
          <div class="finding-header">
            <h3>{{ finding.title }}</h3>
            <span :class="['severity', finding.severity]">{{ finding.severity }}</span>
          </div>
          <p>{{ finding.description }}</p>
          <p v-if="finding.remediation"><strong>Remediation:</strong> {{ finding.remediation }}</p>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.scan-detail { padding: 2rem; max-width: 1000px; margin: 0 auto; }
.header { margin-bottom: 2rem; }
.header button { background: none; border: none; cursor: pointer; color: #666; margin-bottom: 0.5rem; }
.scan-info { background: white; padding: 1.5rem; border-radius: 8px; margin-bottom: 2rem; }
.scan-info p { margin-bottom: 0.5rem; }
.status { padding: 0.25rem 0.5rem; border-radius: 4px; font-size: 0.75rem; }
.status.queued { background: #fef3c7; color: #92400e; }
.status.running { background: #dbeafe; color: #1e40af; }
.status.completed { background: #d1fae5; color: #065f46; }
.status.failed { background: #fee2e2; color: #991b1b; }
.danger-btn { background: #dc2626; color: white; border: none; padding: 0.5rem 1rem; border-radius: 6px; cursor: pointer; margin-top: 1rem; }
.findings h2 { margin-bottom: 1rem; }
.finding-card { background: white; padding: 1.5rem; border-radius: 8px; margin-bottom: 1rem; }
.finding-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem; }
.severity { padding: 0.25rem 0.5rem; border-radius: 4px; font-size: 0.75rem; text-transform: uppercase; }
.severity.critical { background: #7f1d1d; color: white; }
.severity.high { background: #dc2626; color: white; }
.severity.medium { background: #f59e0b; color: white; }
.severity.low { background: #3b82f6; color: white; }
.severity.info { background: #6b7280; color: white; }
</style>