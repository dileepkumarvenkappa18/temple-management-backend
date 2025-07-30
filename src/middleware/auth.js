// src/middleware/authMiddleware.js
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'

// Map numeric role IDs from backend to string roles in frontend
const ROLE_MAP = {
  1: 'superadmin',
  2: 'tenant',
  3: 'devotee',
  4: 'volunteer'
}

/**
 * Authentication middleware
 * Checks if the user is authenticated and has the required role
 */
export default function authMiddleware(to, from, next) {
  const authStore = useAuthStore()
  const toast = useToast()
  
  // Skip authentication check for routes that don't require it
  if (!to.meta.requiresAuth) {
    return next()
  }
  
  // Check if the user is authenticated
  if (!authStore.isAuthenticated) {
    toast.error('Please log in to access this page')
    
    // Redirect to login with the original route as a redirect parameter
    return next({ 
      name: 'Login',
      query: { redirect: to.fullPath }
    })
  }
  
  // Get user role (supporting both string roles and numeric role_id)
  const userRole = getUserRole(authStore.user)
  
  // Check if the route has role restrictions
  if (to.meta.allowedRoles && to.meta.allowedRoles.length > 0) {
    // Check if user's role is allowed for this route
    if (!to.meta.allowedRoles.includes(userRole)) {
      console.error(`Access denied: User has role "${userRole}" but route requires one of: ${to.meta.allowedRoles.join(', ')}`)
      toast.error('You do not have permission to access this page')
      return next({ name: 'Unauthorized' })
    }
  }
  
  // For specific role check (single role requirement)
  if (to.meta.role && userRole !== to.meta.role) {
    console.error(`Access denied: User has role "${userRole}" but route requires "${to.meta.role}"`)
    toast.error('You do not have permission to access this page')
    return next({ name: 'Unauthorized' })
  }
  
  // User is authenticated and authorized, proceed
  next()
}

/**
 * Get user role as string, handling both string roles and numeric role IDs
 * @param {Object} user - User object from auth store
 * @returns {String|null} The user's role as a string
 */
export function getUserRole(user) {
  if (!user) return null
  
  // If user already has string role, use it
  if (typeof user.role === 'string') {
    return user.role
  }
  
  // If user has numeric role_id, map it to string role
  if (user.role_id !== undefined) {
    return ROLE_MAP[user.role_id] || null
  }
  
  return null
}

/**
 * Get the appropriate dashboard route based on user role
 * @param {Object} user - User object from auth store
 * @returns {String} The dashboard route path
 */
export function getDashboardRoute(user) {
  if (!user) return '/'
  
  const role = getUserRole(user)
  
  switch (role) {
    case 'superadmin':
      return '/superadmin/dashboard'
    case 'tenant':
      return '/tenant/dashboard'
    case 'devotee':
      return user.selectedTempleId 
        ? `/entity/${user.selectedTempleId}/devotee/dashboard`
        : '/devotee/temple-selection'
    case 'volunteer':
      return user.selectedTempleId
        ? `/entity/${user.selectedTempleId}/volunteer/dashboard`
        : '/volunteer/temple-selection'
    default:
      return '/'
  }
}