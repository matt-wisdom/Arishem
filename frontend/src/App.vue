<script setup lang="ts">
import Sidebar from './components/Sidebar.vue'
import AppHeader from './components/Header.vue'
import { SignedIn, SignedOut, SignIn } from '@clerk/vue'
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
      <div class="auth-container">
        <div class="scanlines"></div>
        <div class="grid-background"></div>
        <div class="auth-card">
          <SignIn />
        </div>
      </div>
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
</style>