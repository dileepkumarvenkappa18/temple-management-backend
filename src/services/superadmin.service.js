import api, { endpoints } from './api.js'

class SuperAdminService {
  constructor() {
    this.baseURL = '/api/v1/superadmin'
  }

  // ==========================================
  // COUNT ENDPOINTS - EXISTING
  // ==========================================

  /**
   * Get tenant approval counts (pending, approved, rejected)
   * Endpoint: GET /api/v1/superadmin/tenant-approval-count
   * Response: {"approved": 2, "pending": 0, "rejected": 1}
   */
  async getTenantApprovalCounts() {
    try {
      console.log('Service: Fetching tenant approval counts...')
      
      const response = await api.get(`${this.baseURL}/tenant-approval-count`)
      console.log('Service: Tenant counts response:', response)
      
      // Handle the response format
      if (response && typeof response === 'object') {
        return {
          success: true,
          data: {
            pending: response.pending || 0,
            approved: response.approved || 0,
            rejected: response.rejected || 0
          },
          message: 'Tenant approval counts fetched successfully'
        }
      }
      
      return {
        success: false,
        data: { pending: 0, approved: 0, rejected: 0 },
        message: 'Invalid response format for tenant counts'
      }
    } catch (error) {
      console.error('Service: Error fetching tenant approval counts:', error)
      
      // Fallback to mock data for development
      if (error.response?.status === 404) {
        console.warn('Tenant count endpoint not found, using mock data')
        return {
          success: true,
          data: { pending: 2, approved: 3, rejected: 1 }, // Mock data
          message: 'Mock tenant counts loaded (API endpoint not available)'
        }
      }
      
      return {
        success: false,
        data: { pending: 0, approved: 0, rejected: 0 },
        message: error.message || 'Failed to fetch tenant approval counts'
      }
    }
  }

  /**
   * Get temple approval counts (pending_approvals, active_temples, rejected, total_users)
   * Endpoint: GET /api/v1/superadmin/temple-approval-count
   * Response: {"pending_approval": 0, "active_temples": 1, "rejected": 0, "total_devotees": 2}
   */
  async getTempleApprovalCounts() {
    try {
      console.log('Service: Fetching temple approval counts...')
      
      const response = await api.get(`${this.baseURL}/temple-approval-count`)
      console.log('Service: Temple counts response:', response)
      
      // Handle the response format
      if (response && typeof response === 'object') {
        return {
          success: true,
          data: {
            pendingApprovals: response.pending_approval || response.pending_approvals || 0,
            activeTemples: response.active_temples || 0,
            rejectedTemples: response.rejected || 0,
            totalUsers: response.total_devotees || response.total_users || 0,
            newThisWeek: response.new_this_week || 0 // May not be in API response
          },
          message: 'Temple approval counts fetched successfully'
        }
      }
      
      return {
        success: false,
        data: {
          pendingApprovals: 0,
          activeTemples: 0,
          rejectedTemples: 0,
          totalUsers: 0,
          newThisWeek: 0
        },
        message: 'Invalid response format for temple counts'
      }
    } catch (error) {
      console.error('Service: Error fetching temple approval counts:', error)
      
      // Fallback to mock data for development
      if (error.response?.status === 404) {
        console.warn('Temple count endpoint not found, using mock data')
        return {
          success: true,
          data: {
            pendingApprovals: 5,
            activeTemples: 32,
            rejectedTemples: 3,
            totalUsers: 178,
            newThisWeek: 4
          },
          message: 'Mock temple counts loaded (API endpoint not available)'
        }
      }
      
      return {
        success: false,
        data: {
          pendingApprovals: 0,
          activeTemples: 0,
          rejectedTemples: 0,
          totalUsers: 0,
          newThisWeek: 0
        },
        message: error.message || 'Failed to fetch temple approval counts'
      }
    }
  }

  // ==========================================
  // TENANT MANAGEMENT - EXISTING
  // ==========================================

  async getPendingTenants() {
    try {
      console.log('Service: Fetching pending tenants...')
      
      // Updated to use query parameter instead of /pending path
      const response = await api.get(`${this.baseURL}/tenants?status=pending`)
      
      console.log('Service: Raw API response:', response)
      
      // The backend returns pending_tenants wrapped in the response
      return {
        success: true,
        data: response.pending_tenants || response.data || response || [], 
        message: 'Tenants fetched successfully'
      }
    } catch (error) {
      console.error('Service: Error fetching pending tenants:', error)
      
      // If this is a 404, we might need to fall back to mock data
      if (error.response?.status === 404) {
        console.warn('Endpoint not found, using mock data for demonstration')
        return {
          success: true,
          data: this.getMockPendingTenants(),
          message: 'Mock tenant data loaded (API endpoint not available)'
        }
      }
      
      return {
        success: false,
        data: null,
        message: error.message || 'Failed to fetch pending tenants'
      }
    }
  }

