<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const navItems = [
  { name: 'Dashboard', path: '/', icon: 'dashboard' },
  { name: 'LLM Pentests', path: '/llmpentest', icon: 'brain' },
  { name: 'Reports', path: '/reports', icon: 'report' },
  { name: 'Integrations', path: '/integrations', icon: 'integration' },
  { name: 'Alerts', path: '/alerts', icon: 'alert' },
  { name: 'Settings', path: '/settings', icon: 'settings' },
  { name: 'Docs & Help', path: '/docs', icon: 'docs' },
]

const isActive = (path: string) => {
  if (path === '/') return route.path === '/'
  return route.path.startsWith(path)
}
</script>

<template>
  <aside class="sidebar">
    <div class="logo" @click="router.push('/')">
      <div class="logo-icon">
        <svg viewBox="0 0 32 32" fill="none">
          <path d="M16 2L4 8v8c0 7.732 5.12 14.936 12 17 6.88-2.064 12-9.268 12-17V8L16 2z" fill="url(#logoGrad)" />
          <path d="M16 8l-6 3v5c0 4.5 2.64 8.68 6 10 3.36-1.32 6-5.5 6-10v-5l-6-3z" fill="#0a0f1c" />
          <circle cx="16" cy="15" r="3" fill="#3b82f6" />
          <defs>
            <linearGradient id="logoGrad" x1="4" y1="2" x2="28" y2="27" gradientUnits="userSpaceOnUse">
              <stop stop-color="#3b82f6" />
              <stop offset="1" stop-color="#1d4ed8" />
            </linearGradient>
          </defs>
        </svg>
      </div>
      <span class="logo-text">Arishem</span>
    </div>

    <nav class="nav">
      <a
        v-for="item in navItems"
        :key="item.path"
        :href="item.path"
        class="nav-item"
        :class="{ active: isActive(item.path) }"
        @click.prevent="router.push(item.path)"
      >
        <span class="nav-icon">
          <svg v-if="item.icon === 'dashboard'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="3" width="7" height="7" rx="1" />
            <rect x="14" y="3" width="7" height="7" rx="1" />
            <rect x="3" y="14" width="7" height="7" rx="1" />
            <rect x="14" y="14" width="7" height="7" rx="1" />
          </svg>
          <svg v-else-if="item.icon === 'scan'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M9 3H5a2 2 0 00-2 2v4m6-6h10a2 2 0 012 2v4M9 3v18m0 0h10a2 2 0 002-2v-4M9 21H5a2 2 0 01-2-2v-4" />
            <circle cx="17" cy="7" r="3" />
          </svg>
          <svg v-else-if="item.icon === 'brain'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 4.5a2.5 2.5 0 00-4.96-.46 2.5 2.5 0 01-2.92-2.52 2.5 2.5 0 00.52-3.2A2.5 2.5 0 007.5 2.5 2.5 2.5 0 009 5" />
            <path d="M12 4.5a2.5 2.5 0 014.96.46 2.5 2.5 0 012.92-2.52 2.5 2.5 0 00-.52-3.2A2.5 2.5 0 0016.5 2.5 2.5 2.5 0 0015 5" />
            <circle cx="12" cy="12" r="6" />
            <path d="M12 6v6l4 2" />
          </svg>
          <svg v-else-if="item.icon === 'report'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" />
            <path d="M14 2v6h6M16 13H8M16 17H8M10 9H8" />
          </svg>
          <svg v-else-if="item.icon === 'integration'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="3" />
            <path d="M12 2v4M12 18v4M4.93 4.93l2.83 2.83M16.24 16.24l2.83 2.83M2 12h4M18 12h4M4.93 19.07l2.83-2.83M16.24 7.76l2.83-2.83" />
          </svg>
          <svg v-else-if="item.icon === 'alert'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M18 8A6 6 0 006 8c0 7-3 9-3 9h18s-3-2-3-9M13.73 21a2 2 0 01-3.46 0" />
          </svg>
          <svg v-else-if="item.icon === 'settings'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="3" />
            <path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06a1.65 1.65 0 00.33-1.82 1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06a1.65 1.65 0 001.82.33H9a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 001.51 1H21a2 2 0 012 2 2 2 0 01-2 2h-.09a1.65 1.65 0 00-1.51 1z" />
          </svg>
          <svg v-else-if="item.icon === 'docs'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
          </svg>
        </span>
        <span class="nav-text">{{ item.name }}</span>
      </a>
    </nav>

    <div class="sidebar-footer">
      <div class="version">v1.0.0</div>
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  width: 260px;
  background: #04040a;
  border-right: 2px solid var(--accent);
  display: flex;
  flex-direction: column;
  position: fixed;
  height: 100vh;
  z-index: 100;
  box-shadow: 0 0 15px rgba(0, 0, 255, 0.1);
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 24px;
  cursor: pointer;
  border-bottom: 2px solid var(--border-color);
}

.logo-icon {
  width: 40px;
  height: 40px;
  filter: drop-shadow(0 0 5px var(--accent-glow));
}

.logo-icon svg {
  width: 100%;
  height: 100%;
}

.logo-text {
  font-family: 'Orbitron', sans-serif;
  font-size: 22px;
  font-weight: 900;
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-pink) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: 1px;
  text-shadow: 0 0 10px var(--accent-glow);
}

.nav {
  flex: 1;
  padding: 20px 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.nav-item {
  font-family: 'Orbitron', sans-serif;
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 14px 24px;
  border-radius: 0;
  color: #a0aec0;
  text-decoration: none;
  font-weight: 600;
  letter-spacing: 1px;
  text-transform: uppercase;
  font-size: 13px;
  border-left: 3px solid transparent;
  transition: all 0.2s ease;
}

.nav-item:hover {
  background: rgba(0, 255, 204, 0.04);
  color: var(--accent);
  border-left-color: var(--accent);
}

.nav-item.active {
  background: rgba(0, 255, 204, 0.08);
  color: var(--accent);
  border-left-color: var(--accent);
  box-shadow: inset 5px 0 10px -5px var(--accent-glow);
  text-shadow: 0 0 5px var(--accent-glow);
}

.nav-icon {
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.nav-icon svg {
  width: 100%;
  height: 100%;
}

.nav-text {
  font-size: 13px;
}

.sidebar-footer {
  padding: 16px 24px;
  border-top: 1px solid var(--border-color);
  background: rgba(0, 0, 0, 0.2);
}

.version {
  font-family: 'Share Tech Mono', monospace;
  font-size: 12px;
  color: var(--text-muted);
}
</style>