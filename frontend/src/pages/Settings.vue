<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useOrganization } from '@clerk/vue'

const activeTab = ref('general')
const members = ref([
  { id: '1', name: 'John Doe', email: 'john@company.com', role: 'Admin', avatar: 'JD' },
  { id: '2', name: 'Jane Smith', email: 'jane@company.com', role: 'Engineer', avatar: 'JS' },
  { id: '3', name: 'Bob Wilson', email: 'bob@company.com', role: 'Viewer', avatar: 'BW' },
])

const { organization } = useOrganization()
const orgName = computed(() => organization.value?.name || 'Personal Account')

const token = ref('')

onMounted(() => {
  token.value = localStorage.getItem('token') || ''
})

const saveToken = () => {
  localStorage.setItem('token', token.value)
  alert('Authentication JWT Token saved successfully!')
}
</script>

<template>
  <div class="settings-page">
    <Header><template #title>Settings</template></Header>
    <div class="page-content">
      <div class="settings-layout">
        <aside class="settings-nav">
          <button :class="['nav-btn', { active: activeTab === 'general' }]" @click="activeTab = 'general'">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06a1.65 1.65 0 00.33-1.82 1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06a1.65 1.65 0 001.82.33H9a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 001.51 1H21a2 2 0 012 2 2 2 0 01-2 2h-.09a1.65 1.65 0 00-1.51 1z"/></svg>
            General
          </button>
          <button :class="['nav-btn', { active: activeTab === 'members' }]" @click="activeTab = 'members'">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 00-3-3.87M16 3.13a4 4 0 010 7.75"/></svg>
            Members
          </button>
          <button :class="['nav-btn', { active: activeTab === 'billing' }]" @click="activeTab = 'billing'">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="1" y="4" width="22" height="16" rx="2" /><path d="M1 10h22"/></svg>
            Billing
          </button>
          <button :class="['nav-btn', { active: activeTab === 'api' }]" @click="activeTab = 'api'">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 11-7.778 7.778 5.5 5.5 0 017.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/></svg>
            API Keys
          </button>
        </aside>
        
        <main class="settings-content">
          <div v-if="activeTab === 'general'" class="settings-section">
            <h2>General Settings</h2>
            <div class="form-card">
              <div class="form-group">
                <label>Organization Name</label>
                <input type="text" :value="orgName" disabled />
              </div>
              <div class="form-group">
                <label>Timezone</label>
                <select><option>UTC</option><option>America/New_York</option><option>Europe/London</option></select>
              </div>
              <button class="save-btn">Save Changes</button>
            </div>
          </div>
          
          <div v-if="activeTab === 'members'" class="settings-section">
            <div class="section-header">
              <h2>Team Members</h2>
              <button class="btn-primary">Invite Member</button>
            </div>
            <div class="members-list">
              <div v-for="member in members" :key="member.id" class="member-row">
                <div class="member-avatar">{{ member.avatar }}</div>
                <div class="member-info">
                  <h4>{{ member.name }}</h4>
                  <p>{{ member.email }}</p>
                </div>
                <span class="member-role">{{ member.role }}</span>
                <button class="action-btn">Edit</button>
              </div>
            </div>
          </div>
          
          <div v-if="activeTab === 'billing'" class="settings-section">
            <h2>Billing</h2>
            <div class="form-card">
              <p class="placeholder-text">Billing integration coming soon. Contact sales@arishem.com for enterprise plans.</p>
            </div>
          </div>
          
          <div v-if="activeTab === 'api'" class="settings-section">
            <div class="section-header">
              <h2>Authentication & Credentials</h2>
            </div>
            
            <div class="form-card" style="margin-bottom: 24px;">
              <h3 style="margin-bottom: 8px;">Clerk JWT Token</h3>
              <p class="placeholder-text" style="margin-bottom: 16px;">Paste your Clerk JWT token below. It will be sent in the <code>Authorization: Bearer</code> header for all API requests to the backend.</p>
              <div class="form-group">
                <input v-model="token" type="text" placeholder="eyJhbGciOiJSUzI1NiIs..." />
              </div>
              <button class="save-btn" @click="saveToken">Save Token</button>
            </div>

            <div class="section-header" style="margin-top: 32px;">
              <h2>API Keys</h2>
              <button class="btn-primary">Generate Key</button>
            </div>
            <div class="api-keys-list">
              <div class="api-key-card">
                <div class="key-info">
                  <h4>Production Key</h4>
                  <p>arishem_live_****************************</p>
                </div>
                <span class="key-created">Created 30 days ago</span>
              </div>
            </div>
          </div>
        </main>
      </div>
    </div>
  </div>
