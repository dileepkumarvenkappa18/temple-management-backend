// src/router/index.js - UPDATED VERSION WITH AUDIT LOGS REPORT
import { createRouter, createWebHistory } from 'vue-router'
import { checkProfileCompleted, requireAuth } from './guards'
import superadminRoutes from './routes/superadmin'
import MessageComposer from '@/components/communication/MessageComposer.vue'
import { useAuthStore } from '@/stores/auth'
import { setupRouteGuards } from './guards';

// Layouts
import PublicLayout from '@/layouts/PublicLayout.vue'
import AuthLayout from '@/layouts/AuthLayout.vue'
import DashboardLayout from '@/layouts/DashboardLayout.vue'

// Public Views
import LandingPage from '@/views/public/LandingPage.vue'
import TermsPage from '@/views/public/TermsPage.vue'
import PrivacyPage from '@/views/public/PrivacyPage.vue'

// Auth Views
import LoginView from '@/views/auth/LoginView.vue'
import RegisterView from '@/views/auth/RegisterView.vue'
import ForgotPasswordView from '@/views/auth/ForgotPasswordView.vue'
import ResetPasswordView from '@/views/auth/ResetPasswordView.vue'

// Tenant Views
import TenantDashboard from '@/views/tenant/TenantDashboard.vue'
import CreateTemple from '@/views/tenant/CreateTemple.vue'
import ManageTemples from '@/views/tenant/ManageTemples.vue'
import EditTemple from '@/views/tenant/EditTemple.vue'
import TempleDetails from '@/views/tenant/TempleDetails.vue'
import TenantSelectionView from '@/views/tenant/TenantSelectionView.vue'

// Report Views
import TempleRegisterReport from '@/views/tenant/reports/TempleRegisterReport.vue'
import TempleActivitiesReport from '@/views/tenant/reports/TempleActivitiesReport.vue'
import BirthdaysReport from '@/views/tenant/reports/BirthdaysReport.vue'
import AuditLogsReport from '@/views/tenant/reports/AuditLogsView.vue'

// Entity (Temple Admin) Views - Phase 5
import EntityDashboard from '@/views/entity/EntityDashboard.vue'
import DevoteeManagement from '@/views/entity/DevoteeManagement.vue'
import SevaManagement from '@/views/entity/SevaManagement.vue'
import DonationManagement from '@/views/entity/DonationManagement.vue'
import EventManagement from '@/views/entity/EventManagement.vue'
import CommunicationCenter from '@/views/entity/CommunicationCenter.vue'
import VolunteerManagement from '@/views/entity/VolunteerManagement.vue'

// Devotee Views - Phase 6
import TempleSelection from '@/views/devotee/TempleSelection.vue'
import ProfileCreation from '@/views/devotee/ProfileCreation.vue'
import DevoteeDashboard from '@/views/devotee/DevoteeDashboard.vue'
import SevaBooking from '@/views/devotee/SevaBooking.vue'
import DonationHistory from '@/views/devotee/DonationHistory.vue'
import MyEvents from '@/views/devotee/MyEvents.vue'
import ProfileEdit from '@/views/devotee/ProfileEdit.vue'

// Volunteer Views - Phase 7
import VolunteerTempleSelection from '@/views/volunteer/TempleSelection.vue'
import VolunteerDashboard from '@/views/volunteer/VolunteerDashboard.vue'
import MyAssignments from '@/views/volunteer/MyAssignments.vue'
import MySchedule from '@/views/volunteer/MySchedule.vue'

// SuperAdmin Views
import SuperAdminDashboard from '@/views/superadmin/SuperAdminDashboard.vue'

// User Views
import NotificationsView from '@/views/user/NotificationsView.vue'
import UserProfileView from '@/views/user/UserProfileView.vue'
import SettingsView from '@/views/settings/SettingsView.vue'
import SupportView from '@/views/user/SupportView.vue'

// Error Views
import NotFound from '@/views/error/NotFound.vue'
import Unauthorized from '@/views/error/Unauthorized.vue'
import ServerError from '@/views/error/ServerError.vue'