  async getAllTenants(filters = {}) {
    try {
      const params = new URLSearchParams()
      if (filters.status) params.append('status', filters.status)
      if (filters.search) params.append('search', filters.search)
      if (filters.page) params.append('page', filters.page)
      if (filters.limit) params.append('limit', filters.limit)
      if (filters.sortBy) params.append('sortBy', filters.sortBy)
      if (filters.sortOrder) params.append('sortOrder', filters.sortOrder)

      const response = await api.get(`${this.baseURL}/tenants?${params}`)
      
      if (Array.isArray(response.data)) {
        return {
          success: true,
          data: response.data,
          pagination: { total: response.data.length },
          message: 'Tenants fetched successfully'
        }
      } else if (response.data && Array.isArray(response.data.tenants)) {
        return {
          success: true,
          data: response.data.tenants,
          pagination: response.data.pagination || {},
          message: 'Tenants fetched successfully'
        }
      } else if (response.data && Array.isArray(response.data.data)) {
        return {
          success: true,
          data: response.data.data,
          pagination: response.data.pagination || {},
          message: 'Tenants fetched successfully'
        }
      } else {
        return {
          success: false,
          data: [],
          message: 'Unexpected API response format'
        }
      }
    } catch (error) {
      console.error('Error fetching all tenants:', error)
      
      // Use mock data as fallback
      if (error.response?.status === 404) {
        console.warn('Endpoint not found, using mock data for demonstration')
        return {
          success: true,
          data: this.getMockAllTenants(),
          message: 'Mock tenant data loaded (API endpoint not available)'
        }
      }
      
      return {
        success: false,
        data: [],
        message: error.message || 'Failed to fetch tenants'
      }
    }
  }

  async approveTenant(tenantId) {
    try {
      console.log(`Approving tenant ${tenantId}...`)
      
      // Use the correct endpoint from api.js
      const payload = {
        status: "APPROVED"
      }
      
      // Updated to use the actual endpoint from backend
      const response = await api.patch(`${this.baseURL}/tenants/${tenantId}/approval`, payload)
      console.log('Tenant approval response:', response)
      
      return {
        success: true,
        data: response,
        message: 'Tenant approved successfully'
      }
    } catch (error) {
      console.error('Error approving tenant:', error)
      
      // If we're in demo mode, simulate success
      if (error.response?.status === 404) {
        console.warn('Endpoint not found, simulating successful approval')
        
        // Update our mock data to reflect the approval
        this.updateMockTenantStatus(tenantId, 'approved')
        
        return {
          success: true,
          data: { status: 'approved', message: 'Tenant approved successfully (mock)' },
          message: 'Tenant approved successfully (mock)'
        }
      }
      
      return {
        success: false,
        data: null,
        message: error.message || 'Failed to approve tenant'
      }
    }
  }

  async rejectTenant(tenantId, data) {
    try {
      console.log(`Rejecting tenant ${tenantId}...`)
      
      // Use the correct endpoint from api.js
      const payload = {
        status: "REJECTED",
        reason: data.reason || data.notes || ''
      }

      // Updated to use the actual endpoint from backend
      const response = await api.patch(`${this.baseURL}/tenants/${tenantId}/approval`, payload)
      console.log('Tenant rejection response:', response)
      
      return {
        success: true,
        data: response,
        message: 'Tenant rejected successfully'
      }
    } catch (error) {
      console.error('Error rejecting tenant:', error)
      
      // If we're in demo mode, simulate success
      if (error.response?.status === 404) {
        console.warn('Endpoint not found, simulating successful rejection')
        
        // Update our mock data to reflect the rejection
        this.updateMockTenantStatus(tenantId, 'rejected')
        
        return {
          success: true,
          data: { 
            status: 'rejected', 
            reason: data.reason || data.notes,
            message: 'Tenant rejected successfully (mock)' 
          },
          message: 'Tenant rejected successfully (mock)'
        }
      }
      
      return {
        success: false,
        data: null,
        message: error.message || 'Failed to reject tenant'
      }
    }
  }

  async getTenantDetails(tenantId) {
    try {
      const response = await api.get(`${this.baseURL}/tenants/${tenantId}`)
      
      if (response.data && response.data.tenant) {
        return {
          success: true,
          data: response.data.tenant,
          temples: response.data.temples || [],
          message: 'Tenant details fetched successfully'
        }
      } else if (response.data) {
        return {
          success: true,
          data: response.data,
          temples: response.data.temples || [],
          message: 'Tenant details fetched successfully'
        }
      }
      
      return {
        success: false,
        data: null,
        message: 'Unexpected API response format'
      }
    } catch (error) {
      console.error('Error fetching tenant details:', error)
      return {
        success: false,
        data: null,
        message: error.message || 'Failed to fetch tenant details'
      }
    }
  }

  // ==========================================
  // TEMPLE MANAGEMENT - EXISTING
  // ==========================================

