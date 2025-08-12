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

  // Actions
  const clearError = () => {
    error.value = null
  }

  const clearReportData = () => {
    currentReport.value = null
    reportPreview.value = null
    lastReportParams.value = null
    error.value = null
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

  // DEVOTEE BIRTHDAYS REPORT METHODS - NEW
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
    
    // Getters
    hasReportData,
    reportSummary,
    
    // Actions
    clearError,
    clearReportData,
    fetchActivitiesReport,
    downloadActivitiesReport,
    getReportPreview,
    // Temple registered methods
    fetchTempleRegisteredReport,
    downloadTempleRegisteredReport,
    // Devotee birthdays methods - NEW
    fetchDevoteeBirthdaysReport,
    downloadDevoteeBirthdaysReport,
    getDevoteeBirthdaysPreview,
    // Utility methods
    getDefaultDateRange,
    buildReportParams,
    buildDevoteeBirthdaysParams,
    formatReportData
  }
})
