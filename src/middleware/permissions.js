import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import { usePermissions } from '@/composables/usePermissions'

/**
 * Permissions middleware
 * Checks if the user has the required permissions to access the route
 */
export default function permissionsMiddleware(to, from, next) {
  const authStore = useAuthStore()
  const toast = useToast()
  const { hasPermission } = usePermissions()
  
  // Skip permissions check if no permissions are required
  if (!to.meta.permissions) {
    return next()
  }
  
  // Get the required permissions from the route meta
  const requiredPermissions = to.meta.permissions
  
  // Check if the user has at least one of the required permissions
  const hasRequiredPermission = Array.isArray(requiredPermissions)
    ? requiredPermissions.some(permission => hasPermission(permission))
    : hasPermission(requiredPermissions)
  
  if (!hasRequiredPermission) {
    toast.error('You do not have permission to access this page')
    
    // Redirect to unauthorized page
    return next({ name: 'Unauthorized' })
  }
  
  // User has the required permissions, proceed
  next()
}