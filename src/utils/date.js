// src/utils/date.js
// Date utility functions for Temple SaaS Platform

/**
 * Format date to a readable string
 * @param {Date|string} date - Date to format
 * @param {string} format - Format type ('short', 'long', 'medium', 'time', 'datetime')
 * @returns {string} Formatted date string
 */
export const formatDate = (date, format = 'medium') => {
  if (!date) return ''
  
  const dateObj = new Date(date)
  if (isNaN(dateObj.getTime())) return ''

  const options = {
    short: { 
      year: 'numeric', 
      month: 'short', 
      day: 'numeric' 
    },
    medium: { 
      year: 'numeric', 
      month: 'long', 
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

  return dateObj.toLocaleDateString('en-IN', options[format] || options.medium)
}

/**
 * Get relative time (e.g., "2 hours ago", "in 3 days")
 * @param {Date|string} date - Date to compare
 * @returns {string} Relative time string
 */
export const getRelativeTime = (date) => {
  if (!date) return ''
  
  const dateObj = new Date(date)
  if (isNaN(dateObj.getTime())) return ''

  const now = new Date()
  const diffInMs = dateObj.getTime() - now.getTime()
  const diffInSeconds = Math.floor(diffInMs / 1000)
  const diffInMinutes = Math.floor(diffInSeconds / 60)
  const diffInHours = Math.floor(diffInMinutes / 60)
  const diffInDays = Math.floor(diffInHours / 24)

  const rtf = new Intl.RelativeTimeFormat('en', { numeric: 'auto' })

  if (Math.abs(diffInDays) >= 1) {
    return rtf.format(diffInDays, 'day')
  } else if (Math.abs(diffInHours) >= 1) {
    return rtf.format(diffInHours, 'hour')
  } else if (Math.abs(diffInMinutes) >= 1) {
    return rtf.format(diffInMinutes, 'minute')
  } else {
    return rtf.format(diffInSeconds, 'second')
  }
}

/**
 * Check if date is today
 * @param {Date|string} date - Date to check
 * @returns {boolean} True if date is today
 */
export const isToday = (date) => {
  if (!date) return false
  
  const dateObj = new Date(date)
  const today = new Date()
  
  return dateObj.toDateString() === today.toDateString()
}

/**
 * Check if date is within current week
 * @param {Date|string} date - Date to check
 * @returns {boolean} True if date is this week
 */
export const isThisWeek = (date) => {
  if (!date) return false
  
  const dateObj = new Date(date)
  const today = new Date()
  const startOfWeek = new Date(today)
  const endOfWeek = new Date(today)
  
  // Get start of week (Sunday)
  startOfWeek.setDate(today.getDate() - today.getDay())
  startOfWeek.setHours(0, 0, 0, 0)
  
  // Get end of week (Saturday)
  endOfWeek.setDate(today.getDate() + (6 - today.getDay()))
  endOfWeek.setHours(23, 59, 59, 999)
  
  return dateObj >= startOfWeek && dateObj <= endOfWeek
}

/**
 * Check if date is within current month
 * @param {Date|string} date - Date to check
 * @returns {boolean} True if date is this month
 */
export const isThisMonth = (date) => {
  if (!date) return false
  
  const dateObj = new Date(date)
  const today = new Date()
  
  return dateObj.getMonth() === today.getMonth() && 
         dateObj.getFullYear() === today.getFullYear()
}

/**
 * Get start and end of current month
 * @returns {Object} Object with startOfMonth and endOfMonth dates
 */
export const getCurrentMonthRange = () => {
  const now = new Date()
  const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1)
  const endOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0)
  
  return { startOfMonth, endOfMonth }
}

/**
 * Get age from date of birth
 * @param {Date|string} birthDate - Date of birth
 * @returns {number} Age in years
 */
export const getAge = (birthDate) => {
  if (!birthDate) return 0
  
  const birth = new Date(birthDate)
  const today = new Date()
  
  let age = today.getFullYear() - birth.getFullYear()
  const monthDiff = today.getMonth() - birth.getMonth()
  
  if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birth.getDate())) {
    age--
  }
  
  return age
}

/**
 * Format date for input fields (YYYY-MM-DD)
 * @param {Date|string} date - Date to format
 * @returns {string} Date in YYYY-MM-DD format
 */
export const formatForInput = (date) => {
  if (!date) return ''
  
  const dateObj = new Date(date)
  if (isNaN(dateObj.getTime())) return ''
  
  return dateObj.toISOString().split('T')[0]
}

/**
 * Format time for input fields (HH:MM)
 * @param {Date|string} date - Date to format
 * @returns {string} Time in HH:MM format
 */
