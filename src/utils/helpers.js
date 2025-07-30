/**
 * Temple SaaS Platform - Helper Utilities
 * Comprehensive utility functions for multi-tenant temple management system
 */

// ============================================================================
// ROLE & PERMISSION HELPERS
// ============================================================================

/**
 * User roles available in the system
 */
export const USER_ROLES = {
  SUPER_ADMIN: 'superadmin',
  TENANT: 'tenant',
  DEVOTEE: 'devotee',
  VOLUNTEER: 'volunteer'
}

/**
 * Check if user has specific role
 * @param {Object} user - User object
 * @param {string} role - Role to check
 * @returns {boolean}
 */
export const hasRole = (user, role) => {
  return user?.role === role
}

/**
 * Check if user has any of the specified roles
 * @param {Object} user - User object
 * @param {Array} roles - Array of roles to check
 * @returns {boolean}
 */
export const hasAnyRole = (user, roles) => {
  return roles.includes(user?.role)
}

/**
 * Get user dashboard route based on role
 * @param {Object} user - User object
 * @param {string} entityId - Entity ID (for entity-specific roles)
 * @returns {string}
 */
export const getUserDashboardRoute = (user, entityId = null) => {
  const roleRoutes = {
    [USER_ROLES.SUPER_ADMIN]: '/superadmin/dashboard',
    [USER_ROLES.TENANT]: '/tenant/dashboard',
    [USER_ROLES.DEVOTEE]: entityId ? `/entity/${entityId}/devotee/dashboard` : '/temple-selection',
    [USER_ROLES.VOLUNTEER]: entityId ? `/entity/${entityId}/volunteer/dashboard` : '/temple-selection'
  }
  
  return roleRoutes[user?.role] || '/login'
}

/**
 * Check if user can access entity-specific features
 * @param {Object} user - User object
 * @returns {boolean}
 */
export const canAccessEntity = (user) => {
  return hasAnyRole(user, [USER_ROLES.DEVOTEE, USER_ROLES.VOLUNTEER])
}

// ============================================================================
// STATUS & APPROVAL HELPERS
// ============================================================================

/**
 * Status types for entities/temples
 */
export const STATUS_TYPES = {
  PENDING: 'pending',
  APPROVED: 'approved',
  REJECTED: 'rejected',
  ACTIVE: 'active',
  INACTIVE: 'inactive'
}

/**
 * Get status badge classes based on status
 * @param {string} status - Status value
 * @returns {Object}
 */
export const getStatusBadgeClasses = (status) => {
  const statusClasses = {
    [STATUS_TYPES.PENDING]: {
      bg: 'bg-yellow-100',
      text: 'text-yellow-800',
      border: 'border-yellow-200'
    },
    [STATUS_TYPES.APPROVED]: {
      bg: 'bg-green-100',
      text: 'text-green-800',
      border: 'border-green-200'
    },
    [STATUS_TYPES.REJECTED]: {
      bg: 'bg-red-100',
      text: 'text-red-800',
      border: 'border-red-200'
    },
    [STATUS_TYPES.ACTIVE]: {
      bg: 'bg-indigo-100',
      text: 'text-indigo-800',
      border: 'border-indigo-200'
    },
    [STATUS_TYPES.INACTIVE]: {
      bg: 'bg-gray-100',
      text: 'text-gray-800',
      border: 'border-gray-200'
    }
  }
  
  return statusClasses[status] || statusClasses[STATUS_TYPES.PENDING]
}

/**
 * Get human-readable status text
 * @param {string} status - Status value
 * @returns {string}
 */
export const getStatusText = (status) => {
  const statusTexts = {
    [STATUS_TYPES.PENDING]: 'Pending Approval',
    [STATUS_TYPES.APPROVED]: 'Approved',
    [STATUS_TYPES.REJECTED]: 'Rejected',
    [STATUS_TYPES.ACTIVE]: 'Active',
    [STATUS_TYPES.INACTIVE]: 'Inactive'
  }
  
  return statusTexts[status] || 'Unknown'
}

