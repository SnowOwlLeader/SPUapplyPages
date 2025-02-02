import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'login',
      component: LoginView
    },
    {
      path: '/oauth/callback',
      name: 'oauth-callback',
      component: () => import('../views/OAuthCallback.vue')
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/RegisterView.vue')
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const isAuthenticated = localStorage.getItem('token')
  const userInfo = localStorage.getItem('user_info')

  
  // 其他页面的检查
  if (to.name !== 'login' && to.name !== 'oauth-callback' && !isAuthenticated) {
    next({ name: 'login' })
  } else {
    next()
  }
})

export default router 