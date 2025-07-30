// ========================================
// src/router/routes/public.js - Phase 2
// ========================================
import LandingPage from '@/views/public/LandingPage.vue'
import TermsPage from '@/views/public/TermsPage.vue'
import PrivacyPage from '@/views/public/PrivacyPage.vue'

export default [
  {
    path: '',
    name: 'Home',
    component: LandingPage,
    meta: {
      title: 'Welcome',
      requiresAuth: false
    }
  },
  {
    path: 'terms',
    name: 'Terms',
    component: TermsPage,
    meta: {
      title: 'Terms of Service',
      requiresAuth: false
    }
  },
  {
    path: 'privacy',
    name: 'Privacy',
    component: PrivacyPage,
    meta: {
      title: 'Privacy Policy',
      requiresAuth: false
    }
  }
]
