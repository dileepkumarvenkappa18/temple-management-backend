import { api } from './api';

export const auditLogService = {
  /**
   * Get audit logs with optional filters and pagination
   * @param {Object} filters - Filter parameters
   * @param {string} filters.action - Filter by action type
   * @param {string} filters.status - Filter by status (success/failure)
   * @param {string} filters.from_date - Start date for range filter (ISO format)
   * @param {string} filters.to_date - End date for range filter (ISO format)
   * @param {number} page - Page number for pagination
   * @param {number} limit - Number of items per page
   * @returns {Promise} Promise with paginated audit logs
   */
  getAuditLogs(filters = {}, page = 1, limit = 10) {
    return api.get('/api/v1/auditlogs', {
      params: {
        ...filters,
        page,
        limit
      }
    });
  },

  /**
   * Get details for a specific audit log
   * @param {string} id - Audit log ID
   * @returns {Promise} Promise with audit log details
   */
  getAuditLogDetails(id) {
    return api.get(`/api/v1/auditlogs/${id}`);
  }
};

export default auditLogService;