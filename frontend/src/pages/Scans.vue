<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useScansStore } from '@/stores/scans'

const router = useRouter()
const scansStore = useScansStore()
const showNewScanModal = ref(false)
const newScan = ref({ type: 'code', target: '', branch: 'main' })

const createScan = async () => {
  const endpoint = newScan.value.type === 'code' ? '/api/scans/code' : '/api/scans/webapp'
  try {
    const res = await fetch(endpoint, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify(newScan.value)
    })
    if (!res.ok) throw new Error('Failed to create scan')
    showNewScanModal.value = false
    await scansStore.fetchScans()
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => scansStore.fetchScans())
</script>

<template>
  <div class="scans-page">
    <header class="header">
      <h1>Scans</h1>
      <button class="primary-btn" @click="showNewScanModal = true">New Scan</button>
    </header>

    <div v-if="scansStore.loading" class="loading">Loading...</div>
    <div v-else-if="scansStore.error" class="error">{{ scansStore.error }}</div>

    <table v-else class="scans-table">
      <thead>
        <tr>
          <th>Type</th>
          <th>Target</th>
          <th>Status</th>
          <th>Created</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="scan in scansStore.scans" :key="scan.id" @click="router.push(`/scans/${scan.id}`)">
          <td>{{ scan.type }}</td>
          <td>{{ scan.target }}</td>
          <td><span :class="['status', scan.status]">{{ scan.status }}</span></td>
          <td>{{ new Date(scan.created_at).toLocaleDateString() }}</td>
          <td><button @click.stop="router.push(`/scans/${scan.id}`)">View</button></td>
        </tr>
      </tbody>
    </table>

    <div v-if="showNewScanModal" class="modal-overlay" @click.self="showNewScanModal = false">
      <div class="modal">
        <h2>New Scan</h2>
        <form @submit.prevent="createScan">
          <div class="form-group">
            <label>Scan Type</label>
            <select v-model="newScan.type">
              <option value="code">Code Scan</option>
              <option value="webapp">Web App Scan</option>
            </select>
          </div>
          <div class="form-group">
            <label>Target</label>
            <input v-model="newScan.target" placeholder="Repo URL or target URL" required />
          </div>
          <div v-if="newScan.type === 'code'" class="form-group">
            <label>Branch</label>
            <input v-model="newScan.branch" placeholder="main" />
          </div>
          <div class="modal-actions">
            <button type="button" @click="showNewScanModal = false">Cancel</button>
            <button type="submit" class="primary-btn">Create Scan</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.scans-page { padding: 2rem; max-width: 1200px; margin: 0 auto; }
.header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.primary-btn { padding: 0.5rem 1rem; background: #2563eb; color: white; border: none; border-radius: 6px; cursor: pointer; }
.scans-table { width: 100%; background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
.scans-table th, .scans-table td { padding: 1rem; text-align: left; border-bottom: 1px solid #eee; }
.scans-table tr:hover { background: #f9fafb; cursor: pointer; }
.status { padding: 0.25rem 0.5rem; border-radius: 4px; font-size: 0.75rem; }
.status.queued { background: #fef3c7; color: #92400e; }
.status.running { background: #dbeafe; color: #1e40af; }
.status.completed { background: #d1fae5; color: #065f46; }
.status.failed { background: #fee2e2; color: #991b1b; }
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; }
.modal { background: white; padding: 2rem; border-radius: 8px; width: 400px; }
.form-group { margin-bottom: 1rem; }
.form-group label { display: block; margin-bottom: 0.5rem; font-size: 0.875rem; }
.form-group input, .form-group select { width: 100%; padding: 0.5rem; border: 1px solid #ddd; border-radius: 4px; }
.modal-actions { display: flex; justify-content: flex-end; gap: 0.5rem; margin-top: 1.5rem; }
</style>