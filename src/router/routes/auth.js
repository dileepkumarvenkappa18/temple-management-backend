// src/router/routes/auth.js
import LoginView from '@/views/auth/LoginView.vue'
import RegisterView from '@/views/auth/RegisterView.vue'
import ForgotPasswordView from '@/views/auth/ForgotPasswordView.vue'
import ResetPasswordView from '@/views/auth/ResetPasswordView.vue'

export default [
  {
    path: 'login', // No leading slash, as it's relative to parent
    name: 'Login',
    component: LoginView,
    meta: { title: 'Login', requiresAuth: false, layout: 'AuthLayout' }
  },
  {
    path: 'register',
    name: 'Register',
    component: RegisterView,
    meta: { title: 'Register', requiresAuth: false, layout: 'AuthLayout' }
  },
  {
    path: 'forgot-password',
    name: 'ForgotPassword',
    component: ForgotPasswordView,
    meta: { title: 'Forgot Password', requiresAuth: false, layout: 'AuthLayout' }
  },
  {
    path: 'reset-password', // Relative path (no leading slash)
    name: 'ResetPassword',
    component: ResetPasswordView,
    props: route => ({ token: route.query.token }),
    meta: { title: 'Reset Password', requiresAuth: false, layout: 'AuthLayout' }
  }
]