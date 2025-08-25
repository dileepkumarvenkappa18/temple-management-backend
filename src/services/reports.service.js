// src/services/reports.service.js - FIXED VERSION
import api from '@/plugins/axios'

class ReportsService {
  /**
   * Get activities report data (JSON preview)
   */
  async getActivitiesReport(params) {
    const { entityId, type, dateRange = 'weekly', startDate, endDate } = params
    
    if (!entityId || !type) {
      throw new Error('Entity ID and type are required')
    }

    const queryParams = new URLSearchParams({
      type,
      date_range: dateRange
    })

    if (dateRange === 'custom' && startDate && endDate) {
      queryParams.append('start_date', startDate)
      queryParams.append('end_date', endDate)
    }

    try {
      console.log('Making API request:', `/v1/entities/${entityId}/reports/activities?${queryParams}`)
      const response = await api.get(`/v1/entities/${entityId}/reports/activities?${queryParams}`)
      console.log('API Response received:', response)
      
      return response
    } catch (error) {
      console.error('Error fetching activities report:', error)
      throw error
    }
  }

  /**
   * Download activities report in specified format
   */
  async downloadActivitiesReport(params) {
    const { entityId, type, format, dateRange = 'weekly', startDate, endDate } = params
    
    if (!entityId || !type || !format) {
      throw new Error('Entity ID, type, and format are required')
    }

    const queryParams = new URLSearchParams({
      type,
      date_range: dateRange,
      format
    })

    if (dateRange === 'custom' && startDate && endDate) {
      queryParams.append('start_date', startDate)
      queryParams.append('end_date', endDate)
    }

    try {
      const response = await api.get(`/v1/entities/${entityId}/reports/activities?${queryParams}`, {
        responseType: 'blob',
        headers: {
          'Accept': 'application/octet-stream'
        }
      })

      const contentDisposition = response.headers['content-disposition']
      let filename = `${type}_report.${format}`
      
      if (contentDisposition) {
        const filenameMatch = contentDisposition.match(/filename=(.+)/)
        if (filenameMatch) {
          filename = filenameMatch[1].replace(/"/g, '')
        }
      }

      const url = window.URL.createObjectURL(new Blob([response.data]))
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', filename)
      document.body.appendChild(link)
      link.click()
      link.remove()
      window.URL.revokeObjectURL(url)

      return { success: true, filename }
    } catch (error) {
      console.error('Error downloading activities report:', error)
      throw error
    }
  }

  /**
   * Get report preview data for display
   */
  async getReportPreview(params) {
    try {
      const response = await this.getActivitiesReport(params)
      
      console.log('API Response Structure:', response)
      console.log('Response Data:', response.data)
      
      let responseData = response.data
      
      if (responseData && responseData.data) {
        responseData = responseData.data
      }
      
      if (!responseData) {
        console.warn('No data found in response')
        responseData = {}
      }
      
      let previewData = []
      let columns = []

      switch (params.type) {
        case 'events':
          columns = [
            { key: 'title', label: 'Event Title' },
            { key: 'event_type', label: 'Type' },
            { key: 'event_date', label: 'Date' },
            { key: 'location', label: 'Location' },
            { key: 'created_by', label: 'Created By' }
          ]
          previewData = responseData.events || responseData.Events || []
          break
          
        case 'sevas':
          columns = [
            { key: 'name', label: 'Seva Name' },
            { key: 'seva_type', label: 'Type' },
            { key: 'price', label: 'Price' },
            { key: 'date', label: 'Date' },
            { key: 'status', label: 'Status' }
          ]
          previewData = responseData.sevas || responseData.Sevas || []
          break
          
        case 'bookings':
          columns = [
            { key: 'seva_name', label: 'Seva' },
            { key: 'devotee_name', label: 'Devotee' },
            { key: 'devotee_phone', label: 'Phone' },
            { key: 'booking_time', label: 'Booking Time' },
            { key: 'status', label: 'Status' }
          ]
          previewData = responseData.bookings || responseData.Bookings || []
          break
          
        case 'donations':
          columns = [
            { key: 'donor_name', label: 'Donor Name' },
            { key: 'amount', label: 'Amount' },
            { key: 'donation_type', label: 'Type' },
            { key: 'payment_method', label: 'Payment Method' },
            { key: 'status', label: 'Status' },
            { key: 'donation_date', label: 'Donation Date' }
          ]
          previewData = responseData.donations || responseData.Donations || []
          break
      }

      console.log('Extracted preview data:', previewData)
      console.log('Columns:', columns)

      return {
        data: previewData,
        columns,
        totalRecords: previewData.length
      }
    } catch (error) {
      console.error('Error getting report preview:', error)
      throw error
    }
  }

  /**
   * Validate report parameters
   */
  validateReportParams(params) {
    const errors = []

    if (!params.entityId) {
      errors.push('Entity ID is required')
    }

    if (!params.type || !['events', 'sevas', 'bookings', 'donations', 'temple-registered', 'devotee-birthdays'].includes(params.type)) {
      errors.push('Valid report type is required (events, sevas, bookings, donations, temple-registered, devotee-birthdays)')
    }

    if (params.type === 'temple-registered') {
      const validStatuses = ['approved', 'rejected', 'pending']
      if (params.status && !validStatuses.includes(params.status)) {
        errors.push(`Invalid status filter. Allowed values: ${validStatuses.join(', ')}`)
      }
    }

    if (params.dateRange === 'custom') {
      if (!params.startDate || !params.endDate) {
        errors.push('Start date and end date are required for custom date range')
      } else if (new Date(params.startDate) > new Date(params.endDate)) {
        errors.push('Start date must be before end date')
      }
    }

    if (params.format && !['pdf', 'csv', 'excel'].includes(params.format)) {
      errors.push('Invalid format specified')
    }

    return {
      isValid: errors.length === 0,
      errors
    }
  }

  // TEMPLE REGISTERED METHODS
  async getTempleRegisteredReport(params) {
    const { entityId, status, dateRange = 'weekly', startDate, endDate } = params

    if (!entityId) {
      throw new Error('Entity ID is required')
    }

    const queryParams = new URLSearchParams({
      date_range: dateRange,
    })

    if (status) {
      queryParams.append('status', status)
    }

    if (dateRange === 'custom' && startDate && endDate) {
      queryParams.append('start_date', startDate)
      queryParams.append('end_date', endDate)
    }

    try {
      const response = await api.get(`/v1/entities/${entityId}/reports/temple-registered?${queryParams}`)
      return response
    } catch (error) {
      console.error('Error fetching temple-registered report:', error)
      throw error
    }
  }

  async downloadTempleRegisteredReport(params) {
    return this.downloadReport(`/v1/entities/${params.entityId}/reports/temple-registered`, params, 'temple_registered_report')
  }

  async getTempleRegisteredPreview(params) {
    try {
      const response = await this.getTempleRegisteredReport(params)

      let responseData = response.data
      if (responseData && responseData.data) {
        responseData = responseData.data
      }

      const columns = [
        { key: 'temple_name', label: 'Temple Name' },
        { key: 'registration_date', label: 'Registration Date' },
        { key: 'location', label: 'Location' },
        { key: 'status', label: 'Status' }
      ]

      const previewData = responseData.temples || responseData || []

      return {
        data: previewData,
        columns,
        totalRecords: previewData.length || 0
      }
    } catch (error) {
      console.error('Error getting temple-registered preview:', error)
      throw error
    }
  }

  // DEVOTEE BIRTHDAYS METHODS - FIXED
  /**
   * Get devotee birthdays report data (JSON preview)
   */
  async getDevoteeBirthdaysReport(params) {
    const { entityId, dateRange = 'monthly', startDate, endDate } = params

    if (!entityId) {
      throw new Error('Entity ID is required')
    }

    const queryParams = new URLSearchParams({
      date_range: dateRange,
    })

    // FIXED: For birthdays, we need to handle date range differently
    if (dateRange === 'custom' && startDate && endDate) {
      queryParams.append('start_date', startDate)
      queryParams.append('end_date', endDate)
    } else {
      // For birthdays, set appropriate default dates
      const today = new Date()
      let calculatedStartDate, calculatedEndDate
      
      switch (dateRange) {
        case 'weekly':
          calculatedStartDate = new Date()
          calculatedEndDate = new Date()
          calculatedEndDate.setDate(calculatedEndDate.getDate() + 7)
          break
        case 'monthly':
          calculatedStartDate = new Date()
          calculatedEndDate = new Date()
          calculatedEndDate.setDate(calculatedEndDate.getDate() + 30)
          break
        case 'yearly':
          calculatedStartDate = new Date(today.getFullYear(), 0, 1)
          calculatedEndDate = new Date(today.getFullYear(), 11, 31)
          break
        default:
          calculatedStartDate = new Date()
          calculatedEndDate = new Date()
          calculatedEndDate.setDate(calculatedEndDate.getDate() + 30)
      }
      
      queryParams.append('start_date', calculatedStartDate.toISOString().split('T')[0])
      queryParams.append('end_date', calculatedEndDate.toISOString().split('T')[0])
    }

    try {
      console.log('🎂 Making devotee birthdays API request:', `/v1/entities/${entityId}/reports/devotee-birthdays?${queryParams}`)
      const response = await api.get(`/v1/entities/${entityId}/reports/devotee-birthdays?${queryParams}`)
      console.log('✅ Devotee birthdays API Response received:', response)
      return response
    } catch (error) {
      console.error('❌ Error fetching devotee-birthdays report:', error)
      throw error
    }
  }

  /**
   * Download devotee birthdays report in specified format - FIXED
   */
  async downloadDevoteeBirthdaysReport(params) {
    const { entityId, format, dateRange = 'monthly', startDate, endDate } = params

    if (!entityId || !format) {
      throw new Error('Entity ID and format are required')
    }

    const queryParams = new URLSearchParams({
      date_range: dateRange,
      format
    })

    // FIXED: Same date logic as preview method
    if (dateRange === 'custom' && startDate && endDate) {
      queryParams.append('start_date', startDate)
      queryParams.append('end_date', endDate)
    } else {
      const today = new Date()
      let calculatedStartDate, calculatedEndDate
      
      switch (dateRange) {
        case 'weekly':
          calculatedStartDate = new Date()
          calculatedEndDate = new Date()
          calculatedEndDate.setDate(calculatedEndDate.getDate() + 7)
          break
        case 'monthly':
          calculatedStartDate = new Date()
          calculatedEndDate = new Date()
          calculatedEndDate.setDate(calculatedEndDate.getDate() + 30)
          break
        case 'yearly':
          calculatedStartDate = new Date(today.getFullYear(), 0, 1)
          calculatedEndDate = new Date(today.getFullYear(), 11, 31)
          break
        default:
          calculatedStartDate = new Date()
          calculatedEndDate = new Date()
          calculatedEndDate.setDate(calculatedEndDate.getDate() + 30)
      }
      
      queryParams.append('start_date', calculatedStartDate.toISOString().split('T')[0])
      queryParams.append('end_date', calculatedEndDate.toISOString().split('T')[0])
    }

    const apiUrl = `/v1/entities/${entityId}/reports/devotee-birthdays?${queryParams}`
    
    return this.downloadReport(apiUrl, { format }, 'devotee_birthdays_report')
  }

  /**
   * Get devotee birthdays report preview data - FIXED
   */
  async getDevoteeBirthdaysPreview(params) {
    try {
      const response = await this.getDevoteeBirthdaysReport(params)

      console.log('🎂 Devotee birthdays preview response:', response)

      let responseData = response.data
      
      // Handle different response structures
      if (responseData && responseData.data) {
        responseData = responseData.data
      }

      console.log('📊 Processed response data:', responseData)

      const columns = [
        { key: 'full_name', label: 'Full Name' },
        { key: 'date_of_birth', label: 'Date of Birth' },
        { key: 'age', label: 'Age' },
        { key: 'gender', label: 'Gender' },
        { key: 'phone', label: 'Phone' },
        { key: 'email', label: 'Email' },
        { key: 'temple_name', label: 'Temple' },
        { key: 'upcoming_birthday', label: 'Upcoming Birthday' }
      ]

      // FIXED: Handle different response formats more robustly
      let previewData = []
      
      if (Array.isArray(responseData)) {
        previewData = responseData
      } else if (responseData && responseData.devotees) {
        previewData = responseData.devotees
      } else if (responseData && responseData.birthdays) {
        previewData = responseData.birthdays
      } else if (responseData && responseData.devotee_birthdays) {
        previewData = responseData.devotee_birthdays
      } else if (responseData && Array.isArray(responseData.data)) {
        previewData = responseData.data
      } else {
        console.warn('⚠️ No valid devotee birthdays data found in response:', responseData)
        previewData = []
      }

      // Format date fields for better display
      previewData = previewData.map(item => ({
        ...item,
        date_of_birth: item.date_of_birth ? this.formatDate(item.date_of_birth) : '-',
        upcoming_birthday: item.upcoming_birthday ? this.formatDate(item.upcoming_birthday) : '-'
      }))

      console.log('📋 Final preview data:', previewData)

      return {
        data: previewData,
        columns,
        totalRecords: previewData.length || 0
      }
    } catch (error) {
      console.error('❌ Error getting devotee-birthdays preview:', error)
      throw error
    }
  }

  // DEVOTEE LIST METHODS
  async getDevoteeList(params) {
    const { entityId, status = 'all' } = params

    if (!entityId) {
      throw new Error('Entity ID is required')
    }

    const allowedStatuses = ['all', 'active', 'inactive']
    if (!allowedStatuses.includes(status)) {
      throw new Error(`Invalid status. Allowed values: ${allowedStatuses.join(', ')}`)
    }

    const queryParams = new URLSearchParams()
    if (status !== 'all') {
      queryParams.append('status', status)
    }

    try {
      console.log('📋 Making devotee list API request:', `/v1/entities/${entityId}/reports/devotee-list?${queryParams}`)
      const response = await api.get(`/v1/entities/${entityId}/reports/devotee-list?${queryParams}`)
      console.log('✅ Devotee list API Response received:', response)
      return response
    } catch (error) {
      console.error('❌ Error fetching devotee list:', error)
      throw error
    }
  }

  async downloadDevoteeListReport(params) {
    const { entityId, status = 'all', format } = params

    if (!entityId || !format) {
      throw new Error('Entity ID and format are required')
    }

    const queryParams = new URLSearchParams({ format })
    if (status !== 'all') {
      queryParams.append('status', status)
    }

    const apiUrl = `/v1/entities/${entityId}/reports/devotee-list?${queryParams}`
    return this.downloadReport(apiUrl, { format }, `devotee_list_${status}_report`)
  }

  async getDevoteeListPreview(params) {
    try {
      const response = await this.getDevoteeList(params)

      let responseData = response.data
      if (responseData && responseData.data) {
        responseData = responseData.data
      }

      const previewData = responseData.devotees || responseData || []

      const columns = [
        { key: 'full_name', label: 'Full Name' },
        { key: 'phone', label: 'Phone' },
        { key: 'email', label: 'Email' },
        { key: 'status', label: 'Status' },
        { key: 'registration_date', label: 'Registration Date' },
        { key: 'last_login', label: 'Last Login' }
      ]

      return {
        data: previewData,
        columns,
        totalRecords: previewData.length || 0
      }
    } catch (error) {
      console.error('Error getting devotee list preview:', error)
      throw error
    }
  }

  // DEVOTEE PROFILE METHODS
  async getDevoteeProfile(params) {
    const { entityId } = params

    if (!entityId) {
      throw new Error('Entity ID is required')
    }

    try {
      console.log('👤 Making devotee profile API request:', `/v1/entities/${entityId}/reports/devotee-profile`)
      const response = await api.get(`/v1/entities/${entityId}/reports/devotee-profile`)
      console.log('✅ Devotee profile API Response received:', response)
      return response
    } catch (error) {
      console.error('❌ Error fetching devotee profile:', error)
      throw error
    }
  }

  async downloadDevoteeProfileReport(params) {
    const { entityId, format } = params

    if (!entityId || !format) {
      throw new Error('Entity ID and format are required')
    }

    const queryParams = new URLSearchParams({ format })
    const apiUrl = `/v1/entities/${entityId}/reports/devotee-profile?${queryParams}`
    return this.downloadReport(apiUrl, { format }, 'devotee_profile_report')
  }

  async getDevoteeProfilePreview(params) {
    try {
      const response = await this.getDevoteeProfile(params)

      let responseData = response.data
      if (responseData && responseData.data) {
        responseData = responseData.data
      }

      const profileData = responseData.profile || responseData || {}

      return profileData
    } catch (error) {
      console.error('Error getting devotee profile preview:', error)
      throw error
    }
  }

  // UTILITY METHODS - FIXED AND CONSOLIDATED

  /**
   * Generic download method - FIXED
   */
  async downloadReport(apiUrl, params, defaultFilename) {
    try {
      console.log('🔄 Making download request:', apiUrl)
      console.log('📋 Request parameters:', params)

      const response = await api.get(apiUrl, {
        responseType: 'blob',
        headers: {
          'Accept': this.getAcceptHeader(params.format),
          'Cache-Control': 'no-cache'
        }
      })

      console.log('✅ Download response received:', {
        status: response.status,
        contentType: response.headers['content-type'],
        contentLength: response.headers['content-length'],
        contentDisposition: response.headers['content-disposition']
      })

      // Validate response
      if (!response.data || response.data.size === 0) {
        throw new Error('Empty file received from server')
      }

      // Check for error responses disguised as success
      const contentType = response.headers['content-type']
      if (contentType?.includes('text/html') || contentType?.includes('application/json')) {
        const text = await response.data.text()
        console.error('❌ Server returned error response:', text)
        throw new Error(`Server error: ${text.substring(0, 200)}`)
      }

      // Extract filename
      let filename = `${defaultFilename}.${params.format}`
      const contentDisposition = response.headers['content-disposition']
      
      if (contentDisposition) {
        const patterns = [
          /filename[^;=\n]*=\s*"([^"]+)"/i,
          /filename[^;=\n]*=\s*'([^']+)'/i,
          /filename[^;=\n]*=\s*([^;\n]+)/i,
          /filename\*=UTF-8''([^;\n]+)/i
        ]
        
