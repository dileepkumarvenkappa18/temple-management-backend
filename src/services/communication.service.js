// src/services/communication.service.js

import api from './api.js'

/**
 * Communication Service
 * 
 * Handles all communication-related API operations for the temple management system
 * including bulk messaging, templates, notifications, and message history
 */
class CommunicationService {
  
  // ==================== MESSAGE TEMPLATES ====================
  
  /**
   * Get all message templates for an entity
   * @param {string} entityId - Temple entity ID
   * @returns {Promise} Array of message templates
   */
  async getTemplates(entityId) {
    try {
      const response = await api.get(`/entities/${entityId}/communication/templates`)
      return {
        success: true,
        data: response.data,
        message: 'Templates fetched successfully'
      }
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.message || 'Failed to fetch templates',
        status: error.response?.status
      }
    }
  }

  /**
   * Create a new message template
   * @param {string} entityId - Temple entity ID
   * @param {Object} templateData - Template information
   * @returns {Promise} Created template
   */
  async createTemplate(entityId, templateData) {
    try {
      const response = await api.post(`/entities/${entityId}/communication/templates`, {
        name: templateData.name,
        subject: templateData.subject,
        content: templateData.content,
        type: templateData.type, // 'sms', 'email', 'whatsapp', 'notification'
        category: templateData.category, // 'seva', 'donation', 'event', 'general'
        isActive: templateData.isActive || true,
        variables: templateData.variables || [] // Available placeholders like {devotee_name}, {seva_name}
      })
      
      return {
        success: true,
        data: response.data,
        message: 'Template created successfully'
      }
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.message || 'Failed to create template',
        status: error.response?.status
      }
    }
  }

  /**
   * Update an existing message template
   * @param {string} entityId - Temple entity ID
   * @param {string} templateId - Template ID
   * @param {Object} templateData - Updated template information
   * @returns {Promise} Updated template
   */
  async updateTemplate(entityId, templateId, templateData) {
    try {
      const response = await api.put(`/entities/${entityId}/communication/templates/${templateId}`, templateData)
      return {
        success: true,
        data: response.data,
        message: 'Template updated successfully'
      }
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.message || 'Failed to update template',
        status: error.response?.status
      }
    }
  }

  /**
   * Delete a message template
   * @param {string} entityId - Temple entity ID
   * @param {string} templateId - Template ID
   * @returns {Promise} Deletion confirmation
   */
  async deleteTemplate(entityId, templateId) {
    try {
      await api.delete(`/entities/${entityId}/communication/templates/${templateId}`)
      return {
        success: true,
        message: 'Template deleted successfully'
      }
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.message || 'Failed to delete template',
        status: error.response?.status
      }
    }
  }

  // ==================== BULK MESSAGING ====================

  /**
   * Send bulk message to filtered devotees
   * @param {string} entityId - Temple entity ID
   * @param {Object} messageData - Message configuration
   * @returns {Promise} Message sending result
   */
  async sendBulkMessage(entityId, messageData) {
    try {
      const response = await api.post(`/entities/${entityId}/communication/bulk-send`, {
        type: messageData.type, // 'sms', 'email', 'whatsapp', 'notification'
        subject: messageData.subject,
        content: messageData.content,
        templateId: messageData.templateId,
        recipients: {
          type: messageData.recipients.type, // 'all', 'filtered', 'selected'
          filters: messageData.recipients.filters, // Age, gender, seva preferences, etc.
          devoteeIds: messageData.recipients.devoteeIds, // Specific devotee IDs
          excludeIds: messageData.recipients.excludeIds // Devotees to exclude
        },
        scheduling: {
          sendNow: messageData.sendNow || true,
          scheduledAt: messageData.scheduledAt, // Future date/time
          timezone: messageData.timezone || 'Asia/Kolkata'
        },
        attachments: messageData.attachments || [],
        trackOpens: messageData.trackOpens || false,
        trackClicks: messageData.trackClicks || false
      })
      
      return {
        success: true,
        data: response.data,
        message: `Message ${messageData.sendNow ? 'sent' : 'scheduled'} successfully`
      }
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.message || 'Failed to send message',
        status: error.response?.status
      }
    }
  }

  /**
   * Preview bulk message before sending
   * @param {string} entityId - Temple entity ID
   * @param {Object} messageData - Message configuration
   * @returns {Promise} Preview data with recipient count and sample content
   */
  async previewBulkMessage(entityId, messageData) {
    try {
      const response = await api.post(`/entities/${entityId}/communication/preview`, {
        type: messageData.type,
        content: messageData.content,
        templateId: messageData.templateId,
        recipients: messageData.recipients
      })
      
      return {
        success: true,
        data: {
          recipientCount: response.data.recipientCount,
          sampleRecipients: response.data.sampleRecipients,
          renderedContent: response.data.renderedContent,
          estimatedCost: response.data.estimatedCost,
          deliveryTime: response.data.deliveryTime
        }
      }
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.message || 'Failed to preview message',
        status: error.response?.status
      }
    }
  }

  // ==================== MESSAGE HISTORY & LOGS ====================

  /**
   * Get message history for an entity
   * @param {string} entityId - Temple entity ID
   * @param {Object} params - Query parameters
   * @returns {Promise} Message history
   */
  async getMessageHistory(entityId, params = {}) {
    try {
      const queryParams = new URLSearchParams({
        page: params.page || 1,
        limit: params.limit || 20,
        type: params.type || '', // Filter by message type
        status: params.status || '', // 'sent', 'failed', 'scheduled', 'draft'
        dateFrom: params.dateFrom || '',
        dateTo: params.dateTo || '',
        search: params.search || ''
      }).toString()

      const response = await api.get(`/entities/${entityId}/communication/history?${queryParams}`)
      
      return {
        success: true,
        data: response.data.messages,
        pagination: response.data.pagination,
        stats: response.data.stats
      }
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.message || 'Failed to fetch message history',
        status: error.response?.status
      }
    }
  }

  /**
   * Get detailed message log
   * @param {string} entityId - Temple entity ID
   * @param {string} messageId - Message ID
   * @returns {Promise} Detailed message information
   */
  async getMessageDetails(entityId, messageId) {
    try {
      const response = await api.get(`/entities/${entityId}/communication/messages/${messageId}`)
      
      return {
        success: true,
        data: {
          message: response.data.message,
          deliveryStats: response.data.deliveryStats,
          recipients: response.data.recipients,
          timeline: response.data.timeline
        }
      }
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.message || 'Failed to fetch message details',
        status: error.response?.status
      }
    }
  }

  // ==================== RECIPIENT MANAGEMENT ====================

  /**
   * Get devotee list with filters for messaging
   * @param {string} entityId - Temple entity ID
   * @param {Object} filters - Filter criteria
   * @returns {Promise} Filtered devotee list
   */
  async getDevoteesForMessaging(entityId, filters = {}) {
    try {
      const queryParams = new URLSearchParams({
        search: filters.search || '',
        ageFrom: filters.ageFrom || '',
        ageTo: filters.ageTo || '',
        gender: filters.gender || '',
        city: filters.city || '',
        state: filters.state || '',
        sevaPreferences: filters.sevaPreferences || '',
        membershipType: filters.membershipType || '',
        isActive: filters.isActive !== undefined ? filters.isActive : true,
        hasPhone: filters.hasPhone || false,
        hasEmail: filters.hasEmail || false,
        hasWhatsApp: filters.hasWhatsApp || false
      }).toString()

      const response = await api.get(`/entities/${entityId}/devotees/for-messaging?${queryParams}`)
      
      return {
        success: true,
        data: response.data.devotees,
        count: response.data.count,
        filters: response.data.availableFilters
      }
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.message || 'Failed to fetch devotees',
        status: error.response?.status
      }
    }
  }

  // ==================== NOTIFICATIONS ====================

  /**
   * Send notification to specific devotees
   * @param {string} entityId - Temple entity ID
   * @param {Object} notificationData - Notification details
   * @returns {Promise} Notification result
   */
  async sendNotification(entityId, notificationData) {
    try {
      const response = await api.post(`/entities/${entityId}/communication/notifications`, {
        title: notificationData.title,
        body: notificationData.body,
        type: notificationData.type, // 'seva', 'donation', 'event', 'general', 'urgent'
        priority: notificationData.priority || 'normal', // 'low', 'normal', 'high', 'urgent'
        recipients: notificationData.recipients,
        actionUrl: notificationData.actionUrl,
        imageUrl: notificationData.imageUrl,
        expireAt: notificationData.expireAt,
        data: notificationData.data || {} // Additional payload
      })
      
      return {
        success: true,
        data: response.data,
        message: 'Notification sent successfully'
      }
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.message || 'Failed to send notification',
        status: error.response?.status
      }
    }
  }

  /**
   * Get unread notifications for a user
   * @param {string} entityId - Temple entity ID
   * @param {string} userId - User ID
   * @returns {Promise} Unread notifications
   */
  async getUnreadNotifications(entityId, userId) {
    try {
      const response = await api.get(`/entities/${entityId}/users/${userId}/notifications/unread`)
      
      return {
        success: true,
        data: response.data.notifications,
        count: response.data.count
      }
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.message || 'Failed to fetch notifications',
        status: error.response?.status
      }
    }
  }

  /**
   * Mark notification as read
   * @param {string} entityId - Temple entity ID
   * @param {string} userId - User ID
   * @param {string} notificationId - Notification ID
   * @returns {Promise} Success confirmation
   */
  async markNotificationAsRead(entityId, userId, notificationId) {
    try {
      await api.put(`/entities/${entityId}/users/${userId}/notifications/${notificationId}/read`)
      
      return {
        success: true,
        message: 'Notification marked as read'
      }
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.message || 'Failed to mark notification as read',
        status: error.response?.status
      }
    }
  }

  // ==================== COMMUNICATION STATS ====================

  /**
   * Get communication statistics for dashboard
   * @param {string} entityId - Temple entity ID
   * @param {Object} params - Date range and filters
   * @returns {Promise} Communication statistics
   */
  async getCommunicationStats(entityId, params = {}) {
    try {
      const queryParams = new URLSearchParams({
        dateFrom: params.dateFrom || '',
        dateTo: params.dateTo || '',
        type: params.type || ''
      }).toString()

      const response = await api.get(`/entities/${entityId}/communication/stats?${queryParams}`)
      
      return {
        success: true,
        data: {
          totalSent: response.data.totalSent,
          totalDelivered: response.data.totalDelivered,
          totalFailed: response.data.totalFailed,
          openRate: response.data.openRate,
          clickRate: response.data.clickRate,
          byType: response.data.byType,
          byDate: response.data.byDate,
          topTemplates: response.data.topTemplates
        }
      }
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.message || 'Failed to fetch communication stats',
        status: error.response?.status
      }
    }
  }

  // ==================== UTILITY METHODS ====================

  /**
   * Validate message content and recipients
   * @param {Object} messageData - Message data to validate
   * @returns {Object} Validation result
   */
  validateMessage(messageData) {
    const errors = []

    if (!messageData.type) {
      errors.push('Message type is required')
    }

    if (!messageData.content && !messageData.templateId) {
      errors.push('Message content or template is required')
    }

    if (messageData.type === 'email' && !messageData.subject) {
      errors.push('Email subject is required')
    }

    if (!messageData.recipients || 
        (!messageData.recipients.devoteeIds?.length && 
         messageData.recipients.type !== 'all' &&
         !messageData.recipients.filters)) {
      errors.push('Recipients are required')
    }

    if (messageData.content && messageData.content.length > 1600 && messageData.type === 'sms') {
      errors.push('SMS content cannot exceed 1600 characters')
    }

    return {
      isValid: errors.length === 0,
      errors
    }
  }

  /**
   * Replace template variables with actual values
   * @param {string} content - Template content with variables
   * @param {Object} variables - Variable values
   * @returns {string} Rendered content
   */
  renderTemplate(content, variables = {}) {
    let renderedContent = content

    // Common variables
    const defaultVariables = {
      temple_name: variables.templeName || '[Temple Name]',
      devotee_name: variables.devoteeName || '[Devotee Name]',
      current_date: new Date().toLocaleDateString('en-IN'),
      current_time: new Date().toLocaleTimeString('en-IN'),
      ...variables
    }

    // Replace variables in format {variable_name}
    Object.entries(defaultVariables).forEach(([key, value]) => {
      const regex = new RegExp(`\\{${key}\\}`, 'gi')
      renderedContent = renderedContent.replace(regex, value)
    })

    return renderedContent
  }
}

// Export singleton instance
export default new CommunicationService()