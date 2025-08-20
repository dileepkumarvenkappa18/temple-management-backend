import SuperAdminDashboard from '@/views/superadmin/SuperAdminDashboard.vue';
import TenantApprovals from '@/components/superadmin/TenantApprovals.vue'; // Direct import from components

export default [
  {
    path: 'dashboard',
    name: 'SuperAdminDashboard',
    component: () => import('@/views/superadmin/SuperAdminDashboard.vue'),
    meta: { title: 'Super Admin Dashboard', breadcrumb: 'Dashboard' }
  },
  {
    path: 'tenant-approvals',
    name: 'TenantApprovals',
    component: () => import('@/components/superadmin/TenantApprovals.vue'),
    meta: { title: 'Tenant Approvals', breadcrumb: 'Tenant Approvals' }
  },
  {
    path: 'user-management',
    name: 'SuperadminUserManagement',
    component: () => import('@/views/superadmin/UserManagement.vue'),
    meta: { requiresAuth: true, role: 'superadmin' }
  },
  {
    path: 'role-management',
    name: 'SuperadminRoleManagement',
    component: () => import('@/views/superadmin/RoleManagement.vue'),
    meta: { requiresAuth: true, role: 'superadmin' }
  },
  {
    path: 'reset-password',
    name: 'ResetPassword',
    component: () => import('@/views/superadmin/ResetPassword.vue'),
    meta: { requiresAuth: true, role: 'superadmin' }
  },
  {
    path: 'reset-password/:userId',
    name: 'ResetPasswordForm',
    component: () => import('@/views/superadmin/ResetPasswordForm.vue'),
    meta: { requiresAuth: true, role: 'superadmin', title: 'Reset User Password' }
  },
  // New route for tenant assignment
  {
    path: 'users/:userId/assign-tenants',
    name: 'AssignTenants',
    component: () => import('@/views/superadmin/AssignTenantsView.vue'),
    meta: {
      requiresAuth: true,
      role: 'superadmin',
      title: 'Assign Tenants to User'
    }
  },
  // New route for audit logs
  {
    path: 'audit-logs',
    name: 'AuditLogs',
    component: () => import('@/views/superadmin/AuditLogsView.vue'),
    meta: {
      requiresAuth: true,
      role: 'superadmin',
      title: 'Audit Logs',
      breadcrumb: 'Audit Logs'
    }
  },
  // Reports management routes
  {
    path: 'reports',
    name: 'superadmin-reports',
    component: () => import('@/views/superadmin/SuperadminReportsView.vue'),
    meta: {
      requiresAuth: true,
      role: 'superadmin',
      title: 'Reports Management'
    }
  },
  {
    path: 'reports/temple-register',
    name: 'SuperadminTempleRegister',
    component: () => import('@/views/tenant/reports/TempleRegisterReport.vue'),
    meta: {
      requiresAuth: true,
      role: 'superadmin',
      title: 'Temple Register Report'
    }
  },
  {
    path: 'reports/temple-activities',
    name: 'SuperadminTempleActivities',
    component: () => import('@/views/tenant/reports/TempleActivitiesReport.vue'),
    meta: {
      requiresAuth: true,
      role: 'superadmin',
      title: 'Temple Activities Report'
    }
  },
  {
    path: 'reports/birthdays',
    name: 'SuperadminBirthdays',
    component: () => import('@/views/tenant/reports/BirthdaysReport.vue'),
    meta: {
      requiresAuth: true,
      role: 'superadmin',
      title: 'Devotee Birthdays Report'
    }
  },
  {
    path: 'reports/user-details',
    name: 'SuperadminUserDetails',
    component: () => import('@/views/tenant/reports/UserDetailsReport.vue'),
    meta: {
      requiresAuth: true,
      role: 'superadmin',
      title: 'User Details Report'
    }
  },
  {
    path: 'reports/approval-status',
    name: 'SuperadminApprovalStatus',
    component: () => import('@/views/tenant/reports/ApprovalStatusReport.vue'),
    meta: {
      requiresAuth: true,
      role: 'superadmin',
      title: 'Approval Status Report'
    }
  }
];