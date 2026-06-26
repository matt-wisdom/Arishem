<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { Report } from '@/stores/types'

const reports = ref<Report[]>([])
const loading = ref(true)

const fetchReports = async () => {
  try {
    const res = await fetch('/api/reports', {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
    if (!res.ok) throw new Error('Failed to fetch reports')
    reports.value = await res.json()
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const downloadReport = async (report: Report, format: string) => {
  try {
    const res = await fetch(`/api/reports/${report.id}/download?format=${format}`, {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
    if (!res.ok) throw new Error('Failed to get download URL')
    const { url } = await res.json()
    window.open(url, '_blank')
  } catch (e) {
    console.error(e)
  }
}

onMounted(fetchReports)
</script>

<template>
  <div class="reports-page">
    <header class="header">
      <h1>Reports</h1>
    </header>

    <div v-if="loading" class="loading">Loading...</div>
    <table v-else class="reports-table">
      <thead>
        <tr>
          <th>ID</th>
          <th>Format</th>
          <th>Created</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="report in reports" :key="report.id">
          <td>{{ report.id }}</td>
          <td>{{ report.format.toUpperCase() }}</td>
          <td>{{ new Date(report.created_at).toLocaleDateString() }}</td>
          <td>
            <div class="download-btns">
              <button @click="downloadReport(report, 'html')">HTML</button>
              <button @click="downloadReport(report, 'md')">MD</button>
              <button @click="downloadReport(report, 'sarif')">SARIF</button>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.reports-page { padding: 2rem; max-width: 1200px; margin: 0 auto; }
.header { margin-bottom: 2rem; }
.reports-table { width: 100%; background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
.reports-table th, .reports-table td { padding: 1rem; text-align: left; border-bottom: 1px solid #eee; }
.download-btns { display: flex; gap: 0.5rem; }
.download-btns button { padding: 0.25rem 0.5rem; background: #f3f4f6; border: none; border-radius: 4px; cursor: pointer; font-size: 0.75rem; }
.download-btns button:hover { background: #e5e7eb; }
</style>