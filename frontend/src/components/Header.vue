<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useUser, useClerk, useOrganization, useAuth } from '@clerk/vue'
import { useRoute, useRouter } from 'vue-router'

const { user } = useUser()
const clerk = useClerk()
const { organization } = useOrganization()
const { getToken } = useAuth()
const router = useRouter()

const orgName = computed(() => organization.value?.name || 'Personal Account')

const route = useRoute()
const pageTitle = computed(() => {
  if (!route) return 'Dashboard'
  const name = route.name as string
  if (!name) return 'Dashboard'
  const nameMap: Record<string, string> = {
    'dashboard': 'Dashboard',
    'scans': 'Scans',
    'scan-detail': 'Scan Details',
    'llmpentest': 'LLM Pentests',
    'llmpentest-detail': 'LLM Pentest Details',
    'reports': 'Reports',
    'report-detail': 'Report Details',
    'integrations': 'Integrations',
    'alerts': 'Alerts',
    'settings': 'Settings',
    'docs': 'Docs & Help',
  }
  return nameMap[name] || 'Dashboard'
})

const handleSignOut = () => {
  if (confirm('Are you sure you want to sign out?')) {
    clerk.value?.signOut()
  }
}

const showNotifications = ref(false)
const notifications = ref<any[]>([])
const unreadCount = computed(() => notifications.value.filter(n => !n.read).length)

const getHeaders = async () => {
  const tokenFn = typeof getToken.value === 'function' ? getToken.value : getToken
  const token = (await (tokenFn as any)()) || localStorage.getItem('token') || ''
  return { 'Authorization': `Bearer ${token}` }
}

const loadNotifications = async () => {
  try {
    const headers = await getHeaders()
    const [runsRes, scansRes] = await Promise.all([
      fetch('/api/llmpentest', { headers }),
      fetch('/api/scans', { headers })
    ])
    
    let list: any[] = []
    if (runsRes.ok) {
      const runs = await runsRes.json()
      runs.forEach((r: any) => {
        list.push({
          id: r.id,
          type: 'llm',
          target: r.target_endpoint,
          status: r.status,
          time: new Date(r.created_at),
          message: `LLM Pentest run ${r.status} against ${r.target_endpoint}`,
          read: r.status !== 'queued' && r.status !== 'running'
        })
      })
    }
    if (scansRes.ok) {
      const scans = await scansRes.json()
      scans.forEach((s: any) => {
        list.push({
          id: s.id,
          type: 'scan',
          target: s.target,
          status: s.status,
          time: new Date(s.created_at),
          message: `Code scan ${s.status} against ${s.target}`,
          read: s.status !== 'queued' && s.status !== 'running'
        })
      })
    }
    
    list.sort((a, b) => b.time.getTime() - a.time.getTime())
    notifications.value = list.slice(0, 5)
  } catch (e) {
    console.error('Failed to load notifications:', e)
  }
}

const navigateToNotification = (n: any) => {
  showNotifications.value = false
  if (n.type === 'llm') {
    router.push(`/llmpentest/${n.id}`)
  } else {
    router.push('/reports')
  }
}

onMounted(() => {
  loadNotifications()
  setInterval(loadNotifications, 12000)
})
</script>

