/**
 * Formatters Utility
 * Centralized formatting functions for the Multi-Tenant Temple SaaS Platform
 * Handles currency, dates, numbers, text, and temple-specific data formatting
 */

// =============================================================================
// DATE FORMATTERS
// =============================================================================

/**
 * Format date to readable string
 * @param {Date|string} date - Date to format
 * @param {string} format - Format type ('short', 'long', 'time', 'datetime')
 * @returns {string} Formatted date string
 */
export const formatDate = (date, format = 'short') => {
  if (!date) return '-'
  
  const d = new Date(date)
  if (isNaN(d.getTime())) return 'Invalid Date'
  
  const options = {
    short: { 
      year: 'numeric', 
      month: 'short', 
      day: 'numeric' 
    },
    long: { 
      weekday: 'long', 
      year: 'numeric', 
      month: 'long', 
      day: 'numeric' 
    },
    time: { 
      hour: '2-digit', 
      minute: '2-digit', 
      hour12: true 
    },
    datetime: { 
      year: 'numeric', 
      month: 'short', 
      day: 'numeric',
      hour: '2-digit', 
      minute: '2-digit', 
      hour12: true 
    }
  }
  
  return d.toLocaleDateString('en-IN', options[format])
}

/**
 * Format date for temple events and festivals
 * @param {Date|string} date - Event date
 * @returns {string} Formatted event date
 */
export const formatEventDate = (date) => {
  if (!date) return '-'
  
  const d = new Date(date)
  if (isNaN(d.getTime())) return 'Invalid Date'
  
  const today = new Date()
  const tomorrow = new Date(today)
  tomorrow.setDate(tomorrow.getDate() + 1)
  
  if (d.toDateString() === today.toDateString()) {
    return `Today, ${formatDate(date, 'time')}`
  } else if (d.toDateString() === tomorrow.toDateString()) {
    return `Tomorrow, ${formatDate(date, 'time')}`
  } else {
    return formatDate(date, 'datetime')
  }
}

/**
 * Get relative time (e.g., "2 hours ago", "in 3 days")
 * @param {Date|string} date - Date to compare
 * @returns {string} Relative time string
 */
export const formatRelativeTime = (date) => {
  if (!date) return '-'
  
  const d = new Date(date)
  const now = new Date()
  const diffMs = d.getTime() - now.getTime()
  const diffMins = Math.round(diffMs / 60000)
  const diffHours = Math.round(diffMs / 3600000)
  const diffDays = Math.round(diffMs / 86400000)
  
  if (Math.abs(diffMins) < 60) {
    return diffMins > 0 ? `in ${diffMins} minutes` : `${Math.abs(diffMins)} minutes ago`
  } else if (Math.abs(diffHours) < 24) {
    return diffHours > 0 ? `in ${diffHours} hours` : `${Math.abs(diffHours)} hours ago`
  } else {
    return diffDays > 0 ? `in ${diffDays} days` : `${Math.abs(diffDays)} days ago`
  }
}

// =============================================================================
// CURRENCY & NUMBER FORMATTERS
// =============================================================================

/**
 * Format currency in Indian Rupees
 * @param {number} amount - Amount to format
 * @param {boolean} compact - Use compact notation (K, L, Cr)
 * @returns {string} Formatted currency string
 */
export const formatCurrency = (amount, compact = false) => {
  if (amount === null || amount === undefined || isNaN(amount)) return '₹0'
  
  const num = parseFloat(amount)
  
  if (compact && num >= 1000) {
    if (num >= 10000000) { // 1 Crore
      return `₹${(num / 10000000).toFixed(1)}Cr`
    } else if (num >= 100000) { // 1 Lakh
      return `₹${(num / 100000).toFixed(1)}L`
    } else if (num >= 1000) { // 1 Thousand
      return `₹${(num / 1000).toFixed(1)}K`
    }
  }
  
  return new Intl.NumberFormat('en-IN', {
    style: 'currency',
    currency: 'INR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 2
  }).format(num)
}

/**
 * Format numbers with Indian numbering system
 * @param {number} num - Number to format
 * @returns {string} Formatted number string
 */
export const formatNumber = (num) => {
  if (num === null || num === undefined || isNaN(num)) return '0'
  
  return new Intl.NumberFormat('en-IN').format(num)
}

/**
 * Format percentage
 * @param {number} value - Value to format as percentage
 * @param {number} decimals - Number of decimal places
 * @returns {string} Formatted percentage string
 */
