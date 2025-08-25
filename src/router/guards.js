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
 * Authentication Guard - Updated for assigned tenant handling
 * Checks if user is authenticated before accessing protected routes
 */
export function requireAuth(to, from, next) {
  const authStore = useAuthStore()
  
  // Debug auth state
  console.log('🔐 Auth Check:', {
    isAuthenticated: authStore.isAuthenticated,
    user: authStore.user,
    userRole: authStore.userRole,
    needsTenantSelection: authStore.needsTenantSelection,
    assignedTenantId: authStore.assignedTenantId,
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
  
  // Check assigned tenant ID before redirecting to tenant selection
  if (authStore.needsTenantSelection && 
      to.name !== 'TenantSelection' && 
      !to.path.includes('/tenant-selection')) {
    
    console.log('📍 User needs tenant selection, redirecting...')
    console.log('- User role:', authStore.userRole)
    console.log('- Needs tenant selection:', authStore.needsTenantSelection)
    console.log('- Assigned tenant ID:', authStore.assignedTenantId)
    
    // Force redirect to tenant selection page using route name for consistency
    next({ name: 'TenantSelection' });
    return false;
  }
  
  // ALLOW: Superadmin can access tenant selection too
  if (authStore.isSuperAdmin && to.name === 'TenantSelection') {
    console.log('📍 SuperAdmin accessing tenant selection - allowed');
    next();
    return true;
  }
  
  // MODIFIED: Removed the condition that would redirect from temple selection to dashboard
  // when a devotee has already joined a temple. Now they'll always go to temple selection first.
  
  console.log('✅ Auth guard passed for role:', authStore.userRole);
  next()
  return true
}


function normalizeRole(role) {
  if (!role) return ''
  const r = role.toLowerCase().trim()
  if (r === 'monitoringuser' || r === 'monitoring_user') return 'monitoring_user'
  if (r === 'standarduser' || r === 'standard_user') return 'standard_user'
  if (r === 'superadmin' || r === 'super_admin') return 'superadmin'
  return r
}

/**
 * Guest Guard - Updated with assigned tenant awareness
 * Redirects authenticated users away from guest-only pages (login, register)
 */
export function requireGuest(to, from, next) {
  const authStore = useAuthStore()
  
  if (authStore.isAuthenticated) {
    // Use auth store's getDashboardPath method
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

export function setupRouteGuards(router) {
  router.beforeEach((to, from, next) => {
    const authStore = useAuthStore();
    const isAuthenticated = authStore.isAuthenticated;
    const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
    const userRole = authStore.userRole?.toLowerCase() || '';
    
    console.log('🔐 Auth Check:', {
      isAuthenticated,
      userRole,
      route: to.path,
      requiresAuth
    });
    
    // Special check for entity routes with tenant ID
    if (to.path.includes('/entity/') && to.params.id) {
      console.log('🏛️ Entity access check for tenant ID:', to.params.id);
      
      // Store the current entity ID
      localStorage.setItem('current_entity_id', to.params.id);
      
      // For superadmin, always allow access to any entity
      if (userRole.includes('superadmin')) {
        console.log('✅ SuperAdmin accessing entity - allowed');
        return next();
      }
    }
    
    if (requiresAuth && !isAuthenticated) {
      console.log('❌ Authentication required - redirecting to login');
      return next('/auth/login');
    }
    
    // All checks pass, continue
    next();
  });
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
    
    const userRole = authStore.userRole || ''
    const normalizedRole = userRole.toLowerCase().trim()
    
    // BUILD: Complete list of user roles (including variations)
    const userRoles = [normalizedRole]
    
    // Add mapped version of the role if it exists
    if (roleMapping[normalizedRole]) {
      userRoles.push(roleMapping[normalizedRole])
    }
    
    // ENHANCED: For special roles, also check the alternative format
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
  
  const userRole = authStore.userRole || ''
  const normalizedRole = userRole.toLowerCase().trim()
  const normalizedRequiredRole = requiredRole.toLowerCase().trim()
  
  // BUILD: Complete list of user roles (including variations)
  const userRoles = [normalizedRole]
  
  // Add mapped version of the role if it exists
  if (roleMapping[normalizedRole]) {
    userRoles.push(roleMapping[normalizedRole])
  }
  
  // ENHANCED: For special roles, also check the alternative format
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
 * Get default route based on user role - UPDATED with assigned tenant support
 */
export function getDefaultRoute(role) {
  if (!role) return '/';
  
  const normalizedRole = role.toLowerCase().trim()
  const authStore = useAuthStore();
  
  const routes = {
    'tenant': '/tenant/dashboard',
    'templeadmin': '/tenant/dashboard', 
    'devotee': '/devotee/temple-selection',
    'volunteer': '/volunteer/temple-selection',
    'superadmin': '/superadmin/dashboard',
    'super_admin': '/superadmin/dashboard'
  }
  
  // Special handling for standard_user and monitoring_user
  if (normalizedRole === 'standard_user' || normalizedRole === 'standarduser' || 
      normalizedRole === 'monitoring_user' || normalizedRole === 'monitoringuser') {
    
    // Use assigned tenant ID if available
    if (authStore.assignedTenantId) {
      // UPDATED: Redirect to tenant dashboard instead of entity dashboard
      return `/tenant/${authStore.assignedTenantId}/dashboard`;
    }
    return '/tenant-selection';
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