<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuth } from '@clerk/vue'
import { apiUrl, apiFetch } from '@/utils/api'

const { getToken } = useAuth()
const alerts = ref<any[]>([])
const loading = ref(true)
const error = ref<string | null>(null)

const showNewModal = ref(false)
const newRule = ref({
  severity_threshold: 'high' as 'critical' | 'high' | 'medium' | 'low',
  channel: 'webhook' as 'email' | 'slack' | 'webhook',
  target: ''
})

const getHeaders = async () => {
  const tokenFn = typeof getToken.value === 'function' ? getToken.value : getToken
  const token = (await (tokenFn as any)()) || localStorage.getItem('token') || ''
  return { 'Authorization': `Bearer ${token}` }
}

const getTargetField = (rule: any) => {
  if (rule.channel_config) {
    return rule.channel_config.url || rule.channel_config.email || rule.channel_config.webhook_url || ''
  }
  return ''
}

const loadAlerts = async () => {
  loading.value = true
  error.value = null
  try {
    const headers = await getHeaders()
    const res = await apiFetch(apiUrl('/alerts'), { headers })
    if (!res.ok) throw new Error('Failed to fetch alert rules')
    alerts.value = await res.json()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Unknown error'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadAlerts()
})

const handleCreateAlert = async () => {
  if (!newRule.value.target) return
  try {
    // Determine target config key
    const configKey = newRule.value.channel === 'email' ? 'email' : 'url'
    const authHeaders = await getHeaders()
    const res = await apiFetch(apiUrl('/alerts'), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...authHeaders
      },
      body: JSON.stringify({
        severity_threshold: newRule.value.severity_threshold,
        channel: newRule.value.channel,
        channel_config: { [configKey]: newRule.value.target }
      })
    })
    if (!res.ok) throw new Error('Failed to create alert rule')
    showNewModal.value = false
    newRule.value.target = ''
    loadAlerts()
  } catch (e) {
    alert(e instanceof Error ? e.message : 'Unknown error')
  }
}

const handleToggleActive = async (rule: any) => {
  try {
    // Note: backend UpdateAlert expects CreateAlertRequest in body
    const authHeaders = await getHeaders()
    const res = await apiFetch(apiUrl(`/alerts/${rule.id}`), {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        ...authHeaders
      },
      body: JSON.stringify({
        severity_threshold: rule.severity_threshold,
        channel: rule.channel,
        channel_config: rule.channel_config,
        active: !rule.active // toggle
      })
    })
    if (!res.ok) throw new Error('Failed to update alert rule')
    loadAlerts()
  } catch (e) {
    alert(e instanceof Error ? e.message : 'Unknown error')
  }
}

const handleDeleteAlert = async (id: string) => {
  if (!confirm('Are you sure you want to delete this alert rule?')) return
  try {
    const authHeaders = await getHeaders()
    const res = await apiFetch(apiUrl(`/alerts/${id}`), {
      method: 'DELETE',
      headers: authHeaders
    })
    if (!res.ok) throw new Error('Failed to delete alert rule')
    loadAlerts()
  } catch (e) {
    alert(e instanceof Error ? e.message : 'Unknown error')
  }
}

const handleTestAlert = async (id: string) => {
  try {
    const authHeaders = await getHeaders()
    const res = await apiFetch(apiUrl(`/alerts/test/${id}`), {
      method: 'POST',
      headers: authHeaders
    })
    if (!res.ok) throw new Error('Failed to send test alert')
    alert('Test alert dispatched successfully!')
  } catch (e) {
    alert(e instanceof Error ? e.message : 'Unknown error')
  }
}
</script>

