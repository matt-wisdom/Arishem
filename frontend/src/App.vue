<script setup lang="ts">
import { watch } from 'vue'
import Sidebar from './components/Sidebar.vue'
import AppHeader from './components/Header.vue'
import { SignedIn, SignedOut, useAuth } from '@clerk/vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const { isSignedIn, isLoaded } = useAuth()

watch(
  [() => route.path, isSignedIn, isLoaded],
  ([path, signedIn, loaded]) => {
    if (loaded && signedIn === false) {
      if (path !== '/docs' && path !== '/sign-in' && path !== '/sign-up') {
        router.push('/sign-in')
      }
    }
  },
  { immediate: true }
)
</script>

<template>
  <div class="app-container">
    <SignedIn>
      <Sidebar />
      <div class="main-content">
        <AppHeader />
        <main class="content">
          <RouterView />
        </main>
      </div>
    </SignedIn>
    <SignedOut>
      <!-- If visiting /docs, show it directly with public shell -->
      <div v-if="route.path === '/docs'" class="main-content public-layout">
        <header class="public-header">
          <div class="logo" @click="router.push('/sign-in')">
            <img src="@/logo.jpeg" alt="Arishem Logo" class="logo-img" />
            <span class="logo-text">Arishem</span>
          </div>
          <button class="btn-primary sign-in-nav-btn" @click="router.push('/sign-in')">Sign In</button>
        </header>
        <main class="content">
          <RouterView />
        </main>
      </div>
      <!-- Otherwise render RouterView (which renders custom sign-in/sign-up page) -->
      <RouterView v-else />
    </SignedOut>
  </div>
</template>

<style>
@import url('https://fonts.googleapis.com/css2?family=Orbitron:wght@400;500;700;900&family=Rajdhani:wght@500;600;700&family=Share+Tech+Mono&display=swap');

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

:root {
  --bg-primary: #03030c;
  --bg-secondary: #080816;
  --bg-card: rgba(13, 13, 30, 0.85);
  --bg-card-hover: rgba(22, 22, 48, 0.9);
  --border-color: #1a1a3a;
  --text-primary: #f1f5f9;
  --text-secondary: #00ffcc;
  --text-muted: #6272a4;
  --accent: #00ffcc;
  --accent-hover: #33ffdd;
  --accent-glow: rgba(0, 255, 204, 0.25);
  --accent-pink: #ff007f;
  --accent-pink-glow: rgba(255, 0, 127, 0.3);
  --accent-yellow: #ffff00;
  --success: #39ff14;
  --warning: #ff9f00;
  --danger: #ff073a;
  --info: #00e5ff;
  --sidebar-width: 260px;
}

body {
  font-family: 'Rajdhani', sans-serif;
  background: var(--bg-primary);
  color: var(--text-primary);
  min-height: 100vh;
  line-height: 1.6;
  font-size: 16px;
}

.app-container {
  display: flex;
  min-height: 100vh;
}

.main-content {
  flex: 1;
  margin-left: var(--sidebar-width);
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.content {
  flex: 1;
  padding: 24px 32px;
  overflow-y: auto;
  overflow-x: hidden;
  background: radial-gradient(circle at 50% 50%, #0d0d21 0%, #03030c 100%);
  background-image: 
    linear-gradient(rgba(0, 255, 204, 0.02) 1px, transparent 1px), 
    linear-gradient(90deg, rgba(0, 255, 204, 0.02) 1px, transparent 1px);
  background-size: 30px 30px;
  position: relative;
}

/* Cyberpunk Scanline Effect */
.content::after {
  content: " ";
  display: block;
  position: absolute;
  top: 0; left: 0; bottom: 0; right: 0;
  background: linear-gradient(rgba(18, 16, 16, 0) 50%, rgba(0, 0, 0, 0.15) 50%), 
              linear-gradient(90deg, rgba(255, 0, 0, 0.006), rgba(0, 255, 0, 0.003), rgba(0, 0, 255, 0.006));
  background-size: 100% 4px, 3px 100%;
  z-index: 10;
  pointer-events: none;
}

/* Scrollbar */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: var(--bg-secondary);
}

::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 0;
}

::-webkit-scrollbar-thumb:hover {
  background: var(--accent);
  box-shadow: 0 0 10px var(--accent-glow);
}

/* Typography elements */
h1, h2, h3, h4, h5, h6 {
  font-family: 'Orbitron', sans-serif;
  letter-spacing: 1.5px;
  text-transform: uppercase;
  font-weight: 700;
}

/* Cyberpunk Card styles override */
.card, .stat-card {
  border-radius: 0 !important;
  border: 1px solid var(--border-color) !important;
  background: var(--bg-card) !important;
  backdrop-filter: blur(10px);
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.5) !important;
  position: relative;
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1) !important;
}

