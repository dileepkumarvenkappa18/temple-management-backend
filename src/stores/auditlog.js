import { defineStore } from 'pinia';
import { auditLogService } from '@/services/auditlog.service';
import { formatDate } from '@/utils/date';

// Define all possible audit log actions from our system
const auditLogActions = [
  // Superadmin Module
  'TENANT_APPROVED', 'TENANT_REJECTED', 'ENTITY_APPROVED', 'ENTITY_REJECTED',
  'USER_CREATED', 'USER_UPDATED', 'USER_DELETED', 'USER_STATUS_UPDATED',
  'ROLE_CREATED', 'ROLE_UPDATED', 'ROLE_STATUS_UPDATED',
  
  // Entity Module
  'TEMPLE_CREATED', 'TEMPLE_UPDATED', 'TEMPLE_CREATE_FAILED', 'TEMPLE_UPDATE_FAILED',
  
  // Events Module
  'EVENT_CREATED', 'EVENT_UPDATED', 'EVENT_DELETED',
  
  // Seva Module
  'SEVA_CREATED', 'SEVA_UPDATED', 'SEVA_BOOKED', 'SEVA_BOOKING_APPROVED', 
  'SEVA_BOOKING_REJECTED', 'SEVA_BOOKING_FAILED',
  
  // Donations Module
  'DONATION_INITIATED', 'DONATION_SUCCESS', 'DONATION_FAILED', 
  'DONATION_VERIFICATION_FAILED',
  
  // Notification Module
  'TEMPLATE_CREATED', 'TEMPLATE_UPDATED', 'TEMPLATE_DELETED',
  'EMAIL_SENT', 'SMS_SENT', 'WHATSAPP_SENT',
  
  // User Profile Module
  'PROFILE_CREATED', 'PROFILE_UPDATED', 'DEVOTEE_JOINED_TEMPLE', 
  'VOLUNTEER_JOINED_TEMPLE',
  
  // Reports Module
  'DEVOTEE_BIRTHDAYS_REPORT_VIEWED', 'DEVOTEE_BIRTHDAYS_REPORT_DOWNLOADED',
  'DEVOTEE_BIRTHDAYS_REPORT_DOWNLOAD_FAILED', 'TEMPLE_REGISTER_REPORT_VIEWED',
  'TEMPLE_REGISTER_REPORT_DOWNLOADED', 'TEMPLE_REGISTER_REPORT_DOWNLOAD_FAILED',
  'TEMPLE_ACTIVITIES_REPORT_VIEWED', 'TEMPLE_ACTIVITIES_REPORT_DOWNLOADED',
  'TEMPLE_ACTIVITIES_REPORT_DOWNLOAD_FAILED',
  
  // System Actions
  'LOGIN_SUCCESS', 'LOGIN_FAILED', 'LOGOUT',
  'PASSWORD_RESET_REQUESTED', 'PASSWORD_RESET_SUCCESS', 'PASSWORD_RESET_FAILED'
];

export const useAuditLogStore = defineStore('auditLog', {
  state: () => ({
    logs: [],
    selectedLog: null,
    isLoading: false,
    isDetailLoading: false,
    error: null,
    
    // Pagination
    currentPage: 1,
    totalPages: 1,
    limit: 10,
    total: 0,
    
    // Filters
    filters: {
      action: '',
      status: '',
      from_date: '',
      to_date: ''
    },
    
    // Available actions for filtering
    availableActions: auditLogActions
  }),
  
  getters: {
    // Format logs for display
    formattedLogs: (state) => {
      return state.logs.map(log => ({
        ...log,
        formattedDate: formatDate(log.createdAt, 'DD MMM YYYY, HH:mm:ss')
      }));
    },
    
    // Group actions by module for better organization in dropdown
    groupedActions: (state) => {
      const grouped = {
        'Superadmin': state.availableActions.filter(a => 
          a.includes('TENANT_') || a.includes('ENTITY_') || 
          a.includes('USER_') || a.includes('ROLE_')),
          
        'Entity': state.availableActions.filter(a => a.includes('TEMPLE_')),
        
        'Events': state.availableActions.filter(a => a.includes('EVENT_')),
        
        'Seva': state.availableActions.filter(a => a.includes('SEVA_')),
        
        'Donations': state.availableActions.filter(a => a.includes('DONATION_')),
        
        'Notifications': state.availableActions.filter(a => 
          a.includes('TEMPLATE_') || a.includes('_SENT')),
          
        'User Profile': state.availableActions.filter(a => 
          a.includes('PROFILE_') || a.includes('_JOINED_')),
          
        'Reports': state.availableActions.filter(a => a.includes('REPORT_')),
        
        'System': state.availableActions.filter(a => 
          a.includes('LOGIN_') || a.includes('LOGOUT') || a.includes('PASSWORD_RESET_'))
      };
      
      return grouped;
    }
  },
  
  actions: {
    /**
     * Fetch audit logs with current filters and pagination
     */
    async fetchAuditLogs() {
      this.isLoading = true;
      this.error = null;
      
      try {
        const response = await auditLogService.getAuditLogs(
          this.filters,
          this.currentPage,
          this.limit
        );
        
        this.logs = response.data.logs;
        this.totalPages = response.data.total_pages;
        this.total = response.data.total;
        
      } catch (error) {
        this.error = error.message || 'Failed to fetch audit logs';
        console.error('Error fetching audit logs:', error);
      } finally {
        this.isLoading = false;
      }
    },
    
    /**
     * Fetch details for a specific audit log
     * @param {string} id - Audit log ID
     */
    async fetchAuditLogDetails(id) {
      this.isDetailLoading = true;
      this.error = null;
      
      try {
        const response = await auditLogService.getAuditLogDetails(id);
        this.selectedLog = response.data;
      } catch (error) {
        this.error = error.message || 'Failed to fetch audit log details';
        console.error('Error fetching audit log details:', error);
      } finally {
        this.isDetailLoading = false;
      }
    },
    
    /**
     * Set filter values
     * @param {Object} filterValues - New filter values
     */
    setFilters(filterValues) {
      this.filters = { ...this.filters, ...filterValues };
      this.currentPage = 1; // Reset to first page when filters change
      this.fetchAuditLogs();
    },
    
    /**
     * Reset all filters to default values
     */
    resetFilters() {
      this.filters = {
        action: '',
        status: '',
        from_date: '',
        to_date: ''
      };
      this.currentPage = 1;
      this.fetchAuditLogs();
    },
    
    /**
     * Change pagination page
     * @param {number} page - New page number
     */
    setPage(page) {
      this.currentPage = page;
      this.fetchAuditLogs();
    },
    
    /**
     * Change items per page
     * @param {number} limit - New limit value
     */
    setLimit(limit) {
      this.limit = limit;
      this.currentPage = 1; // Reset to first page when limit changes
      this.fetchAuditLogs();
    },
    
    /**
     * Clear selected log detail
     */
    clearSelectedLog() {
      this.selectedLog = null;
    }
  }
});

export default useAuditLogStore;