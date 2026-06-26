<script setup lang="ts">
import { ref, onMounted } from 'vue'

const integrations = ref<{ provider: string; connected: boolean }[]>([])
const loading = ref(true)
const githubToken = ref('')

const fetchIntegrations = async () => {
  try {
    const res = await fetch('/api/integrations/github', {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
    integrations.value = [{ provider: 'github', connected: res.ok }]
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const connectGitHub = async () => {
  if (!githubToken.value) return
  try {
    const res = await fetch('/api/integrations/github', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({ token: githubToken.value })
    })
    if (res.ok) await fetchIntegrations()
  } catch (e) {
    console.error(e)
  }
}

const disconnectGitHub = async () => {
  try {
    const res = await fetch('/api/integrations/github', {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
    if (res.ok) await fetchIntegrations()
  } catch (e) {
    console.error(e)
  }
}

onMounted(fetchIntegrations)
</script>

<template>
  <div class="integrations-page">
    <header class="header">
      <h1>Integrations</h1>
    </header>

    <div class="integration-card">
      <div class="integration-header">
        <h2>GitHub</h2>
        <span v-if="integrations.find(i => i.provider === 'github')?.connected" class="connected">Connected</span>
        <span v-else class="not-connected">Not Connected</span>
      </div>
      <p>Connect your GitHub account to enable repository scanning and GitHub Actions integration.</p>
      <div v-if="!integrations.find(i => i.provider === 'github')?.connected" class="connect-form">
        <input v-model="githubToken" type="password" placeholder="GitHub Personal Access Token" />
        <button @click="connectGitHub">Connect</button>
      </div>
      <button v-else @click="disconnectGitHub" class="disconnect-btn">Disconnect</button>
    </div>
  </div>
</template>

<style scoped>
.integrations-page { padding: 2rem; max-width: 800px; margin: 0 auto; }
.header { margin-bottom: 2rem; }
.integration-card { background: white; padding: 1.5rem; border-radius: 8px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
.integration-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem; }
.connected { color: #059669; }
.not-connected { color: #dc2626; }
.connect-form { display: flex; gap: 0.5rem; margin-top: 1rem; }
.connect-form input { flex: 1; padding: 0.5rem; border: 1px solid #ddd; border-radius: 4px; }
.connect-form button { padding: 0.5rem 1rem; background: #2563eb; color: white; border: none; border-radius: 6px; cursor: pointer; }
.disconnect-btn { margin-top: 1rem; padding: 0.5rem 1rem; background: #dc2626; color: white; border: none; border-radius: 6px; cursor: pointer; }
</style>