.card:hover, .stat-card:hover {
  border-color: var(--accent) !important;
  box-shadow: 0 0 15px var(--accent-glow) !important;
  transform: translateY(-2px);
}

.card::before, .stat-card::before {
  content: '';
  position: absolute;
  top: -1px;
  left: -1px;
  width: 10px;
  height: 10px;
  background: var(--accent);
  clip-path: polygon(0 0, 100% 0, 0 100%);
}

.card::after, .stat-card::after {
  content: '';
  position: absolute;
  bottom: -1px;
  right: -1px;
  width: 8px;
  height: 8px;
  background: var(--accent-pink);
  clip-path: polygon(100% 100%, 100% 0, 0 100%);
}

/* Blocky Forms */
input, select, textarea {
  font-family: 'Rajdhani', sans-serif !important;
  border-radius: 0 !important;
  background: rgba(8, 8, 22, 0.9) !important;
  border: 1px solid var(--border-color) !important;
  color: var(--text-primary) !important;
  font-weight: 500;
  transition: all 0.2s ease;
}

input:focus, select:focus, textarea:focus {
  outline: none;
  border-color: var(--accent) !important;
  box-shadow: 0 0 8px var(--accent-glow);
}

/* Button & Action styles */
button {
  font-family: 'Orbitron', sans-serif;
  cursor: pointer;
  border-radius: 0 !important;
  text-transform: uppercase;
  font-weight: 700;
  letter-spacing: 1px;
  transition: all 0.2s ease;
}

button:active {
  transform: scale(0.98);
}

/* Common Table styling override */
.table-row, tr {
  transition: all 0.2s ease;
}

.table-row:hover, tr:hover {
  background: rgba(0, 255, 204, 0.05) !important;
  box-shadow: inset 3px 0 0 var(--accent);
}

/* Authentication centered panel styles */
.auth-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  width: 100vw;
  background: var(--bg-primary);
  position: fixed;
  top: 0;
  left: 0;
  z-index: 9999;
}

.auth-card {
  position: relative;
  z-index: 10;
  border: 1px solid var(--border-color);
  box-shadow: 0 0 25px rgba(0, 255, 204, 0.15);
}

/* Public Layout & Header styles */
.public-layout {
  margin-left: 0 !important;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.public-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 32px;
  background: var(--bg-primary);
  border-bottom: 2px solid var(--border-color);
}

.public-header .logo {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
}

.public-header .logo-img {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  object-fit: cover;
  border: 1px solid var(--accent);
}

.public-header .logo-text {
  font-family: 'Orbitron', sans-serif;
  font-size: 18px;
  font-weight: 900;
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-pink) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: 1px;
}

.sign-in-nav-btn {
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  padding: 8px 16px;
  border: 1px solid var(--accent);
  background: transparent;
  color: var(--accent);
  transition: all 0.2s ease;
}

.sign-in-nav-btn:hover {
  background: var(--accent);
  color: var(--bg-primary);
  box-shadow: 0 0 10px var(--accent-glow);
}

/* Tooltip styling */
.tooltip-container {
  display: inline-flex;
  align-items: center;
  position: relative;
  cursor: pointer;
  margin-left: 6px;
  vertical-align: middle;
}

.tooltip-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 14px;
  height: 14px;
  background: var(--border-color);
  border: 1px solid var(--text-muted);
  color: var(--text-primary);
  border-radius: 50% !important;
  font-size: 10px;
  font-weight: bold;
  font-family: monospace;
  transition: all 0.2s ease;
}

.tooltip-icon:hover {
  background: var(--accent);
  border-color: var(--accent);
  color: var(--bg-primary);
  box-shadow: 0 0 6px var(--accent-glow);
}

.tooltip-bubble {
  visibility: hidden;
  position: absolute;
  bottom: 135%;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(8, 8, 22, 0.98);
  color: var(--text-primary);
  text-align: center;
  padding: 8px 12px;
  border: 1px solid var(--accent);
  border-radius: 0;
  font-family: 'Rajdhani', sans-serif;
  font-size: 12px;
  line-height: 1.4;
  white-space: normal;
  width: 220px;
  z-index: 9999;
  opacity: 0;
  transition: opacity 0.2s ease, visibility 0.2s ease;
  box-shadow: 0 0 10px rgba(0, 255, 204, 0.2);
  pointer-events: none;
}

/* Tooltip arrow */
.tooltip-bubble::after {
  content: "";
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  border-width: 5px;
  border-style: solid;
  border-color: var(--accent) transparent transparent transparent;
}

.tooltip-container:hover .tooltip-bubble {
  visibility: visible;
  opacity: 1;
}
</style>