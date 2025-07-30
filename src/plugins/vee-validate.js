import { defineRule, configure } from 'vee-validate'
import { required, email, min, max, confirmed, regex } from '@vee-validate/rules'
import { localize } from '@vee-validate/i18n'

// Define validation rules
defineRule('required', required)
defineRule('email', email)
defineRule('min', min)
defineRule('max', max)
defineRule('confirmed', confirmed)
defineRule('regex', regex)

// Custom rules for temple SaaS
defineRule('phone', (value) => {
  if (!value) return true // Optional field
  const phoneRegex = /^[\+]?[1-9][\d]{0,15}$/
  return phoneRegex.test(value) || 'Please enter a valid phone number'
})

defineRule('password', (value) => {
  if (!value) return true
  
  // Password must be at least 8 characters with uppercase, lowercase, number, and special char
  const hasMinLength = value.length >= 8
  const hasUppercase = /[A-Z]/.test(value)
  const hasLowercase = /[a-z]/.test(value)  
  const hasNumber = /\d/.test(value)
  const hasSpecialChar = /[!@#$%^&*(),.?":{}|<>]/.test(value)
  
  if (!hasMinLength) return 'Password must be at least 8 characters long'
  if (!hasUppercase) return 'Password must contain at least one uppercase letter'
  if (!hasLowercase) return 'Password must contain at least one lowercase letter'
  if (!hasNumber) return 'Password must contain at least one number'
  if (!hasSpecialChar) return 'Password must contain at least one special character'
  
  return true
})

defineRule('gotra', (value) => {
  if (!value) return true // Optional field
  // Basic validation for gotra names (letters and spaces only)
  const gotraRegex = /^[a-zA-Z\s]+$/
  return gotraRegex.test(value) || 'Gotra should contain only letters and spaces'
})

defineRule('nakshatra', (value) => {
  if (!value) return true // Optional field
  const nakshatras = [
    'Ashwini', 'Bharani', 'Krittika', 'Rohini', 'Mrigashira', 'Ardra', 'Punarvasu',
    'Pushya', 'Ashlesha', 'Magha', 'Purva Phalguni', 'Uttara Phalguni', 'Hasta',
    'Chitra', 'Swati', 'Vishakha', 'Anuradha', 'Jyeshtha', 'Mula', 'Purva Ashadha',
    'Uttara Ashadha', 'Shravana', 'Dhanishta', 'Shatabhisha', 'Purva Bhadrapada',
    'Uttara Bhadrapada', 'Revati'
  ]
  return nakshatras.includes(value) || 'Please select a valid Nakshatra'
})

defineRule('rashi', (value) => {
  if (!value) return true // Optional field
  const rashis = [
    'Mesha', 'Vrishabha', 'Mithuna', 'Karka', 'Simha', 'Kanya',
    'Tula', 'Vrischika', 'Dhanus', 'Makara', 'Kumbha', 'Meena'
  ]
  return rashis.includes(value) || 'Please select a valid Rashi'
})

defineRule('vedashaka', (value) => {
  if (!value) return true // Optional field
  const vedashakas = [
    'Rigveda - Shakala', 'Rigveda - Baskala', 'Rigveda - Ashvalayana',
    'Yajurveda - Madhyandina', 'Yajurveda - Kanva', 'Yajurveda - Taittiriya',
    'Yajurveda - Maitrayani', 'Samaveda - Kauthuma', 'Samaveda - Ranayaniya',
    'Samaveda - Jaiminiya', 'Atharvaveda - Shaunaka', 'Atharvaveda - Paippalada'
  ]
  return vedashakas.includes(value) || 'Please select a valid Veda Shaka'
})

defineRule('pincode', (value) => {
  if (!value) return true // Optional field
  const pincodeRegex = /^[1-9][0-9]{5}$/
  return pincodeRegex.test(value) || 'Please enter a valid 6-digit pincode'
})

defineRule('templeName', (value) => {
  if (!value) return 'Temple name is required'
  if (value.length < 3) return 'Temple name must be at least 3 characters'
  if (value.length > 100) return 'Temple name must not exceed 100 characters'
  
  // Should contain letters, numbers, spaces, and common temple-related characters
  const templeNameRegex = /^[a-zA-Z0-9\s\-\(\)\&\.]+$/
  return templeNameRegex.test(value) || 'Temple name contains invalid characters'
})

defineRule('amount', (value) => {
  if (!value) return true // Optional field for some cases
  const amount = parseFloat(value)
  if (isNaN(amount)) return 'Please enter a valid amount'
  if (amount <= 0) return 'Amount must be greater than 0'
  if (amount > 1000000) return 'Amount cannot exceed ₹10,00,000'
  
  // Check for valid decimal places (max 2)
  const decimalPlaces = (value.toString().split('.')[1] || '').length
  if (decimalPlaces > 2) return 'Amount can have maximum 2 decimal places'
  
  return true
})

defineRule('fileSize', (file, [maxSize]) => {
  if (!file || !file.size) return true
  const sizeInMB = file.size / 1024 / 1024
  return sizeInMB <= maxSize || `File size should not exceed ${maxSize}MB`
})

defineRule('fileType', (file, allowedTypes) => {
  if (!file || !file.type) return true
  return allowedTypes.includes(file.type) || `File type not allowed. Allowed types: ${allowedTypes.join(', ')}`
})

defineRule('dateRange', (value, [startDate, endDate]) => {
  if (!value) return true
  const date = new Date(value)
  const start = new Date(startDate)
  const end = new Date(endDate)
  
  return (date >= start && date <= end) || `Date must be between ${startDate} and ${endDate}`
})

defineRule('futureDate', (value) => {
  if (!value) return true
  const selectedDate = new Date(value)
  const today = new Date()
  today.setHours(0, 0, 0, 0) // Reset time to start of day
  
  return selectedDate >= today || 'Date cannot be in the past'
})

// Configure validation behavior
configure({
  generateMessage: localize('en', {
    messages: {
      required: 'This field is required',
      email: 'Please enter a valid email address',
      min: 'This field must be at least {length} characters',
      max: 'This field must not exceed {length} characters',
      confirmed: 'Passwords do not match',
      regex: 'This field format is invalid'
    },
    names: {
      email: 'Email',
      password: 'Password',
      password_confirmation: 'Password Confirmation',
      full_name: 'Full Name',
      phone: 'Phone Number',
      temple_name: 'Temple Name',
      address: 'Address',
      city: 'City',
      state: 'State',
      pincode: 'Pincode',
      gotra: 'Gotra',
      nakshatra: 'Nakshatra',
      rashi: 'Rashi',
      vedashaka: 'Veda Shaka',
      amount: 'Amount',
      date_of_birth: 'Date of Birth',
      event_date: 'Event Date',
      seva_date: 'Seva Date'
    }
  }),
  
  // Validate on input change and blur
  validateOnInput: true,
  validateOnBlur: true,
  validateOnChange: true,
  validateOnModelUpdate: true
})

// Export validation schemas for common forms
export const validationSchemas = {
  // Authentication forms
  login: {
    email: 'required|email',
    password: 'required|min:6'
  },
  
  register: {
    full_name: 'required|min:2|max:50',
    email: 'required|email',
    password: 'required|password',
    password_confirmation: 'required|confirmed:@password',
    phone: 'phone',
    role: 'required'
  },
  
  forgotPassword: {
    email: 'required|email'
  },
  
  resetPassword: {
    password: 'required|password',
    password_confirmation: 'required|confirmed:@password'
  },
  
  // Temple creation form
  createTemple: {
    temple_name: 'required|templeName',
    address: 'required|min:10|max:200',
    city: 'required|min:2|max:50',
    state: 'required|min:2|max:50',
    pincode: 'required|pincode',
    contact_person: 'required|min:2|max:50',
    contact_phone: 'required|phone',
    contact_email: 'required|email'
  },
  
  // Devotee profile forms
  personalDetails: {
    full_name: 'required|min:2|max:50',
    date_of_birth: 'required',
    gender: 'required',
    address: 'required|min:10|max:200',
    city: 'required|min:2|max:50',
    state: 'required|min:2|max:50',
    pincode: 'required|pincode',
    phone: 'phone'
  },
  
  spiritualInfo: {
    gotra: 'gotra',
    nakshatra: 'nakshatra',
    rashi: 'rashi',
    vedashaka: 'vedashaka'
  },
  
  // Seva booking form
  sevaBooking: {
    seva_type: 'required',
    seva_date: 'required|futureDate',
    devotee_name: 'required|min:2|max:50',
    phone: 'phone',
    special_requirements: 'max:500'
  },
  
  // Donation form
  donation: {
    amount: 'required|amount',
    donation_type: 'required',
    donor_name: 'required|min:2|max:50',
    phone: 'phone',
    email: 'email',
    purpose: 'max:200'
  },
  
  // Event creation form
  createEvent: {
    title: 'required|min:3|max:100',
    description: 'required|min:10|max:1000',
    event_date: 'required|futureDate',
    event_time: 'required',
    venue: 'required|min:3|max:100',
    max_attendees: 'min:1'
  },
  
  // Communication form
  bulkMessage: {
    subject: 'required|min:3|max:100',
    message: 'required|min:10|max:1000',
    recipient_type: 'required'
  }
}

export default {
  install(app) {
    // Plugin installation is handled by importing this file in main.js
  }
}