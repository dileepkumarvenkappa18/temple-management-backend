// src/utils/permissions.js
// Role-Based Access Control Utility for Multi-Tenant Temple SaaS Platform

/**
 * User Roles Definition
 */
export const ROLES = {
  SUPER_ADMIN: 'superadmin',
  TENANT: 'tenant',
  ENTITY_ADMIN: 'entity_admin', // Temple Admin after approval
  DEVOTEE: 'devotee',
  VOLUNTEER: 'volunteer'
}

/**
 * Permission Categories
 */
export const PERMISSIONS = {
  // Super Admin Permissions
  SUPER_ADMIN: {
    MANAGE_TENANTS: 'manage_tenants',
    APPROVE_TEMPLES: 'approve_temples',
    VIEW_ALL_DATA: 'view_all_data',
    SYSTEM_SETTINGS: 'system_settings'
  },
  
  // Tenant Permissions (Temple Owner)
  TENANT: {
    CREATE_TEMPLE: 'create_temple',
    MANAGE_OWN_TEMPLES: 'manage_own_temples',
    VIEW_APPROVAL_STATUS: 'view_approval_status',
    UPLOAD_DOCUMENTS: 'upload_documents'
  },
  
  // Entity Admin Permissions (Approved Temple Admin)
  ENTITY_ADMIN: {
    MANAGE_DEVOTEES: 'manage_devotees',
    MANAGE_SEVAS: 'manage_sevas',
    MANAGE_DONATIONS: 'manage_donations',
    MANAGE_EVENTS: 'manage_events',
    MANAGE_VOLUNTEERS: 'manage_volunteers',
    SEND_COMMUNICATIONS: 'send_communications',
    VIEW_REPORTS: 'view_reports',
    EXPORT_DATA: 'export_data'
  },
  
  // Devotee Permissions
  DEVOTEE: {
    JOIN_TEMPLE: 'join_temple',
    COMPLETE_PROFILE: 'complete_profile',
    BOOK_SEVA: 'book_seva',
    MAKE_DONATION: 'make_donation',
    VIEW_EVENTS: 'view_events',
    RSVP_EVENTS: 'rsvp_events',
    VIEW_OWN_HISTORY: 'view_own_history'
  },
  
  // Volunteer Permissions
  VOLUNTEER: {
    JOIN_TEMPLE: 'join_temple',
    VIEW_ASSIGNMENTS: 'view_assignments',
    VIEW_EVENTS: 'view_events',
    VIEW_TEMPLE_INFO: 'view_temple_info'
  }
}

/**
 * Role-Permission Mapping
 */
export const ROLE_PERMISSIONS = {
  [ROLES.SUPER_ADMIN]: [
    ...Object.values(PERMISSIONS.SUPER_ADMIN),
    ...Object.values(PERMISSIONS.TENANT),
    ...Object.values(PERMISSIONS.ENTITY_ADMIN),
    ...Object.values(PERMISSIONS.DEVOTEE),
    ...Object.values(PERMISSIONS.VOLUNTEER)
  ],
  
  [ROLES.TENANT]: [
    ...Object.values(PERMISSIONS.TENANT)
  ],
  
  [ROLES.ENTITY_ADMIN]: [
    ...Object.values(PERMISSIONS.ENTITY_ADMIN),
    ...Object.values(PERMISSIONS.DEVOTEE) // Entity admin can do devotee actions too
  ],
  
  [ROLES.DEVOTEE]: [
    ...Object.values(PERMISSIONS.DEVOTEE)
  ],
  
  [ROLES.VOLUNTEER]: [
    ...Object.values(PERMISSIONS.VOLUNTEER)
  ]
}

/**
 * Route Access Control
 */
