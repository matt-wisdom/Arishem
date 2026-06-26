<script setup lang="ts">
import { ref } from 'vue'

const activeTab = ref('general')
</script>

<template>
  <div class="settings-page">
    <header class="header">
      <h1>Settings</h1>
    </header>

    <div class="tabs">
      <button :class="{ active: activeTab === 'general' }" @click="activeTab = 'general'">General</button>
      <button :class="{ active: activeTab === 'members' }" @click="activeTab = 'members'">Members</button>
      <button :class="{ active: activeTab === 'billing' }" @click="activeTab = 'billing'">Billing</button>
    </div>

    <div class="tab-content">
      <div v-if="activeTab === 'general'" class="tab-panel">
        <p>General settings will be available soon.</p>
      </div>
      <div v-if="activeTab === 'members'" class="tab-panel">
        <p>Member management will be available soon.</p>
      </div>
      <div v-if="activeTab === 'billing'" class="tab-panel">
        <p>Billing integration will be available soon.</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.settings-page { padding: 2rem; max-width: 800px; margin: 0 auto; }
.header { margin-bottom: 2rem; }
.tabs { display: flex; gap: 1rem; border-bottom: 1px solid #eee; margin-bottom: 2rem; }
.tabs button { background: none; border: none; padding: 0.75rem 1rem; cursor: pointer; border-bottom: 2px solid transparent; }
.tabs button.active { border-bottom-color: #2563eb; color: #2563eb; }
.tab-content { background: white; padding: 1.5rem; border-radius: 8px; }
</style>