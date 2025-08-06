// src/plugins/axios.js
import axios from 'axios'

// Base URL Configuration
const baseURL = import.meta.env.DEV ? '/api' : import.meta.env.VITE_API_BASE_URL;

// Create axios instance with base configuration
const api = axios.create({
  baseURL,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
})

// Log configuration in development only
if (import.meta.env.DEV) {
  console.log('API Configuration:', {
    baseURL,
    environment: import.meta.env.MODE
  })
}

// Request interceptor - Add auth token and tenant ID
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

    const isAuthEndpoint = config.url.includes('/auth/') || config.url.includes('/v1/auth/')
    
    if (!isAuthEndpoint) {
      const tenantId = localStorage.getItem('current_tenant_id')
      if (tenantId) {
        config.headers['X-Tenant-ID'] = tenantId
        if (import.meta.env.DEV) {
          console.log(`Request with Tenant ID: ${tenantId}`)
        }
      }

      const entityId = localStorage.getItem('current_entity_id')
      if (tenantId) {
  config.headers['X-Tenant-ID'] = tenantId
  if (import.meta.env.DEV) {
    console.log(`Request with Tenant ID: ${tenantId}`)
  }
}
    }
const entityId = localStorage.getItem('current_entity_id')
if (entityId) {
  config.headers['X-Entity-ID'] = entityId
  if (import.meta.env.DEV) {
    console.log(`Request with Entity ID: ${entityId}`)
  }
}

    if (import.meta.env.DEV) {
      console.log(`API Request: ${config.method.toUpperCase()} ${config.url}`)
    }

    return config
  },
  (error) => {
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  (response) => {
    // Log success in development
    if (import.meta.env.DEV) {
      console.log(`API Response: ${response.status} ${response.config.method.toUpperCase()} ${response.config.url}`)
    }
    
    return response
  },
  (error) => {
    // Log error
    console.error('API Error:', error.message)
    
    if (error.response) {
      // Log response error details
      console.error(`Status: ${error.response.status}`, 
                   `Method: ${error.config?.method?.toUpperCase()}`, 
                   `URL: ${error.config?.url}`)
      
      // Handle 401 errors (except during login attempts)
      if (error.response.status === 401) {
        const isLoginAttempt = error.config?.url?.includes('/auth/login') || 
                             error.config?.url?.includes('/v1/auth/login')
        
        if (!isLoginAttempt) {
          // Clear authentication data
          localStorage.removeItem('auth_token')
          localStorage.removeItem('user_data')
          localStorage.removeItem('current_tenant_id')
          delete api.defaults.headers.common['Authorization']
          delete api.defaults.headers.common['X-Tenant-ID']
          
          // Redirect to login if not already there
          if (!window.location.pathname.includes('/login')) {
            window.location.href = '/login?error=session_expired'
          }
        }
      }
    }
    
    return Promise.reject(error)
  }
)

// Auth-specific API endpoints
const authAPI = {
  login: (credentials) => api.post('/v1/auth/login', credentials),
  register: (userData) => api.post('/v1/auth/register', userData),
  logout: () => api.post('/v1/auth/logout'),
  refreshToken: () => api.post('/v1/auth/refresh'),
  getProfile: () => api.get('/v1/auth/me')
}

// Temple/Entity endpoints
const templeAPI = {
  create: (data) => api.post('/v1/tenant/temples', data),
  getAll: () => api.get('/v1/tenant/temples'),
  getById: (id) => api.get(`/v1/tenant/temples/${id}`),
  update: (id, data) => api.put(`/v1/tenant/temples/${id}`, data),
  delete: (id) => api.delete(`/v1/tenant/temples/${id}`)
}

