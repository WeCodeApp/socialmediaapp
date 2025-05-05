import { createRouter, createWebHistory } from 'vue-router'
import LoginForm from '../views/LoginForm.vue'
import SignupForm from '../views/SignupForm.vue'
import DashboardForm from '../views/DashboardForm.vue' // ✅ Corrected import

const routes = [
  {
    path: '/',
    name: 'login',
    component: LoginForm,
  },
  {
    path: '/signup',
    name: 'signup',
    component: SignupForm,
  },
  {
    path: '/dashboard',
    name: 'dashboard', // ✅ Consistent name with the file
    component: DashboardForm,
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

export default router
