import superAdminService from '@/services/superadmin.service';

/**
 * Composable to provide access to Super Admin services
 */
export function useSuperAdminService() {
  const fetchUserById = (userId) => {
    return superAdminService.fetchUserById(userId);
  };

  const fetchAvailableTenants = (userId) => {
    return superAdminService.fetchAvailableTenants(userId);
  };

  const assignTenantsToUser = (userId, tenantIds) => {
    return superAdminService.assignTenantsToUser(userId, tenantIds);
  };

  // Add any additional methods from superadmin.service.js

  return {
    fetchUserById,
    fetchAvailableTenants,
    assignTenantsToUser
  };
}