export const formatPercentage = (value, decimals = 1) => {
  if (value === null || value === undefined || isNaN(value)) return '0%'
  
  return `${parseFloat(value).toFixed(decimals)}%`
}

// =============================================================================
// TEXT FORMATTERS
// =============================================================================

/**
 * Capitalize first letter of each word
 * @param {string} text - Text to capitalize
 * @returns {string} Capitalized text
 */
export const capitalizeWords = (text) => {
  if (!text) return ''
  
  return text.toLowerCase().replace(/\b\w/g, (char) => char.toUpperCase())
}

/**
 * Truncate text with ellipsis
 * @param {string} text - Text to truncate
 * @param {number} length - Maximum length
 * @returns {string} Truncated text
 */
export const truncateText = (text, length = 50) => {
  if (!text) return ''
  
  if (text.length <= length) return text
  
  return text.substring(0, length).trim() + '...'
}

/**
 * Format phone number for Indian mobile numbers
 * @param {string} phone - Phone number to format
 * @returns {string} Formatted phone number
 */
export const formatPhone = (phone) => {
  if (!phone) return ''
  
  // Remove all non-digits
  const cleaned = phone.replace(/\D/g, '')
  
  // Check for Indian mobile number (10 digits)
  if (cleaned.length === 10) {
    return `+91 ${cleaned.substring(0, 5)} ${cleaned.substring(5)}`
  }
  
  // Check for number with country code
  if (cleaned.length === 12 && cleaned.startsWith('91')) {
    const mobile = cleaned.substring(2)
    return `+91 ${mobile.substring(0, 5)} ${mobile.substring(5)}`
  }
  
  return phone // Return original if doesn't match expected format
}

/**
 * Format email to hide part of the address
 * @param {string} email - Email to format
 * @returns {string} Masked email
 */
export const maskEmail = (email) => {
  if (!email || !email.includes('@')) return email
  
  const [username, domain] = email.split('@')
  if (username.length <= 2) return email
  
  const masked = username.substring(0, 2) + '*'.repeat(username.length - 2)
  return `${masked}@${domain}`
}

// =============================================================================
// TEMPLE-SPECIFIC FORMATTERS
// =============================================================================

/**
 * Format seva name with proper capitalization
 * @param {string} sevaName - Seva name to format
 * @returns {string} Formatted seva name
 */
export const formatSevaName = (sevaName) => {
  if (!sevaName) return ''
  
  // Special handling for common Sanskrit/Telugu words
  const specialWords = {
    'archana': 'Archana',
    'abhishekam': 'Abhishekam',
    'kalyanotsavam': 'Kalyanotsavam',
    'annadanam': 'Annadanam',
    'darshan': 'Darshan',
    'pooja': 'Pooja',
    'aarti': 'Aarti',
    'prasadam': 'Prasadam'
  }
  
  return sevaName.toLowerCase().split(' ').map(word => {
    return specialWords[word] || capitalizeWords(word)
  }).join(' ')
}

/**
 * Format devotee status with color coding
 * @param {string} status - Devotee status
 * @returns {object} Status with display text and CSS classes
 */
export const formatDevoteeStatus = (status) => {
  const statusMap = {
    'active': {
      text: 'Active',
      badgeClass: 'bg-green-100 text-green-800',
      dotClass: 'bg-green-500'
    },
    'inactive': {
      text: 'Inactive',
      badgeClass: 'bg-gray-100 text-gray-800',
      dotClass: 'bg-gray-500'
    },
    'pending': {
      text: 'Pending Approval',
      badgeClass: 'bg-yellow-100 text-yellow-800',
      dotClass: 'bg-yellow-500'
    },
    'blocked': {
      text: 'Blocked',
      badgeClass: 'bg-red-100 text-red-800',
      dotClass: 'bg-red-500'
    }
  }
  
  return statusMap[status?.toLowerCase()] || statusMap['inactive']
}

/**
 * Format temple approval status
 * @param {string} status - Temple approval status
 * @returns {object} Status with display text and CSS classes
 */
