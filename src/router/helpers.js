/**
 * Router helper functions
 */

/**
 * Get the post-login redirect URL based on user role or stored redirect
 * @param {Object} user - User object with role information
 * @param {string} storedRedirect - Previously stored redirect path
 * @returns {string} - Redirect path
 */
export function getPostLoginRedirect(user = null, storedRedirect = null) {
  // If there's a stored redirect path, use it
  if (storedRedirect && storedRedirect !== '/login') {
    return storedRedirect;
  }
  
  // Default redirects based on user role
  if (user) {
    switch (user.role) {
      case 'admin':
        return '/admin/dashboard';
      case 'temple-admin':
        return '/temple/dashboard';
      case 'user':
        return '/dashboard';
      case 'devotee':
        return '/devotee/dashboard';
      default:
        return '/dashboard';
    }
  }
  
  // Fallback to default dashboard
  return '/dashboard';
}