<template>
  <div class="alerts-page">
    <Header><template #title>Alerts</template></Header>
    <div class="page-content">
      <div class="toolbar">
        <button class="btn-primary" @click="showNewModal = true">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 5v14M5 12h14" /></svg>
          New Alert Rule
        </button>
      </div>

      <div v-if="loading && alerts.length === 0" class="loading-placeholder">Loading alert rules...</div>
      <div v-else-if="alerts.length === 0" class="empty-placeholder">No alert rules configured.</div>
      <div v-else class="alerts-list">
        <div v-for="alert in alerts" :key="alert.id" class="alert-card">
          <div class="alert-icon" :class="alert.severity_threshold">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9M13.73 21a2 2 0 01-3.46 0" />
            </svg>
          </div>
          <div class="alert-info">
            <h3>Alert for {{ alert.severity_threshold.toUpperCase() }}+</h3>
            <p>{{ alert.channel.toUpperCase() }} → {{ getTargetField(alert) }}</p>
          </div>
          <div class="alert-severity">
            <span :class="['severity-badge', alert.severity_threshold]">{{ alert.severity_threshold }}</span>
          </div>
          <div class="alert-toggle">
            <input type="checkbox" :checked="alert.active" @change="handleToggleActive(alert)" />
          </div>
          <button @click="handleTestAlert(alert.id)" class="action-btn">Test</button>
          <button @click="handleDeleteAlert(alert.id)" class="action-btn" style="border-color: var(--danger); color: var(--danger);">Delete</button>
        </div>
      </div>

      <!-- New Alert Modal -->
      <div v-if="showNewModal" class="modal-overlay" @click.self="showNewModal = false">
        <div class="modal">
          <div class="modal-header">
            <h2>Create Alert Rule</h2>
            <button class="close-btn" @click="showNewModal = false">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 6L6 18M6 6l12 12" />
              </svg>
            </button>
          </div>
          <form @submit.prevent="handleCreateAlert" class="modal-form">
            <div class="form-group">
              <label>Severity Threshold</label>
              <select v-model="newRule.severity_threshold" class="filter-select" style="width: 100%;">
                <option value="critical">Critical</option>
                <option value="high">High</option>
                <option value="medium">Medium</option>
                <option value="low">Low</option>
              </select>
            </div>
            <div class="form-group">
              <label>Alert Channel</label>
              <select v-model="newRule.channel" class="filter-select" style="width: 100%;">
                <option value="email">Email</option>
                <option value="slack">Slack</option>
                <option value="webhook">Webhook</option>
              </select>
            </div>
            <div class="form-group">
              <label>{{ newRule.channel === 'email' ? 'Email Address' : 'Webhook URL' }}</label>
              <input v-model="newRule.target" type="text" placeholder="Enter channel destination..." required />
            </div>
            <div class="modal-actions">
              <button type="button" class="btn-secondary" @click="showNewModal = false">Cancel</button>
              <button type="submit" class="btn-primary">Create Rule</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.alerts-page { min-height: 100%; }
.page-content { padding-top: 24px; }
.toolbar { display: flex; justify-content: flex-end; margin-bottom: 24px; }
.btn-primary { display: flex; align-items: center; gap: 8px; background: linear-gradient(135deg, #f59e0b, #d97706); color: white; border: none; border-radius: 10px; padding: 10px 20px; font-size: 14px; font-weight: 600; }
.btn-primary svg { width: 18px; height: 18px; }
.alerts-list { display: flex; flex-direction: column; gap: 12px; }
.alert-card { display: flex; align-items: center; gap: 20px; background: var(--bg-card); border: 1px solid var(--border-color); border-radius: 14px; padding: 20px 24px; }
.alert-icon { width: 48px; height: 48px; border-radius: 12px; display: flex; align-items: center; justify-content: center; }
.alert-icon.critical { background: rgba(239, 68, 68, 0.15); color: #ef4444; }
.alert-icon.high { background: rgba(249, 115, 22, 0.15); color: #f97316; }
.alert-icon.medium { background: rgba(245, 158, 11, 0.15); color: #f59e0b; }
.alert-icon.low { background: rgba(59, 130, 246, 0.15); color: #3b82f6; }
.alert-icon svg { width: 24px; height: 24px; }
.alert-info { flex: 1; }
.alert-info h3 { font-size: 16px; font-weight: 600; color: var(--text-primary); margin-bottom: 4px; }
.alert-info p { font-size: 13px; color: var(--text-muted); }
.severity-badge { padding: 6px 14px; border-radius: 8px; font-size: 13px; font-weight: 500; text-transform: capitalize; }
.severity-badge.critical { background: rgba(239, 68, 68, 0.15); color: #ef4444; }
.severity-badge.high { background: rgba(249, 115, 22, 0.15); color: #f97316; }
.severity-badge.medium { background: rgba(245, 158, 11, 0.15); color: #f59e0b; }
.severity-badge.low { background: rgba(59, 130, 246, 0.15); color: #3b82f6; }
.alert-toggle input { width: 44px; height: 24px; accent-color: var(--accent); }
.action-btn { padding: 8px 16px; background: var(--bg-secondary); border: 1px solid var(--border-color); border-radius: 8px; color: var(--text-primary); font-size: 13px; }
</style>