// User role-specific endpoints
const userAPI = {
  // Devotee endpoints
  devotee: {
    getProfile: () => api.get('/v1/devotee/profile'),
    updateProfile: (data) => api.put('/v1/devotee/profile', data),
    joinTemple: (templeId) => api.post('/v1/devotee/join-temple', { templeId })
  },
  
  // Volunteer endpoints
  volunteer: {
    getProfile: () => api.get('/v1/volunteer/profile'),
    updateProfile: (data) => api.put('/v1/volunteer/profile', data),
    joinTemple: (templeId) => api.post('/v1/volunteer/join-temple', { templeId })
  }
}

// Admin endpoints
const adminAPI = {
  getTemples: () => api.get('/v1/superadmin/temples'),
  getPendingTemples: () => api.get('/v1/superadmin/temples/pending'),
  approveTemple: (id, notes) => api.post(`/v1/superadmin/temples/${id}/approve`, { notes }),
  rejectTemple: (id, reason) => api.post(`/v1/superadmin/temples/${id}/reject`, { reason })
}

// Event endpoints
const eventAPI = {
  getAll: () => api.get('/v1/events'),
  getUpcoming: () => api.get('/v1/events/upcoming'),
  getById: (id) => api.get(`/v1/events/${id}`),
  getStats: () => api.get('/v1/events/stats'),
  create: (eventData) => api.post('/v1/events', eventData),
  update: (id, eventData) => api.put(`/v1/events/${id}`, eventData),
  delete: (id) => api.delete(`/v1/events/${id}`),
  getRSVPs: (eventId) => api.get(`/v1/event-rsvps/${eventId}`),
  createRSVP: (eventId) => api.post(`/v1/event-rsvps/${eventId}`, {})
}

// Entity Dashboard endpoints
const dashboardAPI = {
  getEntityDashboard: (entityId) => api.get(`/v1/entities/${entityId}/dashboard`),
  getEntityEvents: (entityId, limit = 3) => api.get(`/v1/events/upcoming?entity_id=${entityId}&limit=${limit}`),
  getEntityDonors: (entityId, limit = 5) => api.get(`/v1/donations/top?entity_id=${entityId}&limit=${limit}`),
  getEntityNotifications: (entityId, limit = 5) => api.get(`/v1/notifications?entity_id=${entityId}&limit=${limit}`)
}

// Notification endpoints
const communicationAPI = {
  getTemplates: () => api.get(`/v1/notifications/templates`),
  createTemplate: (data) => api.post(`/v1/notifications/templates`, data),
  updateTemplate: (id, data) => api.put(`/v1/notifications/templates/${id}`, data),
  deleteTemplate: (id) => api.delete(`/v1/notifications/templates/${id}`),

  sendBulk: (entityId, data) => api.post(`/v1/entities/${entityId}/communication/bulk-send`, data),
  previewBulk: (entityId, data) => api.post(`/entities/${entityId}/communication/preview`, data),
  getHistory: (entityId, query) => api.get(`/entities/${entityId}/communication/history?${query}`),
  getMessageDetails: (entityId, messageId) => api.get(`/entities/${entityId}/communication/messages/${messageId}`),

  getDevotees: (entityId, query) => api.get(`/entities/${entityId}/devotees/for-messaging?${query}`),

  sendNotification: (entityId, data) => api.post(`/entities/${entityId}/communication/notifications`, data),
  getUnreadNotifications: (entityId, userId) => api.get(`/entities/${entityId}/users/${userId}/notifications/unread`),
  markNotificationAsRead: (entityId, userId, notificationId) => api.put(`/entities/${entityId}/users/${userId}/notifications/${notificationId}/read`),

  getStats: (entityId, query) => api.get(`/entities/${entityId}/communication/stats?${query}`),

  sendDirectNotification: (payload) => api.post(`/v1/notifications/send`, payload)
};

// Export structured API client
export const apiClient = {
  auth: authAPI,
  temple: templeAPI,
  user: userAPI,
  admin: adminAPI,
  event: eventAPI,
  dashboard: dashboardAPI,
  communication: communicationAPI 
}

// Default export for backward compatibility
export default api