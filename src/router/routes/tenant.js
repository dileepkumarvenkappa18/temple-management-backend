import TenantDashboard from '@/views/tenant/TenantDashboard.vue'
import CreateTemple from '@/views/tenant/CreateTemple.vue'
import ManageTemples from '@/views/tenant/ManageTemples.vue'
import TempleDetails from '@/views/tenant/TempleDetails.vue'
import EditTemple from '@/views/tenant/EditTemple.vue'
import TempleRegisterReport from '@/views/tenant/reports/TempleRegisterReport.vue'
import TempleActivitiesReport from '@/views/tenant/reports/TempleActivitiesReport.vue'
import BirthdaysReport from '@/views/tenant/reports/BirthdaysReport.vue'

export default [
  // Standard tenant dashboard route (kept for backward compatibility)
  {
    path: 'dashboard',
    name: 'TenantDashboard',
    component: TenantDashboard,
    meta: {
      title: 'Tenant Dashboard',
      requiresAuth: true,
      role: 'tenant'
    }
  },
  // Tenant-specific dashboard with tenantId in path
  {
    path: ':tenantId/dashboard',
    name: 'TenantSpecificDashboard',
    component: TenantDashboard,
    meta: {
      title: 'Tenant Dashboard',
      requiresAuth: true,
      role: 'tenant'
    }
  },
  // Standard routes
  {
    path: 'entities/create',
    name: 'CreateTemple',
    component: CreateTemple,
    meta: {
      title: 'Create Temple',
      requiresAuth: true,
      role: 'tenant'
    }
  },
  {
    path: 'entities',
    name: 'ManageTemples',
    component: ManageTemples,
    meta: {
      title: 'Manage Temples',
      requiresAuth: true,
      role: 'tenant'
    }
  },
  {
    path: 'entities/:id',
    name: 'TempleDetails',
    component: TempleDetails,
    meta: {
      title: 'Temple Details',
      requiresAuth: true,
      role: 'tenant'
    }
  },
  {
    path: 'entities/:id/edit',
    name: 'EditTemple',
    component: EditTemple,
    meta: {
      title: 'Edit Temple',
      requiresAuth: true,
      role: 'tenant'
    }
  },
  // NEW REPORT ROUTES
  {
    path: 'reports/temple-register',
    name: 'TempleRegisterReport',
    component: TempleRegisterReport,
    meta: {
      title: 'Temple Register Report',
      requiresAuth: true,
      role: 'tenant'
    }
  },
  {
    path: 'reports/temple-activities',
    name: 'TempleActivitiesReport',
    component: TempleActivitiesReport,
    meta: {
      title: 'Temple Activities Report',
      requiresAuth: true,
      role: 'tenant'
    }
  },
  {
    path: 'reports/birthdays',
    name: 'BirthdaysReport',
    component: BirthdaysReport,
    meta: {
      title: 'Birthdays Report',
      requiresAuth: true,
      role: 'tenant'
    }
  }
]