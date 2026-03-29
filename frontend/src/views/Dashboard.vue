<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '@/api'
import { useAuthStore } from '@/stores/auth'
import { useCache, CACHE_KEYS, CACHE_TTL } from '@/composables/useCache'

const authStore = useAuthStore()
const matches = ref([])
const loading = ref(true)

// 下注對話框狀態
const showBetModal = ref(false)
const selectedMatch = ref(null)
const betAmount = ref(100)
const betChoice = ref('A') // A, B, Draw
const betLoading = ref(false)
const userBetsMap = ref({}) // { matchId: { A: 100, B: 50, Draw: 0, total: 150 } }

const { get, set, invalidate } = useCache()

// 判斷全部賽事是否都尚未開始（賽前狀態）
const allMatchesPending = (list) => {
  return list.every(m => m.status !== 'active' && !new Date(m.startTime) <= new Date())
}

const fetchMatches = async (forceRefresh = false) => {
  // 將 cache 裡的資料路徑者先返回 cache
  if (!forceRefresh) {
    const cached = get(CACHE_KEYS.MATCHES)
    if (cached) {
      matches.value = cached
      loading.value = false
      return
    }
  }

  try {
    const res = await api.get('/matches')
    matches.value = res.data

    // 判斷是否要儲存 cache：如果有任何一場賽事已開您，則不 cache
    const hasStarted = res.data.some(m => new Date(m.startTime) <= new Date())
    const ttl = hasStarted ? CACHE_TTL.MATCHES_ACTIVE : CACHE_TTL.MATCHES_PENDING
    set(CACHE_KEYS.MATCHES, res.data, ttl)
  } catch (error) {
    console.error('Failed to fetch matches', error)
  } finally {
    loading.value = false
  }
}

const fetchUserBets = async () => {
  if (!authStore.isAuthenticated) return
  try {
    const res = await api.get('/user/bets')
    const map = {}
    res.data.forEach(b => {
      if (!map[b.matchId]) map[b.matchId] = { A: 0, B: 0, total: 0 }
      map[b.matchId][b.choice] += b.amount
      map[b.matchId].total += b.amount
    })
    userBetsMap.value = map
  } catch (error) {
    console.error(error)
  }
}

onMounted(() => {
  fetchMatches()
  fetchUserBets()
})

const getPoolPercentage = (pool, choice) => {
  const total = pool.A + pool.B
  if (total === 0) return 0 // 移除預設比例，變成空吧拉
  return Math.round((pool[choice] / total) * 100)
}

const getOdds = (pool, choice) => {
  const total = pool.A + pool.B
  // 防雷：無人下注該區時預設賠率 (也可顯示為 - )
  if (total === 0) return '尚未有注金'
  if (pool[choice] === 0) return '高賠' // 或者可以顯示為極高數值
  // 賠率 = 該比賽的總注金 / 該選項的總注金 (共注池概念：1 賠 X)
  const odds = total / pool[choice]
  return odds.toFixed(2) + 'x'
}

const openBetModal = (match) => {
  if (!authStore.isAuthenticated) {
    alert('請先登入後再進行下注！')
    return
  }
  selectedMatch.value = match
  betChoice.value = 'A'
  betAmount.value = 100
  showBetModal.value = true
}

const submitBet = async () => {
  if (betAmount.value <= 0 || betAmount.value > authStore.currentPoints) {
    alert('點數不足或輸入錯誤！')
    return
  }

  betLoading.value = true
  try {
    const res = await api.post('/bets/', {
      matchId: selectedMatch.value.id,
      choice: betChoice.value,
      amount: betAmount.value
    })
    
    // 更新餘額
    authStore.user.points = res.data.newPoints
    localStorage.setItem('user', JSON.stringify(authStore.user))
    
    alert('下注成功！')
    showBetModal.value = false
    // 下注成功後主動清除 matches cache，確保前端能載到最新的彩池資料
    invalidate(CACHE_KEYS.MATCHES)
    fetchMatches() // 更新彩池比例
    fetchUserBets() // 更新自己的下注狀態
  } catch (error) {
    alert(error.response?.data?.error || '下注失敗')
  } finally {
    betLoading.value = false
  }
}

