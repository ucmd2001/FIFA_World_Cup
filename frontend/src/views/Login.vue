<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const isLoginMode = ref(true)
const username = ref('')
const name = ref('')
const email = ref('')
const password = ref('')
const errorMsg = ref('')

const handleSubmit = async () => {
  errorMsg.value = ''
  if (!username.value || !password.value) {
    errorMsg.value = '請填寫必填欄位'
    return
  }

  try {
    if (isLoginMode.value) {
      await authStore.login(username.value, password.value)
    } else {
      if (!name.value || !email.value) {
        errorMsg.value = '註冊需填寫姓名與 Email'
        return
      }
      await authStore.register(username.value, name.value, email.value, password.value)
      // 註冊後可自動登入或顯示成功訊息
      await authStore.login(username.value, password.value)
    }
    // 成功後回到大廳且頁面重整觸發自動領點
    window.location.href = '/'
  } catch (err) {
    if (err.response?.data?.error) {
      errorMsg.value = err.response.data.error
    } else {
      errorMsg.value = '發生未知的錯誤'
    }
  }
}
</script>

<template>
  <div class="form-container">
    <div class="premium-card auth-card">
      <h2 class="auth-title">{{ isLoginMode ? '進入 GoalBet' : '註冊員工帳號' }}</h2>
      
      <div v-if="errorMsg" class="error-alert">{{ errorMsg }}</div>
      
      <form class="auth-form" @submit.prevent="handleSubmit">
        <div class="form-group">
          <label>帳號 (登入用)</label>
          <input type="text" v-model="username" class="input-form" placeholder="請輸入帳號" required />
        </div>
        
        <template v-if="!isLoginMode">
          <div class="form-group">
            <label>真實姓名</label>
            <input type="text" v-model="name" class="input-form" placeholder="例如：王小明" />
          </div>
          <div class="form-group">
            <label>Email 信箱</label>
            <input type="email" v-model="email" class="input-form" placeholder="請填寫常用信箱" />
          </div>
        </template>

        <div class="form-group">
          <label>密碼</label>
          <input type="password" v-model="password" class="input-form" placeholder="請輸入密碼" required minlength="6"/>
          <small v-if="!isLoginMode" class="hint">第一位註冊的使用者將自動獲得「管理員」權限。</small>
        </div>
        
        <button type="submit" class="btn btn-primary login-btn">
          {{ isLoginMode ? '登入' : '註冊並進入' }}
        </button>
      </form>
      
      <p class="auth-footer">
        {{ isLoginMode ? '還沒有帳號？' : '已經有帳號了？' }}
        <a href="#" class="link" @click.prevent="isLoginMode = !isLoginMode">
          {{ isLoginMode ? '立即註冊' : '我想登入' }}
        </a>
      </p>
    </div>
  </div>
</template>

<style scoped>
.form-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
}

.auth-card {
  width: 100%;
  max-width: 400px;
  padding: 2.5rem 2rem;
}

.auth-title {
  text-align: center;
  margin-bottom: 2rem;
  font-size: 1.5rem;
  background: linear-gradient(135deg, var(--primary), var(--accent));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  font-weight: 700;
}

.error-alert {
  background-color: rgba(239, 68, 68, 0.1);
  color: var(--danger);
  padding: 0.75rem;
  border-radius: 8px;
  margin-bottom: 1rem;
  font-size: 0.875rem;
  text-align: center;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
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

.hint {
  font-size: 0.75rem;
  color: var(--warning);
  margin-top: 0.25rem;
}

.login-btn {
  margin-top: 1rem;
  width: 100%;
  font-size: 1rem;
  padding: 0.875rem;
}

.auth-footer {
  margin-top: 2rem;
  text-align: center;
  font-size: 0.875rem;
  color: var(--text-muted);
}

.link {
  color: var(--primary);
  text-decoration: none;
  font-weight: 500;
}

.link:hover {
  text-decoration: underline;
}
</style>