export const ROUTE_PERMISSIONS = {
  // Super Admin Routes
  '/superadmin': [ROLES.SUPER_ADMIN],
  '/superadmin/dashboard': [ROLES.SUPER_ADMIN],
  
  // Tenant Routes
  '/tenant': [ROLES.TENANT],
  '/tenant/dashboard': [ROLES.TENANT],
  '/tenant/entities/create': [ROLES.TENANT],
  '/tenant/entities': [ROLES.TENANT],
  
  // Entity Admin Routes
  '/entity/:id/dashboard': [ROLES.ENTITY_ADMIN],
  '/entity/:id/devotees': [ROLES.ENTITY_ADMIN],
  '/entity/:id/sevas': [ROLES.ENTITY_ADMIN],
  '/entity/:id/donations': [ROLES.ENTITY_ADMIN],
  '/entity/:id/events': [ROLES.ENTITY_ADMIN],
  '/entity/:id/volunteers': [ROLES.ENTITY_ADMIN],
  '/entity/:id/communication': [ROLES.ENTITY_ADMIN],
  
  // Devotee Routes
  '/entity/:id/devotee/dashboard': [ROLES.DEVOTEE],
  '/entity/:id/devotee/profile': [ROLES.DEVOTEE],
  '/entity/:id/devotee/sevas': [ROLES.DEVOTEE],
  '/entity/:id/devotee/donations': [ROLES.DEVOTEE],
  '/entity/:id/devotee/events': [ROLES.DEVOTEE],
  '/temple-selection': [ROLES.DEVOTEE],
  '/profile-creation': [ROLES.DEVOTEE],
  
  // Volunteer Routes
  '/entity/:id/volunteer/dashboard': [ROLES.VOLUNTEER],
  '/entity/:id/volunteer/assignments': [ROLES.VOLUNTEER],
  '/entity/:id/volunteer/schedule': [ROLES.VOLUNTEER],
  '/volunteer/temple-selection': [ROLES.VOLUNTEER],
  
  // Public Routes (accessible to all)
  '/': ['*'],
  '/login': ['*'],
  '/register': ['*'],
  '/forgot-password': ['*'],
  '/terms': ['*'],
  '/privacy': ['*']
}

/**
 * Check if user has specific permission
 * @param {string} userRole - User's role
 * @param {string} permission - Permission to check
 * @returns {boolean}
 */
export const hasPermission = (userRole, permission) => {
  if (!userRole || !permission) return false
  
  const rolePermissions = ROLE_PERMISSIONS[userRole] || []
  return rolePermissions.includes(permission)
}

/**
 * Check if user has any of the specified permissions
 * @param {string} userRole - User's role
 * @param {string[]} permissions - Array of permissions to check
 * @returns {boolean}
 */
export const hasAnyPermission = (userRole, permissions) => {
  if (!userRole || !Array.isArray(permissions)) return false
  
  return permissions.some(permission => hasPermission(userRole, permission))
}

/**
 * Check if user has all specified permissions
 * @param {string} userRole - User's role
 * @param {string[]} permissions - Array of permissions to check
 * @returns {boolean}
 */
export const hasAllPermissions = (userRole, permissions) => {
  if (!userRole || !Array.isArray(permissions)) return false
  
  return permissions.every(permission => hasPermission(userRole, permission))
}

/**
 * Check if user can access a specific route
 * @param {string} userRole - User's role
 * @param {string} routePath - Route path to check
 * @returns {boolean}
 */
export const canAccessRoute = (userRole, routePath) => {
  if (!userRole || !routePath) return false
  
  // Check exact route match first
  const allowedRoles = ROUTE_PERMISSIONS[routePath]
  if (allowedRoles) {
    return allowedRoles.includes('*') || allowedRoles.includes(userRole)
  }
  
  // Check dynamic routes (with :id parameters)
  const dynamicRouteMatch = Object.keys(ROUTE_PERMISSIONS).find(route => {
    if (!route.includes(':')) return false
    
    const routePattern = route.replace(/:[^/]+/g, '[^/]+')
    const regex = new RegExp(`^${routePattern}$`)
    return regex.test(routePath)
  })
  
  if (dynamicRouteMatch) {
    const allowedRoles = ROUTE_PERMISSIONS[dynamicRouteMatch]
    return allowedRoles.includes('*') || allowedRoles.includes(userRole)
  }
  
  return false
}

/**
 * Get user's dashboard route based on role
 * @param {string} userRole - User's role
 * @param {string} entityId - Entity ID (for entity-specific roles)
 * @returns {string}
 */
export const getDashboardRoute = (userRole, entityId = null) => {
  switch (userRole) {
    case ROLES.SUPER_ADMIN:
      return '/superadmin/dashboard'
    case ROLES.TENANT:
      return '/tenant/dashboard'
    case ROLES.ENTITY_ADMIN:
      return entityId ? `/entity/${entityId}/dashboard` : '/entity/dashboard'
    case ROLES.DEVOTEE:
      return entityId ? `/entity/${entityId}/devotee/dashboard` : '/temple-selection'
    case ROLES.VOLUNTEER:
      return entityId ? `/entity/${entityId}/volunteer/dashboard` : '/volunteer/temple-selection'
    default:
      return '/'
  }
}

