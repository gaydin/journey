import { createRouter, createWebHistory } from 'vue-router'
import useAuthStore from '@/stores/auth.store'
import LoginView from '@/views/LoginView.vue'
import CreateView from '@/views/CreateView.vue'
import SettingsView from '@/views/SettingsView.vue'
import ContentView from '@/views/ContentView.vue'
import SettingsUserView from '@/views/SettingsUserView.vue'
import EditView from '@/views/EditView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: ContentView
    },
    {
      path: '/create',
      component: CreateView
    },
    {
      path: '/edit/:id',
      component: EditView
    },
    {
      path: '/settings',
      component: SettingsView
    },
    {
      path: '/settings/user',
      component: SettingsUserView
    },
    {
      path: '/login',
      component: LoginView
    }
  ]
})

router.beforeEach(async (to) => {
  // redirect to login page if not logged in and trying to access a restricted page
  const publicPages = ['/login']
  const authRequired = !publicPages.includes(to.path)
  const auth = useAuthStore()

  if (authRequired && !auth.isAuthenticated) {
    auth.returnUrl = to.fullPath
    return '/login'
  }
})
export default router