export const formatTimeForInput = (date) => {
  if (!date) return ''
  
  const dateObj = new Date(date)
  if (isNaN(dateObj.getTime())) return ''
  
  return dateObj.toTimeString().slice(0, 5)
}

/**
 * Add days to a date
 * @param {Date|string} date - Starting date
 * @param {number} days - Number of days to add
 * @returns {Date} New date with days added
 */
export const addDays = (date, days) => {
  const dateObj = new Date(date)
  dateObj.setDate(dateObj.getDate() + days)
  return dateObj
}

/**
 * Subtract days from a date
 * @param {Date|string} date - Starting date
 * @param {number} days - Number of days to subtract
 * @returns {Date} New date with days subtracted
 */
export const subtractDays = (date, days) => {
  const dateObj = new Date(date)
  dateObj.setDate(dateObj.getDate() - days)
  return dateObj
}

/**
 * Get difference between two dates in days
 * @param {Date|string} date1 - First date
 * @param {Date|string} date2 - Second date
 * @returns {number} Difference in days
 */
export const getDaysDifference = (date1, date2) => {
  const d1 = new Date(date1)
  const d2 = new Date(date2)
  const timeDiff = d2.getTime() - d1.getTime()
  return Math.ceil(timeDiff / (1000 * 3600 * 24))
}

/**
 * Check if a date is in the past
 * @param {Date|string} date - Date to check
 * @returns {boolean} True if date is in the past
 */
export const isPast = (date) => {
  if (!date) return false
  return new Date(date) < new Date()
}

/**
 * Check if a date is in the future
 * @param {Date|string} date - Date to check
 * @returns {boolean} True if date is in the future
 */
export const isFuture = (date) => {
  if (!date) return false
  return new Date(date) > new Date()
}

/**
 * Get Indian festival/special dates
 * @param {number} year - Year to get festivals for
 * @returns {Array} Array of festival objects with date and name
 */
export const getIndianFestivals = (year = new Date().getFullYear()) => {
  // This is a simplified list - in production, you'd integrate with a proper calendar API
  return [
    { date: `${year}-01-26`, name: 'Republic Day' },
    { date: `${year}-08-15`, name: 'Independence Day' },
    { date: `${year}-10-02`, name: 'Gandhi Jayanti' },
    // Add more festivals based on lunar calendar calculations
  ]
}

/**
 * Format date range
 * @param {Date|string} startDate - Start date
 * @param {Date|string} endDate - End date
 * @param {string} format - Format type
 * @returns {string} Formatted date range
 */
export const formatDateRange = (startDate, endDate, format = 'medium') => {
  if (!startDate && !endDate) return ''
  
  if (!endDate) return formatDate(startDate, format)
  if (!startDate) return `Until ${formatDate(endDate, format)}`
  
  const start = formatDate(startDate, format)
  const end = formatDate(endDate, format)
  
  return `${start} - ${end}`
}

/**
 * Get weekday name
 * @param {Date|string} date - Date to get weekday for
 * @returns {string} Weekday name
 */
export const getWeekday = (date) => {
  if (!date) return ''
  
  const dateObj = new Date(date)
  return dateObj.toLocaleDateString('en-IN', { weekday: 'long' })
}

/**
 * Get month name
 * @param {Date|string} date - Date to get month for
 * @returns {string} Month name
 */
export const getMonth = (date) => {
  if (!date) return ''
  
  const dateObj = new Date(date)
  return dateObj.toLocaleDateString('en-IN', { month: 'long' })
}

/**
 * Check if date is weekend (Saturday or Sunday)
 * @param {Date|string} date - Date to check
 * @returns {boolean} True if date is weekend
 */
export const isWeekend = (date) => {
  if (!date) return false
  
  const dateObj = new Date(date)
  const day = dateObj.getDay()
  return day === 0 || day === 6 // Sunday = 0, Saturday = 6
}

/**
 * Get next occurrence of a specific weekday
 * @param {number} targetDay - Target weekday (0 = Sunday, 1 = Monday, etc.)
 * @param {Date|string} fromDate - Starting date (default: today)
 * @returns {Date} Next occurrence of the weekday
 */
export const getNextWeekday = (targetDay, fromDate = new Date()) => {
  const date = new Date(fromDate)
  const currentDay = date.getDay()
  const daysUntilTarget = (targetDay - currentDay + 7) % 7
  
  if (daysUntilTarget === 0) {
    // If today is the target day, get next week's occurrence
    date.setDate(date.getDate() + 7)
  } else {
    date.setDate(date.getDate() + daysUntilTarget)
  }
  
  return date
}