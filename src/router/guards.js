// src/router/guards.js - UPDATED FOR FIXED ROLE MAPPING
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import { useDevoteeStore } from '@/stores/devotee'

/**
 * Role mapping - map backend roles to frontend roles
 * This will handle the discrepancy between backend and frontend role names
 */
const roleMapping = {
  'templeadmin': 'tenant', // Backend returns 'templeadmin', but routes expect 'tenant'
  'standarduser': 'standard_user',
  'monitoringuser': 'monitoring_user'
}

/**
 * Get mapped role (or original if no mapping exists)
 */
function getMappedRole(role) {
  if (!role) return null;
  return roleMapping[role.toLowerCase()] || role.toLowerCase();
}

/**
 * Helper function for case-insensitive status comparison
 */
function isStatusEqual(status1, status2) {
  return (status1 || '').toString().toLowerCase() === (status2 || '').toString().toLowerCase()
}

/**
 * Authentication Guard - ✅ UPDATED FOR CORRECT ROLE HANDLING
 * Checks if user is authenticated before accessing protected routes
 */
export function requireAuth(to, from, next) {
  const authStore = useAuthStore()
  
  // Debug auth state
  console.log('🔐 Auth Check:', {
    isAuthenticated: authStore.isAuthenticated,
    user: authStore.user,
    userRole: authStore.userRole, // ✅ This should now work correctly
    needsTenantSelection: authStore.needsTenantSelection,
    route: to.path,
    toName: to.name
  })
  
  if (!authStore.isAuthenticated) {
    console.log('⚠️ Not authenticated, redirecting to login')
    next({
      name: 'Login',
      query: { 
        redirect: to.fullPath,
        message: 'Please login to access this page'
      }
    })
    return false
  }
  
  // ✅ SIMPLIFIED: Use the needsTenantSelection computed property from auth store
  if (authStore.needsTenantSelection && 
      to.name !== 'TenantSelection' && 
      !to.path.includes('/tenant-selection')) {
    
    console.log('📍 CRITICAL REDIRECT: User needs tenant selection, redirecting...')
    console.log('- User role:', authStore.userRole)
    console.log('- Needs tenant selection:', authStore.needsTenantSelection)
    
    // Force redirect to tenant selection page using route name for consistency
    next({ name: 'TenantSelection' });
    return false;
  }
  
  // ✅ ALLOW: Superadmin can access tenant selection too
  if (authStore.isSuperAdmin && to.name === 'TenantSelection') {
    console.log('📍 SuperAdmin accessing tenant selection - allowed');
    next();
    return true;
  }
  
  // Devotee role check for temple selection redirection
  if (authStore.isDevotee && to.name !== 'DevoteeTempleSelection' && !to.path.includes('/temple-selection')) {
    const entityId = to.params.id
    // Only redirect if not accessing entity specific routes
    if (!entityId || !to.path.includes(`/entity/${entityId}/`)) {
      next({ name: 'DevoteeTempleSelection' })
      return false
    }
  }
  
  console.log('✅ Auth guard passed for role:', authStore.userRole);
  next()
  return true
}

/**
 * Guest Guard - ✅ UPDATED
 * Redirects authenticated users away from guest-only pages (login, register)
 */
export function requireGuest(to, from, next) {
  const authStore = useAuthStore()
  
  if (authStore.isAuthenticated) {
    // ✅ Use auth store's getDashboardPath method
    const redirectPath = authStore.getDashboardPath(authStore.userRole)
    console.log('🚪 Guest guard: authenticated user redirected to:', redirectPath)
    next({ path: redirectPath })
    return false
  }
  
  next()
  return true
}

/**
 * SPECIFIC GUARD: Check Profile Completed (for devotee routes)
 * Ensures devotee has completed their profile before accessing certain features
 */
export function checkProfileCompleted(to, from, next) {
  const authStore = useAuthStore()
  const devoteeStore = useDevoteeStore()
  const { showToast } = useToast()
  
  if (!authStore.isAuthenticated) {
    next({ name: 'Login' })
    return
  }
  
  if (!authStore.isDevotee) {
    next({ name: 'Unauthorized' })
    return
  }
  
  // Check if profile is completed
  if (!devoteeStore.isProfileComplete && !authStore.user?.profileCompleted) {
    showToast('Please complete your profile first', 'warning')
    
    // If entity ID exists, redirect to entity-specific profile creation
    const entityId = to.params.id
    if (entityId) {
      next({ name: 'DevoteeProfileCreation', params: { id: entityId } })
    } else {
      next({ name: 'DevoteeTempleSelection' })
    }
    return
  }
  
  next()
}

