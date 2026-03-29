/**
 * useCache - 前端 Memory Cache Composable
 *
 * 使用方式：
 *   const { get, set, invalidate } = useCache()
 *   const data = get('leaderboard') ?? await fetchLeaderboard()
 *   set('leaderboard', data, 30 * 60 * 1000)  // TTL 30 分鐘
 *
 * Cache 策略設計：
 *   - leaderboard      TTL 30 min  (更新很慢)
 *   - matches (賽前)   TTL 5 min   (賽前不太變)
 *   - matches (賽中)   TTL 0       (賽中不 cache，每次直接打 API)
 *   - user/bets        TTL 0       (個人資料永遠即時)
 */

// 全域 Map，跨組件共享 - Vue 組件 unmount 後資料仍保留（SPA 行為正確）
const _cacheStore = new Map()

/**
 * @param {string} key - cache 的 key
 * @returns {{ data: any, expiredAt: number } | null} - cache 紀錄或 null
 */
function _getEntry(key) {
    const entry = _cacheStore.get(key)
    if (!entry) return null
    if (Date.now() > entry.expiredAt) {
        _cacheStore.delete(key)
        return null
    }
    return entry
}

export function useCache() {
    /**
     * 取得 cache 資料，若過期或不存在則回傳 null
     * @param {string} key
     * @returns {any | null}
     */
    const get = (key) => {
        const entry = _getEntry(key)
        return entry ? entry.data : null
    }

    /**
     * 寫入 cache 資料
     * @param {string} key
     * @param {any} data
     * @param {number} ttlMs - 存活時間（毫秒），0 表示永不儲存
     */
    const set = (key, data, ttlMs) => {
        if (ttlMs <= 0) return
        _cacheStore.set(key, {
            data,
            expiredAt: Date.now() + ttlMs,
            cachedAt: Date.now(),
        })
    }

    /**
     * 主動清除某個 key 的 cache
     * @param {string} key
     */
    const invalidate = (key) => {
        _cacheStore.delete(key)
    }

    /**
     * 清除所有 cache
     */
    const invalidateAll = () => {
        _cacheStore.clear()
    }

    /**
     * 取得 cache 的剩餘時間（毫秒），若不存在則回傳 0
     * @param {string} key
     * @returns {number}
     */
    const getRemainingTtl = (key) => {
        const entry = _getEntry(key)
        if (!entry) return 0
        return Math.max(0, entry.expiredAt - Date.now())
    }

    /**
     * 取得 cache 的寫入時間，若不存在則回傳 null
     * @param {string} key
     * @returns {Date | null}
     */
    const getCachedAt = (key) => {
        const entry = _getEntry(key)
        if (!entry) return null
        return new Date(entry.cachedAt)
    }

    return {
        get,
        set,
        invalidate,
        invalidateAll,
        getRemainingTtl,
        getCachedAt,
    }
}

// 統一的 Cache Key 與 TTL 常數，方便全專案一致使用
export const CACHE_KEYS = {
    LEADERBOARD: 'leaderboard',
    MATCHES: 'matches',
}

export const CACHE_TTL = {
    LEADERBOARD: 30 * 60 * 1000,     // 30 分鐘
    MATCHES_PENDING: 5 * 60 * 1000,  // 5 分鐘（賽前）
    MATCHES_ACTIVE: 0,                // 不 cache（賽中）
    NO_CACHE: 0,
}