  async getPendingEntities() {
    try {
      console.log('Fetching pending entities from API...')
      
      // Updated to use query parameter instead of /pending path
      const response = await api.get(`${this.baseURL}/entities?status=PENDING`)
      console.log('API response for pending entities:', response)
      
      // Handle the response format properly based on the backend structure
      if (response && response.pending_entities !== undefined) {
        // Backend returns {pending_entities: [...]}
        return {
          success: true,
          data: response.pending_entities || [],
          total: (response.pending_entities || []).length,
          message: 'Pending entities fetched successfully'
        }
      } else if (Array.isArray(response)) {
        return {
          success: true,
          data: response,
          total: response.length,
          message: 'Pending entities fetched successfully'
        }
      } else if (response && response.data && Array.isArray(response.data)) {
        return {
          success: true,
          data: response.data,
          total: response.total || response.data.length,
          message: 'Pending entities fetched successfully'
        }
      } else if (response && Array.isArray(response.entities)) {
        return {
          success: true,
          data: response.entities,
          total: response.total || response.entities.length,
          message: 'Pending entities fetched successfully'
        }
      } else {
        console.warn('API returned unexpected data format:', response)
        // Return empty array as fallback
        return {
          success: true,
          data: [],
          total: 0,
          message: 'No pending entities found'
        }
      }
    } catch (error) {
      console.error('Error fetching pending entities:', error)
      
      // If we're in demo mode, use mock data
      if (error.response?.status === 404) {
        return {
          success: true,
          data: this.getMockPendingEntities(),
          total: this.getMockPendingEntities().length,
          message: 'Mock entities loaded (API endpoint not available)'
        }
      }
      
      return {
        success: false,
        data: [],
        message: error.message || 'Failed to fetch pending entities'
      }
    }
  }

  async approveEntity(entityId, data) {
    try {
      console.log(`Approving entity ${entityId}...`)
      
      const payload = {
        status: "APPROVED",
        notes: data?.notes || ''
      }
      
      // Updated to use the actual endpoint from backend
      const response = await api.patch(`${this.baseURL}/entities/${entityId}/approval`, payload)
      console.log('Entity approval response:', response)
      
      return {
        success: true,
        data: response,
        message: 'Entity approved successfully'
      }
    } catch (error) {
      console.error('Error approving entity:', error)
      
      if (error.response?.status === 404) {
        console.warn('Endpoint not found, simulating successful approval')
        return {
          success: true,
          data: { status: 'approved', message: 'Entity approved successfully (mock)' },
          message: 'Entity approved successfully (mock)'
        }
      }
      
      return {
        success: false,
        data: null,
        message: error.message || 'Failed to approve entity'
      }
    }
  }

  async rejectEntity(entityId, data) {
    try {
      console.log(`Rejecting entity ${entityId}...`)
      
      const payload = {
        status: "REJECTED",
        reason: data.reason || data.notes || ''
      }
      
      // Updated to use the actual endpoint from backend
      const response = await api.patch(`${this.baseURL}/entities/${entityId}/approval`, payload)
      console.log('Entity rejection response:', response)
      
      return {
        success: true,
        data: response,
        message: 'Entity rejected successfully'
      }
    } catch (error) {
      console.error('Error rejecting entity:', error)
      
      if (error.response?.status === 404) {
        console.warn('Endpoint not found, simulating successful rejection')
        return {
          success: true,
          data: { 
            status: 'rejected', 
            reason: data.reason || data.notes,
            message: 'Entity rejected successfully (mock)' 
          },
          message: 'Entity rejected successfully (mock)'
        }
      }
      
      return {
        success: false,
        data: null,
        message: error.message || 'Failed to reject entity'
      }
    }
  }

  // ==========================================
  // USER MANAGEMENT - NEW
  // ==========================================

  /**
   * Get all user roles 
   * Endpoint: GET /api/v1/superadmin/user-roles
   */
  async getUserRoles() {
    try {
      console.log('Service: Fetching user roles...')
      const response = await api.get(`${this.baseURL}/user-roles`)
      console.log('Service: User roles response:', response)
      
      if (response && response.data && Array.isArray(response.data)) {
        return {
          success: true,
          data: response.data,
          message: 'User roles fetched successfully'
        }
      } else if (Array.isArray(response)) {
        return {
          success: true,
          data: response,
          message: 'User roles fetched successfully'
        }
      }
      
      return {
        success: false,
        data: [],
        message: 'Invalid response format for user roles'
      }
    } catch (error) {
      console.error('Service: Error fetching user roles:', error)
      return {
        success: false,
        data: [],
        message: error.message || 'Failed to fetch user roles'
      }
    }
  }

