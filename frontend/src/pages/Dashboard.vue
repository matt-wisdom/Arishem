<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@clerk/vue'

const router = useRouter()
const { getToken } = useAuth()



const recentScans = ref<any[]>([])

const threatDistribution = ref([
  { label: 'Critical', count: 0, color: 'var(--danger)' },
  { label: 'High', count: 0, color: 'var(--warning)' },
  { label: 'Medium', count: 0, color: 'var(--info)' },
  { label: 'Low', count: 0, color: '#a0aec0' },
])

const totalFindings = ref(0)
const loading = ref(true)

const getSeverityClass = (severity: string) => {
  const map: Record<string, string> = {
    critical: 'severity-critical',
    high: 'severity-high',
    medium: 'severity-medium',
    low: 'severity-low',
  }
  return map[severity] || ''
}



const loadDashboardData = async () => {
  loading.value = true
  try {
    const tokenFn = typeof getToken.value === 'function' ? getToken.value : getToken
    const token = (await (tokenFn as any)()) || localStorage.getItem('token') || ''
    const headers = { 'Authorization': `Bearer ${token}` }
    
    // Fetch runs
    const runsRes = await fetch('/api/llmpentest', { headers })
    const runs = runsRes.ok ? await runsRes.json() : []
    
    // Compute unique targets (Protected APIs)
    const targets = new Set<string>()
    runs.forEach((r: any) => targets.add(r.target_endpoint))
    
    // Load findings in background for completed runs
    const findingsList: any[] = []
    
    const fetchPromises = [
      ...runs.filter((r: any) => r.status === 'completed').map(async (r: any) => {
        try {
          const detailRes = await fetch(`/api/llmpentest/${r.id}`, { headers })
          if (detailRes.ok) {
            const data = await detailRes.json()
            if (data.findings) findingsList.push(...data.findings)
          }
        } catch {}
      })
    ]
    
    await Promise.all(fetchPromises)
    
    // Count severities
    let criticalCount = 0
    let highCount = 0
    let mediumCount = 0
    let lowCount = 0
    
    findingsList.forEach(f => {
      if (f.severity === 'critical') criticalCount++
      else if (f.severity === 'high') highCount++
      else if (f.severity === 'medium') mediumCount++
      else if (f.severity === 'low') lowCount++
    })
    
    threatDistribution.value[0].count = criticalCount
    threatDistribution.value[1].count = highCount
    threatDistribution.value[2].count = mediumCount
    threatDistribution.value[3].count = lowCount
    
    totalFindings.value = findingsList.length
    

    
    // Combine recent scans
    const combined = [
      ...runs.map((r: any) => ({
        id: r.id,
        name: r.target_endpoint,
        type: 'LLM Pentest',
        status: r.status,
        date: new Date(r.created_at).toLocaleDateString(),
        created_at: r.created_at,
        findings: null,
        severity: null
      }))
    ]
    
    // Decorate recent with findings counts if details loaded
    combined.forEach(item => {
      const itemFindings = findingsList.filter(f => f.run_id === item.id)
      if (item.status === 'completed') {
        item.findings = itemFindings.length
        let maxSeverity = 'low'
        if (itemFindings.length > 0) {
          for (const f of itemFindings) {
            if (f.severity === 'critical') maxSeverity = 'critical'
            else if (f.severity === 'high' && maxSeverity !== 'critical') maxSeverity = 'high'
            else if (f.severity === 'medium' && maxSeverity !== 'critical' && maxSeverity !== 'high') maxSeverity = 'medium'
          }
          item.severity = maxSeverity
        }
      }
    })
    
    // Sort combined by created_at desc
    combined.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
    recentScans.value = combined.slice(0, 5)
    
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadDashboardData()
})
</script>

