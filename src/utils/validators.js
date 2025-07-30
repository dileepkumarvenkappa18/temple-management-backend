// validators.js - Comprehensive validation utilities for Multi-Tenant Temple SaaS Platform

/**
 * Common validation patterns and utilities
 */

// Regular Expression Patterns
const PATTERNS = {
  email: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
  phone: /^[\+]?[1-9][\d]{0,15}$/,
  indianMobile: /^[6-9]\d{9}$/,
  password: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/,
  alphaNumeric: /^[a-zA-Z0-9\s]+$/,
  alphabetic: /^[a-zA-Z\s]+$/,
  numeric: /^\d+$/,
  url: /^https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)$/,
  pincode: /^[1-9][0-9]{5}$/,
  aadhar: /^\d{12}$/,
  pan: /^[A-Z]{5}[0-9]{4}[A-Z]{1}$/,
  amount: /^\d+(\.\d{1,2})?$/,
  username: /^[a-zA-Z0-9_]{3,20}$/,
  templeName: /^[a-zA-Z0-9\s\-\.,()]+$/,
  address: /^[a-zA-Z0-9\s\-\.,#\/()]+$/
}

// Validation Rules
export const validators = {
  
  // ==================== BASIC VALIDATORS ====================
  
  /**
   * Check if value is required (not empty)
   */
  required: (value, fieldName = 'Field') => {
    if (!value || (typeof value === 'string' && value.trim() === '')) {
      return `${fieldName} is required`
    }
    return true
  },

  /**
   * Validate minimum length
   */
  minLength: (min) => (value, fieldName = 'Field') => {
    if (!value) return true // Let required handle empty values
    if (value.length < min) {
      return `${fieldName} must be at least ${min} characters long`
    }
    return true
  },

  /**
   * Validate maximum length
   */
  maxLength: (max) => (value, fieldName = 'Field') => {
    if (!value) return true
    if (value.length > max) {
      return `${fieldName} must not exceed ${max} characters`
    }
    return true
  },

  /**
   * Validate length range
   */
  lengthRange: (min, max) => (value, fieldName = 'Field') => {
    if (!value) return true
    if (value.length < min || value.length > max) {
      return `${fieldName} must be between ${min} and ${max} characters`
    }
    return true
  },

  // ==================== AUTHENTICATION VALIDATORS ====================

  /**
   * Validate email format
   */
  email: (value, fieldName = 'Email') => {
    if (!value) return true
    if (!PATTERNS.email.test(value)) {
      return `Please enter a valid ${fieldName.toLowerCase()}`
    }
    return true
  },

  /**
   * Validate password strength
   */
  password: (value, fieldName = 'Password') => {
    if (!value) return true
    if (!PATTERNS.password.test(value)) {
      return `${fieldName} must contain at least 8 characters with uppercase, lowercase, number and special character`
    }
    return true
  },

  /**
   * Validate password confirmation
   */
  confirmPassword: (originalPassword) => (value, fieldName = 'Confirm Password') => {
    if (!value) return true
    if (value !== originalPassword) {
      return `${fieldName} does not match`
    }
    return true
  },

  // ==================== CONTACT VALIDATORS ====================

  /**
   * Validate phone number (international format)
   */
  phone: (value, fieldName = 'Phone number') => {
    if (!value) return true
    if (!PATTERNS.phone.test(value)) {
      return `Please enter a valid ${fieldName.toLowerCase()}`
    }
    return true
  },

  /**
   * Validate Indian mobile number
   */
  indianMobile: (value, fieldName = 'Mobile number') => {
    if (!value) return true
    if (!PATTERNS.indianMobile.test(value)) {
      return `Please enter a valid 10-digit Indian ${fieldName.toLowerCase()}`
    }
    return true
  },

  // ==================== PERSONAL INFO VALIDATORS ====================

  /**
   * Validate name (alphabetic with spaces)
   */
  name: (value, fieldName = 'Name') => {
    if (!value) return true
    if (!PATTERNS.alphabetic.test(value)) {
      return `${fieldName} should contain only letters and spaces`
    }
    if (value.trim().length < 2) {
      return `${fieldName} must be at least 2 characters long`
    }
    return true
  },

  /**
   * Validate age
   */
  age: (value, fieldName = 'Age') => {
    if (!value) return true
    const age = parseInt(value)
    if (isNaN(age) || age < 1 || age > 120) {
      return `Please enter a valid ${fieldName.toLowerCase()} between 1 and 120`
    }
    return true
  },

  /**
   * Validate date (not future date for birth dates)
   */
  birthDate: (value, fieldName = 'Birth date') => {
    if (!value) return true
    const date = new Date(value)
    const today = new Date()
    if (date > today) {
      return `${fieldName} cannot be in the future`
    }
    if (date < new Date('1900-01-01')) {
      return `Please enter a valid ${fieldName.toLowerCase()}`
    }
    return true
  },

  // ==================== ADDRESS VALIDATORS ====================

  /**
   * Validate address
   */
  address: (value, fieldName = 'Address') => {
    if (!value) return true
    if (!PATTERNS.address.test(value)) {
      return `${fieldName} contains invalid characters`
    }
    if (value.trim().length < 10) {
      return `${fieldName} should be more descriptive (minimum 10 characters)`
    }
    return true
  },

  /**
   * Validate Indian PIN code
   */
  pincode: (value, fieldName = 'PIN code') => {
    if (!value) return true
    if (!PATTERNS.pincode.test(value)) {
      return `Please enter a valid 6-digit ${fieldName.toLowerCase()}`
    }
    return true
  },

  // ==================== TEMPLE-SPECIFIC VALIDATORS ====================

  /**
   * Validate temple name
   */
  templeName: (value, fieldName = 'Temple name') => {
    if (!value) return true
    if (!PATTERNS.templeName.test(value)) {
      return `${fieldName} contains invalid characters`
    }
    if (value.trim().length < 3) {
      return `${fieldName} must be at least 3 characters long`
    }
    return true
  },

  /**
   * Validate seva name
   */
  sevaName: (value, fieldName = 'Seva name') => {
    if (!value) return true
    if (value.trim().length < 2) {
      return `${fieldName} must be at least 2 characters long`
    }
    if (value.length > 100) {
      return `${fieldName} must not exceed 100 characters`
    }
    return true
  },

  /**
   * Validate donation amount
   */
  donationAmount: (value, fieldName = 'Donation amount') => {
    if (!value) return true
    if (!PATTERNS.amount.test(value)) {
      return `Please enter a valid ${fieldName.toLowerCase()}`
    }
    const amount = parseFloat(value)
    if (amount <= 0) {
      return `${fieldName} must be greater than 0`
    }
    if (amount > 1000000) {
      return `${fieldName} cannot exceed ₹10,00,000`
    }
    return true
  },

  // ==================== DOCUMENT VALIDATORS ====================

  /**
   * Validate Aadhar number
   */
  aadhar: (value, fieldName = 'Aadhar number') => {
    if (!value) return true
    if (!PATTERNS.aadhar.test(value)) {
      return `Please enter a valid 12-digit ${fieldName.toLowerCase()}`
    }
    return true
  },

  /**
   * Validate PAN number
   */
  pan: (value, fieldName = 'PAN number') => {
    if (!value) return true
    if (!PATTERNS.pan.test(value.toUpperCase())) {
      return `Please enter a valid ${fieldName.toLowerCase()} (e.g., ABCDE1234F)`
    }
    return true
  },

  // ==================== FILE VALIDATORS ====================

  /**
   * Validate file upload
   */
  file: (allowedTypes = [], maxSizeMB = 5) => (file, fieldName = 'File') => {
    if (!file) return true
    
    // Check file type
    if (allowedTypes.length > 0) {
      const fileType = file.type.toLowerCase()
      const isValidType = allowedTypes.some(type => 
        fileType.includes(type.toLowerCase())
      )
      if (!isValidType) {
        return `${fieldName} must be of type: ${allowedTypes.join(', ')}`
      }
    }
    
    // Check file size
    const fileSizeMB = file.size / (1024 * 1024)
    if (fileSizeMB > maxSizeMB) {
      return `${fieldName} size must not exceed ${maxSizeMB}MB`
    }
    
    return true
  },

  /**
   * Validate image file
   */
  image: (maxSizeMB = 2) => {
    return validators.file(['jpg', 'jpeg', 'png', 'webp'], maxSizeMB)
  },

  /**
   * Validate document file
   */
  document: (maxSizeMB = 10) => {
    return validators.file(['pdf', 'doc', 'docx', 'jpg', 'jpeg', 'png'], maxSizeMB)
  },

  // ==================== URL VALIDATORS ====================

  /**
   * Validate URL
   */
  url: (value, fieldName = 'URL') => {
    if (!value) return true
    if (!PATTERNS.url.test(value)) {
      return `Please enter a valid ${fieldName.toLowerCase()}`
    }
    return true
  },

  // ==================== SELECTION VALIDATORS ====================

  /**
   * Validate selection from options
   */
  select: (options = []) => (value, fieldName = 'Selection') => {
    if (!value) return true
    if (!options.includes(value)) {
      return `Please select a valid ${fieldName.toLowerCase()}`
    }
    return true
  },

  /**
   * Validate multiple selections
   */
  multiSelect: (options = [], minSelections = 0, maxSelections = null) => (values, fieldName = 'Selections') => {
    if (!values || !Array.isArray(values)) {
      if (minSelections > 0) {
        return `Please make at least ${minSelections} ${fieldName.toLowerCase()}`
      }
      return true
    }
    
    // Check minimum selections
    if (values.length < minSelections) {
      return `Please select at least ${minSelections} ${fieldName.toLowerCase()}`
    }
    
    // Check maximum selections
    if (maxSelections && values.length > maxSelections) {
      return `Please select no more than ${maxSelections} ${fieldName.toLowerCase()}`
    }
    
    // Check if all values are valid options
    const invalidValues = values.filter(value => !options.includes(value))
    if (invalidValues.length > 0) {
      return `Invalid ${fieldName.toLowerCase()}: ${invalidValues.join(', ')}`
    }
    
    return true
  },

  // ==================== DATE/TIME VALIDATORS ====================

  /**
   * Validate date format and range
   */
  date: (minDate = null, maxDate = null) => (value, fieldName = 'Date') => {
    if (!value) return true
    
    const date = new Date(value)
    if (isNaN(date.getTime())) {
      return `Please enter a valid ${fieldName.toLowerCase()}`
    }
    
    if (minDate && date < new Date(minDate)) {
      return `${fieldName} cannot be earlier than ${new Date(minDate).toLocaleDateString()}`
    }
    
    if (maxDate && date > new Date(maxDate)) {
      return `${fieldName} cannot be later than ${new Date(maxDate).toLocaleDateString()}`
    }
    
    return true
  },

  /**
   * Validate future date (for events)
   */
  futureDate: (value, fieldName = 'Date') => {
    if (!value) return true
    const date = new Date(value)
    const today = new Date()
    today.setHours(0, 0, 0, 0)
    
    if (date < today) {
      return `${fieldName} must be today or in the future`
    }
    return true
  },

  /**
   * Validate time format (HH:MM)
   */
  time: (value, fieldName = 'Time') => {
    if (!value) return true
    const timeRegex = /^([01]?[0-9]|2[0-3]):[0-5][0-9]$/
    if (!timeRegex.test(value)) {
      return `Please enter a valid ${fieldName.toLowerCase()} (HH:MM format)`
    }
    return true
  }
}