/**
 * Check if user needs to complete profile
 * @param {object} user - User object
 * @returns {boolean}
 */
export const needsProfileCompletion = (user) => {
  if (user.role !== ROLES.DEVOTEE) return false
  
  // Check if devotee has completed mandatory profile steps
  const requiredFields = [
    'personal_details_completed',
    'spiritual_info_completed',
    'family_lineage_completed',
    'seva_preferences_completed'
  ]
  
  return requiredFields.some(field => !user[field])
}

/**
 * Check if user needs temple selection
 * @param {object} user - User object
 * @returns {boolean}
 */
export const needsTempleSelection = (user) => {
  return (user.role === ROLES.DEVOTEE || user.role === ROLES.VOLUNTEER) && !user.entity_id
}

/**
 * Get available actions for user role
 * @param {string} userRole - User's role
 * @param {string} context - Context (e.g., 'devotee', 'seva', 'donation')
 * @returns {string[]}
 */
export const getAvailableActions = (userRole, context) => {
  const actions = {
    devotee: {
      [ROLES.ENTITY_ADMIN]: ['view', 'edit', 'deactivate', 'export'],
      [ROLES.DEVOTEE]: ['view_own', 'edit_own']
    },
    seva: {
      [ROLES.ENTITY_ADMIN]: ['create', 'view', 'edit', 'delete', 'approve', 'reject'],
      [ROLES.DEVOTEE]: ['view', 'book'],
      [ROLES.VOLUNTEER]: ['view']
    },
    donation: {
      [ROLES.ENTITY_ADMIN]: ['view_all', 'generate_receipt', 'export'],
      [ROLES.DEVOTEE]: ['make', 'view_own', 'download_receipt']
    },
    event: {
      [ROLES.ENTITY_ADMIN]: ['create', 'view', 'edit', 'delete', 'manage_rsvp'],
      [ROLES.DEVOTEE]: ['view', 'rsvp'],
      [ROLES.VOLUNTEER]: ['view']
    }
  }
  
  return actions[context]?.[userRole] || []
}

/**
 * Permission validation middleware helper
 * @param {string} requiredPermission - Required permission
 * @returns {function}
 */
export const requirePermission = (requiredPermission) => {
  return (user) => {
    return hasPermission(user.role, requiredPermission)
  }
}

/**
 * Role validation middleware helper
 * @param {string[]} allowedRoles - Allowed roles
 * @returns {function}
 */
export const requireRole = (allowedRoles) => {
  return (user) => {
    return allowedRoles.includes(user.role)
  }
}

/**
 * Entity access validation
 * @param {object} user - User object
 * @param {string} entityId - Entity ID to access
 * @returns {boolean}
 */
export const canAccessEntity = (user, entityId) => {
  // Super admin can access all entities
  if (user.role === ROLES.SUPER_ADMIN) return true
  
  // Users can only access their own entity
  return user.entity_id === entityId
}

/**
 * Status-based permissions
 */
export const STATUS_PERMISSIONS = {
  TENANT_PENDING: {
    canCreateTemple: false,
    canViewDashboard: true,
    message: 'Your account is pending approval. You will be notified once approved.'
  },
  TENANT_APPROVED: {
    canCreateTemple: true,
    canViewDashboard: true,
    message: 'Your account is approved. You can now create and manage temples.'
  },
  TEMPLE_PENDING: {
    canManageDevotees: false,
    canManageSevas: false,
    message: 'Your temple is pending approval. You will be notified once approved.'
  },
  TEMPLE_APPROVED: {
    canManageDevotees: true,
    canManageSevas: true,
    message: 'Your temple is approved. You can now manage all temple activities.'
  },
  TEMPLE_REJECTED: {
    canManageDevotees: false,
    canManageSevas: false,
    message: 'Your temple application was rejected. Please review the feedback and resubmit.'
  }
}

/**
 * Get status-based permissions
 * @param {string} status - Current status
 * @returns {object}
 */
export const getStatusPermissions = (status) => {
  return STATUS_PERMISSIONS[status] || {}
}

export default {
  ROLES,
  PERMISSIONS,
  ROLE_PERMISSIONS,
  ROUTE_PERMISSIONS,
  hasPermission,
  hasAnyPermission,
  hasAllPermissions,
  canAccessRoute,
  getDashboardRoute,
  needsProfileCompletion,
  needsTempleSelection,
  getAvailableActions,
  requirePermission,
  requireRole,
  canAccessEntity,
  getStatusPermissions
}