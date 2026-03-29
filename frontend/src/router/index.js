import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            name: 'home',
            component: () => import('@/views/Dashboard.vue')
        },
        {
            path: '/login',
            name: 'login',
            component: () => import('@/views/Login.vue')
        },
        {
            path: '/my-bets',
            name: 'myBets',
            component: () => import('@/views/MyBets.vue')
        },
        {
            path: '/leaderboard',
            name: 'leaderboard',
            component: () => import('@/views/Leaderboard.vue')
        },
        {
            path: '/admin',
            name: 'admin',
            component: () => import('@/views/Admin.vue')
        }
    ]
})

export default router