<template>
  <header class="header">
    <div class="header-left">
      <h1 class="page-title">
        {{ pageTitle }}
      </h1>
    </div>

    <div class="header-right">
      <button class="search-btn">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8" />
          <path d="M21 21l-4.35-4.35" />
        </svg>
      </button>

      <div class="notifications-wrapper">
        <button class="notification-btn" @click.stop="showNotifications = !showNotifications">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9" />
            <path d="M13.73 21a2 2 0 01-3.46 0" />
          </svg>
          <span v-if="unreadCount > 0" class="notification-badge">{{ unreadCount }}</span>
        </button>

        <div v-if="showNotifications" class="notifications-dropdown" @click.stop>
          <div class="dropdown-header">
            <h3>Recent Activity</h3>
            <button class="clear-btn" @click="showNotifications = false">Close</button>
          </div>
          <div class="dropdown-body">
            <div v-if="notifications.length === 0" class="empty-notifications">
              No recent scans or runs.
            </div>
            <div 
              v-for="n in notifications" 
              :key="n.id" 
              :class="['notification-item', n.status, { unread: !n.read }]"
              @click="navigateToNotification(n)"
            >
              <div class="notification-icon">
                <span v-if="n.status === 'completed'">🟢</span>
                <span v-else-if="n.status === 'failed'">🔴</span>
                <span v-else-if="n.status === 'cancelled'">🟡</span>
                <span v-else>🔵</span>
              </div>
              <div class="notification-content">
                <p class="message">{{ n.message }}</p>
                <span class="time">{{ n.time.toLocaleTimeString() }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="org-switcher">
        <div class="org-avatar">{{ orgName.charAt(0) }}</div>
        <span class="org-name">{{ orgName }}</span>
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M6 9l6 6 6-6" />
        </svg>
      </div>

      <div class="user-menu" @click="handleSignOut" style="cursor: pointer;" title="Click to Sign Out">
        <div class="user-avatar" v-if="user">
          {{ (user.firstName?.charAt(0) || '') + (user.lastName?.charAt(0) || '') || 'U' }}
        </div>
        <div class="user-avatar" v-else>JD</div>
        <span class="user-name" v-if="user">{{ user.fullName }}</span>
        <span class="user-name" v-else>John Doe</span>
      </div>
    </div>
  </header>
</template>

<style scoped>
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 32px;
  background: var(--bg-secondary);
  border-bottom: 2px solid var(--border-color);
  position: sticky;
  top: 0;
  z-index: 50;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.page-title {
  font-family: 'Orbitron', sans-serif;
  font-size: 20px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1.5px;
  color: var(--accent);
  text-shadow: 0 0 10px var(--accent-glow);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.search-btn,
.notification-btn {
  width: 40px;
  height: 40px;
  border-radius: 0;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  transition: all 0.2s ease;
}

.search-btn:hover,
.notification-btn:hover {
  background: var(--bg-card-hover);
  border-color: var(--accent);
  color: var(--text-primary);
  box-shadow: 0 0 8px var(--accent-glow);
}

.search-btn svg,
.notification-btn svg {
  width: 18px;
  height: 18px;
}

.notification-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  width: 16px;
  height: 16px;
  background: var(--danger);
  border-radius: 0;
  font-family: 'Share Tech Mono', monospace;
  font-size: 10px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 0 8px rgba(255, 7, 58, 0.6);
}

.org-switcher {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 14px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 0;
  cursor: pointer;
  transition: all 0.2s ease;
}

.org-switcher:hover {
  background: var(--bg-card-hover);
  border-color: var(--accent);
  box-shadow: 0 0 8px var(--accent-glow);
}

.org-avatar {
  width: 28px;
  height: 28px;
  background: linear-gradient(135deg, var(--accent), var(--bg-primary));
  border: 1px solid var(--accent);
  border-radius: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: 'Orbitron', sans-serif;
  font-weight: 700;
  font-size: 14px;
  color: #fff;
}

.org-name {
  font-family: 'Orbitron', sans-serif;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
  letter-spacing: 0.5px;
}

.org-switcher svg {
  width: 16px;
  height: 16px;
  color: var(--text-muted);
}

.user-menu {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 12px 6px 6px;
  background: transparent;
  border-radius: 0;
  cursor: pointer;
  transition: all 0.2s ease;
}

.user-menu:hover {
  background: rgba(255, 0, 127, 0.05);
  box-shadow: inset 3px 0 0 var(--accent-pink);
}

.user-avatar {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, var(--accent-pink), var(--bg-primary));
  border: 1px solid var(--accent-pink);
  border-radius: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: 'Orbitron', sans-serif;
  font-weight: 700;
  font-size: 13px;
  color: #fff;
}

.user-name {
  font-family: 'Orbitron', sans-serif;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
  letter-spacing: 0.5px;
}

.notifications-wrapper {
  position: relative;
  display: inline-block;
}

.notifications-dropdown {
  position: absolute;
  top: 50px;
  right: 0;
  width: 320px;
  background: #060613;
  border: 1px solid var(--border-color);
  border-top: 2px solid var(--accent);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.8), 0 0 10px var(--accent-glow);
  z-index: 999;
}

.dropdown-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  background: rgba(0, 0, 0, 0.2);
}

.dropdown-header h3 {
  font-family: 'Orbitron', sans-serif;
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  color: var(--accent);
  margin: 0;
}

.clear-btn {
  background: transparent;
  border: none;
  color: var(--text-muted);
  font-size: 11px;
  cursor: pointer;
}
.clear-btn:hover {
  color: var(--text-primary);
}

.dropdown-body {
  max-height: 300px;
  overflow-y: auto;
}

.empty-notifications {
  padding: 24px;
  text-align: center;
  color: var(--text-muted);
  font-size: 13px;
  font-style: italic;
}

.notification-item {
  display: flex;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  cursor: pointer;
  transition: background 0.2s ease;
}

.notification-item:hover {
  background: rgba(255, 255, 255, 0.03);
}

.notification-item.unread {
  background: rgba(139, 92, 246, 0.05);
}

.notification-icon {
  font-size: 14px;
  display: flex;
  align-items: center;
}

.notification-content {
  flex: 1;
}

.notification-content .message {
  font-size: 13px;
  color: var(--text-secondary);
  margin: 0 0 4px 0;
  line-height: 1.4;
  word-break: break-word;
}

.notification-content .time {
  font-size: 11px;
  color: var(--text-muted);
}
</style>