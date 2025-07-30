// src/utils/constants.js
// Core constants for the Temple SaaS platform

export const USER_ROLES = {
  SUPER_ADMIN: 'super_admin',
  TENANT: 'tenant', // Temple Admin
  DEVOTEE: 'devotee', // End User
  VOLUNTEER: 'volunteer' // End User
}

export const USER_STATUS = {
  PENDING: 'pending',
  APPROVED: 'approved',
  REJECTED: 'rejected',
  ACTIVE: 'active',
  INACTIVE: 'inactive'
}

export const TEMPLE_STATUS = {
  PENDING: 'pending',
  APPROVED: 'approved',
  REJECTED: 'rejected'
}

export const SEVA_STATUS = {
  PENDING: 'pending',
  APPROVED: 'approved',
  REJECTED: 'rejected',
  COMPLETED: 'completed'
}

export const DONATION_STATUS = {
  PENDING: 'pending',
  COMPLETED: 'completed',
  FAILED: 'failed',
  REFUNDED: 'refunded'
}

export const EVENT_STATUS = {
  DRAFT: 'draft',
  PUBLISHED: 'published',
  CANCELLED: 'cancelled',
  COMPLETED: 'completed'
}

export const RSVP_STATUS = {
  ATTENDING: 'attending',
  NOT_ATTENDING: 'not_attending',
  MAYBE: 'maybe'
}

export const NOTIFICATION_TYPES = {
  SUCCESS: 'success',
  ERROR: 'error',
  WARNING: 'warning',
  INFO: 'info'
}

// Route paths for different dashboards
export const DASHBOARD_ROUTES = {
  [USER_ROLES.TENANT]: '/tenant/dashboard',
  [USER_ROLES.DEVOTEE]: '/entity/:id/devotee/dashboard',
  [USER_ROLES.VOLUNTEER]: '/entity/:id/volunteer/dashboard',
  [USER_ROLES.SUPER_ADMIN]: '/admin/dashboard'
}

// Profile completion steps for devotees
export const PROFILE_STEPS = {
  PERSONAL_DETAILS: {
    key: 'personal_details',
    label: 'Personal Details',
    weight: 20,
    fields: ['name', 'dob', 'gender', 'address']
  },
  SPIRITUAL_INFO: {
    key: 'spiritual_info', 
    label: 'Spiritual Info',
    weight: 15,
    fields: ['gotra', 'nakshatra', 'rashi', 'lagna', 'veda_shaka']
  },
  FAMILY_LINEAGE: {
    key: 'family_lineage',
    label: 'Family Lineage', 
    weight: 15,
    fields: ['parents', 'grandparents', 'family_tree']
  },
  SEVA_PREFERENCES: {
    key: 'seva_preferences',
    label: 'Seva Preferences',
    weight: 15,
    fields: ['favorite_sevas', 'donation_interests']
  },
  FAMILY_MEMBERS: {
    key: 'family_members',
    label: 'Family Members',
    weight: 20,
    fields: ['spouse', 'children', 'emergency_contacts']
  },
  NOTES_ATTACHMENTS: {
    key: 'notes_attachments',
    label: 'Notes & Attachments',
    weight: 15,
    fields: ['photos', 'health_info', 'sankalpa_notes']
  }
}