// ==================== COMPOSITE VALIDATORS ====================

/**
 * Combine multiple validators for a single field
 */
export const combineValidators = (...validatorFunctions) => {
  return (value, fieldName) => {
    for (const validator of validatorFunctions) {
      const result = validator(value, fieldName)
      if (result !== true) {
        return result // Return first error found
      }
    }
    return true
  }
}

/**
 * Validate entire form object
 */
export const validateForm = (formData, validationRules) => {
  const errors = {}
  let isValid = true

  for (const [fieldName, rules] of Object.entries(validationRules)) {
    const value = formData[fieldName]
    
    for (const rule of rules) {
      const result = rule(value, fieldName)
      if (result !== true) {
        errors[fieldName] = result
        isValid = false
        break // Stop at first error for this field
      }
    }
  }

  return { isValid, errors }
}

// ==================== FORM-SPECIFIC VALIDATION SCHEMAS ====================

/**
 * Pre-defined validation schemas for common forms
 */
export const validationSchemas = {
  
  // Registration form validation
  registration: {
    fullName: [
      validators.required,
      validators.name,
      validators.lengthRange(2, 50)
    ],
    email: [
      validators.required,
      validators.email
    ],
    password: [
      validators.required,
      validators.password
    ],
    phone: [
      validators.indianMobile
    ],
    role: [
      validators.required,
      validators.select(['tenant', 'devotee', 'volunteer'])
    ]
  },

  // Login form validation
  login: {
    email: [
      validators.required,
      validators.email
    ],
    password: [
      validators.required
    ]
  },

  // Temple creation validation
  templeCreation: {
    name: [
      validators.required,
      validators.templeName,
      validators.lengthRange(3, 100)
    ],
    address: [
      validators.required,
      validators.address
    ],
    city: [
      validators.required,
      validators.name,
      validators.lengthRange(2, 50)
    ],
    state: [
      validators.required,
      validators.name
    ],
    pincode: [
      validators.required,
      validators.pincode
    ],
    contactEmail: [
      validators.required,
      validators.email
    ],
    contactPhone: [
      validators.required,
      validators.indianMobile
    ]
  },

  // Devotee profile validation
  devoteeProfile: {
    fullName: [
      validators.required,
      validators.name,
      validators.lengthRange(2, 50)
    ],
    email: [
      validators.required,
      validators.email
    ],
    phone: [
      validators.required,
      validators.indianMobile
    ],
    dateOfBirth: [
      validators.birthDate
    ],
    address: [
      validators.address
    ]
  },

  // Seva booking validation
  sevaBooking: {
    sevaId: [
      validators.required
    ],
    devoteeId: [
      validators.required
    ],
    bookingDate: [
      validators.required,
      validators.futureDate
    ],
    specialInstructions: [
      validators.maxLength(500)
    ]
  },

  // Donation validation
  donation: {
    amount: [
      validators.required,
      validators.donationAmount
    ],
    donorName: [
      validators.required,
      validators.name
    ],
    donorPhone: [
      validators.indianMobile
    ],
    donorEmail: [
      validators.email
    ],
    purpose: [
      validators.maxLength(200)
    ]
  },

  // Event creation validation
  eventCreation: {
    title: [
      validators.required,
      validators.lengthRange(3, 100)
    ],
    description: [
      validators.required,
      validators.lengthRange(10, 1000)
    ],
    eventDate: [
      validators.required,
      validators.futureDate
    ],
    startTime: [
      validators.required,
      validators.time
    ],
    endTime: [
      validators.required,
      validators.time
    ],
    venue: [
      validators.required,
      validators.lengthRange(3, 200)
    ]
  }
}