<template>
  <div class="dashboard">
    <Header>
      <template #title>Dashboard</template>
    </Header>

    <div class="dashboard-content">

      <!-- Quick Actions -->
      <div class="quick-actions">
        <h2 class="section-title">Quick Actions</h2>
        <div class="action-buttons">
          <button class="action-btn primary" @click="router.push('/llmpentest')">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10" />
              <path d="M12 6v6l4 2" />
            </svg>
            Run LLM Pentest
          </button>
          <button class="action-btn secondary" @click="router.push('/reports')">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" />
              <path d="M14 2v6h6M16 13H8M16 17H8" />
            </svg>
            View Reports
          </button>
        </div>
      </div>

      <!-- Main Grid -->
      <div class="main-grid">
        <!-- Recent Scans -->
        <div class="card scans-card">
          <div class="card-header">
            <h3 class="card-title">Recent Pentests</h3>
            <button class="view-all" @click="router.push('/llmpentest')">View All</button>
          </div>
          <div class="scans-table">
            <div class="table-header">
              <span>Name</span>
              <span>Type</span>
              <span>Status</span>
              <span>Findings</span>
              <span>Time</span>
            </div>
            <div v-for="scan in recentScans" :key="scan.id" class="table-row" @click="router.push(`/llmpentest/${scan.id}`)">
              <span class="scan-name">{{ scan.name }}</span>
              <span class="scan-type">{{ scan.type }}</span>
              <span :class="['scan-status', scan.status]">
                <span class="status-dot"></span>
                {{ scan.status }}
              </span>
              <span v-if="scan.findings !== null" :class="['scan-findings', getSeverityClass(scan.severity)]">
                {{ scan.findings }}
              </span>
              <span v-else class="scan-findings">-</span>
              <span class="scan-date">{{ scan.date }}</span>
            </div>
          </div>
        </div>

        <!-- Threat Distribution -->
        <div class="card threats-card">
          <div class="card-header">
            <h3 class="card-title">Threat Distribution</h3>
          </div>
          <div class="threats-list">
            <div v-for="threat in threatDistribution" :key="threat.label" class="threat-item">
              <div class="threat-info">
                <span class="threat-dot" :style="{ background: threat.color }"></span>
                <span class="threat-label">{{ threat.label }}</span>
              </div>
              <div class="threat-bar-container">
                <div class="threat-bar" :style="{ width: (totalFindings ? (threat.count / totalFindings * 100) : 0) + '%', background: threat.color }"></div>
              </div>
              <span class="threat-count">{{ threat.count }}</span>
            </div>
          </div>
          <div class="threat-summary">
            <div class="summary-item">
              <span class="summary-value">{{ totalFindings }}</span>
              <span class="summary-label">Total Findings</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dashboard {
  min-height: 100%;
}

.dashboard-content {
  padding-top: 24px;
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 32px;
}

.stat-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 0;
  padding: 24px;
  position: relative;
  overflow: hidden;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.4);
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, var(--accent), var(--accent-pink));
  opacity: 0.5;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
  border-color: var(--accent) !important;
  box-shadow: 0 0 15px var(--accent-glow) !important;
}

.stat-card:hover::before {
  opacity: 1;
  height: 4px;
}

.stat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.stat-label {
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  color: var(--text-muted);
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.stat-change {
  display: flex;
  align-items: center;
  gap: 4px;
  font-family: 'Share Tech Mono', monospace;
  font-size: 13px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 0;
}

.stat-change.up {
  color: var(--success);
  background: rgba(57, 255, 20, 0.1);
  border: 1px solid rgba(57, 255, 20, 0.2);
}

.stat-change.down {
  color: var(--danger);
  background: rgba(255, 7, 58, 0.1);
  border: 1px solid rgba(255, 7, 58, 0.2);
}

.stat-change svg {
  width: 14px;
  height: 14px;
}

.stat-value {
  font-family: 'Share Tech Mono', monospace;
  font-size: 38px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
  text-shadow: 0 0 5px rgba(241, 245, 249, 0.2);
}

.stat-icon {
  position: absolute;
  right: 20px;
  bottom: 20px;
  width: 44px;
  height: 44px;
  background: rgba(0, 255, 204, 0.05);
  border: 1px solid var(--border-color);
  border-radius: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-card:hover .stat-icon {
  border-color: var(--accent);
  background: rgba(0, 255, 204, 0.1);
  box-shadow: 0 0 8px var(--accent-glow);
}

.stat-icon svg {
  width: 20px;
  height: 20px;
  color: var(--accent);
}

/* Quick Actions */
.quick-actions {
  margin-bottom: 32px;
}

.section-title {
  font-family: 'Orbitron', sans-serif;
  font-size: 14px;
  font-weight: 700;
  color: var(--accent);
  margin-bottom: 16px;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  text-shadow: 0 0 5px var(--accent-glow);
}

.action-buttons {
  display: flex;
  gap: 12px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  border-radius: 0;
  font-family: 'Orbitron', sans-serif;
  font-size: 13px;
  font-weight: 700;
  border: 1px solid transparent;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.action-btn.primary {
  background: var(--accent-pink);
  color: white;
  border-color: var(--accent-pink);
  box-shadow: 0 0 10px var(--accent-pink-glow);
}

.action-btn.primary:hover {
  transform: translateY(-2px);
  background: #ff1a8c;
  box-shadow: 0 0 20px rgba(255, 0, 127, 0.6);
}

.action-btn.secondary {
  background: var(--bg-card);
  color: var(--accent);
  border-color: var(--accent);
  box-shadow: 0 0 5px var(--accent-glow);
}

.action-btn.secondary:hover {
  background: rgba(0, 255, 204, 0.08);
  box-shadow: 0 0 15px var(--accent-glow);
  transform: translateY(-2px);
}

.action-btn svg {
  width: 16px;
  height: 16px;
}

/* Main Grid */
.main-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 24px;
}

.card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 0;
  overflow: hidden;
  box-shadow: 0 0 15px rgba(0, 0, 0, 0.4);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 2px solid var(--border-color);
  background: rgba(0, 0, 0, 0.2);
}

.card-title {
  font-family: 'Orbitron', sans-serif;
  font-size: 14px;
  font-weight: 700;
  letter-spacing: 1px;
  text-transform: uppercase;
}

.view-all {
  background: transparent;
  border: 1px solid var(--accent);
  color: var(--accent);
  padding: 4px 12px;
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  transition: all 0.2s ease;
}

.view-all:hover {
  background: rgba(0, 255, 204, 0.1);
  box-shadow: 0 0 10px var(--accent-glow);
  text-decoration: none;
}

/* Scans Table */
.scans-table {
  padding: 0 24px 24px;
}

.table-header {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr 1fr;
  padding: 16px;
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  font-weight: 700;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 1px;
  border-bottom: 1px solid var(--border-color);
}

.table-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr 1fr;
  padding: 16px;
  border-radius: 0;
  align-items: center;
  border-bottom: 1px solid var(--border-color);
  transition: all 0.2s ease;
  cursor: pointer;
}