// ============================================================================
// DATA FORMATTING HELPERS
// ============================================================================

/**
 * Format currency amount
 * @param {number} amount - Amount to format
 * @param {string} currency - Currency code (default: INR)
 * @returns {string}
 */
export const formatCurrency = (amount, currency = 'INR') => {
  if (!amount || isNaN(amount)) return '₹0'
  
  const formatter = new Intl.NumberFormat('en-IN', {
    style: 'currency',
    currency: currency,
    minimumFractionDigits: 0,
    maximumFractionDigits: 2
  })
  
  return formatter.format(amount)
}

/**
 * Format large numbers with Indian numbering system
 * @param {number} num - Number to format
 * @returns {string}
 */
export const formatNumber = (num) => {
  if (!num || isNaN(num)) return '0'
  
  return new Intl.NumberFormat('en-IN').format(num)
}

/**
 * Format date in Indian format
 * @param {string|Date} date - Date to format
 * @param {Object} options - Formatting options
 * @returns {string}
 */
export const formatDate = (date, options = {}) => {
  if (!date) return ''
  
  const defaultOptions = {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    timeZone: 'Asia/Kolkata'
  }
  
  const formatOptions = { ...defaultOptions, ...options }
  
  try {
    return new Date(date).toLocaleDateString('en-IN', formatOptions)
  } catch (error) {
    console.error('Invalid date format:', date)
    return ''
  }
}

/**
 * Format date and time
 * @param {string|Date} date - Date to format
 * @returns {string}
 */
export const formatDateTime = (date) => {
  return formatDate(date, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

/**
 * Get relative time (e.g., "2 hours ago")
 * @param {string|Date} date - Date to compare
 * @returns {string}
 */
export const getRelativeTime = (date) => {
  if (!date) return ''
  
  const now = new Date()
  const past = new Date(date)
  const diffInSeconds = Math.floor((now - past) / 1000)
  
  if (diffInSeconds < 60) return 'Just now'
  if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)} minutes ago`
  if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)} hours ago`
  if (diffInSeconds < 2592000) return `${Math.floor(diffInSeconds / 86400)} days ago`
  
  return formatDate(date)
}

// ============================================================================
// STRING MANIPULATION HELPERS
// ============================================================================

/**
 * Capitalize first letter of each word
 * @param {string} str - String to capitalize
 * @returns {string}
 */
export const capitalizeWords = (str) => {
  if (!str) return ''
  return str.replace(/\w\S*/g, (txt) => 
    txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase()
  )
}

/**
 * Truncate text with ellipsis
 * @param {string} text - Text to truncate
 * @param {number} maxLength - Maximum length
 * @returns {string}
 */
export const truncateText = (text, maxLength = 100) => {
  if (!text || text.length <= maxLength) return text
  return text.slice(0, maxLength).trim() + '...'
}

/**
 * Generate initials from name
 * @param {string} name - Full name
 * @returns {string}
 */
export const getInitials = (name) => {
  if (!name) return ''
  return name
    .split(' ')
    .map(word => word.charAt(0).toUpperCase())
    .slice(0, 2)
    .join('')
}

/**
 * Convert string to slug
 * @param {string} str - String to convert
 * @returns {string}
 */
export const slugify = (str) => {
  if (!str) return ''
  return str
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, '')
    .replace(/[\s_-]+/g, '-')
    .replace(/^-+|-+$/g, '')
}

// ============================================================================
// VALIDATION HELPERS
// ============================================================================

/**
 * Validate email format
 * @param {string} email - Email to validate
 * @returns {boolean}
 */
export const isValidEmail = (email) => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

/**
 * Validate Indian mobile number
 * @param {string} mobile - Mobile number to validate
 * @returns {boolean}
 */
export const isValidMobile = (mobile) => {
  const mobileRegex = /^[6-9]\d{9}$/
  return mobileRegex.test(mobile.replace(/\D/g, ''))
}

/**
 * Validate password strength
 * @param {string} password - Password to validate
 * @returns {Object}
 */