export const formatTempleStatus = (status) => {
  const statusMap = {
    'pending': {
      text: 'Pending Review',
      badgeClass: 'bg-yellow-100 text-yellow-800 border-yellow-200',
      iconClass: 'text-yellow-600',
      dotClass: 'bg-yellow-500'
    },
    'approved': {
      text: 'Approved',
      badgeClass: 'bg-green-100 text-green-800 border-green-200',
      iconClass: 'text-green-600',
      dotClass: 'bg-green-500'
    },
    'rejected': {
      text: 'Rejected',
      badgeClass: 'bg-red-100 text-red-800 border-red-200',
      iconClass: 'text-red-600',
      dotClass: 'bg-red-500'
    },
    'draft': {
      text: 'Draft',
      badgeClass: 'bg-gray-100 text-gray-800 border-gray-200',
      iconClass: 'text-gray-600',
      dotClass: 'bg-gray-500'
    }
  }
  
  return statusMap[status?.toLowerCase()] || statusMap['draft']
}

/**
 * Format donation type with appropriate styling
 * @param {string} type - Donation type
 * @returns {object} Type with display text and CSS classes
 */
export const formatDonationType = (type) => {
  const typeMap = {
    'general': {
      text: 'General Donation',
      badgeClass: 'bg-indigo-100 text-indigo-800',
      iconClass: 'text-indigo-600'
    },
    'seva': {
      text: 'Seva Booking',
      badgeClass: 'bg-purple-100 text-purple-800',
      iconClass: 'text-purple-600'
    },
    'event': {
      text: 'Event Contribution',
      badgeClass: 'bg-blue-100 text-blue-800',
      iconClass: 'text-blue-600'
    },
    'annadanam': {
      text: 'Annadanam',
      badgeClass: 'bg-orange-100 text-orange-800',
      iconClass: 'text-orange-600'
    },
    'construction': {
      text: 'Temple Development',
      badgeClass: 'bg-teal-100 text-teal-800',
      iconClass: 'text-teal-600'
    }
  }
  
  return typeMap[type?.toLowerCase()] || typeMap['general']
}

/**
 * Format file size to human readable format
 * @param {number} bytes - File size in bytes
 * @returns {string} Formatted file size
 */
export const formatFileSize = (bytes) => {
  if (!bytes || bytes === 0) return '0 B'
  
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  
  return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${sizes[i]}`
}

/**
 * Format user role for display
 * @param {string} role - User role
 * @returns {object} Role with display text and CSS classes
 */
export const formatUserRole = (role) => {
  const roleMap = {
    'superadmin': {
      text: 'Super Admin',
      badgeClass: 'bg-red-100 text-red-800 border-red-200',
      iconClass: 'text-red-600'
    },
    'tenant': {
      text: 'Temple Admin',
      badgeClass: 'bg-indigo-100 text-indigo-800 border-indigo-200',
      iconClass: 'text-indigo-600'
    },
    'devotee': {
      text: 'Devotee',
      badgeClass: 'bg-green-100 text-green-800 border-green-200',
      iconClass: 'text-green-600'
    },
    'volunteer': {
      text: 'Volunteer',
      badgeClass: 'bg-blue-100 text-blue-800 border-blue-200',
      iconClass: 'text-blue-600'
    }
  }
  
  return roleMap[role?.toLowerCase()] || {
    text: capitalizeWords(role || 'Unknown'),
    badgeClass: 'bg-gray-100 text-gray-800 border-gray-200',
    iconClass: 'text-gray-600'
  }
}

// =============================================================================
// VALIDATION HELPERS
// =============================================================================

/**
 * Check if value exists and is not empty
 * @param {any} value - Value to check
 * @returns {boolean} True if value exists
 */
export const hasValue = (value) => {
  return value !== null && value !== undefined && value !== '' && value !== 0
}

/**
 * Format empty states
 * @param {any} value - Value to check
 * @param {string} fallback - Fallback text
 * @returns {string} Value or fallback
 */
export const formatEmptyState = (value, fallback = 'Not specified') => {
  return hasValue(value) ? value : fallback
}

// =============================================================================
// EXPORT ALL FORMATTERS
// =============================================================================

export default {
  // Date formatters
  formatDate,
  formatEventDate,
  formatRelativeTime,
  
  // Currency & number formatters
  formatCurrency,
  formatNumber,
  formatPercentage,
  
  // Text formatters
  capitalizeWords,
  truncateText,
  formatPhone,
  maskEmail,
  
  // Temple-specific formatters
  formatSevaName,
  formatDevoteeStatus,
  formatTempleStatus,
  formatDonationType,
  formatFileSize,
  formatUserRole,
  
  // Validation helpers
  hasValue,
  formatEmptyState
}