// Permissions matrix
export const PERMISSIONS = {
  // Temple Management
  CREATE_TEMPLE: [USER_ROLES.TENANT],
  MANAGE_TEMPLE: [USER_ROLES.TENANT],
  APPROVE_TEMPLE: [USER_ROLES.SUPER_ADMIN],
  
  // Devotee Management  
  VIEW_DEVOTEES: [USER_ROLES.TENANT],
  MANAGE_DEVOTEES: [USER_ROLES.TENANT],
  EXPORT_DEVOTEES: [USER_ROLES.TENANT],
  
  // Seva Management
  CREATE_SEVA: [USER_ROLES.TENANT],
  MANAGE_SEVA: [USER_ROLES.TENANT],
  BOOK_SEVA: [USER_ROLES.DEVOTEE],
  VOLUNTEER_SEVA: [USER_ROLES.VOLUNTEER],
  
  // Donation Management
  VIEW_DONATIONS: [USER_ROLES.TENANT],
  MAKE_DONATION: [USER_ROLES.DEVOTEE],
  GENERATE_RECEIPTS: [USER_ROLES.TENANT],
  
  // Event Management
  CREATE_EVENT: [USER_ROLES.TENANT],
  MANAGE_EVENT: [USER_ROLES.TENANT],
  RSVP_EVENT: [USER_ROLES.DEVOTEE, USER_ROLES.VOLUNTEER],
  
  // Communication
  SEND_BULK_MESSAGES: [USER_ROLES.TENANT],
  VIEW_MESSAGE_HISTORY: [USER_ROLES.TENANT]
}

// Form validation patterns
export const VALIDATION_PATTERNS = {
  EMAIL: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
  PHONE: /^[\+]?[1-9][\d]{0,15}$/,
  STRONG_PASSWORD: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/,
  NAME: /^[a-zA-Z\s]{2,50}$/,
  PINCODE: /^[1-9][0-9]{5}$/
}

// API endpoints base paths
export const API_ENDPOINTS = {
  AUTH: '/auth',
  TEMPLES: '/temples',
  DEVOTEES: '/devotees', 
  SEVAS: '/sevas',
  DONATIONS: '/donations',
  EVENTS: '/events',
  VOLUNTEERS: '/volunteers',
  COMMUNICATION: '/communication',
  UPLOADS: '/uploads'
}

// Toast/Notification messages
export const MESSAGES = {
  REGISTRATION_SUCCESS: 'Thanks for registering!',
  LOGIN_SUCCESS: 'Welcome back!',
  TEMPLE_CREATED: 'Temple created successfully! Awaiting approval.',
  TEMPLE_APPROVED: 'Temple approved – You can now manage it.',
  TEMPLE_REJECTED: 'Temple rejected – See admin notes and re-upload documents.',
  PROFILE_UPDATED: 'Profile updated successfully!',
  SEVA_BOOKED: 'Seva booked successfully!',
  DONATION_SUCCESS: 'Thank you for your donation!',
  EVENT_RSVP: 'RSVP updated successfully!',
  MESSAGE_SENT: 'Message sent successfully!',
  TEMPLE_JOINED: "You've successfully joined {templeName}!",
  
  // Error messages
  INVALID_CREDENTIALS: 'Invalid email or password.',
  NETWORK_ERROR: 'Network error. Please try again.',
  UNAUTHORIZED: 'You are not authorized to perform this action.',
  TEMPLE_NOT_FOUND: 'Temple not found.',
  PROFILE_INCOMPLETE: 'Please complete your profile to continue.'
}

// File upload constraints
export const UPLOAD_CONSTRAINTS = {
  MAX_FILE_SIZE: 5 * 1024 * 1024, // 5MB
  ALLOWED_IMAGE_TYPES: ['image/jpeg', 'image/png', 'image/webp'],
  ALLOWED_DOCUMENT_TYPES: ['application/pdf', 'image/jpeg', 'image/png'],
  MAX_FILES_PER_UPLOAD: 5
}

// Dashboard widget types
export const WIDGET_TYPES = {
  STATS: 'stats',
  CHART: 'chart', 
  LIST: 'list',
  CALENDAR: 'calendar',
  ACTIVITY: 'activity'
}

// Chart colors for consistency
export const CHART_COLORS = {
  PRIMARY: '#4f46e5', // indigo-600
  SECONDARY: '#6366f1', // indigo-500  
  SUCCESS: '#10b981', // emerald-500
  WARNING: '#f59e0b', // amber-500
  ERROR: '#ef4444', // red-500
  INFO: '#06b6d4' // cyan-500
}