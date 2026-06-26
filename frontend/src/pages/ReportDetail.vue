<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const report = ref<{ id: string; format: string; created_at: string } | null>(null)
const downloadUrl = ref('')
const loading = ref(true)

const fetchReport = async () => {
  try {
    const res = await fetch(`/api/reports/${route.params.id}`, {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
    if (!res.ok) throw new Error('Failed to fetch report')
    report.value = await res.json()
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const getDownloadUrl = async (format: string) => {
  try {
    const res = await fetch(`/api/reports/${route.params.id}/download?format=${format}`, {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
    if (!res.ok) throw new Error('Failed to get download URL')
    const data = await res.json()
    downloadUrl.value = data.url
    window.open(data.url, '_blank')
  } catch (e) {
    console.error(e)
  }
}

onMounted(fetchReport)
</script>

<template>
  <div class="report-detail">
    <header class="header">
      <button @click="router.push('/reports')">&larr; Back to Reports</button>
      <h1>Report Details</h1>
    </header>

    <div v-if="loading" class="loading">Loading...</div>
    <div v-else-if="report" class="content">
      <div class="report-info">
        <p><strong>ID:</strong> {{ report.id }}</p>
        <p><strong>Format:</strong> {{ report.format.toUpperCase() }}</p>
        <p><strong>Created:</strong> {{ new Date(report.created_at).toLocaleString() }}</p>
      </div>

      <div class="download-section">
        <h2>Download Report</h2>
        <div class="download-btns">
          <button @click="getDownloadUrl('html')">Download HTML</button>
          <button @click="getDownloadUrl('md')">Download Markdown</button>
          <button @click="getDownloadUrl('sarif')">Download SARIF</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.report-detail { padding: 2rem; max-width: 800px; margin: 0 auto; }
.header { margin-bottom: 2rem; }
.header button { background: none; border: none; cursor: pointer; color: #666; margin-bottom: 0.5rem; }
.report-info { background: white; padding: 1.5rem; border-radius: 8px; margin-bottom: 2rem; }
.report-info p { margin-bottom: 0.5rem; }
.download-section { background: white; padding: 1.5rem; border-radius: 8px; }
.download-section h2 { margin-bottom: 1rem; }
.download-btns { display: flex; gap: 1rem; }
.download-btns button { padding: 0.75rem 1.5rem; background: #2563eb; color: white; border: none; border-radius: 6px; cursor: pointer; }
.download-btns button:hover { background: #1d4ed8; }
</style>