// ==================== HELPER FUNCTIONS ====================

/**
 * Get password strength score (0-4)
 */
export const getPasswordStrength = (password) => {
  if (!password) return 0
  
  let score = 0
  
  // Length check
  if (password.length >= 8) score++
  if (password.length >= 12) score++
  
  // Character variety checks
  if (/[a-z]/.test(password)) score++
  if (/[A-Z]/.test(password)) score++
  if (/\d/.test(password)) score++
  if (/[@$!%*?&]/.test(password)) score++
  
  return Math.min(score, 4)
}

/**
 * Get password strength text and color
 */
export const getPasswordStrengthInfo = (password) => {
  const strength = getPasswordStrength(password)
  
  const strengthInfo = {
    0: { text: 'Very Weak', color: 'text-red-600', bgColor: 'bg-red-100' },
    1: { text: 'Weak', color: 'text-red-500', bgColor: 'bg-red-50' },
    2: { text: 'Fair', color: 'text-yellow-600', bgColor: 'bg-yellow-100' },
    3: { text: 'Good', color: 'text-blue-600', bgColor: 'bg-blue-100' },
    4: { text: 'Strong', color: 'text-green-600', bgColor: 'bg-green-100' }
  }
  
  return strengthInfo[strength]
}

/**
 * Sanitize input to prevent XSS
 */
export const sanitizeInput = (input) => {
  if (typeof input !== 'string') return input
  
  return input
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#x27;')
    .replace(/\//g, '&#x2F;')
}

/**
 * Format validation error for display
 */
export const formatValidationError = (error) => {
  if (typeof error === 'string') return error
  if (typeof error === 'object' && error.message) return error.message
  return 'Invalid input'
}

export default validators