  /**
   * Get all users with pagination and filters
   * Endpoint: GET /api/v1/superadmin/users
   */
  async getUsers(filters = {}) {
    try {
      console.log('Service: Fetching users...')
      
      const params = new URLSearchParams()
      if (filters.limit) params.append('limit', filters.limit)
      if (filters.page) params.append('page', filters.page)
      if (filters.search) params.append('search', filters.search)
      if (filters.role) params.append('role', filters.role)
      if (filters.status) params.append('status', filters.status)

      const response = await api.get(`${this.baseURL}/users?${params}`)
      console.log('Service: Users response:', response)
      
      if (response && response.data && Array.isArray(response.data)) {
        return {
          success: true,
          data: response.data,
          total: response.total || response.data.length,
          pagination: {
            page: response.page || filters.page || 1,
            limit: response.limit || filters.limit || 10,
            total: response.total || response.data.length
          },
          message: 'Users fetched successfully'
        }
      } else if (Array.isArray(response)) {
        return {
          success: true,
          data: response,
          total: response.length,
          pagination: {
            page: filters.page || 1,
            limit: filters.limit || 10,
            total: response.length
          },
          message: 'Users fetched successfully'
        }
      }
      
      return {
        success: false,
        data: [],
        message: 'Invalid response format for users'
      }
    } catch (error) {
      console.error('Service: Error fetching users:', error)
      return {
        success: false,
        data: [],
        message: error.message || 'Failed to fetch users'
      }
    }
  }

  /**
   * Get user by ID
   * Endpoint: GET /api/v1/superadmin/users/:id
   */
  async getUserById(userId) {
    try {
      console.log(`Service: Fetching user ${userId}...`)
      const response = await api.get(`${this.baseURL}/users/${userId}`)
      console.log('Service: User details response:', response)
      
      if (response && response.data) {
        return {
          success: true,
          data: response.data,
          message: 'User details fetched successfully'
        }
      }
      
      return {
        success: false,
        data: null,
        message: 'Invalid response format for user details'
      }
    } catch (error) {
      console.error('Service: Error fetching user details:', error)
      return {
        success: false,
        data: null,
        message: error.message || 'Failed to fetch user details'
      }
    }
  }

  /**
   * Create new user
   * Endpoint: POST /api/v1/superadmin/users
   */
  async createUser(userData) {
    try {
      console.log('Service: Creating user...', userData)
      const response = await api.post(`${this.baseURL}/users`, userData)
      console.log('Service: Create user response:', response)
      
      return {
        success: true,
        data: response,
        message: 'User created successfully'
      }
    } catch (error) {
      console.error('Service: Error creating user:', error)
      return {
        success: false,
        data: null,
        message: error.response?.data?.error || error.message || 'Failed to create user'
      }
    }
  }

  /**
   * Update user
   * Endpoint: PUT /api/v1/superadmin/users/:id
   */
  async updateUser(userId, userData) {
    try {
      console.log(`Service: Updating user ${userId}...`, userData)
      const response = await api.put(`${this.baseURL}/users/${userId}`, userData)
      console.log('Service: Update user response:', response)
      
      return {
        success: true,
        data: response,
        message: 'User updated successfully'
      }
    } catch (error) {
      console.error('Service: Error updating user:', error)
      return {
        success: false,
        data: null,
        message: error.response?.data?.error || error.message || 'Failed to update user'
      }
    }
  }

  /**
   * Update user status (activate/deactivate)
   * Endpoint: PATCH /api/v1/superadmin/users/:id/status
   */
  async updateUserStatus(userId, status) {
    try {
      console.log(`Service: Updating user ${userId} status to ${status}...`)
      const response = await api.patch(`${this.baseURL}/users/${userId}/status`, { status })
      console.log('Service: Update user status response:', response)
      
      return {
        success: true,
        data: response,
        message: 'User status updated successfully'
      }
    } catch (error) {
      console.error('Service: Error updating user status:', error)
      return {
        success: false,
        data: null,
        message: error.response?.data?.error || error.message || 'Failed to update user status'
      }
    }
  }

  // ==========================================
  // USER-TENANT ASSIGNMENT - NEW
  // ==========================================

  /**
   * Fetch available tenants for assignment to a user
   * @param {string|number} userId - The user ID
   * @returns {Promise} - API response with tenants data
   */
  async getAvailableTenants(userId) {
    try {
      console.log(`Service: Fetching available tenants for user ${userId}...`)
      const response = await api.get(`${this.baseURL}/users/${userId}/available-tenants`)
      console.log('Service: Available tenants response:', response)
      
      if (response && response.data && Array.isArray(response.data)) {
        return {
          success: true,
          data: response.data,
          message: 'Available tenants fetched successfully'
        }
      } else if (Array.isArray(response)) {
        return {
          success: true,
          data: response,
          message: 'Available tenants fetched successfully'
        }
      }
      
      // Fallback to mock data for development
      if (response === null || (typeof response === 'object' && Object.keys(response).length === 0)) {
        console.warn('Empty response, using mock data')
        return {
          success: true,
          data: this.getMockAvailableTenants(),
          message: 'Mock available tenants data loaded (API returned empty)'
        }
      }
      
      return {
        success: false,
        data: [],
        message: 'Invalid response format for available tenants'
      }
    } catch (error) {
      console.error('Service: Error fetching available tenants:', error)
      
      // Fallback to mock data for development
      if (error.response?.status === 404) {
        console.warn('Available tenants endpoint not found, using mock data')
        return {
          success: true,
          data: this.getMockAvailableTenants(),
          message: 'Mock available tenants loaded (API endpoint not available)'
        }
      }
      
      return {
        success: false,
        data: [],
        message: error.message || 'Failed to fetch available tenants'
      }
    }
  }