</template>

<style scoped>
.settings-page { min-height: 100%; }
.page-content { padding-top: 24px; }
.settings-layout { display: flex; gap: 32px; }
.settings-nav { width: 220px; flex-shrink: 0; }
.nav-btn { display: flex; align-items: center; gap: 12px; width: 100%; padding: 14px 16px; background: transparent; border: none; border-radius: 10px; color: var(--text-secondary); font-size: 14px; font-weight: 500; text-align: left; margin-bottom: 4px; }
.nav-btn:hover { background: var(--bg-card); color: var(--text-primary); }
.nav-btn.active { background: var(--bg-card); color: var(--accent); border: 1px solid rgba(59, 130, 246, 0.3); }
.nav-btn svg { width: 18px; height: 18px; }
.settings-content { flex: 1; }
.settings-section h2 { font-size: 20px; font-weight: 600; margin-bottom: 20px; }
.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.section-header h2 { margin-bottom: 0; }
.form-card { background: var(--bg-card); border: 1px solid var(--border-color); border-radius: 14px; padding: 24px; }
.form-group { margin-bottom: 20px; }
.form-group label { display: block; font-size: 14px; font-weight: 500; color: var(--text-secondary); margin-bottom: 8px; }
.form-group input, .form-group select { width: 100%; background: var(--bg-secondary); border: 1px solid var(--border-color); border-radius: 10px; padding: 12px 16px; color: var(--text-primary); font-size: 14px; }
.btn-primary { background: linear-gradient(135deg, var(--accent), #2563eb); color: white; border: none; border-radius: 10px; padding: 10px 20px; font-size: 14px; font-weight: 600; }
.save-btn { background: var(--accent); color: white; border: none; border-radius: 10px; padding: 12px 24px; font-size: 14px; font-weight: 600; }
.members-list { display: flex; flex-direction: column; gap: 12px; }
.member-row { display: flex; align-items: center; gap: 16px; background: var(--bg-card); border: 1px solid var(--border-color); border-radius: 12px; padding: 16px 20px; }
.member-avatar { width: 44px; height: 44px; background: linear-gradient(135deg, #8b5cf6, #6366f1); border-radius: 12px; display: flex; align-items: center; justify-content: center; font-weight: 600; }
.member-info { flex: 1; }
.member-info h4 { font-size: 15px; font-weight: 600; margin-bottom: 2px; }
.member-info p { font-size: 13px; color: var(--text-muted); }
.member-role { padding: 6px 12px; background: var(--bg-secondary); border-radius: 6px; font-size: 13px; font-weight: 500; }
.action-btn { padding: 8px 16px; background: transparent; border: 1px solid var(--border-color); border-radius: 8px; color: var(--text-secondary); font-size: 13px; }
.api-keys-list { display: flex; flex-direction: column; gap: 12px; }
.api-key-card { display: flex; justify-content: space-between; align-items: center; background: var(--bg-card); border: 1px solid var(--border-color); border-radius: 12px; padding: 20px; }
.key-info h4 { font-size: 15px; font-weight: 600; margin-bottom: 4px; }
.key-info p { font-size: 13px; color: var(--text-muted); font-family: monospace; }
.key-created { font-size: 13px; color: var(--text-muted); }
.placeholder-text { color: var(--text-muted); font-size: 14px; }
</style>