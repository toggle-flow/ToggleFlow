import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { authApi } from '@/api/auth'

let initialized: boolean | null = null

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/setup', name: 'setup', component: () => import('@/views/SetupView.vue') },
    { path: '/login', name: 'login', component: () => import('@/views/LoginView.vue') },
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomeView.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach(async (to) => {
  // One-time check: does the backend have any users yet?
  if (initialized === null) {
    try {
      const status = await authApi.status()
      initialized = status.initialized
    } catch {
      initialized = false
    }
  }

  // No superuser yet → force setup
  if (!initialized && to.name !== 'setup') return { name: 'setup' }

  // Already initialized → setup is locked
  if (initialized && to.name === 'setup') return { name: 'login' }

  const auth = useAuthStore()

  // Route requires auth — verify or restore session
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    if (auth.token) {
      try {
        const user = await authApi.me()
        auth.setAuth(auth.token, user)
      } catch {
        auth.logout()
        return { name: 'login' }
      }
    } else {
      return { name: 'login' }
    }
  }

  // Already logged in → bypass login page
  if (to.name === 'login' && auth.isAuthenticated) return { name: 'home' }
})

export default router
