import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'

/**
 * Role-based middleware
 * Checks if the user has the required role to access the route
 */
export default function roleMiddleware(to, from, next) {
  const authStore = useAuthStore()
  const toast = useToast()
  
  // Skip role check if no role is required
  if (!to.meta.role) {
    return next()
  }
  
  // Get the required role from the route meta
  const requiredRole = to.meta.role
  
  // Check if the user has the required role
  if (authStore.user?.role !== requiredRole) {
    toast.error(`Access denied. This area is for ${requiredRole} users only.`)
    
    // Redirect to unauthorized page
    return next({ name: 'Unauthorized' })
  }
  
  // User has the required role, proceed
  next()
}