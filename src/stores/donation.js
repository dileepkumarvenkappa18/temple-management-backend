import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { donationService } from '@/services/donation.service'

export const useDonationStore = defineStore('donation', () => {
  // State
  const donations = ref([])
  const loading = ref(false)
  const error = ref(null)
  const selectedDonation = ref(null)

  // Pagination and filtering state
  const pagination = ref({
    currentPage: 1,
    totalPages: 1,
    totalItems: 0,
    itemsPerPage: 10,
    hasNext: false,
    hasPrevious: false
  })

  const filters = ref({
    dateRange: 'all',
    minAmount: null,
    maxAmount: null,
    donationType: 'all',
    paymentMethod: 'all',
    devoteeId: null,
    status: 'all',
    search: '',
    // Date range specific
    startDate: null,
    endDate: null
  })

  const dashboardData = ref({
    totalAmount: 0,
    averageAmount: 0,
    thisMonth: 0,
    totalDonors: 0,
    completed: 0,
    pending: 0,
    failed: 0,
    totalCount: 0,
  })

  const topDonors = ref([])
  const analytics = ref({
    trends: [],
    byType: [],
    byMethod: [],
  })

  // Getters
  const donationStats = computed(() => {
    return {
      total: dashboardData.value.totalCount || 0,
      completed: dashboardData.value.completed || 0,
      pending: dashboardData.value.pending || 0,
      failed: dashboardData.value.failed || 0,
      totalAmount: dashboardData.value.totalAmount || 0,
      averageAmount: dashboardData.value.averageAmount || 0,
      thisMonth: dashboardData.value.thisMonth || 0,
      totalDonors: dashboardData.value.totalDonors || 0,
    }
  })

  const recentDonations = computed(() =>
    donations.value
      .filter(d => ['success', 'completed'].includes((d.status || '').toLowerCase()))
      .sort((a, b) => new Date(b.donatedAt || b.date) - new Date(a.donatedAt || a.date))
      .slice(0, 5)
  )

  // Actions
  async function fetchMyDonations() {
    loading.value = true
    error.value = null
    try {
      console.log('Store: Fetching my donations...')
      const response = await donationService.getMyDonations()
      console.log('Store: Received donations response:', response)

      // Handle different response structures
      if (response && typeof response === 'object') {
        if (Array.isArray(response)) {
          donations.value = response
          pagination.value.totalItems = response.length
        } else if (response.data && Array.isArray(response.data)) {
          donations.value = response.data
          // Update pagination if provided
          if (response.pagination) {
            pagination.value = {
              ...pagination.value,
              ...response.pagination
            }
          } else {
            pagination.value.totalItems = response.total || response.data.length
          }
        } else if (response.success && Array.isArray(response.data)) {
          donations.value = response.data
          pagination.value.totalItems = response.total || response.data.length
        } else {
          console.warn('Unexpected response format:', response)
          donations.value = []
        }
      } else {
        donations.value = []
      }

      console.log('Store: Final donations:', donations.value)
      console.log('Store: Final pagination:', pagination.value)
      return donations.value
    } catch (err) {
      console.error('Store: Error fetching my donations:', err)
      error.value = err.message || 'Error fetching donations'
      donations.value = []
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchDonations() {
    loading.value = true
    error.value = null
    try {
      // Build API filters from current state
      const apiFilters = {
        page: pagination.value.currentPage,
        limit: pagination.value.itemsPerPage,
        ...(filters.value.status !== 'all' && { status: filters.value.status }),
        ...(filters.value.donationType !== 'all' && { type: filters.value.donationType }),
        ...(filters.value.paymentMethod !== 'all' && { method: filters.value.paymentMethod }),
        ...(filters.value.search && { search: filters.value.search.trim() }),
        ...(filters.value.minAmount !== null && filters.value.minAmount !== '' && { min: filters.value.minAmount }),
        ...(filters.value.maxAmount !== null && filters.value.maxAmount !== '' && { max: filters.value.maxAmount }),
        ...(filters.value.dateRange !== 'all' && { dateRange: filters.value.dateRange }),
        ...(filters.value.startDate && { from: filters.value.startDate }),
        ...(filters.value.endDate && { to: filters.value.endDate }),
        ...(filters.value.devoteeId && { devoteeId: filters.value.devoteeId }),
      }

      console.log('Fetching donations with filters:', apiFilters)
      const response = await donationService.getDonations(apiFilters)

      // Handle response
      if (response && typeof response === 'object') {
        // Set donations
        donations.value = response.data || response.donations || []

        // Update pagination from response
        pagination.value = {
          currentPage: response.currentPage || response.page || pagination.value.currentPage,
          totalPages: response.totalPages || Math.ceil((response.total || 0) / pagination.value.itemsPerPage),
          totalItems: response.total || response.totalItems || 0,
          itemsPerPage: response.limit || response.perPage || pagination.value.itemsPerPage,
          hasNext: response.hasNext || (pagination.value.currentPage < pagination.value.totalPages),
          hasPrevious: response.hasPrevious || (pagination.value.currentPage > 1)
        }

        // Update dashboard data if provided in response
        if (response.stats) {
          dashboardData.value = {
            ...dashboardData.value,
            ...response.stats
          }
        }
      }

      console.log('Updated donations:', donations.value)
      console.log('Updated pagination:', pagination.value)
      return response
    } catch (err) {
      console.error('Error fetching donations:', err)
      error.value = err.message || 'Error fetching donations'
      donations.value = []
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchDashboard() {
    loading.value = true
    error.value = null
    try {
      const response = await donationService.getDashboard()

      // Handle different response structures
      const data = response.data || response
      dashboardData.value = {
        totalAmount: data.totalAmount || 0,
        averageAmount: data.averageAmount || 0,
        thisMonth: data.thisMonth || 0,
        totalDonors: data.totalDonors || data.uniqueDonors || 0,
        completed: data.completed || data.success || 0,
        pending: data.pending || 0,
        failed: data.failed || 0,
        totalCount: data.totalCount || data.total || 0,
      }

      return response
    } catch (err) {
      console.error('Error fetching dashboard:', err)
      error.value = err.message || 'Error fetching dashboard'
      // Set default values on error
      dashboardData.value = {
        totalAmount: 0,
        averageAmount: 0,
        thisMonth: 0,
        totalDonors: 0,
        completed: 0,
        pending: 0,
        failed: 0,
        totalCount: 0,
      }
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchTopDonors(limit = 5) {
    try {
      const response = await donationService.getTopDonors(limit)
      topDonors.value = response.data || response || []
      return response
    } catch (err) {
      console.error('Error fetching top donors:', err)
      error.value = err.message || 'Error fetching top donors'
      topDonors.value = []
      throw err
    }
  }

  async function fetchAnalytics(days = 30) {
    try {
      const response = await donationService.getAnalytics(days)
      analytics.value = response.data || response || { trends: [], byType: [], byMethod: [] }
      return response
    } catch (err) {
      console.error('Error fetching analytics:', err)
      error.value = err.message || 'Error fetching analytics'
      analytics.value = { trends: [], byType: [], byMethod: [] }
      throw err
    }
  }

  async function createDonation(data) {
    loading.value = true
    error.value = null
    try {
      const response = await donationService.createDonation(data)
      return response
    } catch (err) {
      error.value = err.message || 'Error creating donation'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function verifyDonation(paymentData) {
    loading.value = true
    error.value = null
    try {
      const response = await donationService.verifyDonation(paymentData)
      // Refresh donations after verification
      await fetchMyDonations()
      return response
    } catch (err) {
      error.value = err.message || 'Error verifying donation'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function exportDonations(exportFilters = {}, format = 'csv') {
    loading.value = true
    error.value = null
    try {
      const response = await donationService.exportDonations(exportFilters, format)
      return response
    } catch (err) {
      error.value = err.message || 'Error exporting donations'
      throw err
    } finally {
      loading.value = false
    }
  }

  // Filter and pagination methods
  function setFilters(newFilters) {
    filters.value = { ...filters.value, ...newFilters }
    // Reset to first page when filters change
    if (!newFilters.page) {
      pagination.value.currentPage = 1
    }
  }

  function setPage(page) {
    pagination.value.currentPage = page
  }

  function setItemsPerPage(items) {
    pagination.value.itemsPerPage = items
    pagination.value.currentPage = 1 // Reset to first page
  }

  function resetFilters() {
    filters.value = {
      dateRange: 'all',
      minAmount: null,
      maxAmount: null,
      donationType: 'all',
      paymentMethod: 'all',
      devoteeId: null,
      status: 'all',
      search: '',
      startDate: null,
      endDate: null
    }
    pagination.value.currentPage = 1
  }

  function setSelectedDonation(donation) {
    selectedDonation.value = donation
  }

  function resetStore() {
    donations.value = []
    loading.value = false
    error.value = null
    selectedDonation.value = null
    resetFilters()
    pagination.value = {
      currentPage: 1,
      totalPages: 1,
      totalItems: 0,
      itemsPerPage: 10,
      hasNext: false,
      hasPrevious: false
    }
    dashboardData.value = {
      totalAmount: 0,
      averageAmount: 0,
      thisMonth: 0,
      totalDonors: 0,
      completed: 0,
      pending: 0,
      failed: 0,
      totalCount: 0,
    }
    topDonors.value = []
    analytics.value = { trends: [], byType: [], byMethod: [] }
  }

  return {
    // State
    donations,
    loading,
    error,
    selectedDonation,
    filters,
    pagination,
    dashboardData,
    topDonors,
    analytics,

    // Getters
    donationStats,
    recentDonations,

    // Actions
    fetchMyDonations,
    fetchDonations,
    fetchDashboard,
    fetchTopDonors,
    fetchAnalytics,
    createDonation,
    verifyDonation,
    exportDonations,

    // Utilities
    setFilters,
    setPage,
    setItemsPerPage,
    resetFilters,
    setSelectedDonation,
    resetStore,
  }
})
