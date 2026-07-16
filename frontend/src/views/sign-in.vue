<script setup lang="ts">
import { SignIn } from "@clerk/vue";
</script>

<template>
  <div class="auth-container">
    <div class="scanlines"></div>
    <div class="grid-background"></div>
    
    <div class="auth-header">
      <img src="@/logo.jpeg" alt="Arishem Logo" class="auth-logo" />
      <h1 class="auth-title">Arishem</h1>
      <router-link to="/docs" class="auth-docs-link">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
        </svg>
        View Documentation & Help
      </router-link>
    </div>
    
    <div class="auth-card">
      <SignIn />
    </div>
  </div>
</template>

<style scoped>
.auth-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  width: 100vw;
  background: var(--bg-primary);
  position: relative;
  overflow: hidden;
  padding: 24px;
}

.scanlines {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  background: linear-gradient(rgba(18, 16, 16, 0) 50%, rgba(0, 0, 0, 0.15) 50%), 
              linear-gradient(90deg, rgba(255, 0, 0, 0.006), rgba(0, 255, 0, 0.003), rgba(0, 0, 255, 0.006));
  background-size: 100% 4px, 3px 100%;
  pointer-events: none;
  z-index: 1;
}

.grid-background {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  background-image: 
    linear-gradient(rgba(0, 255, 204, 0.02) 1px, transparent 1px), 
    linear-gradient(90deg, rgba(0, 255, 204, 0.02) 1px, transparent 1px);
  background-size: 30px 30px;
  pointer-events: none;
  z-index: 0;
}

.auth-header {
  position: relative;
  z-index: 10;
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 32px;
  text-align: center;
}

.auth-logo {
  width: 80px;
  height: 80px;
  border-radius: 12px;
  border: 1px solid var(--accent);
  box-shadow: 0 0 15px var(--accent-glow);
  object-fit: cover;
  margin-bottom: 16px;
  transition: all 0.3s ease;
}

.auth-logo:hover {
  transform: scale(1.05);
  box-shadow: 0 0 25px var(--accent-glow);
}

.auth-title {
  font-family: 'Orbitron', sans-serif;
  font-size: 36px;
  font-weight: 900;
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-pink) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: 2px;
  text-shadow: 0 0 15px var(--accent-glow);
  margin-bottom: 12px;
}

.auth-docs-link {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--text-muted);
  text-decoration: none;
  font-family: 'Orbitron', sans-serif;
  font-size: 13px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border: 1px solid var(--border-color);
  background: var(--bg-card);
  padding: 8px 16px;
  transition: all 0.2s ease;
}

.auth-docs-link svg {
  width: 14px;
  height: 14px;
}

.auth-docs-link:hover {
  color: var(--accent);
  border-color: var(--accent);
  box-shadow: 0 0 10px var(--accent-glow);
  text-shadow: 0 0 5px var(--accent-glow);
}

.auth-card {
  position: relative;
  z-index: 10;
  border: 1px solid var(--border-color);
  box-shadow: 0 0 25px rgba(0, 255, 204, 0.15);
  border-radius: 8px;
}
</style>