const routes = [
  // Public Routes
  {
    path: '/',
    component: PublicLayout,
    children: [
      {
        path: '',
        name: 'Landing',
        component: LandingPage,
        meta: { 
          title: 'Temple Management SaaS - Manage Your Temple Digitally',
          description: 'Manage your temple, devotees, sevas & donations with ease'
        }
      },
      {
        path: '/terms',
        name: 'Terms',
        component: TermsPage,
        meta: { title: 'Terms of Service' }
      },
      {
        path: '/privacy',
        name: 'Privacy',
        component: PrivacyPage,
        meta: { title: 'Privacy Policy' }
      }
    ]
  },

  // Auth Routes - Direct paths
  {
    path: '/login',
    name: 'Login',
    component: LoginView,
    meta: { 
      title: 'Login',
      layout: 'AuthLayout' 
    }
  },
  {
    path: '/register',
    name: 'Register',
    component: RegisterView,
    meta: { 
      title: 'Register',
      layout: 'AuthLayout' 
    }
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: ForgotPasswordView,
    meta: { 
      title: 'Forgot Password',
      layout: 'AuthLayout' 
    }
  },

  // Reset Password Route
  {
    path: '/auth-pages/reset-password',
    name: 'BackendResetPassword',
    component: ResetPasswordView,
    props: route => ({ token: route.query.token }),
    meta: { 
      title: 'Reset Password',
      layout: 'AuthLayout',
      requiresAuth: false
    }
  },

  // Debug Reset Password Route - For testing
  {
    path: '/debug-reset',
    name: 'DebugReset',
    component: ResetPasswordView,
    meta: { 
      title: 'Debug Reset Password',
      layout: 'AuthLayout',
      requiresAuth: false
    }
  },

  // User Routes (Common for all authenticated users)
  {
    path: '/notifications',
    name: 'Notifications',
    component: NotificationsView,
    meta: { 
      requiresAuth: true,
      title: 'Notifications',
      breadcrumb: 'Notifications',
      layout: 'DashboardLayout'
    }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: UserProfileView, 
    meta: { 
      requiresAuth: true,
      title: 'My Profile',
      breadcrumb: 'My Profile',
      layout: 'DashboardLayout'
    }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: SettingsView,
    meta: {
      requiresAuth: true,
      title: 'Account Settings',
      breadcrumb: 'Settings',
      layout: 'DashboardLayout'
    }
  },
  {
    path: '/support',
    name: 'Support',
    component: SupportView,
    meta: {
      requiresAuth: true,
      title: 'Support Center',
      breadcrumb: 'Support',
      layout: 'DashboardLayout'
    }
  },

  // ✅ CRITICAL FIX: Tenant Selection route - MOVED TO TOP PRIORITY
  {
    path: '/tenant-selection',
    name: 'TenantSelection',
    component: TenantSelectionView,
    meta: { 
      requiresAuth: true,
      allowedRoles: ['superadmin', 'super_admin', 'standard_user', 'standarduser', 'monitoring_user', 'monitoringuser'],
      title: 'Temple Selection',
      layout: 'DashboardLayout'
    }
  },

  // Tenant Routes
  {
    path: '/tenant',
    component: DashboardLayout,
    meta: { 
      requiresAuth: true,
      allowedRoles: ['tenant', 'templeadmin']
    },
    children: [
      {
        path: 'dashboard',
        name: 'TenantDashboard',
        component: TenantDashboard,
        meta: { 
          title: 'Tenant Dashboard',
          breadcrumb: 'Dashboard'
        }
      },
      {
        path: 'entities/create',
        name: 'CreateTemple',
        component: CreateTemple,
        meta: { 
          title: 'Create Temple',
          breadcrumb: 'Create Temple'
        }
      },
      {
        path: 'entities/:id',
        name: 'TempleDetails',
        component: TempleDetails,
        props: true,
        meta: { 
          title: 'Temple Details',
          breadcrumb: 'Temple Details'
        }
      },
      {
        path: 'entities/:id/edit',
        name: 'EditTemple',
        component: EditTemple,
        props: true,
        meta: { 
          title: 'Edit Temple',
          breadcrumb: 'Edit Temple'
        }
      },
      // Tenant-specific dashboard route
      {
        path: ':tenantId/dashboard',
        name: 'TenantSpecificDashboard',
        component: TenantDashboard,
        props: true,
        meta: { 
          title: 'Tenant Dashboard',
          breadcrumb: 'Dashboard'
        }
      },
      // Report routes
      {
        path: 'reports/temple-register',
        name: 'TempleRegisterReport',
        component: TempleRegisterReport,
        meta: {
          title: 'Temple Register Report',
          breadcrumb: 'Temple Register',
          requiresAuth: true,
          allowedRoles: ['tenant', 'templeadmin']
        }
      },
      {
        path: 'reports/temple-activities',
        name: 'TempleActivitiesReport',
        component: TempleActivitiesReport,
        meta: {
          title: 'Temple Activities Report',
          breadcrumb: 'Temple Activities',
          requiresAuth: true,
          allowedRoles: ['tenant', 'templeadmin']
        }
      },
      {
        path: 'reports/birthdays',
        name: 'BirthdaysReport',
        component: BirthdaysReport,
        meta: {
          title: 'Birthdays Report',
          breadcrumb: 'Birthdays',
          requiresAuth: true,
          allowedRoles: ['tenant', 'templeadmin']
        }
      },
      // NEW: Audit Logs Report Route
      {
        path: 'reports/audit-logs',
        name: 'AuditLogsReport',
        component: AuditLogsReport,
        meta: {
          title: 'Audit Logs Report',
          breadcrumb: 'Audit Logs',
          requiresAuth: true,
          allowedRoles: ['tenant', 'templeadmin']
        }
      }
    ]
  },

  // Direct route for tenant-specific dashboard
  {
    path: '/tenant/:tenantId/dashboard',
    name: 'TenantDirectDashboard',
    component: TenantDashboard,
    props: true,
    meta: { 
      requiresAuth: true,
      allowedRoles: ['tenant', 'templeadmin'],
      title: 'Tenant Dashboard',
      breadcrumb: 'Dashboard',
      layout: 'DashboardLayout'
    }
  },

  // Entity (Temple Admin) Routes - Phase 5
  {
    path: '/entity/:id',
    component: DashboardLayout,
    props: true,
    meta: { 
      requiresAuth: true,
      allowedRoles: ['tenant', 'entity_admin', 'templeadmin', 'standard_user', 'standarduser', 'monitoring_user', 'monitoringuser']
    },
    children: [
      {
        path: '/entity/:id/dashboard',
        name: 'EntityDirectDashboard',
        component: () => import('@/views/entity/EntityDashboard.vue'),
        meta: {
          title: 'Entity Dashboard',
          requiresAuth: true,
          allowedRoles: ['tenant', 'entity_admin', 'templeadmin', 'standard_user', 'standarduser', 'monitoring_user', 'monitoringuser']
        }
      },
      {
        path: 'dashboard',
        name: 'EntityDashboard',
        component: EntityDashboard,
        props: true,
        meta: { 
          title: 'Temple Dashboard',
          breadcrumb: 'Dashboard',
          allowedRoles: ['tenant', 'entity_admin', 'templeadmin', 'standard_user', 'standarduser', 'monitoring_user', 'monitoringuser']
        }
      },
      {
        path: 'devotees',
        name: 'DevoteeManagement',
        component: DevoteeManagement,
        props: true,
        meta: { 
          title: 'Devotee Management',
          breadcrumb: 'Devotees',
          allowedRoles: ['tenant', 'entity_admin', 'templeadmin', 'standard_user', 'standarduser', 'monitoring_user', 'monitoringuser']
        }
      },
      {
        path: 'sevas',
        name: 'SevaManagement',
        component: SevaManagement,
        props: true,
        meta: { 
          title: 'Seva Management',
          breadcrumb: 'Sevas',
          allowedRoles: ['tenant', 'entity_admin', 'templeadmin', 'standard_user', 'standarduser']
        }
      },
      {
        path: 'donations',
        name: 'DonationManagement',
        component: DonationManagement,
        props: true,
        meta: { 
          title: 'Donation Management',
          breadcrumb: 'Donations',
          allowedRoles: ['tenant', 'entity_admin', 'templeadmin', 'standard_user', 'standarduser']
        }
      },
      {
        path: 'events',
        name: 'EventManagement',
        component: EventManagement,
        props: true,
        meta: { 
          title: 'Event Management',
          breadcrumb: 'Events',
          allowedRoles: ['tenant', 'entity_admin', 'templeadmin', 'standard_user', 'standarduser']
        }
      },
      {
        path: 'communication',
        name: 'CommunicationCenter',
        component: CommunicationCenter,
        props: true,
        meta: { 
          title: 'Communication Center',
          breadcrumb: 'Communication',
          allowedRoles: ['tenant', 'entity_admin', 'templeadmin', 'standard_user', 'standarduser']
        }
      },
      {
        path: 'message',
        name: 'MessageComposer',
        component: MessageComposer,
        props: true,
        meta: {
          title: 'Compose Message',
          breadcrumb: 'Message Composer',
          allowedRoles: ['tenant', 'entity_admin', 'templeadmin', 'standard_user', 'standarduser']
        }
      },
      {
        path: 'volunteers',
        name: 'VolunteerManagement',
        component: VolunteerManagement,
        props: true,
        meta: { 
          title: 'Volunteer Management',
          breadcrumb: 'Volunteers',
          allowedRoles: ['tenant', 'entity_admin', 'templeadmin', 'standard_user', 'standarduser']
        }
      }
    ]
  },

  // Devotee Routes - Phase 6
  {
    path: '/devotee',
    component: DashboardLayout,
    meta: { 
      requiresAuth: true,
      allowedRoles: ['devotee']
    },
    children: [
      {
        path: 'temple-selection',
        name: 'DevoteeTempleSelection',
        component: TempleSelection,
        meta: { 
          title: 'Select Temple',
          breadcrumb: 'Temple Selection'
        }
      },
      {
        path: 'profile/create',
        name: 'DevoteeProfileCreation',
        component: ProfileCreation,
        meta: { 
          title: 'Complete Profile',
          breadcrumb: 'Profile Creation'
        }
      },
      {
        path: 'profile/edit',
        name: 'DevoteeProfileEdit',
        component: ProfileEdit,
        meta: { 
          title: 'Edit Profile',
          breadcrumb: 'Edit Profile',
          role: 'devotee'
        }
      }
    ]
  },

  // Devotee Entity-Specific Routes
  {
    path: '/entity/:id/devotee',
    component: DashboardLayout,
    props: true,
    meta: { 
      requiresAuth: true,
      allowedRoles: ['devotee']
    },
    children: [
      {
        path: 'profile/create',
        name: 'DevoteeEntityProfileCreation',
        component: ProfileCreation,
        props: true,
        meta: { 
          title: 'Complete Profile',
          breadcrumb: 'Profile Creation'
        }
      },
      {
        path: 'dashboard',
        name: 'DevoteeDashboard',
        component: DevoteeDashboard,
        props: true,
        meta: { 
          title: 'My Dashboard',
          breadcrumb: 'Dashboard'
        }
      },
      {
        path: 'profile/edit',
        name: 'DevoteeProfileEdit',
        component: ProfileEdit,
        props: true,
        meta: { 
          title: 'Edit Profile',
          breadcrumb: 'Edit Profile'
        }
      },
      {
        path: 'seva-booking',
        name: 'SevaBooking',
        component: SevaBooking,
        props: true,
        meta: { 
          title: 'Book Seva',
          breadcrumb: 'Seva Booking'
        }
      },
      {
        path: 'donations',
        name: 'DevoteeDonationHistory',
        component: DonationHistory,
        props: true,
        meta: { 
          title: 'My Donations',
          breadcrumb: 'Donation History'
        }
      },
      {
        path: 'events',
        name: 'DevoteeMyEvents',
        component: MyEvents,
        props: true,
        meta: { 
          title: 'Temple Events',
          breadcrumb: 'My Events'
        }
      },
      {
        path: 'my-seva-bookings',
        name: 'DevoteeMySevaBookings',
        component: () => import('@/views/devotee/MySevaBookings.vue'),
        props: true,
        meta: { 
          title: 'My Seva Bookings',
          breadcrumb: 'My Seva Bookings'
        }
      }
    ]
  },

  // Volunteer Routes - Phase 7
  {
    path: '/volunteer',
    component: DashboardLayout,
    meta: { 
      requiresAuth: true,
      allowedRoles: ['volunteer']
    },
    children: [
      {
        path: 'temple-selection',
        name: 'VolunteerTempleSelection',
        component: VolunteerTempleSelection,
        meta: { 
          title: 'Select Temple',
          breadcrumb: 'Temple Selection'
        }
      }
    ]
  },

  // Volunteer Entity-Specific Routes - Phase 7
  {
    path: '/entity/:id/volunteer',
    component: DashboardLayout,
    props: true,
    meta: { 
      requiresAuth: true,
      allowedRoles: ['volunteer']
    },
    children: [
      {
        path: 'dashboard',
        name: 'VolunteerDashboard',
        component: VolunteerDashboard,
        props: true,
        meta: { 
          title: 'Volunteer Dashboard',
          breadcrumb: 'Dashboard'
        }
      },
      {
        path: 'assignments',
        name: 'MyAssignments',
        component: MyAssignments,
        props: true,
        meta: { 
          title: 'My Assignments',
          breadcrumb: 'Assignments'
        }
      },
      {
        path: 'schedule',
        name: 'MySchedule',
        component: MySchedule,
        props: true,
        meta: { 
          title: 'My Schedule',
          breadcrumb: 'Schedule'
        }
      }
    ]
  },

  // SuperAdmin Routes
  {
    path: '/superadmin',
    component: DashboardLayout,
    meta: { 
      requiresAuth: true,
      allowedRoles: ['superadmin', 'super_admin']
    },
    children: superadminRoutes
  },

  // Convenience Redirects
  {
    path: '/temple-selection',
    redirect: to => {
      return '/devotee/temple-selection'
    }
  },

  // Error Routes
  {
    path: '/unauthorized',
    name: 'Unauthorized',
    component: Unauthorized,
    meta: { title: 'Access Denied' }
  },
  {
    path: '/server-error',
    name: 'ServerError',
    component: ServerError,
    meta: { title: 'Server Error' }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFound,
    meta: { title: 'Page Not Found' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// ✅ FIXED: Simplified and more robust route guards
router.beforeEach((to, from, next) => {
  // Set page title
  if (to.meta.title) {
    document.title = to.meta.title
  }

  console.log('🚀 Route navigation:', {
    to: to.path,
    name: to.name,
    requiresAuth: to.meta.requiresAuth,
    allowedRoles: to.meta.allowedRoles
  })

  // Check if route requires authentication
  if (to.meta.requiresAuth) {
    return requireAuth(to, from, next)
  }
  
  next()
})

// ✅ MOST IMPORTANT FIX: Special redirect guard for standard/monitoring users
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  // Only run for authenticated users
  if (!authStore.isAuthenticated || !authStore.user) {
    next()
    return
  }

  const userRole = (authStore.user.role || '').toLowerCase().trim()
  
  console.log('🛡️ SPECIAL ROLE CHECK:', {
    userRole,
    currentPath: to.path,
    needsTenantSelection: authStore.needsTenantSelection
  })

  // Force standard_user and monitoring_user to tenant selection
  if (authStore.needsTenantSelection && to.path !== '/tenant-selection') {
    console.log('⚠️ REDIRECTING special user to tenant selection')
    next('/tenant-selection')
    return
  }

  next()
})

// Route helpers for components
export function useRouteHelpers() {
  return {
    getTenantRoutes: () => [
      { name: 'TenantDashboard', path: '/tenant/dashboard', label: 'Dashboard' },
      { name: 'ManageTemples', path: '/tenant/entities', label: 'Manage Temples' },
      { name: 'CreateTemple', path: '/tenant/entities/create', label: 'Create Temple' }
    ],
    
    getEntityRoutes: (entityId) => [
      { name: 'EntityDashboard', path: `/entity/${entityId}/dashboard`, label: 'Dashboard' },
      { name: 'DevoteeManagement', path: `/entity/${entityId}/devotees`, label: 'Devotees' },
      { name: 'SevaManagement', path: `/entity/${entityId}/sevas`, label: 'Sevas' },
      { name: 'DonationManagement', path: `/entity/${entityId}/donations`, label: 'Donations' },
      { name: 'EventManagement', path: `/entity/${entityId}/events`, label: 'Events' },
      { name: 'CommunicationCenter', path: `/entity/${entityId}/communication`, label: 'Communication' },
      { name: 'VolunteerManagement', path: `/entity/${entityId}/volunteers`, label: 'Volunteers' }
    ],

    getDevoteeRoutes: (entityId) => [
      { 
        name: 'DevoteeDashboard', 
        path: `/entity/${entityId}/devotee/dashboard`, 
        label: 'Dashboard',
        params: { id: entityId }
      },
      { 
        name: 'DevoteeProfileEdit', 
        path: `/entity/${entityId}/devotee/profile/edit`, 
        label: 'Edit Profile',
        params: { id: entityId }
      },
      { 
        name: 'SevaBooking', 
        path: `/entity/${entityId}/devotee/seva-booking`, 
        label: 'Book Seva',
        params: { id: entityId }
      },
      { 
        name: 'DevoteeDonationHistory', 
        path: `/entity/${entityId}/devotee/donations`, 
        label: 'My Donations',
        params: { id: entityId }
      },
      { 
        name: 'DevoteeMyEvents', 
        path: `/entity/${entityId}/devotee/events`, 
        label: 'Temple Events',
        params: { id: entityId }
      }
    ],

    getVolunteerRoutes: (entityId) => [
      { 
        name: 'VolunteerDashboard', 
        path: `/entity/${entityId}/volunteer/dashboard`, 
        label: 'Dashboard',
        params: { id: entityId }
      },
      { 
        name: 'MyAssignments', 
        path: `/entity/${entityId}/volunteer/assignments`, 
        label: 'My Assignments',
        params: { id: entityId }
      },
      { 
        name: 'MySchedule', 
        path: `/entity/${entityId}/volunteer/schedule`, 
        label: 'My Schedule',
        params: { id: entityId }
      }
    ],
    
    getPublicRoutes: () => [
      { name: 'Landing', path: '/', label: 'Home' },
      { name: 'Terms', path: '/terms', label: 'Terms' },
      { name: 'Privacy', path: '/privacy', label: 'Privacy' }
    ],

    getAuthRoutes: () => [
      { name: 'Login', path: '/login', label: 'Login' },
      { name: 'Register', path: '/register', label: 'Register' },
      { name: 'ForgotPassword', path: '/forgot-password', label: 'Forgot Password' }
    ],
    
    getCommonRoutes: () => [
      { name: 'Notifications', path: '/notifications', label: 'Notifications' },
      { name: 'Profile', path: '/profile', label: 'My Profile' },
      { name: 'Settings', path: '/settings', label: 'Settings' },
      { name: 'Support', path: '/support', label: 'Support' }
    ],

    // ✅ FIXED: Cleaner role-specific redirect logic
    getPostLoginRedirect: (userRole, entityId = null) => {
      const role = (userRole || '').toLowerCase().trim()
      
      // Special admin roles
      if (role === 'standard_user' || role === 'standarduser' || 
          role === 'monitoring_user' || role === 'monitoringuser') {
        return '/tenant-selection'
      }
      
      // Regular roles
      switch (role) {
        case 'superadmin':
        case 'super_admin':
          return '/superadmin/dashboard'
        case 'tenant':
        case 'templeadmin':
          const tenantId = localStorage.getItem('current_tenant_id')
          return tenantId ? `/tenant/${tenantId}/dashboard` : '/tenant/dashboard'
        case 'devotee':
          return entityId ? `/entity/${entityId}/devotee/dashboard` : '/devotee/temple-selection'
        case 'volunteer':
          return entityId ? `/entity/${entityId}/volunteer/dashboard` : '/volunteer/temple-selection'
        case 'entity_admin':
          return entityId ? `/entity/${entityId}/dashboard` : '/tenant/dashboard'
        default:
          return '/'
      }
    }
  }
}

export default router