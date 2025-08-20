// src/router/routes/shared.js
import TenantSelectionView from '@/views/tenant/TenantSelectionView.vue'

export default [
  {
    path: '/tenant-selection',
    name: 'TenantSelection',
    component: TenantSelectionView,
    meta: {
      title: 'Temple Selection',
      requiresAuth: true,
      allowedRoles: ['superadmin', 'super_admin', 'standard_user', 'standarduser', 'monitoring_user', 'monitoringuser'],
      layout: 'DashboardLayout'
    }
  }
]