  /**
   * Assign tenants to a user
   * @param {string|number} userId - The user ID
   * @param {Array} tenantIds - Array of tenant IDs to assign
   * @returns {Promise} - API response
   */
  async assignTenantsToUser(userId, tenantIds) {
    try {
      console.log(`Service: Assigning tenants to user ${userId}...`, tenantIds)
      const response = await api.post(`${this.baseURL}/users/${userId}/assign-tenants`, {
        tenantIds: tenantIds
      })
      console.log('Service: Assign tenants response:', response)
      
      return {
        success: true,
        data: response,
        message: 'Tenants assigned successfully'
      }
    } catch (error) {
      console.error('Service: Error assigning tenants:', error)
      
      // Fallback for development/testing
      if (error.response?.status === 404) {
        console.warn('Assign tenants endpoint not found, simulating successful assignment')
        return {
          success: true,
          data: { 
            userId: userId,
            tenantIds: tenantIds,
            message: 'Tenants assigned successfully (mock)'
          },
          message: 'Tenants assigned successfully (mock)'
        }
      }
      
      return {
        success: false,
        data: null,
        message: error.response?.data?.error || error.message || 'Failed to assign tenants'
      }
    }
  }

  // ==========================================
  // PASSWORD RESET - NEW
  // ==========================================

  /**
   * Search for a user by email
   * Endpoint: GET /api/v1/superadmin/users/search?email=user@example.com
   */
  async searchUserByEmail(email) {
    try {
      console.log(`Service: Searching for user with email ${email}...`)
      const response = await api.get(`${this.baseURL}/users/search?email=${encodeURIComponent(email)}`)
      console.log('Service: User search response:', response)
      
      const user = response.data?.data || response.data || response;

      if (user) {
        return {
          success: true,
          data: user,
          message: 'User found successfully'
        }
      }
      
      return {
        success: false,
        data: null,
        message: 'User not found'
      }
    } catch (error) {
      console.error('Service: Error searching for user:', error)
      
      // For development/testing only
      if (error.response?.status === 404 && email.includes('@')) {
        console.warn('Endpoint not found, using mock data for demonstration')
        return {
          success: true,
          data: {
            id: 1,
            fullName: 'John Doe',
            email: email,
            role: {
              id: 2,
              roleName: 'tenant'
            }
          },
          message: 'Mock user found (API endpoint not available)'
        }
      }
      
      return {
        success: false,
        data: null,
        message: error.response?.data?.error || error.message || 'User not found'
      }
    }
  }

  /**
   * Reset user password
   * Endpoint: POST /api/v1/superadmin/users/:id/reset-password
   */
  async resetUserPassword(userId, newPassword) {
    try {
      console.log(`Service: Resetting password for user ${userId}...`);
      
      // Add sendEmail flag explicitly to ensure email notification is sent
      const response = await api.post(`${this.baseURL}/users/${userId}/reset-password`, { 
        password: newPassword,
        sendEmail: true
      });
      
      console.log('Service: Password reset response:', response);
      
      return {
        success: true,
        message: 'Password reset successfully. A notification email has been sent to the user.'
      };
    } catch (error) {
      console.error('Service: Error resetting password:', error);
      
      // For development/testing only
      if (error.response?.status === 404) {
        console.warn('Endpoint not found, simulating successful password reset');
        return {
          success: true,
          message: 'Password reset successfully (mock). In production, an email would be sent to the user.'
        };
      }
      
      return {
        success: false,
        message: error.response?.data?.error || error.message || 'Failed to reset password'
      };
    }
  }

  // ==========================================
  // ANALYTICS & DASHBOARD - UPDATED
  // ==========================================

  async getDashboardStats(dateRange = {}) {
    try {
      const params = new URLSearchParams()
      if (dateRange.startDate) params.append('startDate', dateRange.startDate)
      if (dateRange.endDate) params.append('endDate', dateRange.endDate)
      if (dateRange.period) params.append('period', dateRange.period)

      console.log('Fetching dashboard stats with params:', params.toString())
      const response = await api.get(`${this.baseURL}/dashboard?${params}`)
      console.log('Dashboard stats response:', response)
      
      if (response && typeof response === 'object' && !Array.isArray(response)) {
        // If the API returns the stats directly
        if (response.pendingApprovals !== undefined || 
            response.activeTemples !== undefined ||
            response.stats !== undefined) {
          return {
            success: true,
            data: response.stats || response,
            message: 'Dashboard stats fetched successfully'
          }
        }
      }
      
      return {
        success: false,
        data: null,
        message: 'Unexpected API response format or no stats available'
      }
    } catch (error) {
      console.error('Error fetching dashboard stats:', error)
      return {
        success: false,
        data: null,
        message: error.message || 'Failed to fetch dashboard stats'
      }
    }
  }

  async getSystemStats(dateRange = {}) {
    // Updated to use the new temple count endpoint
    return this.getTempleApprovalCounts();
  }

