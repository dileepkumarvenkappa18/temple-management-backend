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
  async getAuditLogs(filters = {}, page = 1, limit = 10) {
    console.log('Fetching audit logs with params:', { filters, page, limit });
    
    const params = {
      ...filters,
      page,
      limit
    };
    
    try {
      const response = await api.get('/api/v1/auditlogs', { params });
      console.log('Audit logs response:', response);
      return response;
    } catch (error) {
      console.error('Error in auditLogService.getAuditLogs:', error);
      throw error;
    }
  },

  /**
   * Get details for a specific audit log
   * @param {string|number} id - Audit log ID
   * @returns {Promise} Promise with audit log details
   */
  async getAuditLogDetails(id) {
    console.log('Fetching audit log details:', id);
    try {
      const response = await api.get(`/api/v1/auditlogs/${id}`);
      console.log('Audit log details response:', response);
      return response;
    } catch (error) {
      console.error('Error in auditLogService.getAuditLogDetails:', error);
      throw error;
    }
  }
};

export default auditLogService;