.table-row:hover {
  background: rgba(0, 255, 204, 0.04);
  box-shadow: inset 4px 0 0 var(--accent);
}

.scan-name {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 15px;
}

.scan-type {
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--text-muted);
}

.scan-status {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 0;
}

.scan-status.completed {
  color: var(--success);
}
.scan-status.completed .status-dot {
  background: var(--success);
  box-shadow: 0 0 8px var(--success);
}

.scan-status.running {
  color: var(--accent);
}
.scan-status.running .status-dot {
  background: var(--accent);
  box-shadow: 0 0 8px var(--accent);
  animation: pulse 1.5s infinite;
}

.scan-status.failed {
  color: var(--danger);
}
.scan-status.failed .status-dot {
  background: var(--danger);
  box-shadow: 0 0 8px var(--danger);
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}

.scan-findings {
  font-family: 'Share Tech Mono', monospace;
  font-size: 13px;
  font-weight: 700;
  padding: 2px 10px;
  border-radius: 0;
  text-align: center;
  width: fit-content;
  border: 1px solid transparent;
}

.scan-findings.severity-critical {
  background: rgba(255, 7, 58, 0.1);
  color: var(--danger);
  border-color: rgba(255, 7, 58, 0.3);
  box-shadow: 0 0 8px rgba(255, 7, 58, 0.2);
}

.scan-findings.severity-high {
  background: rgba(255, 159, 0, 0.1);
  color: var(--warning);
  border-color: rgba(255, 159, 0, 0.3);
}

.scan-findings.severity-medium {
  background: rgba(0, 229, 255, 0.1);
  color: var(--info);
  border-color: rgba(0, 229, 255, 0.3);
}

.scan-findings.severity-low {
  background: rgba(98, 114, 164, 0.1);
  color: #a0aec0;
  border-color: rgba(98, 114, 164, 0.3);
}

.scan-date {
  font-family: 'Share Tech Mono', monospace;
  font-size: 13px;
  color: var(--text-muted);
}

/* Threats Card */
.threats-list {
  padding: 20px 24px;
}

.threat-item {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.threat-item:last-child {
  margin-bottom: 0;
}

.threat-info {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100px;
}

.threat-dot {
  width: 8px;
  height: 8px;
  border-radius: 0;
}

.threat-label {
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--text-primary);
}

.threat-bar-container {
  flex: 1;
  height: 6px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid var(--border-color);
  border-radius: 0;
  overflow: hidden;
}

.threat-bar {
  height: 100%;
  border-radius: 0;
  transition: width 0.5s ease;
  box-shadow: 0 0 8px rgba(0, 255, 204, 0.2);
}

.threat-count {
  width: 32px;
  text-align: right;
  font-family: 'Share Tech Mono', monospace;
  font-size: 14px;
  font-weight: 700;
  color: var(--accent);
}

.threat-summary {
  display: flex;
  border-top: 2px solid var(--border-color);
  margin-top: 20px;
  padding-top: 20px;
  background: rgba(0, 0, 0, 0.1);
}

.summary-item {
  flex: 1;
  text-align: center;
  border-right: 1px solid var(--border-color);
}

.summary-item:last-child {
  border-right: none;
}

.summary-value {
  display: block;
  font-family: 'Share Tech Mono', monospace;
  font-size: 26px;
  font-weight: 700;
  color: var(--text-primary);
  text-shadow: 0 0 5px rgba(241, 245, 249, 0.2);
}

.summary-label {
  font-family: 'Orbitron', sans-serif;
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--text-muted);
}

/* Responsive */
@media (max-width: 1400px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 1024px) {
  .main-grid {
    grid-template-columns: 1fr;
  }
}
</style>