<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuth } from '@clerk/vue'
import { apiUrl, apiFetch } from '@/utils/api'

const route = useRoute()
const router = useRouter()
const { getToken } = useAuth()

const scan = ref<any>({
  id: route.params.id,
  target: '',
  type: 'code',
  status: 'queued',
  branch: '',
  created_at: '',
  completed_at: '',
})

const findings = ref<any[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

const getSeverityClass = (s: string) => `severity-${s}`

const getScanName = (target: string) => {
  return target.split('/').pop()?.replace('.git', '') || target
}

const loadScanDetail = async () => {
  loading.value = true
  error.value = null
  try {
    const tokenFn = typeof getToken.value === 'function' ? getToken.value : getToken
    const token = (await (tokenFn as any)()) || localStorage.getItem('token') || ''
    const res = await apiFetch(apiUrl(`/scans/${route.params.id}`), {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (!res.ok) throw new Error('Failed to load scan details')
    const data = await res.json()
    scan.value = data.scan
    findings.value = data.findings || []
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Unknown error'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadScanDetail()
})
</script>

<template>
  <div class="scan-detail">
    <Header>
      <template #title>
        <button class="back-btn" @click="router.push('/scans')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 12H5M12 19l-7-7 7-7"/></svg>
        </button>
        Scan Details
      </template>
    </Header>
    <div class="page-content">
      <div v-if="loading" class="loading-placeholder">Loading scan details...</div>
      <div v-else-if="error" class="error-placeholder">{{ error }}</div>
      <div v-else>
        <div class="scan-header">
          <div class="scan-title">
            <h1>{{ getScanName(scan.target) }}</h1>
            <span :class="['status-badge', scan.status]">{{ scan.status }}</span>
          </div>
          <p class="scan-target">{{ scan.target }} <span v-if="scan.branch">({{ scan.branch }})</span></p>
        </div>

        <div class="stats-row">
          <div class="stat-box">
            <span class="label">Started</span>
            <span class="value">{{ new Date(scan.created_at).toLocaleString() }}</span>
          </div>
          <div class="stat-box">
            <span class="label">Completed</span>
            <span class="value">{{ scan.completed_at ? new Date(scan.completed_at).toLocaleString() : '-' }}</span>
          </div>
          <div class="stat-box">
            <span class="label">Findings</span>
            <span class="value highlight">{{ findings.length }}</span>
          </div>
        </div>

        <div class="findings-section">
          <h2>Findings</h2>
          <div v-if="findings.length === 0" class="empty-placeholder">No findings found for this scan.</div>
          <div v-else class="findings-list">
            <div v-for="finding in findings" :key="finding.id" class="finding-card">
              <div class="finding-header">
                <span :class="['severity-badge', getSeverityClass(finding.severity)]">{{ finding.severity }}</span>
                <h3>{{ finding.title }}</h3>
              </div>
              <p class="finding-desc">{{ finding.description }}</p>
              <div class="finding-remediation">
                <strong>Remediation:</strong> {{ finding.remediation }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.scan-detail { min-height: 100%; }
.page-content { padding-top: 24px; }

.back-btn {
  background: transparent;
  border: 1px solid var(--border-color);
  border-radius: 0;
  color: var(--text-muted);
  cursor: pointer;
  padding: 6px 12px;
  margin-right: 12px;
  transition: all 0.2s ease;
}

.back-btn:hover {
  border-color: var(--accent);
  color: var(--accent);
  box-shadow: 0 0 8px var(--accent-glow);
}

.back-btn svg { width: 16px; height: 16px; }
.scan-header { margin-bottom: 24px; }
.scan-title { display: flex; align-items: center; gap: 16px; margin-bottom: 8px; }

.scan-title h1 {
  font-family: 'Orbitron', sans-serif;
  font-size: 24px;
  font-weight: 700;
  text-transform: uppercase;
  color: var(--text-primary);
  letter-spacing: 1px;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 0;
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border: 1px solid transparent;
}

.status-badge.completed {
  background: rgba(57, 255, 20, 0.05);
  color: var(--success);
  border-color: rgba(57, 255, 20, 0.2);
}

.status-badge.running {
  background: rgba(0, 255, 204, 0.05);
  color: var(--accent);
  border-color: rgba(0, 255, 204, 0.2);
}

.scan-target {
  font-family: 'Share Tech Mono', monospace;
  color: var(--text-muted);
  font-size: 14px;
}

.stats-row { display: flex; gap: 24px; margin-bottom: 32px; }

.stat-box {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 0;
  padding: 20px 24px;
  flex: 1;
  position: relative;
}

.stat-box::before {
  content: '';
  position: absolute;
  top: -1px;
  left: -1px;
  width: 8px;
  height: 8px;
  background: var(--accent);
  clip-path: polygon(0 0, 100% 0, 0 100%);
}

.stat-box .label {
  display: block;
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  font-weight: 700;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 4px;
}

.stat-box .value {
  font-family: 'Share Tech Mono', monospace;
  font-size: 16px;
  font-weight: 700;
}

.stat-box .value.highlight {
  color: var(--accent);
  font-size: 26px;
  text-shadow: 0 0 5px var(--accent-glow);
}

.findings-section h2 {
  font-family: 'Orbitron', sans-serif;
  font-size: 14px;
  font-weight: 700;
  color: var(--accent-pink);
  margin-bottom: 16px;
  text-transform: uppercase;
  letter-spacing: 1px;
  text-shadow: 0 0 5px var(--accent-pink-glow);
}

.findings-list { display: flex; flex-direction: column; gap: 16px; }

.finding-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 0;
  padding: 24px;
  position: relative;
  transition: all 0.2s ease;
}

.finding-card:hover {
  border-color: var(--accent-pink) !important;
  box-shadow: 0 0 12px var(--accent-pink-glow) !important;
}

.finding-card::before {
  content: '';
  position: absolute;
  top: -1px;
  left: -1px;
  width: 10px;
  height: 10px;
  background: var(--accent-pink);
  clip-path: polygon(0 0, 100% 0, 0 100%);
}

.finding-header { display: flex; align-items: center; gap: 12px; margin-bottom: 12px; }

.finding-header h3 {
  font-family: 'Orbitron', sans-serif;
  font-size: 15px;
  font-weight: 700;
  letter-spacing: 0.5px;
}

.severity-badge {
  padding: 4px 10px;
  border-radius: 0;
  font-family: 'Orbitron', sans-serif;
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  border: 1px solid transparent;
}

.severity-badge.severity-critical {
  background: rgba(255, 7, 58, 0.05);
  color: var(--danger);
  border-color: rgba(255, 7, 58, 0.2);
}

.severity-badge.severity-high {
  background: rgba(255, 159, 0, 0.05);
  color: var(--warning);
  border-color: rgba(255, 159, 0, 0.2);
}

.severity-badge.severity-medium {
  background: rgba(0, 229, 255, 0.05);
  color: var(--info);
  border-color: rgba(0, 229, 255, 0.2);
}

.severity-badge.severity-low {
  background: rgba(98, 114, 164, 0.05);
  color: #a0aec0;
  border-color: rgba(98, 114, 164, 0.2);
}

.finding-desc { color: var(--text-primary); font-size: 14px; margin-bottom: 12px; }
.finding-location { margin-bottom: 12px; }

.finding-location code {
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid var(--border-color);
  padding: 6px 12px;
  border-radius: 0;
  font-family: 'Share Tech Mono', monospace;
  font-size: 13px;
  color: var(--accent);
}

.finding-remediation { font-size: 14px; color: var(--text-secondary); }
.finding-remediation strong {
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  font-weight: 700;
  color: var(--accent-pink);
  text-transform: uppercase;
  margin-right: 4px;
}
</style>