  async getRoles() {
    try {
      console.log('Service: Fetching roles...');
      const response = await api.get(`${this.baseURL}/roles`);

      // Log the raw response to confirm its structure
      console.log('Service: Roles response:', response);

      // Check if the response is a valid array
      if (Array.isArray(response)) {
        return {
          success: true,
          data: response, // Use the response array directly
          message: 'Roles fetched successfully'
        };
      }

      // Fallback for an unexpected format
      console.warn('Service: API returned a non-array response for roles:', response);
      return {
        success: false,
        data: [],
        message: 'Unexpected API response format for roles'
      };
    } catch (error) {
      console.error('Service: Error fetching roles:', error);
      return {
        success: false,
        data: [],
        message: error.message || 'Failed to fetch roles'
      };
    }
  }
  
  async createRole(roleData) {
    try {
      console.log('Service: Creating new role...', roleData);
      const response = await api.post(`${this.baseURL}/roles`, roleData);
      console.log('Service: Create role response:', response);
      return {
        success: true,
        data: response.data,
        message: 'Role created successfully'
      };
    } catch (error) {
      console.error('Service: Error creating role:', error);
      return {
        success: false,
        data: null,
        message: error.response?.data?.message || 'Failed to create role'
      };
    }
  }
  
  async updateRole(id, roleData) {
    try {
      console.log(`Service: Updating role with ID ${id}...`, roleData);
      const response = await api.put(`${this.baseURL}/roles/${id}`, roleData);
      console.log('Service: Update role response:', response);
      return {
        success: true,
        data: response.data,
        message: 'Role updated successfully'
      };
    } catch (error) {
      console.error('Service: Error updating role:', error);
      return {
        success: false,
        data: null,
        message: error.response?.data?.message || 'Failed to update role'
      };
    }
  }
  
  async deleteRole(id) {
    try {
      console.log(`Service: Deleting role with ID ${id}...`);
      const response = await api.delete(`${this.baseURL}/roles/${id}`);
      console.log('Service: Delete role response:', response);
      return {
        success: true,
        data: response.data,
        message: 'Role deleted successfully'
      };
    } catch (error) {
      console.error('Service: Error deleting role:', error);
      return {
        success: false,
        data: null,
        message: error.response?.data?.message || 'Failed to delete role'
      };
    }
  }
  
  // ==========================================
  // ACTIVITIES (MOCK)
  // ==========================================
  
  /**
   * Get recent activities - method to support the admin dashboard
   * @param {number} limit - Maximum number of activities to return
   * @returns {Promise<Object>} Recent activities
   */
  async getRecentActivities(limit = 10) {
    try {
      // Try to call the actual endpoint (which may not exist yet)
      try {
        const response = await api.get(`${this.baseURL}/activities?limit=${limit}`)
        
        // Check if valid response format
        if (Array.isArray(response)) {
          return {
            success: true,
            data: response,
            message: 'Recent activities fetched successfully'
          }
        } else if (response && Array.isArray(response.activities)) {
          return {
            success: true,
            data: response.activities,
            message: 'Recent activities fetched successfully'
          }
        }
      } catch (apiError) {
        console.warn('Activities endpoint not yet implemented:', apiError)
        
        // Return mock data instead
        return {
          success: true,
          data: this.getMockActivities(limit),
          message: 'Mock activities generated (endpoint not available)'
        }
      }
      
      return {
        success: false,
        data: [],
        message: 'Unexpected API response format for activities'
      }
    } catch (error) {
      console.error('Error in getRecentActivities:', error)
      return {
        success: false,
        data: [],
        message: error.message || 'Failed to fetch recent activities'
      }
    }
  }

  /**
 * Get all temples for a specific tenant
 * @param {string|number} tenantId - Tenant ID
 * @returns {Promise} - Promise with temples data
 */
async getTemplesByTenant(tenantId) {
  try {
    console.log(`Service: Fetching temples for tenant ${tenantId}...`)
    
    const response = await api.get(`${this.baseURL}/tenants/${tenantId}/temples`)
    console.log('Service: Temples by tenant response:', response)
    
    // Handle various response formats
    let temples = []
    
    if (response && Array.isArray(response)) {
      temples = response
    } else if (response && Array.isArray(response.data)) {
      temples = response.data
    } else if (response && Array.isArray(response.temples)) {
      temples = response.temples
    } else if (response && response.data && Array.isArray(response.data.temples)) {
      temples = response.data.temples
    }
    
    return {
      success: true,
      data: temples,
      message: 'Temples fetched successfully'
    }
  } catch (error) {
    console.error('Service: Error fetching temples by tenant:', error)
    
    // Fallback to mock data for development
    if (error.response?.status === 404) {
      console.warn('Temples by tenant endpoint not found, using mock data')
      return {
        success: true,
        data: this.getMockTemplesByTenant(tenantId),
        message: 'Mock temples loaded (API endpoint not available)'
      }
    }
    
    return {
      success: false,
      data: [],
      message: error.message || 'Failed to fetch temples for this tenant'
    }
  }
}

/**
 * Generate mock temples for a tenant (for development/demo)
 * @param {string|number} tenantId - Tenant ID
 * @returns {Array} - Array of mock temples
 */
getMockTemplesByTenant(tenantId) {
  const mockTemplesMap = {
    // First tenant
    '1': [
      {
        id: 101,
        name: 'Sri Krishna Temple',
        address: '123 Temple Street, Bengaluru',
        status: 'active'
      },
      {
        id: 102,
        name: 'Ganesha Temple',
        address: '456 Divine Road, Bengaluru',
        status: 'active'
      }
    ],
    // Second tenant
    '2': [
      {
        id: 201,
        name: 'Hanuman Temple',
        address: '789 Devotee Lane, Chennai',
        status: 'active'
      }
    ],
    // Third tenant
    '3': [
      {
        id: 301,
        name: 'Shiva Temple',
        address: '101 Spiritual Avenue, Mysore',
        status: 'active'
      },
      {
        id: 302,
        name: 'Durga Temple',
        address: '202 Divine Circle, Mysore',
        status: 'active'
      },
      {
        id: 303,
        name: 'Lakshmi Temple',
        address: '303 Prosperity Road, Mysore',
        status: 'active'
      }
    ]
  }
  
  // Return temples for the specified tenant, or empty array if none
  return mockTemplesMap[tenantId] || []
}
  
