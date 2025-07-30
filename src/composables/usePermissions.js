import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

/**
 * Composable for handling role-based permissions in the Multi-Tenant SaaS Platform
 * Supports: Super Admin, Tenant, Devotee, Volunteer roles
 */
export function usePermissions() {
  const authStore = useAuthStore()

  // Helper function for case-insensitive status comparison
  const isStatusEqual = (status1, status2) => {
    return (status1 || '').toString().toLowerCase() === (status2 || '').toString().toLowerCase()
  }

  // User role and entity information
  const userRole = computed(() => authStore.user?.role)
  const userEntity = computed(() => authStore.user?.entityId)
  const userStatus = computed(() => authStore.user?.status)

  // Role checks
  const isSuperAdmin = computed(() => userRole.value === 'superadmin')
  const isTenant = computed(() => userRole.value === 'tenant')
  const isDevotee = computed(() => userRole.value === 'devotee')
  const isVolunteer = computed(() => userRole.value === 'volunteer')

  // Status checks - Updated to use case-insensitive comparison
  const isApproved = computed(() => isStatusEqual(userStatus.value, 'approved'))
  const isPending = computed(() => isStatusEqual(userStatus.value, 'pending'))
  const isRejected = computed(() => isStatusEqual(userStatus.value, 'rejected'))

  // Entity-level permissions
  const canManageEntity = computed(() => {
    return isTenant.value && isApproved.value
  })

  const canAccessEntityDashboard = computed(() => {
    return (isTenant.value && isApproved.value) || 
           (isDevotee.value && userEntity.value) || 
           (isVolunteer.value && userEntity.value)
  })

  // Super Admin permissions
  const canManageAllTenants = computed(() => isSuperAdmin.value)
  const canApproveTemples = computed(() => isSuperAdmin.value)
  const canViewAllEntities = computed(() => isSuperAdmin.value)

  // Tenant permissions
  const canCreateTemple = computed(() => isTenant.value)
  const canManageTemple = computed(() => isTenant.value && isApproved.value)
  const canUploadDocuments = computed(() => isTenant.value)

  // Entity Admin permissions (for approved tenants managing their temple)
  const canManageDevotees = computed(() => canManageEntity.value)
  const canManageSevas = computed(() => canManageEntity.value)
  const canManageDonations = computed(() => canManageEntity.value)
  const canManageEvents = computed(() => canManageEntity.value)
  const canManageVolunteers = computed(() => canManageEntity.value)
  const canSendBulkMessages = computed(() => canManageEntity.value)
  const canViewReports = computed(() => canManageEntity.value)
  const canExportData = computed(() => canManageEntity.value)

  // Devotee permissions
  const canJoinTemple = computed(() => isDevotee.value || isVolunteer.value)
  const canCompleteProfile = computed(() => isDevotee.value && userEntity.value)
  const canBookSeva = computed(() => isDevotee.value && userEntity.value)
  const canMakeDonation = computed(() => isDevotee.value && userEntity.value)
  const canRSVPEvents = computed(() => isDevotee.value && userEntity.value)
  const canViewOwnHistory = computed(() => isDevotee.value && userEntity.value)

  // Volunteer permissions
  const canViewAssignments = computed(() => isVolunteer.value && userEntity.value)
  const canViewTempleInfo = computed(() => isVolunteer.value && userEntity.value)

  // Navigation permissions
  const canAccessSuperAdminDashboard = computed(() => isSuperAdmin.value)
  const canAccessTenantDashboard = computed(() => isTenant.value)
  const canAccessDevoteeDashboard = computed(() => isDevotee.value && userEntity.value)
  const canAccessVolunteerDashboard = computed(() => isVolunteer.value && userEntity.value)

  // Resource-specific permissions
  const canViewResource = (resourceType, resourceEntityId = null) => {
    switch (resourceType) {
      case 'temple':
        return isSuperAdmin.value || 
               (isTenant.value && (!resourceEntityId || userEntity.value === resourceEntityId))
      
      case 'devotee':
        return isSuperAdmin.value || 
               (canManageEntity.value && userEntity.value === resourceEntityId) ||
               (isDevotee.value && userEntity.value === resourceEntityId)
      
      case 'seva':
        return isSuperAdmin.value || 
               (canManageEntity.value && userEntity.value === resourceEntityId) ||
               ((isDevotee.value || isVolunteer.value) && userEntity.value === resourceEntityId)
      
      case 'donation':
        return isSuperAdmin.value || 
               (canManageEntity.value && userEntity.value === resourceEntityId) ||
               (isDevotee.value && userEntity.value === resourceEntityId)
      
      case 'event':
        return isSuperAdmin.value || 
               (canManageEntity.value && userEntity.value === resourceEntityId) ||
               ((isDevotee.value || isVolunteer.value) && userEntity.value === resourceEntityId)
      
      case 'volunteer':
        return isSuperAdmin.value || 
               (canManageEntity.value && userEntity.value === resourceEntityId) ||
               (isVolunteer.value && userEntity.value === resourceEntityId)
      
      default:
        return false
    }
  }

  const canEditResource = (resourceType, resourceEntityId = null, resourceOwnerId = null) => {
    switch (resourceType) {
      case 'temple':
        return isSuperAdmin.value || 
               (isTenant.value && userEntity.value === resourceEntityId)
      
      case 'devotee':
        return isSuperAdmin.value || 
               (canManageEntity.value && userEntity.value === resourceEntityId) ||
               (isDevotee.value && authStore.user?.id === resourceOwnerId)
      
      case 'seva':
        return isSuperAdmin.value || 
               (canManageEntity.value && userEntity.value === resourceEntityId)
      
      case 'donation':
        return isSuperAdmin.value || 
               (canManageEntity.value && userEntity.value === resourceEntityId)
      
      case 'event':
        return isSuperAdmin.value || 
               (canManageEntity.value && userEntity.value === resourceEntityId)
      
      case 'volunteer':
        return isSuperAdmin.value || 
               (canManageEntity.value && userEntity.value === resourceEntityId)
      
      default:
        return false
    }
  }

  const canDeleteResource = (resourceType, resourceEntityId = null) => {
    switch (resourceType) {
      case 'temple':
        return isSuperAdmin.value
      
      case 'devotee':
      case 'seva':
      case 'donation':
      case 'event':
      case 'volunteer':
        return isSuperAdmin.value || 
               (canManageEntity.value && userEntity.value === resourceEntityId)
      
      default:
        return false
    }
  }

  // Route access permissions
  const canAccessRoute = (routeName) => {
    const routePermissions = {
      // Public routes
      'landing': true,
      'login': true,
      'register': true,
      'forgot-password': true,
      'terms': true,
      'privacy': true,

      // Super Admin routes
      'superadmin-dashboard': isSuperAdmin.value,
      'tenant-approvals': isSuperAdmin.value,
      'temple-approvals': isSuperAdmin.value,

      // Tenant routes
      'tenant-dashboard': isTenant.value,
      'create-temple': isTenant.value,
      'manage-temples': isTenant.value,

      // Entity Admin routes
      'entity-dashboard': canAccessEntityDashboard.value,
      'devotee-management': canManageDevotees.value,
      'seva-management': canManageSevas.value,
      'donation-management': canManageDonations.value,
      'event-management': canManageEvents.value,
      'volunteer-management': canManageVolunteers.value,
      'communication-center': canSendBulkMessages.value,

      // Devotee routes
      'temple-selection': canJoinTemple.value,
      'profile-creation': canCompleteProfile.value,
      'devotee-dashboard': canAccessDevoteeDashboard.value,
      'seva-booking': canBookSeva.value,
      'donation-history': canViewOwnHistory.value,
      'my-events': canRSVPEvents.value,

      // Volunteer routes
      'volunteer-temple-selection': canJoinTemple.value,
      'volunteer-dashboard': canAccessVolunteerDashboard.value,
      'my-assignments': canViewAssignments.value,
      'my-schedule': canViewAssignments.value,
    }

    return routePermissions[routeName] || false
  }

  // Helper function to check multiple permissions
  const hasAnyPermission = (...permissions) => {
    return permissions.some(permission => {
      if (typeof permission === 'string') {
        return canAccessRoute(permission)
      }
      return permission
    })
  }

  const hasAllPermissions = (...permissions) => {
    return permissions.every(permission => {
      if (typeof permission === 'string') {
        return canAccessRoute(permission)
      }
      return permission
    })
  }

  // Get user's allowed navigation items
  const getAllowedNavItems = () => {
    const navItems = []

    if (isSuperAdmin.value) {
      navItems.push(
        { name: 'Dashboard', route: 'superadmin-dashboard', icon: 'dashboard' },
        { name: 'Tenant Approvals', route: 'tenant-approvals', icon: 'check-circle' },
        { name: 'Temple Management', route: 'temple-approvals', icon: 'building' }
      )
    }

    if (isTenant.value) {
      navItems.push(
        { name: 'Dashboard', route: 'tenant-dashboard', icon: 'dashboard' }
      )
      
      if (canCreateTemple.value) {
        navItems.push(
          { name: 'Create Temple', route: 'create-temple', icon: 'plus' },
          { name: 'Manage Temples', route: 'manage-temples', icon: 'building' }
        )
      }
    }

    if (canAccessEntityDashboard.value) {
      navItems.push(
        { name: 'Dashboard', route: 'entity-dashboard', icon: 'dashboard' }
      )

      if (canManageDevotees.value) {
        navItems.push(
          { name: 'Devotees', route: 'devotee-management', icon: 'users' },
          { name: 'Sevas', route: 'seva-management', icon: 'heart' },
          { name: 'Donations', route: 'donation-management', icon: 'gift' },
          { name: 'Events', route: 'event-management', icon: 'calendar' },
          { name: 'Volunteers', route: 'volunteer-management', icon: 'user-plus' },
          { name: 'Communication', route: 'communication-center', icon: 'message-circle' }
        )
      }
    }

    if (isDevotee.value && userEntity.value) {
      navItems.push(
        { name: 'Dashboard', route: 'devotee-dashboard', icon: 'dashboard' },
        { name: 'Book Seva', route: 'seva-booking', icon: 'heart' },
        { name: 'Donations', route: 'donation-history', icon: 'gift' },
        { name: 'Events', route: 'my-events', icon: 'calendar' }
      )
    }

    if (isVolunteer.value && userEntity.value) {
      navItems.push(
        { name: 'Dashboard', route: 'volunteer-dashboard', icon: 'dashboard' },
        { name: 'Assignments', route: 'my-assignments', icon: 'clipboard' },
        { name: 'Schedule', route: 'my-schedule', icon: 'calendar' }
      )
    }

    return navItems
  }

  return {
    // Role checks
    isSuperAdmin,
    isTenant,
    isDevotee,
    isVolunteer,

    // Status checks
    isApproved,
    isPending,
    isRejected,

    // General permissions
    canManageEntity,
    canAccessEntityDashboard,

    // Super Admin permissions
    canManageAllTenants,
    canApproveTemples,
    canViewAllEntities,

    // Tenant permissions
    canCreateTemple,
    canManageTemple,
    canUploadDocuments,

    // Entity Admin permissions
    canManageDevotees,
    canManageSevas,
    canManageDonations,
    canManageEvents,
    canManageVolunteers,
    canSendBulkMessages,
    canViewReports,
    canExportData,

    // Devotee permissions
    canJoinTemple,
    canCompleteProfile,
    canBookSeva,
    canMakeDonation,
    canRSVPEvents,
    canViewOwnHistory,

    // Volunteer permissions
    canViewAssignments,
    canViewTempleInfo,

    // Navigation permissions
    canAccessSuperAdminDashboard,
    canAccessTenantDashboard,
    canAccessDevoteeDashboard,
    canAccessVolunteerDashboard,

    // Resource permissions
    canViewResource,
    canEditResource,
    canDeleteResource,

    // Route permissions
    canAccessRoute,

    // Helper functions
    hasAnyPermission,
    hasAllPermissions,
    getAllowedNavItems,
    isStatusEqual,

    // User info
    userRole,
    userEntity,
    userStatus
  }
}