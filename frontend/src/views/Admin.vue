<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'

const teamA = ref('')
const teamB = ref('')
const startTime = ref('')

const matches = ref([])
const users = ref([])
const loadingMatches = ref(true)
const loadingUsers = ref(true)
const activeTab = ref('matches') // 'matches' | 'users'

// 使用者編輯彈窗狀態
const showUserModal = ref(false)
const selectedUser = ref(null)
const editUserForm = ref({
  name: '',
  email: '',
  points: 0,
  role: 'user'
})

const fetchMatches = async () => {
  try {
    const res = await api.get('/matches')
    matches.value = res.data
  } catch (error) {
    console.error('Failed to fetch matches', error)
  } finally {
    loadingMatches.value = false
  }
}

const fetchUsers = async () => {
  try {
    const res = await api.get('/admin/users')
    users.value = res.data
  } catch (error) {
    console.error('Failed to fetch users', error)
  } finally {
    loadingUsers.value = false
  }
}

const openUserEdit = (user) => {
  selectedUser.value = user
  editUserForm.value = {
    name: user.name,
    email: user.email,
    points: user.points,
    role: user.role
  }
  showUserModal.value = true
}

const submitUserEdit = async () => {
  try {
    await api.put(`/admin/users/${selectedUser.value.id}`, editUserForm.value)
    alert('使用者資訊已更新！')
    showUserModal.value = false
    fetchUsers()
  } catch (error) {
    alert(error.response?.data?.error || '更新失敗')
  }
}

const deleteUser = async (user) => {
  if (!confirm(`確定要刪除使用者 ${user.username} 嗎？此操作不可逆！`)) return
  try {
    await api.delete(`/admin/users/${user.id}`)
    alert('刪除成功！')
    fetchUsers()
  } catch (error) {
    alert('刪除失敗')
  }
}

onMounted(() => {
  fetchMatches()
  fetchUsers()
})

const handleCreateMatch = async () => {
  if (!teamA.value || !teamB.value || !startTime.value) {
    alert('請填寫完整資訊')
    return
  }

  // 將時間轉為 ISO 格式 (UTC) 送給後端
  const isoTime = new Date(startTime.value).toISOString()

  try {
    await api.post('/admin/matches', {
      teamA: teamA.value,
      teamB: teamB.value,
      startTime: isoTime
    })
    alert('新增賽事成功！')
    teamA.value = ''
    teamB.value = ''
    startTime.value = ''
    fetchMatches()
  } catch (error) {
    alert(error.response?.data?.error || '新增失敗')
  }
}

const isSyncing = ref(false)
const handleSyncMatches = async () => {
  isSyncing.value = true
  try {
    const res = await api.post('/admin/matches/sync')
    alert(`成功同步了 ${res.data.count} 場賽事！`)
    fetchMatches()
  } catch (error) {
    alert(error.response?.data?.error || '同步失敗')
  } finally {
    isSyncing.value = false
  }
}

const resolves = ref({})

const handleResolve = async (matchId) => {
  const result = resolves.value[matchId]
  if (!result) {
    alert('請選擇勝負結果！')
    return
  }

  if (!confirm('確認送出結算？送出後即發放獎金且無法更改！')) return

  try {
    await api.post(`/admin/matches/${matchId}/resolve`, { result })
    alert('賽事已成功結算！')
    fetchMatches()
  } catch (error) {
    alert(error.response?.data?.error || '結算失敗')
  }
}