export const validatePasswordStrength = (password) => {
  const criteria = {
    minLength: password.length >= 8,
    hasUpperCase: /[A-Z]/.test(password),
    hasLowerCase: /[a-z]/.test(password),
    hasNumbers: /\d/.test(password),
    hasSpecialChar: /[!@#$%^&*(),.?":{}|<>]/.test(password)
  }
  
  const score = Object.values(criteria).filter(Boolean).length
  
  return {
    ...criteria,
    score,
    strength: score < 3 ? 'weak' : score < 5 ? 'medium' : 'strong',
    isValid: score >= 3
  }
}

// ============================================================================
// ARRAY & OBJECT UTILITIES
// ============================================================================

/**
 * Deep clone an object
 * @param {any} obj - Object to clone
 * @returns {any}
 */
export const deepClone = (obj) => {
  if (obj === null || typeof obj !== 'object') return obj
  if (obj instanceof Date) return new Date(obj.getTime())
  if (obj instanceof Array) return obj.map(item => deepClone(item))
  if (typeof obj === 'object') {
    const clonedObj = {}
    for (const key in obj) {
      if (obj.hasOwnProperty(key)) {
        clonedObj[key] = deepClone(obj[key])
      }
    }
    return clonedObj
  }
}

/**
 * Group array by key
 * @param {Array} array - Array to group
 * @param {string} key - Key to group by
 * @returns {Object}
 */
export const groupBy = (array, key) => {
  return array.reduce((result, item) => {
    const group = item[key]
    if (!result[group]) result[group] = []
    result[group].push(item)
    return result
  }, {})
}

/**
 * Sort array by multiple fields
 * @param {Array} array - Array to sort
 * @param {Array} fields - Fields to sort by (e.g., ['name', '-date'])
 * @returns {Array}
 */
export const sortBy = (array, fields) => {
  return array.sort((a, b) => {
    for (const field of fields) {
      const isDesc = field.startsWith('-')
      const key = isDesc ? field.slice(1) : field
      const direction = isDesc ? -1 : 1
      
      if (a[key] < b[key]) return -1 * direction
      if (a[key] > b[key]) return 1 * direction
    }
    return 0
  })
}

/**
 * Remove duplicates from array
 * @param {Array} array - Array to deduplicate
 * @param {string} key - Key to check for uniqueness (optional)
 * @returns {Array}
 */
export const uniqueBy = (array, key = null) => {
  if (!key) return [...new Set(array)]
  
  const seen = new Set()
  return array.filter(item => {
    const value = item[key]
    if (seen.has(value)) return false
    seen.add(value)
    return true
  })
}

// ============================================================================
// FILE & UPLOAD HELPERS
// ============================================================================

/**
 * Format file size
 * @param {number} bytes - File size in bytes
 * @returns {string}
 */
export const formatFileSize = (bytes) => {
  if (!bytes) return '0 Bytes'
  
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * Get file extension
 * @param {string} filename - Filename
 * @returns {string}
 */
export const getFileExtension = (filename) => {
  return filename.slice((filename.lastIndexOf('.') - 1 >>> 0) + 2)
}

/**
 * Check if file is image
 * @param {string} filename - Filename
 * @returns {boolean}
 */
export const isImageFile = (filename) => {
  const imageExtensions = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp', 'svg']
  const extension = getFileExtension(filename).toLowerCase()
  return imageExtensions.includes(extension)
}

// ============================================================================
// URL & QUERY HELPERS
// ============================================================================

/**
 * Build query string from object
 * @param {Object} params - Query parameters
 * @returns {string}
 */
export const buildQueryString = (params) => {
  const filtered = Object.entries(params)
    .filter(([_, value]) => value !== null && value !== undefined && value !== '')
    .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(value)}`)
  
  return filtered.length > 0 ? `?${filtered.join('&')}` : ''
}

/**
 * Parse query string to object
 * @param {string} queryString - Query string
 * @returns {Object}
 */
export const parseQueryString = (queryString) => {
  const params = new URLSearchParams(queryString)
  const result = {}
  
  for (const [key, value] of params.entries()) {
    result[key] = value
  }
  
  return result
}

// ============================================================================
// LOCAL STORAGE HELPERS (For development - replace with proper state management)
// ============================================================================

/**
 * Set item in localStorage with error handling
 * @param {string} key - Storage key
 * @param {any} value - Value to store
 * @returns {boolean}
 */
export const setStorageItem = (key, value) => {
  try {
    localStorage.setItem(key, JSON.stringify(value))
    return true
  } catch (error) {
    console.error('Failed to set localStorage item:', error)
    return false
  }
}

/**
 * Get item from localStorage with error handling
 * @param {string} key - Storage key
 * @param {any} defaultValue - Default value if not found
 * @returns {any}
 */
export const getStorageItem = (key, defaultValue = null) => {
  try {
    const item = localStorage.getItem(key)
    return item ? JSON.parse(item) : defaultValue
  } catch (error) {
    console.error('Failed to get localStorage item:', error)
    return defaultValue
  }
}

/**
 * Remove item from localStorage
 * @param {string} key - Storage key
 * @returns {boolean}
 */
export const removeStorageItem = (key) => {
  try {
    localStorage.removeItem(key)
    return true
  } catch (error) {
    console.error('Failed to remove localStorage item:', error)
    return false
  }
}

// ============================================================================
// MISCELLANEOUS HELPERS
// ============================================================================

/**
 * Generate random ID
 * @param {number} length - Length of ID
 * @returns {string}
 */
export const generateId = (length = 8) => {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let result = ''
  for (let i = 0; i < length; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  return result
}

/**
 * Debounce function
 * @param {Function} func - Function to debounce
 * @param {number} wait - Wait time in milliseconds
 * @returns {Function}
 */
export const debounce = (func, wait) => {
  let timeout
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout)
      func(...args)
    }
    clearTimeout(timeout)
    timeout = setTimeout(later, wait)
  }
}

/**
 * Throttle function
 * @param {Function} func - Function to throttle
 * @param {number} limit - Time limit in milliseconds
 * @returns {Function}
 */
export const throttle = (func, limit) => {
  let inThrottle
  return function executedFunction(...args) {
    if (!inThrottle) {
      func.apply(this, args)
      inThrottle = true
      setTimeout(() => inThrottle = false, limit)
    }
  }
}

/**
 * Copy text to clipboard
 * @param {string} text - Text to copy
 * @returns {Promise<boolean>}
 */
export const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    return true
  } catch (error) {
    console.error('Failed to copy to clipboard:', error)
    return false
  }
}

/**
 * Download data as file
 * @param {string} data - Data to download
 * @param {string} filename - Filename
 * @param {string} type - MIME type
 */
export const downloadFile = (data, filename, type = 'text/plain') => {
  const blob = new Blob([data], { type })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

// ============================================================================
// NOTIFICATION HELPERS
// ============================================================================

/**
 * Notification types
 */
export const NOTIFICATION_TYPES = {
  SUCCESS: 'success',
  ERROR: 'error',
  WARNING: 'warning',
  INFO: 'info'
}

/**
 * Get notification classes
 * @param {string} type - Notification type
 * @returns {Object}
 */
export const getNotificationClasses = (type) => {
  const classes = {
    [NOTIFICATION_TYPES.SUCCESS]: {
      bg: 'bg-green-100',
      text: 'text-green-800',
      border: 'border-green-200',
      icon: 'text-green-600'
    },
    [NOTIFICATION_TYPES.ERROR]: {
      bg: 'bg-red-100',
      text: 'text-red-800',
      border: 'border-red-200',
      icon: 'text-red-600'
    },
    [NOTIFICATION_TYPES.WARNING]: {
      bg: 'bg-yellow-100',
      text: 'text-yellow-800',
      border: 'border-yellow-200',
      icon: 'text-yellow-600'
    },
    [NOTIFICATION_TYPES.INFO]: {
      bg: 'bg-indigo-100',
      text: 'text-indigo-800',
      border: 'border-indigo-200',
      icon: 'text-indigo-600'
    }
  }
  
  return classes[type] || classes[NOTIFICATION_TYPES.INFO]
}