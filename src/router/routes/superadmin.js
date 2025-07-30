// src/router/routes/superadmin.js
import SuperAdminDashboard from '@/views/superadmin/SuperAdminDashboard.vue'
import TenantApprovals from '@/components/superadmin/TenantApprovals.vue' // Direct import from components

export default [
  {
    path: 'dashboard',
    name: 'SuperAdminDashboard',
    component: SuperAdminDashboard,
    meta: {
      title: 'Super Admin Dashboard',
      breadcrumb: 'Dashboard'
    }
  },
  {
    path: 'tenant-approvals',
    name: 'TenantApprovals',
    component: TenantApprovals, // Use component directly
    meta: {
      title: 'Tenant Approvals',
      breadcrumb: 'Tenant Approvals'
    }
  }
]