const formatDate = (dateStr) => {
  const d = new Date(dateStr)
  const pad = (n) => n.toString().padStart(2, '0')
  return `${d.getFullYear()}/${pad(d.getMonth() + 1)}/${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}
</script>

<template>
  <div class="admin-page">
    <header class="page-header">
      <h1>管理員後台</h1>
      <p class="subtitle">賽事管理與使用者管理中心</p>
    </header>

    <div class="tabs">
      <button class="tab" :class="{ active: activeTab === 'matches' }" @click="activeTab = 'matches'">賽事管理</button>
      <button class="tab" :class="{ active: activeTab === 'users' }" @click="activeTab = 'users'">使用者清單</button>
    </div>

    <!-- 賽事管理 -->
    <div v-if="activeTab === 'matches'" class="admin-grid">
      <!-- 新增賽事表單 -->
      <div class="premium-card">
        <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1.5rem; border-bottom: 1px solid var(--border-color); padding-bottom: 0.5rem;">
          <h2 style="margin: 0; border: none; padding: 0;">新增賽事</h2>
          <button class="btn btn-outline" style="font-size: 0.875rem; padding: 0.25rem 0.75rem; border-color: var(--primary); color: var(--primary);" @click="handleSyncMatches" :disabled="isSyncing">
            {{ isSyncing ? '同步中...' : '🤖 API 自動同步' }}
          </button>
        </div>
        
        <form class="admin-form" @submit.prevent="handleCreateMatch">
          <div class="form-group">
            <label>主場隊伍</label>
            <input type="text" v-model="teamA" class="input-form" placeholder="例如：阿根廷" required />
          </div>
          <div class="form-group">
            <label>客場隊伍</label>
            <input type="text" v-model="teamB" class="input-form" placeholder="例如：法國" required />
          </div>
          <div class="form-group">
            <label>比賽時間</label>
            <input type="datetime-local" v-model="startTime" class="input-form" required />
          </div>
          <button type="submit" class="btn btn-primary" style="margin-top: 1rem;">新增賽事</button>
        </form>
      </div>

      <!-- 結算賽事列表 -->
      <div class="premium-card">
        <h2>待結算賽事</h2>
        
        <div v-if="loadingMatches" class="loading">載入賽事中...</div>
        
        <div class="matches-list" v-else>
          <div class="match-item" v-for="match in matches" :key="match.id">
            <div class="match-info">
              <div class="teams">{{ match.teamA }} VS {{ match.teamB }}</div>
              <div class="time">{{ formatDate(match.startTime) }}</div>
              <span class="status-badge" :class="match.status">{{ match.status }}</span>
            </div>
            
            <div class="resolve-actions" v-if="match.status !== 'ended'">
              <select class="input-form select-result" v-model="resolves[match.id]">
                <option :value="undefined" disabled>選擇結果</option>
                <option value="A">主勝 ({{ match.teamA }})</option>
                <option value="B">客勝 ({{ match.teamB }})</option>
              </select>
              <button 
                class="btn btn-outline" 
                style="color: var(--success); border-color: var(--success);"
                @click="handleResolve(match.id)"
              >結算</button>
            </div>
            
            <div v-else class="resolved-text">
              已結算: {{ match.result === 'A' ? match.teamA : match.teamB }}
            </div>
          </div>
        </div>
      </div>
    </div>
    <!-- 使用者管理清單 -->
    <div v-if="activeTab === 'users'" class="premium-card">
      <div v-if="loadingUsers" style="text-align:center; padding: 2rem;">載入中...</div>
      
      <table class="list-table" v-else>
        <thead class="list-header">
          <tr>
            <th>ID</th>
            <th>帳號</th>
            <th>姓名</th>
            <th>Email</th>
            <th>目前點數</th>
            <th>權限</th>
            <th style="width: 120px;">操作</th>
          </tr>
        </thead>
        <tbody class="list-body">
          <tr class="list-row" v-for="user in users" :key="user.id">
            <td>{{ user.id }}</td>
            <td>{{ user.username }}</td>
            <td>{{ user.name || '-' }}</td>
            <td>{{ user.email || '-' }}</td>
            <td style="font-weight: 700; color: #fbbf24;">{{ user.points }}</td>
            <td>
              <span class="status-badge" :class="{ active: user.role === 'admin' }">{{ user.role }}</span>
            </td>
            <td>
              <div class="user-actions">
                <button class="btn btn-outline edit-btn" @click="openUserEdit(user)">編輯</button>
                <button class="btn btn-outline delete-btn" @click="deleteUser(user)">刪除</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 使用者編輯 Modal -->
    <div class="modal-overlay" v-if="showUserModal" @click.self="showUserModal = false">
      <div class="premium-card modal-content">
        <h2>編輯使用者：{{ selectedUser?.username }}</h2>
        
        <form @submit.prevent="submitUserEdit" class="admin-form" style="margin-top: 1.5rem;">
          <div class="form-group">
            <label>真實姓名</label>
            <input type="text" v-model="editUserForm.name" class="input-form" />
          </div>
          
          <div class="form-group">
            <label>信箱</label>
            <input type="email" v-model="editUserForm.email" class="input-form" />
          </div>
          
          <div class="form-group">
            <label>剩餘積分</label>
            <input type="number" v-model.number="editUserForm.points" class="input-form" />
          </div>

          <div class="form-group">
            <label>系統權限</label>
            <select v-model="editUserForm.role" class="input-form">
              <option value="user">一般員工 (user)</option>
              <option value="admin">管理員 (admin)</option>
            </select>
          </div>

          <div class="modal-actions">
            <button type="button" class="btn btn-outline" @click="showUserModal = false">取消</button>
            <button type="submit" class="btn btn-primary">儲存變更</button>
          </div>
        </form>
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
  color: var(--danger);
}

.subtitle {
  color: var(--text-muted);
}

.tabs {
  display: flex;
  margin-bottom: 2rem;
  border-bottom: 1px solid var(--border-color);
}

.tab {
  padding: 1rem 2rem;
  background: transparent;
  border: none;
  color: var(--text-muted);
  font-weight: 600;
  font-size: 1rem;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: all 0.2s;
  font-family: var(--font-family);
}

.tab:hover {
  color: var(--text-main);
}

.tab.active {
  color: var(--primary);
  border-bottom-color: var(--primary);
}

.admin-grid {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 2rem;
}

h2 {
  font-size: 1.25rem;
  margin-bottom: 1.5rem;
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 0.5rem;
}

.admin-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-size: 0.875rem;
  color: var(--text-muted);
}

.matches-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.match-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background-color: var(--surface-hover);
  border-radius: 8px;
  border: 1px solid var(--border-color);
}

.teams {
  font-weight: 600;
  margin-bottom: 0.25rem;
}

.time {
  font-size: 0.75rem;
  color: var(--text-muted);
  margin-bottom: 0.25rem;
}

.status-badge {
  font-size: 0.75rem;
  padding: 0.1rem 0.5rem;
  border-radius: 999px;
  background: rgba(255,255,255,0.1);
}
.status-badge.active { background: rgba(16, 185, 129, 0.2); color: var(--success); }

.resolve-actions {
  display: flex;
  gap: 0.5rem;
}

.select-result {
  width: 150px;
  padding: 0.5rem;
}

.resolved-text {
  font-weight: 600;
  color: var(--text-muted);
}

/* 復用排行榜表格的樣式 */
.list-table { width: 100%; border-collapse: collapse; }
.list-header { background-color: var(--table-header-bg); border-bottom: 1px solid var(--border-color); }
.list-header th { padding: 1rem; font-weight: 600; color: var(--text-muted); font-size: 0.875rem; text-align: left; }
.list-row { border-bottom: 1px solid var(--table-row-border); transition: background-color 0.2s; color: var(--text-main); }
.list-row:hover { background-color: var(--table-row-hover); }
.list-row td { padding: 1rem; color: var(--text-main); }

.user-actions { display: flex; gap: 0.5rem; }
.user-actions .btn { padding: 0.25rem 0.5rem; font-size: 0.875rem; }
.edit-btn {  }
.delete-btn { color: var(--danger); border-color: var(--danger); }
.delete-btn:hover { background: var(--danger); color: white; }

/* Modal Styles 重用 Dashboard 的 Modal 樣式 */
.modal-overlay {
  position: fixed; top: 0; left: 0; right: 0; bottom: 0;
  background-color: var(--overlay-bg); backdrop-filter: blur(4px);
  display: flex; align-items: center; justify-content: center; z-index: 1000;
}
.modal-content { width: 90%; max-width: 500px; padding: 2rem; }
.modal-content h2 { margin-top: 0; font-size: 1.25rem; border-bottom: 1px solid var(--border-color); padding-bottom: 1rem; }
.modal-actions { display: flex; justify-content: flex-end; gap: 1rem; margin-top: 2rem; }
</style>