        for (const pattern of patterns) {
          const match = contentDisposition.match(pattern)
          if (match) {
            filename = decodeURIComponent(match[1].trim())
            break
          }
        }
      }

      // Create and trigger download
      const blob = new Blob([response.data], {
        type: this.getBlobType(params.format)
      })
      
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', filename)
      link.style.display = 'none'
      
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      
      setTimeout(() => {
        window.URL.revokeObjectURL(url)
      }, 1000)

      console.log('✅ Download completed successfully')
      return { 
        success: true, 
        filename, 
        size: blob.size,
        contentType: contentType
      }

    } catch (error) {
      console.error('❌ Error in download:', error)
      
      if (error.response) {
        console.error('Error response details:', {
          status: error.response.status,
          statusText: error.response.statusText,
          contentType: error.response.headers['content-type']
        })
        
        if (error.response.data instanceof Blob) {
          try {
            const errorText = await error.response.data.text()
            throw new Error(`Server error (${error.response.status}): ${errorText}`)
          } catch (readError) {
            console.error('Failed to read error response:', readError)
          }
        }
        
        throw new Error(`Server returned ${error.response.status}: ${error.response.statusText}`)
      }
      
      throw error
    }
  }

  /**
   * Format date helper
   */
  formatDate(dateString) {
    if (!dateString) return '-'
    try {
      const date = new Date(dateString)
      return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      })
    } catch (error) {
      console.warn('Invalid date format:', dateString)
      return dateString
    }
  }

  /**
   * Get accept header based on format
   */
  getAcceptHeader(format) {
    switch (format?.toLowerCase()) {
      case 'pdf':
        return 'application/pdf'
      case 'excel':
      case 'xlsx':
        return 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
      case 'csv':
        return 'text/csv'
      default:
        return 'application/octet-stream'
    }
  }

  /**
   * Get blob MIME type based on format
   */
  getBlobType(format) {
    switch (format?.toLowerCase()) {
      case 'pdf':
        return 'application/pdf'
      case 'excel':
      case 'xlsx':
        return 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
      case 'csv':
        return 'text/csv'
      default:
        return 'application/octet-stream'
    }
  }
}

export default new ReportsService()