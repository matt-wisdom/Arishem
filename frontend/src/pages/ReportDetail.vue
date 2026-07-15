<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const report = ref({
  id: route.params.id,
  name: 'auth-service-scan',
  type: 'Code Scan',
  format: 'HTML',
  createdAt: '2024-01-15 10:45:23',
  size: '2.4 MB',
})

const downloadFormat = (format: string) => {
  console.log(`Downloading as ${format}`)
}
</script>

<template>
  <div class="report-detail">
    <Header>
      <template #title>
        <button class="back-btn" @click="router.push('/reports')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 12H5M12 19l-7-7 7-7"/></svg>
        </button>
        Report Details
      </template>
    </Header>
    <div class="page-content">
      <div class="report-header">
        <div class="report-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" />
            <path d="M14 2v6h6M16 13H8M16 17H8" />
          </svg>
        </div>
        <div class="report-info">
          <h1>{{ report.name }}</h1>
          <p>{{ report.type }} • {{ report.createdAt }} • {{ report.size }}</p>
        </div>
      </div>

      <div class="download-section">
        <h2>Download Report</h2>
        <div class="download-buttons">
          <button @click="downloadFormat('html')" class="download-btn">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M7 10l5 5 5-5M12 15V3" />
            </svg>
            HTML
          </button>
          <button @click="downloadFormat('md')" class="download-btn">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M7 10l5 5 5-5M12 15V3" />
            </svg>
            Markdown
          </button>
          <button @click="downloadFormat('sarif')" class="download-btn">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M7 10l5 5 5-5M12 15V3" />
            </svg>
            SARIF
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.report-detail { min-height: 100%; }
.page-content { padding-top: 24px; }
.back-btn { background: none; border: none; color: var(--text-muted); cursor: pointer; padding: 8px; margin-right: 8px; }
.back-btn svg { width: 20px; height: 20px; }
.report-header { display: flex; align-items: center; gap: 20px; background: var(--bg-card); border: 1px solid var(--border-color); border-radius: 16px; padding: 24px; margin-bottom: 24px; }
.report-icon { width: 64px; height: 64px; background: rgba(6, 182, 212, 0.15); border-radius: 16px; display: flex; align-items: center; justify-content: center; }
.report-icon svg { width: 32px; height: 32px; color: #06b6d4; }
.report-info h1 { font-size: 24px; font-weight: 700; margin-bottom: 4px; }
.report-info p { color: var(--text-muted); font-size: 14px; }
.download-section h2 { font-size: 18px; font-weight: 600; margin-bottom: 16px; }
.download-buttons { display: flex; gap: 12px; }
.download-btn { display: flex; align-items: center; gap: 10px; padding: 16px 32px; background: var(--bg-card); border: 1px solid var(--border-color); border-radius: 12px; color: var(--text-primary); font-size: 14px; font-weight: 600; }
.download-btn:hover { background: var(--bg-card-hover); border-color: var(--accent); }
.download-btn svg { width: 20px; height: 20px; }
</style>