const formatDate = (dateStr) => {
  const d = new Date(dateStr)
  const pad = (n) => n.toString().padStart(2, '0')
  return `${d.getFullYear()}/${pad(d.getMonth() + 1)}/${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const isMatchStarted = (startTime) => {
  return new Date(startTime) <= new Date()
}

// 將賽事依狀態分類
const activeMatches = computed(() => {
  return matches.value.filter(m => m.status === 'active' || (m.status === 'pending' && isMatchStarted(m.startTime)))
})

const pendingMatches = computed(() => {
  return matches.value.filter(m => m.status === 'pending' && !isMatchStarted(m.startTime))
})

const endedMatches = computed(() => {
  return matches.value.filter(m => m.status === 'ended')
})
</script>

<template>
  <div class="dashboard">
    <header class="page-header">
      <h1>賽事大廳</h1>
      <p class="subtitle">選擇您支持的隊伍，贏得更多積分！</p>
    </header>

    <div v-if="loading" class="loading">載入賽事中...</div>
    
    <div v-else-if="matches.length === 0" class="empty-state premium-card" style="margin-top: 2rem;">
      <h2 style="font-size: 1.5rem; margin-bottom: 0.5rem; color: var(--text-main);">目前沒有賽事</h2>
      <p>管理員尚未建立任何賽事，請稍後再回來查看最新賽程！</p>
    </div>

    <div v-else class="matches-layout">
      
      <!-- 進行中賽事 (開放下注) -->
      <section v-if="activeMatches.length > 0" class="match-section">
        <h2 class="section-title"><span class="live-dot"></span> 進行中 & 開放下注</h2>
        <div class="matches-grid">
          <!-- 卡片模板改用 component 會更好，但為節省時間直接展開 -->
          <div class="premium-card match-card" v-for="match in activeMatches" :key="match.id">
            <div class="match-header">
              <span class="match-time">{{ formatDate(match.startTime) }}</span>
              <span class="match-status active">開放下注</span>
            </div>
            
            <div class="teams-container">
              <div class="team"><div class="team-name">{{ match.teamA }}</div></div>
              <div class="vs">VS</div>
              <div class="team"><div class="team-name">{{ match.teamB }}</div></div>
            </div>

            <div class="pool-bar-container">
              <div class="pool-labels">
                <span>主勝 {{ getPoolPercentage(match.pool, 'A') }}% ({{ getOdds(match.pool, 'A') }})</span>
                <span>客勝 {{ getPoolPercentage(match.pool, 'B') }}% ({{ getOdds(match.pool, 'B') }})</span>
              </div>
              <div class="pool-bar">
                <div class="bar-segment a" :style="{ width: getPoolPercentage(match.pool, 'A') + '%' }"></div>
                <div class="bar-segment b" :style="{ width: getPoolPercentage(match.pool, 'B') + '%' }"></div>
              </div>
            </div>

            <div class="actions">
              <button class="btn btn-primary" @click="openBetModal(match)">
                <template v-if="userBetsMap[match.id]?.total > 0">繼續下注</template>
                <template v-else>立即下注</template>
              </button>
            </div>
          </div>
        </div>
      </section>

      <!-- 即將開始賽事 (未來) -->
      <section v-if="pendingMatches.length > 0" class="match-section">
        <h2 class="section-title" style="color: var(--primary);">⏳ 即將開始 (未來賽事)</h2>
        <div class="matches-grid">
          <div class="premium-card match-card" v-for="match in pendingMatches" :key="match.id">
            <div class="match-header">
              <span class="match-time">{{ formatDate(match.startTime) }}</span>
              <span class="match-status pending">即將開始</span>
            </div>
            
            <div class="teams-container">
              <div class="team"><div class="team-name">{{ match.teamA }}</div></div>
              <div class="vs">VS</div>
              <div class="team"><div class="team-name">{{ match.teamB }}</div></div>
            </div>

            <div class="pool-bar-container">
              <div class="pool-labels">
                <span>主勝 {{ getPoolPercentage(match.pool, 'A') }}% ({{ getOdds(match.pool, 'A') }})</span>
                <span>客勝 {{ getPoolPercentage(match.pool, 'B') }}% ({{ getOdds(match.pool, 'B') }})</span>
              </div>
              <div class="pool-bar">
                <div class="bar-segment a" :style="{ width: getPoolPercentage(match.pool, 'A') + '%' }"></div>
                <div class="bar-segment b" :style="{ width: getPoolPercentage(match.pool, 'B') + '%' }"></div>
              </div>
            </div>

            <div class="actions">
              <button class="btn btn-primary" @click="openBetModal(match)">
                <template v-if="userBetsMap[match.id]?.total > 0">繼續下注</template>
                <template v-else>預先下注</template>
              </button>
            </div>
          </div>
        </div>
      </section>

      <!-- 已結束賽事 -->
      <section v-if="endedMatches.length > 0" class="match-section opacity-75">
        <h2 class="section-title" style="color: var(--text-muted);">🏁 已結束賽事</h2>
        <div class="matches-grid">
          <div class="premium-card match-card" v-for="match in endedMatches" :key="match.id">
            <div class="match-header">
              <span class="match-time">{{ formatDate(match.startTime) }}</span>
              <span class="match-status" style="background: rgba(255,255,255,0.05);">已結束</span>
            </div>
            
            <div class="teams-container">
              <div class="team"><div class="team-name">{{ match.teamA }}</div></div>
              <div class="vs">VS</div>
              <div class="team"><div class="team-name">{{ match.teamB }}</div></div>
            </div>

            <div class="pool-bar-container">
              <div class="pool-labels">
                <span>主勝 {{ getPoolPercentage(match.pool, 'A') }}%</span>
                <span>客勝 {{ getPoolPercentage(match.pool, 'B') }}%</span>
              </div>
              <div class="pool-bar" style="opacity: 0.5;">
                <div class="bar-segment a" :style="{ width: getPoolPercentage(match.pool, 'A') + '%' }"></div>
                <div class="bar-segment b" :style="{ width: getPoolPercentage(match.pool, 'B') + '%' }"></div>
              </div>
            </div>

            <div class="actions">
              <button class="btn btn-primary disabled-btn" disabled>
                已結束
              </button>
            </div>
          </div>
        </div>
      </section>

    </div>

    <!-- 下注對話框 (Modal) -->
    <div class="modal-overlay" v-if="showBetModal" @click.self="showBetModal = false">
      <div class="premium-card modal-content">
        <h2>下注：{{ selectedMatch?.teamA }} VS {{ selectedMatch?.teamB }}</h2>
        
        <div class="form-group" style="margin-top: 1.5rem;">
          <label>選擇預測結果 (即時賠率預估結算)</label>
          <div class="choice-group">
            <button :class="['choice-btn', { active: betChoice === 'A' }]" @click="betChoice = 'A'">
              <div>{{ selectedMatch?.teamA }} 勝</div>
              <div class="odds-text">{{ getOdds(selectedMatch?.pool, 'A') }}</div>
              <div v-if="userBetsMap[selectedMatch?.id]?.A > 0" class="my-bet-text">已押: {{ userBetsMap[selectedMatch?.id].A }}</div>
            </button>
            <button :class="['choice-btn', { active: betChoice === 'B' }]" @click="betChoice = 'B'">
              <div>{{ selectedMatch?.teamB }} 勝</div>
              <div class="odds-text">{{ getOdds(selectedMatch?.pool, 'B') }}</div>
              <div v-if="userBetsMap[selectedMatch?.id]?.B > 0" class="my-bet-text">已押: {{ userBetsMap[selectedMatch?.id].B }}</div>
            </button>
          </div>
        </div>

        <div class="form-group" style="margin-top: 1.5rem;">
          <label>投入積分 (目前餘額: {{ authStore.currentPoints }})</label>
          <input type="number" class="input-form current-bet" v-model="betAmount" min="1" :max="authStore.currentPoints" />
        </div>

        <div class="modal-actions">
          <button class="btn btn-outline" @click="showBetModal = false">取消</button>
          <button class="btn btn-primary" @click="submitBet" :disabled="betLoading">確認下注</button>
        </div>
      </div>
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
  background: linear-gradient(135deg, #f8fafc, #94a3b8);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.subtitle {
  color: var(--text-muted);
  font-size: 1.125rem;
}

.loading, .empty-state {
  text-align: center;
  padding: 4rem;
  color: var(--text-muted);
  font-size: 1.125rem;
}

.matches-layout {
  display: flex;
  flex-direction: column;
  gap: 3rem;
}

.match-section {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.section-title {
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--text-main);
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 0.75rem;
}

.live-dot {
  width: 12px;
  height: 12px;
  background-color: var(--success);
  border-radius: 50%;
  display: inline-block;
  box-shadow: 0 0 8px var(--success);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.7); }
  70% { transform: scale(1); box-shadow: 0 0 0 6px rgba(16, 185, 129, 0); }
  100% { transform: scale(0.95); box-shadow: 0 0 0 0 rgba(16, 185, 129, 0); }
}

.opacity-75 {
  opacity: 0.75;
}

.matches-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 2rem;
}

.match-card {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.match-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.875rem;
}

.match-time {
  color: var(--text-muted);
}

.match-status {
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-weight: 600;
  font-size: 0.75rem;
  background-color: var(--surface-hover);
  color: var(--text-muted);
}

.match-status.active {
  background-color: rgba(16, 185, 129, 0.2);
  color: var(--success);
}
.match-status.pending {
  background-color: rgba(56, 189, 248, 0.2);
  color: var(--primary);
}

.teams-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 0;
}

.team {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  flex: 1;
}

.team-name {
  font-weight: 600;
  font-size: 1.25rem;
  text-align: center;
}

.vs {
  font-weight: 700;
  color: var(--text-muted);
  font-size: 1.5rem;
  padding: 0 1rem;
}

.pool-bar-container {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.pool-labels {
  display: flex;
  justify-content: space-between;
  font-size: 0.75rem;
  color: var(--text-muted);
}

.pool-bar {
  height: 8px;
  width: 100%;
  background-color: var(--surface-hover);
  border-radius: 4px;
  display: flex;
  overflow: hidden;
}

.bar-segment {
  height: 100%;
  transition: width 0.3s ease;
}

.bar-segment.a { background-color: var(--primary); }
.bar-segment.draw { background-color: var(--text-muted); }
.bar-segment.b { background-color: var(--accent); }

.actions {
  display: flex;
  justify-content: center;
  margin-top: 0.5rem;
}

.actions .btn {
  width: 100%;
}

.disabled-btn {
  opacity: 0.5;
  cursor: not-allowed;
  background: var(--surface-hover);
  color: var(--text-muted);
  box-shadow: none;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: var(--overlay-bg);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  width: 90%;
  max-width: 500px;
  padding: 2rem;
}

.modal-content h2 {
  margin-top: 0;
  font-size: 1.25rem;
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 1rem;
}

.choice-group {
  display: flex;
  gap: 0.5rem;
}

.choice-btn {
  flex: 1;
  padding: 0.75rem;
  background-color: var(--bg-color);
  border: 1px solid var(--border-color);
  color: var(--text-muted);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  font-weight: 600;
}

.choice-btn:hover {
  background-color: var(--surface-hover);
}

.choice-btn .odds-text {
  font-size: 0.75rem;
  color: var(--primary);
  margin-top: 0.25rem;
  font-weight: 700;
}

.choice-btn.active {
  background-color: rgba(56, 189, 248, 0.2);
  border-color: var(--primary);
  color: var(--primary);
}

.choice-btn .my-bet-text {
  font-size: 0.75rem;
  color: var(--success);
  margin-top: 0.25rem;
  font-weight: 600;
}

.current-bet {
  font-size: 1.25rem;
  font-weight: 700;
  text-align: center;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
}
</style>
