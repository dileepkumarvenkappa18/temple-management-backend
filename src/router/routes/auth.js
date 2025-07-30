import LoginView from '@/views/auth/LoginView.vue'
import RegisterView from '@/views/auth/RegisterView.vue'
import ForgotPasswordView from '@/views/auth/ForgotPasswordView.vue'

export default [
  {
    path: 'login',
    name: 'Login',
    component: LoginView,
    meta: {
      title: 'Login',
      requiresAuth: false
    }
  },
  {
    path: 'register',
    name: 'Register',
    component: RegisterView,
    meta: {
      title: 'Register',
      requiresAuth: false
    }
  },
  {
    path: 'forgot-password',
    name: 'ForgotPassword',
    component: ForgotPasswordView,
    meta: {
      title: 'Forgot Password',
      requiresAuth: false
    }
  }
]