<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'

const bets = ref([])
const loading = ref(true)

const fetchBets = async () => {
  try {
    const res = await api.get('/user/bets')
    bets.value = res.data
  } catch (error) {
    console.error('Failed to fetch bets', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchBets()
})

const getChoiceText = (choice, match) => {
  if (choice === 'A') return match.teamA + ' 勝'
  if (choice === 'B') return match.teamB + ' 勝'
  return '和局'
}

const formatDate = (dateStr) => {
  const d = new Date(dateStr)
  const pad = (n) => n.toString().padStart(2, '0')
  return `${d.getFullYear()}/${pad(d.getMonth() + 1)}/${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}
</script>

<template>
  <div class="my-bets-page">
    <header class="page-header">
      <h1>我的下注紀錄</h1>
      <p class="subtitle">追蹤您的所有預測與收益。</p>
    </header>

    <div class="premium-card bets-card">
      <div v-if="loading" class="loading">載入紀錄中...</div>
      
      <div v-else-if="bets.length === 0" class="empty-state">
        您還沒有任何下注紀錄
      </div>

      <div v-else class="bets-list">
        <div class="bet-item" v-for="bet in bets" :key="bet.id">
          <div class="bet-info">
            <div class="bet-match">{{ bet.match.teamA }} VS {{ bet.match.teamB }}</div>
            <div class="bet-time">下注時間: {{ formatDate(bet.createdAt) }}</div>
          </div>
          <div class="bet-choice">
            <span class="label">預測:</span>
            <span class="value">{{ getChoiceText(bet.choice, bet.match) }}</span>
          </div>
          <div class="bet-amount">
            <span class="label">投入點數:</span>
            <span class="value">{{ bet.amount }}</span>
          </div>
          
          <div v-if="bet.status === 'pending'" class="bet-status pending">待結算</div>
          <div v-else-if="bet.status === 'won'" class="bet-status won">贏得 {{ bet.payout }} !!</div>
          <div v-else class="bet-status lost">未猜中</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page-header {
  margin-bottom: 2rem;
}

.page-header h1 {
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
}

.subtitle {
  color: var(--text-muted);
}

.bets-card {
  padding: 0;
  overflow: hidden;
  min-height: 200px;
}

.loading, .empty-state {
  text-align: center;
  padding: 4rem;
  color: var(--text-muted);
}

.bets-list {
  display: flex;
  flex-direction: column;
}

.bet-item {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr;
  padding: 1.5rem;
  align-items: center;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.bet-item:last-child {
  border-bottom: none;
}

.bet-match {
  font-weight: 600;
  font-size: 1.125rem;
  margin-bottom: 0.25rem;
}

.bet-time {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.label {
  font-size: 0.75rem;
  color: var(--text-muted);
  display: block;
  margin-bottom: 0.25rem;
}

.value {
  font-weight: 600;
}

.bet-status {
  justify-self: end;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 600;
}

.bet-status.pending {
  background-color: rgba(245, 158, 11, 0.2);
  color: var(--warning);
}
.bet-status.won {
  background-color: rgba(16, 185, 129, 0.2);
  color: var(--success);
}
.bet-status.lost {
  background-color: rgba(239, 68, 68, 0.2);
  color: var(--danger);
}
</style>
