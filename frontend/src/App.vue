<script setup>
import { ref, onMounted } from 'vue'
import { RouterView, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()
const theme = ref(localStorage.getItem('theme') || 'dark')

const toggleTheme = () => {
  theme.value = theme.value === 'dark' ? 'light' : 'dark'
  document.documentElement.setAttribute('data-theme', theme.value)
  localStorage.setItem('theme', theme.value)
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

// 由於後端已經改為 CronJob 自動發放積分，不需要在前端這裡每次呼叫了
onMounted(async () => {
  // 設置初始主題
  document.documentElement.setAttribute('data-theme', theme.value)

  if (authStore.isAuthenticated) {
    await authStore.fetchMe()
  }
})
</script>

<template>
  <div class="app-layout">
    <!-- Navbar -->
    <header class="navbar">
      <div class="nav-brand">🏆 NextLink World Cup</div>
      <nav class="nav-links" v-if="authStore.isAuthenticated">
        <router-link to="/">賽事大廳</router-link>
        <router-link to="/my-bets">我的下注</router-link>
        <router-link to="/leaderboard">排行榜</router-link>
        <router-link to="/admin" v-if="authStore.isAdmin" class="admin-link">管理員</router-link>
        <button class="theme-toggle" @click="toggleTheme" :title="theme === 'dark' ? '切換為亮色模式' : '切換為暗色模式'">
          {{ theme === 'dark' ? '☀️' : '🌙' }}
        </button>
      </nav>
      
      <div class="user-status" v-if="authStore.isAuthenticated">
        <span>Hi, {{ authStore.user?.username }} | </span>
        <span class="points-badge">餘額: {{ authStore.currentPoints }} 點</span>
        <button class="logout-btn" @click="handleLogout">登出</button>
      </div>
      <div v-else style="display: flex; gap: 1rem; align-items: center;">
        <button class="theme-toggle" @click="toggleTheme" :title="theme === 'dark' ? '切換為亮色模式' : '切換為暗色模式'">
          {{ theme === 'dark' ? '☀️' : '🌙' }}
        </button>
        <router-link to="/login" class="btn btn-outline" style="padding: 0.5rem 1rem;">登入</router-link>
      </div>
    </header>

    <!-- Main Content -->
    <main class="main-content">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
.app-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background-color: var(--surface-color);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--border-color);
  position: sticky;
  top: 0;
  z-index: 100;
}

.nav-brand {
  font-size: 1.5rem;
  font-weight: 700;
  color: #38bdf8;
  letter-spacing: -0.025em;
}

.nav-links {
  display: flex;
  gap: 2rem;
}

.nav-links a {
  color: var(--nav-link-color);
  text-decoration: none;
  font-weight: 500;
  transition: color 0.2s;
}

.nav-links a:hover, .nav-links a.router-link-active {
  color: var(--nav-link-active);
}

.admin-link {
  color: var(--danger) !important;
}

.user-status {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 600;
  font-size: 0.875rem;
}

.points-badge {
  background: linear-gradient(135deg, #3b82f6, #8b5cf6);
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  box-shadow: 0 4px 14px 0 rgba(139, 92, 246, 0.39);
}

.logout-btn {
  background: transparent;
  border: 1px solid var(--border-color);
  color: var(--text-muted);
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  cursor: pointer;
  margin-left: 0.5rem;
  transition: all 0.2s;
}

.logout-btn:hover {
  background: var(--surface-hover);
  color: var(--text-main);
}

.theme-toggle {
  background: transparent;
  border: none;
  font-size: 1.25rem;
  cursor: pointer;
  padding: 0.25rem;
  border-radius: 50%;
  transition: transform 0.2s;
}

.theme-toggle:hover {
  transform: scale(1.1);
}

.main-content {
  flex: 1;
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
}
</style>
