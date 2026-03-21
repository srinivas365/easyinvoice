import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../store/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('../views/Login.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/',
      component: () => import('../components/Layout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'Dashboard',
          component: () => import('../views/Dashboard.vue')
        },
        {
          path: 'medicines',
          name: 'Medicines',
          component: () => import('../views/Medicines.vue')
        },
        {
          path: 'medicines/add',
          name: 'AddMedicine',
          component: () => import('../views/AddMedicine.vue')
        },
        {
          path: 'medicines/edit/:id',
          name: 'EditMedicine',
          component: () => import('../views/EditMedicine.vue')
        },
        {
          path: 'sales',
          name: 'Sales',
          component: () => import('../views/Sales.vue')
        },
        {
          path: 'invoices',
          name: 'Invoices',
          component: () => import('../views/Invoices.vue')
        },
        {
          path: 'purchases',
          name: 'Purchases',
          component: () => import('../views/Purchases.vue')
        },
        {
          path: 'settings',
          name: 'Settings',
          component: () => import('../views/Settings.vue')
        }
      ]
    }
  ]
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if (to.path === '/login' && authStore.isAuthenticated) {
    next('/')
  } else {
    next()
  }
})

export default router
