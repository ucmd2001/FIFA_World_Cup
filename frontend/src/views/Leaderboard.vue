<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'
import { useCache, CACHE_KEYS, CACHE_TTL } from '@/composables/useCache'

const leaderboard = ref([])
const loading = ref(true)
const fromCache = ref(false)
const cacheRemainingMin = ref(0)

const { get, set, getCachedAt, getRemainingTtl } = useCache()

const fetchLeaderboard = async (forceRefresh = false) => {
  // 嘗試從 cache 取資料
  if (!forceRefresh) {
    const cached = get(CACHE_KEYS.LEADERBOARD)
    if (cached) {
      leaderboard.value = cached
      fromCache.value = true
      cacheRemainingMin.value = Math.ceil(getRemainingTtl(CACHE_KEYS.LEADERBOARD) / 60000)
      loading.value = false
      return
    }
  }

  // Cache miss -> 打 API
  fromCache.value = false
  try {
    const res = await api.get('/leaderboard')
    leaderboard.value = res.data
    set(CACHE_KEYS.LEADERBOARD, res.data, CACHE_TTL.LEADERBOARD)
    cacheRemainingMin.value = Math.ceil(CACHE_TTL.LEADERBOARD / 60000)
  } catch (error) {
    console.error('Failed to fetch leaderboard', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchLeaderboard()
})
</script>


<template>
  <div class="leaderboard-page">
    <header class="page-header">
      <h1>積分排行榜</h1>
      <p class="subtitle">看看誰是世足預測王！</p>
      <div class="cache-info" v-if="!loading">
        <span v-if="fromCache" class="cache-badge">
          📆 快取資料（剩餘 {{ cacheRemainingMin }} 分鐘到期）
        </span>
        <span v-else class="fresh-badge">
          ⚡ 即時資料
        </span>
        <button class="refresh-btn" @click="fetchLeaderboard(true)" title="手動刷新">
          🔄
        </button>
      </div>
    </header>

    <div class="premium-card leaderboard-card">
      <div v-if="loading" style="text-align:center; padding: 2rem;">載入中...</div>
      
      <table class="list-table" v-else>
        <thead class="list-header">
          <tr>
            <th class="col-rank">排名</th>
            <th class="col-user">員工</th>
            <th class="col-points">總積分</th>
          </tr>
        </thead>
        
        <tbody class="list-body">
          <tr class="list-row" v-for="(user, idx) in leaderboard" :key="user.id">
            <td class="col-rank">
              <span :class="['rank-badge', `rank-${idx + 1}`]" v-if="idx < 3">{{ idx + 1 }}</span>
              <span v-else>{{ idx + 1 }}</span>
            </td>
            <td class="col-user">
              <div class="avatar">{{ user.username.charAt(0).toUpperCase() }}</div>
              <span>{{ user.username }}</span>
            </td>
            <td class="col-points highlight-points">{{ user.points }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
.page-header {
  margin-bottom: 2rem;
  text-align: center;
}

.page-header h1 {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  background: linear-gradient(135deg, #f59e0b, #fbbf24);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.cache-info {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  margin-top: 0.5rem;
  font-size: 0.8rem;
}

.cache-badge {
  color: var(--text-muted);
  background: rgba(255,255,255,0.05);
  padding: 0.2rem 0.6rem;
  border-radius: 9999px;
  border: 1px solid var(--border-color);
}

.fresh-badge {
  color: var(--success);
  font-weight: 600;
}

.refresh-btn {
  background: transparent;
  border: none;
  font-size: 1rem;
  cursor: pointer;
  opacity: 0.6;
  transition: opacity 0.2s, transform 0.2s;
  padding: 0;
}

.refresh-btn:hover {
  opacity: 1;
  transform: rotate(30deg);
}

.subtitle {
  color: var(--text-muted);
  font-size: 1.125rem;
}

.leaderboard-card {
  max-width: 800px;
  margin: 0 auto;
  padding: 0;
  overflow: hidden;
}

.list-table {
  width: 100%;
  border-collapse: collapse;
}

.list-header {
  background-color: var(--table-header-bg);
  border-bottom: 1px solid var(--border-color);
}

.list-header th {
  padding: 1rem 1.5rem;
  font-weight: 600;
  color: var(--text-muted);
  font-size: 0.875rem;
  text-align: left;
}

.list-row {
  border-bottom: 1px solid var(--table-row-border);
  transition: background-color 0.2s;
}

.list-row:hover {
  background-color: var(--table-row-hover);
}

.list-row td {
  padding: 1rem 1.5rem;
}

.col-rank {
  width: 80px;
  text-align: center !important;
  font-weight: 600;
}

.col-user {
  display: flex;
  align-items: center;
  gap: 1rem;
  font-weight: 500;
}

.col-points {
  width: 120px;
  text-align: right !important;
  font-weight: 700;
}

.highlight-points {
  color: #fbbf24;
}

.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary), var(--accent));
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.875rem;
  font-weight: bold;
}

.rank-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  color: white;
}

.rank-1 { background: linear-gradient(135deg, #fbbf24, #f59e0b); box-shadow: 0 0 10px rgba(245, 158, 11, 0.5); }
.rank-2 { background: linear-gradient(135deg, #94a3b8, #64748b); }
.rank-3 { background: linear-gradient(135deg, #b45309, #78350f); }
</style>
