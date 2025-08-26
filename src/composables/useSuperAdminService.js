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

  /**
   * Fetch reports for a specific tenant as a superadmin
   * @param {Object} params Report parameters including entityId and type
   * @returns {Promise} Promise resolving to report data
   */
  const fetchTenantReportPreview = (params) => {
    return superAdminService.fetchTenantReportPreview(params);
  };

  /**
   * Download reports for a specific tenant as a superadmin
   * @param {Object} params Report parameters including entityId, type, and format
   * @returns {Promise} Promise resolving to download result
   */
  const downloadTenantReport = (params) => {
    return superAdminService.downloadTenantReport(params);
  };

  /**
   * Get tenant entity details for reports
   * @param {string|number} tenantId Tenant ID
   * @returns {Promise} Promise resolving to tenant details
   */
  const getTenantEntityDetails = (tenantId) => {
    return superAdminService.getTenantEntityDetails(tenantId);
  };

  return {
    fetchUserById,
    fetchAvailableTenants,
    assignTenantsToUser,
    fetchTenantReportPreview,
    downloadTenantReport,
    getTenantEntityDetails
  };
}