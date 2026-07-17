<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useAuth } from '@clerk/vue'
import { apiUrl, apiFetch } from '@/utils/api'

const { getToken } = useAuth()
const reports = ref<any[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

const currentPage = ref(1)
const itemsPerPage = 6

const totalPages = computed(() => {
  return Math.ceil(reports.value.length / itemsPerPage)
})

const paginatedReports = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage
  const end = start + itemsPerPage
  return reports.value.slice(start, end)
})

const getHeaders = async () => {
  const tokenFn = typeof getToken.value === 'function' ? getToken.value : getToken
  const token = (await (tokenFn as any)()) || localStorage.getItem('token') || ''
  return { 'Authorization': `Bearer ${token}` }
}

const loadReports = async () => {
  loading.value = true
  error.value = null
  try {
    const headers = await getHeaders()
    const res = await apiFetch(apiUrl('/reports'), { headers })
    if (!res.ok) throw new Error('Failed to fetch reports')
    reports.value = await res.json()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Unknown error'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadReports()
})

const getReportName = (report: any) => {
  if (report.storage_key) {
    const parts = report.storage_key.split('/')
    return parts[parts.length - 1] || 'report-' + report.id
  }
  return 'report-' + report.id
}

const getReportType = (report: any) => {
  return report.scan_id ? 'Scan Report' : 'LLM Pentest Report'
}

const downloadReport = async (id: string, format: string) => {
  try {
    const headers = await getHeaders()
    const res = await apiFetch(apiUrl(`/reports/${id}/download?format=${format}`), { headers })
    if (!res.ok) throw new Error('Failed to generate download link')
    const data = await res.json()
    if (data.url) {
      window.open(data.url, '_blank')
    }
  } catch (e) {
    alert(e instanceof Error ? e.message : 'Failed to download report')
  }
}
</script>

<template>
  <div class="reports-page">
    <Header><template #title>Reports</template></Header>
    <div class="page-content">
      <div v-if="loading && reports.length === 0" class="loading-placeholder">Loading reports...</div>
      <div v-else-if="reports.length === 0" class="empty-placeholder">No reports generated yet.</div>
      <div v-else class="reports-grid">
        <div v-for="report in paginatedReports" :key="report.id" class="report-card">
          <div class="report-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" />
              <path d="M14 2v6h6M16 13H8M16 17H8" />
            </svg>
          </div>
          <div class="report-info">
            <h3>{{ getReportName(report) }}</h3>
            <p>{{ getReportType(report) }} • {{ new Date(report.created_at).toLocaleString() }}</p>
          </div>
          <div class="report-meta">
            <span class="format-badge">{{ report.format.toUpperCase() }}</span>
          </div>
          <div class="report-actions">
            <button @click="downloadReport(report.id, 'html')" class="action-btn">HTML</button>
            <button @click="downloadReport(report.id, 'md')" class="action-btn">MD</button>
            <button @click="downloadReport(report.id, 'sarif')" class="action-btn">SARIF</button>
          </div>
        </div>

        <!-- Pagination Controls -->
        <div v-if="totalPages > 1" class="pagination-controls" style="display: flex; justify-content: center; align-items: center; gap: 16px; margin-top: 32px;">
          <button 
            class="btn-secondary" 
            :disabled="currentPage === 1" 
            @click="currentPage--"
            style="padding: 8px 16px; font-size: 13px;"
          >
            Previous
          </button>
          <span style="font-family: 'Share Tech Mono', monospace; color: var(--text-muted); font-size: 14px;">
            Page <strong style="color: var(--accent);">{{ currentPage }}</strong> of <strong>{{ totalPages }}</strong>
          </span>
          <button 
            class="btn-secondary" 
            :disabled="currentPage === totalPages" 
            @click="currentPage++"
            style="padding: 8px 16px; font-size: 13px;"
          >
            Next
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.reports-page { min-height: 100%; }
.page-content { padding-top: 24px; }
.reports-grid { display: flex; flex-direction: column; gap: 12px; }
.report-card { display: flex; align-items: center; gap: 20px; background: var(--bg-card); border: 1px solid var(--border-color); border-radius: 14px; padding: 20px 24px; transition: all 0.2s ease; }
.report-card:hover { background: var(--bg-card-hover); }
.report-icon { width: 48px; height: 48px; background: rgba(6, 182, 212, 0.15); border-radius: 12px; display: flex; align-items: center; justify-content: center; }
.report-icon svg { width: 24px; height: 24px; color: #06b6d4; }
.report-info { flex: 1; }
.report-info h3 { font-size: 16px; font-weight: 600; color: var(--text-primary); margin-bottom: 4px; }
.report-info p { font-size: 13px; color: var(--text-muted); }
.report-meta { display: flex; align-items: center; gap: 12px; }
.format-badge { padding: 4px 10px; background: var(--bg-secondary); border-radius: 6px; font-size: 12px; font-weight: 600; color: var(--text-secondary); }
.size { font-size: 13px; color: var(--text-muted); width: 70px; text-align: right; }
.report-actions { display: flex; gap: 8px; }
.action-btn { padding: 8px 16px; background: var(--bg-secondary); border: 1px solid var(--border-color); border-radius: 8px; color: var(--text-primary); font-size: 13px; font-weight: 500; }
.action-btn:hover { background: var(--bg-card-hover); border-color: var(--accent); }
</style>