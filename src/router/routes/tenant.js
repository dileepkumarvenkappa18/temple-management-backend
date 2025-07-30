import TenantDashboard from '@/views/tenant/TenantDashboard.vue'
import CreateTemple from '@/views/tenant/CreateTemple.vue'
import ManageTemples from '@/views/tenant/ManageTemples.vue'
import TempleDetails from '@/views/tenant/TempleDetails.vue'
import EditTemple from '@/views/tenant/EditTemple.vue'
import EntityDashboard from '@/views/entity/EntityDashboard.vue' // Import the EntityDashboard component

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
  // NEW ROUTE: Tenant-specific dashboard with tenantId in path
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
  // Standard routes (no changes needed)
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
  }
]