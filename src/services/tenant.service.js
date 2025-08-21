import api from '@/plugins/axios'
import { useAuthStore } from '@/stores/auth'

const tenantService = {
  /**
   * Get tenants available for selection based on user role
   * @returns {Promise<Array>} List of tenants for selection
   */
  async getTenantsForSelection() {
    // Get current user role to handle special cases
    const authStore = useAuthStore()
    const userRole = authStore.userRole?.toLowerCase() || ''
    const isMonitoringUser = userRole.includes('monitoring')
    
    try {
      console.log('📡 Fetching tenants for selection')
      
      // For monitoring users, directly use mock data until backend permission is fixed
    //   if (isMonitoringUser) {
    //     console.log('👀 Monitoring user detected - using direct data access')
    //     return this.getMockTenants()
    //   }
      
      // For other users, make the API call
      const response = await api.get('/v1/tenants/selection')
      console.log('📥 Tenants selection response:', response)
      
      // Extract data from response
      let tenantData = response.data || response
      if (!Array.isArray(tenantData)) {
        if (tenantData.data && Array.isArray(tenantData.data)) {
          tenantData = tenantData.data
        }
      }
      
      if (!Array.isArray(tenantData)) {
        console.error('🚨 Could not extract array from response:', response)
        return []
      }
      
      return tenantData.map(tenant => this.normalizeTenantData(tenant))
    } catch (error) {
      console.error('❌ Error fetching tenants for selection:', error)
      console.error('Error response:', error.response?.data)
      
      // Specific handling for 403 errors
      if (error.response?.status === 403) {
        console.warn('⚠️ Access forbidden - using mock data instead')
      } else {
        console.warn('⚠️ API error - using mock data as fallback')
      }
      
      // Return mock data for development until API is ready
      return this.getMockTenants()
    }
  },
  
  /**
   * Get mock tenants for development/testing
   * @returns {Array} List of mock tenants
   */
  getMockTenants() {
    console.log('📋 Using mock tenant data')
    return [
      {
        id: 1,
        name: 'Bangalore Temple Trust',
        email: 'admin@bangaloretemple.com',
        location: 'Bengaluru, Karnataka',
        status: 'active',
        templesCount: 5,
        imageUrl: null
      },
      {
        id: 2,
        name: 'Mumbai Temples Association',
        email: 'info@mumbaitemples.org',
        location: 'Mumbai, Maharashtra',
        status: 'active',
        templesCount: 8,
        imageUrl: null
      },
      {
        id: 3,
        name: 'Madurai Temple Management',
        email: 'admin@maduraitemples.com',
        location: 'Madurai, Tamil Nadu',
        status: 'active',
        templesCount: 3,
        imageUrl: null
      },
      {
        id: 4,
        name: 'Puri Temple Network',
        email: 'contact@puritemples.org',
        location: 'Puri, Odisha',
        status: 'pending',
        templesCount: 2,
        imageUrl: null
      }
    ].map(tenant => this.normalizeTenantData(tenant))
  },
  
  /**
   * Get a specific tenant by ID
   * @param {number|string} id - Tenant ID
   * @returns {Promise<Object>} Tenant details
   */
  async getTenantById(id) {
    try {
      console.log(`📡 Fetching tenant with ID: ${id}`)
      
      const response = await api.get(`/v1/entities/${id}`)
      console.log('📥 Tenant by ID response:', response)
      
      return this.normalizeTenantData(response.data || response)
    } catch (error) {
      console.error(`❌ Error fetching tenant ID ${id}:`, error)
      console.error('Error response:', error.response?.data)
      
      // Return mock data for the requested ID
      const mockData = this.getMockTenants().find(t => Number(t.id) === Number(id))
      if (mockData) return mockData
      
      throw error
    }
  },
  
  /**
   * Normalize tenant data from backend
   * @param {Object} tenant - Raw tenant data from backend
   * @returns {Object} - Normalized tenant data
   */
  normalizeTenantData(tenant) {
    if (!tenant) return null
    
    return {
      id: Number(tenant.id) || 0,
      name: tenant.name || tenant.tenantName || tenant.temple_name || 'Unknown Tenant',
      email: tenant.email || '',
      location: tenant.location || tenant.templeAddress || tenant.temple_address || '',
      status: (tenant.status || 'active').toLowerCase(),
      templesCount: tenant.templesCount || 0,
      imageUrl: tenant.imageUrl || null
    }
  }
}

export default tenantService