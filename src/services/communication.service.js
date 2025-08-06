// src/services/communication.service.js
import { apiClient } from '@/plugins/axios';

const communicationService = {
  async getTemplates() {
    try {
      const response = await apiClient.communication.getTemplates();
      return { success: true, data: response.data };
    } catch (error) {
      return this.handleError(error);
    }
  },

  async createTemplate(data) {
    try {
      const response = await apiClient.communication.createTemplate(data);
      return { success: true, data: response.data };
    } catch (error) {
      return this.handleError(error);
    }
  },

  async updateTemplate(templateId, data) {
    try {
      const response = await apiClient.communication.updateTemplate(templateId, data);
      return { success: true, data: response.data };
    } catch (error) {
      return this.handleError(error);
    }
  },

  async deleteTemplate(templateId) {
    try {
      await apiClient.communication.deleteTemplate(templateId);
      return { success: true, message: 'Deleted successfully' };
    } catch (error) {
      return this.handleError(error);
    }
  },

  async sendBulkMessage(entityId, data) {
    try {
      const response = await apiClient.communication.sendBulk(entityId, data);
      return { success: true, data: response.data };
    } catch (error) {
      return this.handleError(error);
    }
  },

  async previewBulkMessage(entityId, data) {
    try {
      const response = await apiClient.communication.previewBulk(entityId, data);
      return { success: true, data: response.data };
    } catch (error) {
      return this.handleError(error);
    }
  },

  async getMessageHistory(entityId, query) {
    try {
      const response = await apiClient.communication.getHistory(entityId, query);
      return { success: true, data: response.data };
    } catch (error) {
      return this.handleError(error);
    }
  },

  async getMessageDetails(entityId, messageId) {
    try {
      const response = await apiClient.communication.getMessageDetails(entityId, messageId);
      return { success: true, data: response.data };
    } catch (error) {
      return this.handleError(error);
    }
  },

  async getDevoteesForMessaging(entityId, query) {
    try {
      const response = await apiClient.communication.getDevotees(entityId, query);
      return { success: true, data: response.data };
    } catch (error) {
      return this.handleError(error);
    }
  },

  async sendDirectNotification(payload) {
  try {
    const response = await apiClient.communication.sendDirectNotification(payload);
    return { success: true, data: response.data };
  } catch (error) {
    console.error('Error in sendDirectNotification:', error);
    return { 
      success: false, 
      error: error.response?.data?.error || error.response?.data?.message || error.message || 'Server error' 
    };
  }
},

  async getUnreadNotifications(entityId, userId) {
    try {
      const response = await apiClient.communication.getUnreadNotifications(entityId, userId);
      return { success: true, data: response.data };
    } catch (error) {
      return this.handleError(error);
    }
  },

  async markNotificationAsRead(entityId, userId, notificationId) {
    try {
      await apiClient.communication.markNotificationAsRead(entityId, userId, notificationId);
      return { success: true, message: 'Marked as read' };
    } catch (error) {
      return this.handleError(error);
    }
  },

  async getCommunicationStats(entityId, query) {
    try {
      const response = await apiClient.communication.getStats(entityId, query);
      return { success: true, data: response.data };
    } catch (error) {
      return this.handleError(error);
    }
  },

  async sendDirectNotification(payload) {
    try {
      const response = await apiClient.communication.sendDirectNotification(payload);
      return { success: true, data: response.data };
    } catch (error) {
      console.error('Error in sendDirectNotification:', error);
      return this.handleError(error);
    }
  },

  validateMessage(messageData) {
    const errors = [];
    if (!messageData.type) errors.push('Message type is required');
    if (!messageData.content && !messageData.templateId) errors.push('Content or template required');
    if (messageData.type === 'email' && !messageData.subject) errors.push('Email subject required');
    if (!messageData.recipients) errors.push('Recipients required');
    return { isValid: errors.length === 0, errors };
  },

  renderTemplate(content, variables = {}) {
    let output = content;
    const merged = {
      temple_name: '[Temple Name]',
      devotee_name: '[Devotee Name]',
      current_date: new Date().toLocaleDateString('en-IN'),
      current_time: new Date().toLocaleTimeString('en-IN'),
      ...variables
    };
    Object.entries(merged).forEach(([key, value]) => {
      output = output.replace(new RegExp(`\\{${key}\\}`, 'gi'), value);
    });
    return output;
  },

  handleError(error) {
    console.error('API Error:', error);
    
    // Extract error message
    let errorMessage = 'An unexpected error occurred';
    
    if (error.response) {
      // The request was made and the server responded with a status code
      // that falls out of the range of 2xx
      if (error.response.data && error.response.data.error) {
        errorMessage = error.response.data.error;
      } else if (error.response.data && error.response.data.message) {
        errorMessage = error.response.data.message;
      } else if (error.response.status === 500) {
        errorMessage = 'Server error: Could not process your request';
      } else if (error.response.status === 404) {
        errorMessage = 'Service not found';
      } else if (error.response.status === 401) {
        errorMessage = 'Authentication error';
      }
    } else if (error.request) {
      // The request was made but no response was received
      errorMessage = 'No response from server. Check your connection.';
    } else if (error.message) {
      // Something happened in setting up the request that triggered an Error
      errorMessage = error.message;
    }
    
    return {
      success: false,
      error: errorMessage,
      status: error.response?.status || 500
    };
  }
};

export default communicationService;