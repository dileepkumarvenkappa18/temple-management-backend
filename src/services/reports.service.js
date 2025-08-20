// src/services/reports.service.js
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

    if (!params.type || !['events', 'sevas', 'bookings', 'temple-registered', 'devotee-birthdays'].includes(params.type)) {
      errors.push('Valid report type is required (events, sevas, bookings, temple-registered, devotee-birthdays)')
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

  /**
   * ENHANCED Download temple registered report with detailed debugging
   */
  async downloadTempleRegisteredReport(params) {
    const { entityId, status, format, dateRange = 'weekly', startDate, endDate } = params

    if (!entityId || !format) {
      throw new Error('Entity ID and format are required')
    }

    const queryParams = new URLSearchParams({
      date_range: dateRange,
      format
    })

    if (status) {
      queryParams.append('status', status)
    }

    if (dateRange === 'custom' && startDate && endDate) {
      queryParams.append('start_date', startDate)
      queryParams.append('end_date', endDate)
    }

    const apiUrl = `/v1/entities/${entityId}/reports/temple-registered?${queryParams}`
    
    try {
      console.log('🔄 Making download request:', apiUrl)
      console.log('📋 Request parameters:', params)
      console.log('🎯 Expected format:', format)

      const response = await api.get(apiUrl, {
        responseType: 'blob',
        headers: {
          'Accept': this.getAcceptHeader(format),
          'Cache-Control': 'no-cache'
        }
      })

      console.log('✅ Response received:', {
        status: response.status,
        statusText: response.statusText,
        headers: Object.fromEntries(Object.entries(response.headers)),
        contentType: response.headers['content-type'],
        contentLength: response.headers['content-length'],
        contentDisposition: response.headers['content-disposition']
      })

      // Check if we actually got a blob
      if (!response.data) {
        throw new Error('No response data received from server')
      }

      // Log blob details
      console.log('📊 Blob details:', {
        size: response.data.size,
        type: response.data.type,
        constructor: response.data.constructor.name
      })

      // Check if blob is empty
      if (response.data.size === 0) {
        throw new Error('Empty file received from server')
      }

      // For debugging, let's peek into the blob content for small files
      if (response.data.size < 1000) {
        try {
          const arrayBuffer = await response.data.arrayBuffer()
          const uint8Array = new Uint8Array(arrayBuffer)
          const firstBytes = Array.from(uint8Array.slice(0, 50))
          console.log('🔍 First 50 bytes:', firstBytes)
          
          // Try to read as text to see if it's actually text content
          const textContent = new TextDecoder().decode(uint8Array.slice(0, 200))
          console.log('📄 Content preview:', textContent)
          
          // Check if it starts with typical file signatures
          const signature = this.detectFileSignature(uint8Array)
          console.log('🎭 Detected file type:', signature)
          
          // Recreate blob from arrayBuffer for download
          response.data = new Blob([arrayBuffer], { type: this.getBlobType(format) })
        } catch (error) {
          console.warn('⚠️ Could not analyze blob content:', error)
        }
      }

      // Check content type against expected
      const contentType = response.headers['content-type']
      const expectedTypes = this.getExpectedContentType(format)
      
      console.log('🎯 Content type check:', {
        received: contentType,
        expected: expectedTypes,
        matches: expectedTypes.some(type => contentType?.includes(type))
      })

      // Check for error responses disguised as success
      if (contentType?.includes('text/html') || contentType?.includes('application/json')) {
        const text = await response.data.text()
        console.error('❌ Server returned error response instead of file:', text)
        throw new Error(`Server returned ${contentType} instead of ${format} file. Response: ${text.substring(0, 200)}`)
      }

      // Extract filename from Content-Disposition header
      const contentDisposition = response.headers['content-disposition']
      let filename = `temple_registered_report.${format}`
      
      if (contentDisposition) {
        console.log('📁 Content-Disposition header:', contentDisposition)
        
        // Try multiple patterns for filename extraction
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
            console.log('📁 Extracted filename:', filename)
            break
          }
        }
      } else {
        console.warn('⚠️ No Content-Disposition header found')
      }

      // Validate file extension matches format
      const fileExtension = filename.split('.').pop()?.toLowerCase()
      const expectedExtension = this.getExpectedExtension(format)
      
      if (fileExtension !== expectedExtension) {
        console.warn('⚠️ File extension mismatch:', {
          filename,
          detectedExtension: fileExtension,
          expectedExtension
        })
        
        // Fix the filename if needed
        if (!filename.endsWith(`.${expectedExtension}`)) {
          filename = filename.replace(/\.[^.]+$/, `.${expectedExtension}`)
          console.log('🔧 Corrected filename:', filename)
        }
      }

      // Create download with explicit MIME type
      const blob = new Blob([response.data], {
        type: this.getBlobType(format)
      })
      
      console.log('💾 Final download blob:', {
        size: blob.size,
        type: blob.type
      })

      const url = window.URL.createObjectURL(blob)
      
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', filename)
      link.style.display = 'none'
      
      // Add some attributes for better compatibility
      link.setAttribute('type', this.getBlobType(format))
      link.setAttribute('data-format', format)
      
      document.body.appendChild(link)
      
      // Log the download attempt
      console.log('🚀 Initiating download:', {
        url: url.substring(0, 50) + '...',
        filename,
        mimeType: this.getBlobType(format)
      })
      
      link.click()
      
      // Clean up
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
      console.error('❌ Error downloading temple-registered report:', error)
      
      // Enhanced error reporting
      if (error.response) {
        console.error('📋 Error response details:', {
          status: error.response.status,
          statusText: error.response.statusText,
          headers: error.response.headers,
          contentType: error.response.headers['content-type']
        })
        
        // Try to read error response
        if (error.response.data instanceof Blob) {
          try {
            const errorText = await error.response.data.text()
            console.error('📄 Error response body:', errorText)
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

  // DEVOTEE BIRTHDAYS METHODS - NEW
  /**
   * Get devotee birthdays report data (JSON preview)
   */
  async getDevoteeBirthdaysReport(params) {
    const { entityId, dateRange = 'weekly', startDate, endDate } = params

    if (!entityId) {
      throw new Error('Entity ID is required')
    }

    const queryParams = new URLSearchParams({
      date_range: dateRange,
    })

    if (dateRange === 'custom' && startDate && endDate) {
      queryParams.append('start_date', startDate)
      queryParams.append('end_date', endDate)
    }

    try {
      console.log('Making devotee birthdays API request:', `/v1/entities/${entityId}/reports/devotee-birthdays?${queryParams}`)
      const response = await api.get(`/v1/entities/${entityId}/reports/devotee-birthdays?${queryParams}`)
      console.log('Devotee birthdays API Response received:', response)
      return response
    } catch (error) {
      console.error('Error fetching devotee-birthdays report:', error)
      throw error
    }
  }

  /**
   * Download devotee birthdays report in specified format
   */
  async downloadDevoteeBirthdaysReport(params) {
    const { entityId, format, dateRange = 'weekly', startDate, endDate } = params

    if (!entityId || !format) {
      throw new Error('Entity ID and format are required')
    }

    const queryParams = new URLSearchParams({
      date_range: dateRange,
      format
    })

    if (dateRange === 'custom' && startDate && endDate) {
      queryParams.append('start_date', startDate)
      queryParams.append('end_date', endDate)
    }

    const apiUrl = `/v1/entities/${entityId}/reports/devotee-birthdays?${queryParams}`
    
    try {
      console.log('🔄 Making devotee birthdays download request:', apiUrl)
      console.log('📋 Request parameters:', params)
      console.log('🎯 Expected format:', format)

      const response = await api.get(apiUrl, {
        responseType: 'blob',
        headers: {
          'Accept': this.getAcceptHeader(format),
          'Cache-Control': 'no-cache'
        }
      })

      console.log('✅ Devotee birthdays response received:', {
        status: response.status,
        statusText: response.statusText,
        contentType: response.headers['content-type'],
        contentLength: response.headers['content-length'],
        contentDisposition: response.headers['content-disposition']
      })

      // Check if we actually got a blob
      if (!response.data) {
        throw new Error('No response data received from server')
      }

      // Check if blob is empty
      if (response.data.size === 0) {
        throw new Error('Empty file received from server')
      }

      // Check content type against expected
      const contentType = response.headers['content-type']
      const expectedTypes = this.getExpectedContentType(format)
      
      console.log('🎯 Content type check:', {
        received: contentType,
        expected: expectedTypes,
        matches: expectedTypes.some(type => contentType?.includes(type))
      })

      // Check for error responses disguised as success
      if (contentType?.includes('text/html') || contentType?.includes('application/json')) {
        const text = await response.data.text()
        console.error('❌ Server returned error response instead of file:', text)
        throw new Error(`Server returned ${contentType} instead of ${format} file. Response: ${text.substring(0, 200)}`)
      }

      // Extract filename from Content-Disposition header
      const contentDisposition = response.headers['content-disposition']
      let filename = `devotee_birthdays_report.${format}`
      
      if (contentDisposition) {
        console.log('📁 Content-Disposition header:', contentDisposition)
        
        // Try multiple patterns for filename extraction
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
            console.log('📁 Extracted filename:', filename)
            break
          }
        }
      } else {
        console.warn('⚠️ No Content-Disposition header found')
      }

      // Validate file extension matches format
      const fileExtension = filename.split('.').pop()?.toLowerCase()
      const expectedExtension = this.getExpectedExtension(format)
      
      if (fileExtension !== expectedExtension) {
        console.warn('⚠️ File extension mismatch:', {
          filename,
          detectedExtension: fileExtension,
          expectedExtension
        })
        
        // Fix the filename if needed
        if (!filename.endsWith(`.${expectedExtension}`)) {
          filename = filename.replace(/\.[^.]+$/, `.${expectedExtension}`)
          console.log('🔧 Corrected filename:', filename)
        }
      }

      // Create download with explicit MIME type
      const blob = new Blob([response.data], {
        type: this.getBlobType(format)
      })
      
      console.log('💾 Final download blob:', {
        size: blob.size,
        type: blob.type
      })

      const url = window.URL.createObjectURL(blob)
      
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', filename)
      link.style.display = 'none'
      
      // Add some attributes for better compatibility
      link.setAttribute('type', this.getBlobType(format))
      link.setAttribute('data-format', format)
      
      document.body.appendChild(link)
      
      // Log the download attempt
      console.log('🚀 Initiating download:', {
        url: url.substring(0, 50) + '...',
        filename,
        mimeType: this.getBlobType(format)
      })
      
      link.click()
      
      // Clean up
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
      console.error('❌ Error downloading devotee-birthdays report:', error)
      
      // Enhanced error reporting
      if (error.response) {
        console.error('📋 Error response details:', {
          status: error.response.status,
          statusText: error.response.statusText,
          headers: error.response.headers,
          contentType: error.response.headers['content-type']
        })
        
        // Try to read error response
        if (error.response.data instanceof Blob) {
          try {
            const errorText = await error.response.data.text()
            console.error('📄 Error response body:', errorText)
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
   * Get devotee birthdays report preview data
   */
  async getDevoteeBirthdaysPreview(params) {
    try {
      const response = await this.getDevoteeBirthdaysReport(params)

      console.log('Devotee birthdays preview response:', response)

      let responseData = response.data
      
      // Handle different response structures
      if (responseData && responseData.data) {
        responseData = responseData.data
      }

      // Log the actual response structure for debugging
      console.log('Processed response data:', responseData)

      const columns = [
        { key: 'full_name', label: 'Full Name' },
        { key: 'date_of_birth', label: 'Date of Birth' },
        { key: 'gender', label: 'Gender' },
        { key: 'phone', label: 'Phone' },
        { key: 'email', label: 'Email' },
        { key: 'temple_name', label: 'Temple' },
        { key: 'member_since', label: 'Member Since' }
      ]

      // Handle different response formats
      let previewData = []
      
      if (Array.isArray(responseData)) {
        // Direct array response
        previewData = responseData
      } else if (responseData && responseData.devotee_birthdays) {
        // Wrapped in devotee_birthdays property
        previewData = responseData.devotee_birthdays
      } else if (responseData && Array.isArray(responseData.data)) {
        // Wrapped in data property as array
        previewData = responseData.data
      } else {
        // Default to empty array if no valid data found
        console.warn('No valid devotee birthdays data found in response:', responseData)
        previewData = []
      }

      console.log('Final preview data:', previewData)

      return {
        data: previewData,
        columns,
        totalRecords: previewData.length || 0
      }
    } catch (error) {
      console.error('Error getting devotee-birthdays preview:', error)
      throw error
    }
  }
  

  // HELPER METHODS FOR ENHANCED DEBUGGING
  
  detectFileSignature(uint8Array) {
    const signatures = {
      'PDF': [0x25, 0x50, 0x44, 0x46], // %PDF
      'ZIP/Excel': [0x50, 0x4B, 0x03, 0x04], // PK.. (Excel files are ZIP-based)
      'Excel Old': [0xD0, 0xCF, 0x11, 0xE0], // Old Excel format
      'CSV/Text': [] // Text files don't have a specific signature
    }
    
    for (const [type, signature] of Object.entries(signatures)) {
      if (signature.length === 0) continue // Skip text check
      if (uint8Array.length >= signature.length) {
        const matches = signature.every((byte, index) => uint8Array[index] === byte)
        if (matches) return type
      }
    }
    
    // Check if it looks like text/CSV
    const firstChunk = new TextDecoder().decode(uint8Array.slice(0, 100))
    if (/^[a-zA-Z0-9\s,";'\r\n-]+$/.test(firstChunk)) {
      return 'CSV/Text'
    }
    
    return 'Unknown'
  }

  getExpectedExtension(format) {
    switch (format.toLowerCase()) {
      case 'pdf':
        return 'pdf'
      case 'excel':
      case 'xlsx':
        return 'xlsx'
      case 'csv':
        return 'csv'
      default:
        return format.toLowerCase()
    }
  }

  getAcceptHeader(format) {
    switch (format.toLowerCase()) {
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

  getExpectedContentType(format) {
    switch (format.toLowerCase()) {
      case 'pdf':
        return ['application/pdf']
      case 'excel':
      case 'xlsx':
        return [
          'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
          'application/vnd.ms-excel',
          'application/octet-stream' // Some servers use generic type
        ]
      case 'csv':
        return ['text/csv', 'application/csv', 'text/plain']
      default:
        return ['application/octet-stream']
    }
  }

  getBlobType(format) {
    switch (format.toLowerCase()) {
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