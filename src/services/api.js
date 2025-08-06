import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'

// API Configuration - FIXED to use proper URL structure
const API_CONFIG = {
  BASE_URL: import.meta.env.DEV 
    ? ''  // Relative for Vite dev proxy
    : (import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'),
  TIMEOUT: 30000,
  RETRY_ATTEMPTS: 3,
  RETRY_DELAY: 1000
}

console.log('🔧 API Config:', API_CONFIG)

// Create axios instance
const axiosInstance = axios.create({
  baseURL: API_CONFIG.BASE_URL,
  timeout: API_CONFIG.TIMEOUT,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
})

// Request interceptor - Add auth token and tenant context
axiosInstance.interceptors.request.use(
  (config) => {
    console.log(`🚀 API Request: ${config.method?.toUpperCase()} ${config.url}`)

    let token = localStorage.getItem('access_token')

    if (!token) {
      try {
        const authStore = useAuthStore()
        token = authStore.token
      } catch (e) {
        console.warn('Could not access auth store:', e)
      }
    }

    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    } else {
      console.warn('No authentication token for:', config.url)
    }

    try {
      const authStore = useAuthStore()
      if (authStore.user?.entityId && !config.url.includes('/superadmin/')) {
        config.headers['X-Tenant-ID'] = authStore.user.entityId
      }
    } catch (e) {
      console.warn('Could not set X-Tenant-ID:', e)
    }

    config.metadata = { startTime: new Date() }
    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor
axiosInstance.interceptors.response.use(
  (response) => {
    const duration = new Date() - response.config.metadata.startTime
    console.log(`✅ API Response: ${response.config.method?.toUpperCase()} ${response.config.url} - ${duration}ms`)
    return response
  },
  async (error) => {
    console.error('❌ API Error:', error)

    let showToast
    try {
      const toastService = useToast()
      showToast = toastService.showToast || toastService.error || toastService
    } catch {
      showToast = () => {}
    }

    let authStore
    try {
      authStore = useAuthStore()
    } catch {}

    if (error.response) {
      const { status, data } = error.response

      switch (status) {
        case 401:
          showToast({ type: 'error', title: 'Authentication Required', message: 'Your session has expired. Please login again.' })
          if (authStore) authStore.logout()
          else localStorage.removeItem('auth_token')
          window.location.href = '/login'
          break

        case 403:
          showToast({ type: 'warning', title: 'Access Denied', message: data.message || 'You do not have permission.' })
          break

        case 404:
          showToast({ type: 'error', title: 'Not Found', message: data?.message || 'API endpoint not found.' })
          console.warn('404 URL:', error.config.url)
          break

        case 422:
          console.warn('Validation errors:', data.errors)
          break

        case 429:
          showToast({ type: 'warning', title: 'Too Many Requests', message: 'Please wait before retrying.' })
          break

        case 500:
          showToast({ type: 'error', title: 'Server Error', message: 'Internal error, please try again later.' })
          break

        default:
          showToast({ type: 'error', title: 'Request Failed', message: data?.message || 'Unexpected error occurred.' })
      }
    } else if (error.request) {
      showToast({ type: 'error', title: 'Network Error', message: 'Check your internet connection.' })
    } else {
      showToast({ type: 'error', title: 'Request Error', message: 'Request setup failed.' })
    }

    return Promise.reject(error)
  }
)

// Retry logic
const retryRequest = async (config, attempt = 1) => {
  try {
    return await axiosInstance(config)
  } catch (error) {
    if (attempt < API_CONFIG.RETRY_ATTEMPTS && error.response?.status >= 500) {
      console.log(`Retrying request (attempt ${attempt + 1})`)
      await new Promise(resolve => setTimeout(resolve, API_CONFIG.RETRY_DELAY * attempt))
      return retryRequest(config, attempt + 1)
    }
    throw error
  }
}

// Core API
export const api = {
  async get(url, config = {}) {
    try {
      const response = await retryRequest({ method: 'get', url, ...config })
      return response.data
    } catch (error) {
      throw this.handleError(error)
    }
  },

  async post(url, data = {}, config = {}) {
    try {
      const response = await axiosInstance.post(url, data, config)
      return response.data
    } catch (error) {
      throw this.handleError(error)
    }
  },

  async put(url, data = {}, config = {}) {
    try {
      const response = await axiosInstance.put(url, data, config)
      return response.data
    } catch (error) {
      throw this.handleError(error)
    }
  },

  async patch(url, data = {}, config = {}) {
    try {
      const response = await axiosInstance.patch(url, data, config)
      return response.data
    } catch (error) {
      throw this.handleError(error)
    }
  },

  async delete(url, config = {}) {
    try {
      const response = await axiosInstance.delete(url, config)
      return response.data
    } catch (error) {
      throw this.handleError(error)
    }
  },

  async upload(url, formData, onProgress = null) {
    try {
      const config = {
        headers: { 'Content-Type': 'multipart/form-data' }
      }
      if (onProgress) {
        config.onUploadProgress = (e) => {
          const percent = Math.round((e.loaded * 100) / e.total)
          onProgress(percent)
        }
      }
      const response = await axiosInstance.post(url, formData, config)
      return response.data
    } catch (error) {
      throw this.handleError(error)
    }
  },

  async download(url, filename = null) {
    try {
      const response = await axiosInstance.get(url, { responseType: 'blob' })
      const blob = new Blob([response.data])
      const downloadUrl = URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = downloadUrl
      link.download = filename || 'download'
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      URL.revokeObjectURL(downloadUrl)
      return response.data
    } catch (error) {
      throw this.handleError(error)
    }
  },

  async batch(requests) {
    try {
      const promises = requests.map(({ method, url, data, config }) =>
        this[method](url, data, config)
      )
      return await Promise.allSettled(promises)
    } catch (error) {
      throw this.handleError(error)
    }
  },

  handleError(error) {
    const err = {
      message: 'Unexpected error occurred',
      status: null,
      data: null,
      errors: null
    }
    if (error.response) {
      err.message = error.response.data?.message || error.message
      err.status = error.response.status
      err.data = error.response.data
      err.errors = error.response.data?.errors
    } else if (error.request) {
      err.message = 'Network error - check your internet connection'
    } else {
      err.message = error.message
    }
    return err
  }
}

// API Endpoints
export const endpoints = {
  auth: {
    login: '/api/v1/auth/login',
    register: '/api/v1/auth/register',
    logout: '/api/v1/auth/logout',
    refresh: '/api/v1/auth/refresh',
    forgotPassword: '/api/v1/auth/forgot-password',
    resetPassword: '/api/v1/auth/reset-password',
    verifyEmail: '/api/v1/auth/verify-email'
  },

  superAdmin: {
    dashboard: '/api/v1/superadmin/dashboard',
    tenants: '/api/v1/superadmin/tenants',
    pendingTenants: '/api/v1/superadmin/tenants/pending',
    updateTenantApproval: (id) => `/api/v1/superadmin/tenants/${id}/approval`,
    pendingEntities: '/api/v1/superadmin/entities/pending',
    updateEntityApproval: (id) => `/api/v1/superadmin/entities/${id}/approval`,
    stats: '/api/v1/superadmin/stats'
  },

  tenant: {
    dashboard: '/api/v1/tenant/dashboard',
    profile: '/api/v1/tenant/profile',
    temples: '/api/v1/tenant/temples',
    createTemple: '/api/v1/tenant/temples',
    updateTemple: (id) => `/api/v1/tenant/temples/${id}`,
    deleteTemple: (id) => `/api/v1/tenant/temples/${id}`,
    documents: '/api/v1/tenant/documents',
    uploadDocument: '/api/v1/tenant/documents/upload'
  },

  entity: (entityId) => ({
    dashboard: `/api/v1/entities/${entityId}/dashboard`,
    devotees: `/api/v1/entities/${entityId}/devotees`,
    devotee: (devoteeId) => `/api/v1/entities/${entityId}/devotees/${devoteeId}`,
    sevas: `/api/v1/entities/${entityId}/sevas`,
    seva: (sevaId) => `/api/v1/entities/${entityId}/sevas/${sevaId}`,
    donations: `/api/v1/entities/${entityId}/donations`,
    donation: (donationId) => `/api/v1/entity/${entityId}/donations/${donationId}`,
    events: `/api/v1/entities/${entityId}/events`,
    event: (eventId) => `/api/v1/entities/${entityId}/events/${eventId}`,
    volunteers: `/api/v1/entities/${entityId}/volunteers`,
    volunteer: (volunteerId) => `/api/v1/entities/${entityId}/volunteers/${volunteerId}`,
    communication: `/api/v1/entities/${entityId}/communication`,
    templates: `/api/v1/entities/${entityId}/communication/templates`,
    bulkMessage: `/api/v1/entities/${entityId}/communication/bulk`,
    reports: `/api/v1/entities/${entityId}/reports`,
    settings: `/api/v1/entities/${entityId}/settings`
  }),

  donations: {
    my: '/api/v1/donations/my',
  },

  // ✅ Updated: Use Devotee-accessible endpoint for temple listings
  public: {
    temples: '/api/v1/devotee/available-temples',
    temple: (id) => `/api/v1/devotee/available-temples/${id}`,
    searchTemples: '/api/v1/devotee/available-temples/search',
    joinTemple: '/api/v1/devotee/available-temples/join'
  },

  user: {
    profile: '/api/v1/user/profile',
    preferences: '/api/v1/user/preferences',
    notifications: '/api/v1/user/notifications',
    markNotificationRead: (id) => `/api/v1/user/notifications/${id}/read`,
    avatar: '/api/v1/user/avatar'
  }
}

// API Utils
export const apiUtils = {
  buildQuery(params) {
    const query = new URLSearchParams()
    Object.entries(params).forEach(([key, value]) => {
      if (value !== null && value !== undefined && value !== '') {
        query.append(key, value)
      }
    })
    return query.toString()
  },

  async paginate(endpoint, params = {}) {
    const query = this.buildQuery(params)
    const url = query ? `${endpoint}?${query}` : endpoint
    return api.get(url)
  },

  async search(endpoint, searchTerm, filters = {}) {
    const params = { search: searchTerm, ...filters }
    return this.paginate(endpoint, params)
  },

  async bulkOperation(endpoint, operation, ids) {
    return api.post(`${endpoint}/bulk`, { operation, ids })
  }
}

export default api