  // ==========================================
  // MOCK DATA HELPERS - EXISTING
  // ==========================================
  
  // Mock tenants data for development/demo
  // Store it as a class property so we can update it
  _mockTenants = [
    {
      ID: 1,
      FullName: "Naresh V",
      Email: "nareshvn4n@gmail.com",
      PasswordHash: "$2a$10$Q4IcyohhMOT49iyx0nRYqOaIVnZrh0b7nYn9CHZsiQEbV1rB6Rz4q",
      Phone: null,
      RoleID: 2,
      Role: {
        ID: 0,
        RoleName: "templeadmin",
        Description: "Temple Administrator",
        CanRegisterPublicly: true,
        CreatedAt: "2025-01-01T00:00:00Z",
        UpdatedAt: "2025-01-01T00:00:00Z"
      },
      EntityID: null,
      Status: "pending",
      EmailVerified: false,
      EmailVerifiedAt: null,
      CreatedAt: "2025-07-12T17:00:16.325394Z",
      UpdatedAt: "2025-07-12T17:00:16.325394Z",
      DeletedAt: null
    },
    {
      ID: 2,
      FullName: "GANESH",
      Email: "ganesh123@gmail.com",
      PasswordHash: "$2a$10$qvfEW9znvJEH19YDGiTyROSmTM9nzPjtpHXL/MDZ/3NoDqTsW8O5q",
      Phone: null,
      RoleID: 2,
      Role: {
        ID: 0,
        RoleName: "templeadmin",
        Description: "Temple Administrator",
        CanRegisterPublicly: true,
        CreatedAt: "2025-01-01T00:00:00Z",
        UpdatedAt: "2025-01-01T00:00:00Z"
      },
      EntityID: null,
      Status: "pending",
      EmailVerified: false,
      EmailVerifiedAt: null,
      CreatedAt: "2025-07-12T17:28:03.293811Z",
      UpdatedAt: "2025-07-12T17:28:03.293811Z",
      DeletedAt: null
    },
    {
      ID: 3,
      FullName: "Sharath Kumar",
      Email: "sharath@example.com",
      PasswordHash: "$2a$10$Q4IcyohhMOT49iyx0nRYqOaIVnZrh0b7nYn9CHZsiQEbV1rB6Rz4q",
      Phone: "9876543210",
      RoleID: 2,
      Role: {
        ID: 0,
        RoleName: "templeadmin",
        Description: "Temple Administrator",
        CanRegisterPublicly: true,
        CreatedAt: "2025-01-01T00:00:00Z",
        UpdatedAt: "2025-01-01T00:00:00Z"
      },
      EntityID: null,
      Status: "approved", // one approved for demo
      EmailVerified: true,
      EmailVerifiedAt: "2025-07-10T10:00:00Z",
      CreatedAt: "2025-07-10T10:00:00Z",
      UpdatedAt: "2025-07-10T10:00:00Z",
      DeletedAt: null
    },
    {
      ID: 4,
      FullName: "Rajesh K",
      Email: "rajesh@example.com",
      PasswordHash: "$2a$10$Q4IcyohhMOT49iyx0nRYqOaIVnZrh0b7nYn9CHZsiQEbV1rB6Rz4q",
      Phone: "8765432109",
      RoleID: 2,
      Role: {
        ID: 0,
        RoleName: "templeadmin",
        Description: "Temple Administrator",
        CanRegisterPublicly: true,
        CreatedAt: "2025-01-01T00:00:00Z",
        UpdatedAt: "2025-01-01T00:00:00Z"
      },
      EntityID: null,
      Status: "rejected", // one rejected for demo
      EmailVerified: false,
      EmailVerifiedAt: null,
      CreatedAt: "2025-07-09T15:30:00Z",
      UpdatedAt: "2025-07-09T17:45:00Z",
      DeletedAt: null
    }
  ]
  
