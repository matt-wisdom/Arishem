<script setup lang="ts">
import { ref } from 'vue'

const integrations = ref([
  // { id: 'github', name: 'GitHub', description: 'Connect repositories for ', icon: 'github' },
  { id: 'slack', name: 'Slack', description: 'Receive alerts in Slack channels', icon: 'slack' },
  { id: 'jira', name: 'Jira', description: 'Create tickets from findings', icon: 'jira' },
])
</script>

<template>
  <div class="integrations-page">
    <Header><template #title>Integrations</template></Header>
    <div class="page-content">
      <div class="integrations-grid">
        <div v-for="int in integrations" :key="int.id" class="integration-card">
          <div class="int-icon" :class="int.icon">
            <svg v-if="int.icon === 'github'" viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/>
            </svg>
            <svg v-else-if="int.icon === 'slack'" viewBox="0 0 24 24" fill="currentColor">
              <path d="M5.042 15.165a2.528 2.528 0 01-2.52 2.523A2.528 2.528 0 010 15.165a2.527 2.527 0 012.522-2.52h2.52v2.52zm1.271 0a2.527 2.527 0 012.521-2.52 2.527 2.527 0 012.521 2.52v6.313A2.528 2.528 0 018.834 24a2.528 2.528 0 01-2.521-2.522v-6.313zM8.834 5.042a2.528 2.528 0 01-2.521-2.52A2.528 2.528 0 018.834 0a2.528 2.528 0 012.521 2.522v2.52H8.834zm0 1.271a2.528 2.528 0 012.521 2.521 2.528 2.528 0 01-2.521 2.521H2.522A2.528 2.528 0 010 8.834a2.528 2.528 0 012.522-2.521h6.312zm10.124 2.521a2.528 2.528 0 012.52-2.521A2.528 2.528 0 0124 8.834a2.528 2.528 0 01-2.522 2.521h-2.52V8.834zm-1.271 0a2.528 2.528 0 01-2.521 2.521 2.528 2.528 0 01-2.521-2.521V2.522A2.528 2.528 0 0115.166 0a2.528 2.528 0 012.521 2.522v6.312zm-2.521 10.124a2.528 2.528 0 012.521 2.52A2.528 2.528 0 0115.166 24a2.528 2.528 0 01-2.521-2.522v-2.52h2.521zm0-1.271a2.528 2.528 0 01-2.521-2.521 2.528 2.528 0 012.521-2.521h6.312A2.528 2.528 0 0124 15.166a2.528 2.528 0 01-2.522 2.521h-6.312z"/>
            </svg>
            <svg v-else viewBox="0 0 24 24" fill="currentColor">
              <path d="M11.571 11.513H0a5.218 5.218 0 005.232 5.215h2.13v2.057A5.215 5.215 0 0012.575 24V12.851a1.338 1.338 0 00-1.004-1.338zm5.723-5.756H5.736a5.215 5.215 0 005.215 5.214h2.129v2.058a5.218 5.218 0 005.215 5.214V7.089a1.332 1.332 0 00-1.001-1.332zM23.013 0H11.455a5.215 5.215 0 005.215 5.215h2.129v2.057A5.215 5.215 0 0024 12.574V1.332A1.338 1.338 0 0023.013 0z"/>
            </svg>
          </div>
          <div class="int-info">
            <h3>{{ int.name }}</h3>
            <p>{{ int.description }}</p>
          </div>
          <div class="int-status">
            <span class="not-implemented-badge">Not yet implemented</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.integrations-page { min-height: 100%; }
.page-content { padding-top: 24px; }
.integrations-grid { display: flex; flex-direction: column; gap: 16px; }
.integration-card { display: flex; align-items: center; gap: 20px; background: var(--bg-card); border: 1px solid var(--border-color); border-radius: 14px; padding: 24px; }
.int-icon { width: 56px; height: 56px; border-radius: 14px; display: flex; align-items: center; justify-content: center; }
.int-icon.github { background: #24292e; color: white; }
.int-icon.slack { background: #4A154B; color: white; }
.int-icon.jira { background: #0052CC; color: white; }
.int-icon svg { width: 28px; height: 28px; }
.int-info { flex: 1; }
.int-info h3 { font-size: 16px; font-weight: 600; color: var(--text-primary); margin-bottom: 4px; }
.int-info p { font-size: 13px; color: var(--text-muted); }
.not-implemented-badge { padding: 6px 14px; background: rgba(255, 159, 0, 0.15); color: var(--warning); border-radius: 8px; font-size: 13px; font-weight: 500; }
</style>