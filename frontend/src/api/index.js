import axios from 'axios'

const api = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api',
    timeout: 10000,
})

// 請求攔截器：自動帶上 JWT Token
api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token')
        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`
        }
        return config
    },
    (error) => Promise.reject(error)
)

// 回應攔截器：處理資料包裝與 401 未授權等情況
api.interceptors.response.use(
    (response) => {
        // 自動解除後端的標準化回傳外殼 { code, message, data }
        if (response.data && response.data.code !== undefined && response.data.data !== undefined) {
            // 將內層的實際資料提升為 response.data
            response.data = response.data.data
        }
        return response
    },
    (error) => {
        if (error.response && error.response.status === 401) {
            localStorage.removeItem('token')
            localStorage.removeItem('user')
            // 如果不是在登入頁，就重導向回登入
            if (window.location.pathname !== '/login') {
                window.location.href = '/login'
            }
        }
        return Promise.reject(error)
    }
)

export default api
