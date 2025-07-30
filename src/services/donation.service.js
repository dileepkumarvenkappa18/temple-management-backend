// src/services/donation.service.js
import api from '@/plugins/axios'

export const donationService = {
  // Create a new donation (initiate payment process)
  async createDonation(donationData) {
    try {
      const response = await api.post('/v1/donations/', {
        amount: donationData.amount,
        donationType: donationData.donationType,
        note: donationData.purpose,
        referenceID: donationData.referenceID || ''
      })
      return response
    } catch (error) {
      console.error('Error creating donation:', error)
      throw error
    }
  },

  // Verify payment after successful Razorpay transaction
  async verifyDonation(paymentData) {
    try {
      const response = await api.post('/v1/donations/verify', {
        paymentID: paymentData.razorpay_payment_id,
        orderID: paymentData.razorpay_order_id,
        razorpaySig: paymentData.razorpay_signature
      })
      return response
    } catch (error) {
      console.error('Error verifying donation:', error)
      throw error
    }
  },

  // Get all donations for an entity admin (with pagination and filters)
  async getDonations(filters = {}) {
    try {
      // Convert filters into query parameters
      const queryParams = new URLSearchParams()
      
      if (filters.page) queryParams.append('page', filters.page)
      if (filters.limit) queryParams.append('limit', filters.limit)
      if (filters.status) queryParams.append('status', filters.status)
      if (filters.from) queryParams.append('from', filters.from)
      if (filters.to) queryParams.append('to', filters.to)
      if (filters.type) queryParams.append('type', filters.type)
      if (filters.method) queryParams.append('method', filters.method)
      if (filters.min) queryParams.append('min', filters.min)
      if (filters.max) queryParams.append('max', filters.max)
      if (filters.search) queryParams.append('search', filters.search)
      
      const queryString = queryParams.toString()
      const url = `/v1/donations/${queryString ? `?${queryString}` : ''}`
      
      const response = await api.get(url)
      return response
    } catch (error) {
      console.error('Error fetching donations:', error)
      throw error
    }
  },

  // Get donations for the current logged-in devotee
  async getMyDonations() {
    try {
      const response = await api.get('/v1/donations/my')
      return response
    } catch (error) {
      console.error('Error fetching my donations:', error)
      throw error
    }
  },

  // Get donation dashboard data for entity admins
 async getDashboard() {
  try {
    const entityId = localStorage.getItem('current_entity_id')
    if (!entityId) {
      throw new Error('Entity ID is missing in localStorage')
    }

    const response = await api.get(`/v1/donations/dashboard?entity_id=${entityId}`)
    return response
  } catch (error) {
    console.error('Error fetching donation dashboard:', error)
    throw error
  }
},

  // Get top donors for an entity
  async getTopDonors() {
    try {
      const response = await api.get('/v1/donations/top')
      return response
    } catch (error) {
      console.error('Error fetching top donors:', error)
      throw error
    }
  },

  // Generate receipt for a donation
  async generateReceipt(donationId) {
    try {
      const response = await api.get(`/v1/donations/${donationId}/receipt`, {
        responseType: 'blob'
      })
      return response
    } catch (error) {
      console.error('Error generating receipt:', error)
      throw error
    }
  },

  // Get donation types (for dropdown options)
  async getDonationTypes() {
    try {
      const response = await api.get('/v1/donations/types')
      return response
    } catch (error) {
      console.error('Error fetching donation types:', error)
      return [
        { value: 'general', label: 'General Donation' },
        { value: 'seva', label: 'Seva Donation' },
        { value: 'festival', label: 'Festival Donation' },
        { value: 'construction', label: 'Construction Fund' },
        { value: 'annadanam', label: 'Annadanam' },
        { value: 'education', label: 'Education Fund' },
        { value: 'maintenance', label: 'Maintenance' }
      ]
    }
  }
}

export default donationService