  // Method to get pending tenants
  getMockPendingTenants() {
    return this._mockTenants.filter(t => t.Status === 'pending');
  }
  
  // Method to get all tenants
  getMockAllTenants() {
    return this._mockTenants;
  }
  
  // Method to update tenant status in mock data
  updateMockTenantStatus(tenantId, status) {
    // Make sure ID is treated as a number for comparison
    const id = parseInt(tenantId);
    
    this._mockTenants = this._mockTenants.map(tenant => {
      if (tenant.ID === id) {
        return {
          ...tenant,
          Status: status,
          UpdatedAt: new Date().toISOString()
        };
      }
      return tenant;
    });
    
    console.log(`Mock tenant ${tenantId} status updated to ${status}`);
  }
  
  // Generate mock entities data
  getMockPendingEntities() {
    return [
      {
        ID: 1,
        Name: "Sri Venkateswara Temple",
        Description: "Famous temple dedicated to Lord Venkateswara",
        Type: "Hindu",
        Address: "123 Temple Street",
        City: "Bengaluru",
        State: "Karnataka",
        Country: "India",
        Zip: "560001",
        Phone: "9876543210",
        Email: "info@svtemple.com",
        Website: "www.svtemple.com",
        TenantID: 1,
        Status: "pending",
        CreatedBy: 1,
        CreatedAt: "2025-07-11T10:00:00Z",
        UpdatedAt: "2025-07-11T10:00:00Z"
      },
      {
        ID: 2,
        Name: "Ganesh Mandir",
        Description: "Temple dedicated to Lord Ganesha",
        Type: "Hindu",
        Address: "456 Temple Lane",
        City: "Chennai",
        State: "Tamil Nadu",
        Country: "India",
        Zip: "600001",
        Phone: "8765432109",
        Email: "info@ganeshmandir.com",
        Website: "www.ganeshmandir.com",
        TenantID: 2,
        Status: "pending",
        CreatedBy: 2,
        CreatedAt: "2025-07-12T11:30:00Z",
        UpdatedAt: "2025-07-12T11:30:00Z"
      }
    ];
  }
  
  // Generate mock available tenants for user assignment
  getMockAvailableTenants() {
    return [
      {
        id: 1,
        userId: "T001",
        name: "Sri Krishna Trust",
        temple: {
          id: 101,
          name: "Sri Krishna Temple",
          address: "123 Temple Street, Bengaluru, KA"
        }
      },
      {
        id: 2,
        userId: "T002",
        name: "Ganesh Temple Trust",
        temple: {
          id: 102,
          name: "Ganesh Temple",
          address: "456 Mandir Road, Mysore, KA"
        }
      },
      {
        id: 3,
        userId: "T003",
        name: "Shiva Temple Trust",
        temple: {
          id: 103,
          name: "Shiva Temple",
          address: "789 Divine Lane, Hassan, KA"
        }
      },
      {
        id: 4,
        userId: "T004",
        name: "Lakshmi Temple Association",
        temple: {
          id: 104,
          name: "Lakshmi Temple",
          address: "101 Prosperity Avenue, Mangalore, KA"
        }
      },
      {
        id: 5,
        userId: "T005",
        name: "Saraswati Education Trust",
        temple: {
          id: 105,
          name: "Saraswati Temple",
          address: "202 Knowledge Street, Hubli, KA"
        }
      }
    ];
  }
  
  // Generate mock activities data
  getMockActivities(limit = 10) {
    const activities = [
      {
        id: 1,
        type: 'approval',
        description: 'Approved temple admin registration',
        timestamp: new Date(Date.now() - 3600000).toISOString() // 1 hour ago
      },
      {
        id: 2,
        type: 'approval',
        description: 'Approved new temple registration',
        timestamp: new Date(Date.now() - 7200000).toISOString() // 2 hours ago
      },
      {
        id: 3,
        type: 'new_application',
        description: 'Created new temple',
        timestamp: new Date(Date.now() - 86400000).toISOString() // 1 day ago
      },
      {
        id: 4,
        type: 'rejection',
        description: 'Rejected temple admin registration due to incomplete information',
        timestamp: new Date(Date.now() - 172800000).toISOString() // 2 days ago
      },
      {
        id: 5,
        type: 'new_application',
        description: 'Received donation of ₹5000',
        timestamp: new Date(Date.now() - 259200000).toISOString() // 3 days ago
      },
      {
        id: 6,
        type: 'new_application',
        description: 'Created new temple event: Ganesh Chaturthi Celebration',
        timestamp: new Date(Date.now() - 345600000).toISOString() // 4 days ago
      },
      {
        id: 7,
        type: 'new_application',
        description: 'Booked seva: Abhishekam',
        timestamp: new Date(Date.now() - 432000000).toISOString() // 5 days ago
      }
    ];
    
    return activities.slice(0, limit);
  }
}

// Create and export singleton instance
const superAdminService = new SuperAdminService()
export default superAdminService