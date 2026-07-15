<script setup lang="ts">
import { computed } from 'vue'
import { useUser, useClerk, useOrganization } from '@clerk/vue'
import { useRoute } from 'vue-router'

const { user } = useUser()
const clerk = useClerk()
const { organization } = useOrganization()

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
  }
  return nameMap[name] || 'Dashboard'
})

const handleSignOut = () => {
  if (confirm('Are you sure you want to sign out?')) {
    clerk.value?.signOut()
  }
}
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

      <button class="notification-btn">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9" />
          <path d="M13.73 21a2 2 0 01-3.46 0" />
        </svg>
        <span class="notification-badge">3</span>
      </button>

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
</style>