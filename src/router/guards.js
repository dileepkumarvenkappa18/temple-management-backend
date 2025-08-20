// src/router/guards.js
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
  'monitoringuser': 'monitoring_user' // Update to map standarduser to standard_user
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
 * Authentication Guard
 * Checks if user is authenticated before accessing protected routes
 */
export function requireAuth(to, from, next) {
  const authStore = useAuthStore()
  
  // Debug auth state
  console.log('🔐 Auth Check:', {
    isAuthenticated: authStore.isAuthenticated,
    user: authStore.user,
    role: authStore.user?.role,
    mappedRole: getMappedRole(authStore.user?.role),
    route: to.path,
    toName: to.name,
    toPath: to.path.includes('/tenant-selection')
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
  
  // Get user role for role-based checks
  const userRole = authStore.user?.role || '';
  const normalizedRole = userRole.toLowerCase();
  
  // CRITICAL FIX: Check for both formats of roles (with underscore and without)
  const isStandardUser = normalizedRole === 'standard_user' || normalizedRole === 'standarduser';
  const isMonitoringUser = normalizedRole === 'monitoring_user' || normalizedRole === 'monitoringuser';
  
  // Direct standard and monitoring users to tenant selection if not already there
  if ((isStandardUser || isMonitoringUser) && 
      to.name !== 'TenantSelection' && 
      !to.path.includes('/tenant-selection')) {
    
    console.log(`📍 Role ${userRole} redirecting to tenant selection`);
    
    // Force redirect to tenant selection page
    next({ path: '/tenant-selection' });
    return false;
  }
  
  // Devotee role check for temple selection redirection
  if (normalizedRole === 'devotee' && to.name !== 'DevoteeTempleSelection' && !to.path.includes('/temple-selection')) {
    const entityId = to.params.id
    // Only redirect if not accessing entity specific routes
    if (!entityId || !to.path.includes(`/entity/${entityId}/`)) {
      next({ name: 'DevoteeTempleSelection' })
      return false
    }
  }
  
  next()
  return true
}

/**
 * Guest Guard
 * Redirects authenticated users away from guest-only pages (login, register)
 */
export function requireGuest(to, from, next) {
  const authStore = useAuthStore()
  
  if (authStore.isAuthenticated) {
    // Use mapped role for redirection
    const mappedRole = getMappedRole(authStore.user?.role)
    const redirectPath = getDefaultRoute(mappedRole)
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
  
  const mappedRole = getMappedRole(authStore.user?.role)
  
  if (mappedRole !== 'devotee') {
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
 * Role-based Access Guard
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
    
    const userRole = authStore.user?.role || '';
    const normalizedRole = userRole.toLowerCase();
    
    // IMPROVED: Check for both versions of role names
    const userRoles = [normalizedRole];
    
    // Add mapped version of the role if it exists
    if (roleMapping[normalizedRole]) {
      userRoles.push(roleMapping[normalizedRole]);
    }
    
    // For 'standard_user' or 'monitoring_user', also check the alternative format
    if (normalizedRole === 'standard_user') userRoles.push('standarduser');
    if (normalizedRole === 'monitoring_user') userRoles.push('monitoringuser');
    if (normalizedRole === 'standarduser') userRoles.push('standard_user');
    if (normalizedRole === 'monitoringuser') userRoles.push('monitoring_user');
    
    // Check if any of the user's roles are included in required roles
    const hasRequiredRole = roles.some(role => userRoles.includes(role.toLowerCase()));
    
    if (!hasRequiredRole) {
      showToast(`Access denied: User has role "${userRole}" but route requires one of: ${roles.join(', ')}`, 'error')
      next({ name: 'Unauthorized' })
      return
    }
    
    next()
  }
}

/**
 * SPECIFIC GUARD: Check Role (for route files)
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
  
  const userRole = authStore.user?.role || '';
  const normalizedRole = userRole.toLowerCase();
  
  // IMPROVED: Check for both versions of role names
  const userRoles = [normalizedRole];
  
  // Add mapped version of the role if it exists
  if (roleMapping[normalizedRole]) {
    userRoles.push(roleMapping[normalizedRole]);
  }
  
  // For standard_user and monitoring_user, also check the alternative format
  if (normalizedRole === 'standard_user') userRoles.push('standarduser');
  if (normalizedRole === 'monitoring_user') userRoles.push('monitoringuser');
  if (normalizedRole === 'standarduser') userRoles.push('standard_user');
  if (normalizedRole === 'monitoringuser') userRoles.push('monitoring_user');
  
  // Compare user roles with required role
  if (!userRoles.includes(requiredRole.toLowerCase())) {
    showToast(`Access denied: User has role "${userRole}" but route requires: ${requiredRole}`, 'error')
    next({ name: 'Unauthorized' })
    return
  }
  
  next()
}

/**
 * Get default route based on user role
 */
export function getDefaultRoute(role) {
  if (!role) return '/';
  
  const normalizedRole = role.toLowerCase();
  
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