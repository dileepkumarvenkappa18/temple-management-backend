// src/stores/reports.js
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import reportsService from '@/services/reports.service'

export const useReportsStore = defineStore('reports', () => {
  // State
  const loading = ref(false)
  const downloadLoading = ref(false)
  const error = ref(null)
  const currentReport = ref(null)
  const reportPreview = ref(null)
  const lastReportParams = ref(null)
  
  // MISSING REACTIVE VARIABLES - ADDED
  const devoteeList = ref([])
  const devoteeProfile = ref(null)
  const devoteeListStatus = ref('all')

  // Getters
  const hasReportData = computed(() => {
    return reportPreview.value && reportPreview.value.data && reportPreview.value.data.length > 0
  })

  const reportSummary = computed(() => {
    if (!reportPreview.value) return null
    
    return {
      totalRecords: reportPreview.value.totalRecords || 0,
      type: lastReportParams.value?.type || 'Unknown',
      dateRange: lastReportParams.value?.dateRange || 'Unknown',
      entityId: lastReportParams.value?.entityId || 'Unknown'
    }
  })

  // MISSING COMPUTED - ADDED
  const filteredDevoteeList = computed(() => {
    if (!devoteeList.value || !Array.isArray(devoteeList.value)) {
      return []
    }
    
    if (devoteeListStatus.value === 'all') {
      return devoteeList.value
    }
    
    return devoteeList.value.filter(devotee => {
      if (devoteeListStatus.value === 'active') {
        return devotee.status === 'active' || devotee.is_active === true
      } else if (devoteeListStatus.value === 'inactive') {
        return devotee.status === 'inactive' || devotee.is_active === false
      } else if (devoteeListStatus.value === 'new') {
        // Assuming 'new' means recently joined (within last 30 days)
        const thirtyDaysAgo = new Date()
        thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30)
        const joinDate = new Date(devotee.created_at || devotee.join_date)
        return joinDate >= thirtyDaysAgo
      }
      return true
    })
  })

  // Actions
  const clearError = () => {
    error.value = null
  }

  const clearReportData = () => {
    currentReport.value = null
    reportPreview.value = null
    lastReportParams.value = null
    error.value = null
    // ADDED: Clear devotee-specific data
    devoteeList.value = []
    devoteeProfile.value = null
    devoteeListStatus.value = 'all'
  }

  /**
   * Fetch activities report data (JSON preview)
   */
  const fetchActivitiesReport = async (params) => {
    try {
      loading.value = true
      error.value = null

      // Validate parameters
      const validation = reportsService.validateReportParams(params)
      if (!validation.isValid) {
        throw new Error(validation.errors.join(', '))
      }

      // Store params for reference
      lastReportParams.value = { ...params }

      // Fetch report data
      const response = await reportsService.getActivitiesReport(params)
      currentReport.value = response

      // Get formatted preview
      const preview = await reportsService.getReportPreview(params)
      reportPreview.value = preview

      return response
    } catch (err) {
      error.value = err.message || 'Failed to fetch activities report'
      console.error('Error in fetchActivitiesReport:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Download activities report in specified format
   */
  const downloadActivitiesReport = async (params) => {
    try {
      downloadLoading.value = true
      error.value = null

      // Validate parameters
      const validation = reportsService.validateReportParams(params)
      if (!validation.isValid) {
        throw new Error(validation.errors.join(', '))
      }

      // Download report
      const result = await reportsService.downloadActivitiesReport(params)
      
      // Store successful download params
      lastReportParams.value = { ...params }

      return result
    } catch (err) {
      error.value = err.message || 'Failed to download activities report'
      console.error('Error in downloadActivitiesReport:', err)
      throw err
    } finally {
      downloadLoading.value = false
    }
  }

  /**
   * Get report preview with loading state
   */
  const getReportPreview = async (params) => {
    try {
      loading.value = true
      error.value = null

      const preview = await reportsService.getReportPreview(params)
      reportPreview.value = preview
      lastReportParams.value = { ...params }

      return preview
    } catch (err) {
      error.value = err.message || 'Failed to get report preview'
      console.error('Error in getReportPreview:', err)
      throw err
    } finally {
      loading.value = false
    }
  }
/**
 * Fetch audit logs report data (JSON preview)
 */
const fetchAuditLogsReport = async (params) => {
  try {
    loading.value = true
    error.value = null

    // Add report type for summary
    lastReportParams.value = { ...params, type: 'audit-logs' }

    const response = await reportsService.getAuditLogsReport(params)
    currentReport.value = response

    const preview = await reportsService.getAuditLogsPreview(params)
    reportPreview.value = preview

    return response
  } catch (err) {
    error.value = err.message || 'Failed to fetch audit logs report'
    console.error('Error in fetchAuditLogsReport:', err)
    throw err
  } finally {
    loading.value = false
  }
}

/**
 * Download audit logs report in specified format
 */
const downloadAuditLogsReport = async (params) => {
  try {
    downloadLoading.value = true
    error.value = null

    const result = await reportsService.downloadAuditLogsReport(params)
    lastReportParams.value = { ...params, type: 'audit-logs' }

    return result
  } catch (err) {
    error.value = err.message || 'Failed to download audit logs report'
    console.error('Error in downloadAuditLogsReport:', err)
    throw err
  } finally {
    downloadLoading.value = false
  }
}

  // TEMPLE REGISTERED REPORT METHODS
  /**
   * Fetch temple registered report data (JSON preview)
   */
  const fetchTempleRegisteredReport = async (params) => {
    try {
      loading.value = true
      error.value = null

      // Store params for reference
      lastReportParams.value = { ...params, type: 'temple-registered' }

      // Fetch report data
      const response = await reportsService.getTempleRegisteredReport(params)
      currentReport.value = response

      // Get formatted preview
      const preview = await reportsService.getTempleRegisteredPreview(params)
      reportPreview.value = preview

      return response
    } catch (err) {
      error.value = err.message || 'Failed to fetch temple registered report'
      console.error('Error in fetchTempleRegisteredReport:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Download temple registered report in specified format
   */
  const downloadTempleRegisteredReport = async (params) => {
    try {
      downloadLoading.value = true
      error.value = null

      // Download report
      const result = await reportsService.downloadTempleRegisteredReport(params)
      
      // Store successful download params
      lastReportParams.value = { ...params, type: 'temple-registered' }

      return result
    } catch (err) {
      error.value = err.message || 'Failed to download temple registered report'
      console.error('Error in downloadTempleRegisteredReport:', err)
      throw err
    } finally {
      downloadLoading.value = false
    }
  }

  // DEVOTEE BIRTHDAYS REPORT METHODS
  /**
   * Fetch devotee birthdays report data (JSON preview)
   */
  const fetchDevoteeBirthdaysReport = async (params) => {
    try {
      loading.value = true
      error.value = null

      // Store params for reference
      lastReportParams.value = { ...params, type: 'devotee-birthdays' }

      // Fetch report data
      const response = await reportsService.getDevoteeBirthdaysReport(params)
      currentReport.value = response

      // Get formatted preview
      const preview = await reportsService.getDevoteeBirthdaysPreview(params)
      reportPreview.value = preview

      return response
    } catch (err) {
      error.value = err.message || 'Failed to fetch devotee birthdays report'
      console.error('Error in fetchDevoteeBirthdaysReport:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Download devotee birthdays report in specified format
   */
  const downloadDevoteeBirthdaysReport = async (params) => {
    try {
      downloadLoading.value = true
      error.value = null

      // Download report
      const result = await reportsService.downloadDevoteeBirthdaysReport(params)
      
      // Store successful download params
      lastReportParams.value = { ...params, type: 'devotee-birthdays' }

      return result
    } catch (err) {
      error.value = err.message || 'Failed to download devotee birthdays report'
      console.error('Error in downloadDevoteeBirthdaysReport:', err)
      throw err
    } finally {
      downloadLoading.value = false
    }
  }

  /**
   * Get devotee birthdays preview with loading state
   */
  const getDevoteeBirthdaysPreview = async (params) => {
    try {
      loading.value = true
      error.value = null

      const preview = await reportsService.getDevoteeBirthdaysPreview(params)
      reportPreview.value = preview
      lastReportParams.value = { ...params, type: 'devotee-birthdays' }

      return preview
    } catch (err) {
      error.value = err.message || 'Failed to get devotee birthdays preview'
      console.error('Error in getDevoteeBirthdaysPreview:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // DEVOTEE LIST REPORT METHODS
  const fetchDevoteeListReport = async (params) => {
    try {
      loading.value = true
      error.value = null
      const { entityId, status = 'all' } = params
      lastReportParams.value = { ...params, type: 'devotee-list' }
      
      const response = await reportsService.getDevoteeList({ entityId, status })
      currentReport.value = response
      
      // Handle nested response data
      let responseData = response.data
      if (responseData && responseData.data) {
        responseData = responseData.data
      }
      
      // FIXED: Set devoteeList properly
      devoteeList.value = responseData.devotees || responseData || []
      devoteeListStatus.value = status
      
      const preview = await reportsService.getDevoteeListPreview({ entityId, status })
      reportPreview.value = preview
      
      return response
    } catch (err) {
      error.value = err.message || 'Failed to fetch devotee list report'
      console.error('Error in fetchDevoteeListReport:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  const downloadDevoteeListReport = async (params) => {
    try {
      downloadLoading.value = true
      error.value = null
      const { format } = params
      if (!format) {
        throw new Error('Format is required for download')
      }
      const result = await reportsService.downloadDevoteeListReport(params)
      lastReportParams.value = { ...params, type: 'devotee-list' }
      return result
    } catch (err) {
      error.value = err.message || 'Failed to download devotee list report'
      console.error('Error in downloadDevoteeListReport:', err)
      throw err
    } finally {
      downloadLoading.value = false
    }
  }

  const setDevoteeListStatus = async (entityId, status) => {
    try {
      loading.value = true
      error.value = null
      devoteeListStatus.value = status
      await fetchDevoteeListReport({ entityId, status })
    } catch (err) {
      error.value = err.message || 'Failed to filter devotee list'
      console.error('Error in setDevoteeListStatus:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // DEVOTEE PROFILE METHODS
  const fetchDevoteeProfile = async (params) => {
    try {
      loading.value = true
      error.value = null
      const { entityId } = params
      lastReportParams.value = { ...params, type: 'devotee-profile' }
      
      const response = await reportsService.getDevoteeProfile({ entityId })
      currentReport.value = response
      
      const preview = await reportsService.getDevoteeProfilePreview({ entityId })
      devoteeProfile.value = preview
      
      return response
    } catch (err) {
      error.value = err.message || 'Failed to fetch devotee profile'
      console.error('Error in fetchDevoteeProfile:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  const downloadDevoteeProfileReport = async (params) => {
    try {
      downloadLoading.value = true
      error.value = null
      const { format } = params
      if (!format) {
        throw new Error('Format is required for download')
      }
      const result = await reportsService.downloadDevoteeProfileReport(params)
      lastReportParams.value = { ...params, type: 'devotee-profile' }
      return result
    } catch (err) {
      error.value = err.message || 'Failed to download devotee profile report'
      console.error('Error in downloadDevoteeProfileReport:', err)
      throw err
    } finally {
      downloadLoading.value = false
    }
  }

  const getDevoteeListPreview = async (params) => {
    try {
      loading.value = true
      error.value = null
      const preview = await reportsService.getDevoteeListPreview(params)
      reportPreview.value = preview
      lastReportParams.value = { ...params, type: 'devotee-list' }
      return preview
    } catch (err) {
      error.value = err.message || 'Failed to get devotee list preview'
      console.error('Error in getDevoteeListPreview:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // ADDED: Missing getDevoteeProfilePreview method
  const getDevoteeProfilePreview = async (params) => {
    try {
      loading.value = true
      error.value = null
      const preview = await reportsService.getDevoteeProfilePreview(params)
      reportPreview.value = preview
      lastReportParams.value = { ...params, type: 'devotee-profile' }
      return preview
    } catch (err) {
      error.value = err.message || 'Failed to get devotee profile preview'
      console.error('Error in getDevoteeProfilePreview:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Generate default date range based on preset
   */
  const getDefaultDateRange = (preset) => {
    const today = new Date()
    const startDate = new Date()
    
    switch (preset) {
      case 'daily':
        // Today only
        startDate.setDate(today.getDate())
        break
      case 'weekly':
        // Last 7 days
        startDate.setDate(today.getDate() - 7)
        break
      case 'monthly':
        // Last 30 days
        startDate.setDate(today.getDate() - 30)
        break
      case 'yearly':
        // Last 365 days
        startDate.setDate(today.getDate() - 365)
        break
      default:
        // Default to weekly
        startDate.setDate(today.getDate() - 7)
    }

    return {
      startDate: startDate.toISOString().split('T')[0],
      endDate: today.toISOString().split('T')[0]
    }
  }

  /**
   * Build report parameters from component state
   */
  const buildReportParams = (componentState) => {
    const {
      selectedTemple,
      activityType,
      activeFilter,
      selectedFormat,
      startDate,
      endDate
    } = componentState

    // Use default date range if custom dates aren't provided
    let dates = { startDate, endDate }
    if (activeFilter !== 'custom' || !startDate || !endDate) {
      dates = getDefaultDateRange(activeFilter)
    }

    return {
      entityId: selectedTemple === 'all' ? 'all' : selectedTemple.toString(),
      type: activityType,
      dateRange: activeFilter,
      format: selectedFormat,
      startDate: dates.startDate,
      endDate: dates.endDate
    }
  }

  /**
   * Build devotee birthdays report parameters from component state
   */
  const buildDevoteeBirthdaysParams = (componentState) => {
    const {
      selectedTemple,
      activeFilter,
      selectedFormat,
      startDate,
      endDate
    } = componentState

    // Use default date range if custom dates aren't provided
    let dates = { startDate, endDate }
    if (activeFilter !== 'custom' || !startDate || !endDate) {
      dates = getDefaultDateRange(activeFilter)
    }

    return {
      entityId: selectedTemple === 'all' ? 'all' : selectedTemple.toString(),
      dateRange: activeFilter,
      format: selectedFormat,
      startDate: dates.startDate,
      endDate: dates.endDate
    }
  }

  const buildDevoteeListParams = (componentState) => {
    const {
      selectedTemple,
      selectedFormat,
      status = 'all'
    } = componentState
    return {
      entityId: selectedTemple === 'all' ? 'all' : selectedTemple.toString(),
      status,
      format: selectedFormat
    }
  }

  const buildDevoteeProfileParams = (componentState) => {
    const {
      selectedTemple,
      selectedFormat
    } = componentState
    return {
      entityId: selectedTemple === 'all' ? 'all' : selectedTemple.toString(),
      format: selectedFormat
    }
  }

  /**
   * Format report data for display
   */
  const formatReportData = (data, type) => {
    if (!data || !Array.isArray(data)) return []

    return data.map(item => {
      // Common formatting for all types
      const formatted = { ...item }

      // Format dates
      if (item.created_at) {
        formatted.created_at = new Date(item.created_at).toLocaleDateString()
      }
      if (item.updated_at) {
        formatted.updated_at = new Date(item.updated_at).toLocaleDateString()
      }

      // Type-specific formatting
      switch (type) {
        case 'events':
          if (item.event_date) {
            formatted.event_date = new Date(item.event_date).toLocaleDateString()
          }
          break
        case 'sevas':
          if (item.date) {
            formatted.date = new Date(item.date).toLocaleDateString()
          }
          if (item.price) {
            formatted.price = `₹${item.price.toFixed(2)}`
          }
          break
        case 'bookings':
          if (item.booking_time) {
            formatted.booking_time = new Date(item.booking_time).toLocaleString()
          }
          break
        case 'devotee-birthdays':
          if (item.date_of_birth) {
            formatted.date_of_birth = new Date(item.date_of_birth).toLocaleDateString()
          }
          if (item.member_since) {
            formatted.member_since = new Date(item.member_since).toLocaleDateString()
          }
          break
      }

      return formatted
    })
  }

  return {
    // State
    loading,
    downloadLoading,
    error,
    currentReport,
    reportPreview,
    lastReportParams,
    devoteeList,
    devoteeProfile,
    devoteeListStatus,
    
    // Getters
    hasReportData,
    reportSummary,
    filteredDevoteeList,

    //audit lgs
    fetchAuditLogsReport,
  downloadAuditLogsReport,
    
    // Actions
    clearError,
    clearReportData,
    fetchActivitiesReport,
    downloadActivitiesReport,
    getReportPreview,
    // Temple registered methods
    fetchTempleRegisteredReport,
    downloadTempleRegisteredReport,
    // Devotee birthdays methods
    fetchDevoteeBirthdaysReport,
    downloadDevoteeBirthdaysReport,
    getDevoteeBirthdaysPreview,
    // Devotee list methods
    fetchDevoteeListReport,
    downloadDevoteeListReport,
    getDevoteeListPreview,
    setDevoteeListStatus,
    // Devotee profile methods
    fetchDevoteeProfile,
    downloadDevoteeProfileReport,
    getDevoteeProfilePreview,

    // Utility methods
    getDefaultDateRange,
    buildReportParams,
    buildDevoteeBirthdaysParams,
    formatReportData,
    buildDevoteeListParams,
    buildDevoteeProfileParams
  }
})