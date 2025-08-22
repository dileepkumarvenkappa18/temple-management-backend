import api from '@/plugins/axios'
import { useAuthStore } from '@/stores/auth'

/**
 * Get all temples for a tenant
 * @param {string|number} tenantId - The tenant ID
 * @returns {Promise} - Promise with temple data
 */
export const getTemples = async (tenantId) => {
  try {
    console.log('📡 Making API call to fetch available temples')
    console.log('🔍 Search params:', { tenantId, headers: api.defaults.headers.common })
    
    // Get current URL path to determine context
    const currentPath = window.location.pathname
    console.log('📍 Current path:', currentPath)

    // Set the tenant ID header if provided
    if (tenantId) {
      api.defaults.headers.common['X-Tenant-ID'] = tenantId
    }
    
    // Get auth store to check user role
    const authStore = useAuthStore()
    const userRole = (authStore.userRole || '').toLowerCase()
    const isMonitoringUser = userRole.includes('monitoring') || userRole.includes('monitoringuser')
    const isStandardUser = userRole.includes('standard') || userRole.includes('standarduser')
    
    // Different endpoints based on the context
    let endpoint = '/v1/entities'
    
    // For standard/monitoring users, use tenant endpoint with the assigned tenant ID
    if ((isMonitoringUser || isStandardUser) && tenantId) {
      console.log('🔒 Standard/Monitoring user accessing temples for tenant ID:', tenantId)
      endpoint = '/v1/entities'
    } 
    // If accessing as tenant admin or with a specific tenant ID
    else if (tenantId) {
      console.log('🔒 Using admin endpoint: /v1/entities')
      endpoint = '/v1/entities'
    } else {
      // Default to user context (for devotees/volunteers)
      console.log('👤 Using user endpoint: /v1/entities/available')
      endpoint = '/v1/entities/available'
    }
    
    console.log(`🔍 Making GET request to: ${endpoint}`)
    const response = await api.get(endpoint)
    
    // Handle various response formats
    let temples = response.data
    if (response.data && response.data.data && Array.isArray(response.data.data)) {
      temples = response.data.data
    } else if (!Array.isArray(temples)) {
      console.warn('⚠️ Unexpected response format:', response.data)
      temples = []
    }
    
    console.log(`✅ Received ${temples.length} temples:`, temples)
    return temples
  } catch (error) {
    console.error('❌ Error fetching temples:', error)
    console.error('Error response:', error.response?.data)
    
    // Return empty array on error
    return []
  }
}

/**
 * Get a specific temple by ID
 * @param {string|number} id - Temple ID
 * @returns {Promise} - Promise with temple data
 */
export const getTempleById = async (id) => {
  try {
    console.log(`📡 Fetching temple with ID: ${id}`)
    const response = await api.get(`/v1/entities/${id}`)
    console.log('📥 Temple details response:', response.data)
    return response.data
  } catch (error) {
    console.error(`❌ Error fetching temple ID ${id}:`, error)
    console.error('Error response:', error.response?.data)
    throw error
  }
}

/**
 * Create a new temple
 * @param {Object} templeData - Temple data to create
 * @returns {Promise} - Promise with created temple data
 */
export const createTemple = async (templeData) => {
  try {
    const response = await api.post('/v1/entities', templeData)
    return response.data
  } catch (error) {
    console.error('Error creating temple:', error)
    throw error
  }
}

/**
 * Update an existing temple
 * @param {string|number} id - Temple ID to update
 * @param {Object} templeData - Updated temple data
 * @returns {Promise} - Promise with updated temple data
 */
export const updateTemple = async (id, templeData) => {
  try {
    const response = await api.put(`/v1/entities/${id}`, templeData)
    return response.data
  } catch (error) {
    console.error('Error updating temple:', error)
    throw error
  }
}

/**
 * Delete a temple
 * @param {string|number} id - Temple ID to delete
 * @returns {Promise} - Promise with deletion status
 */
export const deleteTemple = async (id) => {
  try {
    const response = await api.delete(`/v1/entities/${id}`)
    return response.data
  } catch (error) {
    console.error('Error deleting temple:', error)
    throw error
  }
}

export default {
  getTemples,
  getTempleById,
  createTemple,
  updateTemple,
  deleteTemple
}