/**
 * Role-based Access Guard - ✅ ENHANCED VERSION
 * Checks if user has required role to access the route
 */
export function requireRole(roles) {
  return (to, from, next) => {
    const authStore = useAuthStore()
    const { showToast } = useToast()
    
    if (!authStore.isAuthenticated) {
      next({
        name: 'Login',
        query: { redirect: to.fullPath }
      })
      return
    }
    
    const userRole = authStore.userRole || ''
    const normalizedRole = userRole.toLowerCase().trim()
    
    // ✅ BUILD: Complete list of user roles (including variations)
    const userRoles = [normalizedRole]
    
    // Add mapped version of the role if it exists
    if (roleMapping[normalizedRole]) {
      userRoles.push(roleMapping[normalizedRole])
    }
    
    // ✅ ENHANCED: For special roles, also check the alternative format
    const roleVariations = {
      'standard_user': ['standarduser'],
      'monitoring_user': ['monitoringuser'],
      'standarduser': ['standard_user'],
      'monitoringuser': ['monitoring_user'],
      'super_admin': ['superadmin'],
      'superadmin': ['super_admin']
    }
    
    if (roleVariations[normalizedRole]) {
      userRoles.push(...roleVariations[normalizedRole])
    }
    
    // Normalize required roles for comparison
    const normalizedRequiredRoles = roles.map(role => role.toLowerCase().trim())
    
    // Check if any of the user's roles are included in required roles
    const hasRequiredRole = normalizedRequiredRoles.some(requiredRole => 
      userRoles.includes(requiredRole)
    )
    
    console.log('🎭 Role check:', {
      userRole,
      userRoles,
      requiredRoles: normalizedRequiredRoles,
      hasAccess: hasRequiredRole
    })
    
    if (!hasRequiredRole) {
      showToast(`Access denied: User has role "${userRole}" but route requires one of: ${roles.join(', ')}`, 'error')
      next({ name: 'Unauthorized' })
      return
    }
    
    next()
  }
}

/**
 * SPECIFIC GUARD: Check Role (for route files) - ✅ ENHANCED VERSION
 * This is the specific function your route files are importing
 */
export function checkRole(to, from, next, requiredRole) {
  const authStore = useAuthStore()
  const { showToast } = useToast()
  
  if (!authStore.isAuthenticated) {
    next({
      name: 'Login',
      query: { redirect: to.fullPath }
    })
    return
  }
  
  const userRole = authStore.userRole || ''
  const normalizedRole = userRole.toLowerCase().trim()
  const normalizedRequiredRole = requiredRole.toLowerCase().trim()
  
  // ✅ BUILD: Complete list of user roles (including variations)
  const userRoles = [normalizedRole]
  
  // Add mapped version of the role if it exists
  if (roleMapping[normalizedRole]) {
    userRoles.push(roleMapping[normalizedRole])
  }
  
  // ✅ ENHANCED: For special roles, also check the alternative format
  const roleVariations = {
    'standard_user': ['standarduser'],
    'monitoring_user': ['monitoringuser'],
    'standarduser': ['standard_user'],
    'monitoringuser': ['monitoring_user'],
    'super_admin': ['superadmin'],
    'superadmin': ['super_admin']
  }
  
  if (roleVariations[normalizedRole]) {
    userRoles.push(...roleVariations[normalizedRole])
  }
  
  // Compare user roles with required role
  if (!userRoles.includes(normalizedRequiredRole)) {
    showToast(`Access denied: User has role "${userRole}" but route requires: ${requiredRole}`, 'error')
    next({ name: 'Unauthorized' })
    return
  }
  
  next()
}

/**
 * Get default route based on user role - ✅ ENHANCED VERSION
 */
export function getDefaultRoute(role) {
  if (!role) return '/';
  
  const normalizedRole = role.toLowerCase().trim()
  
  const routes = {
    'tenant': '/tenant/dashboard',
    'templeadmin': '/tenant/dashboard', 
    'devotee': '/devotee/temple-selection',
    'volunteer': '/volunteer/temple-selection',
    'superadmin': '/superadmin/dashboard',
    'super_admin': '/superadmin/dashboard',
    'standard_user': '/tenant-selection',
    'standarduser': '/tenant-selection',
    'monitoring_user': '/tenant-selection',
    'monitoringuser': '/tenant-selection'
  }
  
  return routes[normalizedRole] || '/'
}

// Export utility functions
export {
  getMappedRole,
  isStatusEqual
}

// Export the default object with all guard functions
export default {
  requireAuth,
  requireGuest,
  requireRole,
  checkRole,
  checkProfileCompleted,
  getDefaultRoute
}
