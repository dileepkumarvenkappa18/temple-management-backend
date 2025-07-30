import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { donationService } from '@/services/donation.service'

export const useDonationStore = defineStore('donation', () => {
  // State
  const donations = ref([])
  const loading = ref(false)
  const error = ref(null)
  const selectedDonation = ref(null)
  const filters = ref({
    dateRange: 'all',
    minAmount: null,
    maxAmount: null,
    donationType: 'all',
    devoteeId: null,
    status: 'all',
    page: 1,
    limit: 10
  })

  // Dashboard data
  const dashboardData = ref({
  totalAmount: 0,
  averageAmount: 0,
  thisMonth: 0,
  totalDonors: 0,
  completed: 0,
  pending: 0
})
  const topDonors = ref([])

  // Getters
  const totalDonations = computed(() => {
    return filteredDonations.value.reduce((sum, donation) => sum + donation.amount, 0)
  })

  const recentDonations = computed(() => {
    return donations.value
      .filter(d => d.status === 'SUCCESS' || d.status === 'completed')
      .sort((a, b) => new Date(b.donatedAt || b.date) - new Date(a.donatedAt || a.date))
      .slice(0, 5)
  })

  const filteredDonations = computed(() => {
    let filtered = donations.value

    // Date range filter
    if (filters.value.dateRange !== 'all') {
      const now = new Date()
      const filterDate = new Date()
      
      switch (filters.value.dateRange) {
        case 'today':
          filterDate.setHours(0, 0, 0, 0)
          filtered = filtered.filter(d => new Date(d.donatedAt || d.date) >= filterDate)
          break
        case 'week':
          filterDate.setDate(now.getDate() - 7)
          filtered = filtered.filter(d => new Date(d.donatedAt || d.date) >= filterDate)
          break
        case 'month':
          filterDate.setMonth(now.getMonth() - 1)
          filtered = filtered.filter(d => new Date(d.donatedAt || d.date) >= filterDate)
          break
        case 'year':
          filterDate.setFullYear(now.getFullYear() - 1)
          filtered = filtered.filter(d => new Date(d.donatedAt || d.date) >= filterDate)
          break
      }
    }

    // Amount range filter
    if (filters.value.minAmount) {
      filtered = filtered.filter(d => d.amount >= filters.value.minAmount)
    }
    if (filters.value.maxAmount) {
      filtered = filtered.filter(d => d.amount <= filters.value.maxAmount)
    }

    // Donation type filter
    if (filters.value.donationType !== 'all') {
      filtered = filtered.filter(d => d.donationType === filters.value.donationType)
    }

    // Devotee filter
    if (filters.value.devoteeId) {
      filtered = filtered.filter(d => d.userID === filters.value.devoteeId)
    }

    // Status filter
    if (filters.value.status !== 'all') {
      filtered = filtered.filter(d => d.status.toLowerCase() === filters.value.status.toLowerCase())
    }

    return filtered
  })

  const donationStats = computed(() => {
    if (dashboardData.value) {
      return {
        total: dashboardData.value.total_count || 0,
        completed: dashboardData.value.total_count || 0,
        pending: 0,
        totalAmount: dashboardData.value.total_donations || 0,
        averageAmount: dashboardData.value.average_donation || 0,
        thisMonth: dashboardData.value.this_month || 0,
        totalDonors: dashboardData.value.total_donors || 0
      }
    }

    // Fallback calculation from donations array
    const completed = donations.value.filter(d => d.status === 'SUCCESS' || d.status === 'completed')
    const pending = donations.value.filter(d => d.status === 'PENDING' || d.status === 'pending')
    
    return {
      total: donations.value.length,
      completed: completed.length,
      pending: pending.length,
      totalAmount: completed.reduce((sum, d) => sum + d.amount, 0),
      averageAmount: completed.length > 0 ? Math.round(completed.reduce((sum, d) => sum + d.amount, 0) / completed.length) : 0,
      thisMonth: 0, // Will be filled from API
      totalDonors: new Set(donations.value.map(d => d.userID || d.devoteeId)).size
    }
  })

  // Actions
  const fetchDonations = async (entityId) => {
    loading.value = true
    error.value = null
    
    try {
      const apiFilters = {
        page: filters.value.page,
        limit: filters.value.limit,
        status: filters.value.status !== 'all' ? filters.value.status : '',
        type: filters.value.donationType !== 'all' ? filters.value.donationType : '',
        min: filters.value.minAmount,
        max: filters.value.maxAmount
      }
      
      const response = await donationService.getDonations(apiFilters)
      donations.value = response.data || []
      
      return response
    } catch (err) {
      error.value = err.message
      console.error('Error fetching donations:', err)
    } finally {
      loading.value = false
    }
  }

  const fetchMyDonations = async () => {
    loading.value = true
    error.value = null
    
    try {
      const response = await donationService.getMyDonations()
      donations.value = response || []
      return response
    } catch (err) {
      error.value = err.message
      console.error('Error fetching my donations:', err)
    } finally {
      loading.value = false
    }
  }

  const fetchDashboard = async () => {
    loading.value = true
    error.value = null
    
    try {
      const response = await donationService.getDashboard()
      dashboardData.value = response
      return response
    } catch (err) {
      error.value = err.message
      console.error('Error fetching donation dashboard:', err)
    } finally {
      loading.value = false
    }
  }

  const fetchTopDonors = async () => {
    loading.value = true
    error.value = null
    
    try {
      const response = await donationService.getTopDonors()
      topDonors.value = response?.top_donors || []
      return response
    } catch (err) {
      error.value = err.message
      console.error('Error fetching top donors:', err)
    } finally {
      loading.value = false
    }
  }

  const createDonation = async (donationData) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await donationService.createDonation(donationData)
      return response
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const verifyDonation = async (paymentData) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await donationService.verifyDonation(paymentData)
      return response
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const generateReceipt = async (donationId) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await donationService.generateReceipt(donationId)
      
      // Create and download the PDF blob
      const url = window.URL.createObjectURL(new Blob([response]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', `donation-receipt-${donationId}.pdf`);
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      
      return response
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const setFilters = (newFilters) => {
    filters.value = { ...filters.value, ...newFilters }
  }

  const clearFilters = () => {
    filters.value = {
      dateRange: 'all',
      minAmount: null,
      maxAmount: null,
      donationType: 'all',
      devoteeId: null,
      status: 'all',
      page: 1,
      limit: 10
    }
  }

  const setSelectedDonation = (donation) => {
    selectedDonation.value = donation
  }

  // Return public API
  return {
    // State
    donations,
    loading,
    error,
    selectedDonation,
    filters,
    dashboardData,
    topDonors,
    
    // Getters
    totalDonations,
    recentDonations,
    filteredDonations,
    donationStats,
    
    // Actions
    fetchDonations,
    fetchMyDonations,
    fetchDashboard,
    fetchTopDonors,
    createDonation,
    verifyDonation,
    generateReceipt,
    setFilters,
    clearFilters,
    setSelectedDonation
  }
})