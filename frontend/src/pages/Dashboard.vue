<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const stats = ref({
  scansThisMonth: 0,
  findingsThisMonth: 0,
  activeAlerts: 0,
  integrations: 0
})

onMounted(() => {
  // TODO: Fetch stats from API
})
</script>

<template>
  <div class="dashboard">
    <header class="header">
      <h1>Dashboard</h1>
      <p>Welcome to Arishem</p>
    </header>

    <div class="stats-grid">
      <div class="stat-card">
        <h3>Scans This Month</h3>
        <p class="stat-value">{{ stats.scansThisMonth }}</p>
      </div>
      <div class="stat-card">
        <h3>Findings This Month</h3>
        <p class="stat-value">{{ stats.findingsThisMonth }}</p>
      </div>
      <div class="stat-card">
        <h3>Active Alerts</h3>
        <p class="stat-value">{{ stats.activeAlerts }}</p>
      </div>
      <div class="stat-card">
        <h3>Integrations</h3>
        <p class="stat-value">{{ stats.integrations }}</p>
      </div>
    </div>

    <section class="quick-actions">
      <h2>Quick Actions</h2>
      <div class="actions">
        <button @click="router.push('/scans')">Run New Scan</button>
        <button @click="router.push('/llmpentest')">Run LLM Pentest</button>
        <button @click="router.push('/reports')">View Reports</button>
      </div>
    </section>
  </div>
</template>

<style scoped>
.dashboard {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.header {
  margin-bottom: 2rem;
}

.header h1 {
  font-size: 2rem;
  color: #1a1a1a;
}

.header p {
  color: #666;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: white;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.stat-card h3 {
  font-size: 0.875rem;
  color: #666;
  margin-bottom: 0.5rem;
}

.stat-value {
  font-size: 2rem;
  font-weight: 600;
  color: #1a1a1a;
}

.quick-actions h2 {
  margin-bottom: 1rem;
}

.actions {
  display: flex;
  gap: 1rem;
}

.actions button {
  padding: 0.75rem 1.5rem;
  background: #2563eb;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.875rem;
}

.actions button:hover {
  background: #1d4ed8;
}
</style>