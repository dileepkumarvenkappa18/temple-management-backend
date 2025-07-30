// src/services/seva.service.js
import api from '@/plugins/axios'

class SevaService {
  /**
   * Get all sevas for a specific entity (temple)
   * @param {string} entityId - Temple ID
   * @param {Object} params - Query parameters
   * @returns {Promise<Object>} Seva list with pagination
   */
  async getSevas(entityId, params = {}) {
    try {
      // Remove redundant /api prefix, as it's already in the baseURL
      const response = await api.get('/v1/sevas', { 
        params: { ...params, entity_id: entityId } 
      })
      
      return {
        success: true,
        data: response.data || [],
        pagination: response.data.pagination || {},
        total: response.data.total || 0
      }
    } catch (error) {
      console.error('Error fetching sevas:', error)
      return {
        success: false,
        error: error.response?.data?.error || 'Failed to fetch sevas',
        data: []
      }
    }
  }

  /**
   * Get seva by ID
   * @param {string} entityId - Temple ID
   * @param {string} sevaId - Seva ID
   * @returns {Promise<Object>} Seva details
   */
  async getSevaById(entityId, sevaId) {
    try {
      const response = await api.get(`/v1/sevas/${sevaId}`, {
        params: { entity_id: entityId }
      })
      return {
        success: true,
        data: response.data || null
      }
    } catch (error) {
      console.error('Error fetching seva:', error)
      return {
        success: false,
        error: error.response?.data?.error || 'Failed to fetch seva details',
        data: null
      }
    }
  }

  /**
   * Create new seva
   * @param {Object} sevaData - Seva information with entity_id included
   * @returns {Promise<Object>} Created seva
   */
  async createSeva(sevaData) {
    try {
      console.log('Creating seva with data:', sevaData)
      const response = await api.post('/v1/sevas', sevaData)
      
      return {
        success: true,
        data: response.data,
        message: 'Seva created successfully'
      }
    } catch (error) {
      console.error('Error creating seva:', error)
      return {
        success: false,
        error: error.response?.data?.error || 'Failed to create seva',
        errors: error.response?.data?.errors || {}
      }
    }
  }

  /**
   * Update seva
   * @param {string} sevaId - Seva ID
   * @param {Object} sevaData - Updated seva information
   * @returns {Promise<Object>} Updated seva
   */
  async updateSeva(sevaId, sevaData) {
    try {
      console.log('Updating seva with ID:', sevaId, 'Data:', sevaData)
      const response = await api.put(`/v1/sevas/${sevaId}`, sevaData)
      
      return {
        success: true,
        data: response.data,
        message: 'Seva updated successfully'
      }
    } catch (error) {
      console.error('Error updating seva:', error)
      return {
        success: false,
        error: error.response?.data?.error || 'Failed to update seva',
        errors: error.response?.data?.errors || {}
      }
    }
  }

  /**
   * Delete seva
   * @param {string} sevaId - Seva ID
   * @returns {Promise<Object>} Success status
   */
  async deleteSeva(sevaId) {
    try {
      await api.delete(`/v1/sevas/${sevaId}`)
      
      return {
        success: true,
        message: 'Seva deleted successfully'
      }
    } catch (error) {
      console.error('Error deleting seva:', error)
      return {
        success: false,
        error: error.response?.data?.error || 'Failed to delete seva'
      }
    }
  }

  /**
   * Get seva bookings for entity (temple)
   * @returns {Promise<Object>} Entity seva bookings
   */
  async getEntityBookings() {
    try {
      const response = await api.get('/v1/sevas/entity-bookings')
      
      return {
        success: true,
        data: response.data || [],
        pagination: response.data.pagination || {},
        total: response.data.total || 0
      }
    } catch (error) {
      console.error('Error fetching entity bookings:', error)
      return {
        success: false,
        error: error.response?.data?.error || 'Failed to fetch bookings',
        data: []
      }
    }
  }

  /**
   * Update booking status (approve/reject)
   * @param {string} bookingId - Booking ID
   * @param {string} status - New status (approved/rejected)
   * @returns {Promise<Object>} Updated booking
   */
  async updateBookingStatus(bookingId, status) {
    try {
      const response = await api.patch(`/v1/sevas/bookings/${bookingId}/status`, { 
        status 
      })
      
      return {
        success: true,
        data: response.data,
        message: `Booking ${status} successfully`
      }
    } catch (error) {
      console.error('Error updating booking status:', error)
      return {
        success: false,
        error: error.response?.data?.error || 'Failed to update booking status'
      }
    }
  }
}

// Export singleton instance
export const sevaService = new SevaService()
export default sevaService