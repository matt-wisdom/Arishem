<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useScansStore } from '../stores/scans'
import { useAuth } from '@clerk/vue'

const router = useRouter()
const scansStore = useScansStore()
const { scans, loading } = storeToRefs(scansStore)
const { getToken } = useAuth()

const searchQuery = ref('')
const typeFilter = ref('All Types')
const statusFilter = ref('All Status')

const showNewScanModal = ref(false)
const newScan = ref({ type: 'code' as 'code' | 'webapp', target: '', branch: 'main' })

const findingsCounts = ref<Record<string, { count: number, severity?: string }>>({})

const fetchFindingsCount = async (id: string) => {
  try {
    const tokenFn = typeof getToken.value === 'function' ? getToken.value : getToken
    const token = (await (tokenFn as any)()) || localStorage.getItem('token') || ''
    const res = await fetch(`/api/scans/${id}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (res.ok) {
      const data = await res.json()
      const count = data.findings ? data.findings.length : 0
      let maxSeverity = 'low'
      if (data.findings && data.findings.length > 0) {
        for (const f of data.findings) {
          if (f.severity === 'critical') maxSeverity = 'critical'
          else if (f.severity === 'high' && maxSeverity !== 'critical') maxSeverity = 'high'
          else if (f.severity === 'medium' && maxSeverity !== 'critical' && maxSeverity !== 'high') maxSeverity = 'medium'
        }
        findingsCounts.value[id] = { count, severity: maxSeverity }
      } else {
        findingsCounts.value[id] = { count }
      }
    }
  } catch (e) {
    console.error(e)
  }
}

const loadData = async () => {
  await scansStore.fetchScans()
  scans.value.forEach(s => {
    if (s.status === 'completed') {
      fetchFindingsCount(s.id)
    }
  })
}

onMounted(() => {
  loadData()
})

const getScanName = (target: string) => {
  return target.split('/').pop()?.replace('.git', '') || target
}

const getStatusClass = (status: string) => {
  const map: Record<string, string> = {
    completed: 'status-completed',
    running: 'status-running',
    queued: 'status-queued',
    failed: 'status-failed',
  }
  return map[status] || ''
}

const getTypeIcon = (type: string) => {
  return type === 'code' ? 'code' : 'web'
}

const handleCreateScan = async () => {
  try {
    await scansStore.createScan(newScan.value.type, newScan.value.target, newScan.value.branch)
    showNewScanModal.value = false
    newScan.value.target = ''
    newScan.value.branch = 'main'
    loadData()
  } catch (e) {
    alert(e instanceof Error ? e.message : 'Failed to create scan')
  }
}

const handleCancelScan = async (id: string) => {
  if (confirm('Are you sure you want to cancel this scan?')) {
    try {
      await scansStore.cancelScan(id)
      loadData()
    } catch (e) {
      alert(e instanceof Error ? e.message : 'Failed to cancel scan')
    }
  }
}

const filteredScans = computed(() => {
  return scans.value.filter(scan => {
    const name = getScanName(scan.target).toLowerCase()
    const target = scan.target.toLowerCase()
    const query = searchQuery.value.toLowerCase()
    const matchesSearch = name.includes(query) || target.includes(query)

    const matchesType = typeFilter.value === 'All Types' ||
      (typeFilter.value === 'Code Scan' && scan.type === 'code') ||
      (typeFilter.value === 'Webapp Scan' && scan.type === 'webapp')

    const matchesStatus = statusFilter.value === 'All Status' ||
      scan.status.toLowerCase() === statusFilter.value.toLowerCase()

    return matchesSearch && matchesType && matchesStatus
  })
})
</script>

<template>
  <div class="scans-page">
    <Header>
      <template #title>Scans</template>
    </Header>

    <div class="page-content">
      <!-- Toolbar -->
      <div class="toolbar">
        <div class="search-box">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8" />
            <path d="M21 21l-4.35-4.35" />
          </svg>
          <input v-model="searchQuery" type="text" placeholder="Search scans..." />
        </div>
        <div class="toolbar-actions">
          <select v-model="typeFilter" class="filter-select">
            <option>All Types</option>
            <option>Code Scan</option>
            <option>Webapp Scan</option>
          </select>
          <select v-model="statusFilter" class="filter-select">
            <option>All Status</option>
            <option>Completed</option>
            <option>Running</option>
            <option>Queued</option>
            <option>Failed</option>
          </select>
          <button class="btn-primary" @click="showNewScanModal = true">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 5v14M5 12h14" />
            </svg>
            New Scan
          </button>
        </div>
      </div>

      <!-- Scans List -->
      <div class="scans-list">
        <div v-if="loading && filteredScans.length === 0" class="loading-placeholder">Loading scans...</div>
        <div v-else-if="filteredScans.length === 0" class="empty-placeholder">No scans found</div>

        <div v-for="scan in filteredScans" :key="scan.id" class="scan-card" @click="router.push(`/scans/${scan.id}`)">
          <div class="scan-icon" :class="getTypeIcon(scan.type)">
            <svg v-if="scan.type === 'code'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M16 18l6-6-6-6M8 6l-6 6 6 6" />
            </svg>
            <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10" />
              <path d="M2 12h20M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z" />
            </svg>
          </div>
          <div class="scan-info">
            <h3 class="scan-name">{{ getScanName(scan.target) }}</h3>
            <p class="scan-meta">{{ scan.type === 'code' ? 'Code Scan' : 'Webapp Scan' }} • {{ new Date(scan.created_at).toLocaleString() }}</p>
          </div>
          <div class="scan-status">
            <span :class="['status-badge', getStatusClass(scan.status)]">
              <span class="status-dot"></span>
              {{ scan.status }}
            </span>
          </div>
          <div v-if="findingsCounts[scan.id]" class="scan-findings">
            <span :class="['findings-count', findingsCounts[scan.id].severity ? `severity-${findingsCounts[scan.id].severity}` : 'neutral']">
              {{ findingsCounts[scan.id].count }}
            </span>
            <span class="findings-label">findings</span>
          </div>
          <div v-else class="scan-findings">
            <span class="findings-count neutral">-</span>
            <span class="findings-label">findings</span>
          </div>
          <div class="scan-actions">
            <button v-if="scan.status === 'queued' || scan.status === 'running'" class="action-icon" @click.stop="handleCancelScan(scan.id)" title="Cancel Scan">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- New Scan Modal -->
      <div v-if="showNewScanModal" class="modal-overlay" @click.self="showNewScanModal = false">
        <div class="modal">
          <div class="modal-header">
            <h2>Create New Scan</h2>
            <button class="close-btn" @click="showNewScanModal = false">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 6L6 18M6 6l12 12" />
              </svg>
            </button>
          </div>
          <form @submit.prevent="handleCreateScan" class="modal-form">
            <div class="form-group">
              <label>Scan Type</label>
              <div class="type-selector">
                <button type="button" :class="['type-btn', { active: newScan.type === 'code' }]" @click="newScan.type = 'code'">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M16 18l6-6-6-6M8 6l-6 6 6 6" />
                  </svg>
                  Code Scan
                </button>
                <button type="button" :class="['type-btn', { active: newScan.type === 'webapp' }]" @click="newScan.type = 'webapp'">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10" />
                    <path d="M2 12h20M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z" />
                  </svg>
                  Webapp Scan
                </button>
              </div>
            </div>
            <div class="form-group">
              <label>Target</label>
              <input v-model="newScan.target" type="text" placeholder="Repository URL or target website" required />
            </div>
            <div v-if="newScan.type === 'code'" class="form-group">
              <label>Branch</label>
              <input v-model="newScan.branch" type="text" placeholder="main" />
            </div>
            <div class="modal-actions">
              <button type="button" class="btn-secondary" @click="showNewScanModal = false">Cancel</button>
              <button type="submit" class="btn-primary">Create Scan</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.scans-page {
  min-height: 100%;
}

.page-content {
  padding-top: 24px;
}

/* Toolbar */
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 12px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 0;
  padding: 12px 16px;
  width: 320px;
  transition: all 0.2s ease;
}

.search-box:focus-within {
  border-color: var(--accent);
  box-shadow: 0 0 10px var(--accent-glow);
}

.search-box svg {
  width: 18px;
  height: 18px;
  color: var(--text-muted);
}

.search-box input {
  background: none;
  border: none !important;
  outline: none;
  color: var(--text-primary);
  font-family: 'Rajdhani', sans-serif;
  font-size: 14px;
  width: 100%;
  padding: 0 !important;
}

.search-box input::placeholder {
  color: var(--text-muted);
}

.toolbar-actions {
  display: flex;
  gap: 12px;
}

.filter-select {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 0;
  padding: 10px 16px;
  color: var(--text-primary);
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.filter-select:hover {
  border-color: var(--accent);
  box-shadow: 0 0 8px var(--accent-glow);
}

.btn-primary {
  display: flex;
  align-items: center;
  gap: 8px;
  background: var(--accent-pink);
  color: white;
  border: 1px solid var(--accent-pink);
  border-radius: 0;
  padding: 10px 20px;
  font-family: 'Orbitron', sans-serif;
  font-size: 13px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1px;
  box-shadow: 0 0 10px var(--accent-pink-glow);
  transition: all 0.2s ease;
}

.btn-primary:hover {
  transform: translateY(-2px);
  background: #ff1a8c;
  box-shadow: 0 0 18px rgba(255, 0, 127, 0.6);
}

.btn-primary svg {
  width: 18px;
  height: 18px;
}

/* Scans List */
.scans-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.scan-card {
  display: flex;
  align-items: center;
  gap: 20px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 0;
  padding: 20px 24px;
  cursor: pointer;
  position: relative;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
}

.scan-card:hover {
  background: rgba(0, 255, 204, 0.04);
  border-color: var(--accent) !important;
  transform: translateX(6px);
  box-shadow: inset 4px 0 0 var(--accent), 0 0 15px var(--accent-glow) !important;
}

.scan-card::before {
  content: '';
  position: absolute;
  top: -1px;
  left: -1px;
  width: 8px;
  height: 8px;
  background: var(--accent);
  clip-path: polygon(0 0, 100% 0, 0 100%);
}

.scan-icon {
  width: 44px;
  height: 44px;
  border-radius: 0;
  border: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: center;
}

.scan-card:hover .scan-icon {
  border-color: var(--accent);
  box-shadow: 0 0 8px var(--accent-glow);
}

.scan-icon.code {
  background: rgba(139, 92, 246, 0.05);
  color: #8b5cf6;
}

.scan-icon.web {
  background: rgba(6, 182, 212, 0.05);
  color: #06b6d4;
}

.scan-icon svg {
  width: 20px;
  height: 20px;
}

.scan-info {
  flex: 1;
}

.scan-name {
  font-family: 'Orbitron', sans-serif;
  font-size: 15px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 4px;
  letter-spacing: 0.5px;
}

.scan-meta {
  font-family: 'Share Tech Mono', monospace;
  font-size: 13px;
  color: var(--text-muted);
}

.scan-status {
  width: 120px;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 4px 12px;
  border-radius: 0;
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border: 1px solid transparent;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 0;
}

.status-completed {
  background: rgba(57, 255, 20, 0.05);
  color: var(--success);
  border-color: rgba(57, 255, 20, 0.2);
}

.status-completed .status-dot {
  background: var(--success);
  box-shadow: 0 0 5px var(--success);
}

.status-running {
  background: rgba(0, 255, 204, 0.05);
  color: var(--accent);
  border-color: rgba(0, 255, 204, 0.2);
}

.status-running .status-dot {
  background: var(--accent);
  box-shadow: 0 0 5px var(--accent);
  animation: pulse 1.5s infinite;
}

.status-queued {
  background: rgba(255, 159, 0, 0.05);
  color: var(--warning);
  border-color: rgba(255, 159, 0, 0.2);
}

.status-queued .status-dot {
  background: var(--warning);
}

.status-failed {
  background: rgba(255, 7, 58, 0.05);
  color: var(--danger);
  border-color: rgba(255, 7, 58, 0.2);
}

.status-failed .status-dot {
  background: var(--danger);
  box-shadow: 0 0 5px var(--danger);
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}

.scan-findings {
  width: 80px;
  text-align: center;
}

.findings-count {
  display: block;
  font-family: 'Share Tech Mono', monospace;
  font-size: 22px;
  font-weight: 700;
}

.findings-count.severity-critical { color: var(--danger); text-shadow: 0 0 5px var(--danger); }
.findings-count.severity-high { color: var(--warning); }
.findings-count.severity-medium { color: var(--info); }
.findings-count.severity-low { color: #a0aec0; }
.findings-count.neutral { color: var(--text-muted); }

.findings-label {
  font-family: 'Orbitron', sans-serif;
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  color: var(--text-muted);
}

.scan-actions {
  width: 40px;
}

.action-icon {
  width: 32px;
  height: 32px;
  background: transparent;
  border: 1px solid var(--border-color);
  border-radius: 0;
  color: var(--text-muted);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.action-icon:hover {
  background: rgba(0, 255, 204, 0.05);
  border-color: var(--accent);
  color: var(--accent);
  box-shadow: 0 0 8px var(--accent-glow);
}

.action-icon svg {
  width: 16px;
  height: 16px;
}

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: #060613;
  border: 2px solid var(--accent);
  border-radius: 0;
  width: 500px;
  overflow: hidden;
  box-shadow: 0 0 25px var(--accent-glow);
  position: relative;
}

.modal::after {
  content: '';
  position: absolute;
  bottom: 0;
  right: 0;
  width: 20px;
  height: 20px;
  background: var(--accent-pink);
  clip-path: polygon(100% 100%, 100% 0, 0 100%);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 2px solid var(--border-color);
  background: rgba(0, 0, 0, 0.3);
}

.modal-header h2 {
  font-family: 'Orbitron', sans-serif;
  font-size: 16px;
  font-weight: 700;
  text-transform: uppercase;
  color: var(--accent);
  text-shadow: 0 0 5px var(--accent-glow);
}

.close-btn {
  width: 32px;
  height: 32px;
  background: transparent;
  border: 1px solid var(--border-color);
  border-radius: 0;
  color: var(--text-muted);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.close-btn:hover {
  border-color: var(--danger);
  color: var(--danger);
  box-shadow: 0 0 8px rgba(255, 7, 58, 0.4);
}

.close-btn svg {
  width: 16px;
  height: 16px;
}

.modal-form {
  padding: 24px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--text-muted);
  margin-bottom: 8px;
}

.form-group input {
  width: 100%;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid var(--border-color);
  border-radius: 0;
  padding: 12px 16px;
  color: var(--text-primary);
  font-family: 'Rajdhani', sans-serif;
  font-size: 15px;
}

.form-group input:focus {
  outline: none;
  border-color: var(--accent);
  box-shadow: 0 0 8px var(--accent-glow);
}

.type-selector {
  display: flex;
  gap: 12px;
}

.type-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 14px;
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid var(--border-color);
  border-radius: 0;
  color: var(--text-muted);
  font-family: 'Orbitron', sans-serif;
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  transition: all 0.2s ease;
}

.type-btn:hover {
  border-color: var(--text-primary);
  color: var(--text-primary);
}

.type-btn.active {
  border-color: var(--accent);
  background: rgba(0, 255, 204, 0.08);
  color: var(--accent);
  box-shadow: 0 0 8px var(--accent-glow);
}

.type-btn svg {
  width: 18px;
  height: 18px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 28px;
}

.btn-secondary {
  background: transparent;
  color: var(--text-muted);
  border: 1px solid var(--border-color);
  border-radius: 0;
  padding: 10px 20px;
  font-family: 'Orbitron', sans-serif;
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  transition: all 0.2s ease;
}

.btn-secondary:hover {
  border-color: var(--text-primary);
  color: var(--text-primary);
}
</style>