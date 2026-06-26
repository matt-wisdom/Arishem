<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { AlertRule } from '@/stores/types'

const alertRules = ref<AlertRule[]>([])
const loading = ref(true)
const showNewModal = ref(false)
const newRule = ref({ severity_threshold: 'high', channel: 'email', channel_config: {} as Record<string, string> })

const fetchAlerts = async () => {
  try {
    const res = await fetch('/api/alerts', {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
    if (!res.ok) throw new Error('Failed to fetch alerts')
    alertRules.value = await res.json()
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const createRule = async () => {
  try {
    const res = await fetch('/api/alerts', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify(newRule.value)
    })
    if (!res.ok) throw new Error('Failed to create rule')
    showNewModal.value = false
    await fetchAlerts()
  } catch (e) {
    console.error(e)
  }
}

const deleteRule = async (id: string) => {
  if (!confirm('Are you sure?')) return
  try {
    const res = await fetch(`/api/alerts/${id}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
    if (res.ok) await fetchAlerts()
  } catch (e) {
    console.error(e)
  }
}

const testAlert = async (id: string) => {
  try {
    const res = await fetch(`/api/alerts/test/${id}`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
    if (res.ok) alert('Test alert sent!')
  } catch (e) {
    console.error(e)
  }
}

onMounted(fetchAlerts)
</script>

<template>
  <div class="alerts-page">
    <header class="header">
      <h1>Alert Rules</h1>
      <button class="primary-btn" @click="showNewModal = true">New Rule</button>
    </header>

    <div v-if="loading" class="loading">Loading...</div>
    <table v-else class="alerts-table">
      <thead>
        <tr>
          <th>Severity Threshold</th>
          <th>Channel</th>
          <th>Active</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="rule in alertRules" :key="rule.id">
          <td>{{ rule.severity_threshold }}</td>
          <td>{{ rule.channel }}</td>
          <td>{{ rule.active ? 'Yes' : 'No' }}</td>
          <td>
            <button @click="testAlert(rule.id)">Test</button>
            <button @click="deleteRule(rule.id)" class="danger-btn">Delete</button>
          </td>
        </tr>
      </tbody>
    </table>

    <div v-if="showNewModal" class="modal-overlay" @click.self="showNewModal = false">
      <div class="modal">
        <h2>New Alert Rule</h2>
        <form @submit.prevent="createRule">
          <div class="form-group">
            <label>Severity Threshold</label>
            <select v-model="newRule.severity_threshold">
              <option value="critical">Critical</option>
              <option value="high">High</option>
              <option value="medium">Medium</option>
              <option value="low">Low</option>
            </select>
          </div>
          <div class="form-group">
            <label>Channel</label>
            <select v-model="newRule.channel">
              <option value="email">Email</option>
              <option value="slack">Slack</option>
              <option value="webhook">Webhook</option>
            </select>
          </div>
          <div v-if="newRule.channel === 'email'" class="form-group">
            <label>Email</label>
            <input v-model="newRule.channel_config.address" placeholder="email@example.com" />
          </div>
          <div v-if="newRule.channel === 'slack'" class="form-group">
            <label>Slack Webhook URL</label>
            <input v-model="newRule.channel_config.url" placeholder="https://hooks.slack.com/..." />
          </div>
          <div v-if="newRule.channel === 'webhook'" class="form-group">
            <label>Webhook URL</label>
            <input v-model="newRule.channel_config.url" placeholder="https://example.com/webhook" />
          </div>
          <div class="modal-actions">
            <button type="button" @click="showNewModal = false">Cancel</button>
            <button type="submit" class="primary-btn">Create</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.alerts-page { padding: 2rem; max-width: 1000px; margin: 0 auto; }
.header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.primary-btn { padding: 0.5rem 1rem; background: #2563eb; color: white; border: none; border-radius: 6px; cursor: pointer; }
.alerts-table { width: 100%; background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
.alerts-table th, .alerts-table td { padding: 1rem; text-align: left; border-bottom: 1px solid #eee; }
.danger-btn { background: #dc2626; color: white; border: none; padding: 0.25rem 0.5rem; border-radius: 4px; cursor: pointer; margin-left: 0.5rem; }
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; }
.modal { background: white; padding: 2rem; border-radius: 8px; width: 400px; }
.form-group { margin-bottom: 1rem; }
.form-group label { display: block; margin-bottom: 0.5rem; font-size: 0.875rem; }
.form-group input, .form-group select { width: 100%; padding: 0.5rem; border: 1px solid #ddd; border-radius: 4px; }
.modal-actions { display: flex; justify-content: flex-end; gap: 0.5rem; margin-top: 1.5rem; }
</style>