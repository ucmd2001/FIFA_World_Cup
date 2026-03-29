import { defineStore } from 'pinia'
import api from '@/api'

export const useAuthStore = defineStore('auth', {
    state: () => ({
        user: JSON.parse(localStorage.getItem('user')) || null,
        token: localStorage.getItem('token') || null,
    }),
    getters: {
        isAuthenticated: (state) => !!state.token,
        isAdmin: (state) => state.user?.role === 'admin',
        currentPoints: (state) => state.user?.points || 0,
    },
    actions: {
        async login(username, password) {
            const res = await api.post('/auth/login', { username, password })
            this.token = res.data.token
            this.user = res.data.user
            localStorage.setItem('token', this.token)
            localStorage.setItem('user', JSON.stringify(this.user))
        },
        async register(username, name, email, password) {
            await api.post('/auth/register', { username, name, email, password })
        },
        logout() {
            this.token = null
            this.user = null
            localStorage.removeItem('token')
            localStorage.removeItem('user')
        },
        async fetchMe() {
            if (!this.token) return
            try {
                const res = await api.get('/user/me')
                this.user = res.data
                localStorage.setItem('user', JSON.stringify(this.user))
            } catch (err) {
                console.error('Failed to fetch user', err)
            }
        },
        async claimDaily() {
            const res = await api.post('/user/claim-daily')
            this.user.points = res.data.points
            localStorage.setItem('user', JSON.stringify(this.user))
            return res.data
        }
    }
})
